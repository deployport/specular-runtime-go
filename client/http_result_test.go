package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// BookThrottledProblem stands for a declared error annotated
// @specular:Status(429). The generator emits HTTPStatus for such a struct and
// for no other.
type BookThrottledProblem struct {
	Message string `json:"message"`
}

func (e *BookThrottledProblem) Error() string {
	return e.Message
}

// StructPath returns the struct path of the struct
func (e *BookThrottledProblem) StructPath() StructPath {
	return *NewStructPath(*testPackagePath, "bookthrottledproblem")
}

// HTTPStatus is what the generator emits for an annotated declared error.
func (e *BookThrottledProblem) HTTPStatus() int {
	return http.StatusTooManyRequests
}

// A declared error that was never annotated must keep answering 400. This is
// the compatibility guarantee of the whole annotation: every error shipped
// before it existed behaves exactly as it did.
func TestHTTPStatusCodeUnannotatedDeclaredErrorIs400(t *testing.T) {
	result := HTTPResultForErrorStruct(&BookCreationProblem{Message: "nope"})
	require.Equal(t, http.StatusBadRequest, result.HTTPStatusCode())
}

// An annotated declared error answers the status it declared.
func TestHTTPStatusCodeAnnotatedDeclaredErrorCarriesItsOwnStatus(t *testing.T) {
	result := HTTPResultForErrorStruct(&BookThrottledProblem{Message: "slow down"})
	require.Equal(t, http.StatusTooManyRequests, result.HTTPStatusCode())
}

// A built-in error wins over both. Its status comes from the reserved table and
// nothing generated may override it.
func TestHTTPStatusCodeBuiltinErrorWinsOverCarrier(t *testing.T) {
	builtin := NewError()
	builtin.Message = "internal"
	builtin.HTTPStatusCode = http.StatusInternalServerError
	result := HTTPResultForError(builtin)
	require.Equal(t, http.StatusInternalServerError, result.HTTPStatusCode())
}

// A regular output keeps 200, and it keeps it even when the struct carries a
// status, because a 4xx may never accompany an envelope that is not kind=error.
func TestHTTPStatusCodeRegularOutputIs200(t *testing.T) {
	require.Equal(t, http.StatusOK, HTTPResultForStruct(&BookCreateOutput{ID: "1"}).HTTPStatusCode())
	require.Equal(t, http.StatusOK, HTTPResultForStruct(&BookThrottledProblem{Message: "slow down"}).HTTPStatusCode())
}

// The status reaches the wire, and the envelope still marks the error. The
// status is advisory; the content type plus kind stays authoritative.
func TestWriteResponseSendsDeclaredStatusAndErrorEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	result := HTTPResultForErrorStruct(&BookThrottledProblem{Message: "slow down"})
	require.NoError(t, result.WriteResponse(rec))

	res := rec.Result()
	defer func() { _ = res.Body.Close() }()
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "kind=error")
	require.Contains(t, res.Header.Get("Content-Type"), "bookthrottledproblem")
}

// An unannotated declared error reaches the wire as 400 with the same envelope.
func TestWriteResponseSendsBadRequestForUnannotatedError(t *testing.T) {
	rec := httptest.NewRecorder()
	result := HTTPResultForErrorStruct(&BookCreationProblem{Message: "nope"})
	require.NoError(t, result.WriteResponse(rec))

	res := rec.Result()
	defer func() { _ = res.Body.Close() }()
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "kind=error")
}

// The interface the generator targets must be satisfied by the method the
// generator emits. A rename on either side breaks this test rather than
// silently dropping every declared status back to 400.
func TestHTTPStatusCarrierIsSatisfiedByTheGeneratedMethod(t *testing.T) {
	var carrier HTTPStatusCarrier = &BookThrottledProblem{}
	require.Equal(t, http.StatusTooManyRequests, carrier.HTTPStatus())

	_, ok := Struct(&BookCreationProblem{}).(HTTPStatusCarrier)
	require.False(t, ok, "an unannotated error must not satisfy HTTPStatusCarrier")
}
