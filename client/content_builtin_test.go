package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestContentRequireTimeProperty(t *testing.T) {
	content := NewContent()
	content.SetProperty("time", time.Date(1999, 7, 7, 15, 35, 0, 0, time.UTC).Format(time.RFC3339))
	var tm time.Time
	err := ContentRequireTimeProperty(content, "time", &tm)
	if err != nil {
		t.Fatal(err)
	}
	if tm != time.Date(1999, 7, 7, 15, 35, 0, 0, time.UTC) {
		t.Fatalf("unexpected time: %v", tm)
	}
}

func TestContentOptionalTimeProperty(t *testing.T) {
	content := NewContent()
	content.SetProperty("time", time.Date(1999, 7, 7, 15, 35, 0, 0, time.UTC).Format(time.RFC3339))
	var tm *time.Time
	err := ContentOptionalTimeProperty(content, "time", &tm)
	require.NoError(t, err)
	if *tm != time.Date(1999, 7, 7, 15, 35, 0, 0, time.UTC) {
		t.Fatalf("unexpected time: %v", tm)
	}
	tm = nil
	content.SetProperty("time", nil)
	err = ContentOptionalTimeProperty(content, "time", &tm)
	if err != nil {
		t.Fatal(err)
	}
	require.Nil(t, tm)
}

func TestContentOptionalArrayBuiltin(t *testing.T) {
	content := NewContent()
	content.SetProperty("stringArray", []string{"a", "b", "c"})
	var tm []string
	err := ContentOptionalBuiltinArrayProperty(content, "stringArray", &tm)
	require.NoError(t, err)
	require.Len(t, tm, 3)
}

func TestContentRequireArrayBuiltin(t *testing.T) {
	content := NewContent()
	content.SetProperty("stringArray", []string{"a", "b", "c"})
	var tm []string
	err := ContentRequireBuiltinArrayProperty(content, "stringArray", &tm)
	require.NoError(t, err)
	require.Len(t, tm, 3)
	require.Equal(t, "a", tm[0])
}

func TestContentOptionalArrayBuiltinNumeric(t *testing.T) {
	content := NewContent()
	content.SetProperty("floatArray", []float64{1.1, 2.1, 1.3})
	var tm []int32
	err := ContentOptionalBuiltinNumericArrayProperty(content, "floatArray", &tm)
	require.NoError(t, err)
	require.Len(t, tm, 3)
	require.Equal(t, int32(1), tm[0])
	require.Equal(t, int32(2), tm[1])
}

func TestContentRequireArrayBuiltinNumeric(t *testing.T) {
	content := NewContent()
	content.SetProperty("floatArray", []any{1.1, 2.1, 1.3})
	var tm []int64
	err := ContentRequireBuiltinNumericArrayProperty(content, "floatArray", &tm)
	require.NoError(t, err)
	require.Len(t, tm, 3)
	require.Equal(t, int64(1), tm[0])
	require.Equal(t, int64(2), tm[1])
}
