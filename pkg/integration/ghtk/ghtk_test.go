package ghtk

import (
	"testing"

	ghtkClient "etop.vn/backend/pkg/integration/ghtk/client"
	"github.com/stretchr/testify/assert"
)

func TestGenerateServiceID(t *testing.T) {
	g := newServiceIDGenerator(100)
	{
		id, err := g.GenerateServiceID('D', ghtkClient.TransportRoad)
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 8)
		assert.Equal(t, id[1], byte('D'))
		assert.Equal(t, id[5], byte('R'))

		code, transport, err := ParseServiceID(id)
		assert.NoError(t, err)
		assert.Equal(t, code, byte('D'))
		assert.Equal(t, transport, ghtkClient.TransportRoad)
	}
	{
		id, err := g.GenerateServiceID('S', ghtkClient.TransportFly)
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 8)
		assert.Equal(t, id[2], byte('S'))
		assert.Equal(t, id[6], byte('F'))

		code, transport, err := ParseServiceID(id)
		assert.NoError(t, err)
		assert.Equal(t, code, byte('S'))
		assert.Equal(t, transport, ghtkClient.TransportFly)
	}
}
