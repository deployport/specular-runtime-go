package client

import (
	"strings"
)

// ModulePath is the name of the module
type ModulePath struct {
	namespace string
	name      string
}

// Namespace returns the namespace of the module
func (m ModulePath) Namespace() string {
	return m.namespace
}

// Name returns the name of the module
func (m ModulePath) Name() string {
	return m.name
}

// ModulePathFromTrustedValues creates a new ModulePath from namespace and name
func ModulePathFromTrustedValues(namespace, name string) *ModulePath {
	return &ModulePath{
		namespace: namespace,
		name:      name,
	}
}

// ParseModulePath parses a module path from a string
func ParseModulePath(namespace, module string) (*ModulePath, error) {
	return &ModulePath{
		namespace: strings.ToLower(namespace),
		name:      strings.ToLower(module),
	}, nil
}

// MIMETreeName returns the mime tree of the module as in application/spec.<namespace>.<module>
func (m *ModulePath) MIMETreeName() string {
	var sb strings.Builder
	m.AppendMIMETreeName(&sb)
	return sb.String()
}

const modulePathSeparator = "."

// AppendMIMETreeName adds the mime trree of the module to the string
func (m *ModulePath) AppendMIMETreeName(sb *strings.Builder) {
	sb.WriteString(ModuleMIMEApplicationSubType)
	sb.WriteString(modulePathSeparator)
	sb.WriteString(m.namespace)
	sb.WriteString(modulePathSeparator)
	sb.WriteString(m.name)
}

// String returns the string representation of the module
func (m *ModulePath) String() string {
	return m.MIMETreeName()
}

// Equal returns true if the module paths are equal
func (m ModulePath) Equal(other ModulePath) bool {
	return m.namespace == other.namespace && m.name == other.name
}
