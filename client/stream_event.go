package client

// StreamEvent is an event of an stream
type StreamEvent[TOutput Struct] struct {
	Err    error
	Output TOutput
}

// GetErr returns the error of the event
func (e StreamEvent[TOutput]) GetErr() error {
	return e.Err
}

// GetOutput returns the output of the event
func (e StreamEvent[TOutput]) GetOutput() TOutput {
	return e.Output
}
