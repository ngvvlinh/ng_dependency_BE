package l

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func match(t *testing.T, p, name string, match bool, expectedLevel int) {
	lvl, r, err := parsePattern(p)
	require.NoError(t, err)
	require.Equal(t, match, r.MatchString(name), "pattern")
	require.Equal(t, expectedLevel, lvl, "level")
}

func TestMatch(t *testing.T) {
	match(t, "*", "foo/service", true, 1)
	match(t, "foo/*", "foo/service", true, 1)
	match(t, "foo/*", "foo/cmd/foo", true, 1)
	match(t, "foo/*", "bar/service", false, 1)

	match(t, "*/service", "foo/service", true, 1)
	match(t, "*/service", "bar/service", true, 1)
	match(t, "bar/sample", "bar/sample", true, 1)
	match(t, "bar/sample", "foo/service", false, 1)

	match(t, "hello:2", "hello", true, 2)
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		input  string
		expect zapcore.Level
		ok     bool
	}{
		{"info", zapcore.InfoLevel, true},
		{"debug", zapcore.DebugLevel, true},
		{"debug-0", 0, false},
		{"debug-1", -1, true},
		{"debug-3", -3, true},
		{"DEBUG-9", -9, true},
		{"DEBUG-10", 0, false},

		// must be "debug" or "DEBUG"
		{"Debug-3", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			lvl, ok := unmarshalLevel(tt.input)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.expect, lvl)
		})
	}
}
