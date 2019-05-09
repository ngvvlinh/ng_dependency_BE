package sender

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildJSON(t *testing.T) {
	var msgs = [][]byte{
		[]byte(`{"a":"b"}`),
		[]byte(`{"c":"d"}`),
	}
	id := int64(1234567890123456789)
	out := buildJSON(id, msgs)

	var callback struct {
		ID      string              `json:"id"`
		Changes []map[string]string `json:"changes"`
	}
	err := json.Unmarshal(out, &callback)
	require.NoError(t, err)
	require.Equal(t, callback.ID, strconv.FormatInt(id, 10))
	require.EqualValues(t, callback.Changes, []map[string]string{
		{"a": "b"},
		{"c": "d"},
	})
}
