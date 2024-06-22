package formats_test

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
	formats "go.deployport.com/specular-runtime/client/formats3"
)

// func TestReader(t *testing.T) {
// 	r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
// 		{
// 			"name": "johan",
// 			"color": "blue",
// 			"deletedAt": null,
// 			"count": 1,
// 			"container1": {
// 				"type": "image"
// 			},
// 			"container2": {
// 				"type": "video"
// 			},
// 			"contList": [
// 				"item1",
// 				{
// 					"type": "image"
// 				}
// 			],
// 			"afterCont1": {
// 				"type": "image"
// 			}
// 		}
// 	`)))
// 	readNames := []string{}
// 	readStrings := map[string]string{}
// 	deletedAtRead := true
// 	var countRead json.Number
// 	err := r.Read(func(p *formats.ReaderProp) error {
// 		log.Printf("root prop %s", p.Name)
// 		readNames = append(readNames, p.Name)
// 		switch p.Name {
// 		case "name", "color":
// 			v, err := p.Value.String()
// 			require.NoError(t, err)
// 			readStrings[p.Name] = *v
// 			require.False(t, p.Value.IsNull())
// 			return nil
// 		case "deletedAt":
// 			v := p.Value.IsNull()
// 			require.True(t, v)
// 			deletedAtRead = true
// 			return nil
// 		case "count":
// 			v, err := p.Value.Number()
// 			require.NoError(t, err)
// 			countRead = *v
// 			require.False(t, p.Value.IsNull())
// 			return nil
// 		case "container1":
// 			subReader, err := p.Value.Object()
// 			require.NoError(t, err)
// 			require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
// 				log.Printf("sub prop %s", p.Name)
// 				if p.Name == "type" {
// 					v, err := p.Value.String()
// 					require.NoError(t, err)
// 					require.Equal(t, "image", *v)
// 					return nil
// 				}
// 				return formats.ErrUnknownField
// 			}))
// 			return nil
// 		case "container2":
// 			subReader, err := p.Value.Object()
// 			require.NoError(t, err)
// 			require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
// 				log.Printf("sub prop %s", p.Name)
// 				if p.Name == "type" {
// 					v, err := p.Value.String()
// 					require.NoError(t, err)
// 					require.Equal(t, "video", *v)
// 					return nil
// 				}
// 				return formats.ErrUnknownField
// 			}))
// 			return nil
// 		case "contList":
// 			subReader, err := p.Value.Array()
// 			require.NoError(t, err)
// 			require.NoError(t, subReader.Read(func(i *formats.ReaderItem) error {
// 				log.Printf("sub item %d", i.Index)
// 				if i.Index == 0 {
// 					v, err := i.Value.String()
// 					require.NoError(t, err)
// 					require.Equal(t, "item1", *v)
// 					return nil
// 				}
// 				subReader, err := i.Value.Object()
// 				require.NoError(t, err)
// 				require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
// 					log.Printf("sub prop %s", p.Name)
// 					if p.Name == "type" {
// 						v, err := p.Value.String()
// 						require.NoError(t, err)
// 						require.Equal(t, "image", *v)
// 						return nil
// 					}
// 					return formats.ErrUnknownField
// 				}))
// 				return nil
// 			}))
// 			return nil
// 		case "afterCont1":
// 			subReader, err := p.Value.Object()
// 			require.NoError(t, err)
// 			require.NoError(t, subReader.Read(func(p *formats.ReaderProp) error {
// 				log.Printf("sub prop %s", p.Name)
// 				if p.Name == "type" {
// 					v, err := p.Value.String()
// 					require.NoError(t, err)
// 					require.Equal(t, "image", *v)
// 					return nil
// 				}
// 				return formats.ErrUnknownField
// 			}))
// 			return nil
// 		default:
// 			return formats.ErrUnknownField
// 		}
// 	})
// 	require.NoError(t, err)
// 	require.Contains(t, readNames, "name")
// 	require.Contains(t, readNames, "color")
// 	require.Equal(t, "johan", readStrings["name"])
// 	require.Equal(t, "blue", readStrings["color"])
// 	require.True(t, deletedAtRead)
// 	require.Equal(t, "1", countRead.String())
// }

