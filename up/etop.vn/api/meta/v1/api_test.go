package v1

import (
	"encoding/json"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	someUUID := "650fa5f8-b89e-4990-bef8-6b1060a6df9e"
	u, err := uuid.FromString(someUUID)
	require.NoError(t, err)

	t.Run("unmarshal json", func(t *testing.T) {
		tests := []struct {
			name   string
			input  string
			error  string
			expect UUID
		}{
			{
				name:   "null",
				input:  "null",
				error:  "",
				expect: UUID{},
			},
			{
				name:   "empty string",
				input:  `""`,
				error:  "",
				expect: UUID{},
			},
			{
				name:   "invalid",
				input:  "100",
				error:  "invalid uuid format",
				expect: UUID{},
			},
			{
				name:   "with value",
				input:  `"` + someUUID + `"`,
				error:  "",
				expect: UUID{Data: u[:]},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var v UUID
				err := json.Unmarshal([]byte(tt.input), &v)
				if tt.error != "" {
					require.EqualError(t, err, tt.error)
					return
				}

				require.NoError(t, err)
				if tt.expect.Data == nil {
					require.Nil(t, v.Data)
				} else {
					require.EqualValues(t, tt.expect.Data, v.Data)
				}
			})
		}
	})
	t.Run("marshal json", func(t *testing.T) {
		tests := []struct {
			name   string
			input  UUID
			expect string
		}{
			{
				name:   "nil",
				input:  UUID{},
				expect: "null",
			},
			{
				name:   "zero",
				input:  UUID{Data: uuid.Nil[:]},
				expect: "null",
			},
			{
				name:   "with value",
				input:  UUID{Data: u[:]},
				expect: `"` + someUUID + `"`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				data, err := json.Marshal(tt.input)
				require.NoError(t, err)
				require.EqualValues(t, tt.expect, data)
			})
		}
	})
}
