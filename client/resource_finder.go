package client

// ResourceFinder is a container of resources
type ResourceFinder interface {
	// FindResource returns a resource by name or nil if not found
	FindResource(name string) *Resource
}
