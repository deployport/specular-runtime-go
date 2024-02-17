package server

import "log/slog"

func createDefaultLogger() *slog.Logger {
	return slog.Default()
}

// WithServiceLogger is a function that sets the logger
func WithServiceLogger(logger *slog.Logger) ServiceOptionFunc {
	return func(s *Service) {
		s.logger = logger
	}
}

type logging struct {
	logger *slog.Logger
}

func (s *logging) initLogging() {
	if s.logger != nil {
		return
	}
	s.logger = createDefaultLogger()
}
