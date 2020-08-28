package shipnow

import (
	"time"

	"o.o/api/main/ordering/types"
	v1 "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipnow

type ShipnowFulfillment struct {
	ID                         dot.ID
	ShopID                     dot.ID
	PartnerID                  dot.ID
	PickupAddress              *types.Address
	DeliveryPoints             []*DeliveryPoint
	Carrier                    v1.ShipnowCarrier
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
	OrderIDs             []dot.ID
	ShippingCreatedAt    time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	EtopPaymentStatus    status4.Status
	CODEtopTransferedAt  time.Time
	ShippingPickingAt    time.Time
	ShippingDeliveringAt time.Time
	ShippingDeliveredAt  time.Time
	ShippingCancelledAt  time.Time
	ShippingSharedLink   string
	CancelReason         string
	ConnectionID         dot.ID
	ConnectionMethod     connection_type.ConnectionMethod

	FeeLines        []*shippingtypes.ShippingFeeLine
	CarrierFeeLines []*shippingtypes.ShippingFeeLine
	ExternalID      string
	TotalFee        int
	Coupon          string
}

type DeliveryPoint = shipnowtypes.DeliveryPoint

type SyncStates struct {
	SyncAt    time.Time
	TrySyncAt time.Time
	Error     *meta.Error
}

type ShipnowOrderReservationEvent struct {
	meta.EventMeta

	ShipnowFulfillmentID dot.ID
	OrderIDs             []dot.ID
}

type ShipnowOrderChangedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentID dot.ID
	OldOrderIDs          []dot.ID
	OrderIDs             []dot.ID
}

type ShipnowCancelledEvent struct {
	meta.EventMeta

	ShipnowFulfillmentID dot.ID
	OrderIDs             []dot.ID
	ExternalShipnowID    string
	CarrierServiceCode   string
	CancelReason         string
}

type ShipnowValidateConfirmedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentID dot.ID
	OrderIDs             []dot.ID
}

type ShipnowExternalCreatedEvent struct {
	meta.EventMeta

	ShipnowFulfillmentID dot.ID
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
