package server

import (
	"context"

	"go.deployport.com/specular-runtime/client"
)

// ServiceUnhandledError is the error for unhandled error
type ServiceUnhandledError struct {
	// PreOperationExecution is the pre operation execution
	// Full OperationContext might be available in context via OperationExecutionFromContext
	*PreOperationExecution
	Err error
}

// ServiceUnhandledErrorHandler is the handler for unhandled error in the service.
// It's called when an error is not handled by the service.
// The handler can return a new struct to be returned to the caller, or nil to return CallErrorCodeInternalServiceError
// Full OperationContext might be available in context via OperationExecutionFromContext
type ServiceUnhandledErrorHandler func(ctx context.Context, error ServiceUnhandledError) client.Struct

// WithServiceUnhandledErrorHandler is a function that sets the unhandled error handler as an error handler
func WithServiceUnhandledErrorHandler(handler ServiceUnhandledErrorHandler) ServiceOptionFunc {
	return func(s *Service) {
		s.unhandledErrorHandler = handler
	}
}

// unhandledErrorHandling is the unhandled error handling support for the service
type unhandledErrorHandling struct {
	unhandledErrorHandler ServiceUnhandledErrorHandler
}

// handleUnhandledError handles an unhandled error and returns the struct to be returned to the caller
// the handler can not be set or return nil
func (ue *unhandledErrorHandling) handleUnhandledError(opx *PreOperationExecution, err error) client.Struct {
	h := ue.unhandledErrorHandler
	if h == nil {
		return nil
	}
	return h(opx.Ctx, ServiceUnhandledError{
		PreOperationExecution: opx,
		Err:                   err,
	})
}
