package handler

import (
	"context"
	"fmt"

	cm "o.o/api/top/types/common"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	fbmessaging "o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/mq"
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

	cmds := buildCommentNotifyCmd(ctx, &cmt)
	if cmds == nil {
		return mq.CodeIgnore, nil
	}
	if err := CreateNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func buildCommentNotifyCmd(ctx context.Context, cmt *fbmessaging.FbExternalComment) []*notifiermodel.CreateNotificationArgs {
	page, err := getPage(cmt.ExternalPageID)
	if err != nil || page == nil {
		return nil
	}
	shopID := page.ShopID
	userIDs, err := getUserIDsWithShopID(ctx, shopID)
	if err != nil || userIDs == nil {
		return nil
	}
	title := fmt.Sprintf("Bình luận mới đến %s", page.ExternalName)
	message := fmt.Sprintf("%s: %s", cmt.ExternalFrom.Name, cmt.ExternalMessage)

	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     page.ShopID,
		Title:      title,
		Message:    message,
		SendNotify: true,
		Entity:     notifiermodel.NotiFaboComment,
		EntityID:   cmt.ID,
		Meta:       cm.Empty{},
	}
	return buildNotifyCmds(args)
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

	cmds := buildMessageNotiCmd(ctx, &message)
	if cmds == nil {
		return mq.CodeIgnore, nil
	}
	if err := CreateNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}
	return mq.CodeOK, nil
}

func buildMessageNotiCmd(ctx context.Context, msg *fbmessaging.FbExternalMessage) []*notifiermodel.CreateNotificationArgs {
	page, err := getPage(msg.ExternalPageID)
	if err != nil || page == nil {
		return nil
	}
	userIDs, err := getUserIDsWithShopID(ctx, page.ShopID)
	if err != nil || userIDs == nil {
		return nil
	}
	title := fmt.Sprintf("Tin nhắn mới đến %s", page.ExternalName)
	message := fmt.Sprintf("%s: %s", msg.ExternalFrom.Name, msg.ExternalMessage)

	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     page.ShopID,
		Title:      title,
		Message:    message,
		SendNotify: true,
		Entity:     notifiermodel.NotiFaboMessage,
		EntityID:   msg.ID,
		Meta:       cm.Empty{},
	}
	return buildNotifyCmds(args)
}
