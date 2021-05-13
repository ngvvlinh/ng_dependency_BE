package handler

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/setting"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessagingmodel "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/com/fabo/pkg/fbclient"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
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

	// hide comments have "Op" were "insert"
	if event.Op == pgevent.OpInsert {
		go func(externalComment *fbmessaging.FbExternalComment) {
			h.hideComments(ctx, externalComment)
		}(fbExternalComment)
	}

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
	// Have two types of comment: comment on personal page and fan page.
	// With comment on fan page that has external_page_id
	// With comment on personal page that has external_owner_post_id (external_user_id created post)
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
			ExternalID: query.Result.ExternalOwnerPostID,
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

func (h *Handler) hideComments(ctx context.Context, externalComment *fbmessaging.FbExternalComment) {
	externalPageID, externalCommentID := externalComment.ExternalPageID, externalComment.ExternalID
	// ignore comment from page
	if externalComment.ExternalFrom != nil && externalComment.ExternalFrom.ID == externalPageID {
		return
	}

	// get shopSetting
	getFbPageQuery := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: externalPageID,
	}

	if err := h.fbPagingQuery.Dispatch(ctx, getFbPageQuery); err != nil {
		ll.Warn("fb_page not found", l.String("external_id", externalPageID))
		return
	}
	fbExternalPage := getFbPageQuery.Result
	shopID := fbExternalPage.ShopID

	if fbExternalPage.Status != status3.P {
		return
	}

	getShopSettingQuery := &setting.GetShopSettingQuery{
		ShopID: shopID,
	}
	if err := h.settingQuery.Dispatch(ctx, getShopSettingQuery); err != nil {
		ll.Warn("shop_setting not found", l.Int64("id", shopID.Int64()))
		return
	}
	shopSetting := getShopSettingQuery.Result

	// Check field HideAllComments into shopSetting
	if !shopSetting.HideAllComments.Valid || !shopSetting.HideAllComments.Bool {
		return
	}

	// Get accessToken of page
	getFbPageInternalQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := h.fbPagingQuery.Dispatch(ctx, getFbPageInternalQuery); err != nil {
		ll.Warn("fb_page_internal not found", l.String("external_id", externalPageID))
		return
	}
	accessToken := getFbPageInternalQuery.Result.Token

	// Hide comment
	hideCommentRequest := &fbclient.HideOrUnHideCommentRequest{
		AccessToken: accessToken,
		PageID:      externalPageID,
		CommentID:   externalCommentID,
		IsHidden:    true,
	}
	if _, err := h.fbClient.CallAPIHideAndUnHideComment(hideCommentRequest); err != nil {
		ll.Warn("hide comment error", l.String("external_id", externalCommentID))
		return
	}
}
