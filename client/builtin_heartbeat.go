package client

// NewHeartbeat creates a new Heartbeat
func NewHeartbeat() *Heartbeat {
	s := &Heartbeat{}
	return s
}

// Heartbeat entity
type Heartbeat struct {
}

// Hydrate deserializes the content into the struct
func (e *Heartbeat) Hydrate(ctx *HydratationContext) error {
	// if err := ContentRequireStringProperty(ctx.Content(), "message", &e.Message); err != nil {
	// 	return err
	// }
	return nil
}

// Dehydrate serializes the struct into the content
func (e *Heartbeat) Dehydrate(ctx *DehydrationContext) (err error) {
	// ctx.Content().SetProperty("message", e.Message)
	return nil
}

// StructPath returns the struct path of the struct
func (e *Heartbeat) StructPath() StructPath {
	return *NewStructPath(*BuiltinPackage().Path(), "hb")
}
