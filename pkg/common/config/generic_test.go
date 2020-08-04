package cc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestGenericConfig(t *testing.T) {
	input := `
alice:
    name: Alice B.
    age: 20
bob:
    name: Bobby C.
    university:
      name: overflow
`
	var alice struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}
	var bob struct {
		Name       string `yaml:"name"`
		University struct {
			Name string `yaml:"name"`
		} `yaml:"university"`
	}

	var gc GenericConfig
	gc.Register("alice", &alice)
	gc.Register("bob", &bob)

	err := yaml.Unmarshal([]byte(input), &gc)
	require.NoError(t, err)

	assert.Equal(t, "Alice B.", alice.Name)
	assert.Equal(t, 20, alice.Age)
	assert.Equal(t, "Bobby C.", bob.Name)
	assert.Equal(t, "overflow", bob.University.Name)
}
