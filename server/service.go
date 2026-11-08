package server

import (
	"log/slog"
	"net/http"
	"time"

	"go.deployport.com/specular-runtime/client"
)

const heartbeatInterval = 5 * time.Second

// WithServicePackage is a function that sets the package
func WithServicePackage(pk *client.Package) ServiceOptionFunc {
	return func(s *Service) {
		s.pk = pk
	}
}

// ServiceOptionFunc is a function that applies an option to the ServiceOptions
type ServiceOptionFunc func(*Service)

// Service is the server for a service
type Service struct {
	pk *client.Package
	opsHandlers
	logging
	unhandledErrorHandling
	authentication
	// telemetry is nil unless WithTelemetry configured a write URL.
	telemetry *telemetryRecorder
}

// NewService creates a new Service
func NewService(opts ...ServiceOptionFunc) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	s.initLogging()
	s.initTelemetry()
	return s
}

// initTelemetry gives the recorder the service logger and starts its interval
// loop. A recorder with no write URL is nil, so nothing is started and nothing
// is pushed.
func (s *Service) initTelemetry() {
	recorder := s.telemetry
	if recorder == nil {
		return
	}
	recorder.logger = s.logger
	go recorder.run()
}

// ServeHTTP implements http.Handler for servicing a package. Every operation of
// every service passes through here, which is why one measurement here covers
// every endpoint.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := s.logger
	recorder := s.telemetry
	var started time.Time
	if recorder != nil {
		started = recorder.now()
	}

	writer := newStatusResponseWriter(w)
	popx, pre, preResult := s.buildHTTPOperationExecution(logger, r)
	if preResult != nil {
		_ = preResult.WriteResponse(writer)
		// pre is set as soon as the operation is known, so a request refused by
		// authentication is still measured against its operation. A path that
		// names no operation leaves pre nil and creates no series.
		if recorder != nil && pre != nil {
			recorder.observe(pre.Operation, writer.status, recorder.now().Sub(started))
		}
		return
	}
	octx := newServiceOperationContext(logger, popx, &s.opsHandlers)
	octx.serve(writer)
	if recorder != nil {
		// A handler that panics unwinds past this line and is not recorded.
		// Guessing a class for it would report whatever the writer saw before
		// the panic, and a partial write is not a status class.
		recorder.observe(popx.Operation, writer.status, recorder.now().Sub(started))
	}
}

// buildHTTPOperationExecution resolves the request to an operation. It returns
// the pre execution as soon as the operation is known, even when the request is
// then refused or its body is rejected, so the caller can measure it.
func (s *Service) buildHTTPOperationExecution(
	logger *slog.Logger,
	r *http.Request,
) (*OperationExecution, *PreOperationExecution, *client.HTTPResult) {
	pk := s.pk
	path, herr := parseOperationExecutionPathFromURIPath(r.URL.Path)
	if herr != nil {
		return nil, nil, client.HTTPResultForError(herr)
	}
	ctx := r.Context()
	pre, herr := buildPreOperationExecution(
		ctx,
		logger,
		&s.unhandledErrorHandling,
		pk,
		path,
	)
	if herr != nil {
		return nil, nil, client.HTTPResultForError(herr)
	}
	r = r.WithContext(pre.Ctx)
	//  authenticate
	if handler := s.authentication.onAuthenticate; handler != nil {
		params := &AuthenticationParams{
			OperationExecution: pre,
			Request:            r,
			Logger:             logger,
		}
		if err := handler(params); err != nil {
			return nil, pre, pre.BuildFinalResult(nil, err)
		}
		// context possible replaced as part of authentication, we must re-adopt
		pre.Ctx = params.Request.Context()
	}

	opx, herr := buildOperationExecution(logger, pk, pre, r.Body)
	if herr != nil {
		return nil, pre, client.HTTPResultForError(herr)
	}
	return opx, pre, nil
}
