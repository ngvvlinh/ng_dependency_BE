package payment

import (
	"encoding/json"
	"time"

	"etop.vn/api/top/types/etc/payment_provider"
	"etop.vn/api/top/types/etc/payment_state"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/capi/dot"
)

type Payment struct {
	ID              dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	PaymentProvider payment_provider.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
