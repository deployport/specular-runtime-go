package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Reader represents a JSON reader
type Reader struct {
	jr *json.Decoder
}

// NewReaderJSON creates a new JSON reader
func NewReaderJSON(jsonStream io.Reader) *Reader {
	return &Reader{
		jr: json.NewDecoder(jsonStream),
	}
}

// ReadPropFunc is a function type for reading properties
type ReadPropFunc func(p *ReaderProp) error

// ReaderProp represents a property in a JSON object.
// The given property instance is reused and should not be cached outside the function.
// The function should return an error if the property could not be read.
// The function should return nil if the property was read successfully.
func (r *Reader) Read(readProp ReadPropFunc) error {
	dec := r.jr
	dec.UseNumber()
	readCount := 0
	var currentProp ReaderProp
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		readCount++
		if readCount == 0 && tk != json.Delim('{') {
			return fmt.Errorf("expected object in JSON root, got %v", tk)
		}
		if currentProp.IsEmpty() {
			// check if token is an string
			if name, ok := tk.(string); ok {
				// reading prop
				currentProp.Name = name
				continue
			}
		}
		if !currentProp.IsEmpty() && currentProp.Value.IsEmpty() {
			// start reading a value
			log.Printf("reading value for prop %s, %#v", currentProp.Name, tk)
			if v, ok := tk.(string); ok {
				currentProp.Value.s = &v
			} else if tk == nil {
				currentProp.Value.null = true
			} else if v, ok := tk.(json.Number); ok {
				currentProp.Value.number = &v
			} else {
				return fmt.Errorf("expected value for prop %s, got %T(%v)", currentProp.Name, tk, tk)
			}
			if err := readProp(&currentProp); err != nil {
				return fmt.Errorf("failed to read prop %s, %w", currentProp.Name, err)
			}
			currentProp.Reset()
			continue
		}
		if readCount > 0 && tk == json.Delim('}') {
			return nil
		}
	}
	if readCount == 0 {
		return fmt.Errorf("no properties found")
	}
	return nil
}

// ValueError represents an error related to a value.
type ValueError struct {
	msg string
}

// NewValueError creates a new ValueError
func NewValueError(msg string) *ValueError {
	return &ValueError{msg: msg}
}

// Error returns the error message.
func (e *ValueError) Error() string {
	return e.msg
}

// Is checks if the target error is a ValueError.
func (e *ValueError) Is(target error) bool {
	_, ok := target.(*ValueError)
	return ok
}

// IsValueError checks if the given error is a ValueError.
func IsValueError(err error) bool {
	_, ok := err.(*ValueError)
	return ok
}
