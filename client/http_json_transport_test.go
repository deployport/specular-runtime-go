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

// TypeFQTN returns the fully qualified domain name of the type
func (t *BookCreateInput) TypeFQTN() TypeFQTN {
	return NewTypeFQTN("deployport/test", "BookCreateInput")
}
func (t *BookCreateInput) Hydrate(ctx *HydratationContext) error {
	t.Name = ctx.Content().GetProperty("name").(string)
	return nil
}

func (t *BookCreateInput) Dehydrate(ctx *DehydrationContext) error {
	ctx.Content().SetStruct(t.TypeFQTN().String())
	ctx.Content().SetProperty("name", t.Name)
	return nil
}

type BookCreateOutput struct {
	ID string `json:"id"`
}

// TypeFQTN returns the fully qualified domain name of the type
func (t *BookCreateOutput) TypeFQTN() TypeFQTN {
	return NewTypeFQTN("deployport/test", "BookCreateOutput")
}

func (t *BookCreateOutput) Hydrate(ctx *HydratationContext) error {
	t.ID = ctx.Content().GetProperty("id").(string)
	return nil
}

func (t *BookCreateOutput) Dehydrate(ctx *DehydrationContext) error {
	ctx.Content().SetStruct(t.TypeFQTN().String())
	ctx.Content().SetProperty("id", t.ID)
	return nil
}

type BookCreationProblem struct {
	Message string `json:"message"`
}

func (e *BookCreationProblem) Error() string {
	return e.Message
}

// TypeFQTN returns the fully qualified domain name of the type
func (e *BookCreationProblem) TypeFQTN() TypeFQTN {
	return NewTypeFQTN("deployport/test", "BookCreationProblem")
}

// Hydrate hydrates the type from the content
func (e *BookCreationProblem) Hydrate(ctx *HydratationContext) error {
	e.Message = ctx.Content().GetProperty("message").(string)
	return nil
}

// Dehydrate dehydrates the type into the content
func (e *BookCreationProblem) Dehydrate(ctx *DehydrationContext) error {
	ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// Is returns true if the error is of the same type
func (e *BookCreationProblem) Is(err error) bool {
	_, ok := err.(*BookCreationProblem)
	return ok
}

func testPackage() (*Package, error) {
	pk := NewPackage("deployport/test")
	input, err := pk.NewType("BookCreateInput", TypeBuilder(func() Struct {
		return &BookCreateInput{}
	}))
	if err != nil {
		return nil, err
	}
	output, err := pk.NewType("BookCreateOutput", TypeBuilder(func() Struct {
		return &BookCreateOutput{}
	}))
	if err != nil {
		return nil, err
	}
	prob, err := pk.NewType("BookCreationProblem", TypeBuilder(func() Struct {
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
	out := NewContent()
	out.SetStruct("deployport/test:BookCreateOutput")
	out.SetProperty("id", "id123")
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
	out := NewContent()
	out.SetStruct("deployport/test:BookCreationProblem")
	out.SetProperty("message", "too bad")
	tst := httpTransportTest(t, httpTransportTestParams{
		Output: out,
	})
	res := tst.Response
	err := tst.ResponseError
	require.NotNil(t, err)
	require.Nil(t, res)
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
	Output Content
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
		w.Header().Set("Content-Type", "specular/struct")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}))
	defer ts.Close()

	transport, err := NewHTTPJSONTransport(ts.URL)
	require.NoError(t, err)
	rs := pk.FindResource("Book")
	op := rs.FindOperation("Create")
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
