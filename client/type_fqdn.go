package client

import (
	"fmt"
	"strings"
)

// TypeFQTN is the fully qualified name of type (<package-path>:<type-name>)
type TypeFQTN struct {
	packagePath string
	typeName    string
	s           string
}

// NewTypeFQTN creates a new TypeFQTN
func NewTypeFQTN(packagePath, typeName string) TypeFQTN {
	return TypeFQTN{
		packagePath: packagePath,
		typeName:    typeName,
		s:           packagePath + ":" + typeName,
	}
}

// PackagePath returns the package path of the type
func (t *TypeFQTN) PackagePath() string {
	return t.packagePath
}

// TypeName returns the type name of the type
func (t *TypeFQTN) TypeName() string {
	return t.typeName
}

// String returns the string representation of the TypeFQTN
func (t TypeFQTN) String() string {
	return t.s
}

// TypeFQTNFromString creates a new TypeFQTN from string using format (<package-path>:<type-name>)
func TypeFQTNFromString(s string) (TypeFQTN, error) {
	if !strings.ContainsRune(s, ':') {
		return TypeFQTN{}, &InvalidTypeFQTNError{
			message: fmt.Sprintf("invalid type FQDN, %s", s),
		}
	}
	packagePath := s[0:strings.IndexRune(s, ':')]
	typeName := s[strings.IndexRune(s, ':')+1:]
	return NewTypeFQTN(packagePath, typeName), nil
}

// TypeFQTNString is a string type for TypeFQTN
type TypeFQTNString string
