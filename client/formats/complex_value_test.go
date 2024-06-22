package formats_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client/formats"
)

func TestComplexValueTest(t *testing.T) {
	t.Run("Object", func(t *testing.T) {
		v := formats.NewComplexValueObject(map[string]formats.ComplexValue{})
		require.Equal(t, formats.ValueTypeObject, v.Type())
		require.NotNil(t, v)
		require.False(t, v.IsEmpty())
		_, err := v.Object()
		require.NoError(t, err)
		_, err = v.Array()
		require.Error(t, err)
	})

	t.Run("Array", func(t *testing.T) {
		v := formats.NewComplexValueArray([]formats.ComplexValue{})
		require.Equal(t, formats.ValueTypeArray, v.Type())
		require.NotNil(t, v)
		require.False(t, v.IsEmpty())
		_, err := v.Object()
		require.Error(t, err)
		_, err = v.Array()
		require.NoError(t, err)
	})

	t.Run("SimpleValue", func(t *testing.T) {
		v := formats.NewComplexValue(formats.NewValueNull())
		require.NotNil(t, v)
		require.False(t, v.IsEmpty())
		require.Equal(t, formats.ValueTypeNull, v.Type())
	})
}
