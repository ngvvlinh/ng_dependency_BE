package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAction(t *testing.T) {
	t.Run("Test Check()", func(t *testing.T) {
		roles := Roles{"owner"}
		assert.True(t, roles.Check("shop/customer:view"))
	})

	t.Run("Test Check() 1", func(t *testing.T) {
		roles := Roles{"accountant"}
		assert.False(t, roles.Check("shop/supplier:delete"))
	})

	t.Run("Test Check() 2", func(t *testing.T) {
		roles := Roles{"accountant"}
		assert.True(t, roles.Check("shop/supplier:view"))
	})
}
