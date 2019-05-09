package cmsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDial(t *testing.T) {
	parts := reAddress.FindStringSubmatch(`[extopvn:asia-southeast1:etoppg1]:5432`)
	require.Equal(t, len(parts), 2)
	require.Equal(t, parts[1], `extopvn:asia-southeast1:etoppg1`)
}
