package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client"
)

func TestOperationExecutionPath(t *testing.T) {
	// Routing uses the last two path segments; any leading prefix is ignored.
	t.Run("valid", func(t *testing.T) {
		fixtures := map[string]OperationExecutionPath{
			"/res1/op1":                {ResourceName: "res1", OperationName: "op1"},
			"/api/users/create":        {ResourceName: "users", OperationName: "create"},
			"/admin/api/users/create":  {ResourceName: "users", OperationName: "create"},
			"/admin/api/users/create/": {ResourceName: "users", OperationName: "create"}, // trailing slash
			"/a/b/c/d/users/create":    {ResourceName: "users", OperationName: "create"}, // any depth of prefix
		}
		for path, want := range fixtures {
			t.Run("fixture "+path, func(t *testing.T) {
				p, err := parseOperationExecutionPathFromURIPath(path)
				require.Nil(t, err)
				require.NotNil(t, p)
				require.Equal(t, want.ResourceName, p.ResourceName)
				require.Equal(t, want.OperationName, p.OperationName)
			})
		}
	})

	t.Run("too few segments", func(t *testing.T) {
		for _, fixture := range []string{"", "/", "/create", "/onlyone/"} {
			t.Run("fixture "+fixture, func(t *testing.T) {
				p, err := parseOperationExecutionPathFromURIPath(fixture)
				require.Nil(t, p)
				require.NotNil(t, err)
				require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
				require.Equal(t, "expected a valid path to execute a resource operation", err.Message)
				require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
			})
		}
	})

	t.Run("resource space", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/a /op")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.HTTPStatusCode)
		require.Equal(t, "invalid resource name", err.Message)
		require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
	})

	t.Run("empty resource segment", func(t *testing.T) {
		p, err := parseOperationExecutionPathFromURIPath("/api//create")
		require.Nil(t, p)
		require.NotNil(t, err)
		require.Equal(t, "invalid resource name", err.Message)
		require.Equal(t, client.CallErrorCodeMalformedRequest, err.ErrorCode)
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
}
