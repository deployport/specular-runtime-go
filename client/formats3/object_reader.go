package formats

import (
	"fmt"
	"io"

	"github.com/valyala/fastjson"
)

// ObjectReader represents a JSON reader
type ObjectReader struct {
	logger
	stream             io.Reader
	jr                 *fastjson.Value
	openingTokenRead   bool
	parseUnknownFields bool
	// unknownFields      ComplexValueFields
}

// NewObjectReaderJSON creates a new JSON reader
func NewObjectReaderJSON(jsonStream io.Reader) *ObjectReader {
	return &ObjectReader{
		stream: jsonStream,
	}
}

func newSubObjectReader(jr *fastjson.Value, logPrefix string) *ObjectReader {
	r := &ObjectReader{
		jr:               jr,
		openingTokenRead: true,
	}
	r.logger.logPrefix = logPrefix
	return r
}

// UseUnknownFields allows the reader to parse unknown fields
func (r *ObjectReader) UseUnknownFields() {
	r.parseUnknownFields = true
}

// // UnknownFields returns the unknown fields, if any. Otherwise nil
// // UseUnknownFields must have been called prior reading
// func (r *ObjectReader) UnknownFields() ComplexValueFields {
// 	return r.unknownFields
// }

// ReadPropFunc is a function type for reading properties
type ReadPropFunc func(p *ReaderProp) error

// ReaderProp represents a property in a JSON object.
// The given property instance is reused and should not be cached outside the function.
// The function should return an error if the property could not be read.
// The function should return nil if the property was read successfully.
// If the reader doesn't know the field, it must not read the field and return ErrUnknownField
func (r *ObjectReader) Read(readProp ReadPropFunc) (err error) {
	if r.stream != nil {
		allJSON, err := io.ReadAll(r.stream)
		if err != nil {
			return fmt.Errorf("failed to read all JSON, %w", err)
		}
		r.jr, err = fastjson.Parse(string(allJSON))
		if err != nil {
			return fmt.Errorf("failed to parse JSON, %w", err)
		}
		r.stream = nil
	}

	dec := r.jr
	// dec.UseNumber()
	obj, err := dec.Object()
	if err != nil {
		return fmt.Errorf("failed to decode object, %w", err)
	}

	// if !r.openingTokenRead {
	// 	tk, err := dec.Token() // read '{'
	// 	if err != nil {
	// 		return fmt.Errorf("failed to decode token, %w", err)
	// 	}
	// 	if tk != json.Delim('{') {
	// 		return fmt.Errorf("expected object in JSON root, got %v", tk)
	// 	}
	// }
	err = nil
	var currentProp ReaderProp
	obj.Visit(func(rawKey []byte, v *fastjson.Value) {
		if err != nil {
			return
		}
		currentProp.Name = rawKey
		currentProp.Value = v
		if err := readProp(&currentProp); err != nil && err != ErrUnknownField {
			err = fmt.Errorf("failed to read prop %s, %w", currentProp.Name, err)
			return
		}
		// } else if err == ErrUnknownField {
		// 	if r.unknownFields == nil && r.parseUnknownFields {
		// 		r.unknownFields = make(ComplexValueFields)
		// 	}
		// 	var finalValue ComplexValue
		// 	if !currentProp.Value.Value.IsEmpty() {
		// 		finalValue.Value = currentProp.Value.Value
		// 	} else if vr, err := currentProp.Value.Object(); !IsValueError(err) {
		// 		// we always parse objects so the json decoder moves forward
		// 		if err := vr.Read(unknownPropReader); err != nil {
		// 			err = fmt.Errorf("failed to read unknown prop %s, %w", currentProp.Name, err)
		// 			return
		// 		}
		// 		finalValue.object = vr.UnknownFields()
		// 	} else if vr, err := currentProp.Value.Array(); !IsValueError(err) {
		// 		// we always parse arrays so the json decoder moves forward
		// 		if err := vr.Read(unknownItemReader); err != nil {
		// 			err = fmt.Errorf("failed to load unknown array prop %s, %w", currentProp.Name, err)
		// 			return
		// 		}
		// 		finalValue.array = vr.UnknownItems()
		// 	}
		// 	// we always parse unknown fields so the parser moves on, but we don't retain values if they are empty or
		// 	// unknown fields are not enabled
		// 	if r.unknownFields != nil && !finalValue.IsEmpty() {
		// 		r.unknownFields[currentProp.Name] = finalValue
		// 	}
		// }
	})
	// for dec.More() {
	// 	tk, err := dec.Token()
	// 	if err != nil {
	// 		return fmt.Errorf("failed to decode token, %w", err)
	// 	}
	// 	if r.isLoggingEnabled() {
	// 		r.logf("read token %v", tk)
	// 	}
	// 	if !currentProp.IsEmpty() && currentProp.Value.IsEmpty() {
	// 		// start reading a value
	// 		if r.isLoggingEnabled() {
	// 			r.logf("reading value for prop %s, %#v", currentProp.Name, tk)
	// 		}
	// 		if v, ok := tk.(string); ok {
	// 			currentProp.Value.s = &v
	// 		} else if tk == nil {
	// 			currentProp.Value.null = true
	// 		} else if v, ok := tk.(json.Number); ok {
	// 			currentProp.Value.number = &v
	// 		} else if v, ok := tk.(json.Delim); ok {
	// 			if v == json.Delim('{') {
	// 				// read sub object
	// 				subReader := newSubObjectReader(dec, "reader for "+currentProp.Name)
	// 				if r.parseUnknownFields {
	// 					subReader.UseUnknownFields()
	// 				}
	// 				currentProp.Value.object = subReader
	// 			} else if v == json.Delim('[') {
	// 				subReader := newSubArrayReader(dec, "reader for "+currentProp.Name)
	// 				if r.parseUnknownFields {
	// 					subReader.UseUnknownItems()
	// 				}
	// 				currentProp.Value.array = subReader
	// 			}
	// 		} else {
	// 			return fmt.Errorf("expected value for prop %s, got %T(%v)", currentProp.Name, tk, tk)
	// 		}

	// 		currentProp.Reset()
	// 		continue
	// 	}
	// }
	// tk, err := dec.Token() // read '}'
	// if err != nil {
	// 	return fmt.Errorf("failed to decode token, %w", err)
	// }
	// if tk == json.Delim('}') {
	// 	return nil
	// }
	return
}

func unknownPropReader(p *ReaderProp) error {
	return ErrUnknownField
}

// func unknownItemReader(i *ReaderItem) error {
// 	return ErrUnknownField
// }
