package handler

import (
	"context"
	"fmt"
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/eventhandler/notifier"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	oneSignalNotifier *notifier.Notifier
)

// TODO(vu): remove this
func Init(n *notifier.Notifier) {
	oneSignalNotifier = n
}

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
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func SendNotification(ctx context.Context, notify *notifiermodel.Notification) error {
	if notify == nil {
		return nil
	}

	disable, err := isDisableTopicNotify(ctx, notify)
	if err != nil {
		return err
	}
	if disable {
		return nil
	}

	return sendToOneSignal(ctx, notify)
}

func isDisableTopicNotify(ctx context.Context, notify *notifiermodel.Notification) (bool, error) {
	userNotifySetting, err := userNotifySettingStore(ctx).ByUserID(notify.UserID).GetUserNotifySetting()
	if err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return false, nil
		}
		return false, err
	}
	for _, topic := range userNotifySetting.DisableTopics {
		if topic == notify.TopicType {
			return true, nil
		}
	}
	return false, nil
}

func sendToOneSignal(ctx context.Context, noti *notifiermodel.Notification) error {
	if noti.UserID == 0 {
		accUsers, err := accountUserStore(ctx).ListAccountUsers()
		if err != nil {
			return err
		}

		for _, accUser := range accUsers {
			err := _sendToOneSignal(ctx, accUser.UserID, noti)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return _sendToOneSignal(ctx, noti.UserID, noti)
}

func _sendToOneSignal(ctx context.Context, userID dot.ID, notify *notifiermodel.Notification) error {
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
	deviceIDs := filterDevicesByConfig(devices, notify.AccountID)

	data := notifiermodel.PrepareNotiData(&notifiermodel.NotiDataAddition{
		Entity:   notify.Entity,
		EntityID: notify.EntityID.String(),
		NotiID:   notify.ID.String(),
		ShopID:   notify.AccountID.String(),
		MetaData: notify.MetaData,
	})

	webUrl := buildNotificationURL(notify)
	now := time.Now()
	cmd := &notifiermodel.SendNotificationCommand{
		Request: &notifiermodel.CreateNotificationRequest{
			Title:             notify.Title,
			Content:           notify.Message,
			Data:              data,
			ExternalDeviceIDs: deviceIDs,
			WebURL:            webUrl,
		},
	}

	updateNotify := &notifiermodel.Notification{
		ID:                notify.ID,
		SyncStatus:        status3.P,
		ExternalServiceID: notifiermodel.ExternalServiceOneSignalID,
		SyncedAt:          now,
	}
	if err = oneSignalNotifier.CreateNotification(ctx, cmd); err != nil {
		updateNotify.SyncStatus = status3.N
	}

	if cmd.Result.Recipients < len(deviceIDs) {
		handleNotifyError(cmd)
	}

	updateNotify.ExternalNotiID = cmd.Result.ID
	err = notifyStore.UpdateNotification(updateNotify)
	return err
}

// Handle error for OneSignal
// https://documentation.onesignal.com/reference/create-notification#results---create-notification
// "invalid_player_ids" will only return device ids associated with the targeted app_id if you include ids that are not found within the app_id they will be ignored and not returned within this field
// => When this error happens, delete the device in our database
func handleNotifyError(cmd *notifiermodel.SendNotificationCommand) {
	errVal, ok := cmd.Result.Errors.(map[string]interface{})
	if !ok {
		return
	}

	if invalidPlayerIDs, ok := errVal["invalid_player_ids"]; ok {
		ids, ok := invalidPlayerIDs.([]interface{})
		if !ok {
			return
		}

		for _, id := range ids {
			err := deviceStore.DeleteDeviceByExternalID(fmt.Sprintf("%v", id), notifiermodel.ExternalServiceOneSignalID)
			if err != nil {
				ll.Error(fmt.Sprintf("delete device got error: %v", err))
			}
		}
	}
}

func filterDevicesByConfig(devices []*notifiermodel.Device, accountID dot.ID) (deviceIDs []string) {
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
