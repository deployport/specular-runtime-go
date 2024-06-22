package formats

// ReaderValue is a value to be read in a reader
type ReaderValue struct {
	Value
	object *ObjectReader
	array  *ArrayReader
}

// Reset resets the reader value
func (v *ReaderValue) Reset() {
	v.Value.Reset()
	v.object = nil
	v.array = nil
}

// IsEmpty checks if the reader value is empty
func (v *ReaderValue) IsEmpty() bool {
	return v.Value.IsEmpty() && v.object == nil && v.array == nil
}

// Object returns the object value
func (v *ReaderValue) Object() (*ObjectReader, error) {
	if v.object == nil {
		return nil, NewValueError("expected object")
	}
	return v.object, nil
}

// Array returns the array value
func (v *ReaderValue) Array() (*ArrayReader, error) {
	if v.array == nil {
		return nil, NewValueError("expected array")
	}
	return v.array, nil
}
