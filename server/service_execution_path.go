package server

import (
	"net/http"
	"strings"

	"go.deployport.com/specular-runtime/client"
)

// OperationExecutionPath is the resource and operation name
type OperationExecutionPath struct {
	ResourceName  string
	OperationName string
}

// parseOperationExecutionPathFromURIPath resolves the resource and operation
// from the request URI path.
//
// Specular routes on the LAST TWO path segments only: the second-to-last segment
// is the resource (its package-unique name, always a single segment, since
// nested resources are concatenated without a separator) and the last segment is
// the operation. Any leading prefix segments (for example a mount path like
// "/admin/api") are ignored for routing
func parseOperationExecutionPathFromURIPath(path string) (*OperationExecutionPath, *client.Error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 2 {
		return nil, &client.Error{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "expected a valid path to execute a resource operation",
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}
	resourceName := segments[len(segments)-2]
	operationName := segments[len(segments)-1]

	if resourceName == "" || strings.Contains(resourceName, " ") {
		return nil, &client.Error{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "invalid resource name",
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}
	if operationName == "" || strings.Contains(operationName, " ") {
		return nil, &client.Error{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "invalid operation name",
			Resource:       resourceName,
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}
	return &OperationExecutionPath{
		ResourceName:  resourceName,
		OperationName: operationName,
	}, nil
}
