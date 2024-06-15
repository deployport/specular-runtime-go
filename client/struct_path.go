package client

import "strings"

// StructPathID is the id of the struct path
// used in maps where a unique value per struct path is necessary
type StructPathID string

// StructPath is the path of the struct
type StructPath struct {
	module ModulePath
	name   StructName
}

// NewStructPath creates a new struct path
func NewStructPath(module ModulePath, name StructName) *StructPath {
	return &StructPath{
		module: module,
		name:   name,
	}
}

// String returns the string representation of the struct path
// as in <modulePath>.<name>
func (s StructPath) String() string {
	return s.MIMEName()
}

// ID returns the id of the struct
func (s StructPath) ID() StructPathID {
	return StructPathID(s.String())
}

// MIMEName returns the mime name of the struct
func (s StructPath) MIMEName() string {
	var sb strings.Builder
	s.AppendMIMEName(&sb)
	return sb.String()
}

// MIMENameJSONHTTP returns the mime name of the struct +json media type sub type
func (s StructPath) MIMENameJSONHTTP() string {
	var sb strings.Builder
	s.AppendMIMEName(&sb)
	sb.WriteString("+json")
	return sb.String()
}

// AppendMIMEName adds the mime name of the struct to the string builder
func (s *StructPath) AppendMIMEName(sb *strings.Builder) {
	s.module.AppendMIMETreeName(sb)
	sb.WriteString(".")
	sb.WriteString(s.name.String())
}

// Module returns the module of the struct
func (s StructPath) Module() ModulePath {
	return s.module
}

// Name returns the name of the struct
func (s StructPath) Name() StructName {
	return s.name
}
