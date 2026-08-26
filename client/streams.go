package client

// UnwrapStreamHandler is a function that unwraps a stream handler with type input and output
func UnwrapStreamHandler[TOut Struct](from <-chan StreamEvent[Struct]) chan StreamEvent[TOut] {
	to := make(chan StreamEvent[TOut])
	go func() {
		defer close(to)
		for ev := range from {
			// Output is nil on error events
			out, _ := ev.Output.(TOut)
			to <- StreamEvent[TOut]{
				Output: out,
				Err:    ev.Err,
			}
		}
	}()
	return to
}
