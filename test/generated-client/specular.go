package testpackage

import (
	"context"
	"errors"
	"time"

	clientruntime "go.deployport.com/specular-runtime/client"
)

// NewBody creates a new Body
func NewBody() *Body {
	s := &Body{}
	return s
}

// Body entity
type Body struct {
	BodyType                     BodyType   `json:"bodyType,omitempty"`
	BodyTypeNullable             *BodyType  `json:"bodyTypeNullable,omitempty"`
	ContentLengthFloat32         float32    `json:"contentLengthFloat32,omitempty"`
	ContentLengthFloat64         float64    `json:"contentLengthFloat64,omitempty"`
	ContentLengthFloat64Nullable *float64   `json:"contentLengthFloat64Nullable,omitempty"`
	ContentLengthInt32           int32      `json:"contentLengthInt32,omitempty"`
	ContentLengthInt64           int64      `json:"contentLengthInt64,omitempty"`
	ContentLengthUint32          uint32     `json:"contentLengthUint32,omitempty"`
	ContentLengthUint64          uint64     `json:"contentLengthUint64,omitempty"`
	ContentLengthUint64Nullable  uint64     `json:"contentLengthUint64Nullable,omitempty"`
	CreatedAt                    time.Time  `json:"createdAt,omitempty"`
	CreatedAtNullable            *time.Time `json:"createdAtNullable,omitempty"`
	FileData                     []byte     `json:"fileData,omitempty"`
	FileDataNullable             []byte     `json:"fileDataNullable,omitempty"`
	MessageString                string     `json:"messageString,omitempty"`
	MessageStringNullable        *string    `json:"messageStringNullable,omitempty"`
}

// GetBodyType returns the value for the field bodyType
func (e *Body) GetBodyType() BodyType {
	return e.BodyType
}

// SetBodyType sets the value for the field bodyType
func (e *Body) SetBodyType(bodyType BodyType) {
	e.BodyType = bodyType
}

// GetBodyTypeNullable returns the value for the field bodyTypeNullable
func (e *Body) GetBodyTypeNullable() *BodyType {
	return e.BodyTypeNullable
}

// SetBodyTypeNullable sets the value for the field bodyTypeNullable
func (e *Body) SetBodyTypeNullable(bodyTypeNullable *BodyType) {
	e.BodyTypeNullable = bodyTypeNullable
}

// GetContentLengthFloat32 returns the value for the field contentLengthFloat32
func (e *Body) GetContentLengthFloat32() float32 {
	return e.ContentLengthFloat32
}

// SetContentLengthFloat32 sets the value for the field contentLengthFloat32
func (e *Body) SetContentLengthFloat32(contentLengthFloat32 float32) {
	e.ContentLengthFloat32 = contentLengthFloat32
}

// GetContentLengthFloat64 returns the value for the field contentLengthFloat64
func (e *Body) GetContentLengthFloat64() float64 {
	return e.ContentLengthFloat64
}

// SetContentLengthFloat64 sets the value for the field contentLengthFloat64
func (e *Body) SetContentLengthFloat64(contentLengthFloat64 float64) {
	e.ContentLengthFloat64 = contentLengthFloat64
}

// GetContentLengthFloat64Nullable returns the value for the field contentLengthFloat64Nullable
func (e *Body) GetContentLengthFloat64Nullable() *float64 {
	return e.ContentLengthFloat64Nullable
}

// SetContentLengthFloat64Nullable sets the value for the field contentLengthFloat64Nullable
func (e *Body) SetContentLengthFloat64Nullable(contentLengthFloat64Nullable *float64) {
	e.ContentLengthFloat64Nullable = contentLengthFloat64Nullable
}

// GetContentLengthInt32 returns the value for the field contentLengthInt32
func (e *Body) GetContentLengthInt32() int32 {
	return e.ContentLengthInt32
}

// SetContentLengthInt32 sets the value for the field contentLengthInt32
func (e *Body) SetContentLengthInt32(contentLengthInt32 int32) {
	e.ContentLengthInt32 = contentLengthInt32
}

