package formats

import (
	"encoding/json"
	"fmt"
	"io"
)

// ObjectReader represents a JSON reader
type ObjectReader struct {
	logger
	jr                    *json.Decoder
	firstTokenAlreadyRead bool
	parseUnknownFields    bool
	unknownFields         map[string]ComplexValue
}

// NewObjectReaderJSON creates a new JSON reader
func NewObjectReaderJSON(jsonStream io.Reader) *ObjectReader {
	return &ObjectReader{
		jr: json.NewDecoder(jsonStream),
	}
}

func newSubObjectReader(jr *json.Decoder, logPrefix string) *ObjectReader {
	r := &ObjectReader{
		jr:                    jr,
		firstTokenAlreadyRead: true,
	}
	r.logger.logPrefix = logPrefix
	return r
}

// UseUnknownFields allows the reader to parse unknown fields
func (r *ObjectReader) UseUnknownFields() {
	r.parseUnknownFields = true
}

// UnknownFields returns the unknown fields, if any. Otherwise nil
// UseUnknownFields must have been called prior reading
func (r *ObjectReader) UnknownFields() map[string]ComplexValue {
	return r.unknownFields
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
						subReader.UseUnknownItems()
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
					r.unknownFields = make(map[string]ComplexValue)
				}
				var finalValue ComplexValue
				if !currentProp.Value.Value.IsEmpty() {
					finalValue.Value = currentProp.Value.Value
				} else if vr, err := currentProp.Value.Object(); !IsValueError(err) {
					// we always parse objects so the json decoder moves forward
					if err := vr.Read(unknownPropReader); err != nil {
						return fmt.Errorf("failed to read unknown prop %s, %w", currentProp.Name, err)
					}
					finalValue.object = vr.UnknownFields()
				} else if vr, err := currentProp.Value.Array(); !IsValueError(err) {
					// we always parse arrays so the json decoder moves forward
					if err := vr.Read(unknownItemReader); err != nil {
						return fmt.Errorf("failed to load unknown array prop %s, %w", currentProp.Name, err)
					}
					finalValue.array = vr.UnknownItems()
				}
				// we always parse unknown fields so the parser moves on, but we don't retain values if they are empty or
				// unknown fields are not enabled
				if r.unknownFields != nil && !finalValue.IsEmpty() {
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

func unknownPropReader(p *ReaderProp) error {
	return ErrUnknownField
}

func unknownItemReader(i *ReaderItem) error {
	return ErrUnknownField
}
