package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbConversation struct {
	ID                   dot.ID
	ExternalID           string
	FbPageID             dot.ID
	ExternalLink         string
	ExternalMessageCount int
	ExternalUpdatedTime  time.Time
	CreatedTime          time.Time `sq:"create"`
	UpdatedTime          time.Time `sq:"update"`
	DeletedTime          time.Time
}
