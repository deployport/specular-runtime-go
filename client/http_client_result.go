package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClientResult is the result of a HTTP client
type HTTPClientResult struct {
	Single Struct
	Stream *clientStream
}

// ParseHTTPClientResult parses a HTTP client result
// may return UnexpectedClientResultError if the client doesn't understand the server result
// may return io.EOF if the content type is empty (friendly end of stream)
// note that HTTPResult can also have a Struct representing the error of the spec
func ParseHTTPClientResult(structFinder *StructDefinitionFinder, header http.Header, body io.ReadCloser) (*HTTPClientResult, error) {
	contentTypeHeader := header.Get("Content-Type")
	if contentTypeHeader == "" {
		return nil, io.EOF
	}
	contentMime, err := NewMIME(contentTypeHeader)
	if err != nil {
		return nil, err
	}
	streamHeader, err := contentMime.StreamHeader()
	if err != nil {
		return nil, fmt.Errorf("invalid stream header, %w", err)
	}
	if streamHeader != nil {
		stream := newClientStream(streamHeader.Boundary, structFinder, body)
		return &HTTPClientResult{
			Stream: stream,
		}, nil
	}
	structPath, err := contentMime.StructPath()
	if err != nil {
		return nil, err
	}
	if structPath != nil {
		sd, err := structFinder.Find(*structPath)
		if err != nil {
			return nil, err
		}
		dec := json.NewDecoder(body)
		sst := sd.TypeBuilder()()
		if err := dec.Decode(&sst); err != nil {
			return nil, err
		}
		return &HTTPClientResult{
			Single: sst,
		}, nil
	}
	var partialResponseBody bytes.Buffer
	if _, err := partialResponseBody.ReadFrom(io.LimitReader(body, 100)); err != nil {
		return nil, fmt.Errorf("failed to read partial response body, %w", err)
	}
	return nil, NewUnexpectedClientResultError(contentTypeHeader, partialResponseBody.String())
}

// UnexpectedClientResultError is an error for unexpected client result
type UnexpectedClientResultError struct {
	ContentType string
	PartialBody string
	Message     string
}

// Error returns the error message
func (e *UnexpectedClientResultError) Error() string {
	return e.Message
}

// NewUnexpectedClientResultError creates a new UnexpectedClientResultError
func NewUnexpectedClientResultError(contentType, partialBody string) *UnexpectedClientResultError {
	message := fmt.Sprintf("unexpected server response, %s, body %s", contentType, partialBody)
	return &UnexpectedClientResultError{
		ContentType: contentType,
		PartialBody: partialBody,
		Message:     message,
	}
}
