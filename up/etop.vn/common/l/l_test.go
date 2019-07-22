package l

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func match(t *testing.T, p, name string) bool {
	r, err := parsePattern(p)
	require.NoError(t, err)

	return r.MatchString(name)
}

func TestMatch(t *testing.T) {
	matched := match(t, "*", "foo/service")
	require.True(t, matched)

	matched = match(t, "foo/*", "foo/service")
	require.True(t, matched)

	matched = match(t, "foo/*", "foo/cmd/foo")
	require.True(t, matched)

	matched = match(t, "foo/*", "bar/service")
	require.False(t, matched)

	matched = match(t, "*/service", "foo/service")
	require.True(t, matched)

	matched = match(t, "*/service", "bar/service")
	require.True(t, matched)

	matched = match(t, "bar/sample", "bar/sample")
	require.True(t, matched)

	matched = match(t, "bar/sample", "foo/service")
	require.False(t, matched)
}
