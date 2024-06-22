package formats_test

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client/formats"
)

func TestReader(t *testing.T) {
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
				"item1",
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
		case "contList":
			subReader, err := p.Value.Array()
			require.NoError(t, err)
			require.NoError(t, subReader.Read(func(i *formats.ReaderItem) error {
				log.Printf("sub item %d", i.Index)
				if i.Index == 0 {
					v, err := i.Value.String()
					require.NoError(t, err)
					require.Equal(t, "item1", *v)
					return nil
				}
				subReader, err := i.Value.Object()
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
			}))
			return nil
		case "afterCont1":
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

func TestReaderUnknownsSimple(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 123.45
		}
	`)))
		r.UseUnknownFields()
		readNames := []string{}
		readStrings := map[string]string{}
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
			default:
				return formats.ErrUnknownField
			}
		})
		require.NoError(t, err)
		require.Contains(t, readNames, "name")
		require.Contains(t, readNames, "color")
		require.Equal(t, "johan", readStrings["name"])
		require.Equal(t, "blue", readStrings["color"])
		unknowns := r.UnknownFields()
		require.NotNil(t, unknowns)
		require.Contains(t, unknowns, "deletedAt")
		require.Contains(t, unknowns, "count")
		deletedAtValue := unknowns["deletedAt"]
		require.True(t, deletedAtValue.IsNull())
		countValue := unknowns["count"]
		count, err := countValue.Number()
		require.NoError(t, err)
		require.Equal(t, "123.45", count.String())
	})
	t.Run("not parsing", func(t *testing.T) {
		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 1
		}
	`)))
		readNames := []string{}
		readStrings := map[string]string{}
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
			default:
				return formats.ErrUnknownField
			}
		})
		require.NoError(t, err)
		require.Contains(t, readNames, "name")
		require.Contains(t, readNames, "color")
		require.Equal(t, "johan", readStrings["name"])
		require.Equal(t, "blue", readStrings["color"])
		unknowns := r.UnknownFields()
		require.Nil(t, unknowns, "unknowns should be nil because UseUnknownFields was not called")
	})
}

func TestReaderUnknownsComplex(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"container1": {
				"type": "pdf"
			},
			"count": 123.4,
			"fruits": ["apple", "banana", null, 55454],
			"contList": [
				{
					"name": "container one"
				},
				[1]
			],
			"knownList": ["ok"]
		}
	`)))
		r.UseUnknownFields()
		readNames := []string{}
		readStrings := map[string]string{}
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
			case "knownList":
				// we parse knownList right after the container object
				// to make sure the json decoder for sub-objects moves forward
				v, err := p.Value.Array()
				require.NoError(t, err)
				err = v.Read(func(i *formats.ReaderItem) error {
					log.Printf("known item %d", i.Index)
					if i.Index == 0 {
						v, err := i.Value.String()
						require.NoError(t, err)
						require.Equal(t, "ok", *v)
						return nil
					}
					return formats.ErrUnknownField
				})
				require.NoError(t, err)
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
		unknowns := r.UnknownFields()
		require.NotNil(t, unknowns)
		require.Contains(t, unknowns, "deletedAt")
		require.Contains(t, unknowns, "count")
		require.Contains(t, unknowns, "container1")
		countVal := unknowns["count"]
		count, err := countVal.Number()
		require.NoError(t, err)
		require.Equal(t, "123.4", count.String())
		container1Value := unknowns["container1"]
		container1Object, err := container1Value.Object()
		container1Type := container1Object["type"]
		require.NoError(t, err)
		v, err := container1Type.String()
		require.NoError(t, err)
		require.Equal(t, "pdf", *v)
		require.Contains(t, unknowns, "fruits")
		fruitsValue := unknowns["fruits"]
		fruits, err := fruitsValue.Array()
		require.NoError(t, err)
		require.Len(t, fruits, 4)
		fruits0 := fruits[0]
		fruits0Value, err := fruits0.String()
		require.NoError(t, err)
		require.Equal(t, "apple", *fruits0Value)
		fruits1 := fruits[1]
		fruits1Value, err := fruits1.String()
		require.NoError(t, err)
		require.Equal(t, "banana", *fruits1Value)
		fruits2 := fruits[2]
		require.True(t, fruits2.IsNull())
		fruits3 := fruits[3]
		fruits3Value, err := fruits3.Number()
		require.NoError(t, err)
		require.Equal(t, "55454", fruits3Value.String())
		require.Contains(t, unknowns, "contList")
		contListValue := unknowns["contList"]
		contList, err := contListValue.Array()
		require.NoError(t, err)
		require.Len(t, contList, 2)
		contList0 := contList[0]
		contList0Object, err := contList0.Object()
		contList0Name := contList0Object["name"]
		v, err = contList0Name.String()
		require.NoError(t, err)
		require.Equal(t, "container one", *v)
		nestedArray := contList[1]
		nestedArrayValue, err := nestedArray.Array()
		require.NoError(t, err)
		require.Len(t, nestedArrayValue, 1)
		nestedArrayValue0, err := nestedArrayValue[0].Number()
		require.NoError(t, err)
		require.Equal(t, "1", nestedArrayValue0.String())
	})
	t.Run("not parsing", func(t *testing.T) {
		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"container1": {
				"type": "pdf"
			},
			"count": 123.4
		}
	`)))
		readNames := []string{}
		readStrings := map[string]string{}
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
			case "count":
				// we parse count right after the container object
				// to make sure the json decoder for sub-objects moves forward
				v, err := p.Value.Number()
				require.NoError(t, err)
				require.Equal(t, "123.4", v.String())
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
		unknowns := r.UnknownFields()
		require.Nil(t, unknowns)
	})
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
