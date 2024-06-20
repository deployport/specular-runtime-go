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

// UseUnknownFields allows the reader to parse unknown fields
func (r *ArrayReader) UseUnknownFields() {
	r.parseUnknownFields = true
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
						subReader.UseUnknownFields()
					}
					currentItem.Value.array = subReader
				}
			} else {
				return fmt.Errorf("expected value for item %d, got %T(%v)", currentItem.Index, tk, tk)
			}
			if err := readItem(&currentItem); err != nil {
				return fmt.Errorf("failed to read prop %d, %w", currentItem.Index, err)
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
