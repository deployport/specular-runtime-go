package formats_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	formats "go.deployport.com/specular-runtime/client/formats"
)

func TestReader(t *testing.T) {
	r := formats.NewReaderJSON(bytes.NewReader([]byte(`{"name":"johan", "color": "blue", "deletedAt": null, "count": 1}`)))
	readNames := []string{}
	readStrings := map[string]string{}
	deletedAtRead := true
	var countRead json.Number
	err := r.Read(func(p *formats.ReaderProp) error {
		readNames = append(readNames, p.Name)
		if p.Name == "name" || p.Name == "color" {
			v, err := p.Value.String()
			require.NoError(t, err)
			readStrings[p.Name] = *v
			require.False(t, p.Value.IsNull())
		}
		if p.Name == "deletedAt" {
			v := p.Value.IsNull()
			require.True(t, v)
			deletedAtRead = true
		}
		if p.Name == "count" {
			v, err := p.Value.Number()
			require.NoError(t, err)
			countRead = *v
			require.False(t, p.Value.IsNull())
		}
		return nil
	})
	require.NoError(t, err)
	require.Contains(t, readNames, "name")
	require.Contains(t, readNames, "color")
	require.Equal(t, "johan", readStrings["name"])
	require.Equal(t, "blue", readStrings["color"])
	require.True(t, deletedAtRead)
	require.Equal(t, "1", countRead.String())
}
