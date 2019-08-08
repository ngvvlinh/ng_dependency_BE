package modelx

import (
	"etop.vn/backend/com/main/shipnow/model"
	cm "etop.vn/backend/pkg/common"
)

type UpdateShipnowFulfillmentArgs struct {
	Fulfilment *model.ShipnowFulfillment
}

type GetByIDArgs struct {
	ID int64

	ShopID    int64
	PartnerID int64
}

type GetShipnowFulfillmentsArgs struct {
	ShopID  int64
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
