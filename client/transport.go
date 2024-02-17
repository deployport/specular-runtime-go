package client

import "context"

// Transport is the transport for the client
type Transport interface {
	// Perform performs the request for an operation
	Execute(ctx context.Context, req *Request) (Struct, error)

	// Stream performs the request for an operation and returns a channel with the events
	Stream(ctx context.Context, req *Request) (<-chan StreamEvent[Struct], error)
}

// Request is the request to send
type Request struct {
	// The operation name
	Operation *Operation

	// The request body
	Input Struct
}

// Response is the response from a request
type Response struct {
	// Output is the instance of the struct containing the output of the operation
	Output Type

	// When output is not set, Problem is the instance of the struct containing the problem
	Problem Type
}
