package sqlstore

import (
	"context"

	notimodel "o.o/backend/com/handler/notifier/model"
)

func CreateDevice(ctx context.Context, args *notimodel.CreateDeviceArgs) (*notimodel.Device, error) {
	device, err := deviceStore.CreateDevice(args)
	return device, err
}

func DeleteDevice(ctx context.Context, device *notimodel.Device) error {
	return deviceStore.DeleteDevice(device)
}

func GetNotification(ctx context.Context, args *notimodel.GetNotificationArgs) (noti *notimodel.Notification, err error) {
	return notificationStore.GetNotification(args)
}

func GetNotifications(ctx context.Context, args *notimodel.GetNotificationsArgs) (notis []*notimodel.Notification, err error) {
	notis, err = notificationStore.GetNotifications(args)
	return
}

func UpdateNotifications(ctx context.Context, args *notimodel.UpdateNotificationsArgs) error {
	return notificationStore.UpdateNotifications(args)
}

func CreateNotifications(ctx context.Context, args *notimodel.CreateNotificationsArgs) (int, int, error) {
	return notificationStore.CreateNotifications(args)
}
