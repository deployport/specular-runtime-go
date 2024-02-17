package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client"
)

func TestOperationExecutionPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fixtures := []string{"/", ""}
		for _, fixture := range fixtures {
			t.Run("fixture "+fixture, func(t *testing.T) {
				p, err := parseOperationExecutionPathFromURIPath("/")
				require.Nil(t, p)
				require.NotNil(t, err)
				require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
				require.Equal(t, "expected a valid path to execute a resource operation", err.Message)
				require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
			})
		}
	})
	t.Run("resource empty", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("//")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
		require.Equal(t, "resource not found", err.Message)
		require.Equal(t, client.CallErrorCodeResourceNotFound, err.ErrorCode)
	})

	t.Run("resource space", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/a /")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
		require.Equal(t, "invalid resource name", err.Message)
		require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
	})
	t.Run("operation empty", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/res/")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
		require.Equal(t, "operation not found", err.Message)
		require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
		require.Equal(t, "res", err.Resource)
	})
	t.Run("operation space", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/res/op ")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
		require.Equal(t, "invalid operation name", err.Message)
		require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
		require.Equal(t, "res", err.Resource)
	})
	t.Run("valid", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/res1/op1")
		require.Nil(t, err)
		require.NotNil(t, p)
		require.Equal(t, "res1", p.ResourceName)
		require.Equal(t, "op1", p.OperationName)
	})
}
