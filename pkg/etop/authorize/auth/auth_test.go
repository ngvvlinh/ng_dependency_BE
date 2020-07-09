package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"o.o/backend/pkg/etop/authorize/authcommon"
)

func TestCheckAction(t *testing.T) {
	t.Run("Test Check()", func(t *testing.T) {
		Init(authcommon.CommonPolicy)
		roles := Roles{"owner"}
		assert.True(t, roles.Check("shop/customer:view"))
	})

	t.Run("Test Check() 1", func(t *testing.T) {
		Init(authcommon.CommonPolicy)
		roles := Roles{"accountant"}
		assert.False(t, roles.Check("shop/supplier:delete"))
	})

	t.Run("Test Check() 2", func(t *testing.T) {
		Init(authcommon.CommonPolicy)
		roles := Roles{"accountant"}
		assert.True(t, roles.Check("shop/supplier:view"))
	})

	t.Run("Test Invalid Config Policy - Empty Role", func(t *testing.T) {
		defer func() { recover() }()
		buildMapRoleActions(`
		p, permission1, admin, user
		p, permission5, user, ,
		`)
		t.Errorf("Should panic")
	})

	t.Run("Test Invalid Config Policy - Role Not Set", func(t *testing.T) {
		defer func() { recover() }()
		buildMapRoleActions(`
		p, permission1,
		p, permission5, user
		`)
		t.Errorf("Should panic")
	})

	t.Run("Test Invalid Config Policy - Prefix Not Set", func(t *testing.T) {
		defer func() { recover() }()
		buildMapRoleActions(`
		permission1, admin
		p, permission5, user
		`)
		t.Errorf("Should panic")
	})

	t.Run("Test Invalid Config Policy - Duplicate Action", func(t *testing.T) {
		defer func() { recover() }()
		buildMapRoleActions(`
		p, permission1, admin
		p, permission2, admin
		p, permission3, admin
		p, permission4, admin
		p, permission1, user
		`)
		t.Errorf("Should panic")
	})

	t.Run("BuildMapRoleActions", func(t *testing.T) {
		expected := map[string][]string{
			"admin": []string{"permission1", "permission2", "permission3", "permission4"},
			"user":  []string{"permission1", "permission5"},
			"owner": []string{"permission2"},
			"loser": []string{"permission3", "permission4"},
		}
		actual := buildMapRoleActions(`
		p, permission1, admin, user
		p, permission2, admin, owner
		p, permission3, admin, loser
		p, permission4, admin, loser
		p, permission5, user
		`)
		assert.EqualValues(t, expected, actual)
	})

	t.Run("ListActionsByRoles", func(t *testing.T) {
		Init(`
		p, permission1, admin, user
		p, permission2, admin
		`)
		expected := []string{"permission1", "permission2"}
		actual := ListActionsByRoles([]string{"admin"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission1"}
		actual = ListActionsByRoles([]string{"user"})
		assert.EqualValues(t, []string{"permission1"}, ListActionsByRoles([]string{"user"}))
	})

	t.Run("Test ListActionsByRoles Complicated Case", func(t *testing.T) {
		Init(`
		p, permission0, admin, owner
		p, permission1, admin, owner, user
		p, permission2, admin, user
		p, permission3, user
		p, permission4, ac_admin
		p, permission5, sa_admin
		`)
		expected := []string{"permission0", "permission1", "permission2"}
		actual := ListActionsByRoles([]string{"admin"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission1", "permission2", "permission3"}
		actual = ListActionsByRoles([]string{"user"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission0", "permission1"}
		actual = ListActionsByRoles([]string{"owner"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission0", "permission1", "permission2", "permission3"}
		actual = ListActionsByRoles([]string{"admin", "user"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission5"}
		actual = ListActionsByRoles([]string{"sa_admin"})
		assert.EqualValues(t, expected, actual)
	})
}
