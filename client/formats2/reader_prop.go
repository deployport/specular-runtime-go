package formats

// ReaderProp represents a property in a JSON object
type ReaderProp struct {
	Name  string
	Value ReaderValue
}

// Reset resets the reader prop
func (p *ReaderProp) Reset() {
	p.Name = ""
	p.Value.Reset()
}

// IsEmpty checks if the reader prop is empty
func (p *ReaderProp) IsEmpty() bool {
	return p.Name == "" && p.Value.IsEmpty()
}
