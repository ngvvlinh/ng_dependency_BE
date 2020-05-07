package handler

import (
	"context"
	"fmt"
	"strings"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/com/handler/pgevent"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/common/l"
)

func (h *Handler) HandleFbConversationEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFbConversationEvent", l.Object("pgevent", event))
	var history fbmessagingmodel.FbExternalConversationHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("FbConversation not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)

	query := &fbmessaging.GetFbExternalConversationByIDQuery{
		ID: id,
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("FbConversation not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	result := convertpb.PbFbExternalConversationEvent(query.Result, event.Op.String())
	queryPage := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: query.Result.ExternalPageID,
	}
	if err := h.fbPagingQuery.Dispatch(ctx, queryPage); err != nil {
		ll.Warn("fb_conversation not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	result.FbPageID = queryPage.Result.ID
	topic := h.prefix + event.Table + "_fabo"
	d, ok := pgevent.TopicMap[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(query.Result.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}

func (h *Handler) HandleFbConversationFaboEvent(ctx context.Context, event *pgevent.PgEventFabo) (mq.Code, error) {
	title := "fabo/conversation/" + strings.ToLower(event.PgEventConversation.Op)
	eventComment := eventstream.Event{
		Type:    title,
		UserID:  event.PgEventConversation.FbPageID,
		Payload: event.PgEventConversation.FbEventConversation,
	}
	publisher.Publish(eventComment)
	return mq.CodeOK, nil
}
