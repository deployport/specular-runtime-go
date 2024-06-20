package formats

import "fmt"

// ErrUnknownField is an error for unknown fields
var ErrUnknownField = fmt.Errorf("unknown field")

// ValueError represents an error related to a value.
type ValueError struct {
	msg string
}

// NewValueError creates a new ValueError
func NewValueError(msg string) *ValueError {
	return &ValueError{msg: msg}
}

// Error returns the error message.
func (e *ValueError) Error() string {
	return e.msg
}

// Is checks if the target error is a ValueError.
func (e *ValueError) Is(target error) bool {
	_, ok := target.(*ValueError)
	return ok
}

// IsValueError checks if the given error is a ValueError.
func IsValueError(err error) bool {
	_, ok := err.(*ValueError)
	return ok
}
