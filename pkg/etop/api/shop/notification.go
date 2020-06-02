package shop

import (
	"context"

	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/sqlstore"
)

type NotificationService struct{}

func (s *NotificationService) Clone() *NotificationService { res := *s; return &res }

func (s *NotificationService) CreateDevice(ctx context.Context, q *CreateDeviceEndpoint) error {
	cmd := &notimodel.CreateDeviceArgs{
		UserID:           q.Context.UserID,
		AccountID:        q.Context.Shop.ID,
		DeviceID:         q.DeviceId,
		DeviceName:       q.DeviceName,
		ExternalDeviceID: q.ExternalDeviceId,
	}
	device, err := sqlstore.CreateDevice(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = convertpb.PbDevice(device)
	return nil
}

func (s *NotificationService) DeleteDevice(ctx context.Context, q *DeleteDeviceEndpoint) error {
	device := &notimodel.Device{
		DeviceID:         q.DeviceId,
		ExternalDeviceID: q.ExternalDeviceId,
		AccountID:        q.Context.Shop.ID,
		UserID:           q.Context.UserID,
	}
	if err := sqlstore.DeleteDevice(ctx, device); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{
		Deleted: 1,
	}
	return nil
}

func (s *NotificationService) GetNotification(ctx context.Context, q *GetNotificationEndpoint) error {
	query := &notimodel.GetNotificationArgs{
		AccountID: q.Context.Shop.ID,
		ID:        q.Id,
	}
	noti, err := sqlstore.GetNotification(ctx, query)
	if err != nil {
		return err
	}
	q.Result = convertpb.PbNotification(noti)
	return nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, q *GetNotificationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &notimodel.GetNotificationsArgs{
		Paging:    paging,
		AccountID: q.Context.Shop.ID,
	}
	notis, err := sqlstore.GetNotifications(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &etop.NotificationsResponse{
		Notifications: convertpb.PbNotifications(notis),
		Paging:        cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *NotificationService) UpdateNotifications(ctx context.Context, q *UpdateNotificationsEndpoint) error {
	cmd := &notimodel.UpdateNotificationsArgs{
		IDs:    q.Ids,
		IsRead: q.IsRead,
	}
	if err := sqlstore.UpdateNotifications(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: len(q.Ids),
	}
	return nil
}
