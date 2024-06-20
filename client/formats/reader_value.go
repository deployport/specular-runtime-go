package formats

// ReaderValue is a value to be read in a reader
type ReaderValue struct {
	Value
	object *Reader
}

// Reset resets the reader value
func (v *ReaderValue) Reset() {
	v.Value.Reset()
	v.object = nil
}

// IsEmpty checks if the reader value is empty
func (v *ReaderValue) IsEmpty() bool {
	return v.Value.IsEmpty() && v.object == nil
}

// Object returns the object value
func (v *ReaderValue) Object() (*Reader, error) {
	if v.object == nil {
		return nil, NewValueError("expected object")
	}
	return v.object, nil
}
