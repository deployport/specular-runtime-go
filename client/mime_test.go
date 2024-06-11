package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMIME(t *testing.T) {
	mime, err := NewMIME("application/spec.myns.mymod")
	require.NoError(t, err)
	app, err := mime.StructPath()
	require.NoError(t, err)
	require.Equal(t, "mymod", app.Module().Name())
	require.Equal(t, "myns", app.Module().Namespace())
}
