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
