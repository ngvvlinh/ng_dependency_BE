package types

import (
	"time"

	"o.o/api/main/ordering/types"
	v1 "o.o/api/main/shipnow/carrier/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
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
