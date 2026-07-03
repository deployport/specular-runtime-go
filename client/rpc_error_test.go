package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// serveRawResponse starts a server that replies with a fixed content type and body,
// then executes the Book/Create operation against it and returns the result.
func serveRawResponse(t *testing.T, contentType, body string) (Struct, error) {
	ctx := context.TODO()
	pk, err := testPackage()
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(body))
	}))
	defer ts.Close()

	transport, err := NewHTTPJSONTransport(ts.URL)
	require.NoError(t, err)
	rs := pk.FindResource("Book")
	require.NotNil(t, rs)
	op := rs.FindOperation("Create")
	require.NotNil(t, op)
	return transport.Execute(ctx, &Request{
		Operation: op,
		Input:     &BookCreateInput{},
	})
}

// An error type the client was not generated with must surface as a generic,
// catchable error carrying the type name and raw payload, never a decode failure.
func TestHTTPJSONTransportUnknownErrorType(t *testing.T) {
	unknown := NewStructPath(*testPackagePath, "futureproblem")
	res, err := serveRawResponse(t,
		unknown.MIMENameJSONHTTP()+"; kind=error",
		`{"message":"quota exceeded","code":"QUOTA_EXCEEDED","limit":100}`,
	)
	require.Nil(t, res)
	require.Error(t, err)

	var ue *UnknownRPCError
	require.True(t, errors.As(err, &ue), "expected an *UnknownRPCError, got %T", err)
	require.Equal(t, "quota exceeded", ue.Error())
	require.Equal(t, "quota exceeded", ue.Message)
	require.Equal(t, "QUOTA_EXCEEDED", ue.Code)
	require.Equal(t, unknown.String(), ue.Path.String())
	require.JSONEq(t, `{"message":"quota exceeded","code":"QUOTA_EXCEEDED","limit":100}`, string(ue.Payload))

	// It must also satisfy the RPCError contract.
	var rpcErr RPCError
	require.True(t, errors.As(err, &rpcErr))
}

// Without the kind=error marker an unknown type stays a hard error: the client
// must not silently treat an unrecognized output as an error.
func TestHTTPJSONTransportUnknownTypeWithoutErrorMarker(t *testing.T) {
	unknown := NewStructPath(*testPackagePath, "futureoutput")
	res, err := serveRawResponse(t,
		unknown.MIMENameJSONHTTP(),
		`{"id":"123"}`,
	)
	require.Nil(t, res)
	require.Error(t, err)
	require.True(t, IsTypeNotFoundError(err), "expected TypeNotFoundError, got %T", err)

	var ue *UnknownRPCError
	require.False(t, errors.As(err, &ue))
}

// The server marks user-defined error envelopes with kind=error and HTTP 400,
// while regular outputs stay unmarked at 200.
func TestHTTPResultErrorEnvelopeMarking(t *testing.T) {
	errRes := HTTPResultForErrorStruct(&BookCreationProblem{Message: "bad"})
	require.Equal(t, http.StatusBadRequest, errRes.HTTPStatusCode())
	require.Contains(t, errRes.MimeType(), "; kind=error")

	okRes := HTTPResultForStruct(&BookCreateOutput{ID: "1"})
	require.Equal(t, http.StatusOK, okRes.HTTPStatusCode())
	require.NotContains(t, okRes.MimeType(), "kind=error")

	// Built-in errors keep their own status and are also marked as errors.
	builtin := HTTPResultForError(&Error{HTTPStatusCode: http.StatusForbidden, ErrorCode: CallErrorCodeAccessDenied})
	require.Equal(t, http.StatusForbidden, builtin.HTTPStatusCode())
	require.Contains(t, builtin.MimeType(), "; kind=error")
}
