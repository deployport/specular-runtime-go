package formats

// ComplexValue represents an unknown value
type ComplexValue struct {
	Value
	object map[string]ComplexValue
	array  []ComplexValue
}

// NewComplexValueObject creates a new complex value array
func NewComplexValueObject(object map[string]ComplexValue) *ComplexValue {
	return &ComplexValue{
		object: object,
	}
}

// NewComplexValueArray creates a new complex value array
func NewComplexValueArray(array []ComplexValue) *ComplexValue {
	return &ComplexValue{
		array: array,
	}
}

// NewComplexValue creates a new complex value
func NewComplexValue(value Value) *ComplexValue {
	return &ComplexValue{
		Value: value,
	}
}

// IsEmpty returns true if the value is empty
func (v *ComplexValue) IsEmpty() bool {
	return v.Value.IsEmpty() && v.object == nil && v.array == nil
}

// Object returns the object value
func (v *ComplexValue) Object() (map[string]ComplexValue, error) {
	if v.object == nil {
		return nil, NewValueError("expected object")
	}
	return v.object, nil
}

// Array returns the array value
func (v *ComplexValue) Array() ([]ComplexValue, error) {
	if v.array == nil {
		return nil, NewValueError("expected array")
	}
	return v.array, nil
}

// Type returns the value type
func (v *ComplexValue) Type() ValueType {
	if v.object != nil {
		return ValueTypeObject
	}
	if v.array != nil {
		return ValueTypeArray
	}
	return v.Value.Type()
}
