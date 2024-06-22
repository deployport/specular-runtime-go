package formats

// ValueType represents a value type
type ValueType int

const (
	// ValueTypeUnknown (-1) represents an unknown value type
	ValueTypeUnknown ValueType = iota - 1
	// ValueTypeNull (0) represents a null value type
	ValueTypeNull
	// ValueTypeString (1) represents a string value type
	ValueTypeString
	// ValueTypeNumber (2) represents a number value type
	ValueTypeNumber
	// ValueTypeObject (3) represents an object value type
	ValueTypeObject
	// ValueTypeArray (4) represents an array value type
	ValueTypeArray
)
