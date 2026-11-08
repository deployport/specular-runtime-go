package server

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"go.deployport.com/specular-runtime/client"
)

type serviceOperationStreamServeContext struct {
	*serviceOperationContext
	flusher http.Flusher
	writer  *streamPartWriter
	tpChan  <-chan client.StreamEvent[client.Struct]
}

// streamPartWriter writes multipart/mixed with the same bytes as
// multipart.Writer, but it writes the delimiter that ends a part
// right after the part's content instead of when the next part starts.
// A reader can only finish a part once it sees that delimiter, so
// writing it late held every part back until the next one began.
type streamPartWriter struct {
	w        io.Writer
	boundary string
	// ended is true once the close delimiter is written
	ended bool
}

func newStreamPartWriter(w io.Writer) *streamPartWriter {
	return &streamPartWriter{
		w:        w,
		boundary: multipart.NewWriter(io.Discard).Boundary(),
	}
}

// open writes the first delimiter, after which the next part's headers follow
func (sw *streamPartWriter) open() error {
	_, err := fmt.Fprintf(sw.w, "--%s\r\n", sw.boundary)
	return err
}

// writePart writes one part and the delimiter that ends it
func (sw *streamPartWriter) writePart(contentType string, content func(io.Writer) error) error {
	if _, err := fmt.Fprintf(sw.w, "Content-Type: %s\r\n\r\n", contentType); err != nil {
		return err
	}
	if content != nil {
		if err := content(sw.w); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(sw.w, "\r\n--%s\r\n", sw.boundary)
	return err
}

// close ends the body. The writer always sits where a part's headers
// begin, so it writes the end marker (a part with an empty Content-Type,
// which clients read as the end of the stream) and the close delimiter.
func (sw *streamPartWriter) close() error {
	if sw.ended {
		return nil
	}
	sw.ended = true
	_, err := fmt.Fprintf(sw.w, "Content-Type: \r\n\r\n\r\n--%s--\r\n", sw.boundary)
	return err
}

func serveServiceOperationStream(
	octx *serviceOperationContext,
	w http.ResponseWriter,
) {
	handler := octx.streamHandler
	if handler == nil {
		_ = octx.opx.BuildFinalResult(nil, fmt.Errorf("stream handler unset, service unable to serve requests")).WriteResponse(w)
		return
	}
	tpChan, err := handler.HandleStream(octx.Context(), octx.opx)
	if err != nil {
		_ = octx.opx.BuildFinalResult(nil, err).WriteResponse(w)
		return
	}

	flusher, _ := w.(http.Flusher)
	soctx := &serviceOperationStreamServeContext{
		serviceOperationContext: octx,
		flusher:                 flusher,
		tpChan:                  tpChan,
	}
	soctx.prepareMultipart(w)
	soctx.serveStream()
}

func (octx *serviceOperationStreamServeContext) prepareMultipart(w http.ResponseWriter) {
	// Server.WriteTimeout is set once when the headers are read and never extended
	// Clearing it is best effort and a wrapped ResponseWriter may not support it
	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	octx.writer = newStreamPartWriter(w)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "multipart/mixed; boundary="+octx.writer.boundary)
	w.WriteHeader(http.StatusOK)
	_ = octx.writer.open()
}

func (octx *serviceOperationStreamServeContext) flush() {
	if octx.flusher != nil {
		octx.flusher.Flush()
	}
}

// streamResult receives a HTTPResult and streams it to the client using a HTTP multipart response
// result could be nil to indicate the end of the stream
func (octx *serviceOperationStreamServeContext) streamResult(result *client.HTTPResult) error {
	if result == nil {
		// The end marker: an empty Content-Type, then the close delimiter
		err := octx.writer.close()
		octx.flush()
		if err != nil {
			return fmt.Errorf("failed to end stream: %w", err)
		}
		return nil
	}
	if err := octx.writer.writePart(result.MimeType(), result.WriteContent); err != nil {
		return fmt.Errorf("failed to write part: %w", err)
	}
	octx.flush()
	return nil
}

func (octx *serviceOperationStreamServeContext) sendHeartbeatAndFlush() error {
	return octx.streamResult(client.HTTPResultForHeartbeat())
}

// serveStream serves an streamed operation
func (octx *serviceOperationStreamServeContext) serveStream() {
	defer octx.writer.close()

	// sendStreamedType := func(output Struct) error {
	// 	part, err := writer.CreatePart(textproto.MIMEHeader{
	// 		"Content-Type": {"application/json"},
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create part, %w", err)
	// 	}
	// 	if output != nil {
	// 		res, err := ContentFromStruct(output)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to create content from streamed outout struct: %w", err)
	// 		}
	// 		buf := bytes.NewBuffer(nil)
	// 		enc := json.NewEncoder(buf)
	// 		if err := enc.Encode(res); err != nil {
	// 			return fmt.Errorf("failed to encode streamed outout struct: %w", err)
	// 		}
	// 		// log.Printf("writing part: %s", buf.String())
	// 		if _, err := part.Write(buf.Bytes()); err != nil {
	// 			return fmt.Errorf("failed to write part: %w", err)
	// 		}
	// 		if err := sendHeartbeatAndFlush(); err != nil {
	// 			return fmt.Errorf("failed to send heartbeat after result: %w", err)
	// 		}
	// 	} else {
	// 		log.Printf("writing part nil")
	// 		if _, err := part.Write(nil); err != nil {
	// 			return fmt.Errorf("failed to write part: %w", err)
	// 		}
	// 		octx.flush()
	// 	}
	// 	return nil
	// }
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()
	for {
		select {
		case <-octx.Context().Done():
			return
		case <-heartbeatTicker.C:
			if err := octx.sendHeartbeatAndFlush(); err != nil {
				return
			}
		case ev, ok := <-octx.tpChan:
			heartbeatTicker.Reset(heartbeatInterval)
			if !ok {
				// channel closed, end of stream
				_ = octx.streamResult(nil)
				return
			}
			var outputResult *client.HTTPResult
			if o := ev.Output; o != nil {
				outputResult = client.HTTPResultForStruct(o)
			}
			if err := octx.streamResult(octx.opx.BuildFinalResult(outputResult, ev.Err)); err != nil {
				return
			}
		}
	}
}
