package shipnow

import (
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
)

type ShipnowFulfillment = shipnowv1.ShipnowFulfillment
type DeliveryPoint = shipnowtypes.DeliveryPoint // re-export

type ShipnowEvent = shipnowv1.ShipnowEvent
type IsEventData = shipnowv1.IsEventData
type EventType = shipnowv1.EventType

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

const Ahamove = "ahamove"

type ShipnowOrderReservationEvent = shipnowv1.ShipnowOrderReservationEvent
type ShipnowOrderChangedEvent = shipnowv1.ShipnowOrderChangedEvent
type ShipnowCancelledEvent = shipnowv1.ShipnowCancelledEvent
type ShipnowValidatedEvent = shipnowv1.ShipnowValidatedEvent
