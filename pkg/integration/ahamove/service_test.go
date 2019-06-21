package ahamove

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateServiceID(t *testing.T) {
	g := newServiceIDGenerator(100)
	{
		id, err := g.GenerateServiceID('D', SGNBIKE)
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 8)
		assert.Equal(t, id[1], byte('D'))

		code, serviceID, err := parseServiceCode(id)
		assert.NoError(t, err)
		assert.Equal(t, code, byte('D'))
		assert.Equal(t, serviceID, SGNBIKE)
	}
}
