package client

import "strings"

// StructName is the name of the struct
type StructName string

// StructNameFromTrustedValue creates a new StructName from a string
func StructNameFromTrustedValue(name string) StructName {
	return StructName(name)
}

// ParseStructName parses a string and returns a StructName
func ParseStructName(name string) (StructName, error) {
	// TODO: validate name
	return StructName(strings.ToLower(name)), nil
}

// String returns the string representation of the struct name
func (s StructName) String() string {
	return string(s)
}
