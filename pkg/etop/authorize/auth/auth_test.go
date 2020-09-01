package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"o.o/backend/pkg/etop/authorize/authcommon"
)

func TestCheckAction(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		authorizer := New(authcommon.CommonPolicy)
		t.Run("Test Check()", func(t *testing.T) {
			roles := Roles{"owner"}
			assert.True(t, authorizer.CheckSingle(roles, "shop/customer:view"))
		})
		t.Run("Test Check() 1", func(t *testing.T) {
			roles := Roles{"accountant"}
			assert.False(t, authorizer.CheckSingle(roles, "shop/supplier:delete"))
		})
		t.Run("Test Check() 2", func(t *testing.T) {
			roles := Roles{"accountant"}
			assert.True(t, authorizer.CheckSingle(roles, "shop/supplier:view"))
		})
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
			"admin": {"permission1", "permission2", "permission3", "permission4"},
			"user":  {"permission1", "permission5"},
			"owner": {"permission2"},
			"loser": {"permission3", "permission4"},
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
		authorizer := New(`
		p, permission1, admin, user
		p, permission2, admin
		`)
		expected := []string{"permission1", "permission2"}
		actual := authorizer.ListActionsByRoles([]string{"admin"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission1"}
		actual = authorizer.ListActionsByRoles([]string{"user"})
		assert.EqualValues(t, []string{"permission1"}, authorizer.ListActionsByRoles([]string{"user"}))
	})

	t.Run("Test ListActionsByRoles Complicated Case", func(t *testing.T) {
		authorizer := New(`
		p, permission0, admin, owner
		p, permission1, admin, owner, user
		p, permission2, admin, user
		p, permission3, user
		p, permission4, ac_admin
		p, permission5, sa_admin
		`)
		expected := []string{"permission0", "permission1", "permission2"}
		actual := authorizer.ListActionsByRoles([]string{"admin"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission1", "permission2", "permission3"}
		actual = authorizer.ListActionsByRoles([]string{"user"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission0", "permission1"}
		actual = authorizer.ListActionsByRoles([]string{"owner"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission0", "permission1", "permission2", "permission3"}
		actual = authorizer.ListActionsByRoles([]string{"admin", "user"})
		assert.EqualValues(t, expected, actual)

		expected = []string{"permission5"}
		actual = authorizer.ListActionsByRoles([]string{"sa_admin"})
		assert.EqualValues(t, expected, actual)
	})
}
