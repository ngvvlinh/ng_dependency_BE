package types

import (
	"time"

	"etop.vn/api/main/ordering/types"
	v1 "etop.vn/api/main/shipnow/carrier/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
)

type DeliveryPoint struct {
	ShippingAddress *types.Address
	Lines           []*types.ItemLine
	ShippingNote    string
	OrderId         dot.ID
	OrderCode       string
	shippingtypes.WeightInfo
	shippingtypes.ValueInfo
	TryOn try_on.TryOnCode
}

type ShipnowService struct {
	Carrier            v1.Carrier
	Name               string
	Code               string
	Fee                int
	ExpectedPickupAt   time.Time
	ExpectedDeliveryAt time.Time
	Description        string
}
