package sqlstore

import (
	"context"

	"o.o/api/etelecom/usersetting"
	"o.o/backend/com/etelecom/usersetting/convert"
	"o.o/backend/com/etelecom/usersetting/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

type UserSettingStore struct {
	ft    UserSettingFilters
	query func() cmsql.QueryInterface
	preds []interface{}
}

type UserSettingStoreFactory func(ctx context.Context) *UserSettingStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewUserSettingStore(db *cmsql.Database) UserSettingStoreFactory {
	return func(ctx context.Context) *UserSettingStore {
		return &UserSettingStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *UserSettingStore) ID(id dot.ID) *UserSettingStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *UserSettingStore) GetUserSettingDB() (*model.UserSetting, error) {
	var res model.UserSetting
	err := s.query().Where(s.preds).ShouldGet(&res)
	return &res, err
}

func (s *UserSettingStore) GetUserSetting() (*usersetting.UserSetting, error) {
	settingDB, err := s.GetUserSettingDB()
	if err != nil {
		return nil, err
	}
	var res usersetting.UserSetting
	if err = scheme.Convert(settingDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *UserSettingStore) CreateUserSetting(setting *usersetting.UserSetting) (*usersetting.UserSetting, error) {
	var settingDB model.UserSetting
	if err := scheme.Convert(setting, &settingDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&settingDB); err != nil {
		return nil, err
	}
	return s.ID(setting.ID).GetUserSetting()
}

func (s *UserSettingStore) UpdateUserSetting(setting *usersetting.UserSetting) error {
	var settingDB model.UserSetting
	if err := scheme.Convert(setting, &settingDB); err != nil {
		return err
	}
	return s.query().Where(s.preds).ShouldUpdate(&settingDB)
}
