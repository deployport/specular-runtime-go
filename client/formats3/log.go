package formats

// LoggerPrintfFunc is a function type for logging
type LoggerPrintfFunc func(format string, v ...any)

var globalLogger LoggerPrintfFunc = nil

// EnableLogger enables or disables the logger via log.Printf
func EnableLogger(f LoggerPrintfFunc) {
	globalLogger = f
}

type logger struct {
	logPrefix string
}

// isLoggingEnabled returns true if logging is enabled
func (r *logger) isLoggingEnabled() bool {
	return globalLogger != nil
}

// logf prints a log message
func (r *logger) logf(format string, args ...interface{}) {
	if !r.isLoggingEnabled() {
		return
	}
	if r.logPrefix != "" {
		format = "(" + r.logPrefix + ") " + format
	} else {
		format = "(root reader) " + format
	}
	printf := globalLogger
	if printf == nil {
		return
	}
	printf(format, args...)
}
