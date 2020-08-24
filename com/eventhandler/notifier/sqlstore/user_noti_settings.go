package sqlstore

import (
	"context"

	"github.com/lib/pq"

	"o.o/api/main/notify"
	"o.o/api/meta"
	"o.o/backend/com/eventhandler/notifier/convert"
	"o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type UserNotiSettingStore struct {
	query      cmsql.QueryFactory
	preds      []interface{}
	userNotiFt UserNotiSettingFilters

	filter         meta.Filters
	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

type UserNotiSettingStoreFactory func(context.Context) *UserNotiSettingStore

func NewUserNotiSettingStore(db *cmsql.Database) UserNotiSettingStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserNotiSettingStore {
		return &UserNotiSettingStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

func (s *UserNotiSettingStore) CreateUserNotifySettingDB(setting *model.UserNotiSetting) error {
	return s.query().ShouldInsert(setting)
}

func (s *UserNotiSettingStore) CreateUserNotifySetting(setting *notify.UserNotiSetting) error {
	settingDB := convert.Convert_api_UserNotiSetting_To_model_UserNotiSetting(setting)
	return s.CreateUserNotifySettingDB(settingDB)
}

func (s *UserNotiSettingStore) GetUserNotifySettingDB() (*model.UserNotiSetting, error) {
	var userSetting model.UserNotiSetting
	err := s.query().Where(s.preds).ShouldGet(&userSetting)
	if err != nil {
		return nil, err
	}
	return &userSetting, nil
}

func (s *UserNotiSettingStore) ByUserID(userID dot.ID) *UserNotiSettingStore {
	s.preds = append(s.preds, s.userNotiFt.ByUserID(userID))
	return s
}

func (s *UserNotiSettingStore) ByUserIDs(userIDs []dot.ID) *UserNotiSettingStore {
	s.preds = append(s.preds, sq.In("user_id", userIDs))
	return s
}

func (s *UserNotiSettingStore) GetUserNotifySetting() (*notify.UserNotiSetting, error) {
	settingDB, err := s.GetUserNotifySettingDB()
	if err != nil {
		return nil, err
	}

	return convert.Convert_model_UserNotiSetting_To_api_UserNotiSetting(settingDB), nil
}

func (s *UserNotiSettingStore) UpdateSettingDB(setting *model.UserNotiSetting) (int, error) {
	return s.query().Where(s.preds).Update(setting)
}

func (s *UserNotiSettingStore) UpdateSetting(setting *notify.UserNotiSetting) error {
	settingDB := convert.Convert_api_UserNotiSetting_To_model_UserNotiSetting(setting)
	_, err := s.UpdateSettingDB(settingDB)
	return err
}

func (s *UserNotiSettingStore) UpdateDisableTopic(disableTopics []string) error {
	return s.query().Table("user_noti_setting").Where(s.preds).ShouldUpdateMap(map[string]interface{}{
		"disable_topics": pq.StringArray(disableTopics),
	})
}
