package inventory

import (
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/backend/pkg/etop/authorize/auth"
)

func CheckRoleAutoInventoryVoucher(roles auth.Roles, autoInventoryVoucher inventory_auto.AutoInventoryVoucher) inventory_auto.AutoInventoryVoucher {
	roleLevel := inventory_auto.Unknown
	if roles.Check("shop/inventory:create") == true {
		roleLevel = inventory_auto.Create
	}
	if roles.Check("shop/inventory:confirm") == true {
		roleLevel = inventory_auto.Confirm
	}
	if roleLevel < autoInventoryVoucher {
		return roleLevel
	}
	return autoInventoryVoucher
}