// GetContentLengthInt64 returns the value for the field contentLengthInt64
func (e *Body) GetContentLengthInt64() int64 {
	return e.ContentLengthInt64
}

// SetContentLengthInt64 sets the value for the field contentLengthInt64
func (e *Body) SetContentLengthInt64(contentLengthInt64 int64) {
	e.ContentLengthInt64 = contentLengthInt64
}

// GetContentLengthUint32 returns the value for the field contentLengthUint32
func (e *Body) GetContentLengthUint32() uint32 {
	return e.ContentLengthUint32
}

// SetContentLengthUint32 sets the value for the field contentLengthUint32
func (e *Body) SetContentLengthUint32(contentLengthUint32 uint32) {
	e.ContentLengthUint32 = contentLengthUint32
}

// GetContentLengthUint64 returns the value for the field contentLengthUint64
func (e *Body) GetContentLengthUint64() uint64 {
	return e.ContentLengthUint64
}

// SetContentLengthUint64 sets the value for the field contentLengthUint64
func (e *Body) SetContentLengthUint64(contentLengthUint64 uint64) {
	e.ContentLengthUint64 = contentLengthUint64
}

// GetContentLengthUint64Nullable returns the value for the field contentLengthUint64Nullable
func (e *Body) GetContentLengthUint64Nullable() uint64 {
	return e.ContentLengthUint64Nullable
}

// SetContentLengthUint64Nullable sets the value for the field contentLengthUint64Nullable
func (e *Body) SetContentLengthUint64Nullable(contentLengthUint64Nullable uint64) {
	e.ContentLengthUint64Nullable = contentLengthUint64Nullable
}

// GetCreatedAt returns the value for the field createdAt
func (e *Body) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// SetCreatedAt sets the value for the field createdAt
func (e *Body) SetCreatedAt(createdAt time.Time) {
	e.CreatedAt = createdAt
}

// GetCreatedAtNullable returns the value for the field createdAtNullable
func (e *Body) GetCreatedAtNullable() *time.Time {
	return e.CreatedAtNullable
}

// SetCreatedAtNullable sets the value for the field createdAtNullable
func (e *Body) SetCreatedAtNullable(createdAtNullable *time.Time) {
	e.CreatedAtNullable = createdAtNullable
}

// GetFileData returns the value for the field fileData
func (e *Body) GetFileData() []byte {
	return e.FileData
}

// SetFileData sets the value for the field fileData
func (e *Body) SetFileData(fileData []byte) {
	e.FileData = fileData
}

// GetFileDataNullable returns the value for the field fileDataNullable
func (e *Body) GetFileDataNullable() []byte {
	return e.FileDataNullable
}

// SetFileDataNullable sets the value for the field fileDataNullable
func (e *Body) SetFileDataNullable(fileDataNullable []byte) {
	e.FileDataNullable = fileDataNullable
}

// GetMessageString returns the value for the field messageString
func (e *Body) GetMessageString() string {
	return e.MessageString
}

// SetMessageString sets the value for the field messageString
func (e *Body) SetMessageString(messageString string) {
	e.MessageString = messageString
}

// GetMessageStringNullable returns the value for the field messageStringNullable
func (e *Body) GetMessageStringNullable() *string {
	return e.MessageStringNullable
}

// SetMessageStringNullable sets the value for the field messageStringNullable
func (e *Body) SetMessageStringNullable(messageStringNullable *string) {
	e.MessageStringNullable = messageStringNullable
}

