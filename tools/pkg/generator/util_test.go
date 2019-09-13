package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var empty = struct{}{}

func TestParseDirectives(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		results []Directive
		err     string
	}{
		{
			name:  "simple",
			input: "+gen:sample",
			results: []Directive{
				{Raw: "+gen:sample", Cmd: "gen:sample"},
			},
		},
		{
			name:  "simple with arg (and extra space)",
			input: "+gen:sample=10 ",
			results: []Directive{
				{Raw: "+gen:sample=10", Cmd: "gen:sample", Arg: "10"},
			},
		},
		{
			name:  "multiple",
			input: "+gen:sample=10   +one  +another=10-13 ",
			results: []Directive{
				{Raw: "+gen:sample=10", Cmd: "gen:sample", Arg: "10"},
				{Raw: "+one", Cmd: "one"},
				{Raw: "+another=10-13", Cmd: "another", Arg: "10-13"},
			},
		},
		{
			name:  "with colon",
			input: "+gen:validate: 0 < $ && $ <= 10 + 20 ",
			results: []Directive{
				{Raw: "+gen:validate: 0 < $ && $ <= 10 + 20",
					Cmd: "gen:validate", Arg: "0 < $ && $ <= 10 + 20"},
			},
		},
		{
			// looks like two directives but it's actually one, because of
			// ending ":"
			name:  "with colon 2",
			input: "+gen:sample: +optional",
			results: []Directive{
				{Raw: "+gen:sample: +optional", Cmd: "gen:sample", Arg: "+optional"},
			},
		},
		{
			name:  "no plug sign (error)",
			input: "gen:sample",
			err:   "invalid directive",
		},
		{
			name:  "no ending sign (error)",
			input: "+gen:sample 10",
			err:   "invalid directive",
		},
		{
			name:  "extra space (error)",
			input: "+gen:sample=10 20",
			err:   "invalid directive",
		},
		{
			name:  "extra space 2 (error)",
			input: "+gen:sample= 20",
			err:   "invalid directive",
		},
		{
			name:  "no arg (error)",
			input: "+gen:sample=",
			err:   "invalid directive",
		},
		{
			name:  "no arg 2 (error)",
			input: "+gen:sample= +optional",
			err:   "invalid directive",
		},
		{
			name:  "no arg 3 (error)",
			input: "+gen:sample:",
			err:   "invalid directive",
		},
		{
			name:  "no name",
			input: "+gen:sample=10 + +another",
			err:   "invalid directive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			directives, err := parseDirective(tt.input, nil)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, len(tt.results), len(directives), "expect %#v, got %#v", tt.results, directives)
			for i, d := range directives {
				require.Equal(t, tt.results[i], d)
			}
		})
	}
}

func TestPatterns(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		patterns := []string{
			"example.com/world/water/...",
		}
		ps := parsePatterns(patterns)
		assert.Equal(t, 0, len(ps.paths))
		expectedPrefix := map[string]int{
			"example.com/":             1,
			"example.com/world/":       1,
			"example.com/world/water/": 2,
		}
		assert.EqualValues(t, expectedPrefix, ps.prefixes)
	})
	t.Run("eliminate path", func(t *testing.T) {
		patterns := []string{
			"example.com/underworld/...",
			"example.com/world/water/...",
			"example.com/overworld/cloud",
			"example.com/underworld/underwater", // eliminated
			"example.com/the/other/world",
			"example.com/underworld", // eliminated
		}
		ps := parsePatterns(patterns)
		expectedPath := map[string]struct{}{
			"example.com/overworld/cloud": empty,
			"example.com/the/other/world": empty,
		}
		expectedPrefix := map[string]int{
			"example.com/":             1,
			"example.com/world/":       1,
			"example.com/world/water/": 2,
			"example.com/underworld/":  2,
		}
		assert.EqualValues(t, expectedPath, ps.paths)
		assert.EqualValues(t, expectedPrefix, ps.prefixes)

		t.Run("match", func(t *testing.T) {
			assert.False(t, ps.match("example.com"))
			assert.False(t, ps.match("example.com/world"))
			assert.True(t, ps.match("example.com/world/water"))
			assert.True(t, ps.match("example.com/world/water/fish"))
			assert.False(t, ps.match("example.com/world/fire"))
		})
	})
	t.Run("eliminate prefix", func(t *testing.T) {
		patterns := []string{
			"example.com/overworld/...",
			"example.com/world/water/...", // eliminated
			"example.com/world/...",
			"example.com/the/other/world",
		}
		ps := parsePatterns(patterns)
		expectedPath := map[string]struct{}{
			"example.com/the/other/world": empty,
		}
		expectedPrefix := map[string]int{
			"example.com/":           1,
			"example.com/world/":     2,
			"example.com/overworld/": 2,
		}
		assert.EqualValues(t, expectedPath, ps.paths)
		assert.EqualValues(t, expectedPrefix, ps.prefixes)
	})
}
