package query

import (
	"testing"

	"github.com/stretchr/testify/require"

	"o.o/api/meta"
	"o.o/api/supporting/crm/vtiger"
	"o.o/backend/com/supporting/crm/vtiger/mapping"
)

func TestBuildVtigerQuery(t *testing.T) {
	fieldMap := mapping.ConfigGroup{
		"one": "field_1",
		"two": "field_2",
	}
	tests := []struct {
		name      string
		module    string
		condition map[string]string
		orderBy   *vtiger.OrderBy
		paging    *meta.Paging

		expected string
		err      string
	}{
		{
			name:     "Only module",
			module:   "sample",
			expected: "SELECT * FROM sample;",
		},
		{
			name:   "All fields",
			module: "sample",
			condition: map[string]string{
				"one":  "10",
				"nine": "text",
				"ten":  "100",
			},
			orderBy: &vtiger.OrderBy{"created_at", "DESC"},
			paging: &meta.Paging{
				Offset: 20,
				Limit:  15,
			},
			expected: "SELECT * FROM sample WHERE nine = 'text' AND field_1 = '10' AND ten = '100' ORDER BY created_at DESC LIMIT 20, 15;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := buildVtigerQuery(fieldMap, tt.module, tt.condition, tt.orderBy, tt.paging)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.Equal(t, tt.expected, output)
		})
	}
}
