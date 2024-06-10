package client

// StructName is the name of the struct
type StructName string

// String returns the string representation of the struct name
func (s StructName) String() string {
	return string(s)
}
