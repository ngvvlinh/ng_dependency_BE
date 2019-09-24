package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"etop.vn/backend/pkg/common/scheme"
)

func TestConvert(t *testing.T) {
	t.Run("A to B", func(t *testing.T) {
		var b B
		a := &A{
			Value:   10,
			Int:     100,
			String:  "hello",
			Strings: []string{"one", "two"},
			C:       &C0{-10},
			Cs:      []*C0{{-100}, {-200}},
		}
		err := scheme.Convert(a, &b)
		require.NoError(t, err)
		assert.Equal(t, b.Value, "10")
		assert.Equal(t, b.Int, int32(100))
		assert.Equal(t, b.String, S("hello"))
		assert.EqualValues(t, b.Strings, []string{"one", "two"})
		assert.Equal(t, b.C.Value, "-10")
		assert.EqualValues(t, b.Cs, []*C1{{"-100"}, {"-200"}})
	})
	t.Run("[]*A to []*B", func(t *testing.T) {
		var bs []*B
		as := []*A{{Value: 10}, {Value: 20}}
		err := scheme.Convert(as, &bs)
		require.NoError(t, err)
		require.Len(t, bs, 2)
		require.Equal(t, bs[0].Value, "10")
		require.Equal(t, bs[1].Value, "20")
	})
}
