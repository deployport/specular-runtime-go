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

// StructDefinitionFinder allows to find a struct definition in builtin or user package
type StructDefinitionFinder struct {
	builtinPackage *Package
	userPackage    *Package
}

// Builtin returns the builtin package of the StructDefinitionFinder
func (f *StructDefinitionFinder) Builtin() *Package {
	return f.builtinPackage
}

// User returns the user package of the StructDefinitionFinder
func (f *StructDefinitionFinder) User() *Package {
	return f.userPackage
}

// NewMultiPackageStructFinder is a finder of struct definitions in multiple packages
func NewMultiPackageStructFinder(userPkg *Package) *StructDefinitionFinder {
	return &StructDefinitionFinder{
		userPackage:    userPkg,
		builtinPackage: BuiltinPackage(),
	}
}

// Find finds a struct definition by its struct path, may return TypeNotFoundError
func (f *StructDefinitionFinder) Find(sp StructPath) (*StructDefinition, error) {
	if sp.Module().Equal(*f.builtinPackage.Path()) {
		return f.builtinPackage.TypeByPath(sp)
	}
	return f.userPackage.TypeByPath(sp)
}
