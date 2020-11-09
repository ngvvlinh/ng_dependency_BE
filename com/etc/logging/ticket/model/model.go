package model

import (
	"encoding/json"
	"time"

	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

// +sqlgen
type TicketProviderWebhook struct {
	ID             dot.ID
	TicketProvider string
	Data           json.RawMessage
	ExternalStatus string
	ExternalType   string
	ClientID       string
	CreatedAt      time.Time `sq:"create"`
	// Error: save error logs when process webhook
	Error        *model.Error
	ConnectionID dot.ID
}
