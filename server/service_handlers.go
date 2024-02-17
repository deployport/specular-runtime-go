package server

import (
	"context"

	"go.deployport.com/specular-runtime/client"
)

type opsHandlers struct {
	opHandler     OperationHandler
	streamHandler StreamHandler
}

// OperationHandler is the interface for handling an operation
type OperationHandler interface {
	HandleOperation(ctx context.Context, opx *OperationExecution) (client.Struct, error)
}

// OperationHandlerFunc is a function that implements OperationHandler
type OperationHandlerFunc func(ctx context.Context, opx *OperationExecution) (client.Struct, error)

// HandleOperation implements OperationHandler
func (f OperationHandlerFunc) HandleOperation(ctx context.Context, op *OperationExecution) (client.Struct, error) {
	return f(ctx, op)
}

// WithOperationHandler is a function that sets the operation handler
func WithOperationHandler(handler OperationHandler) ServiceOptionFunc {
	return func(s *Service) {
		s.opHandler = handler
	}
}

// StreamHandler is the interface for handling a stream
type StreamHandler interface {
	HandleStream(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error)
}

// StreamHandlerFunc is a function that implements StreamHandler
type StreamHandlerFunc func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error)

// HandleStream implements StreamHandler
func (f StreamHandlerFunc) HandleStream(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
	return f(ctx, opx)
}

// WithStreamHandler is a function that sets the stream handler
func WithStreamHandler(handler StreamHandler) ServiceOptionFunc {
	return func(s *Service) {
		s.streamHandler = handler
	}
}

// CallTypedStreamHandler is a function that calls a stream handler with typed input and output
func CallTypedStreamHandler[TInput client.Struct, TOutput client.Struct](
	ctx context.Context,
	input TInput,
	handler func(ctx context.Context, input TInput) (<-chan client.StreamEvent[TOutput], error),
) (<-chan client.StreamEvent[client.Struct], error) {
	to := make(chan client.StreamEvent[client.Struct])
	from, err := handler(ctx, input)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(to)
		for ev := range from {
			to <- client.StreamEvent[client.Struct]{
				Err:    ev.Err,
				Output: ev.Output,
			}
		}
	}()
	return to, nil
}
