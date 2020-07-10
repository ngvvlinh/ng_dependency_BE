package model

import (
	"encoding/json"
	"time"

	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

// +sqlgen
type FbWebhookLog struct {
	ID         dot.ID
	PageID     string
	Type       string
	ExternalID string
	Data       json.RawMessage
	Error      *model.Error
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
}
