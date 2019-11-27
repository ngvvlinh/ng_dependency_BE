package model

import (
	"encoding/json"
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShippingProviderWebhook(&ShippingProviderWebhook{})

type ShippingProviderWebhook struct {
	ID               dot.ID
	ShippingProvider string
	Data             json.RawMessage
	ShippingCode     string
	ShippingState    string

	ExternalShippingState    string
	ExternalShippingSubState string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
