package model

import (
	"encoding/json"
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

type PaymentAction string

var (
	PaymentActionValidate PaymentAction = "validate"
	PaymentActionResult   PaymentAction = "result"
)

var _ = sqlgenPayment(&Payment{})

type Payment struct {
	ID   dot.ID
	Data json.RawMessage
	// Mã từ eTop gửi sang đối tác
	OrderID         string
	PaymentProvider string
	Action          PaymentAction
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}
