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

// parseOperationExecutionPathFromURIPath returns the operation execution header from the URI path
func parseOperationExecutionPathFromURIPath(path string) (*OperationExecutionPath, *client.HTTPError) {
	// path is of the form /<resource-name>/<operation-name>
	// we need to extract the resource name and operation name
	// we can do this by finding the first slash and then taking the first part as the resource name and the second part as the operation name
	rp := strings.TrimPrefix(path, "/")
	slashPathIndex := strings.Index(rp, "/")
	if slashPathIndex == -1 {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "expected a valid path to execute a resource operation",
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}
	// resource name is the first part of the path and is the package unique resource name
	resourceName := rp[0:slashPathIndex]

	// check if resource name has empty spaces
	if strings.Contains(resourceName, " ") {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "invalid resource name",
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}

	if resourceName == "" {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "resource not found",
			ErrorCode:      client.CallErrorCodeResourceNotFound,
		}
	}
	operationName := rp[slashPathIndex+1:]
	if operationName == "" {
		return nil, &client.HTTPError{
			HTTPStatusCode: http.StatusNotFound,
			Message:        "operation not found",
			Resource:       resourceName,
			ErrorCode:      client.CallErrorCodeMalformedRequest,
		}
	}

	// check if operation name has empty spaces
	if strings.Contains(operationName, " ") {
		return nil, &client.HTTPError{
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
