package api

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop-handler/webhook/sender"

	pbcm "etop.vn/backend/pb/common"
	wraphandler "etop.vn/backend/wrapper/services/handler"
)

var whsender *sender.WebhookSender

func init() {
	bus.AddHandlers("handler",
		VersionInfo,
		ResetState,
	)
}

func Init(s *sender.WebhookSender) {
	whsender = s
}

func VersionInfo(ctx context.Context, q *wraphandler.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service:   "etop-event-handler",
		Version:   "0.1",
		UpdatedAt: nil,
	}
	return nil
}

func ResetState(ctx context.Context, q *wraphandler.ResetStateEndpoint) error {
	if q.AccountId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "invalid account_id")
	}

	err := whsender.ResetState(q.AccountId)
	q.Result = &pbcm.Empty{}
	return err
}
