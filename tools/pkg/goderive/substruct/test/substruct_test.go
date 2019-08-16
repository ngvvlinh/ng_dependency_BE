package test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstruct(t *testing.T) {
	foo := &Foo{
		A:    "A",
		I:    1,
		SS:   []string{"a", "b"},
		PS:   ps("ps"),
		PSS:  pss("a", "b"),
		PSPS: psps("a", "b"),
	}
	expected := `{
		"A": "A",
		"I": 1,
		"SS": ["a","b"],
		"PS": "ps",
		"PSS": ["a","b"],
		"PSPS": ["a","b"]
	}`
	foo1 := NewFoo1FromFoo(foo)
	data, _ := json.Marshal(foo1)
	assert.Equal(t, clean(expected), string(data))
}

func ps(s string) *string {
	return &s
}

func pss(ss ...string) *[]string {
	return &ss
}

func psps(ss ...string) *[]*string {
	sps := make([]*string, len(ss))
	for i, s := range ss {
		sps[i] = ps(s)
	}
	return &sps
}

var reClean = regexp.MustCompile(`\s+`)

func clean(s string) string {
	return reClean.ReplaceAllString(s, "")
}