// Hydrate implements struct hydrate
func (e *Body) Hydrate(ctx *clientruntime.HydratationContext) error {
	if err := clientruntime.ContentRequireStringProperty(ctx.Content(), "bodyType", &e.BodyType); err != nil {
		return err
	}
	if err := clientruntime.ContentOptionalStringProperty(ctx.Content(), "bodyTypeNullable", &e.BodyTypeNullable); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthFloat32", &e.ContentLengthFloat32); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthFloat64", &e.ContentLengthFloat64); err != nil {
		return err
	}
	if err := clientruntime.ContentOptionalNumericProperty(ctx.Content(), "contentLengthFloat64Nullable", &e.ContentLengthFloat64Nullable); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthInt32", &e.ContentLengthInt32); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthInt64", &e.ContentLengthInt64); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthUint32", &e.ContentLengthUint32); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthUint64", &e.ContentLengthUint64); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "contentLengthUint64Nullable", &e.ContentLengthUint64Nullable); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireTimeProperty(ctx.Content(), "createdAt", &e.CreatedAt); err != nil {
		return err
	}
	if err := clientruntime.ContentOptionalTimeProperty(ctx.Content(), "createdAtNullable", &e.CreatedAtNullable); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireBinaryProperty(ctx.Content(), "fileData", &e.FileData); err != nil {
		return err
	}
	if err := clientruntime.ContentOptionalBinaryProperty(ctx.Content(), "fileDataNullable", &e.FileDataNullable); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireStringProperty(ctx.Content(), "messageString", &e.MessageString); err != nil {
		return err
	}
	if err := clientruntime.ContentOptionalStringProperty(ctx.Content(), "messageStringNullable", &e.MessageStringNullable); err != nil {
		return err
	}
	return nil
}

// Dehydrate implements struct dehydrate
func (e *Body) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	ctx.Content().SetProperty("bodyType", e.BodyType)
	ctx.Content().SetProperty("bodyTypeNullable", e.BodyTypeNullable)
	ctx.Content().SetProperty("contentLengthFloat32", e.ContentLengthFloat32)
	ctx.Content().SetProperty("contentLengthFloat64", e.ContentLengthFloat64)
	ctx.Content().SetProperty("contentLengthFloat64Nullable", e.ContentLengthFloat64Nullable)
	ctx.Content().SetProperty("contentLengthInt32", e.ContentLengthInt32)
	ctx.Content().SetProperty("contentLengthInt64", e.ContentLengthInt64)
	ctx.Content().SetProperty("contentLengthUint32", e.ContentLengthUint32)
	ctx.Content().SetProperty("contentLengthUint64", e.ContentLengthUint64)
	ctx.Content().SetProperty("contentLengthUint64Nullable", e.ContentLengthUint64Nullable)
	ctx.Content().SetProperty("createdAt", e.CreatedAt)
	ctx.Content().SetProperty("createdAtNullable", e.CreatedAtNullable)
	ctx.Content().SetProperty("fileData", e.FileData)
	ctx.Content().SetProperty("fileDataNullable", e.FileDataNullable)
	ctx.Content().SetProperty("messageString", e.MessageString)
	ctx.Content().SetProperty("messageStringNullable", e.MessageStringNullable)
	return nil
}

// StructPath returns StructPath
func (e *Body) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathBody.Path()
}

// NewResponse creates a new Response
func NewResponse() *Response {
	s := &Response{}
	return s
}

// Response entity
type Response struct {
	Body *Body `json:"body,omitempty"`
}

// GetBody returns the value for the field body
func (e *Response) GetBody() *Body {
	return e.Body
}

// SetBody sets the value for the field body
func (e *Response) SetBody(body *Body) {
	e.Body = body
}

// Hydrate implements struct hydrate
func (e *Response) Hydrate(ctx *clientruntime.HydratationContext) error {
	if err := clientruntime.ContentObjectProperty(ctx.Content(), SpecularMeta().Body().TypeBuilder(), ctx.Package(), "body", &e.Body); err != nil {
		return err
	}
	return nil
}

// Dehydrate implements struct dehydrate
func (e *Response) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {

	var fieldContentBody clientruntime.Content
	if v := e.Body; v != nil {
		fieldContentBody = clientruntime.NewContent()
		if err := v.Dehydrate(ctx.CopyWithContent(fieldContentBody)); err != nil {
			return clientruntime.NewPropertyDehydrationError("body", err)
		}
	}
	ctx.Content().SetProperty("body", fieldContentBody.Map())
	return nil
}

// StructPath returns StructPath
func (e *Response) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathResponse.Path()
}

// NewNotFoundProblem creates a new NotFoundProblem
func NewNotFoundProblem() *NotFoundProblem {
	s := &NotFoundProblem{}
	return s
}

