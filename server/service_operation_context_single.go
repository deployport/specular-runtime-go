package server

import (
	"fmt"
	"net/http"

	"go.deployport.com/specular-runtime/client"
)

type serviceOperationSingleServeContext struct {
	*serviceOperationContext
	w http.ResponseWriter
}

func serveServiceOperationSingle(
	octx *serviceOperationContext,
	w http.ResponseWriter,
) {
	soctx := &serviceOperationSingleServeContext{
		serviceOperationContext: octx,
		w:                       w,
	}
	soctx.serveSingle()
}

// handleOperation handles an operation
func (octx *serviceOperationSingleServeContext) serveSingle() {
	result := octx.opx.BuildFinalResult(octx.executeSingle())
	_ = result.WriteResponse(octx.w)
}

// executeSingle executes the operation and expects a single result, err could be a struct
func (octx *serviceOperationSingleServeContext) executeSingle() (result *client.HTTPResult, err error) {
	handler := octx.opHandler
	if handler == nil {
		return nil, fmt.Errorf("operation handler unset, service unable to serve requests")
	}
	output, err := handler.HandleOperation(octx.Context(), octx.opx)
	if err != nil {
		return nil, err
	}
	return client.HTTPResultForStruct(output), nil
}
