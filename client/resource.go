package client

// Resource is the metadata of resource
type Resource struct {
	pkg               *Package
	name              string
	packageUniqueName string
	ops               map[string]*Operation
	list              []*Operation
	parent            *Resource
	subResources      map[string]*Resource
}

func newResource(pkg *Package, name string, packageUniqueName string, parent *Resource) *Resource {
	return &Resource{
		pkg:               pkg,
		name:              name,
		packageUniqueName: packageUniqueName,
		ops:               make(map[string]*Operation),
		list:              make([]*Operation, 0, 50),
		parent:            parent,
		subResources:      make(map[string]*Resource),
	}
}

// Name returns the name of resource
func (res *Resource) Name() string {
	return res.name
}

// PackageUniqueName returns the unique name of the resource in the package
func (res *Resource) PackageUniqueName() string {
	return res.packageUniqueName
}

// Package returns the service of resource
func (res *Resource) Package() *Package {
	return res.pkg
}

// NewSubResource registers a new sub-resource with a name
func (res *Resource) NewSubResource(name string) (*Resource, error) {
	resources := res.subResources
	if _, ok := resources[name]; ok {
		return nil, &ResourceAlreadyExistsError{
			message: "resource already exists",
		}
	}
	nr := newResource(res.Package(), name, res.PackageUniqueName()+name, res.parent)
	resources[name] = nr
	if err := res.pkg.registerUniqueResource(nr); err != nil {
		return nil, err
	}
	return nr, nil
}

// FindResource returns a resource by name or ResourceNotFoundError error
func (res *Resource) FindResource(name string) *Resource {
	if res, ok := res.subResources[name]; ok {
		return res
	}
	return nil
}

// NewOperation creates a new Operation
func (res *Resource) NewOperation(name string) (*Operation, error) {
	if _, ok := res.ops[name]; ok {
		return nil, &OperationAlreadyExistsError{
			message: "operation already exists",
		}
	}
	op := newOperation(res, name)
	res.ops[name] = op
	res.list = append(res.list, op)
	return op, nil
}

// FindOperation finds an operation by name, returns OperationNotFoundError if not found
func (res *Resource) FindOperation(name string) *Operation {
	if op, ok := res.ops[name]; ok {
		return op
	}
	return nil
}

// Operations returns all operations
func (res *Resource) Operations() []*Operation {
	return res.list
}