// NotFoundProblem entity
type NotFoundProblem struct {
	Detail  string `json:"detail,omitempty"`
	Status  int32  `json:"status,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error implements the error interface
func (e *NotFoundProblem) Error() string {
	return e.GetMessage()
}

// Is indicates whether the given error chain contains an error of type [NotFoundProblem]
func (e *NotFoundProblem) Is(err error) bool {
	_, ok := err.(*NotFoundProblem)
	return ok
}

// IsNotFoundProblem indicates whether the given error chain contains an error of type [NotFoundProblem]
func IsNotFoundProblem(err error) bool {
	return errors.Is(err, &NotFoundProblem{})
}

// GetDetail returns the value for the field detail
func (e *NotFoundProblem) GetDetail() string {
	return e.Detail
}

// SetDetail sets the value for the field detail
func (e *NotFoundProblem) SetDetail(detail string) {
	e.Detail = detail
}

// GetStatus returns the value for the field status
func (e *NotFoundProblem) GetStatus() int32 {
	return e.Status
}

// SetStatus sets the value for the field status
func (e *NotFoundProblem) SetStatus(status int32) {
	e.Status = status
}

// GetTitle returns the value for the field title
func (e *NotFoundProblem) GetTitle() string {
	return e.Title
}

// SetTitle sets the value for the field title
func (e *NotFoundProblem) SetTitle(title string) {
	e.Title = title
}

// GetMessage returns the value for the field message
func (e *NotFoundProblem) GetMessage() string {
	return e.Message
}

// SetMessage sets the value for the field message
func (e *NotFoundProblem) SetMessage(message string) {
	e.Message = message
}

// Hydrate implements struct hydrate
func (e *NotFoundProblem) Hydrate(ctx *clientruntime.HydratationContext) error {
	if err := clientruntime.ContentRequireStringProperty(ctx.Content(), "detail", &e.Detail); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireNumericProperty(ctx.Content(), "status", &e.Status); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireStringProperty(ctx.Content(), "title", &e.Title); err != nil {
		return err
	}
	if err := clientruntime.ContentRequireStringProperty(ctx.Content(), "message", &e.Message); err != nil {
		return err
	}
	return nil
}

// Dehydrate implements struct dehydrate
func (e *NotFoundProblem) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	ctx.Content().SetProperty("detail", e.Detail)
	ctx.Content().SetProperty("status", e.Status)
	ctx.Content().SetProperty("title", e.Title)
	ctx.Content().SetProperty("message", e.Message)
	return nil
}

// StructPath returns StructPath
func (e *NotFoundProblem) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathNotFoundProblem.Path()
}

// NewTestHTTPGetInput creates a new TestHTTPGetInput
func NewTestHTTPGetInput() *TestHTTPGetInput {
	s := &TestHTTPGetInput{}
	return s
}

// TestHTTPGetInput entity
type TestHTTPGetInput struct {
}

// Hydrate implements struct hydrate
func (e *TestHTTPGetInput) Hydrate(ctx *clientruntime.HydratationContext) error {
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPGetInput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPGetInput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPGetInput.Path()
}

// NewTestHTTPGetOutput creates a new TestHTTPGetOutput
func NewTestHTTPGetOutput() *TestHTTPGetOutput {
	s := &TestHTTPGetOutput{}
	return s
}

// TestHTTPGetOutput entity
type TestHTTPGetOutput struct {
	Response *Response `json:"response,omitempty"`
}

// GetResponse returns the value for the field response
func (e *TestHTTPGetOutput) GetResponse() *Response {
	return e.Response
}

// SetResponse sets the value for the field response
func (e *TestHTTPGetOutput) SetResponse(response *Response) {
	e.Response = response
}

// Hydrate implements struct hydrate
func (e *TestHTTPGetOutput) Hydrate(ctx *clientruntime.HydratationContext) error {
	if err := clientruntime.ContentObjectProperty(ctx.Content(), SpecularMeta().Response().TypeBuilder(), ctx.Package(), "response", &e.Response); err != nil {
		return err
	}
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPGetOutput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {

	var fieldContentResponse clientruntime.Content
	if v := e.Response; v != nil {
		fieldContentResponse = clientruntime.NewContent()
		if err := v.Dehydrate(ctx.CopyWithContent(fieldContentResponse)); err != nil {
			return clientruntime.NewPropertyDehydrationError("response", err)
		}
	}
	ctx.Content().SetProperty("response", fieldContentResponse.Map())
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPGetOutput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPGetOutput.Path()
}

// NewTestHTTPOtherInput creates a new TestHTTPOtherInput
func NewTestHTTPOtherInput() *TestHTTPOtherInput {
	s := &TestHTTPOtherInput{}
	return s
}

// TestHTTPOtherInput entity
type TestHTTPOtherInput struct {
}

// Hydrate implements struct hydrate
func (e *TestHTTPOtherInput) Hydrate(ctx *clientruntime.HydratationContext) error {
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPOtherInput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPOtherInput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPOtherInput.Path()
}

// NewTestHTTPOtherOutput creates a new TestHTTPOtherOutput
func NewTestHTTPOtherOutput() *TestHTTPOtherOutput {
	s := &TestHTTPOtherOutput{}
	return s
}

// TestHTTPOtherOutput entity
type TestHTTPOtherOutput struct {
}

// Hydrate implements struct hydrate
func (e *TestHTTPOtherOutput) Hydrate(ctx *clientruntime.HydratationContext) error {
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPOtherOutput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPOtherOutput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPOtherOutput.Path()
}

// NewTestHTTPWatchChangesInput creates a new TestHTTPWatchChangesInput
func NewTestHTTPWatchChangesInput() *TestHTTPWatchChangesInput {
	s := &TestHTTPWatchChangesInput{}
	return s
}

// TestHTTPWatchChangesInput entity
type TestHTTPWatchChangesInput struct {
}

// Hydrate implements struct hydrate
func (e *TestHTTPWatchChangesInput) Hydrate(ctx *clientruntime.HydratationContext) error {
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPWatchChangesInput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPWatchChangesInput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPWatchChangesInput.Path()
}

// NewTestHTTPWatchChangesOutput creates a new TestHTTPWatchChangesOutput
func NewTestHTTPWatchChangesOutput() *TestHTTPWatchChangesOutput {
	s := &TestHTTPWatchChangesOutput{}
	return s
}

// TestHTTPWatchChangesOutput entity
type TestHTTPWatchChangesOutput struct {
	Response *Response `json:"response,omitempty"`
}

// GetResponse returns the value for the field response
func (e *TestHTTPWatchChangesOutput) GetResponse() *Response {
	return e.Response
}

// SetResponse sets the value for the field response
func (e *TestHTTPWatchChangesOutput) SetResponse(response *Response) {
	e.Response = response
}

// Hydrate implements struct hydrate
func (e *TestHTTPWatchChangesOutput) Hydrate(ctx *clientruntime.HydratationContext) error {
	if err := clientruntime.ContentObjectProperty(ctx.Content(), SpecularMeta().Response().TypeBuilder(), ctx.Package(), "response", &e.Response); err != nil {
		return err
	}
	return nil
}

// Dehydrate implements struct dehydrate
func (e *TestHTTPWatchChangesOutput) Dehydrate(ctx *clientruntime.DehydrationContext) (err error) {

	var fieldContentResponse clientruntime.Content
	if v := e.Response; v != nil {
		fieldContentResponse = clientruntime.NewContent()
		if err := v.Dehydrate(ctx.CopyWithContent(fieldContentResponse)); err != nil {
			return clientruntime.NewPropertyDehydrationError("response", err)
		}
	}
	ctx.Content().SetProperty("response", fieldContentResponse.Map())
	return nil
}

// StructPath returns StructPath
func (e *TestHTTPWatchChangesOutput) StructPath() clientruntime.StructPath {
	return *localSpecularMeta.structPathTestHTTPWatchChangesOutput.Path()
}

// BodyType entity
type BodyType string

const BodyTypeNormal BodyType = "normal"
const BodyTypeSpecial BodyType = "special"

// String returns the string representation of the enum
func (en BodyType) String() string {
	return string(en)
}

// Optional returns the optional value
func (en BodyType) Optional() *BodyType {
	c := en
	return &c
}

var packagePath = clientruntime.ModulePathFromTrustedValues(
	"speculargo",
	"testpackage",
)

func newSpecularPackage() (pk *clientruntime.Package, err error) {
	pk = clientruntime.NewPackage(packagePath)
	localSpecularMeta.structPathBody, err = pk.NewType(
		"Body",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewBody()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathResponse, err = pk.NewType(
		"Response",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewResponse()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathNotFoundProblem, err = pk.NewType(
		"NotFoundProblem",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewNotFoundProblem()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPGetInput, err = pk.NewType(
		"TestHTTPGetInput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPGetInput()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPGetOutput, err = pk.NewType(
		"TestHTTPGetOutput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPGetOutput()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPOtherInput, err = pk.NewType(
		"TestHTTPOtherInput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPOtherInput()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPOtherOutput, err = pk.NewType(
		"TestHTTPOtherOutput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPOtherOutput()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPWatchChangesInput, err = pk.NewType(
		"TestHTTPWatchChangesInput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPWatchChangesInput()
		}),
	)
	if err != nil {
		return nil, err
	}
	localSpecularMeta.structPathTestHTTPWatchChangesOutput, err = pk.NewType(
		"TestHTTPWatchChangesOutput",
		clientruntime.TypeBuilder(func() clientruntime.Struct {
			return NewTestHTTPWatchChangesOutput()
		}),
	)
	if err != nil {
		return nil, err
	}
	resTestHTTP, err := pk.NewResource("TestHTTP")
	if err != nil {
		return nil, err
	}
	_ = resTestHTTP

	var op *clientruntime.Operation
	op, err = resTestHTTP.NewOperation("Get")
	if err != nil {
		return nil, err
	}

	op.SetInput(SpecularMeta().TestHTTPGetInput())
	op.SetOutput(SpecularMeta().TestHTTPGetInput())

	op, err = resTestHTTP.NewOperation("Other")
	if err != nil {
		return nil, err
	}

	op.SetInput(SpecularMeta().TestHTTPOtherInput())
	op.SetOutput(SpecularMeta().TestHTTPOtherInput())

	op, err = resTestHTTP.NewOperation("WatchChanges")
	if err != nil {
		return nil, err
	}

	op.SetInput(SpecularMeta().TestHTTPWatchChangesInput())
	op.SetOutput(SpecularMeta().TestHTTPWatchChangesInput())
	op.SetStreamed()

	return pk, nil
}

// TestHTTPResource is the TestHTTPResource resource client
type TestHTTPResource struct {
	transport    clientruntime.Transport
	res          *clientruntime.Resource
	get          *clientruntime.Operation
	other        *clientruntime.Operation
	watchChanges *clientruntime.Operation
}

func newTestHTTPResource(
	transport clientruntime.Transport,
	finder clientruntime.ResourceFinder,
) (*TestHTTPResource, error) {
	res := finder.FindResource("TestHTTP")
	r := &TestHTTPResource{
		transport: transport,
		res:       res,
	}
	r.get = res.FindOperation("Get")
	r.other = res.FindOperation("Other")
	r.watchChanges = res.FindOperation("WatchChanges")
	return r, nil
}

// Get entity
func (res *TestHTTPResource) Get(ctx context.Context, input *TestHTTPGetInput) (*TestHTTPGetOutput, error) {
	o, err := res.transport.Execute(ctx, &clientruntime.Request{
		Operation: res.get,
		Input:     input,
	})
	if err != nil {
		return nil, err
	}
	output := o.(*TestHTTPGetOutput)
	return output, nil
}

// Other entity
func (res *TestHTTPResource) Other(ctx context.Context, input *TestHTTPOtherInput) (*TestHTTPOtherOutput, error) {
	o, err := res.transport.Execute(ctx, &clientruntime.Request{
		Operation: res.other,
		Input:     input,
	})
	if err != nil {
		return nil, err
	}
	output := o.(*TestHTTPOtherOutput)
	return output, nil
}

// WatchChanges entity
func (res *TestHTTPResource) WatchChanges(ctx context.Context, input *TestHTTPWatchChangesInput) (<-chan clientruntime.StreamEvent[*TestHTTPWatchChangesOutput], error) {
	ch, err := res.transport.Stream(ctx, &clientruntime.Request{
		Operation: res.watchChanges,
		Input:     input,
	})
	if err != nil {
		return nil, err
	}
	return clientruntime.UnwrapStreamHandler[*TestHTTPWatchChangesOutput](ch), nil
}

// Client is the main client of the package
type Client struct {
	transport clientruntime.Transport
	pk        *clientruntime.Package
	TestHTTP  *TestHTTPResource
}

// WithTransport configures the transport in the client
func WithTransport(transport clientruntime.Transport) clientruntime.OptionFunc {
	return func(o any) error {
		c, ok := o.(*Client)
		if !ok {
			return nil
		}
		c.transport = transport
		return nil
	}
}

// NewEndpointTransport returns a transport with the package base endpoint
func NewEndpointTransport(options ...clientruntime.Option) (clientruntime.Transport, error) {
	o := options
	return clientruntime.NewHTTPJSONTransport(
		"http://localhost:8080",
		o...,
	)
}

// NewClient returns a new instance of Client
func NewClient(options ...clientruntime.Option) (*Client, error) {
	pk := SpecularMeta().Module()
	c := &Client{
		pk: pk,
	}
	if err := clientruntime.ApplyOptions(c, options...); err != nil {
		return nil, err
	}
	if c.transport == nil {
		t, err := NewEndpointTransport(options...)
		if err != nil {
			return nil, err
		}
		c.transport = t
	}
	transport := c.transport
	var err error
	c.TestHTTP, err = newTestHTTPResource(transport, pk)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func init() {
	initSpecularMeta()
}

func initSpecularMeta() {
	pk, err := newSpecularPackage()
	if err != nil {
		panic(errors.New("failed to initialize shared allow package speculargo/testpackage"))
	}
	localSpecularMeta.mod = pk
}

// SpecularMetaInfo defines metadata of the specular module
type SpecularMetaInfo struct {
	mod                                  *clientruntime.Package
	structPathBody                       *clientruntime.StructDefinition
	structPathResponse                   *clientruntime.StructDefinition
	structPathNotFoundProblem            *clientruntime.StructDefinition
	structPathTestHTTPGetInput           *clientruntime.StructDefinition
	structPathTestHTTPGetOutput          *clientruntime.StructDefinition
	structPathTestHTTPOtherInput         *clientruntime.StructDefinition
	structPathTestHTTPOtherOutput        *clientruntime.StructDefinition
	structPathTestHTTPWatchChangesInput  *clientruntime.StructDefinition
	structPathTestHTTPWatchChangesOutput *clientruntime.StructDefinition
}

// Module returns the module definition
func (m *SpecularMetaInfo) Module() *clientruntime.Package {
	return m.mod
}

// Body allows easy access to structure
func (m *SpecularMetaInfo) Body() *clientruntime.StructDefinition {
	return m.structPathBody
}

// Response allows easy access to structure
func (m *SpecularMetaInfo) Response() *clientruntime.StructDefinition {
	return m.structPathResponse
}

// NotFoundProblem allows easy access to structure
func (m *SpecularMetaInfo) NotFoundProblem() *clientruntime.StructDefinition {
	return m.structPathNotFoundProblem
}

// TestHTTPGetInput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPGetInput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPGetInput
}

// TestHTTPGetOutput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPGetOutput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPGetOutput
}

// TestHTTPOtherInput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPOtherInput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPOtherInput
}

// TestHTTPOtherOutput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPOtherOutput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPOtherOutput
}

// TestHTTPWatchChangesInput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPWatchChangesInput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPWatchChangesInput
}

// TestHTTPWatchChangesOutput allows easy access to structure
func (m *SpecularMetaInfo) TestHTTPWatchChangesOutput() *clientruntime.StructDefinition {
	return m.structPathTestHTTPWatchChangesOutput
}

var localSpecularMeta *SpecularMetaInfo = &SpecularMetaInfo{}

// SpecularMeta returns metadata of the specular module
func SpecularMeta() *SpecularMetaInfo {
	return localSpecularMeta
}
