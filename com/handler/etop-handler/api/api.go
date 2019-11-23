package api

import (
	"context"

	pbcm "etop.vn/api/pb/common"
	"etop.vn/backend/com/handler/etop-handler/webhook/sender"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
)

var whsender *sender.WebhookSender

func init() {
	bus.AddHandlers("handler",
		miscService.VersionInfo,
		webhookService.ResetState,
	)
}

type MiscService struct{}
type WebhookService struct{}

var miscService = &MiscService{}
var webhookService = &WebhookService{}

func Init(s *sender.WebhookSender) {
	whsender = s
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service:   "etop-event-handler",
		Version:   "0.1",
		UpdatedAt: nil,
	}
	return nil
}

func (s *WebhookService) ResetState(ctx context.Context, q *ResetStateEndpoint) error {
	if q.AccountId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "invalid account_id")
	}

	err := whsender.ResetState(q.AccountId)
	q.Result = &pbcm.Empty{}
	return err
}