// func TestReaderUnknownsSimple(t *testing.T) {
// 	t.Run("parsing", func(t *testing.T) {
// 		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
// 		{
// 			"name": "johan",
// 			"color": "blue",
// 			"deletedAt": null,
// 			"count": 123.45
// 		}
// 	`)))
// 		r.UseUnknownFields()
// 		readNames := []string{}
// 		readStrings := map[string]string{}
// 		err := r.Read(func(p *formats.ReaderProp) error {
// 			log.Printf("root prop %s", p.Name)
// 			readNames = append(readNames, p.Name)
// 			switch p.Name {
// 			case "name", "color":
// 				v, err := p.Value.String()
// 				require.NoError(t, err)
// 				readStrings[p.Name] = *v
// 				require.False(t, p.Value.IsNull())
// 				return nil
// 			default:
// 				return formats.ErrUnknownField
// 			}
// 		})
// 		require.NoError(t, err)
// 		require.Contains(t, readNames, "name")
// 		require.Contains(t, readNames, "color")
// 		require.Equal(t, "johan", readStrings["name"])
// 		require.Equal(t, "blue", readStrings["color"])
// 		unknowns := r.UnknownFields()
// 		require.NotNil(t, unknowns)
// 		require.Contains(t, unknowns, "deletedAt")
// 		require.Contains(t, unknowns, "count")
// 		deletedAtValue := unknowns["deletedAt"]
// 		require.True(t, deletedAtValue.IsNull())
// 		countValue := unknowns["count"]
// 		count, err := countValue.Number()
// 		require.NoError(t, err)
// 		require.Equal(t, "123.45", count.String())
// 	})
// 	t.Run("not parsing", func(t *testing.T) {
// 		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
// 		{
// 			"name": "johan",
// 			"color": "blue",
// 			"deletedAt": null,
// 			"count": 1
// 		}
// 	`)))
// 		readNames := []string{}
// 		readStrings := map[string]string{}
// 		err := r.Read(func(p *formats.ReaderProp) error {
// 			log.Printf("root prop %s", p.Name)
// 			readNames = append(readNames, p.Name)
// 			switch p.Name {
// 			case "name", "color":
// 				v, err := p.Value.String()
// 				require.NoError(t, err)
// 				readStrings[p.Name] = *v
// 				require.False(t, p.Value.IsNull())
// 				return nil
// 			default:
// 				return formats.ErrUnknownField
// 			}
// 		})
// 		require.NoError(t, err)
// 		require.Contains(t, readNames, "name")
// 		require.Contains(t, readNames, "color")
// 		require.Equal(t, "johan", readStrings["name"])
// 		require.Equal(t, "blue", readStrings["color"])
// 		unknowns := r.UnknownFields()
// 		require.Nil(t, unknowns, "unknowns should be nil because UseUnknownFields was not called")
// 	})
// }

