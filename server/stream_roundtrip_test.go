package server

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

// A part reaches the caller when it is sent, not when the next part starts
func TestAPartArrivesBeforeTheNextPartStarts(t *testing.T) {
	pk, op := streamTestPackage(t)
	release := make(chan struct{})
	defer close(release)
	svc := NewService(
		WithServicePackage(pk),
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			out := make(chan client.StreamEvent[client.Struct])
			go func() {
				defer close(out)
				select {
				case out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: "a"}}:
				case <-ctx.Done():
					return
				}
				// Nothing else for longer than the bound below
				select {
				case <-time.After(2 * time.Second):
				case <-release:
				case <-ctx.Done():
				}
			}()
			return out, nil
		})),
	)
	ts := httptest.NewServer(svc)
	defer ts.Close()

	transport, err := client.NewHTTPJSONTransport(ts.URL, client.WithHTTPJSONTransportClient(ts.Client()))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	events, err := transport.Stream(ctx, &client.Request{Operation: op, Input: &BookCreateInput{}})
	require.NoError(t, err)

	select {
	case ev, ok := <-events:
		require.True(t, ok, "stream closed before the first part")
		require.NoError(t, ev.Err)
		require.Equal(t, "a", ev.Output.(*BookCreateOutput).ID)
		t.Logf("first part after %v", time.Since(start))
	case <-time.After(time.Second):
		// The handler holds this part for two seconds, so a part that waits
		// for the next one can never arrive inside this bound.
		t.Fatalf("first part not received within 1s")
	}
}

// postStream sends what a client sends, without the client, so a test can read
// the wire itself.
func postStream(t *testing.T, ts *httptest.Server, op *client.Operation) *http.Response {
	t.Helper()
	input := &BookCreateInput{}
	uri := ts.URL + "/" + op.Resource().PackageUniqueName() + "/" + op.Name()
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader("{}"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", input.StructPath().MIMENameJSONHTTP())
	res, err := ts.Client().Do(req)
	require.NoError(t, err)
	return res
}

// The delimiter that ends a part is written with the part, not when the next
// part starts. A multipart reader can only finish a part once it sees that
// delimiter, so a reader that finishes a part must not wait for the handler to
// send the next one. This is the framing the client depends on: the delimiter
// used to be written when the next part started, and every part then waited
// for the next part or the five second heartbeat.
func TestAStandardReaderEndsAPartBeforeTheNextOneStarts(t *testing.T) {
	pk, op := streamTestPackage(t)
	var once sync.Once
	release := make(chan struct{})
	released := func() { once.Do(func() { close(release) }) }
	defer released()

	svc := NewService(
		WithServicePackage(pk),
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			out := make(chan client.StreamEvent[client.Struct])
			go func() {
				defer close(out)
				select {
				case out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: "a"}}:
				case <-ctx.Done():
					return
				}
				// Longer than the bound below, so a part whose delimiter waits
				// for the next part cannot pass, and under the five second
				// heartbeat, so a heartbeat cannot end the part either.
				select {
				case <-time.After(2 * time.Second):
				case <-release:
				case <-ctx.Done():
					return
				}
				select {
				case out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: "b"}}:
				case <-ctx.Done():
				}
			}()
			return out, nil
		})),
	)
	ts := httptest.NewServer(svc)
	defer ts.Close()

	res := postStream(t, ts, op)
	defer res.Body.Close()

	mediaType, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	require.NoError(t, err)
	require.Equal(t, "multipart/mixed", mediaType)

	reader := multipart.NewReader(res.Body, params["boundary"])
	started := time.Now()
	part, err := reader.NextPart()
	require.NoError(t, err)
	require.Equal(t, (&BookCreateOutput{}).StructPath().MIMENameJSONHTTP(), part.Header.Get("Content-Type"))
	content, err := io.ReadAll(part)
	require.NoError(t, err)
	took := time.Since(started)

	require.Equal(t, `{"id":"a"}`, strings.TrimSpace(string(content)))
	require.Less(t, took, 1500*time.Millisecond,
		"the first part took %v, so its delimiter arrived with the next part", took)

	released()
	second, err := reader.NextPart()
	require.NoError(t, err)
	secondContent, err := io.ReadAll(second)
	require.NoError(t, err)
	require.Equal(t, `{"id":"b"}`, strings.TrimSpace(string(secondContent)))
}

// A stream that sends nothing still ends with the part that carries no
// Content-Type, which is what the specification says ends a stream. The
// closing delimiter alone would be a different wire.
func TestAnEmptyStreamEndsWithThePartThatCarriesNoContentType(t *testing.T) {
	pk, op := streamTestPackage(t)
	svc := NewService(
		WithServicePackage(pk),
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			out := make(chan client.StreamEvent[client.Struct])
			close(out)
			return out, nil
		})),
	)
	ts := httptest.NewServer(svc)
	defer ts.Close()

	res := postStream(t, ts, op)
	defer res.Body.Close()

	mediaType, params, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	require.NoError(t, err)
	require.Equal(t, "multipart/mixed", mediaType)

	reader := multipart.NewReader(res.Body, params["boundary"])
	marker, err := reader.NextPart()
	require.NoError(t, err, "a stream ends with a part, not with the closing delimiter alone")
	require.Empty(t, marker.Header.Get("Content-Type"))
	content, err := io.ReadAll(marker)
	require.NoError(t, err)
	require.Empty(t, content)

	_, err = reader.NextPart()
	require.Equal(t, io.EOF, err)
}
