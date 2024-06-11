package client

import "errors"

// NewError creates a new Error
func NewError() *Error {
	s := &Error{}
	return s
}

// Error entity
type Error struct {
	Message string
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.GetMessage()
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

// GetMessage returns the value for the field message
func (e *Error) GetMessage() string {
	return e.Message
}

// SetMessage sets the value for the field message
func (e *Error) SetMessage(message string) {
	e.Message = message
}

// Hydrate deserializes the content into the struct
func (e *Error) Hydrate(ctx *HydratationContext) error {
	if err := ContentRequireStringProperty(ctx.Content(), "message", &e.Message); err != nil {
		return err
	}
	return nil
}

// Dehydrate serializes the struct into the content
func (e *Error) Dehydrate(ctx *DehydrationContext) (err error) {
	ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// TypeFQTN returns the Allow Typq Fully Qualified Type Name
func (e *Error) TypeFQTN() TypeFQTN {
	// TODO: get rid of the hardcoded value
	return NewTypeFQTN("proto/error", "Error")
}
