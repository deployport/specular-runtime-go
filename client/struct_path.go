package client

import "strings"

// StructPath is the path of the struct
type StructPath struct {
	module ModulePath
	name   StructName
}

// String returns the string representation of the struct path
// as in <modulePath>.<name>
func (s *StructPath) String() string {
	return s.MIMEName()
}

// MIMEName returns the mime name of the struct
func (s *StructPath) MIMEName() string {
	var sb strings.Builder
	s.AppendMIMEName(&sb)
	return sb.String()
}

// AppendMIMEName adds the mime name of the struct to the string builder
func (s *StructPath) AppendMIMEName(sb *strings.Builder) {
	s.module.AppendMIMETreeName(sb)
	sb.WriteString(".")
	sb.WriteString(s.name.String())
}

// Module returns the module of the struct
func (s *StructPath) Module() ModulePath {
	return s.module
}

// Name returns the name of the struct
func (s *StructPath) Name() StructName {
	return s.name
}
