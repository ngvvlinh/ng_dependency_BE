package api

import (
	"context"

	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/eventhandler/webhook/sender"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
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

func (s *MiscService) Clone() *MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop-event-handler",
		Version: "0.1",
	}
	return nil
}

func (s *WebhookService) Clone() *WebhookService {
	res := *s
	return &res
}

func (s *WebhookService) ResetState(ctx context.Context, q *ResetStateEndpoint) error {
	if q.AccountId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "invalid account_id")
	}

	err := whsender.ResetState(q.AccountId)
	q.Result = &pbcm.Empty{}
	return err
}
