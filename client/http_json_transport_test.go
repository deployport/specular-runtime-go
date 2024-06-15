package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type BookCreateInput struct {
	Name string `json:"name"`
}

// StructPath returns the struct path of the struct
func (e *BookCreateInput) StructPath() StructPath {
	return *NewStructPath(*testPackagePath, "bookcreateinput")
}

func (e *BookCreateInput) Hydrate(ctx *HydratationContext) error {
	e.Name = ctx.Content().GetProperty("name").(string)
	return nil
}

func (e *BookCreateInput) Dehydrate(ctx *DehydrationContext) error {
	// ctx.Content().SetStruct(t.TypeFQTN().String())
	ctx.Content().SetProperty("name", e.Name)
	return nil
}

type BookCreateOutput struct {
	ID string `json:"id"`
}

// StructPath returns the struct path of the struct
func (e *BookCreateOutput) StructPath() StructPath {
	return *NewStructPath(*testPackagePath, "bookcreateoutput")
}

func (e *BookCreateOutput) Hydrate(ctx *HydratationContext) error {
	e.ID = ctx.Content().GetProperty("id").(string)
	return nil
}

func (e *BookCreateOutput) Dehydrate(ctx *DehydrationContext) error {
	// ctx.Content().SetStruct(t.TypeFQTN().String())
	ctx.Content().SetProperty("id", e.ID)
	return nil
}

type BookCreationProblem struct {
	Message string `json:"message"`
}

func (e *BookCreationProblem) Error() string {
	return e.Message
}

// StructPath returns the struct path of the struct
func (e *BookCreationProblem) StructPath() StructPath {
	return *NewStructPath(*testPackagePath, "bookcreationproblem")
}

// Hydrate hydrates the type from the content
func (e *BookCreationProblem) Hydrate(ctx *HydratationContext) error {
	e.Message = ctx.Content().GetProperty("message").(string)
	return nil
}

// Dehydrate dehydrates the type into the content
func (e *BookCreationProblem) Dehydrate(ctx *DehydrationContext) error {
	// ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// Is returns true if the error is of the same type
func (e *BookCreationProblem) Is(err error) bool {
	_, ok := err.(*BookCreationProblem)
	return ok
}

var testPackagePath = ModulePathFromTrustedValues("deployport", "test")

func testPackage() (*Package, error) {
	pk := NewPackage(testPackagePath)
	input, err := pk.NewType("bookcreateinput", TypeBuilder(func() Struct {
		return &BookCreateInput{}
	}))
	if err != nil {
		return nil, err
	}
	output, err := pk.NewType("bookcreateoutput", TypeBuilder(func() Struct {
		return &BookCreateOutput{}
	}))
	if err != nil {
		return nil, err
	}
	prob, err := pk.NewType("bookcreationproblem", TypeBuilder(func() Struct {
		return &BookCreationProblem{}
	}))
	if err != nil {
		return nil, err
	}
	rs, err := pk.NewResource("Book")
	if err != nil {
		return nil, err
	}
	op, err := rs.NewOperation("Create")
	if err != nil {
		return nil, err
	}
	op.SetInput(input)
	op.SetOutput(output)
	op.RegisterProblemType(prob)
	return pk, nil
}

func TestHTTPJSONTransportSuccessResponse(t *testing.T) {
	out := &BookCreateOutput{}
	out.ID = "id123"
	tst := httpTransportTest(t, httpTransportTestParams{
		Output: out,
	})
	res := tst.Response
	err := tst.ResponseError
	require.NoError(t, err)
	require.NotNil(t, res)
	output := res.(*BookCreateOutput)
	require.Equal(t, "id123", output.ID)
}

func TestHTTPJSONTransportProblem(t *testing.T) {
	out := &BookCreationProblem{}
	out.Message = "too bad"
	tst := httpTransportTest(t, httpTransportTestParams{
		Output: out,
	})
	res := tst.Response
	err := tst.ResponseError
	require.NotNil(t, err)
	require.Nil(t, res)
	require.ErrorIs(t, err, &BookCreationProblem{})
	problem := err.(*BookCreationProblem)
	require.Equal(t, "too bad", problem.Message)
}

func TestWithEndpointURL(t *testing.T) {
	transport := newNoServerTestTransport(t, WithEndpointURL("http://override:3020"))
	require.Equal(t, "http://override:3020", transport.EndpointURL)
}

func newNoServerTestTransport(t *testing.T, options ...Option) *HTTPJSONTransport {
	tr, err := NewHTTPJSONTransport("http://localhost:8080", options...)
	require.NoError(t, err)
	return tr
}

type httpTransportTestParams struct {
	Output Struct
}

func httpTransportTest(t *testing.T, params httpTransportTestParams) *httpTransportTestState {
	ctx := context.TODO()
	pk, err := testPackage()
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := params.Output
		b, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", res.StructPath().MIMENameJSONHTTP())
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}))
	defer ts.Close()

	transport, err := NewHTTPJSONTransport(ts.URL)
	require.NoError(t, err)
	rs := pk.FindResource("Book")
	require.NotNil(t, rs)
	op := rs.FindOperation("Create")
	require.NotNil(t, op)
	res, err := transport.Execute(ctx, &Request{
		Operation: op,
		Input:     &BookCreateInput{},
	})
	return &httpTransportTestState{
		Response:      res,
		ResponseError: err,
	}
}

type httpTransportTestState struct {
	Response      Struct
	ResponseError error
}
