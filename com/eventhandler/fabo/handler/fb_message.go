package handler

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/mq"
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
	//get avatars
	{
		var externalUserIDs []string
		if fbMessage.ExternalFrom != nil {
			externalUserIDs = append(externalUserIDs, fbMessage.ExternalFrom.ID)
		}
		for _, externalTo := range fbMessage.ExternalTo {
			externalUserIDs = append(externalUserIDs, externalTo.ID)
		}

		listFbExternalUsersQuery := &fbusering.ListFbExternalUsersByExternalIDsQuery{
			ExternalIDs: externalUserIDs,
		}
		if err := h.fbuserQuery.Dispatch(ctx, listFbExternalUsersQuery); err != nil {
			return mq.CodeStop, err
		}

		fbExternalUsers := listFbExternalUsersQuery.Result
		mapExternalUserIDAndImageURL := make(map[string]string)
		for _, fbExternalUser := range fbExternalUsers {
			if fbExternalUser.ExternalInfo != nil {
				mapExternalUserIDAndImageURL[fbExternalUser.ExternalID] = fbExternalUser.ExternalInfo.ImageURL
			}
		}

		if fbMessage.ExternalFrom != nil {
			fbMessage.ExternalFrom.ImageURL = mapExternalUserIDAndImageURL[fbMessage.ExternalFrom.ID]
		}
		for _, externalTo := range fbMessage.ExternalTo {
			externalTo.ImageURL = mapExternalUserIDAndImageURL[externalTo.ID]
		}
	}

	result := PbFbExternalMessageEvent(fbMessage, event.Op.String())
	queryPage := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: query.Result.ExternalPageID,
	}
	if err := h.fbPagingQuery.Dispatch(ctx, queryPage); err != nil {
		ll.Warn("fb_page not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	result.FbPageID = queryPage.Result.ID
	result.ShopID = queryPage.Result.ShopID

	topic := h.prefix + event.Table + "_fabo" // TODO(vu): refactor
	d, ok := mapTopics[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(fbMessage.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}
