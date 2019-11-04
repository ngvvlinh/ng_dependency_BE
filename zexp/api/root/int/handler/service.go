package handler

import (
	"context"

	cm "etop.vn/backend/pb/common"
	handler "etop.vn/backend/pb/services/handler"
)

// +gen:apix

// +apix:path=/handler.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/handler.Webhook
type WebhookAPI interface {
	ResetState(context.Context, *handler.ResetStateRequest) (*cm.Empty, error)
}
