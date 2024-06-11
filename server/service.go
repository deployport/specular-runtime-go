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
}

// NewService creates a new Service
func NewService(opts ...ServiceOptionFunc) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	s.initLogging()
	return s
}

// ServeHTTP implements http.Handler for servicing a package
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := s.logger

	popx, preResult := s.buildHTTPOperationExecution(logger, r)
	if preResult != nil {
		_ = preResult.WriteResponse(w)
		return
	}
	octx := newServiceOperationContext(logger, popx, &s.opsHandlers)
	octx.serve(w)
}

func (s *Service) buildHTTPOperationExecution(
	logger *slog.Logger,
	r *http.Request,
) (*OperationExecution, *client.HTTPResult) {
	pk := s.pk
	path, herr := parseOperationExecutionPathFromURIPath(r.URL.Path)
	if herr != nil {
		return nil, client.HTTPResultForError(herr)
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
		return nil, client.HTTPResultForError(herr)
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
			return nil, pre.BuildFinalResult(nil, err)
		}
		// context possible replaced as part of authentication, we must re-adopt
		pre.Ctx = params.Request.Context()
	}

	opx, herr := buildOperationExecution(logger, pk, pre, r.Body)
	if herr != nil {
		return nil, client.HTTPResultForError(herr)
	}
	return opx, nil
}
