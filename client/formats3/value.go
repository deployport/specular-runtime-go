package formats

import (
	"github.com/valyala/fastjson"
)

// Value represents a value in a JSON object
type Value struct {
	v *fastjson.Value
}

// // NewValueString creates a new value with a string
// func NewValueString(s string) Value {
// 	return Value{s: &s}
// }

// // NewValueNumber creates a new value with a number
// func NewValueNumber(n json.Number) Value {
// 	return Value{number: &n}
// }

// // NewValueNull creates a new value with a null
// func NewValueNull() Value {
// 	return Value{null: true}
// }

// // Reset resets the value
// func (v *Value) Reset() {
// 	v.s = nil
// 	v.null = false
// 	v.number = nil
// }

// // IsEmpty checks if the value is empty
// func (v *Value) IsEmpty() bool {
// 	return v.s == nil && v.number == nil && !v.null
// }
