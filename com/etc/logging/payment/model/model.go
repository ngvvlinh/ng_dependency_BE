package model

import (
	"encoding/json"
	"time"

	"o.o/api/top/types/etc/payment_provider"
	"o.o/capi/dot"
)

type PaymentAction string

var (
	PaymentActionUnknown  PaymentAction = "unknown"
	PaymentActionValidate PaymentAction = "validate"
	PaymentActionResult   PaymentAction = "result"
)

// +sqlgen
type Payment struct {
	ID   dot.ID
	Data json.RawMessage
	// Mã từ eTop gửi sang đối tác
	OrderID         string
	PaymentProvider payment_provider.PaymentProvider
	Action          PaymentAction
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}
