package client

// StructDefinition is the metadata of type
type StructDefinition struct {
	pkg   *Package
	build TypeBuilder
	path  *StructPath
}

// NewStructDefinition creates a new Type
func NewStructDefinition(pkg *Package, name StructName, build TypeBuilder) *StructDefinition {
	return &StructDefinition{
		pkg:   pkg,
		build: build,
		path:  NewStructPath(*pkg.Path(), name),
	}
}

// // Name returns the name of type
// func (t *StructDefinition) Name() string {
// 	return t.name
// }

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

// Path returns the path of type
func (t *StructDefinition) Path() *StructPath {
	return t.path
}
