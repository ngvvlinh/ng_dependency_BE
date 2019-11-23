package handler

import (
	"context"

	cm "etop.vn/api/pb/common"
	handler "etop.vn/api/pb/services/handler"
)

// +gen:apix
// +gen:apix:doc-path=services/handler

// +apix:path=/handler.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/handler.Webhook
type WebhookService interface {
	ResetState(context.Context, *handler.ResetStateRequest) (*cm.Empty, error)
}
