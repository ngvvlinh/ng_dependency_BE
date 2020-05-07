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

func (h *Handler) HandleFbMessageEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFbMessageEvent", l.Object("pgevent", event))
	var history fbmessagingmodel.FbExternalMessageHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("FbMessage not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)

	query := &fbmessaging.GetFbExternalMessageByIDQuery{
		ID: id,
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("fb_message not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	fbMessage := query.Result
	result := convertpb.PbFbExternalMessageEvent(fbMessage, event.Op.String())
	queryPage := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: query.Result.ExternalPageID,
	}
	if err := h.fbPagingQuery.Dispatch(ctx, queryPage); err != nil {
		ll.Warn("fb_message not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	result.FbPageID = queryPage.Result.ID
	topic := h.prefix + event.Table + "_fabo"
	d, ok := pgevent.TopicMap[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(fbMessage.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}

func (h *Handler) HandleFbMessageFaboEvent(ctx context.Context, event *pgevent.PgEventFabo) (mq.Code, error) {
	title := "fabo/message/" + strings.ToLower(event.PgEventMessage.Op)
	eventComment := eventstream.Event{
		Type:    title,
		UserID:  event.PgEventMessage.FbPageID,
		Payload: event.PgEventMessage.FbEventMessage,
	}
	publisher.Publish(eventComment)
	return mq.CodeOK, nil
}
