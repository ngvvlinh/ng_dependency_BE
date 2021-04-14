package payment

import (
	"encoding/json"
	"time"

	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/capi/dot"
)

// +gen:event:topic=event/payment

type Payment struct {
	ID              dot.ID
	ShopID          dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	PaymentProvider payment_provider.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PaymentStatusUpdatedEvent struct {
	ID            dot.ID
	ShopID        dot.ID
	PaymentStatus status4.Status
}
