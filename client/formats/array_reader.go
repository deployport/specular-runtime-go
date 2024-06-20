package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
)

// ArrayReader represents a JSON reader for arrays
type ArrayReader struct {
	jr                    *json.Decoder
	firstTokenAlreadyRead bool
	logPrefix             string
	parseUnknownFields    bool
	unknownItems          []UnknownValue
}

// NewArrayReaderJSON creates a new JSON reader for arrays
func NewArrayReaderJSON(jsonStream io.Reader) *ArrayReader {
	return &ArrayReader{
		jr: json.NewDecoder(jsonStream),
	}
}

func newSubArrayReader(jr *json.Decoder, logPrefix string) *ArrayReader {
	return &ArrayReader{
		jr:                    jr,
		firstTokenAlreadyRead: true,
		logPrefix:             logPrefix,
	}
}

// UseUnknownItems allows the reader to parse unknown fields
func (r *ArrayReader) UseUnknownItems() {
	r.parseUnknownFields = true
}

// UnknownItems returns the unknown items, if any. Otherwise nil
// UseUnknownItems must have been called prior reading
func (r *ArrayReader) UnknownItems() []UnknownValue {
	return r.unknownItems
}

// logf prints a log message
func (r *ArrayReader) logf(format string, args ...interface{}) {
	if r.logPrefix != "" {
		format = "(" + r.logPrefix + ") " + format
	} else {
		format = "(root reader) " + format
	}
	log.Printf(format, args...)
}

// ReadItemFunc is a function type for reading array items
type ReadItemFunc func(i *ReaderItem) error

// Read reads the array and calls the readItem function for each item.
// The given item instance is reused and should not be retained outside the function.
// The function should return an error if the item could not be read.
// The function should return nil if the item was read successfully.
func (r *ArrayReader) Read(readItem ReadItemFunc) error {
	dec := r.jr
	dec.UseNumber()

	if !r.firstTokenAlreadyRead {
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
					if r.parseUnknownFields {
						subReader.UseUnknownFields()
					}
					currentItem.Value.object = subReader
				} else if v == json.Delim('[') {
					subReader := newSubArrayReader(dec, "reader for "+strconv.FormatInt(int64(currentItem.Index), 10))
					if r.parseUnknownFields {
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
				if r.unknownItems == nil && r.parseUnknownFields {
					r.unknownItems = make([]UnknownValue, 0, 100)
				}
				var finalValue UnknownValue
				if !currentItem.Value.Value.IsEmpty() {
					finalValue.Value = currentItem.Value.Value
				} else if vr, err := currentItem.Value.Object(); !IsValueError(err) {
					// we always parse objects so the json decoder moves forward
					if err != nil {
						return fmt.Errorf("failed to read unknown item %v, %w", currentItem.Index, err)
					}
					if err := vr.Read(unknownPropReader); err != nil {
						return fmt.Errorf("failed to read unknown item %v, %w", currentItem.Index, err)
					}
					finalValue.Object = vr.UnknownFields()
				} else if vr, err := currentItem.Value.Array(); !IsValueError(err) {
					// we always parse arrays so the json decoder moves forward
					if err != nil {
						return fmt.Errorf("failed to load unknown object item %d, %w", currentItem.Index, err)
					}
					if err := vr.Read(unknownItemReader); err != nil {
						return fmt.Errorf("failed to load unknown array item %d, %w", currentItem.Index, err)
					}
					finalValue.Array = vr.UnknownItems()
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
