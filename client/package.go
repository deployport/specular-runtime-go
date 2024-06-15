package client

import "fmt"

// Package is the metadata of package
type Package struct {
	AnnotationContainer
	path            *ModulePath
	resources       map[string]*Resource
	uniqueResources map[string]*Resource

	// all types, including imported from other packages
	allTypes map[StructPathID]*StructDefinition

	localTypes map[StructName]*StructDefinition
}

// NewPackage creates a new Package
func NewPackage(path *ModulePath) *Package {
	return &Package{
		path:            path,
		resources:       map[string]*Resource{},
		uniqueResources: map[string]*Resource{},
		allTypes:        map[StructPathID]*StructDefinition{},
		localTypes:      map[StructName]*StructDefinition{},
	}
}

// Path returns the path of package
func (s *Package) Path() *ModulePath {
	return s.path
}

// NewType creates a new type with a name, may return TypeAlreadyExistsError
func (s *Package) NewType(name string, builder TypeBuilder) (*StructDefinition, error) {
	nm, err := ParseStructName(name)
	if err != nil {
		return nil, err
	}
	tp := NewStructDefinition(s, nm, builder)
	return tp, s.registerType(tp)
}

// RegisterType registers a types with a name or TypeAlreadyExistsError error
func (s *Package) registerType(tp *StructDefinition) error {
	// check if types already exists
	isLocalType := tp.Path().Module().Equal(*s.path)
	if isLocalType {
		name := tp.path.Name()
		if _, ok := s.localTypes[name]; ok {
			return &TypeAlreadyExistsError{
				message: fmt.Sprintf("type already exists, duplicate name %s", name.String()),
			}
		}
		s.localTypes[name] = tp
	}
	s.allTypes[tp.Path().ID()] = tp
	return nil
}

// FindResource returns a resource by name or nil
func (s *Package) FindResource(name string) *Resource {
	if res, ok := s.resources[name]; ok {
		return res
	}
	return nil
}

// FindUniqueResource returns a resource by name or nil
func (s *Package) FindUniqueResource(uniqueResourceName string) *Resource {
	if res, ok := s.uniqueResources[uniqueResourceName]; ok {
		return res
	}
	return nil
}

// TypeByName returns a type by name or TypeNotFoundError error
// only local types are returned if found
func (s *Package) TypeByName(name StructName) (*StructDefinition, error) {
	if tp, ok := s.localTypes[name]; ok {
		return tp, nil
	}
	return nil, &TypeNotFoundError{
		message: "type not found " + name.String(),
	}
}

// TypeByPath returns a type by path or TypeNotFoundError error
// all types, including imported from other packages are matched
func (s *Package) TypeByPath(path StructPath) (*StructDefinition, error) {
	if tp, ok := s.allTypes[path.ID()]; ok {
		return tp, nil
	}
	return nil, &TypeNotFoundError{
		message: "type not found " + path.String(),
	}
}

// NewResource registers a new resource with a name
func (s *Package) NewResource(name string) (*Resource, error) {
	if _, ok := s.resources[name]; ok {
		return nil, &ResourceAlreadyExistsError{
			message: "resource already exists",
		}
	}
	res := newResource(s, name, name, nil)
	s.resources[name] = res
	if err := s.registerUniqueResource(res); err != nil {
		return nil, err
	}
	return res, nil
}

// registerUniqueResource registers a new resource with a unique name
func (s *Package) registerUniqueResource(res *Resource) error {
	s.uniqueResources[res.PackageUniqueName()] = res
	return nil
}

// Types returns all types registered in the package
func (s *Package) Types() []*StructDefinition {
	typesMap := s.localTypes
	types := make([]*StructDefinition, 0, len(typesMap))
	for _, tp := range typesMap {
		types = append(types, tp)
	}
	return types
}

// Import imports all types from another package
func (s *Package) Import(in *Package) error {
	for _, tp := range in.Types() {
		if err := s.registerType(tp); err != nil {
			return err
		}
	}
	return nil
}
