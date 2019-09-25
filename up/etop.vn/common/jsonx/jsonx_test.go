package jsonx

import (
	"reflect"
	"testing"

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

var withCache int

func reset() {
	if withCache == 0 {
		recognizedTypes = make(map[reflect.Type]modeType)
	}
}

func expectFastpath(t *testing.T, expected int) {
	if withCache == 2 {
		require.Equal(t, expected, fastpath)
	}
}

func TestValidate(t *testing.T) {
	t.Run("without cache", func(t *testing.T) {
		withCache = 0
		testValidate(t)
	})
	t.Run("with cache 1", func(t *testing.T) {
		// execute the first time to populate cache
		withCache = 1
		testValidate(t)
	})
	t.Run("with cache 2", func(t *testing.T) {
		// execute the second time and test for fastpath
		withCache = 2
		testValidate(t)
	})
}

func testValidate(t *testing.T) {
	t.Run("no tag (error)", func(t *testing.T) {
		reset()
		_, err := validate(Invalid{})
		require.EqualError(t, err, "field String of type Invalid must have json tag")
	})
	t.Run("nil", func(t *testing.T) {
		reset()
		mode, err := validate(nil)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 1)
	})
	t.Run("simple safe struct", func(t *testing.T) {
		reset()
		mode, err := validate(A{})
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 2)
	})
	t.Run("simple pointer to safe struct with nil value", func(t *testing.T) {
		reset()
		mode, err := validate((*A)(nil))
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 2)
	})
	t.Run("simple pointer to safe struct", func(t *testing.T) {
		reset()
		mode, err := validate(&A{})
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 2)
	})
	t.Run("complex safe struct with nil value", func(t *testing.T) {
		reset()
		mode, err := validate((*B)(nil))
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 2)
	})
	t.Run("complex safe struct", func(t *testing.T) {
		reset()
		mode, err := validate(&B{})
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 2)
	})
	t.Run("revalidate struct with nil value", func(t *testing.T) {
		reset()
		mode, err := validate((*C)(nil))
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
	})
	t.Run("revalidate struct with empty interface", func(t *testing.T) {
		reset()
		mode, err := validate(&C{})
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
	})
	t.Run("revalidate struct with valid interface", func(t *testing.T) {
		reset()
		mode, err := validate(&C{Interface: &C{}})
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
	})
	t.Run("revalidate struct with invalid interface (error)", func(t *testing.T) {
		reset()
		_, err := validate(&C{Interface: &Invalid{}})
		require.EqualError(t, err, "field Interface of type C: field String of type Invalid must have json tag")
	})
	t.Run("revalidate indirect struct with nil value", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			mode, err := validate((*D)(nil))
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
		t.Run("slice", func(t *testing.T) {
			mode, err := validate((*E)(nil))
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
	})
	t.Run("revalidate indirect struct with empty interface", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{},
			}
			mode, err := validate(value)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{},
			}
			mode, err := validate(value)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
	})
	t.Run("revalidate indirect struct with valid interface", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{
					"one": &C{},
				},
			}
			mode, err := validate(value)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{
					&C{},
				},
			}
			mode, err := validate(value)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
		})
	})
	t.Run("revalidate indirect struct with invalid interface (error)", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{
					"one":   &C{},
					"two":   &C{Interface: &C{}},
					"three": &C{Interface: &Invalid{}},
				},
			}
			_, err := validate(value)
			require.EqualError(t, err, "field MapC of type D: field Interface of type C: field String of type Invalid must have json tag")
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{
					&C{},
					&C{Interface: &C{}},
					&C{Interface: &Invalid{}},
				},
			}
			_, err := validate(value)
			require.EqualError(t, err, "field SliceC of type E: field Interface of type C: field String of type Invalid must have json tag")
		})
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
	t.Run("nil", func(t *testing.T) {
		data, err := Marshal(nil)
		require.NoError(t, err)
		require.Equal(t, `null`, string(data))
	})
	t.Run("interface of nil value", func(t *testing.T) {
		data, err := Marshal((*A)(nil))
		require.NoError(t, err)
		require.Equal(t, `null`, string(data))
	})
	t.Run("invalid", func(t *testing.T) {
		value := C{Interface: Invalid{}}
		_, err := validate(value)
		require.EqualError(t, err, "field Interface of type C: field String of type Invalid must have json tag")
		require.Panics(t, func() {
			_, _ = Marshal(value)
		})
	})
}
