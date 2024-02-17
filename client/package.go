package client

// Package is the metadata of package
type Package struct {
	AnnotationContainer
	path            string
	resources       map[string]*Resource
	uniqueResources map[string]*Resource
	types           map[TypeFQTNString]*StructDefinition
}

// NewPackage creates a new Package
func NewPackage(path string) *Package {
	return &Package{
		path:            path,
		resources:       map[string]*Resource{},
		uniqueResources: map[string]*Resource{},
		types:           map[TypeFQTNString]*StructDefinition{},
	}
}

// Path returns the path of package
func (s *Package) Path() string {
	return s.path
}

// NewType creates a new type with a name
func (s *Package) NewType(name string, builder TypeBuilder) (*StructDefinition, error) {
	tp := NewStructDefinition(s, name, builder)
	return tp, s.RegisterType(tp)
}

// RegisterType registers a types with a name or TypeAlreadyExistsError error
func (s *Package) RegisterType(tp *StructDefinition) error {
	fqdn := TypeFQTNString(tp.FQDN().String())
	// check if types already exists
	if _, ok := s.types[fqdn]; ok {
		return &TypeAlreadyExistsError{
			message: "types already exists",
		}
	}
	s.types[fqdn] = tp
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

// TypeByFQDN returns a type by name or TypeNotFoundError error
func (s *Package) TypeByFQDN(name TypeFQTN) (*StructDefinition, error) {
	if tp, ok := s.types[TypeFQTNString(name.String())]; ok {
		return tp, nil
	}
	return nil, &TypeNotFoundError{
		message: "type not found " + name.String(),
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

// Types returns all types
func (s *Package) Types() []*StructDefinition {
	types := make([]*StructDefinition, 0, len(s.types))
	for _, tp := range s.types {
		types = append(types, tp)
	}
	return types
}

// Import imports all types from another package
func (s *Package) Import(in *Package) error {
	for _, tp := range in.Types() {
		if err := s.RegisterType(tp); err != nil {
			return err
		}
	}
	return nil
}
