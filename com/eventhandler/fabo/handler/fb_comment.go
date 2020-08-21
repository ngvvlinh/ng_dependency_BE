package handler

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleFbCommentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFbCommentEvent", l.Object("pgevent", event))
	var history fbmessagingmodel.FbExternalCommentHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("FbComment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)

	query := &fbmessaging.GetFbExternalCommentByIDQuery{
		ID: id,
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("fb_comment not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	fbExternalComment := query.Result

	if fbExternalComment.ExternalFrom != nil {
		getFbExternalUserQuery := &fbusering.GetFbExternalUserByExternalIDQuery{
			ExternalID: fbExternalComment.ExternalFrom.ID,
		}
		_err := h.fbuserQuery.Dispatch(ctx, getFbExternalUserQuery)
		if _err != nil && cm.ErrorCode(_err) != cm.NotFound {
			return mq.CodeStop, _err
		}

		fbExternalUser := getFbExternalUserQuery.Result
		if fbExternalUser != nil && fbExternalUser.ExternalInfo != nil {
			fbExternalComment.ExternalFrom.ImageURL = fbExternalUser.ExternalInfo.ImageURL
		}
	}

	var fbParentExternalComment *fbmessaging.FbExternalComment
	if fbExternalComment.ExternalParent != nil && fbExternalComment.ExternalParent.ID != "" {
		queryParentComment := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
			ExternalID: fbExternalComment.ExternalParent.ID,
		}
		if err := h.fbMessagingQuery.Dispatch(ctx, queryParentComment); err != nil && cm.ErrorCode(err) != cm.NotFound {
			ll.Warn("fb_parent_comment not found", l.Int64("rid", event.RID), l.String("external_id", fbExternalComment.ExternalParent.ID))
			return mq.CodeIgnore, nil
		}
		fbParentExternalComment = queryParentComment.Result
	}

	result := PbFbExternalCommentEvent(fbExternalComment, fbParentExternalComment, event.Op.String())
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
	d, ok := mapTopics[event.Table]
	if !ok {
		return mq.CodeIgnore, fmt.Errorf("table not found in TopicMap: %v", event.Table)
	}
	partition := int(query.Result.ID.Int64() % int64(d.Partitions))

	h.producer.SendJSON(topic, partition, event.EventKey, result)
	return mq.CodeOK, nil
}
