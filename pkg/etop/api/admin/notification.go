package admin

import (
	"context"

	"o.o/api/top/int/admin"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/pkg/etop/sqlstore"
)

type NotificationService struct{}

func (s *NotificationService) Clone() *NotificationService {
	res := *s
	return &res
}

func (s *NotificationService) CreateNotifications(ctx context.Context, q *CreateNotificationsEndpoint) error {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       q.AccountIds,
		Title:            q.Title,
		Message:          q.Message,
		EntityID:         q.EntityId,
		Entity:           q.Entity,
		SendAll:          q.SendAll,
		SendNotification: true,
	}
	created, errored, err := sqlstore.CreateNotifications(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = &admin.CreateNotificationsResponse{
		Created: created,
		Errored: errored,
	}
	return nil
}
