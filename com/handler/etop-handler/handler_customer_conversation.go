package handler

import (
	"context"
	"fmt"
	"strings"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/com/handler/pgevent"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleFbCustomerConversationEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFbCustomerConversationEvent", l.Object("pgevent", event))
	var history fbmessagingmodel.FbCustomerConversationHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("FbCustomerConversation not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)

	query := &fbmessaging.GetFbCustomerConversationByIDQuery{
		ID: id,
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("FbCustomerConversation not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	fbCustomerConversation := query.Result
	// get avatars
	{
		var externalUserIDs []string
		if fbCustomerConversation.ExternalUserID != "" {
			externalUserIDs = append(externalUserIDs, fbCustomerConversation.ExternalUserID)
		}
		if fbCustomerConversation.ExternalFrom != nil {
			externalUserIDs = append(externalUserIDs, fbCustomerConversation.ExternalFrom.ID)
		}

		listFbExternalUserQuery := &fbusering.ListFbExternalUsersByExternalIDsQuery{
			ExternalIDs: externalUserIDs,
		}
		if err := h.fbuserQuery.Dispatch(ctx, listFbExternalUserQuery); err != nil {
			return mq.CodeStop, err
		}
		fbExternalUsers := listFbExternalUserQuery.Result

		mapExternalUserIDAndImageURl := make(map[string]string)
		for _, fbExternalUser := range fbExternalUsers {
			if fbExternalUser.ExternalInfo != nil {
				mapExternalUserIDAndImageURl[fbExternalUser.ExternalID] = fbExternalUser.ExternalInfo.ImageURL
			}
		}

		if fbCustomerConversation.ExternalUserID != "" {
			fbCustomerConversation.ExternalUserPictureURL = mapExternalUserIDAndImageURl[fbCustomerConversation.ExternalUserID]
		}
		if fbCustomerConversation.ExternalFrom != nil {
			fbCustomerConversation.ExternalFrom.ImageURL = mapExternalUserIDAndImageURl[fbCustomerConversation.ExternalFrom.ID]
		}
	}

	result := convertpb.PbFbCustomerConversationEvent(fbCustomerConversation, event.Op.String())
	queryPage := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: query.Result.ExternalPageID,
	}
	if err := h.fbPagingQuery.Dispatch(ctx, queryPage); err != nil {
		ll.Warn("fb_page not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	result.FbPageID = queryPage.Result.ID
	queryUser := &identity.GetUsersByAccountQuery{
		AccountID: queryPage.Result.ShopID,
	}
	if err := h.indentityQuery.Dispatch(ctx, queryUser); err != nil {
		ll.Warn("user not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	var userIDs []dot.ID
	for _, user := range queryUser.Result {
		userIDs = append(userIDs, user.UserID)
	}
	result.UserIDs = userIDs
	topic := h.prefix + event.Table + "_fabo"
	d, ok := pgevent.TopicMap[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(query.Result.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}

func (h *Handler) HandleFbCustomerConversationFaboEvent(ctx context.Context, event *pgevent.PgEventFabo) (mq.Code, error) {
	title := "fabo/customer_conversation/" + strings.ToLower(event.PgEventCustomerConversation.Op)
	for _, userID := range event.PgEventCustomerConversation.UserIDs {
		eventCustomerConversation := eventstream.Event{
			Type:    title,
			UserID:  userID,
			Payload: event.PgEventCustomerConversation.FbEventCustomerConversation,
		}
		publisher.Publish(eventCustomerConversation)
	}
	return mq.CodeOK, nil
}
