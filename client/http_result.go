package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPResult is the result to send over a HTTP Response
type HTTPResult struct {
	// Struct not nil indicates the result is a struct
	Struct Struct
}

// HTTPResultForStruct returns a HTTPResult for a struct
func HTTPResultForStruct(s Struct) *HTTPResult {
	return &HTTPResult{
		Struct: s,
	}
}

// HTTPResultForError returns a HTTPResult for an error
func HTTPResultForError(err Struct) *HTTPResult {
	return HTTPResultForStruct(err)
}

// HTTPResultForHeartbeat returns a HTTPResult for a heartbeat
func HTTPResultForHeartbeat() *HTTPResult {
	return &HTTPResult{
		Struct: &Heartbeat{},
	}
}

// MimeType returns the mime type of the result based on the field set in the following order of priority: struct, err, heartbeat
func (r *HTTPResult) MimeType() string {
	return r.Struct.StructPath().MIMENameJSONHTTP()
}

// HTTPStatusCode returns the HTTP status code of the result
func (r *HTTPResult) HTTPStatusCode() int {
	if e, ok := r.Struct.(*Error); ok {
		return e.HTTPStatusCode
	}
	return http.StatusOK
}

// WriteResponse writes the result to the HTTP response
func (r *HTTPResult) WriteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", r.MimeType())
	w.WriteHeader(r.HTTPStatusCode())
	return r.WriteContent(w)
}

// WriteContent writes the content to the io.Writer
func (r *HTTPResult) WriteContent(w io.Writer) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(r.Struct); err != nil {
		return fmt.Errorf("failed to encode content: %w", err)
	}
	return nil
}
