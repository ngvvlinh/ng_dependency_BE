package sqlstore

import (
	"context"
	"time"

	"o.o/api/top/types/etc/account_type"
	"o.o/backend/com/eventhandler/notifier/model"
	com "o.o/backend/com/main"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	cmstore "o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type NotificationStore struct {
	db *cmsql.Database

	AccountUserStore sqlstore.AccountUserStoreInterface
}

func NewNotificationStore(
	db com.NotifierDB,
	accountUserStore sqlstore.AccountUserStoreInterface,
) *NotificationStore {
	model.SQLVerifySchema(db)
	return &NotificationStore{
		db: db,

		AccountUserStore: accountUserStore,
	}
}

func (s *NotificationStore) CreateNotification(args *model.CreateNotificationArgs) (*model.Notification, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}

	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Message")
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
		UserID:           args.UserID,
		Title:            args.Title,
		Message:          args.Message,
		Entity:           args.Entity,
		EntityID:         args.EntityID,
		IsRead:           false,
		AccountID:        args.AccountID,
		SendNotification: args.SendNotification,
		MetaData:         args.MetaData,
		TopicType:        args.TopicType,
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
		accountIDs = []dot.ID{}
		deviceStore := NewDeviceStore(s.db)
		userIDs, _err := deviceStore.GetAllUsers()
		if _err != nil {
			return 0, 0, _err
		}
		query := &identitymodelx.GetAllAccountUsersQuery{
			UserIDs: userIDs,
			Type:    account_type.Shop.Wrap(),
		}
		if err := s.AccountUserStore.GetAllAccountUsers(context.Background(), query); err != nil {
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
		go func(aID dot.ID) (_err error) {
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
				MetaData:         args.MetaData,
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

func (s *NotificationStore) GetNotifications(args *model.GetNotificationsArgs) (notis []*model.Notification, err error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	x := s.db.Table("notification").Where("account_id = ?", args.AccountID)
	if args.Paging != nil && len(args.Paging.Sort) == 0 {
		args.Paging.Sort = []string{"-updated_at"}
	}
	{
		x1 := x.Clone()
		x1, err = cmstore.LimitSort(x1, cmstore.ConvertPaging(args.Paging), map[string]string{"created_at": "", "updated_at": ""})
		if err != nil {
			return nil, err
		}
		if err = x1.Find((*model.Notifications)(&notis)); err != nil {
			return nil, err
		}
	}
	return
}
