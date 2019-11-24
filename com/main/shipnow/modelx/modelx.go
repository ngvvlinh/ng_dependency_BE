package modelx

import (
	"etop.vn/backend/com/main/shipnow/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

type UpdateShipnowFulfillmentArgs struct {
	Fulfilment *model.ShipnowFulfillment
}

type GetByIDArgs struct {
	ID dot.ID

	ShopID    dot.ID
	PartnerID dot.ID
}

type GetShipnowFulfillmentsArgs struct {
	ShopID  dot.ID
	Paging  cm.Paging
	Filters []cm.Filter

	Result struct {
		Total               int
		ShipnowFulfillments []*model.ShipnowFulfillment
	}
}

type GetByShippingCodeArgs struct {
	ShippingCode string
}
