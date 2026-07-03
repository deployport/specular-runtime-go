package server

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client"
)

func TestService(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)

	resource := pk.FindResource("Book")
	opCreate := resource.FindOperation("Create")

	svc := NewService(
		WithServicePackage(pk),
		WithOperationHandler(OperationHandlerFunc(func(ctx context.Context, op *OperationExecution) (client.Struct, error) {
			// op.Operation.Annotations()
			if op.Operation == opCreate {
				input := op.Input.(*BookCreateInput)
				if input.Name == "inval1d" {
					return nil, &BookCreationProblem{
						Message: "invalid name",
					}
				}
				return &BookCreateOutput{
					ID: "123",
				}, nil
			}
			return nil, nil
		})),
	)
	server := httptest.NewServer(svc)
	defer server.Close()

	// test resource not found
	resp, err := server.Client().Post(server.URL+"/invalidres/bar", "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// test operation not found
	resp, err = server.Client().Post(server.URL+"/Book/invalidop", "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	transport, err := client.NewHTTPJSONTransport(server.URL, client.WithHTTPJSONTransportClient(server.Client()))
	require.NoError(t, err)
	client, err := client.NewClient(client.WithTransport(transport))
	require.NoError(t, err)

	ctx := context.TODO()
	input := &BookCreateInput{
		Name: "foo",
	}
	output, err := client.Execute(ctx, opCreate, input)
	require.NoError(t, err)
	// check output is of type BookCreateOutput
	require.IsType(t, &BookCreateOutput{}, output)
	// check output is correct
	require.Equal(t, "123", output.(*BookCreateOutput).ID)

	// check inval1d input
	input = &BookCreateInput{
		Name: "inval1d",
	}
	_, err = client.Execute(ctx, opCreate, input)
	require.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	// check output is of type BookCreationProblem
	var creationProblem *BookCreationProblem
	require.True(t, errors.As(err, &creationProblem))
	// check output is correct
	require.Equal(t, "invalid name", creationProblem.Message)
}

// TestServiceRoutesIgnoringPathPrefix verifies the same server routes correctly
// whether the client mounts it at the root, /api, or /admin/api. routing is on
// the last two path segments, so any mount prefix is transparent.
func TestServiceRoutesIgnoringPathPrefix(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	resource := pk.FindResource("Book")
	opCreate := resource.FindOperation("Create")

	svc := NewService(
		WithServicePackage(pk),
		WithOperationHandler(OperationHandlerFunc(func(ctx context.Context, op *OperationExecution) (client.Struct, error) {
			return &BookCreateOutput{ID: "123"}, nil
		})),
	)
	server := httptest.NewServer(svc)
	defer server.Close()

	for _, prefix := range []string{"", "/api", "/admin/api"} {
		t.Run("prefix "+prefix, func(t *testing.T) {
			transport, err := client.NewHTTPJSONTransport(server.URL+prefix, client.WithHTTPJSONTransportClient(server.Client()))
			require.NoError(t, err)
			cl, err := client.NewClient(client.WithTransport(transport))
			require.NoError(t, err)
			output, err := cl.Execute(context.TODO(), opCreate, &BookCreateInput{Name: "foo"})
			require.NoError(t, err)
			require.IsType(t, &BookCreateOutput{}, output)
			require.Equal(t, "123", output.(*BookCreateOutput).ID)
		})
	}
}
