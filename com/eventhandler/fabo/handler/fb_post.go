package handler

import (
	"fmt"

	"golang.org/x/net/context"
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/common/l"
)

func (h *Handler) HandleFbPostEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFbPostEvent", l.Object("pgevent", event))
	var history fbmessagingmodel.FbExternalPostHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("FbPost not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)

	query := &fbmessaging.GetFbExternalPostByIDQuery{
		ID: id,
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("fb_post not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	fbExternalPost := query.Result

	if fbExternalPost.ExternalFrom != nil {
		getFbExternalUserQuery := &fbusering.GetFbExternalUserByExternalIDQuery{
			ExternalID: fbExternalPost.ExternalFrom.ID,
		}
		_err := h.fbuserQuery.Dispatch(ctx, getFbExternalUserQuery)
		if _err != nil && cm.ErrorCode(_err) != cm.NotFound {
			return mq.CodeStop, _err
		}

		fbExternalUser := getFbExternalUserQuery.Result
		if fbExternalUser != nil && fbExternalUser.ExternalInfo != nil {
			fbExternalPost.ExternalFrom.ImageURL = fbExternalUser.ExternalInfo.ImageURL
		}
	}

	result := PbFbExternalPostEvent(fbExternalPost, event.Op.String())
	// Have two types of post: post on personal page and fan page.
	// With post on fan page that has external_page_id
	// With post on personal page that has external_user_id
	if query.Result.ExternalPageID != "" {
		queryPage := &fbpaging.GetFbExternalPageByExternalIDQuery{
			ExternalID: query.Result.ExternalPageID,
		}
		if err := h.fbPagingQuery.Dispatch(ctx, queryPage); err != nil {
			ll.Warn("fb_page not found", l.Int64("rid", event.RID), l.ID("id", id))
			return mq.CodeIgnore, nil
		}
		result.FbPageID = queryPage.Result.ID
		result.ShopID = queryPage.Result.ShopID
	} else {
		queryUserConnected := &fbusering.GetFbExternalUserConnectedByExternalIDQuery{
			ExternalID: query.Result.ExternalUserID,
		}
		if err := h.fbuserQuery.Dispatch(ctx, queryUserConnected); err != nil {
			ll.Warn("fb_external_user_connected not found", l.Int64("rid", event.RID), l.ID("id", id))
			return mq.CodeIgnore, nil
		}

		if queryUserConnected.Result.ShopID == 0 {
			ll.Warn("shop_id is empty", l.Int64("rid", event.RID), l.ID("id", id))
			return mq.CodeIgnore, nil
		}

		result.ShopID = queryUserConnected.Result.ShopID
	}

	topic := h.prefix + event.Table + "_fabo"
	d, ok := mapTopics[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(query.Result.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}
