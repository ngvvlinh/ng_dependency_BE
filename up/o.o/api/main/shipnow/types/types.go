package types

import (
	"time"

	connectiontypes "o.o/api/main/connectioning/types"
	"o.o/api/main/ordering/types"
	v1 "o.o/api/main/shipnow/carrier/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipnow_state"
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
	TryOn         try_on.TryOnCode
	ShippingState shipnow_state.State
}

type ShipnowService struct {
	Carrier            v1.ShipnowCarrier
	Name               string
	Code               string
	Fee                int
	ExpectedPickupAt   time.Time
	ExpectedDeliveryAt time.Time
	Description        string
	ConnectionInfo     *connectiontypes.ConnectionInfo
}
