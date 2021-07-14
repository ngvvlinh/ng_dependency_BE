package account_user

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"o.o/api/top/types/etc/shop_user_role"
)

func TestHasPermission(t *testing.T) {
	for _, tt := range []struct {
		name         string
		currentRoles []string
		actionRoles  []shop_user_role.UserRole
		expected     bool
	}{
		{
			name:         "Same owner",
			currentRoles: []string{"owner"},
			actionRoles:  []shop_user_role.UserRole{shop_user_role.Owner},
			expected:     false,
		}, {
			name:         "Same staff_management",
			currentRoles: []string{"staff_management"},
			actionRoles:  []shop_user_role.UserRole{shop_user_role.StaffManagement},
			expected:     false,
		}, {
			name:         "Role does not exist",
			currentRoles: []string{"no_existed_role"},
			actionRoles:  nil,
			expected:     false,
		}, {
			name:         "Mutiple current roles accepted",
			currentRoles: []string{"owner", "staff_management"},
			actionRoles:  []shop_user_role.UserRole{shop_user_role.StaffManagement, shop_user_role.Accountant, shop_user_role.InventoryManagement},
			expected:     true,
		}, {
			name:         "Current role staff & action roles not accepted",
			currentRoles: []string{"staff_management"},
			actionRoles:  []shop_user_role.UserRole{shop_user_role.StaffManagement, shop_user_role.Accountant, shop_user_role.InventoryManagement},
			expected:     false,
		}, {
			name:         "Multiple current role & action roles not accepted",
			currentRoles: []string{"accountant", "salesman"},
			actionRoles:  []shop_user_role.UserRole{shop_user_role.StaffManagement, shop_user_role.InventoryManagement},
			expected:     false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ok := hasPermision(tt.currentRoles, tt.actionRoles)
			assert.Equal(t, tt.expected, ok)
		})
	}
}
