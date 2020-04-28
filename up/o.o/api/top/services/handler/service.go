package handler

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=services/handler

// +apix:path=/handler.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/handler.Webhook
type WebhookService interface {
	ResetState(context.Context, *ResetStateRequest) (*cm.Empty, error)
}
