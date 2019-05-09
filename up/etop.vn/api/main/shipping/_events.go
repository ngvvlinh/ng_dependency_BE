package shipping

import (
	"time"

	"etop.vn/api/meta"
)

type FulfillmentEvent struct {
	// Auto-increment id taken from database.
	ID int64 `json:"id"`

	CorrelationID meta.UUID `json:"correlation_id"`

	CreatedAt time.Time `json:"created_at"`

	OrderID int64 `json:"order_id"`

	FulfillmentID int64 `json:"fulfillment_id"`

	Type EventType `json:"type"`

	Data FulfillmentEventData `json:"data"`
}

func NewFulfillmentEvent(
	correlationID meta.UUID,
	orderID int64,
	fulfillmentID int64,
	data FulfillmentEventData,
) *FulfillmentEvent {
	if correlationID.IsZero() {
		correlationID = meta.NewUUID()
	}
	return &FulfillmentEvent{
		ID:            0,
		CorrelationID: correlationID,
		CreatedAt:     time.Now(),
		OrderID:       orderID,
		FulfillmentID: fulfillmentID,
		Type:          data.eventType(),
		Data:          data,
	}
}

type EventType string

const (
	FulfillmentCreated          EventType = "FulfillmentCreated"
	FulfillmentConfirmed        EventType = "FulfillmentConfirmed"
	FulfillmentConfirmedFailure EventType = "FulfillmentConfirmedFailure"
	FulfillmentCancelled        EventType = "FulfillmentCancelled"
	FulfillmentCancelledFailure EventType = "FulfillmentCancelledFailure"
	FulfillmentCompensated      EventType = "FulfillmentCompensated"
	ShippingStateChanged        EventType = "ShippingStateChanged"
	ShippingFeeChanged          EventType = "ShippingFeeChanged"
	ActualCODAmountChanged      EventType = "ActualCODAmountChanged"
	PartialDelivered            EventType = "PartialDelivered"
	AdminNoteChanged            EventType = "AdminNoteChanged"
	WeightAndSizeChanged        EventType = "WeightAndSizeChanged"
	ChargeableWeightChanged     EventType = "ChargeableWeightChanged"
	PickupAddressChanged        EventType = "PickupAddressChanged"
	ShippingAddressChanged      EventType = "ShippingAddressChanged"
	ReturnAddressChanged        EventType = "ReturnAddressChanged"
	InfoChanged                 EventType = "InfoChanged"

	// TODO: money transaction events
)

type FulfillmentEventData interface {
	eventType() EventType
}

func (_ *FulfillmentCreatedData) eventType() EventType         { return FulfillmentCreated }
func (_ FulfillmentConfirmedData) eventType() EventType        { return FulfillmentConfirmed }
func (_ FulfillmentConfirmedFailureData) eventType() EventType { return FulfillmentConfirmedFailure }
func (_ FulfillmentCancelledData) eventType() EventType        { return FulfillmentCancelled }
func (_ FulfillmentCancelledFailureData) eventType() EventType { return FulfillmentCancelledFailure }
func (_ FulfillmentCompensatedData) eventType() EventType      { return FulfillmentCompensated }
func (_ ShippingStateChangedData) eventType() EventType        { return ShippingStateChanged }
func (_ ShippingFeeChangedData) eventType() EventType          { return ShippingFeeChanged }
func (_ ActualCODAmountChangedData) eventType() EventType      { return ActualCODAmountChanged }
func (_ PartialDeliveredData) eventType() EventType            { return PartialDelivered }
func (_ AdminNoteChangedData) eventType() EventType            { return AdminNoteChanged }
func (_ WeightAndSizeChangedData) eventType() EventType        { return WeightAndSizeChanged }
func (_ ChargeableWeightChangedData) eventType() EventType     { return ChargeableWeightChanged }
func (_ PickupAddressChangedData) eventType() EventType        { return PickupAddressChanged }
func (_ ShippingAddressChangedData) eventType() EventType      { return ShippingAddressChanged }
func (_ ReturnAddressChangedData) eventType() EventType        { return ReturnAddressChanged }
func (_ InfoChangedData) eventType() EventType                 { return InfoChanged }

type FulfillmentCreatedData struct {
}

type FulfillmentConfirmedData struct {
}

type FulfillmentConfirmedFailureData struct {
}

type FulfillmentCancelledData struct {
}

type FulfillmentCancelledFailureData struct {
}

type FulfillmentCompensatedData struct {
}

type ShippingStateChangedData struct {
}

type ShippingFeeChangedData struct {
}

type ActualCODAmountChangedData struct {
}

type PartialDeliveredData struct {
	ActualCODAmount *int32 `json:"actual_cod_amount"`
}

type AdminNoteChangedData struct {
	AdminNote string `json:"admin_note"`
}

type WeightAndSizeChangedData struct {
	Weight *int32 `json:"weight"`
	Length *int32 `json:"length"`
	Width  *int32 `json:"width"`
	Height *int32 `json:"height"`

	ChargeableWeight *int32 `json:"chargeable_weight"`
}

type ChargeableWeightChangedData struct {
	ChargeableWeight *int32 `json:"chargeable_weight"`
}

type PickupAddressChangedData struct {
	Address Address `json:"address"`

	PrevAddress *Address `json:"prev_address"`
}

type ShippingAddressChangedData struct {
	Address Address `json:"address"`

	PrevAddress *Address `json:"prev_address"`
}

type ReturnAddressChangedData struct {
	Address Address `json:"address"`

	PrevAddress *Address `json:"prev_address"`
}

type InfoChangedData struct {
}
