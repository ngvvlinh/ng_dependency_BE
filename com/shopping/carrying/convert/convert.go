package convert

import (
	"time"

	"etop.vn/api/shopping/carrying"
	_ "etop.vn/backend/com/shopping/carrying/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/shopping/carrying/model -> etop.vn/api/shopping/carrying
// +gen:convert: etop.vn/api/shopping/carrying

func createShopCarrier(args *carrying.CreateCarrierArgs, out *carrying.ShopCarrier) {
	apply_carrying_CreateCarrierArgs_carrying_ShopCarrier(args, out)
	out.ID = cm.NewID()
}

func updateShopCarrier(args *carrying.UpdateCarrierArgs, out *carrying.ShopCarrier) {
	apply_carrying_UpdateCarrierArgs_carrying_ShopCarrier(args, out)
	out.UpdatedAt = time.Now()
}
