package formats

// UnknownValue represents an unknown value
type UnknownValue struct {
	Value
	object map[string]UnknownValue
	array  []UnknownValue
}

// IsEmpty returns true if the value is empty
func (v *UnknownValue) IsEmpty() bool {
	return v.Value.IsEmpty() && v.object == nil && v.array == nil
}

// Object returns the object value
func (v *UnknownValue) Object() (map[string]UnknownValue, error) {
	if v.object == nil {
		return nil, NewValueError("expected object")
	}
	return v.object, nil
}

// Array returns the array value
func (v *UnknownValue) Array() ([]UnknownValue, error) {
	if v.array == nil {
		return nil, NewValueError("expected array")
	}
	return v.array, nil
}

// Type returns the value type
func (v *UnknownValue) Type() ValueType {
	if v.object != nil {
		return ValueTypeObject
	}
	if v.array != nil {
		return ValueTypeArray
	}
	return v.Value.Type()
}
