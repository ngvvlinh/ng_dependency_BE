package handler

import (
	"context"
	"fmt"

	"o.o/api/main/authorization"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessaging "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

func HandleCommentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	// Only handle on create comment
	if event.Op != pgevent.OpInsert {
		return mq.CodeOK, nil
	}

	var history fbmessaging.FbExternalCommentHistory
	if ok, err := historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fb comment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := history.ID().Int64().Apply(0)

	var cmt fbmessaging.FbExternalComment
	if ok, err := x.Where("id = ?", id).Get(&cmt); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fb comment not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}

	// If comment from page owner, ignore them
	externalFrom := cmt.ExternalFrom
	if cmt.ExternalPageID == externalFrom.ID {
		return mq.CodeIgnore, nil
	}

	cmds := buildCommentNotifyCmds(ctx, &cmt)
	if cmds == nil {
		return mq.CodeIgnore, nil
	}
	if err := createNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func buildFbExternalCommentNotifyMessage(cmt *fbmessaging.FbExternalComment) string {
	if cmt.ExternalAttachment != nil {
		attachment := cmt.ExternalAttachment
		switch {
		case attachment.Type == "sticker":
			return fmt.Sprintf("%s đã comment bằng sticker.", cmt.ExternalFrom.Name)
		case attachment.Type == "photo":
			return fmt.Sprintf("%s đã comment bằng hình ảnh.", cmt.ExternalFrom.Name)
		case attachment.Type == "video_inline":
			return fmt.Sprintf("%s đã comment bằng video.", cmt.ExternalFrom.Name)
		}
	}

	return fmt.Sprintf("%s: %s.", cmt.ExternalFrom.Name, cmt.ExternalMessage)
}

func buildCommentNotifyCmds(ctx context.Context, cmt *fbmessaging.FbExternalComment) []*notifiermodel.CreateNotificationArgs {
	page, err := getPage(cmt.ExternalPageID)
	if err != nil || page == nil {
		return nil
	}
	shopID := page.ShopID
	userIDs, err := filterRecipient(ctx, shopID, notifyTopicRolesMap[TopicFBComment])
	if err != nil || userIDs == nil {
		return nil
	}
	title := fmt.Sprintf("Bình luận mới đến %s", page.ExternalName)
	message := buildFbExternalCommentNotifyMessage(cmt)

	conversation, err := getCustomerConversation(ctx, cmt.ExternalPostID, cmt.ExternalUserID)
	if err != nil {
		return nil
	}
	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     page.ShopID,
		Title:      title,
		Message:    message,
		SendNotify: true,
		Entity:     notifiermodel.NotiFaboComment,
		EntityID:   cmt.ID,
		Meta: map[string]interface{}{
			"conversation_id":        conversation.ID, // TODO(nakhoa17): refactor
			"fb_external_comment_id": cmt.ExternalID,
			"fb_external_post_id":    cmt.ExternalPostID,
		},
		TopicType: TopicFBComment,
	}
	return buildNotifyCmds(args)
}

func getCustomerConversation(ctx context.Context, externalID string, externalUserID string) (*fbmessaging.FbCustomerConversation, error) {
	return customerConversationStore(ctx).
		ExternalID(externalID).
		FbExternalUserID(externalUserID).
		GetFbCustomerConversationDB()
}

func filterRecipient(ctx context.Context, shopID dot.ID, roles []authorization.Role) ([]dot.ID, error) {
	accUsers, err := accountUserStore(ctx).ByAccountID(shopID).ListAccountUser()
	if err != nil {
		return nil, err
	}

	var res []dot.ID
	for _, acc := range accUsers {
		if doesContainRole(acc.Roles, roles) {
			res = append(res, acc.UserID)
		}
	}

	return res, nil
}

func doesContainRole(userRoles []string, permRoles []authorization.Role) bool {
	for _, uRole := range userRoles {
		for _, pRole := range permRoles {
			if uRole == string(pRole) {
				return true
			}
		}
	}
	return false
}

func HandleMessageEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	// Only handle on create message
	if event.Op != pgevent.OpInsert {
		return mq.CodeOK, nil
	}

	var history fbmessaging.FbExternalMessageHistory
	if ok, err := historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fb message not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := history.ID().Int64().Apply(0)

	var message fbmessaging.FbExternalMessage
	if ok, err := x.Where("id = ?", id).Get(&message); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fb message not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}

	// If comment from page owner, ignore them
	externalFrom := message.ExternalFrom
	if message.ExternalPageID == externalFrom.ID {
		return mq.CodeOK, nil
	}

	cmds := buildMessageNotifyCmds(ctx, &message)
	if cmds == nil {
		return mq.CodeIgnore, nil
	}
	if err := createNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func buildFbExternalMessageNotifyMessage(msg *fbmessaging.FbExternalMessage) string {
	if msg.ExternalSticker != "" {
		return fmt.Sprintf("%s đã gửi cho bạn một sticker.", msg.ExternalFrom.Name)
	}

	if msg.ExternalAttachments != nil || len(msg.ExternalAttachments) > 0 {
		return fmt.Sprintf("%s đã gửi cho bạn một hình ảnh.", msg.ExternalFrom.Name)
	}

	return fmt.Sprintf("%s: %s", msg.ExternalFrom.Name, msg.ExternalMessage)
}

func buildMessageNotifyCmds(ctx context.Context, msg *fbmessaging.FbExternalMessage) []*notifiermodel.CreateNotificationArgs {
	page, err := getPage(msg.ExternalPageID)
	if err != nil || page == nil {
		return nil
	}
	userIDs, err := filterRecipient(ctx, page.ShopID, notifyTopicRolesMap[TopicFBMessage])
	if err != nil || userIDs == nil {
		return nil
	}
	title := fmt.Sprintf("Tin nhắn mới đến %s", page.ExternalName)
	message := buildFbExternalMessageNotifyMessage(msg)

	conversation, err := getCustomerConversation(ctx, msg.ExternalConversationID, msg.ExternalFrom.ID)
	if err != nil {
		return nil
	}
	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     page.ShopID,
		Title:      title,
		Message:    message,
		SendNotify: true,
		Entity:     notifiermodel.NotiFaboMessage,
		EntityID:   msg.ID,
		Meta: map[string]interface{}{
			"conversation_id":             conversation.ID, // TODO(nakhoa17): refactor
			"fb_external_message_id":      msg.ExternalID,
			"fb_external_conversation_id": msg.ExternalConversationID,
		},
		TopicType: TopicFBMessage,
	}
	return buildNotifyCmds(args)
}
