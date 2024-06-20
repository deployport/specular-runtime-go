package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// ObjectReader represents a JSON reader
type ObjectReader struct {
	jr                    *json.Decoder
	firstTokenAlreadyRead bool
	logPrefix             string
	parseUnknownFields    bool
	unknownFields         map[string]interface{}
}

// NewObjectReaderJSON creates a new JSON reader
func NewObjectReaderJSON(jsonStream io.Reader) *ObjectReader {
	return &ObjectReader{
		jr: json.NewDecoder(jsonStream),
	}
}

func newSubObjectReader(jr *json.Decoder, logPrefix string) *ObjectReader {
	return &ObjectReader{
		jr:                    jr,
		firstTokenAlreadyRead: true,
		logPrefix:             logPrefix,
	}
}

// UseUnknownFields allows the reader to parse unknown fields
func (r *ObjectReader) UseUnknownFields() {
	r.parseUnknownFields = true
}

// UnknownFields returns the unknown fields, if any. Otherwise nil
// UseUnknownFields must have been called prior reading
func (r *ObjectReader) UnknownFields() map[string]interface{} {
	return r.unknownFields
}

// logf prints a log message
func (r *ObjectReader) logf(format string, args ...interface{}) {
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
func (r *ObjectReader) Read(readProp ReadPropFunc) error {
	dec := r.jr
	dec.UseNumber()

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
					subReader := newSubObjectReader(dec, "reader for "+currentProp.Name)
					if r.parseUnknownFields {
						subReader.UseUnknownFields()
					}
					currentProp.Value.object = subReader
				} else if v == json.Delim('[') {
					subReader := newSubArrayReader(dec, "reader for "+currentProp.Name)
					if r.parseUnknownFields {
						subReader.UseUnknownFields()
					}
					currentProp.Value.array = subReader
				}
			} else {
				return fmt.Errorf("expected value for prop %s, got %T(%v)", currentProp.Name, tk, tk)
			}
			if err := readProp(&currentProp); err != nil && err != ErrUnknownField {
				return fmt.Errorf("failed to read prop %s, %w", currentProp.Name, err)
			} else if err == ErrUnknownField {
				if r.unknownFields == nil && r.parseUnknownFields {
					r.unknownFields = make(map[string]interface{})
				}
				var finalValue any
				if !currentProp.Value.Value.IsEmpty() {
					finalValue = currentProp.Value.Value
				}
				if r.unknownFields != nil {
					r.unknownFields[currentProp.Name] = finalValue
				}
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
	return nil
}
