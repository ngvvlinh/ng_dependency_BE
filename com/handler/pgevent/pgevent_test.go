package pgevent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEventPayload(t *testing.T) {
	data := "name:123:UPDATE:456"
	t.Run("name:123:UPDATE:456", func(t *testing.T) {
		event, err := ParseEventPayload(data)
		require.Nil(t, err)
		require.EqualValues(t, &PgEvent{
			EventKey:  "name:456",
			ID:        456,
			RID:       123,
			Table:     "name",
			Op:        "UPDATE",
			Keys:      nil,
			Timestamp: event.Timestamp,
		}, event)
	})

	data = "name:123:UPDATE:sh456:p789"
	t.Run(data, func(t *testing.T) {
		event, err := ParseEventPayload(data)
		require.Nil(t, err)
		require.EqualValues(t, &PgEvent{
			EventKey: "name:sh456",
			ID:       456,
			RID:      123,
			Table:    "name",
			Op:       "UPDATE",
			Keys: map[string]int64{
				"sh": 456,
				"p":  789,
			},
			Timestamp: event.Timestamp,
		}, event)
	})

	data = "name:123:UPDATE:sh456:789"
	t.Run(data+" (error)", func(t *testing.T) {
		_, err := ParseEventPayload(data)
		require.EqualError(t, err, "Empty key (789)")
	})
}
