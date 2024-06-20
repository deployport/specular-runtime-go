package formats_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	formats "go.deployport.com/specular-runtime/client/formats"
)

func TestReader(t *testing.T) {
	r := formats.NewReaderJSON(bytes.NewReader([]byte(`{"name":"johan", "color": "blue"}`)))
	readNames := []string{}
	readStrings := map[string]string{}
	err := r.Read(func(p *formats.ReaderProp) error {
		readNames = append(readNames, p.Name)
		if p.Name == "name" || p.Name == "color" {
			v, err := p.Value.String()
			require.NoError(t, err)
			readStrings[p.Name] = *v
		}
		return nil
	})
	require.NoError(t, err)
	require.Contains(t, readNames, "name")
	require.Contains(t, readNames, "color")
	require.Equal(t, "johan", readStrings["name"])
	require.Equal(t, "blue", readStrings["color"])
}
