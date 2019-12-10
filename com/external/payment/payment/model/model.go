package model

import (
	"encoding/json"
	"time"

	"etop.vn/api/top/types/etc/payment_provider"
	"etop.vn/api/top/types/etc/payment_state"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenPayment(&Payment{})

type Payment struct {
	ID              dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	PaymentProvider payment_provider.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}
