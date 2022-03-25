package query

import (
	"context"

	"o.o/api/shopping/setting"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/setting/sqlstore"
	"o.o/backend/com/shopping/setting/util"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

var _ setting.QueryService = &ShopSettingQuery{}

type ShopSettingQuery struct {
	store sqlstore.ShopSettingStoreFactory
	util  *util.ShopSettingUtil
}

func NewShopSettingQuery(
	db com.MainDB,
	util *util.ShopSettingUtil,
) *ShopSettingQuery {
	return &ShopSettingQuery{
		store: sqlstore.NewShopSettingStore(db),
		util:  util,
	}
}

func ShopSettingQueryMessageBus(q *ShopSettingQuery) setting.QueryBus {
	b := bus.New()
	return setting.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (s ShopSettingQuery) GetShopSetting(
	ctx context.Context, args *setting.GetShopSettingArgs,
) (shopSetting *setting.ShopSetting, err error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}

	shopSetting, _err := s.util.GetShopSetting(args.ShopID)
	if _err != nil {
		return nil, _err
	}
	if shopSetting != nil {
		return shopSetting, nil
	}

	shopSetting, err = s.store(ctx).ShopID(args.ShopID).GetShopSetting()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		shopSetting = &setting.ShopSetting{
			ShopID: args.ShopID,
		}
	case cm.NoError:
		// no-op
	default:
		return nil, err
	}

	if _err := s.util.SetShopSetting(args.ShopID, *shopSetting); _err != nil {
		return nil, _err
	}

	return shopSetting, nil
}

func (s ShopSettingQuery) GetShopSettingDirectShipment(
	ctx context.Context, args *setting.GetShopSettingArgs,
) (shopSetting *setting.ShopSetting, err error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}
	shopSetting, err = s.store(ctx).ShopID(args.ShopID).GetShopSetting()
	if err != nil {
		return nil, err
	}
	return shopSetting, nil
}
