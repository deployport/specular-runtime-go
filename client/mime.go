package client

import (
	"fmt"
	"mime"
	"strings"
)

// MIME is the media type of the package
type MIME struct {
	// type in lowercase and without parameters
	tp     string
	params map[string]string
}

// NewMIME creates a new MIME from a string
func NewMIME(rawMime string) (MIME, error) {
	mediaType, params, err := mime.ParseMediaType(rawMime)
	if err != nil {
		return MIME{}, fmt.Errorf("error parsing media type, %w", err)
	}
	return MIME{
		tp:     strings.ToLower(mediaType),
		params: params,
	}, nil
}

// String returns the string representation of the MIME
func (m MIME) String() string {
	return string(m.tp)
}

// Type returns the type of the MIME
func (m MIME) Type() string {
	return m.tp
}

// Parameter returns the parameter of the MIME
func (m MIME) Parameter(key string) (v string, ok bool) {
	v, ok = m.params[key]
	return
}

// ModuleMIMEApplicationSubType is the mime application sub type for the module
// as in application/spec
const ModuleMIMEApplicationSubType = "application/spec"

// // IsModuleMIMEApplicationSubType returns true if the MIME is the module MIME application sub type
// func (m MIME) IsModuleMIMEApplicationSubType() bool {
// 	return strings.HasPrefix(m.String(), ModuleMIMEApplicationSubType)
// }

// StreamHeader is the header of a stream
type StreamHeader struct {
	Boundary string
}

// StreamHeader returns the Multipart header, returns nil if content type is not multipart
// and returns an error if the multipart is missing its boundary
func (m MIME) StreamHeader() (*StreamHeader, error) {
	if !strings.EqualFold(m.tp, "multipart/mixed") {
		return nil, nil
	}
	boundary, ok := m.Parameter("boundary")
	if !ok {
		return nil, fmt.Errorf("boundary missing in multipart stream")
	}
	return &StreamHeader{
		Boundary: boundary,
	}, nil
}

// StructPath returns the struct path of the media type, or nil if not a struct
// or an error if it's an invalid struct path
// parsing from media type like "application/spec.<namespace>.<module>.<struct>"
func (m MIME) StructPath() (*StructPath, error) {
	if !strings.HasPrefix(m.String(), ModuleMIMEApplicationSubType) {
		return nil, nil
	}
	parts := strings.Split(m.tp, ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid struct path")
	}
	modulePath, err := ParseModulePath(parts[0], parts[1])
	if err != nil {
		return nil, err
	}
	name := StructNameFromTrustedValue(parts[2])
	return NewStructPath(*modulePath, name), nil
}
