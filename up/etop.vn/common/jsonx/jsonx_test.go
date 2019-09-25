package jsonx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type A struct {
	Int    int    `json:"int,omitempty"`
	String string `json:"string"`
	Byte   byte   `json:"byte,omitempty"`
}

type B struct {
	MapString map[string]string `json:"map_string"`
	MapStruct map[string]A      `json:"map_struct"`
	MapPtrA   map[string]*A     `json:"map_ptr_a"`
	MapPtrB   map[string]*B     `json:"map_ptr_b"` // recursive

	SliceString []string `json:"slice_string"`
	SliceStruct []A      `json:"slice_struct"`
	SlicePtrA   []*A     `json:"slice_ptr_a"`
	SlicePtrB   []*B     `json:"slice_ptr_b"` // recursive

	Func func() `json:"-"` // ignored
}

type C struct {
	String    string      `json:"string"`
	Interface interface{} `json:"interface"` // revalidate
}

type D struct {
	MapC map[string]*C `json:"map_c"` // revalidate
}

type E struct {
	SliceC []*C `json:"slice_c"` // revalidate
}

type Invalid struct {
	String string
}

func runValidate(v interface{}) (modeType, error) {
	return validateTag(reflect.ValueOf(v), reflect.TypeOf(v))
}

var shouldReset bool

func reset() {
	if shouldReset {
		recognizedTypes = make(map[reflect.Type]modeType)
	}
}

func TestValidate(t *testing.T) {
	t.Run("without cache", func(t *testing.T) {
		shouldReset = true
		testValidate(t)
	})
	t.Run("with cache", func(t *testing.T) {
		shouldReset = false
		testValidate(t)
	})
}

func testValidate(t *testing.T) {
	t.Run("no tag (error)", func(t *testing.T) {
		reset()
		_, err := runValidate(Invalid{})
		assert.EqualError(t, err, "field String of type Invalid must have json tag")
	})
	t.Run("simple struct", func(t *testing.T) {
		reset()
		mode, err := runValidate(A{})
		require.NoError(t, err)
		assert.Equal(t, safe, mode)
	})
	t.Run("simple pointer to struct", func(t *testing.T) {
		reset()
		mode, err := runValidate(&A{})
		require.NoError(t, err)
		assert.Equal(t, safe, mode)
	})
	t.Run("complex safe struct", func(t *testing.T) {
		reset()
		mode, err := runValidate(&B{})
		require.NoError(t, err)
		assert.Equal(t, safe, mode)
	})
	t.Run("revalidate struct with empty interface", func(t *testing.T) {
		reset()
		mode, err := runValidate(&C{})
		require.NoError(t, err)
		assert.Equal(t, revalidate, mode)
	})
	t.Run("revalidate struct with valid interface", func(t *testing.T) {
		reset()
		mode, err := runValidate(&C{Interface: &C{}})
		require.NoError(t, err)
		assert.Equal(t, revalidate, mode)
	})
	t.Run("revalidate struct with invalid interface (error)", func(t *testing.T) {
		reset()
		_, err := runValidate(&C{Interface: &Invalid{}})
		require.EqualError(t, err, "field Interface of type C: field String of type Invalid must have json tag")
	})
	t.Run("revalidate indirect struct with empty interface", func(t *testing.T) {
		reset()
		{
			value := &D{
				MapC: map[string]*C{},
			}
			mode, err := runValidate(value)
			require.NoError(t, err)
			assert.Equal(t, revalidate, mode)
		}
		{
			value := &E{
				SliceC: []*C{},
			}
			mode, err := runValidate(value)
			require.NoError(t, err)
			assert.Equal(t, revalidate, mode)
		}
	})
	t.Run("revalidate indirect struct with valid interface", func(t *testing.T) {
		reset()
		{
			value := &D{
				MapC: map[string]*C{
					"one": &C{},
				},
			}
			mode, err := runValidate(value)
			require.NoError(t, err)
			assert.Equal(t, revalidate, mode)
		}
		{
			value := &E{
				SliceC: []*C{
					&C{},
				},
			}
			mode, err := runValidate(value)
			require.NoError(t, err)
			assert.Equal(t, revalidate, mode)
		}
	})
	t.Run("revalidate indirect struct with invalid interface (error)", func(t *testing.T) {
		reset()
		{
			value := &D{
				MapC: map[string]*C{
					"one":   &C{},
					"two":   &C{Interface: &C{}},
					"three": &C{Interface: &Invalid{}},
				},
			}
			_, err := runValidate(value)
			require.EqualError(t, err, "field MapC of type D: field Interface of type C: field String of type Invalid must have json tag")
		}
		{
			value := &E{
				SliceC: []*C{
					&C{},
					&C{Interface: &C{}},
					&C{Interface: &Invalid{}},
				},
			}
			_, err := runValidate(value)
			require.EqualError(t, err, "field SliceC of type E: field Interface of type C: field String of type Invalid must have json tag")
		}
	})
}

func TestMarshal(t *testing.T) {
	reset()
	EnableValidation()
	t.Run("ok", func(t *testing.T) {
		data, err := Marshal(A{String: "one"})
		require.NoError(t, err)
		require.Equal(t, `{"string":"one"}`, string(data))
	})
	t.Run("invalid", func(t *testing.T) {
		value := C{Interface: Invalid{}}
		_, err := runValidate(value)
		require.EqualError(t, err, "field Interface of type C: field String of type Invalid must have json tag")

		assert.Panics(t, func() {
			_, _ = Marshal(value)
		})
	})
}
