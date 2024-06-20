package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Reader struct {
	jr *json.Decoder
}

func NewReaderJSON(jsonStream io.Reader) *Reader {
	return &Reader{
		jr: json.NewDecoder(jsonStream),
	}
}

type ReaderProp struct {
	Name  string
	Value *Value
}

type ReadPropFunc func(p *ReaderProp) error

func (r *Reader) Read(readProp ReadPropFunc) error {
	dec := r.jr
	readCount := 0
	var currentProp *ReaderProp
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		readCount++
		if readCount == 0 && tk != json.Delim('{') {
			return fmt.Errorf("expected object in JSON root, got %v", tk)
		}
		if currentProp == nil {
			// check if token is an string
			if name, ok := tk.(string); ok {
				// reading prop
				currentProp = &ReaderProp{Name: name}
				continue
			}
		}
		if currentProp != nil && currentProp.Value == nil {
			// start reading a value
			log.Printf("reading value for prop %s, %#v", currentProp.Name, tk)
			if v, ok := tk.(string); ok {
				currentProp.Value = &Value{s: &v}
			} else {
				return fmt.Errorf("expected string value for prop %s, got %v", currentProp.Name, tk)
			}
			if err := readProp(currentProp); err != nil {
				return fmt.Errorf("failed to read prop %s, %w", currentProp.Name, err)
			}
			currentProp = nil
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

type Value struct {
	s *string
}

func (v *Value) String() (*string, error) {
	if v.s == nil {
		return nil, NewValueError("expected string")
	}
	return v.s, nil
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
