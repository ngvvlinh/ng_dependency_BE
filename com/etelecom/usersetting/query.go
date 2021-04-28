package usersetting

import (
	"context"

	"o.o/api/etelecom/usersetting"
	"o.o/backend/com/etelecom/usersetting/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

type QueryService struct {
	userSettingStore sqlstore.UserSettingStoreFactory
}

var _ usersetting.QueryService = &QueryService{}

func NewQueryService(db com.EtelecomDB) *QueryService {
	return &QueryService{
		userSettingStore: sqlstore.NewUserSettingStore(db),
	}
}

func QueryServiceMessageBus(s *QueryService) usersetting.QueryBus {
	b := bus.New()
	return usersetting.NewQueryServiceHandler(s).RegisterHandlers(b)
}

func (q *QueryService) GetUserSetting(ctx context.Context, userID dot.ID) (*usersetting.UserSetting, error) {
	return q.userSettingStore(ctx).ID(userID).GetUserSetting()
}

func (q *QueryService) ListUserSettings(ctx context.Context, args *usersetting.ListUserSettingsArgs) (*usersetting.ListUserSettingsResponse, error) {
	query := q.userSettingStore(ctx).WithPaging(args.Paging)
	if len(args.UserIDs) > 0 {
		query = query.IDs(args.UserIDs)
	}
	userSettings, err := query.ListUserSetting()
	if err != nil {
		return nil, err
	}
	return &usersetting.ListUserSettingsResponse{
		UserSettings: userSettings,
		Paging:       query.GetPaging(),
	}, nil
}
