package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbCustomerConversationSearch struct {
	ID                   dot.ID
	ExternalUserNameNorm string
	ExternalPageID       string
	CreatedAt            time.Time `sq:"create"`
}

// +sqlgen
type FbExternalCommentSearch struct {
	ID                  dot.ID
	ExternalMessageNorm string
	ExternalPageID      string
	ExternalPostID      string
	ExternalUserID      string
	CreatedAt           time.Time `sq:"create"`
}

// +sqlgen
type FbExternalMessageSearch struct {
	ID                     dot.ID
	ExternalMessageNorm    string
	ExternalPageID         string
	ExternalConversationID string
	CreatedAt              time.Time `sq:"create"`
}
