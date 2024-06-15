package server

import "go.deployport.com/specular-runtime/client"

type BookCreateInput struct {
	Name string `json:"name"`
}

// StructPath returns the struct path of the struct
func (e *BookCreateInput) StructPath() client.StructPath {
	return *client.NewStructPath(*testModulePath, "bookcreateinput")
}

func (e *BookCreateInput) Hydrate(ctx *client.HydratationContext) error {
	e.Name = ctx.Content().GetProperty("name").(string)
	return nil
}

func (e *BookCreateInput) Dehydrate(ctx *client.DehydrationContext) error {
	ctx.Content().SetProperty("name", e.Name)
	return nil
}

type BookCreateOutput struct {
	ID string `json:"id"`
}

// StructPath returns the struct path of the struct
func (e *BookCreateOutput) StructPath() client.StructPath {
	return *client.NewStructPath(*testModulePath, "bookcreateoutput")
}

func (e *BookCreateOutput) Hydrate(ctx *client.HydratationContext) error {
	e.ID = ctx.Content().GetProperty("id").(string)
	return nil
}

func (e *BookCreateOutput) Dehydrate(ctx *client.DehydrationContext) error {
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
func (e *BookCreationProblem) StructPath() client.StructPath {
	return *client.NewStructPath(*testModulePath, "bookcreationproblem")
}

// Hydrate hydrates the type from the content
func (e *BookCreationProblem) Hydrate(ctx *client.HydratationContext) error {
	e.Message = ctx.Content().GetProperty("message").(string)
	return nil
}

// Dehydrate dehydrates the type into the content
func (e *BookCreationProblem) Dehydrate(ctx *client.DehydrationContext) error {
	// ctx.Content().SetStruct(e.TypeFQTN().String())
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// Is returns true if the error is of the same type
func (e *BookCreationProblem) Is(err error) bool {
	_, ok := err.(*BookCreationProblem)
	return ok
}

var testModulePath = client.ModulePathFromTrustedValues("deployport", "test")

func testPackage() (*client.Package, error) {
	pk := client.NewPackage(testModulePath)
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
