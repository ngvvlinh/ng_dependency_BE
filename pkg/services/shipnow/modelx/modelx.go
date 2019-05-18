package modelx

import "etop.vn/backend/pkg/services/shipnow/model"

type CreateShipNowFulfillmentArgs struct {
	ShopID              int64
	ShopAddressID       int64
	OrderIDs            []int64
	ShippingNote        string
	ShippingServiceCode string
	ShippingServiceFee  int
	Carrier             string

	Result *model.ShipnowFulfillment
}
