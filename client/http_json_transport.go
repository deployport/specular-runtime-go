package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

// Submission is the structure with the pre submission information
type Submission struct {
	// HTTPBody is the buffer that will be sent to the HTTP Request body
	// This body is the JSON marshalled struct input of the operation
	HTTPBody []byte
	// HTTPHeaders is the HTTP Headers that will be sent to the HTTP Request
	HTTPRequest      *http.Request
	OperationRequest *Request
	Logger           *slog.Logger
}

// HTTPJSONTransport is the HTTP JSON transport
type HTTPJSONTransport struct {
	// The endpoint URL
	EndpointURL string
	Client      *http.Client

	onSubmission func(ctx context.Context, sub *Submission) error
	Logger       *slog.Logger
}

// // HTTPJSONTransportOption is the option for HTTPJSONTransport
// type HTTPJSONTransportOption func(*HTTPJSONTransport) error

// WithHTTPJSONTransportClient sets the HTTP client
func WithHTTPJSONTransportClient(client *http.Client) OptionFunc {
	return func(o any) error {
		t, ok := o.(*HTTPJSONTransport)
		if !ok {
			return nil
		}
		t.Client = client
		return nil
	}
}

// OnSubmissionHandler is the function that is called before the submission to the HTTP Client
// It's the last chance to modify the request with HTTP Headers before it's sent to the HTTP Client
// If it returns an error, the request is not sent to the HTTP Client and returned to the caller
// If it returns nil, the request is sent to the HTTP Client
type OnSubmissionHandler func(ctx context.Context, sub *Submission) error

// WithOnSubmission sets the OnSubmissionHandler on the HTTPJSONTransport
// If there is already a OnSubmissionHandler, it will be executed before the new one
func WithOnSubmission(fn func(ctx context.Context, sub *Submission) error) OptionFunc {
	return func(o any) error {
		t, ok := o.(*HTTPJSONTransport)
		if !ok {
			return nil
		}
		if previous := t.onSubmission; previous != nil {
			t.onSubmission = func(ctx context.Context, sub *Submission) error {
				if err := previous(ctx, sub); err != nil {
					return err
				}
				return fn(ctx, sub)
			}
			return nil
		}
		t.onSubmission = fn
		return nil
	}
}

// WithLogger sets the logger
func WithLogger(logger *slog.Logger) OptionFunc {
	return func(o any) error {
		t, ok := o.(*HTTPJSONTransport)
		if !ok {
			return nil
		}
		t.Logger = logger
		return nil
	}
}

// WithEndpointURL overrides the endpoint URL in the transport
func WithEndpointURL(endpointURL string) OptionFunc {
	return func(o any) error {
		t, ok := o.(*HTTPJSONTransport)
		if !ok {
			return nil
		}
		t.EndpointURL = endpointURL
		return nil
	}
}

