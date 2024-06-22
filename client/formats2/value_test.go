package formats_test

import (
	"reflect"
	"testing"

	require "github.com/stretchr/testify/require"
	formats "go.deployport.com/specular-runtime/client/formats2"
)

func TestValue(t *testing.T) {
	// use reflection to count the fields and check the field names

	v := formats.Value{}
	reflectType := reflect.TypeOf(v)

	require.Equal(t, 3, reflectType.NumField(), "update IsEmpty and Reset when fields are added or removed")
	// check field s
	_, ok := reflectType.FieldByName("s")
	require.True(t, ok, "field s not found")
	// check field null
	_, ok = reflectType.FieldByName("null")
	require.True(t, ok, "field null not found")
	// check field number
	_, ok = reflectType.FieldByName("number")
	require.True(t, ok, "field number not found")
}

func TestValueMarshalJSON(t *testing.T) {
	v := formats.NewValueString("test")
	require.Equal(t, v.Type(), formats.ValueTypeString)
	b, err := v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`"test"`), b)

	v = formats.NewValueNumber("123")
	require.Equal(t, v.Type(), formats.ValueTypeNumber)
	b, err = v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`123`), b)

	v = formats.NewValueNull()
	require.Equal(t, v.Type(), formats.ValueTypeNull)
	b, err = v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`null`), b)

	v = formats.Value{}
	require.Equal(t, v.Type(), formats.ValueTypeUnknown)
	_, err = v.String()
	require.Error(t, err)
	_, err = v.Number()
	require.Error(t, err)
	b, err = v.MarshalJSON()
	require.NoError(t, err)
}
