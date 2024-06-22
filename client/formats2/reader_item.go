package formats

// ReaderItem represents an item in a JSON array
type ReaderItem struct {
	Index int
	Value ReaderValue
}

// NewReaderItem creates a new reader item
func NewReaderItem() ReaderItem {
	r := ReaderItem{}
	r.Reset()
	return r
}

// Reset resets the reader item to its default state
func (p *ReaderItem) Reset() {
	p.Index = -1
	p.Value.Reset()
}

// IsDefault checks if the reader prop is at default state
func (p *ReaderItem) IsDefault() bool {
	return p.Index == -1 && p.Value.IsEmpty()
}
