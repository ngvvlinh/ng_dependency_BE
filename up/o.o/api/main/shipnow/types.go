package shipnow

import (
	"time"

	"o.o/api/main/ordering/types"
	v1 "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipnow

type ShipnowFulfillment struct {
	Id                         dot.ID
	ShopId                     dot.ID
	PartnerId                  dot.ID
	PickupAddress              *types.Address
	DeliveryPoints             []*DeliveryPoint
	Carrier                    v1.Carrier
	ShippingServiceCode        string
	ShippingServiceFee         int
	ShippingServiceName        string
	ShippingServiceDescription string
	shippingtypes.WeightInfo
	ValueInfo            shippingtypes.ValueInfo
	ShippingNote         string
	RequestPickupAt      time.Time
	Status               status5.Status
	ShippingStatus       status5.Status
	ShippingCode         string
	ShippingState        shipnow_state.State
	ConfirmStatus        status3.Status
	OrderIds             []dot.ID
	ShippingCreatedAt    time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	EtopPaymentStatus    status4.Status
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
	meta.EventMeta

	ShipnowFulfillmentId dot.ID
	OrderIds             []dot.ID
}

type ShipnowOrderChangedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentId dot.ID
	OldOrderIds          []dot.ID
	OrderIds             []dot.ID
}

type ShipnowCancelledEvent struct {
	meta.EventMeta

	ShipnowFulfillmentId dot.ID
	OrderIds             []dot.ID
	ExternalShipnowId    string
	CarrierServiceCode   string
	CancelReason         string
}

type ShipnowValidateConfirmedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentId dot.ID
	OrderIds             []dot.ID
}

type ShipnowExternalCreatedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentId dot.ID
}

func ShipnowStatus(state shipnow_state.State, paymentStatus status4.Status) status5.Status {
	switch state {
	case shipnow_state.StateDefault,
		shipnow_state.StateCreated,
		shipnow_state.StateAssigning,
		shipnow_state.StatePicking,
		shipnow_state.StateDelivering:
		return status5.S
	case shipnow_state.StateCancelled:
		return status5.N
	case shipnow_state.StateDelivered:
		if paymentStatus == status4.P {
			return status5.P
		}
		return status5.S
	default:
		return status5.Z
	}
	return status5.Z
}
