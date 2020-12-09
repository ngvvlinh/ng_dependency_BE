package query

import (
	"context"

	"o.o/api/shopping/setting"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/setting/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

var _ setting.QueryService = &ShopSettingQuery{}

type ShopSettingQuery struct {
	store sqlstore.ShopSettingStoreFactory
}

func NewShopSettingQuery(db com.MainDB) *ShopSettingQuery {
	return &ShopSettingQuery{
		store: sqlstore.NewShopSettingStore(db),
	}
}

func ShopSettingQueryMessageBus(q *ShopSettingQuery) setting.QueryBus {
	b := bus.New()
	return setting.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (s ShopSettingQuery) GetShopSetting(
	ctx context.Context, args *setting.GetShopSettingArgs,
) (*setting.ShopSetting, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}

	shopSetting, err := s.store(ctx).GetShopSetting()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		return &setting.ShopSetting{
			ShopID: args.ShopID,
		}, nil
	case cm.NoError:
		return shopSetting, nil
	default:
		return nil, err
	}
}
