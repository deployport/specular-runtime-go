package formats

import "testing"

func TestError(t *testing.T) {
	t.Run("ValueError", func(t *testing.T) {
		err := NewValueError("test")
		if err.Error() != "test" {
			t.Errorf("expected error message to be 'test', got %s", err.Error())
		}
		if !err.Is(NewValueError("test")) {
			t.Errorf("expected error to be a ValueError")
		}
		if !IsValueError(NewValueError("test")) {
			t.Errorf("expected error to be a ValueError")
		}
	})
}
