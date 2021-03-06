package admin

import (
	"context"

	"o.o/api/top/int/admin"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	notistore "o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/pkg/etop/authorize/session"
)

type NotificationService struct {
	session.Session

	NotificationStore *notistore.NotificationStore
}

func (s *NotificationService) Clone() admin.NotificationService {
	res := *s
	return &res
}

func (s *NotificationService) CreateNotifications(ctx context.Context, q *admin.CreateNotificationsRequest) (*admin.CreateNotificationsResponse, error) {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       q.AccountIds,
		Title:            q.Title,
		Message:          q.Message,
		EntityID:         q.EntityId,
		Entity:           q.Entity,
		SendAll:          q.SendAll,
		SendNotification: true,
	}
	created, errored, err := s.NotificationStore.CreateNotifications(cmd)
	if err != nil {
		return nil, err
	}
	result := &admin.CreateNotificationsResponse{
		Created: created,
		Errored: errored,
	}
	return result, nil
}
