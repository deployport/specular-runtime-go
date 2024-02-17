package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPError is an error that can be returned by the HTTP transport
type HTTPError struct {
	Message        string        `json:"message"`
	Resource       string        `json:"resource"`
	Operation      string        `json:"operation"`
	HTTPStatusCode int           `json:"-"`
	ErrorCode      CallErrorCode `json:"code"`
}

// Error returns the error message
func (err *HTTPError) Error() string {
	return fmt.Sprintf("%s: %s, while attempting to execute %s/%s", err.ErrorCode, err.Message, err.Resource, err.Operation)
}

// Is implements errors.Is
func (err *HTTPError) Is(target error) bool {
	_, ok := target.(*HTTPError)
	return ok
}

// HTTPErrorInternalServiceError returns an internal service error HTTPError
func HTTPErrorInternalServiceError(resourceName string, operationName string) HTTPError {
	return HTTPError{
		Message:        "Internal Service Error",
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorCode:      CallErrorCodeInternalServiceError,
		Resource:       resourceName,
		Operation:      operationName,
	}
}

// writeContent writes the content to the io.Writer
func (err *HTTPError) writeContent(w io.Writer) error {
	enc := json.NewEncoder(w)
	if err := enc.Encode(err); err != nil {
		return fmt.Errorf("failed to encode http error: %w", err)
	}
	return nil
}

// CallErrorCode is the error code for a call
type CallErrorCode string

// CallErrorCodeResourceNotFound is returned when a resource is not found
const CallErrorCodeResourceNotFound CallErrorCode = "RESOURCE_NOT_FOUND"

// CallErrorCodeOperationNotFound is returned when an operation is not found
const CallErrorCodeOperationNotFound CallErrorCode = "OPERATION_NOT_FOUND"

// CallErrorCodeInvalidInput is returned when the input is invalid
const CallErrorCodeInvalidInput CallErrorCode = "INVALID_INPUT"

// CallErrorCodeInvalidStruct is returned when the input is invalid
const CallErrorCodeInvalidStruct CallErrorCode = "INVALID_STRUCT"

// CallErrorCodeInternalServiceError is returned when the service returns an internal error
const CallErrorCodeInternalServiceError CallErrorCode = "INTERNAL_SERVICE_ERROR"

// CallErrorCodeMalformedRequest is returned when the request is invalid
const CallErrorCodeMalformedRequest CallErrorCode = "MALFORMED_REQUEST"

// CallErrorCodeAccessDenied is returned when the authentication fails
const CallErrorCodeAccessDenied CallErrorCode = "ACCESS_DENIED"
