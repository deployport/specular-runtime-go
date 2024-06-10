package server

import "go.deployport.com/specular-runtime/client"

type BookCreateInput struct {
	Name string `json:"name"`
}

// TypeFQTN returns the fully qualified domain name of the type
func (t *BookCreateInput) TypeFQTN() client.TypeFQTN {
	return client.NewTypeFQTN("deployport/test", "BookCreateInput")
}
func (t *BookCreateInput) Hydrate(ctx *client.HydratationContext) error {
	t.Name = ctx.Content().GetProperty("name").(string)
	return nil
}

func (t *BookCreateInput) Dehydrate(ctx *client.DehydrationContext) error {
	ctx.Content().SetStruct(t.TypeFQTN().String())
	ctx.Content().SetProperty("name", t.Name)
	return nil
}

type BookCreateOutput struct {
	ID string `json:"id"`
}

// TypeFQTN returns the fully qualified domain name of the type
func (t *BookCreateOutput) TypeFQTN() client.TypeFQTN {
	return client.NewTypeFQTN("deployport/test", "BookCreateOutput")
}

func (t *BookCreateOutput) Hydrate(ctx *client.HydratationContext) error {
	t.ID = ctx.Content().GetProperty("id").(string)
	return nil
}

func (t *BookCreateOutput) Dehydrate(ctx *client.DehydrationContext) error {
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
func (e *BookCreationProblem) TypeFQTN() client.TypeFQTN {
	return client.NewTypeFQTN("deployport/test", "BookCreationProblem")
}

// Hydrate hydrates the type from the content
func (e *BookCreationProblem) Hydrate(ctx *client.HydratationContext) error {
	e.Message = ctx.Content().GetProperty("message").(string)
	return nil
}

// Dehydrate dehydrates the type into the content
func (e *BookCreationProblem) Dehydrate(ctx *client.DehydrationContext) error {
	ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// Is returns true if the error is of the same type
func (e *BookCreationProblem) Is(err error) bool {
	_, ok := err.(*BookCreationProblem)
	return ok
}

func testPackage() (*client.Package, error) {
	pk := client.NewPackage(client.ModulePathFromTrustedValues("deployport", "test"))
	input, err := pk.NewType("BookCreateInput", client.TypeBuilder(func() client.Struct {
		return &BookCreateInput{}
	}))
	if err != nil {
		return nil, err
	}
	output, err := pk.NewType("BookCreateOutput", client.TypeBuilder(func() client.Struct {
		return &BookCreateOutput{}
	}))
	if err != nil {
		return nil, err
	}
	prob, err := pk.NewType("BookCreationProblem", client.TypeBuilder(func() client.Struct {
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
