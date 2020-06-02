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
	for _, userID := range event.PgEventCustomerConversation.UserIDs {
		eventCustomerConversation := eventstream.Event{
			Type:    title,
			UserID:  userID,
			Payload: event.PgEventCustomerConversation.FbEventCustomerConversation,
		}
		h.publisher.Publish(eventCustomerConversation)
	}
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbConversationFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/conversation/" + strings.ToLower(event.PgEventConversation.Op)
	for _, userID := range event.PgEventConversation.UserIDs {
		eventConversation := eventstream.Event{
			Type:    title,
			UserID:  userID,
			Payload: event.PgEventConversation.FbEventConversation,
		}
		h.publisher.Publish(eventConversation)
	}
	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbCommentFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/comment/" + strings.ToLower(event.PgEventComment.Op)
	for _, userID := range event.PgEventComment.UserIDs {
		eventComment := eventstream.Event{
			Type:    title,
			UserID:  userID,
			Payload: event.PgEventComment.FbEventComment,
		}
		h.publisher.Publish(eventComment)
	}

	return mq.CodeOK, nil
}

func (h *Publisher) HandleFbMessageFaboEvent(ctx context.Context, event *types.FaboEvent) (mq.Code, error) {
	title := "fabo/message/" + strings.ToLower(event.PgEventMessage.Op)
	for _, userID := range event.PgEventMessage.UserIDs {
		eventMessage := eventstream.Event{
			Type:    title,
			UserID:  userID,
			Payload: event.PgEventMessage.FbEventMessage,
		}
		h.publisher.Publish(eventMessage)
	}
	return mq.CodeOK, nil
}
