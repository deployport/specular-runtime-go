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
	// IsError marks the result as an error envelope. Error envelopes are sent
	// with the `kind=error` media-type parameter and an error HTTP status code
	// so clients can surface them as errors without knowing the concrete type.
	IsError bool
}

// HTTPResultForStruct returns a HTTPResult for a regular output struct
func HTTPResultForStruct(s Struct) *HTTPResult {
	return &HTTPResult{
		Struct: s,
	}
}

// HTTPResultForErrorStruct returns a HTTPResult for a user-defined error struct
func HTTPResultForErrorStruct(s Struct) *HTTPResult {
	return &HTTPResult{
		Struct:  s,
		IsError: true,
	}
}

// HTTPResultForError returns a HTTPResult for an error
func HTTPResultForError(err Struct) *HTTPResult {
	return &HTTPResult{
		Struct:  err,
		IsError: true,
	}
}

// HTTPResultForHeartbeat returns a HTTPResult for a heartbeat
func HTTPResultForHeartbeat() *HTTPResult {
	return &HTTPResult{
		Struct: &Heartbeat{},
	}
}

// errorMIMEParameter marks an envelope as carrying an error struct. It is
// appended to the Content-Type of error results so clients can detect an error
// independently of whether they can resolve its concrete type.
const errorMIMEParameter = "; kind=error"

// MimeType returns the mime type of the result based on the field set in the following order of priority: struct, err, heartbeat
func (r *HTTPResult) MimeType() string {
	mime := r.Struct.StructPath().MIMENameJSONHTTP()
	if r.IsError {
		mime += errorMIMEParameter
	}
	return mime
}

// HTTPStatusCarrier is implemented by a user-defined error struct that declares
// its own HTTP status with @specular:Status. The generator emits the method
// only for an annotated struct, so an error that was never annotated does not
// satisfy this interface and keeps the 400 default.
type HTTPStatusCarrier interface {
	// HTTPStatus returns the advisory status the error is served with.
	HTTPStatus() int
}

// HTTPStatusCode returns the result's HTTP status code, which is advisory: the
// envelope (content type plus kind) is authoritative. Built-in errors carry
// their own status, a user-defined error carries the status it declared with
// @specular:Status, every other user-defined error uses 400, and regular
// outputs use 200.
func (r *HTTPResult) HTTPStatusCode() int {
	// The built-in branch stays first on purpose. A built-in error already
	// carries a status from the reserved table, and nothing generated may
	// override it.
	if e, ok := r.Struct.(*Error); ok {
		return e.HTTPStatusCode
	}
	if r.IsError {
		// The declared status is read inside the error branch, never outside
		// it, so a 4xx or a 5xx can only ever accompany kind=error. A struct
		// sent as a regular output keeps 200 even if it carries a status.
		if s, ok := r.Struct.(HTTPStatusCarrier); ok {
			return s.HTTPStatus()
		}
		return http.StatusBadRequest
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
