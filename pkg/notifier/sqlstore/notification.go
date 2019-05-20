package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/pkg/common/bus"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/notifier/model"
)

type NotificationStore struct {
	db cmsql.Database
}

func NewNotificationStore(db cmsql.Database) *NotificationStore {
	return &NotificationStore{
		db: db,
	}
}

func (s *NotificationStore) CreateNotification(args *model.CreateNotificationArgs) (*model.Notification, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.Title == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Title")
	}
	if args.Message == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Message")
	}
	id := cm.NewID()
	noti := &model.Notification{
		ID:               id,
		Title:            args.Title,
		Message:          args.Message,
		Entity:           args.Entity,
		EntityID:         args.EntityID,
		IsRead:           false,
		AccountID:        args.AccountID,
		SendNotification: args.SendNotification,
	}
	if err := s.db.Table("notification").ShouldInsert(noti); err != nil {
		return nil, err
	}
	res, err := s.GetNotification(&model.GetNotificationArgs{
		AccountID: args.AccountID,
		ID:        id,
	})
	return res, err
}

func (s *NotificationStore) CreateNotifications(args *model.CreateNotificationsArgs) (created int, errored int, err error) {
	if len(args.AccountIDs) == 0 && !args.SendAll {
		return 0, 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Account IDs")
	}
	if args.Title == "" {
		return 0, 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Title")
	}
	if args.Message == "" {
		return 0, 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Message")
	}

	accountIDs := args.AccountIDs
	if args.SendAll {
		accountIDs = []int64{}
		deviceStore := NewDeviceStore(s.db)
		userIDs, _err := deviceStore.GetAllUsers()
		if _err != nil {
			return 0, 0, _err
		}
		query := &etopmodel.GetAllAccountUsersQuery{
			UserIDs: userIDs,
			Type:    etopmodel.TypeShop,
		}
		if err := bus.Dispatch(context.Background(), query); err != nil {
			return 0, 0, err
		}
		for _, accountUser := range query.Result {
			accountIDs = append(accountIDs, accountUser.AccountID)
		}
	}

	maxGoroutines := 8
	chCreate := make(chan error, maxGoroutines)
	guard := make(chan int, maxGoroutines)

	for i, accountID := range accountIDs {
		guard <- i
		go func(aID int64) (_err error) {
			defer func() {
				<-guard
				chCreate <- _err
			}()
			noti := &model.Notification{
				ID:               cm.NewID(),
				Title:            args.Title,
				Message:          args.Message,
				Entity:           args.Entity,
				EntityID:         args.EntityID,
				IsRead:           false,
				AccountID:        aID,
				SendNotification: args.SendNotification,
			}
			return s.db.Table("notification").ShouldInsert(noti)
		}(accountID)
	}

	for i, l := 0, len(accountIDs); i < l; i++ {
		err := <-chCreate
		if err == nil {
			created++
		} else {
			errored++
		}
	}
	ll.S.Infof("create notifications :: created %v/%v, errored %v/%v",
		created, len(accountIDs),
		errored, len(accountIDs))
	return
}

func (s *NotificationStore) GetNotification(args *model.GetNotificationArgs) (*model.Notification, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	var noti = new(model.Notification)
	err := s.db.Table("notification").Where("account_id = ? AND id = ?", args.AccountID, args.ID).ShouldGet(noti)
	return noti, err
}

func (s *NotificationStore) UpdateNotification(m *model.Notification) error {
	if m.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	err := s.db.Table("notification").Where("id = ?", m.ID).ShouldUpdate(m)
	return err
}

func (s *NotificationStore) UpdateNotifications(args *model.UpdateNotificationsArgs) error {
	if len(args.IDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing IDs")
	}
	noti := &model.Notification{
		IsRead: args.IsRead,
	}
	if args.IsRead {
		noti.SeenAt = time.Now()
	}
	err := s.db.Table("notification").In("id", args.IDs).ShouldUpdate(noti)
	return err
}

func (s *NotificationStore) GetNotifications(args *model.GetNotificationsArgs) (notis []*model.Notification, total int, err error) {
	if args.AccountID == 0 {
		return nil, 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	x := s.db.Table("notification").Where("account_id = ?", args.AccountID)
	if args.Paging != nil && len(args.Paging.Sort) == 0 {
		args.Paging.Sort = []string{"-updated_at"}
	}
	{
		x1 := x.Clone()
		x1, err = LimitSort(x1, args.Paging, Ms{"created_at": "", "updated_at": ""})
		if err != nil {
			return nil, 0, err
		}
		if err = x1.Find((*model.Notifications)(&notis)); err != nil {
			return nil, 0, err
		}
	}
	_total, err := x.Count(&model.Notification{})
	total = int(_total)
	return
}
