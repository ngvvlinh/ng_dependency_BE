package usersetting

import (
	"context"

	"o.o/api/etelecom/usersetting"
	"o.o/backend/com/etelecom/usersetting/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type UserSettingAggregate struct {
	userSettingStore sqlstore.UserSettingStoreFactory
}

var _ usersetting.Aggregate = &UserSettingAggregate{}

func NewUserSettingAggregate(dbEtelecom com.EtelecomDB) *UserSettingAggregate {
	return &UserSettingAggregate{
		userSettingStore: sqlstore.NewUserSettingStore(dbEtelecom),
	}
}

func AggerateMessageBus(a *UserSettingAggregate) usersetting.CommandBus {
	b := bus.New()
	return usersetting.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *UserSettingAggregate) UpdateUserSetting(ctx context.Context, args *usersetting.UpdateUserSettingArgs) (*usersetting.UserSetting, error) {
	_, err := a.userSettingStore(ctx).ID(args.UserID).GetUserSetting()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		// create new one
		userSetting := &usersetting.UserSetting{
			ID:                  args.UserID,
			ExtensionChargeType: args.ExtensionChargeType,
		}
		return a.userSettingStore(ctx).CreateUserSetting(userSetting)
	case cm.NoError:
		// update existed setting
		update := &usersetting.UserSetting{
			ExtensionChargeType: args.ExtensionChargeType,
		}
		err = a.userSettingStore(ctx).ID(args.UserID).UpdateUserSetting(update)
		if err != nil {
			return nil, err
		}
		return a.userSettingStore(ctx).ID(args.UserID).GetUserSetting()
	default:
		return nil, err
	}
}
