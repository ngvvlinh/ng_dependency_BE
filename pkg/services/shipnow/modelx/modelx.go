package modelx

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/services/shipnow/model"
)

type UpdateShipnowFulfillmentArgs struct {
	Fulfilment *model.ShipnowFulfillment
}

type GetByIDArgs struct {
	ID int64

	ShopID    int64
	PartnerID int64
}

type GetActiveShipnowFulfillmentsByOrderIDArgs struct {
	OrderID                     int64
	ExcludeShipnowFulfillmentID int64
}

type GetShipnowFulfillmentsArgs struct {
	ShopID  int64
	Paging  *cm.Paging
	Filters []cm.Filter

	Result struct {
		Total               int
		ShipnowFulfillments []*model.ShipnowFulfillment
	}
}