// func TestReaderUnknownsComplex(t *testing.T) {
// 	t.Run("parsing", func(t *testing.T) {
// 		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
// 		{
// 			"name": "johan",
// 			"color": "blue",
// 			"deletedAt": null,
// 			"container1": {
// 				"type": "pdf"
// 			},
// 			"count": 123.4,
// 			"fruits": ["apple", "banana", null, 55454],
// 			"contList": [
// 				{
// 					"name": "container one"
// 				},
// 				[1]
// 			],
// 			"knownList": ["ok"]
// 		}
// 	`)))
// 		r.UseUnknownFields()
// 		readNames := []string{}
// 		readStrings := map[string]string{}
// 		err := r.Read(func(p *formats.ReaderProp) error {
// 			log.Printf("root prop %s", p.Name)
// 			readNames = append(readNames, p.Name)
// 			switch p.Name {
// 			case "name", "color":
// 				v, err := p.Value.String()
// 				require.NoError(t, err)
// 				readStrings[p.Name] = *v
// 				require.False(t, p.Value.IsNull())
// 				return nil
// 			case "knownList":
// 				// we parse knownList right after the container object
// 				// to make sure the json decoder for sub-objects moves forward
// 				v, err := p.Value.Array()
// 				require.NoError(t, err)
// 				err = v.Read(func(i *formats.ReaderItem) error {
// 					log.Printf("known item %d", i.Index)
// 					if i.Index == 0 {
// 						v, err := i.Value.String()
// 						require.NoError(t, err)
// 						require.Equal(t, "ok", *v)
// 						return nil
// 					}
// 					return formats.ErrUnknownField
// 				})
// 				require.NoError(t, err)
// 				return nil
// 			default:
// 				return formats.ErrUnknownField
// 			}
// 		})
// 		require.NoError(t, err)
// 		require.Contains(t, readNames, "name")
// 		require.Contains(t, readNames, "color")
// 		require.Equal(t, "johan", readStrings["name"])
// 		require.Equal(t, "blue", readStrings["color"])
// 		unknowns := r.UnknownFields()
// 		require.NotNil(t, unknowns)
// 		require.Contains(t, unknowns, "deletedAt")
// 		require.Contains(t, unknowns, "count")
// 		require.Contains(t, unknowns, "container1")
// 		countVal := unknowns["count"]
// 		count, err := countVal.Number()
// 		require.NoError(t, err)
// 		require.Equal(t, "123.4", count.String())
// 		container1Value := unknowns["container1"]
// 		container1Object, err := container1Value.Object()
// 		container1Type := container1Object["type"]
// 		require.NoError(t, err)
// 		v, err := container1Type.String()
// 		require.NoError(t, err)
// 		require.Equal(t, "pdf", *v)
// 		require.Contains(t, unknowns, "fruits")
// 		fruitsValue := unknowns["fruits"]
// 		fruits, err := fruitsValue.Array()
// 		require.NoError(t, err)
// 		require.Len(t, fruits, 4)
// 		fruits0 := fruits[0]
// 		fruits0Value, err := fruits0.String()
// 		require.NoError(t, err)
// 		require.Equal(t, "apple", *fruits0Value)
// 		fruits1 := fruits[1]
// 		fruits1Value, err := fruits1.String()
// 		require.NoError(t, err)
// 		require.Equal(t, "banana", *fruits1Value)
// 		fruits2 := fruits[2]
// 		require.True(t, fruits2.IsNull())
// 		fruits3 := fruits[3]
// 		fruits3Value, err := fruits3.Number()
// 		require.NoError(t, err)
// 		require.Equal(t, "55454", fruits3Value.String())
// 		require.Contains(t, unknowns, "contList")
// 		contListValue := unknowns["contList"]
// 		contList, err := contListValue.Array()
// 		require.NoError(t, err)
// 		require.Len(t, contList, 2)
// 		contList0 := contList[0]
// 		contList0Object, err := contList0.Object()
// 		contList0Name := contList0Object["name"]
// 		v, err = contList0Name.String()
// 		require.NoError(t, err)
// 		require.Equal(t, "container one", *v)
// 		nestedArray := contList[1]
// 		nestedArrayValue, err := nestedArray.Array()
// 		require.NoError(t, err)
// 		require.Len(t, nestedArrayValue, 1)
// 		nestedArrayValue0, err := nestedArrayValue[0].Number()
// 		require.NoError(t, err)
// 		require.Equal(t, "1", nestedArrayValue0.String())
// 	})
// 	t.Run("not parsing", func(t *testing.T) {
// 		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
// 		{
// 			"name": "johan",
// 			"color": "blue",
// 			"deletedAt": null,
// 			"container1": {
// 				"type": "pdf"
// 			},
// 			"count": 123.4
// 		}
// 	`)))
// 		readNames := []string{}
// 		readStrings := map[string]string{}
// 		err := r.Read(func(p *formats.ReaderProp) error {
// 			log.Printf("root prop %s", p.Name)
// 			readNames = append(readNames, p.Name)
// 			switch p.Name {
// 			case "name", "color":
// 				v, err := p.Value.String()
// 				require.NoError(t, err)
// 				readStrings[p.Name] = *v
// 				require.False(t, p.Value.IsNull())
// 				return nil
// 			case "count":
// 				// we parse count right after the container object
// 				// to make sure the json decoder for sub-objects moves forward
// 				v, err := p.Value.Number()
// 				require.NoError(t, err)
// 				require.Equal(t, "123.4", v.String())
// 				return nil
// 			default:
// 				return formats.ErrUnknownField
// 			}
// 		})
// 		require.NoError(t, err)
// 		require.Contains(t, readNames, "name")
// 		require.Contains(t, readNames, "color")
// 		require.Equal(t, "johan", readStrings["name"])
// 		require.Equal(t, "blue", readStrings["color"])
// 		unknowns := r.UnknownFields()
// 		require.Nil(t, unknowns)
// 	})
// }

