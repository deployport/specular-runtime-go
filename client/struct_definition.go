package client

// StructDefinition is the metadata of type
type StructDefinition struct {
	name  string
	pkg   *Package
	build TypeBuilder
	fqdn  TypeFQTN
}

// NewStructDefinition creates a new Type
func NewStructDefinition(pkg *Package, name string, build TypeBuilder) *StructDefinition {
	return &StructDefinition{
		name:  name,
		pkg:   pkg,
		build: build,
		fqdn:  NewTypeFQTN(pkg.Path(), name),
	}
}

// Name returns the name of type
func (t *StructDefinition) Name() string {
	return t.name
}

// Package returns the package of type
func (t *StructDefinition) Package() *Package {
	return t.pkg
}

// TypeBuilder returns the build an instance(pointer) of a type
func (t *StructDefinition) TypeBuilder() TypeBuilder {
	return t.build
}

// TypeBuilder is a function that returns an instance of the type
type TypeBuilder func() Struct

// FQDN returns the fully qualified name
func (t *StructDefinition) FQDN() TypeFQTN {
	return t.fqdn
}
