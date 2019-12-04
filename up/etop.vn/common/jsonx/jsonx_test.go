package jsonx

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MyTime time.Time

func (t MyTime) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

func (t *MyTime) UnmarshalJSON(data []byte) error {
	tt := (*time.Time)(t)
	return tt.UnmarshalJSON(data)
}

type A struct {
	Int     int        `json:"int,omitempty"`
	String  string     `json:"string"`
	Byte    byte       `json:"byte,omitempty"`
	Time    time.Time  `json:"time"`
	PtrTime *time.Time `json:"ptr_time,omitempty"`
	MyTime  MyTime     `json:"my_time"`
}

type A1 = A

type B struct {
	A
	*A1

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

type InvalidTime struct {
	Time time.Time
}

var withCache int
var withRoute routeType

func reset() {
	enabledMode = 0
	if withCache == 0 {
		recognizedTypes = make(map[reflect.Type]modeType)
	}
}

func expectFastpath(t *testing.T, expected int, fastpath int) {
	if withCache == 2 {
		require.Equal(t, expected, fastpath)
	}
}

func TestValidate(t *testing.T) {
	t.Run("without cache", func(t *testing.T) {
		withCache = 0
		testValidateWithRoutes(t)
	})
	t.Run("with cache 1", func(t *testing.T) {
		// execute the first time to populate cache
		withCache = 1
		testValidateWithRoutes(t)
	})
	t.Run("with cache 2", func(t *testing.T) {
		// execute the second time and test for fastpath
		withCache = 2
		testValidateWithRoutes(t)
	})
}

func testValidateWithRoutes(t *testing.T) {
	t.Run("with marshal", func(t *testing.T) {
		withRoute = marshal
		testValidate(t)
	})
	t.Run("with unmarshal", func(t *testing.T) {
		withRoute = unmarshal
		testValidate(t)
	})
}

func testValidate(t *testing.T) {
	t.Run("no tag (error)", func(t *testing.T) {
		reset()
		_, _, err := validate(Invalid{}, withRoute)
		require.EqualError(t, err, "field String of type etop.vn/common/jsonx.Invalid must have json tag")
	})
	t.Run("no tag with time (error)", func(t *testing.T) {
		reset()
		_, _, err := validate(InvalidTime{}, withRoute)
		require.EqualError(t, err, "field Time of type etop.vn/common/jsonx.InvalidTime must have json tag")
	})
	t.Run("nil", func(t *testing.T) {
		reset()
		fastpath, mode, err := validate(nil, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		expectFastpath(t, 1, fastpath)
	})
	t.Run("simple safe struct", func(t *testing.T) {
		reset()
		value := A{}
		fastpath, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(value)])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("simple pointer to safe struct with nil value", func(t *testing.T) {
		reset()
		value := (*A)(nil)
		fastpath, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(value).Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("simple pointer to safe struct", func(t *testing.T) {
		reset()
		value := &A{}
		fastpath, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(value).Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("complex safe struct with nil value", func(t *testing.T) {
		reset()
		value := (*B)(nil)
		fastpath, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(value).Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("complex safe struct", func(t *testing.T) {
		reset()
		value := &B{}
		fastpath, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(value).Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("double pointer to safe struct with nil value", func(t *testing.T) {
		reset()
		value0 := (*B)(nil)
		fastpath, mode, err := validate(&value0, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(&value0).Elem().Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("double pointer to safe struct", func(t *testing.T) {
		reset()
		value0 := &B{}
		fastpath, mode, err := validate(&value0, withRoute)
		require.NoError(t, err)
		require.Equal(t, safe, mode)
		require.Equal(t, safe, recognizedTypes[reflect.TypeOf(&value0).Elem().Elem()])
		expectFastpath(t, 2, fastpath)
	})
	t.Run("revalidate struct with nil value", func(t *testing.T) {
		reset()
		value := (*C)(nil)
		_, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
		require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
	})
	t.Run("revalidate struct with empty interface", func(t *testing.T) {
		reset()
		value := &C{}
		_, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
		require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
	})
	t.Run("revalidate struct with valid interface", func(t *testing.T) {
		reset()
		value := &C{Interface: &C{}}
		_, mode, err := validate(value, withRoute)
		require.NoError(t, err)
		require.Equal(t, revalidate, mode)
		require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
	})
	t.Run("revalidate struct with invalid interface (error)", func(t *testing.T) {
		reset()
		_, _, err := validate(&C{Interface: &Invalid{}}, withRoute)
		require.EqualError(t, err, "field Interface of type etop.vn/common/jsonx.C: field String of type etop.vn/common/jsonx.Invalid must have json tag")
	})
	t.Run("revalidate indirect struct with nil value", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := (*D)(nil)
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
		t.Run("slice", func(t *testing.T) {
			value := (*E)(nil)
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
	})
	t.Run("revalidate indirect struct with empty interface", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{},
			}
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{},
			}
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
	})
	t.Run("revalidate indirect struct with valid interface", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{
					"one": {},
				},
			}
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{
					{},
				},
			}
			_, mode, err := validate(value, withRoute)
			require.NoError(t, err)
			require.Equal(t, revalidate, mode)
			require.Equal(t, revalidate, recognizedTypes[reflect.TypeOf(value).Elem()])
		})
	})
	t.Run("revalidate indirect struct with invalid interface (error)", func(t *testing.T) {
		reset()
		t.Run("map", func(t *testing.T) {
			value := &D{
				MapC: map[string]*C{
					"one":   {},
					"two":   {Interface: &C{}},
					"three": {Interface: &Invalid{}},
				},
			}
			_, _, err := validate(value, withRoute)
			require.EqualError(t, err, "field MapC of type etop.vn/common/jsonx.D: field Interface of type etop.vn/common/jsonx.C: field String of type etop.vn/common/jsonx.Invalid must have json tag")
		})
		t.Run("slice", func(t *testing.T) {
			value := &E{
				SliceC: []*C{
					{},
					{Interface: &C{}},
					{Interface: &Invalid{}},
				},
			}
			_, _, err := validate(value, withRoute)
			require.EqualError(t, err, "field SliceC of type etop.vn/common/jsonx.E: field Interface of type etop.vn/common/jsonx.C: field String of type etop.vn/common/jsonx.Invalid must have json tag")
		})
	})
}

func TestMarshal(t *testing.T) {
	reset()
	EnableValidation(Panicking)
	t.Run("marshal ok", func(t *testing.T) {
		data, err := Marshal(A{String: "one"})
		require.NoError(t, err)
		require.Equal(t, `{"string":"one","time":"0001-01-01T00:00:00Z","my_time":"0001-01-01T00:00:00Z"}`, string(data))
	})
	t.Run("unmarshal ok", func(t *testing.T) {
		data := `{"string":"one","time":"0001-01-01T00:00:00Z","my_time":"0001-01-01T00:00:00Z"}`
		var a A
		err := Unmarshal([]byte(data), &a)
		require.NoError(t, err)
		require.Equal(t, "one", a.String)
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
	t.Run("invalid (panic)", func(t *testing.T) {
		value := C{Interface: Invalid{}}
		_, _, err := validate(value, marshal)
		require.EqualError(t, err, "field Interface of type etop.vn/common/jsonx.C: field String of type etop.vn/common/jsonx.Invalid must have json tag")
		require.Panics(t, func() {
			_, _ = Marshal(value)
		})
	})

	t.Run("invalid (error)", func(t *testing.T) {
		reset()
		EnableValidation(Warning)
		value := C{Interface: Invalid{}}
		_, _ = Marshal(value)
		errs := GetErrors()
		require.Len(t, errs, 1)
		require.Equal(t, errs[0].Message, "field Interface of type etop.vn/common/jsonx.C: field String of type etop.vn/common/jsonx.Invalid must have json tag")
	})
}