// func TestReaderEmpty(t *testing.T) {
// 	r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`{}`)))
// 	readNames := []string{}
// 	err := r.Read(func(p *formats.ReaderProp) error {
// 		readNames = append(readNames, p.Name)
// 		t.Fail()
// 		return formats.ErrUnknownField
// 	})
// 	require.NoError(t, err)
// 	require.Empty(t, readNames)
// }

// func TestArrayReader(t *testing.T) {
// 	t.Run("simple", func(t *testing.T) {
// 		r := formats.NewArrayReaderJSON(bytes.NewReader([]byte(`["item1", "item2", "item3"]`)))
// 		readItems := []string{}
// 		err := r.Read(func(i *formats.ReaderItem) error {
// 			v, err := i.Value.String()
// 			require.NoError(t, err)
// 			readItems = append(readItems, *v)
// 			return nil
// 		})
// 		require.NoError(t, err)
// 		require.Len(t, readItems, 3)
// 		require.Contains(t, readItems, "item1")
// 		require.Contains(t, readItems, "item2")
// 		require.Contains(t, readItems, "item3")
// 	})
// 	t.Run("failing to read token", func(t *testing.T) {
// 		r := formats.NewArrayReaderJSON(newfailReader(bytes.NewReader([]byte(`["item1", "item2", "item3"]`)), 1))
// 		err := r.Read(func(i *formats.ReaderItem) error {
// 			return nil
// 		})
// 		require.Error(t, err)
// 	})
// 	t.Run("object stream", func(t *testing.T) {
// 		r := formats.NewArrayReaderJSON(bytes.NewReader([]byte(`{}`)))
// 		err := r.Read(func(i *formats.ReaderItem) error {
// 			return nil
// 		})
// 		require.Error(t, err)
// 	})
// 	t.Run("array item token error", func(t *testing.T) {
// 		r := formats.NewArrayReaderJSON(bytes.NewReader([]byte(`[:]`)))
// 		err := r.Read(func(i *formats.ReaderItem) error {
// 			return nil
// 		})
// 		require.ErrorContains(t, err, "failed to decode array item token")
// 	})
// 	t.Run("array end token error", func(t *testing.T) {
// 		r := formats.NewArrayReaderJSON(bytes.NewReader([]byte(`["1"}`)))
// 		err := r.Read(func(i *formats.ReaderItem) error {
// 			return nil
// 		})
// 		require.ErrorContains(t, err, "failed to decode end of array token")
// 	})
// }

// // failReader is an io.Reader that fails after N bytes.
// type failReader struct {
// 	r         io.Reader // The underlying reader.
// 	n         int       // Number of bytes after which to fail.
// 	bytesRead int       // Counter for the number of bytes read.
// }

// // NewfailReader creates a new failReader.
// func newfailReader(r io.Reader, n int) *failReader {
// 	return &failReader{
// 		r: r,
// 		n: n,
// 	}
// }

// // Read reads data into p. It fails with an error after N bytes.
// func (f *failReader) Read(p []byte) (int, error) {
// 	return 0, fmt.Errorf("intentional error after %d bytes", f.n)
// 	// log.Printf("cap: %v", cap(p))
// 	// if f.bytesRead >= f.n {
// 	// 	return 0, errors.New("custom error after N bytes")
// 	// }

// 	// n, err := f.r.Read(p[:f.n-f.bytesRead])
// 	// f.bytesRead += n
// 	// if f.bytesRead > f.n {
// 	// 	// Simulate an error after reading more than N bytes.
// 	// 	return n, errors.New("custom error after N bytes")
// 	// }
// 	// return n, err
// }

var runTimes = 1

var nameBytes = []byte("name")
var colorBytes = []byte("color")
var deletedAtBytes = []byte("deletedAt")
var countBytes = []byte("count")

