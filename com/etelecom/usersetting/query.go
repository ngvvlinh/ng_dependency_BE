package usersetting

import (
	"context"

	"o.o/api/etelecom/usersetting"
	"o.o/backend/com/etelecom/usersetting/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
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
	setting, err := q.userSettingStore(ctx).ID(userID).GetUserSetting()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return &usersetting.UserSetting{ID: userID}, nil
	case cm.NoError:
		return setting, nil
	default:
		return nil, err
	}
}
