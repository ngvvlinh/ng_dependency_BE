package api

import (
	"regexp"
)

// /v1 /v1a, /v1beta, /v1/foo
var reVx = regexp.MustCompile(`[a-z0-9]+/v[0-9]+[A-z]*(/[_0-9A-z]+)?$`)
