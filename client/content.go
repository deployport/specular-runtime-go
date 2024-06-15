package client

// Content is an incoming our outgoing struct
type Content map[string]any

// NewContent creates a new Content
func NewContent() Content {
	c := Content{}
	return c
}

// SetProperty sets a property value
func (c Content) SetProperty(name string, value any) {
	c[name] = value
}

// GetProperty sets a property value, nil if not found
func (c Content) GetProperty(name string) any {
	v, ok := c[name]
	if !ok {
		return nil
	}
	return v
}

// Map returns the underlying map
func (c Content) Map() map[string]any {
	return c
}

// StructFromContent hydrates a struct from a content
func StructFromContent(content Content, pkg *Package, st Struct) error {
	return st.Hydrate(NewHydratationContext(pkg, content))
}
