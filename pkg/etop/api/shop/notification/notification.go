package notification

import (
	"context"

	"o.o/api/top/int/etop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	notistore "o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type NotificationService struct {
	session.Session

	NotificationStore *notistore.NotificationStore
	DeviceStore       *notistore.DeviceStore
}

func (s *NotificationService) Clone() api.NotificationService { res := *s; return &res }

func (s *NotificationService) CreateDevice(ctx context.Context, q *etop.CreateDeviceRequest) (*etop.Device, error) {
	cmd := &notimodel.CreateDeviceArgs{
		UserID:           s.SS.Claim().UserID,
		AccountID:        s.SS.Shop().ID,
		DeviceID:         q.DeviceId,
		DeviceName:       q.DeviceName,
		ExternalDeviceID: q.ExternalDeviceId,
	}
	device, err := s.DeviceStore.CreateDevice(cmd)
	if err != nil {
		return nil, err
	}
	result := convertpb.PbDevice(device)
	return result, nil
}

func (s *NotificationService) DeleteDevice(ctx context.Context, q *etop.DeleteDeviceRequest) (*pbcm.DeletedResponse, error) {
	device := &notimodel.Device{
		DeviceID:         q.DeviceId,
		ExternalDeviceID: q.ExternalDeviceId,
		AccountID:        s.SS.Shop().ID,
		UserID:           s.SS.Claim().UserID,
	}
	if err := s.DeviceStore.DeleteDevice(device); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{
		Deleted: 1,
	}
	return result, nil
}

func (s *NotificationService) GetNotification(ctx context.Context, q *pbcm.IDRequest) (*etop.Notification, error) {
	query := &notimodel.GetNotificationArgs{
		AccountID: s.SS.Shop().ID,
		ID:        q.Id,
	}
	noti, err := s.NotificationStore.GetNotification(query)
	if err != nil {
		return nil, err
	}
	result := convertpb.PbNotification(noti)
	return result, nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, q *etop.GetNotificationsRequest) (*etop.NotificationsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &notimodel.GetNotificationsArgs{
		Paging:    paging,
		AccountID: s.SS.Shop().ID,
	}
	if q.Filter != nil {
		query.Filter = &notimodel.GetNotificationFilterArgs{Entity: q.Filter.Entity}
	}
	notis, err := s.NotificationStore.GetNotifications(query)
	if err != nil {
		return nil, err
	}
	result := &etop.NotificationsResponse{
		Notifications: convertpb.PbNotifications(notis),
		Paging:        cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *NotificationService) UpdateNotifications(ctx context.Context, q *etop.UpdateNotificationsRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &notimodel.UpdateNotificationsArgs{
		IDs:    q.Ids,
		IsRead: q.IsRead,
	}
	if err := s.NotificationStore.UpdateNotifications(cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: len(q.Ids),
	}
	return result, nil
}
