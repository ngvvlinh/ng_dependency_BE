package payment

import (
	"encoding/json"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/capi/dot"
)

type Payment struct {
	ID              dot.ID
	Amount          int
	Status          etop.Status4
	State           PaymentState
	PaymentProvider PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type (
	PaymentState    string
	PaymentProvider string
	PaymentSource   string
)

var (
	PaymentStateDefault   PaymentState = "default"
	PaymentStateCreated   PaymentState = "created"
	PaymentStatePending   PaymentState = "pending"
	PaymentStateSuccess   PaymentState = "success"
	PaymentStateFailed    PaymentState = "failed"
	PaymentStateUnknown   PaymentState = "unknown"
	PaymentStateCancelled PaymentState = "cancelled"

	PaymentProviderVTPay PaymentProvider = "vtpay"

	PaymentSourceOrder PaymentSource = "order"
)
