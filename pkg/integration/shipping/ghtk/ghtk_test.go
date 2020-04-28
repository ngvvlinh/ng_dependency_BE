package ghtk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
)

func TestGenerateServiceID(t *testing.T) {
	g := newServiceIDGenerator(100)
	{
		id, err := g.GenerateServiceID('D', ghtkclient.TransportRoad)
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 8)
		assert.Equal(t, id[1], byte('D'))
		assert.Equal(t, id[5], byte('R'))

		code, transport, err := ParseServiceID(id)
		assert.NoError(t, err)
		assert.Equal(t, code, byte('D'))
		assert.Equal(t, transport, ghtkclient.TransportRoad)
	}
	{
		id, err := g.GenerateServiceID('S', ghtkclient.TransportFly)
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 8)
		assert.Equal(t, id[2], byte('S'))
		assert.Equal(t, id[6], byte('F'))

		code, transport, err := ParseServiceID(id)
		assert.NoError(t, err)
		assert.Equal(t, code, byte('S'))
		assert.Equal(t, transport, ghtkclient.TransportFly)
	}
}
