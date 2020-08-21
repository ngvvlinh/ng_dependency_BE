package webhook

import (
	"context"

	"o.o/api/top/external/types"
	"o.o/backend/com/eventhandler/handler/intctl"
	"o.o/backend/com/eventhandler/webhook/sender"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	callbackmodel "o.o/backend/pkg/etc/xmodel/callback/model"
	callbackstore "o.o/backend/pkg/etc/xmodel/callback/sqlstore"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
)

type Producer mq.Producer

type Service struct {
	db         *cmsql.Database
	producer   mq.Producer
	redisStore redis.Store
}

func New(db com.MainDB, p Producer, r redis.Store) *Service {
	return &Service{
		db:         db,
		producer:   p,
		redisStore: r,
	}
}

func (s *Service) CreateWebhook(ctx context.Context, accountID dot.ID, r *types.CreateWebhookRequest) (*types.Webhook, error) {
	n, err := callbackstore.Webhook(ctx, s.db).AccountID(accountID).Count()
	if err != nil {
		return nil, err
	}
	if n >= 5 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bạn đã tạo quá nhiều webhook. Vui lòng xoá webhook cũ để tạo lại.")
	}

	item := convertpb.CreateWebhookRequestToModel(r, accountID)
	err = callbackstore.Webhook(ctx, s.db).Create(item)
	if err != nil {
		return nil, err
	}

	item, err = callbackstore.Webhook(ctx, s.db).ID(item.ID).AccountID(accountID).Get()
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

	err := callbackstore.Webhook(ctx, s.db).ID(r.Id).SoftDelete()
	if err != nil {
		return nil, err
	}

	items, err := callbackstore.Webhook(ctx, s.db).AccountID(accountID).List()
	if err != nil {
		return nil, err
	}
	resp := &types.WebhooksResponse{
		Webhooks: convertpb.PbWebhooks(items, s.loadWebhookStates(items)),
	}
	return resp, nil
}

func (s *Service) GetWebhooks(ctx context.Context, accountID dot.ID) (*types.WebhooksResponse, error) {
	items, err := callbackstore.Webhook(ctx, s.db).AccountID(accountID).List()
	resp := &types.WebhooksResponse{
		Webhooks: convertpb.PbWebhooks(items, s.loadWebhookStates(items)),
	}
	return resp, err
}

func (s *Service) loadWebhookStates(webhooks []*callbackmodel.Webhook) []sender.WebhookStates {
	res := make([]sender.WebhookStates, len(webhooks))
	for i, item := range webhooks {
		res[i] = sender.LoadWebhookStates(s.redisStore, item.ID)
	}
	return res
}
