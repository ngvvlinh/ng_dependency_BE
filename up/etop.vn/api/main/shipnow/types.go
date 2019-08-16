package shipnow

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	v1 "etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
)

type ShipnowFulfillment struct {
	Id                         int64
	ShopId                     int64
	PartnerId                  int64
	PickupAddress              *types.Address
	DeliveryPoints             []*DeliveryPoint
	Carrier                    v1.Carrier
	ShippingServiceCode        string
	ShippingServiceFee         int32
	ShippingServiceName        string
	ShippingServiceDescription string
	shippingtypes.WeightInfo
	ValueInfo            shippingtypes.ValueInfo
	ShippingNote         string
	RequestPickupAt      time.Time
	Status               etop.Status5
	ShippingStatus       etop.Status5
	ShippingCode         string
	ShippingState        shipnowtypes.State
	ConfirmStatus        etop.Status3
	OrderIds             []int64
	ShippingCreatedAt    time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	EtopPaymentStatus    etop.Status4
	CodEtopTransferedAt  time.Time
	ShippingPickingAt    time.Time
	ShippingDeliveringAt time.Time
	ShippingDeliveredAt  time.Time
	ShippingCancelledAt  time.Time
	ShippingSharedLink   string
	CancelReason         string
}

type DeliveryPoint = shipnowtypes.DeliveryPoint

type SyncStates struct {
	SyncAt    time.Time
	TrySyncAt time.Time
	Error     *meta.Error
}

type ShipnowOrderReservationEvent struct {
	ShipnowFulfillmentId int64
	OrderIds             []int64
}

type ShipnowOrderChangedEvent struct {
	ShipnowFulfillmentId int64
	OldOrderIds          []int64
	OrderIds             []int64
}

type ShipnowCancelledEvent struct {
	ShipnowFulfillmentId int64
	OrderIds             []int64
	ExternalShipnowId    string
	CarrierServiceCode   string
	CancelReason         string
}

type ShipnowValidateConfirmedEvent struct {
	ShipnowFulfillmentId int64
	OrderIds             []int64
}

type ShipnowCreateExternalEvent struct {
	ShipnowFulfillmentId int64
}

func ShipnowStatus(state shipnowtypes.State, paymentStatus etop.Status4) etop.Status5 {
	switch state {
	case shipnowtypes.StateDefault:
	case shipnowtypes.StateCreated:
	case shipnowtypes.StateAssigning:
	case shipnowtypes.StatePicking:
	case shipnowtypes.StateDelivering:
		return etop.S5SuperPos
	case shipnowtypes.StateCancelled:
		return etop.S5Negative
	case shipnowtypes.StateDelivered:
		if paymentStatus == etop.S4Positive {
			return etop.S5Positive
		}
		return etop.S5SuperPos
	default:
		return etop.S5Zero
	}
	return etop.S5Zero
}
