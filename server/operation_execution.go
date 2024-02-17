package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"go.deployport.com/specular-runtime/client"
)

// PreOperationExecution is the execution of an operation
type PreOperationExecution struct {
	*unhandledErrorHandling
	Ctx       context.Context
	Logger    *slog.Logger
	Path      OperationExecutionPath
	Operation *client.Operation
}

func (pre *PreOperationExecution) createInternalServiceError() client.HTTPError {
	return client.HTTPErrorInternalServiceError(pre.Operation.Resource().Name(), pre.Operation.Name())
}

// OperationExecution is the execution of an operation
type OperationExecution struct {
	PreOperationExecution
	Input client.Struct
}

type operationExecutionContextKey struct{}

var operationExecutionContextKeyValue = operationExecutionContextKey{}

// ContextWithOperationExecution returns a context with the operation execution embedded as key
func ContextWithOperationExecution(ctx context.Context, opx *OperationExecution) context.Context {
	return context.WithValue(ctx, operationExecutionContextKeyValue, opx)
}

// OperationExecutionFromContext returns the operation execution from the context, otherwise nil
func OperationExecutionFromContext(ctx context.Context) *OperationExecution {
	opx, _ := ctx.Value(operationExecutionContextKeyValue).(*OperationExecution)
	return opx
}

func buildPreOperationExecution(
	ctx context.Context,
	logger *slog.Logger,
	unhandledErrorHandling *unhandledErrorHandling,
	pk *client.Package,
	path *OperationExecutionPath,
) (*PreOperationExecution, *client.HTTPError) {
	resourceName := path.ResourceName
	operationName := path.OperationName
	resource := pk.FindUniqueResource(resourceName)
	if resource == nil {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "resource not found",
			Resource:       resourceName,
			Operation:      operationName,
			ErrorCode:      client.CallErrorCodeResourceNotFound,
		}
	}

	op := resource.FindOperation(operationName)
	if op == nil {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "operation not found",
			Resource:       resourceName,
			Operation:      operationName,
			ErrorCode:      client.CallErrorCodeOperationNotFound,
		}
	}
	return &PreOperationExecution{
		unhandledErrorHandling: unhandledErrorHandling,
		Ctx:                    ctx,
		Logger:                 logger,
		Path:                   *path,
		Operation:              op,
	}, nil
}

func buildOperationExecution(
	logger *slog.Logger,
	pk *client.Package,
	pre *PreOperationExecution,
	requestBody io.Reader,
) (*OperationExecution, *client.HTTPError) {
	resourceName := pre.Path.ResourceName
	operationName := pre.Path.OperationName
	logger.Debug("operation call", "resource-name", resourceName, "operation-name", operationName)
	op := pre.Operation
	inputContent := client.NewContent()
	dec := json.NewDecoder(requestBody)
	if err := dec.Decode(&inputContent); err != nil {
		logger.Warn("failed to decode input content", "err", err)
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        "expected valid JSON body",
			Resource:       resourceName,
			Operation:      operationName,
			ErrorCode:      client.CallErrorCodeInvalidInput,
		}
	}
	input := op.Input().TypeBuilder()()
	if err := client.StructFromContent(inputContent, pk, input); err != nil {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusBadRequest,
			Message:        fmt.Sprintf("expected a valid input of type %s", input.TypeFQTN()),
			Resource:       resourceName,
			Operation:      operationName,
			ErrorCode:      client.CallErrorCodeInvalidStruct,
		}
	}
	opx := &OperationExecution{
		PreOperationExecution: *pre,
		Input:                 input,
	}
	pre.Ctx = ContextWithOperationExecution(pre.Ctx, opx)
	return opx, nil
}

// BuildFinalResult builds the final result for the execution of an operation
// if errOrStruct is an error, it will be handled using the unhandledErrorHandling
// if errOrStruct is a struct, it will be returned as is
// if errOrStruct is nil, the result will be returned as is
func (pre *PreOperationExecution) BuildFinalResult(
	result *client.HTTPResult,
	errOrStruct error,
) *client.HTTPResult {
	if errOrStruct != nil {
		if errStruct, ok := errOrStruct.(client.Struct); ok {
			return client.HTTPResultForStruct(errStruct)
		}
		err := errOrStruct
		if IsAuthenticationError(err) {
			return client.HTTPResultForError(client.HTTPError{
				Message:        err.Error(),
				Resource:       pre.Operation.Resource().Name(),
				Operation:      pre.Operation.Name(),
				ErrorCode:      client.CallErrorCodeAccessDenied,
				HTTPStatusCode: http.StatusForbidden,
			})
		}
		handledErrStruct := pre.handleUnhandledError(pre, err)
		if handledErrStruct != nil {
			return client.HTTPResultForStruct(handledErrStruct)
		}
		pre.Logger.Warn("internal server error", "err", err)
		return client.HTTPResultForError(pre.createInternalServiceError())
	}
	return result
}
