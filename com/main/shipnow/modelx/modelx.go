package modelx

import (
	"o.o/backend/com/main/shipnow/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
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
