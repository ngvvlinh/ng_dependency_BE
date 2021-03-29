package sqlstore

import (
	"context"

	"o.o/api/etelecom/usersetting"
	"o.o/api/meta"
	"o.o/backend/com/etelecom/usersetting/convert"
	"o.o/backend/com/etelecom/usersetting/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type UserSettingStore struct {
	ft    UserSettingFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	sqlstore.Paging
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

func (s *UserSettingStore) WithPaging(paging meta.Paging) *UserSettingStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *UserSettingStore) IDs(ids []dot.ID) *UserSettingStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
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

func (s *UserSettingStore) ListUserSettingDB() (res []*model.UserSetting, err error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortUserSetting)
	if err != nil {
		return nil, err
	}
	if err = query.Find((*model.UserSettings)(&res)); err != nil {
		return nil, err
	}
	s.Paging.Apply(res)
	return
}

func (s *UserSettingStore) ListUserSetting() ([]*usersetting.UserSetting, error) {
	userSettingsDB, err := s.ListUserSettingDB()
	if err != nil {
		return nil, err
	}
	var res []*usersetting.UserSetting
	if err = scheme.Convert(userSettingsDB, &res); err != nil {
		return nil, err
	}
	return res, err
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
