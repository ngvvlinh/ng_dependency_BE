package handler

import (
	"context"

	"github.com/Shopify/sarama"

	"o.o/backend/com/eventhandler/handler/intctl"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
)

func (h *Handler) RegisterTo(intctlHandler *intctl.Handler) {
	intctlHandler.Subscribe(intctl.ChannelReloadWebhook, h.handleReloadWebhook)
}

func (h *Handler) handleReloadWebhook(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
	var v intctl.ReloadWebhook
	err := jsonx.Unmarshal(msg.Value, &v)
	if err != nil {
		return mq.CodeStop, nil
	}
	if v.AccountID == 0 {
		ll.Error("webhook/reload: account_id is empty")
		return mq.CodeStop, nil
	}
	if err := h.sender.Reload(ctx, v.AccountID); err != nil {
		ll.Error("webhook/reload: account_id is empty")
		return mq.CodeStop, nil
	}
	return mq.CodeOK, nil
}
