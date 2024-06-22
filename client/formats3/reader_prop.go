package formats

import "github.com/valyala/fastjson"

// ReaderProp represents a property in a JSON object
type ReaderProp struct {
	Name []byte
	*fastjson.Value
}
