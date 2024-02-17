package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPResultMimeType is the mime type of the result
type HTTPResultMimeType string

const (
	// HTTPResultMimeTypeStruct is the mime type for a struct
	HTTPResultMimeTypeStruct HTTPResultMimeType = "specular/struct"
	// HTTPResultMimeTypeError is the mime type for an error
	HTTPResultMimeTypeError HTTPResultMimeType = "specular/error"
	// HTTPResultMimeTypeHeartbeat is the mime type for a heartbeat
	HTTPResultMimeTypeHeartbeat HTTPResultMimeType = "specular/heartbeat"
)

// String returns the string representation of the mime type
func (m HTTPResultMimeType) String() string {
	return string(m)
}

// HTTPResult is the result to send over a HTTP Response
type HTTPResult struct {
	// Struct not nil indicates the result is a struct
	Struct Struct
	// Err not nil indicates the result is an error, the error is JSON serializable
	Err       *HTTPError
	Heartbeat bool
}

// HTTPResultForStruct returns a HTTPResult for a struct
func HTTPResultForStruct(s Struct) *HTTPResult {
	return &HTTPResult{
		Struct: s,
	}
}

// HTTPResultForError returns a HTTPResult for an error
func HTTPResultForError(err HTTPError) *HTTPResult {
	return &HTTPResult{
		Err: &err,
	}
}

// HTTPResultForHeartbeat returns a HTTPResult for a heartbeat
func HTTPResultForHeartbeat() *HTTPResult {
	return &HTTPResult{
		Heartbeat: true,
	}
}

// MimeType returns the mime type of the result based on the field set in the following order of priority: struct, err, heartbeat
func (r *HTTPResult) MimeType() HTTPResultMimeType {
	if r.Struct != nil {
		return HTTPResultMimeTypeStruct
	}
	if r.Err != nil {
		return HTTPResultMimeTypeError
	}
	if r.Heartbeat {
		return HTTPResultMimeTypeHeartbeat
	}
	return ""
}

// HTTPStatusCode returns the HTTP status code of the result
func (r *HTTPResult) HTTPStatusCode() int {
	if r.Err != nil {
		return r.Err.HTTPStatusCode
	}
	return http.StatusOK
}

// WriteResponse writes the result to the HTTP response
func (r *HTTPResult) WriteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", r.MimeType().String())
	w.WriteHeader(r.HTTPStatusCode())
	return r.WriteContent(w)
}

// WriteContent writes the content to the io.Writer
func (r *HTTPResult) WriteContent(w io.Writer) error {
	if r.Err != nil {
		return r.Err.writeContent(w)
	}
	if r.Struct != nil {
		return r.writeStructContent(w)
	}
	if r.Heartbeat {
		return r.writeHeartbeatContent(w)
	}
	return fmt.Errorf("unwritable http result")
}

func (r *HTTPResult) writeStructContent(w io.Writer) error {
	content := NewContent()
	err := r.Struct.Dehydrate(NewDehydrationContext(content))
	if err != nil {
		return fmt.Errorf("failed to create content from struct: %w", err)
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(content); err != nil {
		return fmt.Errorf("failed to encode content: %w", err)
	}
	return nil
}

var heartbeatContent = []byte(" ")

func (r *HTTPResult) writeHeartbeatContent(w io.Writer) error {
	if _, err := w.Write(heartbeatContent); err != nil {
		return fmt.Errorf("failed to write heartbeat: %w", err)
	}
	return nil
}
