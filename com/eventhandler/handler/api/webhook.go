package api

import (
	"context"

	api "o.o/api/top/services/handler"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/eventhandler/webhook/sender"
	cm "o.o/backend/pkg/common"
)

type WebhookService struct {
	Sender *sender.WebhookSender
}

func (s *WebhookService) Clone() api.WebhookService {
	res := *s
	return &res
}

func (s *WebhookService) ResetState(ctx context.Context, q *api.ResetStateRequest) (*pbcm.Empty, error) {
	if q.AccountId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "invalid account_id")
	}

	err := s.Sender.ResetState(q.AccountId)
	return &pbcm.Empty{}, err
}
