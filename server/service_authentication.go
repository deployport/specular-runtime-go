package server

import (
	"errors"
	"log/slog"
	"net/http"
)

// AuthenticationParams is the parameters for authentication
type AuthenticationParams struct {
	OperationExecution *PreOperationExecution
	Request            *http.Request
	Logger             *slog.Logger
}
type authentication struct {
	onAuthenticate func(params AuthenticationParams) error
}

// AuthenticationError is an error that occurs during authentication
type AuthenticationError struct {
	// Message is the error message
	Message string
}

// Error returns the error message
func (err *AuthenticationError) Error() string {
	return err.Message
}

// Is implements errors.Is
func (err *AuthenticationError) Is(target error) bool {
	_, ok := target.(*AuthenticationError)
	return ok
}

// IsAuthenticationError returns true if the error is an authentication error
func IsAuthenticationError(err error) bool {
	return errors.Is(err, &AuthenticationError{})
}

// WithAuthentication is a function that sets the authentication function
func WithAuthentication(authFunc func(params AuthenticationParams) error) ServiceOptionFunc {
	return func(s *Service) {
		s.onAuthenticate = authFunc
	}
}
