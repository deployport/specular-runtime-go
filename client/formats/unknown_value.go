package formats

// UnknownValue represents an unknown value
type UnknownValue struct {
	Value
	Object map[string]UnknownValue
	Array  []UnknownValue
}

// IsEmpty returns true if the value is empty
func (v *UnknownValue) IsEmpty() bool {
	return v.Value.IsEmpty() && v.Object == nil && v.Array == nil
}
