package formats_test

import (
	"reflect"
	"testing"

	require "github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client/formats"
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
	b, err := v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`"test"`), b)

	v = formats.NewValueNumber("123")
	b, err = v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`123`), b)

	v = formats.NewValueNull()
	b, err = v.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte(`null`), b)
}
