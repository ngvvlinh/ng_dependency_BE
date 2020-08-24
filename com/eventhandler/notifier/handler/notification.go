package handler

import (
	"context"
	"fmt"
	"time"

	"o.o/api/top/types/etc/status3"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
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
	if noti.UserID == 0 {
		userIds, err := getUserIDsWithShopID(ctx, noti.AccountID)
		if err != nil {
			return err
		}
		for _, userID := range userIds {
			err := _sendToOneSignal(ctx, userID, noti)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return _sendToOneSignal(ctx, noti.UserID, noti)
}

func _sendToOneSignal(ctx context.Context, userID dot.ID, noti *notifiermodel.Notification) error {
	args := &notifiermodel.GetDevicesArgs{
		UserID:            userID,
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

	data := notifiermodel.PrepareNotiData(&notifiermodel.NotiDataAddition{
		Entity:   noti.Entity,
		EntityID: noti.EntityID.String(),
		NotiID:   noti.ID.String(),
		ShopID:   noti.AccountID.String(),
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
		SyncStatus:        status3.P,
		ExternalNotiID:    cmd.Result.ID,
		ExternalServiceID: notifiermodel.ExternalServiceOneSignalID,
		SyncedAt:          now,
	}
	if err := notiStore.UpdateNotification(updateNoti); err != nil {
		return err
	}
	return nil
}

func FilterDevicesByConfig(devices []*notifiermodel.Device, accountID dot.ID) (deviceIDs []string) {
	for _, device := range devices {
		if device.Config == nil {
			deviceIDs = append(deviceIDs, device.ExternalDeviceID)
			continue
		}
		if device.Config.Mute {
			continue
		}
		if device.Config.SubcribeAllShop || cm.IDsContain(device.Config.SubcribeShopIDs, accountID) {
			deviceIDs = append(deviceIDs, device.ExternalDeviceID)
		}
	}
	return
}

func buildNotificationURL(noti *notifiermodel.Notification) string {
	return fmt.Sprintf("%v/notifications/%v?shop_id=%v", cm.MainSiteBaseURL(), noti.ID, noti.AccountID)
}
