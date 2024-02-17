package client

// Type is the interface for all types instance, mostly structs and annotations
type Type interface {
	TypeFQTN() TypeFQTN
}
