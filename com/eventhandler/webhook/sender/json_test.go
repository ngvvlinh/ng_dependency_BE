package sender

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"o.o/capi/dot"
)

func TestBuildJSON(t *testing.T) {
	var msgs = [][]byte{
		[]byte(`{"a":"b"}`),
		[]byte(`{"c":"d"}`),
	}
	id := dot.ID(1234567890123456789)
	out, inner := buildJSON(id, msgs)

	var callback struct {
		ID      string              `json:"id"`
		Changes []map[string]string `json:"changes"`
	}
	err := json.Unmarshal(out, &callback)
	require.NoError(t, err)
	require.Equal(t, callback.ID, id.String())
	require.EqualValues(t, callback.Changes, []map[string]string{
		{"a": "b"},
		{"c": "d"},
	})
	expected := `[{"a":"b"},{"c":"d"}]`
	require.Equal(t, expected, string(inner))
}
