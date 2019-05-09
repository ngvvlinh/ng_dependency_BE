package model_log

import (
	"encoding/json"
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate bash -c "rm derived.gen.go || true"
//go:generate ../../../scripts/derive.sh

var _ = sqlgenShippingProviderWebhook(&ShippingProviderWebhook{})

type ShippingProviderWebhook struct {
	ID                       int64
	ShippingProvider         model.ShippingProvider
	Data                     json.RawMessage
	ShippingCode             string
	ShippingState            string
	ExternalShippingState    string
	ExternalShippingSubState string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
