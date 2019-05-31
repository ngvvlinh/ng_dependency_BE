package shipnow

import (
	"etop.vn/api/main/etop"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
)

type ShipnowFulfillment = shipnowv1.ShipnowFulfillment
type DeliveryPoint = shipnowtypes.DeliveryPoint // re-export

type ShipnowEvent = shipnowv1.ShipnowEvent
type IsEventData = shipnowv1.IsEventData
type EventType = shipnowv1.EventType
type SyncStates = shipnowv1.SyncStates

func NewShipnowEvent(
	correlationID meta.UUID,
	shipnowFulfillmentID int64,
	data IsEventData,
) *ShipnowEvent {
	meta.AutoFill(&correlationID)
	return &ShipnowEvent{
		Id:            0,
		Uuid:          meta.NewUUID(),
		CorrelationId: correlationID,
		Type:          int32(data.GetEnumTag()),
		Data:          &shipnowv1.EventData{Data: data},
	}
}

type CreatedData = shipnowv1.CreatedData
type ConfirmationRequestedData = shipnowv1.ConfirmationRequestedData
type ConfirmationAcceptedData = shipnowv1.ConfirmationAcceptedData
type ConfirmationRejectedData = shipnowv1.ConfirmationRejectedData
type CancellationRequestedData = shipnowv1.CancellationRequestedData
type CancellationAcceptedData = shipnowv1.CancellationAcceptedData
type CancellationRejectedData = shipnowv1.CancellationRejectedData

type ShipnowOrderReservationEvent = shipnowv1.ShipnowOrderReservationEvent
type ShipnowOrderChangedEvent = shipnowv1.ShipnowOrderChangedEvent
type ShipnowCancelledEvent = shipnowv1.ShipnowCancelledEvent
type ShipnowValidateConfirmedEvent = shipnowv1.ShipnowValidatedConfirmedEvent
type ShipnowCreateExternalEvent = shipnowv1.ShipnowCreateExternalEvent

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
