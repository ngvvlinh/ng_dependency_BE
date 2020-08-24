package notifier

import (
	"context"

	"o.o/api/main/notify"
	"o.o/api/main/shipnow/carrier"
	"o.o/backend/com/eventhandler/notifier/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ notify.Aggregate = &Aggregate{}

func NewNotiAggregate(db com.MainDB, carrierManager carrier.Manager) *Aggregate {
	return &Aggregate{
		db:                   db,
		userNotiSettingStore: sqlstore.NewUserNotiSettingStore(db),
	}
}

func NewNotiAggregateMessageBus(s *Aggregate) notify.CommandBus {
	b := bus.New()
	return notify.NewAggregateHandler(s).RegisterHandlers(b)
}

type Aggregate struct {
	db                   *cmsql.Database
	userNotiSettingStore sqlstore.UserNotiSettingStoreFactory
}

func (s *Aggregate) CreateUserNotifySetting(
	ctx context.Context,
	args *notify.CreateUserNotifySettingArgs,
) (*notify.UserNotiSetting, error) {
	setting := &notify.UserNotiSetting{
		UserID:        args.UserID,
		DisableTopics: args.DisableTopics,
	}
	err := s.userNotiSettingStore(ctx).CreateUserNotifySetting(setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (s *Aggregate) GetOrCreateUserNotifySetting(ctx context.Context, args *notify.GetOrCreateUserNotifySettingArgs) (*notify.UserNotiSetting, error) {
	setting, err := s.userNotiSettingStore(ctx).ByUserID(args.UserID).GetUserNotifySetting()
	if err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			_setting := &notify.UserNotiSetting{
				UserID:        args.UserID,
				DisableTopics: args.DisableTopics,
			}
			err = s.userNotiSettingStore(ctx).CreateUserNotifySetting(_setting)
			if err != nil {
				return nil, err
			}
			return _setting, nil
		}
		return nil, err
	}
	return setting, nil
}

func (s *Aggregate) DisableTopic(ctx context.Context, args *notify.DisableTopicArgs) (*notify.UserNotiSetting, error) {
	getOrCreateArgs := &notify.GetOrCreateUserNotifySettingArgs{
		UserID:        args.UserID,
		DisableTopics: []string{},
	}
	setting, err := s.GetOrCreateUserNotifySetting(ctx, getOrCreateArgs)
	if err != nil {
		return nil, err
	}

	_topics := append(setting.DisableTopics, args.Topic)
	err = s.userNotiSettingStore(ctx).ByUserID(args.UserID).UpdateDisableTopic(_topics)
	if err != nil {
		return nil, err
	}
	setting.DisableTopics = _topics
	return setting, nil
}

func (s *Aggregate) EnableTopic(ctx context.Context, args *notify.EnableTopicArgs) (*notify.UserNotiSetting, error) {
	getOrCreateArgs := &notify.GetOrCreateUserNotifySettingArgs{
		UserID:        args.UserID,
		DisableTopics: []string{},
	}
	setting, err := s.GetOrCreateUserNotifySetting(ctx, getOrCreateArgs)
	if err != nil {
		return nil, err
	}
	newDisableTopics, contains := removeTopicIfContains(setting.DisableTopics, args.Topic)
	if !contains {
		return setting, nil
	}
	if err = s.userNotiSettingStore(ctx).ByUserID(args.UserID).UpdateDisableTopic(newDisableTopics); err != nil {
		return nil, err
	}
	setting.DisableTopics = newDisableTopics
	return setting, nil
}

func removeTopicIfContains(topics []string, it string) ([]string, bool) {
	var _topics []string
	isContains := false
	for _, _it := range topics {
		if _it == it {
			isContains = true
			continue
		}
		_topics = append(_topics, _it)
	}
	return _topics, isContains
}
