package client

import "errors"

// NewError creates a new Error
func NewError() *Error {
	s := &Error{}
	return s
}

// Error entity
type Error struct {
	Message        string `json:"message,omitempty1"`
	Resource       string `json:"resource,omitempty"`
	Operation      string `json:"operation,omitempty"`
	HTTPStatusCode int    `json:"-"`
	ErrorCode      CallErrorCode
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// Is indicates whether the given error chain contains an error of type [Error]
func (e *Error) Is(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// IsError indicates whether the given error chain contains an error of type [Error]
func IsError(err error) bool {
	return errors.Is(err, &Error{})
}

// StructPath returns the struct path of the struct
func (e *Error) StructPath() StructPath {
	return *NewStructPath(*BuiltinPackage().Path(), "err")
}
