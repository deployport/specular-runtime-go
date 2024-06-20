package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Reader represents a JSON reader
type Reader struct {
	jr                    *json.Decoder
	firstTokenAlreadyRead bool
	logPrefix             string
}

// NewReaderJSON creates a new JSON reader
func NewReaderJSON(jsonStream io.Reader) *Reader {
	return &Reader{
		jr: json.NewDecoder(jsonStream),
	}
}

func newSubReader(jr *json.Decoder, logPrefix string) *Reader {
	return &Reader{
		jr:                    jr,
		firstTokenAlreadyRead: true,
		logPrefix:             logPrefix,
	}
}

// logf prints a log message
func (r *Reader) logf(format string, args ...interface{}) {
	if r.logPrefix != "" {
		format = "(" + r.logPrefix + ") " + format
	} else {
		format = "(root reader) " + format
	}
	log.Printf(format, args...)
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
	// readCount := 0

	if !r.firstTokenAlreadyRead {
		tk, err := dec.Token() // read '{'
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		if tk != json.Delim('{') {
			return fmt.Errorf("expected object in JSON root, got %v", tk)
		}
	}
	var currentProp ReaderProp
	for dec.More() {
		tk, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to decode token, %w", err)
		}
		r.logf("read token %v", tk)
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
			r.logf("reading value for prop %s, %#v", currentProp.Name, tk)
			if v, ok := tk.(string); ok {
				currentProp.Value.s = &v
			} else if tk == nil {
				currentProp.Value.null = true
			} else if v, ok := tk.(json.Number); ok {
				currentProp.Value.number = &v
			} else if v, ok := tk.(json.Delim); ok {
				if v == json.Delim('{') {
					// read sub object
					subReader := newSubReader(dec, "reader for "+currentProp.Name)
					currentProp.Value.object = subReader
				} else if v == json.Delim('[') {
					return fmt.Errorf("array not supported")
				}
			} else {
				return fmt.Errorf("expected value for prop %s, got %T(%v)", currentProp.Name, tk, tk)
			}
			if err := readProp(&currentProp); err != nil {
				return fmt.Errorf("failed to read prop %s, %w", currentProp.Name, err)
			}
			currentProp.Reset()
			continue
		}
	}
	tk, err := dec.Token() // read '}'
	if err != nil {
		return fmt.Errorf("failed to decode token, %w", err)
	}
	if tk == json.Delim('}') {
		return nil
	}
	// if readCount == 0 {
	// 	return fmt.Errorf("no properties found")
	// }
	return nil
}
