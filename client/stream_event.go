package client

// StreamEvent is an event of an stream
type StreamEvent[TOutput Struct] struct {
	Err    error
	Output TOutput
}
