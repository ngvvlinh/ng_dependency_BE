package model

import (
	"encoding/json"
	"time"

	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

// +sqlgen
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
	// Error: save error logs when process webhook
	Error *model.Error
}
