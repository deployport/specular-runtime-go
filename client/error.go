package client

import "fmt"

// OperationAlreadyExistsError is the error for operation already exists
type OperationAlreadyExistsError struct {
	message string
}

// Error returns the error message
func (e *OperationAlreadyExistsError) Error() string {
	return e.message
}

// Is checks if the error is OperationAlreadyExistsError
func (e *OperationAlreadyExistsError) Is(err error) bool {
	_, ok := err.(*OperationAlreadyExistsError)
	return ok
}

// OperationNotFoundError is the error for operation not found
type OperationNotFoundError struct {
	message string
}

// Error returns the error message
func (e *OperationNotFoundError) Error() string {
	return e.message
}

// Is checks if the error is OperationNotFoundError
func (e *OperationNotFoundError) Is(err error) bool {
	_, ok := err.(*OperationNotFoundError)
	return ok
}

// TypeAlreadyExistsError is the error for type already exists
type TypeAlreadyExistsError struct {
	message string
}

// Error returns the error message
func (e *TypeAlreadyExistsError) Error() string {
	return e.message
}

// Is checks if the error is TypeAlreadyExistsError
func (e *TypeAlreadyExistsError) Is(err error) bool {
	_, ok := err.(*TypeAlreadyExistsError)
	return ok
}

// TypeNotFoundError is the error for type not found
type TypeNotFoundError struct {
	message string
}

// Error returns the error message
func (e *TypeNotFoundError) Error() string {
	return e.message
}

// Is checks if the error is TypeNotFoundError
func (e *TypeNotFoundError) Is(err error) bool {
	_, ok := err.(*TypeNotFoundError)
	return ok
}

// ResourceAlreadyExistsError is the error for resource already exists
type ResourceAlreadyExistsError struct {
	message string
}

// Error returns the error message
func (e *ResourceAlreadyExistsError) Error() string {
	return e.message
}

// Is checks if the error is ResourceAlreadyExistsError
func (e *ResourceAlreadyExistsError) Is(err error) bool {
	_, ok := err.(*ResourceAlreadyExistsError)
	return ok
}

// InvalidTypeFQTNError is the error for invalid type FQDN
type InvalidTypeFQTNError struct {
	message string
}

// Error returns the error message
func (e *InvalidTypeFQTNError) Error() string {
	return e.message
}

// Is checks if the error is InvalidTypeFQTNError
func (e *InvalidTypeFQTNError) Is(err error) bool {
	_, ok := err.(*InvalidTypeFQTNError)
	return ok
}

// ResourceNotFoundError is the error for resource not found
type ResourceNotFoundError struct {
	message string
}

// Error returns the error message
func (e *ResourceNotFoundError) Error() string {
	return e.message
}

// Is checks if the error is ResourceNotFoundError
func (e *ResourceNotFoundError) Is(err error) bool {
	_, ok := err.(*ResourceNotFoundError)
	return ok
}

// ArrayItemMissingError is the error for array item not found but required in schema
type ArrayItemMissingError struct {
	message string
}

// NewArrayItemMissingError returns a new ArrayItemMissingError
func NewArrayItemMissingError(propertyName string, index int) *ArrayItemMissingError {
	return &ArrayItemMissingError{
		message: fmt.Sprintf("%s[%d] missing", propertyName, index),
	}
}

// Error returns the error message
func (e *ArrayItemMissingError) Error() string {
	return e.message
}

// Is checks if the error is ArrayItemMissingError
func (e *ArrayItemMissingError) Is(err error) bool {
	_, ok := err.(*ArrayItemMissingError)
	return ok
}

// PropertyRequiredError is the error for required property not found in schema
type PropertyRequiredError struct {
	message string
}

// NewPropertyRequiredError returns a new PropertyRequiredError
func NewPropertyRequiredError(propertyName string) *PropertyRequiredError {
	return &PropertyRequiredError{
		message: fmt.Sprintf("%s is required", propertyName),
	}
}

// Error returns the error message
func (e *PropertyRequiredError) Error() string {
	return e.message
}

// Is checks if the error is PropertyRequiredError
func (e *PropertyRequiredError) Is(err error) bool {
	_, ok := err.(*PropertyRequiredError)
	return ok
}

// PropertyDehydrationError is the error the occurs while dehydrating a property
type PropertyDehydrationError struct {
	message    string
	innerError error
}

// NewPropertyDehydrationError returns a new PropertyDehydrationError
func NewPropertyDehydrationError(propertyName string, err error) *PropertyDehydrationError {
	return &PropertyDehydrationError{
		message:    fmt.Sprintf("%s failed to dehydrate, %s", propertyName, err.Error()),
		innerError: err,
	}
}

// Error returns the error message
func (e *PropertyDehydrationError) Error() string {
	return e.message
}

// Is checks if the error is PropertyDehydrationError
func (e *PropertyDehydrationError) Is(err error) bool {
	_, ok := err.(*PropertyDehydrationError)
	return ok
}

// Uwrap returns the inner error
func (e *PropertyDehydrationError) Uwrap() error {
	return e.innerError
}

// PropertyItemDehydrationError is the error the occurs while dehydrating an item in an array
type PropertyItemDehydrationError struct {
	message    string
	innerError error
}

// NewPropertyItemDehydrationError returns a new PropertyItemDehydrationError
func NewPropertyItemDehydrationError(propertyName string, index int, err error) *PropertyItemDehydrationError {
	return &PropertyItemDehydrationError{
		message:    fmt.Sprintf("%s[%d] failed to dehydrate, %s", propertyName, index, err.Error()),
		innerError: err,
	}
}

// Error returns the error message
func (e *PropertyItemDehydrationError) Error() string {
	return e.message
}

// Is checks if the error is PropertyItemDehydrationError
func (e *PropertyItemDehydrationError) Is(err error) bool {
	_, ok := err.(*PropertyItemDehydrationError)
	return ok
}

// Uwrap returns the inner error
func (e *PropertyItemDehydrationError) Uwrap() error {
	return e.innerError
}
