package server

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"go.deployport.com/specular-runtime/client"
)

type serviceOperationStreamServeContext struct {
	*serviceOperationContext
	flusher http.Flusher
	writer  *multipart.Writer
	tpChan  <-chan client.StreamEvent[client.Struct]
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
	octx.writer = multipart.NewWriter(w)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Connection", "Keep-Alive")
	contentType := octx.writer.FormDataContentType()
	contentType = strings.Replace(contentType, "multipart/form-data", "multipart/mixed", 1)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
}

func (octx *serviceOperationStreamServeContext) flush() {
	if octx.flusher != nil {
		octx.flusher.Flush()
	}
}

// streamResult receives a HTTPResult and streams it to the client using a HTTP multipart response
// result could be nil to indicate the end of the stream
func (octx *serviceOperationStreamServeContext) streamResult(result *client.HTTPResult) error {
	mime := ""
	if result != nil {
		mime = result.MimeType().String()
	}
	part, err := octx.writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {mime},
	})
	if err != nil {
		return fmt.Errorf("failed to create part, %w", err)
	}
	if result != nil {
		if err := result.WriteContent(part); err != nil {
			return fmt.Errorf("failed to write part: %w", err)
		}
	}
	octx.flush()
	return nil
}

func (octx *serviceOperationStreamServeContext) sendHeartbeatAndFlush() error {
	return octx.streamResult(client.HTTPResultForHeartbeat())
}

// serveStream serves an streamed operation
func (octx *serviceOperationStreamServeContext) serveStream() {
	defer octx.writer.Close()

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
