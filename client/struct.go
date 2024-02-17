package client

// Struct is the interface for all struct instance
type Struct interface {
	Type
	Hydrate(ctx *HydratationContext) error
	Dehydrate(ctx *DehydrationContext) error
}
