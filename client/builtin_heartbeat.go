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
	ctx.Content().SetStruct(e.TypeFQTN().String())
	// ctx.Content().SetProperty("message", e.Message)
	return nil
}

// TypeFQTN returns the Allow Typq Fully Qualified Type Name
func (e *Heartbeat) TypeFQTN() TypeFQTN {
	// TODO: get rid of the hardcoded value
	return NewTypeFQTN("proto", "hb")
}
