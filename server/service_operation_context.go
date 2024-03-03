package server

import (
	"context"
	"log/slog"
	"net/http"
)

type serviceOperationContext struct {
	opx    *OperationExecution
	logger *slog.Logger
	*opsHandlers
}

func newServiceOperationContext(
	logger *slog.Logger,
	opx *OperationExecution,
	opsHandlers *opsHandlers,
) *serviceOperationContext {
	return &serviceOperationContext{
		opx:         opx,
		logger:      logger,
		opsHandlers: opsHandlers,
	}
}

func (octx *serviceOperationContext) serve(w http.ResponseWriter) {
	if octx.opx.Operation.IsStreamed() {
		serveServiceOperationStream(octx, w)
	} else {
		serveServiceOperationSingle(octx, w)
	}
}

// Context returns the context
func (octx *serviceOperationContext) Context() context.Context {
	return octx.opx.Ctx
}
