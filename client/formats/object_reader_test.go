package formats_test

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client/formats"
)

func TestObjectReader(t *testing.T) {
	r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 1,
			"container1": {
				"type": "image"
			},
			"container2": {
				"type": "video"
			},
			"contList": [
				{
					"type": "image"
				}
			],
			"afterCont1": {
				"type": "image"
			}
		}
	`)))
	readNames := []string{}
	readStrings := map[string]string{}
	deletedAtRead := true
	var countRead json.Number
	err := r.Read(func(p *formats.ReaderProp) error {
		log.Printf("root prop %s", p.Name)
		readNames = append(readNames, p.Name)
		switch p.Name {
		case "name", "color":
			v, err := p.Value.String()
			require.NoError(t, err)
			readStrings[p.Name] = *v
			require.False(t, p.Value.IsNull())
			return nil
		case "deletedAt":
			v := p.Value.IsNull()
			require.True(t, v)
			deletedAtRead = true
			return nil
		case "count":
			v, err := p.Value.Number()
			require.NoError(t, err)
			countRead = *v
			require.False(t, p.Value.IsNull())
			return nil
		case "container1":
			subReader, err := p.Value.Object()
			require.NoError(t, err)
			require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
				log.Printf("sub prop %s", p.Name)
				if p.Name == "type" {
					v, err := p.Value.String()
					require.NoError(t, err)
					require.Equal(t, "image", *v)
					return nil
				}
				return formats.ErrUnknownField
			}))
			return nil
		case "container2":
			subReader, err := p.Value.Object()
			require.NoError(t, err)
			require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
				log.Printf("sub prop %s", p.Name)
				if p.Name == "type" {
					v, err := p.Value.String()
					require.NoError(t, err)
					require.Equal(t, "video", *v)
					return nil
				}
				return formats.ErrUnknownField
			}))
			return nil
		default:
			return formats.ErrUnknownField
		}
	})
	require.NoError(t, err)
	require.Contains(t, readNames, "name")
	require.Contains(t, readNames, "color")
	require.Equal(t, "johan", readStrings["name"])
	require.Equal(t, "blue", readStrings["color"])
	require.True(t, deletedAtRead)
	require.Equal(t, "1", countRead.String())
}

func TestReaderEmpty(t *testing.T) {
	r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`{}`)))
	readNames := []string{}
	err := r.Read(func(p *formats.ReaderProp) error {
		readNames = append(readNames, p.Name)
		t.Fail()
		return formats.ErrUnknownField
	})
	require.NoError(t, err)
	require.Empty(t, readNames)
}