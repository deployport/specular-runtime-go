package client

import (
	"errors"
	"testing"
	"time"
)

// The concrete type a generated client unwraps to
type streamProbeOutput struct{}

func (o *streamProbeOutput) StructPath() StructPath { return StructPath{} }

// A bare type assertion panics on the nil Output
func TestAnErrorEventDoesNotPanicTheProcess(t *testing.T) {
	from := make(chan StreamEvent[Struct], 1)
	from <- StreamEvent[Struct]{Err: errors.New("the stream failed")}
	close(from)

	to := UnwrapStreamHandler[*streamProbeOutput](from)

	select {
	case ev, ok := <-to:
		if !ok {
			t.Fatal("the error event was not delivered")
		}
		if ev.Err == nil {
			t.Fatal("the error was dropped")
		}
		if ev.Output != nil {
			t.Fatalf("Output = %v, want nil", ev.Output)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was delivered")
	}
}

// Without this the test above passes on a version that drops every event
func TestAValueEventStillUnwraps(t *testing.T) {
	want := &streamProbeOutput{}
	from := make(chan StreamEvent[Struct], 1)
	from <- StreamEvent[Struct]{Output: want}
	close(from)

	to := UnwrapStreamHandler[*streamProbeOutput](from)

	select {
	case ev, ok := <-to:
		if !ok {
			t.Fatal("the event was dropped")
		}
		if ev.Err != nil {
			t.Fatalf("Err = %v, want nil", ev.Err)
		}
		if ev.Output != want {
			t.Fatalf("Output = %v, want the sent value", ev.Output)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was delivered")
	}
}
