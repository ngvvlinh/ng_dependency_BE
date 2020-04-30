package convert

import (
	"time"

	"o.o/api/shopping/carrying"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/shopping/carrying/model  -> o.o/api/shopping/carrying
// +gen:convert: o.o/api/shopping/carrying

func createShopCarrier(args *carrying.CreateCarrierArgs, out *carrying.ShopCarrier) {
	apply_carrying_CreateCarrierArgs_carrying_ShopCarrier(args, out)
	out.ID = cm.NewID()
}

func updateShopCarrier(args *carrying.UpdateCarrierArgs, out *carrying.ShopCarrier) {
	apply_carrying_UpdateCarrierArgs_carrying_ShopCarrier(args, out)
	out.UpdatedAt = time.Now()
}
