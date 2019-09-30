package convert

import (
	"time"

	"etop.vn/api/shopping/vendoring"
	_ "etop.vn/backend/com/shopping/vendoring/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/shopping/vendoring/model -> etop.vn/api/shopping/vendoring
// +gen:convert: etop.vn/api/shopping/vendoring

func createShopVendor(args *vendoring.CreateVendorArgs, out *vendoring.ShopVendor) {
	apply_vendoring_CreateVendorArgs_vendoring_ShopVendor(args, out)
	out.ID = cm.NewID()
}

func updateShopVendor(args *vendoring.UpdateVendorArgs, out *vendoring.ShopVendor) {
	apply_vendoring_UpdateVendorArgs_vendoring_ShopVendor(args, out)
	out.UpdatedAt = time.Now()
}
