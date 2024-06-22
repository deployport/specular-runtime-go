package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// ArrayReader represents a JSON reader for arrays
type ArrayReader struct {
	logger
	jr                *json.Decoder
	openingTokenRead  bool
	parseUnknownItems bool
	unknownItems      ComplexValueItems
}

// NewArrayReaderJSON creates a new JSON reader for an array in the stream
func NewArrayReaderJSON(jsonStream io.Reader) *ArrayReader {
	r := &ArrayReader{
		jr: json.NewDecoder(jsonStream),
	}
	return r
}

func newSubArrayReader(jr *json.Decoder, logPrefix string) *ArrayReader {
	r := &ArrayReader{
		jr:               jr,
		openingTokenRead: true,
	}
	r.logger.logPrefix = logPrefix
	return r
}

// UseUnknownItems allows the reader to parse unknown fields
func (r *ArrayReader) UseUnknownItems() {
	r.parseUnknownItems = true
}

// UnknownItems returns the unknown items, if any. Otherwise nil
// UseUnknownItems must have been called prior reading
func (r *ArrayReader) UnknownItems() ComplexValueItems {
	return r.unknownItems
}

// ReadItemFunc is a function type for reading array items
type ReadItemFunc func(i *ReaderItem) error

// Read reads the array and calls the readItem function for each item.
// The given item instance is reused and should not be retained outside the function.
// The function should return an error if the item could not be read.
// The function should return nil if the item was read successfully.
// If the reader doesn't know the item, it must not read the item and return ErrUnknownField
func (r *ArrayReader) Read(readItem ReadItemFunc) error {
	dec := r.jr
	dec.UseNumber()

	if !r.openingTokenRead {
		tk, err := dec.Token() // read '['
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		if tk != json.Delim('[') {
			return fmt.Errorf("expected array in JSON root, got %v", tk)
		}
	}
	currentItem := NewReaderItem()
	index := 0
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		currentItem.Index = index
		r.logf("read token %v", tk)
		if !currentItem.IsDefault() && currentItem.Value.IsEmpty() {
			// start reading a value
			r.logf("reading value for index %d, %#v", currentItem.Index, tk)
			if v, ok := tk.(string); ok {
				currentItem.Value.s = &v
			} else if tk == nil {
				currentItem.Value.null = true
			} else if v, ok := tk.(json.Number); ok {
				currentItem.Value.number = &v
			} else if v, ok := tk.(json.Delim); ok {
				if v == json.Delim('{') {
					subReader := newSubObjectReader(dec, "reader for "+strconv.FormatInt(int64(currentItem.Index), 10))
					if r.parseUnknownItems {
						subReader.UseUnknownFields()
					}
					currentItem.Value.object = subReader
				} else if v == json.Delim('[') {
					subReader := newSubArrayReader(dec, "reader for "+strconv.FormatInt(int64(currentItem.Index), 10))
					if r.parseUnknownItems {
						subReader.UseUnknownItems()
					}
					currentItem.Value.array = subReader
				}
			} else {
				return fmt.Errorf("expected value for item %d, got %T(%v)", currentItem.Index, tk, tk)
			}
			if err := readItem(&currentItem); err != nil && err != ErrUnknownField {
				return fmt.Errorf("failed to read prop %d, %w", currentItem.Index, err)
			} else if err == ErrUnknownField {
				if r.unknownItems == nil && r.parseUnknownItems {
					r.unknownItems = make(ComplexValueItems, 0, 100)
				}
				var finalValue ComplexValue
				if !currentItem.Value.Value.IsEmpty() {
					finalValue.Value = currentItem.Value.Value
				} else if vr, err := currentItem.Value.Object(); !IsValueError(err) {
					// we always parse objects so the json decoder moves forward
					if err := vr.Read(unknownPropReader); err != nil {
						return fmt.Errorf("failed to read unknown item %v, %w", currentItem.Index, err)
					}
					finalValue.object = vr.UnknownFields()
				} else if vr, err := currentItem.Value.Array(); !IsValueError(err) {
					// we always parse arrays so the json decoder moves forward
					if err := vr.Read(unknownItemReader); err != nil {
						return fmt.Errorf("failed to load unknown array item %d, %w", currentItem.Index, err)
					}
					finalValue.array = vr.UnknownItems()
				}
				// we always parse unknown fields so the parser moves on, but we don't retain values if they are empty or
				// unknown fields are not enabled
				if r.unknownItems != nil && !finalValue.IsEmpty() {
					r.unknownItems = append(r.unknownItems, finalValue)
				}
			}
			currentItem.Reset()
			index++
			continue
		}
	}
	tk, err := dec.Token() // read ']'
	if err != nil {
		return fmt.Errorf("failed to decode token, %w", err)
	}
	if tk == json.Delim(']') {
		return nil
	}
	return nil
}
