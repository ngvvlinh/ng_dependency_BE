package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldResembleSlice(t *testing.T) {
	t.Run("Remove Space 1", func(t *testing.T) {
		code := NormalizeExternalCode("    a  -----------   b  ----------   c   ddddd   333  aaaa    ")
		assert.EqualValues(t, "a-b-c-ddddd-333-aaaa", code)
	})
	t.Run("Remove Space 1", func(t *testing.T) {
		code := NormalizeExternalCode(" a  -----------   b  ----------   c   ddddd   333  aaaa ")
		assert.EqualValues(t, "a-b-c-ddddd-333-aaaa", code)
	})
	t.Run("Remove Space 2", func(t *testing.T) {
		code := NormalizeExternalCode(" PRD013336 000232 22 ")
		assert.EqualValues(t, "PRD013336-000232-22", code)
	})
	t.Run("Remove Space 2", func(t *testing.T) {
		code := NormalizeExternalCode("")
		assert.EqualValues(t, "", code)
	})
	t.Run("Remove Space 3", func(t *testing.T) {
		code := NormalizeExternalCode("---   ----")
		assert.EqualValues(t, "", code)
	})
	t.Run("Remove Space 4", func(t *testing.T) {
		code := NormalizeExternalCode("-----!@&*#R^&%!@#-----")
		assert.EqualValues(t, "R", code)
	})
}
