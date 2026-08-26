package server

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client"
)

// testPackage with Book.Create marked streamed
func streamTestPackage(t *testing.T) (*client.Package, *client.Operation) {
	t.Helper()
	pk, err := testPackage()
	require.NoError(t, err)
	op := pk.FindResource("Book").FindOperation("Create")
	op.SetStreamed()
	return pk, op
}

// streamOf serves ids one at a time, waiting gap before each
func streamOf(t *testing.T, pk *client.Package, gap time.Duration, ids ...string) *Service {
	t.Helper()
	return NewService(
		WithServicePackage(pk),
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			out := make(chan client.StreamEvent[client.Struct])
			go func() {
				defer close(out)
				for _, id := range ids {
					select {
					case <-time.After(gap):
					case <-ctx.Done():
						return
					}
					select {
					case out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: id}}:
					case <-ctx.Done():
						return
					}
				}
			}()
			return out, nil
		})),
	)
}

// drain returns every id the client received
func drain(t *testing.T, ts *httptest.Server, op *client.Operation) []string {
	t.Helper()
	transport, err := client.NewHTTPJSONTransport(ts.URL, client.WithHTTPJSONTransportClient(ts.Client()))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	events, err := transport.Stream(ctx, &client.Request{Operation: op, Input: &BookCreateInput{}})
	require.NoError(t, err)

	var got []string
	for ev := range events {
		require.NoError(t, ev.Err)
		got = append(got, ev.Output.(*BookCreateOutput).ID)
	}
	return got
}

// The real server runtime and client transport over multipart/mixed
func TestAStreamedOperationDeliversItsParts(t *testing.T) {
	pk, op := streamTestPackage(t)
	ts := httptest.NewServer(streamOf(t, pk, 0, "a", "b", "c"))
	defer ts.Close()

	require.Equal(t, []string{"a", "b", "c"}, drain(t, ts, op))
}

// Server.WriteTimeout would otherwise cap the stream
func TestAStreamOutlivesTheServerWriteTimeout(t *testing.T) {
	pk, op := streamTestPackage(t)
	ts := httptest.NewUnstartedServer(streamOf(t, pk, 150*time.Millisecond, "a", "b", "c"))
	// Well under the time the handler needs
	ts.Config.WriteTimeout = 200 * time.Millisecond
	ts.Start()
	defer ts.Close()

	require.Equal(t, []string{"a", "b", "c"}, drain(t, ts, op))
}
