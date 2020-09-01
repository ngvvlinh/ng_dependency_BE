package inventory

import (
	"o.o/api/top/types/etc/inventory_auto"
)

func CheckRoleAutoInventoryVoucher(checker func(action string) bool, autoInventoryVoucher inventory_auto.AutoInventoryVoucher) inventory_auto.AutoInventoryVoucher {
	roleLevel := inventory_auto.Unknown
	if checker("shop/inventory:create") == true {
		roleLevel = inventory_auto.Create
	}
	if checker("shop/inventory:confirm") == true {
		roleLevel = inventory_auto.Confirm
	}
	if roleLevel < autoInventoryVoucher {
		return roleLevel
	}
	return autoInventoryVoucher
}
