package client

// DehydrationContext is the context for dehydrate
type DehydrationContext struct {
	content Content
}

// Content returns the content of the context
func (ctx *DehydrationContext) Content() Content {
	return ctx.content
}

// NewDehydrationContext creates a new DehydrationContext
func NewDehydrationContext(content Content) *DehydrationContext {
	return &DehydrationContext{
		content: content,
	}
}

// CopyWithContent returns a new DehydrationContext with new content
func (ctx *DehydrationContext) CopyWithContent(content Content) *DehydrationContext {
	return NewDehydrationContext(content)
}
