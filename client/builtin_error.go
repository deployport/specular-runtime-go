package client

import "errors"

// NewError creates a new Error
func NewError() *Error {
	s := &Error{}
	return s
}

// Error entity
type Error struct {
	Message        string
	Resource       string
	Operation      string
	HTTPStatusCode int `json:"-"`
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

// Hydrate deserializes the content into the struct
func (e *Error) Hydrate(ctx *HydratationContext) error {
	if err := ContentRequireStringProperty(ctx.Content(), "message", &e.Message); err != nil {
		return err
	}
	if err := ContentRequireStringProperty(ctx.Content(), "operation", &e.Operation); err != nil {
		return err
	}
	if err := ContentRequireStringProperty(ctx.Content(), "resource", &e.Resource); err != nil {
		return err
	}
	if err := ContentRequireStringProperty(ctx.Content(), "errorCode", &e.ErrorCode); err != nil {
		return err
	}
	return nil
}

// Dehydrate serializes the struct into the content
func (e *Error) Dehydrate(ctx *DehydrationContext) (err error) {
	// ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	ctx.Content().SetProperty("operation", e.Operation)
	ctx.Content().SetProperty("resource", e.Resource)
	ctx.Content().SetProperty("errorCode", e.ErrorCode)
	return nil
}

// StructPath returns the struct path of the struct
func (e *Error) StructPath() StructPath {
	return *NewStructPath(*BuiltinPackage().Path(), "err")
}
