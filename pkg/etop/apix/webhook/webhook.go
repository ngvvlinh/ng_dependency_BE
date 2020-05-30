package webhook

import (
	"context"

	"o.o/api/top/external/types"
	"o.o/backend/com/handler/etop-handler/intctl"
	"o.o/backend/com/handler/etop-handler/webhook/sender"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type Producer mq.Producer

type Service struct {
	producer   mq.Producer
	redisStore redis.Store
}

func New(p Producer, r redis.Store) *Service {
	return &Service{
		producer:   p,
		redisStore: r,
	}
}

func (s *Service) CreateWebhook(ctx context.Context, accountID dot.ID, r *types.CreateWebhookRequest) (*types.Webhook, error) {
	n, err := sqlstore.Webhook(ctx).AccountID(accountID).Count()
	if err != nil {
		return nil, err
	}
	if n >= 5 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bạn đã tạo quá nhiều webhook. Vui lòng xoá webhook cũ để tạo lại.")
	}

	item := convertpb.CreateWebhookRequestToModel(r, accountID)
	err = sqlstore.Webhook(ctx).Create(item)
	if err != nil {
		return nil, err
	}

	item, err = sqlstore.Webhook(ctx).ID(item.ID).AccountID(accountID).Get()
	if err != nil {
		return nil, err
	}

	resp := convertpb.PbWebhook(item, sender.LoadWebhookStates(s.redisStore, item.ID))

	event := &intctl.ReloadWebhook{
		AccountID: accountID,
	}
	s.producer.SendJSON(0, intctl.NewKey(intctl.ChannelReloadWebhook), event)
	return resp, err
}

func (s *Service) DeleteWebhook(ctx context.Context, accountID dot.ID, r *types.DeleteWebhookRequest) (*types.WebhooksResponse, error) {
	event := &intctl.ReloadWebhook{
		AccountID: accountID,
	}
	// always send events after deleting webhooks
	defer s.producer.SendJSON(0, intctl.NewKey(intctl.ChannelReloadWebhook), event)

	err := sqlstore.Webhook(ctx).ID(r.Id).SoftDelete()
	if err != nil {
		return nil, err
	}

	items, err := sqlstore.Webhook(ctx).AccountID(accountID).List()
	if err != nil {
		return nil, err
	}
	resp := &types.WebhooksResponse{
		Webhooks: convertpb.PbWebhooks(items, s.loadWebhookStates(items)),
	}
	return resp, nil
}

func (s *Service) GetWebhooks(ctx context.Context, accountID dot.ID) (*types.WebhooksResponse, error) {
	items, err := sqlstore.Webhook(ctx).AccountID(accountID).List()
	resp := &types.WebhooksResponse{
		Webhooks: convertpb.PbWebhooks(items, s.loadWebhookStates(items)),
	}
	return resp, err
}

func (s *Service) loadWebhookStates(webhooks []*model.Webhook) []sender.WebhookStates {
	res := make([]sender.WebhookStates, len(webhooks))
	for i, item := range webhooks {
		res[i] = sender.LoadWebhookStates(s.redisStore, item.ID)
	}
	return res
}
