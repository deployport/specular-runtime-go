package client

// UnwrapStreamHandler is a function that unwraps a stream handler with type input and output
func UnwrapStreamHandler[TOut Struct](from <-chan StreamEvent[Struct]) chan StreamEvent[TOut] {
	to := make(chan StreamEvent[TOut])
	go func() {
		defer close(to)
		for ev := range from {
			to <- StreamEvent[TOut]{
				Output: ev.Output.(TOut),
				Err:    ev.Err,
			}
		}
	}()
	return to
}
