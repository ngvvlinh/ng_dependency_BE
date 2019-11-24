package model

import (
	"encoding/json"
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenPayment(&Payment{})

type Payment struct {
	ID              dot.ID
	Amount          int
	Status          int
	State           string
	PaymentProvider string
	ExternalTransID string
	ExternalData    json.RawMessage
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}
