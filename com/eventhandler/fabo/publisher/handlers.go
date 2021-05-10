package publisher

import (
	"context"
	"strings"

	"o.o/backend/com/eventhandler/fabo/types"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/eventstream"
)

func (h *Publisher) HandleFbCustomerConversationFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/customer_conversation/" + strings.ToLower(event.PgEventCustomerConversation.Op)
	eventCustomerConversation := eventstream.Event{
		Type:      title,
		AccountID: event.PgEventCustomerConversation.ShopID,
		Payload:   event.PgEventCustomerConversation.FbEventCustomerConversation,
	}
	h.publisher.Publish(eventCustomerConversation)
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbConversationFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/conversation/" + strings.ToLower(event.PgEventConversation.Op)
	eventConversation := eventstream.Event{
		Type:      title,
		AccountID: event.PgEventConversation.ShopID,
		Payload:   event.PgEventConversation.FbEventConversation,
	}
	h.publisher.Publish(eventConversation)
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbCommentFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/comment/" + strings.ToLower(event.PgEventComment.Op)
	eventComment := eventstream.Event{
		Type:      title,
		AccountID: event.PgEventComment.ShopID,
		Payload:   event.PgEventComment.FbEventComment,
	}
	h.publisher.Publish(eventComment)
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbMessageFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/message/" + strings.ToLower(event.PgEventMessage.Op)
	eventMessage := eventstream.Event{
		Type:      title,
		AccountID: event.PgEventMessage.ShopID,
		Payload:   event.PgEventMessage.FbEventMessage,
	}
	h.publisher.Publish(eventMessage)
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbPostFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/post/" + strings.ToLower(event.PgEventPost.Op)
	eventMessage := eventstream.Event{
		Type:      title,
		AccountID: event.PgEventPost.ShopID,
		Payload:   event.PgEventPost.FbEventPost,
	}
	h.publisher.Publish(eventMessage)
	return mq.CodeOK, nil
}
