package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/notifier/model"
	"etop.vn/backend/pkg/pgevent"
)

func HandleNotificationEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	if event.Op != pgevent.OpInsert {
		return mq.CodeIgnore, nil
	}
	id := event.ID
	var noti model.Notification
	if ok, err := xNotifier.Where("id = ?", id).Get(&noti); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Notification not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}
	if !noti.SendNotification {
		return mq.CodeIgnore, nil
	}
	if err := SendNotification(ctx, &noti); err != nil {
		return mq.CodeRetry, nil
	}

	return mq.CodeOK, nil
}

func SendNotification(ctx context.Context, noti *model.Notification) error {
	if noti == nil {
		return nil
	}
	if err := sendToOneSignal(ctx, noti); err != nil {
		return err
	}
	return nil
}

func sendToOneSignal(ctx context.Context, noti *model.Notification) error {
	args := &model.GetDevicesArgs{
		AccountID:         noti.AccountID,
		ExternalServiceID: model.ExternalServiceOneSignalID,
	}
	deviceIDs, err := deviceStore.GetExternalDeviceIDs(args)
	if err != nil {
		return err
	}
	if len(deviceIDs) == 0 {
		return nil
	}
	data := model.PrepareNotiData(model.NotiDataAddition{
		Entity:   noti.Entity,
		EntityID: strconv.FormatInt(noti.EntityID, 10),
		NotiID:   strconv.FormatInt(noti.ID, 10),
	})

	webUrl := buildNotificationURL(noti)
	now := time.Now()
	cmd := &model.SendNotificationCommand{
		Request: &model.CreateNotificationRequest{
			Title:             noti.Title,
			Content:           noti.Message,
			Data:              data,
			ExternalDeviceIDs: deviceIDs,
			WebURL:            webUrl,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	// Update external_noti_id and sync_status
	updateNoti := &model.Notification{
		ID:                noti.ID,
		SyncStatus:        etopmodel.S3Positive,
		ExternalNotiID:    cmd.Result.ID,
		ExternalServiceID: model.ExternalServiceOneSignalID,
		SyncedAt:          now,
	}
	if err := notiStore.UpdateNotification(updateNoti); err != nil {
		return err
	}
	return nil
}

func buildNotificationURL(noti *model.Notification) string {
	return fmt.Sprintf("%v/notifications/%v", cm.MainSiteBaseURL(), noti.ID)
}
