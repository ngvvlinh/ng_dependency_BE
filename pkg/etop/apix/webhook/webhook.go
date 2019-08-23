package webhook

import (
	"context"

	"etop.vn/backend/com/handler/etop-handler/intctl"
	"etop.vn/backend/com/handler/etop-handler/webhook/sender"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"

	pbexternal "etop.vn/backend/pb/external"
)

var producer mq.Producer
var redisStore redis.Store

func Init(p mq.Producer, r redis.Store) {
	producer = p
	redisStore = r
}

func CreateWebhook(ctx context.Context, accountID int64, r *pbexternal.CreateWebhookRequest) (*pbexternal.Webhook, error) {
	n, err := sqlstore.Webhook(ctx).AccountID(accountID).Count()
	if err != nil {
		return nil, err
	}
	if n >= 5 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bạn đã tạo quá nhiều webhook. Vui lòng xoá webhook cũ để tạo lại.")
	}

	item := r.ToModel(accountID)
	err = sqlstore.Webhook(ctx).Create(item)
	if err != nil {
		return nil, err
	}

	item, err = sqlstore.Webhook(ctx).ID(item.ID).AccountID(accountID).Get()
	if err != nil {
		return nil, err
	}

	resp := pbexternal.PbWebhook(item, sender.LoadWebhookStates(redisStore, item.ID))

	event := &intctl.ReloadWebhook{
		AccountID: accountID,
	}
	producer.SendJSON(0, intctl.NewKey(intctl.ChannelReloadWebhook), event)
	return resp, err
}

func DeleteWebhook(ctx context.Context, accountID int64, r *pbexternal.DeleteWebhookRequest) (*pbexternal.WebhooksResponse, error) {
	event := &intctl.ReloadWebhook{
		AccountID: accountID,
	}
	// always send events after deleting webhooks
	defer producer.SendJSON(0, intctl.NewKey(intctl.ChannelReloadWebhook), event)

	err := sqlstore.Webhook(ctx).ID(r.Id).SoftDelete()
	if err != nil {
		return nil, err
	}

	items, err := sqlstore.Webhook(ctx).AccountID(accountID).List()
	resp := &pbexternal.WebhooksResponse{
		Webhooks: pbexternal.PbWebhooks(items, loadWebhookStates(items)),
	}
	return resp, nil
}

func GetWebhooks(ctx context.Context, accountID int64) (*pbexternal.WebhooksResponse, error) {
	items, err := sqlstore.Webhook(ctx).AccountID(accountID).List()
	resp := &pbexternal.WebhooksResponse{
		Webhooks: pbexternal.PbWebhooks(items, loadWebhookStates(items)),
	}
	return resp, err
}

func loadWebhookStates(webhooks []*model.Webhook) []sender.WebhookStates {
	res := make([]sender.WebhookStates, len(webhooks))
	for i, item := range webhooks {
		res[i] = sender.LoadWebhookStates(redisStore, item.ID)
	}
	return res
}
