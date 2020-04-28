package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"o.o/backend/pkg/common/randgenerator"
)

func TestGenerateServiceID(t *testing.T) {
	g := randgenerator.NewGenerator(100)
	{
		id, err := GenerateServiceID(g, "Nhanh", "53320")
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 7)
		assert.Equal(t, id[1], byte('N'))
	}
	{
		id, err := GenerateServiceID(g, "Chuẩn", "53320")
		t.Log(id)
		assert.NoError(t, err)
		assert.Equal(t, len(id), 7)
		assert.Equal(t, id[1], byte('C'))
	}
}
