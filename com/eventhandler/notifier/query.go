package notifier

import (
	"context"

	"o.o/api/main/notify"
	"o.o/backend/com/eventhandler/notifier/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
)

var _ notify.QueryService = &QueryService{}

type QueryService struct {
	userNotiSettingStore sqlstore.UserNotiSettingStoreFactory
}

func QueryServiceNotifyBus(q *QueryService) notify.QueryBus {
	b := bus.New()
	return notify.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		sqlstore.NewUserNotiSettingStore(db),
	}
}

func (s *QueryService) GetUserNotifySetting(ctx context.Context, args *notify.GetUserNotiSettingArgs) (*notify.UserNotiSetting, error) {
	return s.userNotiSettingStore(ctx).ByUserID(args.UserID).GetUserNotifySetting()
}
