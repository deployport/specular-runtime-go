package client

// NewHeartbeat creates a new Heartbeat
func NewHeartbeat() *Heartbeat {
	s := &Heartbeat{}
	return s
}

// Heartbeat entity
type Heartbeat struct{}

// StructPath returns the struct path of the struct
func (e *Heartbeat) StructPath() StructPath {
	return *NewStructPath(*BuiltinPackage().Path(), "hb")
}
