package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	notifiermodel "etop.vn/backend/com/handler/notifier/model"
	"etop.vn/backend/com/handler/pgevent"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/mq"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

func HandleNotificationEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	if event.Op != pgevent.OpInsert {
		return mq.CodeIgnore, nil
	}
	id := event.ID
	var noti notifiermodel.Notification
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

func SendNotification(ctx context.Context, noti *notifiermodel.Notification) error {
	if noti == nil {
		return nil
	}
	if err := sendToOneSignal(ctx, noti); err != nil {
		return err
	}
	return nil
}

func sendToOneSignal(ctx context.Context, noti *notifiermodel.Notification) error {
	cmdUser := &etopmodel.GetAccountUserQuery{
		AccountID:       noti.AccountID,
		FindByAccountID: true,
	}
	if err := bus.Dispatch(ctx, cmdUser); err != nil {
		return err
	}
	args := &notifiermodel.GetDevicesArgs{
		UserID:            cmdUser.Result.UserID,
		ExternalServiceID: notifiermodel.ExternalServiceOneSignalID,
	}
	devices, err := deviceStore.GetDevices(args)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return nil
	}
	deviceIDs := FilterDevicesByConfig(devices, noti.AccountID)

	data := notifiermodel.PrepareNotiData(notifiermodel.NotiDataAddition{
		Entity:   noti.Entity,
		EntityID: strconv.FormatInt(noti.EntityID, 10),
		NotiID:   strconv.FormatInt(noti.ID, 10),
		ShopID:   strconv.FormatInt(noti.AccountID, 10),
		MetaData: noti.MetaData,
	})

	webUrl := buildNotificationURL(noti)
	now := time.Now()
	cmd := &notifiermodel.SendNotificationCommand{
		Request: &notifiermodel.CreateNotificationRequest{
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
	// UpdateInfo external_noti_id and sync_status
	updateNoti := &notifiermodel.Notification{
		ID:                noti.ID,
		SyncStatus:        etopmodel.S3Positive,
		ExternalNotiID:    cmd.Result.ID,
		ExternalServiceID: notifiermodel.ExternalServiceOneSignalID,
		SyncedAt:          now,
	}
	if err := notiStore.UpdateNotification(updateNoti); err != nil {
		return err
	}
	return nil
}

func FilterDevicesByConfig(devices []*notifiermodel.Device, accountID int64) (deviceIDs []string) {
	for _, device := range devices {
		if device.Config == nil {
			deviceIDs = append(deviceIDs, device.ExternalDeviceID)
			continue
		}
		if device.Config.Mute {
			continue
		}
		if device.Config.SubcribeAllShop || cm.ContainInt64(device.Config.SubcribeShopIDs, accountID) {
			deviceIDs = append(deviceIDs, device.ExternalDeviceID)
		}
	}
	return
}

func buildNotificationURL(noti *notifiermodel.Notification) string {
	return fmt.Sprintf("%v/notifications/%v?shop_id=%v", cm.MainSiteBaseURL(), noti.ID, noti.AccountID)
}
