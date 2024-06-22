package formats

import "encoding/json"

// Value represents a value in a JSON object
type Value struct {
	s      *string
	null   bool
	number *json.Number
}

// NewValueString creates a new value with a string
func NewValueString(s string) Value {
	return Value{s: &s}
}

// NewValueNumber creates a new value with a number
func NewValueNumber(n json.Number) Value {
	return Value{number: &n}
}

// NewValueNull creates a new value with a null
func NewValueNull() Value {
	return Value{null: true}
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

// MarshalJSON marshals the value to JSON
func (v *Value) MarshalJSON() ([]byte, error) {
	if v.s != nil {
		return json.Marshal(*v.s)
	}
	if v.number != nil {
		return json.Marshal(*v.number)
	}
	return json.Marshal(nil)
}

// Type returns the value type
func (v *Value) Type() ValueType {
	if v.s != nil {
		return ValueTypeString
	}
	if v.number != nil {
		return ValueTypeNumber
	}
	if v.null {
		return ValueTypeNull
	}
	return ValueTypeUnknown
}
