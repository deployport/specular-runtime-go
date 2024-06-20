package formats

import "encoding/json"

// Value represents a value in a JSON object
type Value struct {
	s      *string
	null   bool
	number *json.Number
}

// Reset resets the value
func (v *Value) Reset() {
	v.s = nil
	v.null = false
	v.number = nil
}

// IsEmpty checks if the value is empty
func (v *Value) IsEmpty() bool {
	return v.s == nil && v.number == nil && !v.null
}

func (v *Value) String() (*string, error) {
	if v.s == nil {
		return nil, NewValueError("expected string")
	}
	return v.s, nil
}

// Number returns the number value
func (v *Value) Number() (*json.Number, error) {
	if v.number == nil {
		return nil, NewValueError("expected number")
	}
	return v.number, nil
}

// IsNull checks if the value is null
func (v *Value) IsNull() bool {
	return v.null
}
