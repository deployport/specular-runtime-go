package client

// HydratationContext is the context for hydrate
type HydratationContext struct {
	content Content
	pkg     *Package
}

// Content returns the content of the context
func (ctx *HydratationContext) Content() Content {
	return ctx.content
}

// Package returns the package of the context
func (ctx *HydratationContext) Package() *Package {
	return ctx.pkg
}

// NewHydratationContext creates a new HydratationContext
func NewHydratationContext(pkg *Package, content Content) *HydratationContext {
	return &HydratationContext{
		content: content,
		pkg:     pkg,
	}
}

// CopyWithContent returns a new HydratationContext a new content
func (ctx *HydratationContext) CopyWithContent(content Content) HydratationContext {
	return *NewHydratationContext(ctx.pkg, content)
}
