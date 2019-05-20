package modelx

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
)

type GetOrderQuery struct {
	OrderID    int64
	ExternalID string
	ShopID     int64
	SupplierID int64
	PartnerID  int64
	Code       string

	// If true, don't filter order lines from other suppliers
	AllSuppliers       bool
	IncludeFulfillment bool

	Result struct {
		Order        *model.Order
		Fulfillments []*shipmodel.Fulfillment
	}
}

type OrderWithFulfillments struct {
	*model.Order

	Fulfillments []*shipmodel.Fulfillment
}

type GetOrdersQuery struct {
	ShopIDs    []int64 // MixedAccount
	SupplierID int64
	PartnerID  int64

	// If true, don't filter order lines from other suppliers
	AllSuppliers bool

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	Status  *int

	Result struct {
		Orders []OrderWithFulfillments
		Shops  []*model.Shop
		Total  int
	}
}
