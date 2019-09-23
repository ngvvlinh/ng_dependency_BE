package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"etop.vn/backend/pkg/common/scheme"
)

func TestConvert(t *testing.T) {
	t.Run("A to B", func(t *testing.T) {
		var b B
		err := scheme.Convert(&A{10}, &b)
		require.NoError(t, err)
		require.Equal(t, b.Value, "10")
	})
	t.Run("[]*A to []*B", func(t *testing.T) {
		var bs []*B
		as := []*A{{10}, {20}}
		err := scheme.Convert(as, &bs)
		require.NoError(t, err)
		require.Len(t, bs, 2)
		require.Equal(t, bs[0].Value, "10")
		require.Equal(t, bs[1].Value, "20")
	})
}