// NewHTTPJSONTransport creates a new HTTP JSON transport
func NewHTTPJSONTransport(endpointURL string, options ...Option) (*HTTPJSONTransport, error) {
	t := &HTTPJSONTransport{
		EndpointURL: endpointURL,
	}
	if err := ApplyOptions(t, options...); err != nil {
		return nil, err
	}
	if t.Client == nil {
		t.Client = http.DefaultClient
	}
	if t.Logger == nil {
		t.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	return t, nil
}

// handleOnSubmission handles the pre submission
func (t *HTTPJSONTransport) handleOnSubmission(ctx context.Context, pre *Submission) error {
	if t.onSubmission != nil {
		return t.onSubmission(ctx, pre)
	}
	return nil
}

// Execute executes the request
func (t *HTTPJSONTransport) Execute(ctx context.Context, req *Request) (Struct, error) {
	hr, err := t.createRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := t.Client.Do(hr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pk := req.Operation.Resource().Package()
	structFinder := NewMultiPackageStructFinder(pk)
	cr, err := ParseHTTPClientResult(structFinder, res.Header, res.Body)
	if err != nil {
		return nil, err
	}
	if mr := cr.Stream; mr != nil {
		mr.Close()
		return nil, fmt.Errorf("unexpected multipart result")
	}
	sr := cr.Single
	if sr == nil {
		return nil, fmt.Errorf("server result expected")
	}
	if err, ok := sr.(error); ok {
		return nil, err
	}
	return sr, nil
}

// Stream executes the request and streams its results
func (t *HTTPJSONTransport) Stream(ctx context.Context, req *Request) (<-chan StreamEvent[Struct], error) {
	hr, err := t.createRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := t.Client.Do(hr)
	if err != nil {
		return nil, err
	}

	pk := req.Operation.Resource().Package()
	structFinder := NewMultiPackageStructFinder(pk)
	cr, err := ParseHTTPClientResult(structFinder, res.Header, res.Body)
	if err != nil {
		return nil, err
	}
	sr := cr.Single
	if st := sr; st != nil {
		if sterr, ok := st.(error); ok {
			return nil, sterr
		}
		return nil, fmt.Errorf("unexpected struct result before stream, %v", st)
	}
	stream := cr.Stream
	if stream == nil {
		return nil, fmt.Errorf("server stream expected")
	}
	go stream.stream(ctx)
	return stream.outputChan, nil
}

func (t *HTTPJSONTransport) createRequest(ctx context.Context, req *Request) (*http.Request, error) {
	inputContent := NewContent()
	err := req.Input.Dehydrate(NewDehydrationContext(inputContent))
	if err != nil {
		return nil, err
	}
	inputBytes, err := json.Marshal(inputContent)
	if err != nil {
		return nil, err
	}
	// log.Printf("input bytes: %v", string(inputBytes))

	operationEndpointURI := t.EndpointURL + "/" + req.Operation.Resource().PackageUniqueName() + "/" + req.Operation.Name()

	hr, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		operationEndpointURI,
		bytes.NewReader(inputBytes),
	)
	if err != nil {
		return nil, err
	}
	submission := &Submission{
		OperationRequest: req,
		HTTPBody:         inputBytes,
		HTTPRequest:      hr,
		Logger:           t.Logger,
	}
	if err := t.handleOnSubmission(ctx, submission); err != nil {
		return nil, err
	}

	return submission.HTTPRequest, nil
}

type clientStream struct {
	reader       *multipart.Reader
	outputChan   chan StreamEvent[Struct]
	structFinder *StructDefinitionFinder
	closer       io.Closer
}

func newClientStream(boundary string, structFinder *StructDefinitionFinder, closer io.ReadCloser) *clientStream {
	reader := multipart.NewReader(closer, boundary)
	outputChan := make(chan StreamEvent[Struct],
		1, // we want to at least keep a message received until the caller can start consuming it
	)
	return &clientStream{
		reader:       reader,
		outputChan:   outputChan,
		structFinder: structFinder,
		closer:       closer,
	}
}

// Close closes the stream
func (cs *clientStream) Close() error {
	return cs.closer.Close()
}
func (cs *clientStream) stream(ctx context.Context) {
	defer func() {
		// log.Printf("closing client stream")
		cs.Close()
		close(cs.outputChan)
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			more := cs.nextEvent()
			if !more {
				return
			}
		}
	}
}

func (cs *clientStream) nextEvent() (more bool) {
	// log.Printf("reading next event")
	res, err := cs.fetchNextPart()
	if _, ok := res.(*Heartbeat); ok {
		return true
	}
	if errors.Is(err, io.EOF) {
		// log.Printf("no more parts")
		return false
	}
	if err != nil {
		cs.outputChan <- StreamEvent[Struct]{
			Output: nil,
			Err:    err,
		}
		return
	}
	more = true
	cs.outputChan <- StreamEvent[Struct]{
		Output: res,
		Err:    nil,
	}
	// log.Printf("more: %v", more)
	return
}

func (cs *clientStream) fetchNextPart() (Struct, error) {
	// log.Printf("fetching next part")
	part, err := cs.reader.NextPart()
	if err != nil {
		return nil, err
	}
	// log.Printf("next part has been fetched")
	defer part.Close()
	result, err := ParseHTTPClientResult(cs.structFinder, http.Header(part.Header), part)
	if err != nil {
		// includes io.EOF when necessary
		return nil, err
	}
	sr := result.Single
	if sr != nil {
		if err, ok := sr.(error); ok {
			return nil, err
		}
		return sr, nil
	} else if stream := result.Stream; stream != nil {
		stream.Close()
		return nil, fmt.Errorf("unexpected stream in stream")
	}
	return nil, fmt.Errorf("expected stream part")
}