func equalBytes(s1, s2 []byte) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func BenchmarkObjectReader(b *testing.B) {
	for i := 0; i < runTimes; i++ {
		tp := containerStruct{}
		r := formats.NewObjectReaderJSON(bytes.NewReader([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 1
		}
	`)))
		err := r.Read(func(p *formats.ReaderProp) error {
			// log.Printf("root prop %s", p.Name)
			if equalBytes(p.Name, nameBytes) {
				tp.Name = string(p.Value.GetStringBytes())
				return nil
			} else if equalBytes(p.Name, colorBytes) {
				tp.Color = string(p.Value.GetStringBytes())
				return nil
			} else if equalBytes(p.Name, deletedAtBytes) {
				_ = p.Value.Type() == fastjson.TypeNull
				return nil
			} else if equalBytes(p.Name, countBytes) {
				tp.Count = p.Value.GetInt()
				return nil
			} else {
				return formats.ErrUnknownField
			}
		})
		require.NoError(b, err)
		i++
	}
}

type containerStruct struct {
	Name      string  `json:"name"`
	Color     string  `json:"color"`
	DeletedAt *string `json:"deletedAt"`
	Count     int     `json:"count"`
}

func BenchmarkJSONSerialization(b *testing.B) {
	for i := 0; i < runTimes; i++ {
		tp := containerStruct{}
		err := json.Unmarshal([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 1
		}`), &tp)
		require.NoError(b, err)
		i++
	}

}

func BenchmarkFastJSONSerialization(b *testing.B) {
	var p fastjson.Parser
	for i := 0; i < runTimes; i++ {

		v, err := p.Parse(`{
                "str": "bar",
                "int": 123,
                "float": 1.23,
                "bool": true,
                "arr": [1, "foo", {}]
        }`)
		// v.GetObject().Visit()
		require.NoError(b, err)
		// fmt.Printf("foo=%s\n", v.GetStringBytes("str"))
		// fmt.Printf("int=%d\n", v.GetInt("int"))
		// fmt.Printf("float=%f\n", v.GetFloat64("float"))
		// fmt.Printf("bool=%v\n", v.GetBool("bool"))
		// fmt.Printf("arr.1=%s\n", v.GetStringBytes("arr", "1"))
		tp := struct {
			Name      string  `json:"name"`
			Color     string  `json:"color"`
			DeletedAt *string `json:"deletedAt"`
			Count     int     `json:"count"`
		}{}

		tp.Name = string(v.GetStringBytes("name"))
		tp.Color = string(v.GetStringBytes("color"))
		if vv := v.GetStringBytes("deletedAt"); vv != nil {
			// tp.DeletedAt = string(vv)
		}
		tp.Count = v.GetInt("count")
		i++
	}

}

type containerType string

const (
	containerTypeDefault containerType = ""
	containerTypeImage   containerType = "image"
	containerTypeVideo   containerType = "video"
)

type containerStruct2 struct {
	Name      string         `json:"name,omitempty"`
	Color     string         `json:"color,omitempty"`
	DeletedAt *string        `json:"deletedAt,omitempty"`
	Count     int            `json:"count,omitempty"`
	Type      containerType  `json:"type,omitempty"`
	TypeOpt   *containerType `json:"typeOpt,omitempty"`
}

// ensureDefaults sets the default values for the struct.
func (c *containerStruct2) SpecularEnsureDefaults() {
	if c.Type == containerTypeDefault {
		c.Type = containerTypeImage
	}
}

type aliasContainerStruct2 containerStruct2

func (c *containerStruct2) UnmarshalJSON(data []byte) error {
	// Step 2: Unmarshal the data into the alias, leveraging the default unmarshaling.
	var alias aliasContainerStruct2
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	(*containerStruct2)(&alias).SpecularEnsureDefaults()
	*c = containerStruct2(alias)
	return nil
}

// MarshalJSON marshals the struct into JSON.
func (c containerStruct2) MarshalJSON() ([]byte, error) {
	log.Printf("marshal JSON")
	alias := aliasContainerStruct2(c)
	if alias.Type == containerTypeImage {
		alias.Type = containerTypeDefault
	}
	// Step 2: Marshal the alias, leveraging the default marshaling.
	return json.Marshal(alias)
}

func TestJSON2Serialization(t *testing.T) {
	tp := containerStruct2{}
	err := json.Unmarshal([]byte(`
		{
			"name": "johan",
			"color": "blue",
			"deletedAt": null,
			"count": 1
		}`), &tp)
	require.NoError(t, err)
	require.Equal(t, containerTypeImage, tp.Type)
}

func TestJSON2MarshalSerialization(t *testing.T) {
	tp := &containerStruct2{
		Type: containerTypeImage,
	}
	m, err := json.Marshal(tp)
	require.NoError(t, err)
	require.Equal(t, `{}`, string(m))
}
