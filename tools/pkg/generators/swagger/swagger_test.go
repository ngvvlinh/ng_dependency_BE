package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseMethodDescriptionText(t *testing.T) {
	tests := []struct {
		name                    string
		input                   string
		wantSummary             string
		wantDescription         string
		wantFormatedDescription string
		wantDeprecated          bool
		wantError               string
	}{
		{
			name:      "no colon",
			input:     "deprecated invalid\n",
			wantError: "invalid keyword, must contain @ or : (deprecated invalid)",
		},
		{
			name:                    "single",
			input:                   "deprecated\n",
			wantSummary:             "",
			wantDescription:         "deprecated\n",
			wantFormatedDescription: "**Deprecated:**\n",
			wantDeprecated:          true,
		},
		{
			name:                    "with colon",
			input:                   "DEPRECATED: no longer use\n",
			wantSummary:             "",
			wantDescription:         "DEPRECATED: no longer use\n",
			wantFormatedDescription: "**Deprecated:** no longer use\n",
			wantDeprecated:          true,
		},
		{
			name:                    "with paragraph",
			input:                   "a paragraph\ndeprecated: no longer use\nanother paragraph\n",
			wantSummary:             "",
			wantDescription:         "a paragraph\ndeprecated: no longer use\nanother paragraph\n",
			wantFormatedDescription: "a paragraph\n**Deprecated:** no longer use\nanother paragraph\n",
			wantDeprecated:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseItemDescriptionText(tt.input)
			if tt.wantError != "" {
				require.EqualError(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantSummary, got.Summary, "summary")
			assert.Equal(t, tt.wantDescription, got.Description, "description")
			assert.Equal(t, tt.wantFormatedDescription, got.FormattedDescription, "formatedDescription")
			assert.Equal(t, tt.wantDeprecated, got.Deprecated, "deprecated")
		})
	}
}
