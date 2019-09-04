package affiliate

import (
	"context"

	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/com/services/affiliate/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
	"etop.vn/common/l"

	"etop.vn/api/services/affiliate"
)

var _ affiliate.Aggregate = &Aggregate{}

var ll = l.New()

type Aggregate struct {
	commissionSetting sqlstore.CommissionSettingStoreFactory
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{
		commissionSetting: sqlstore.NewShopCommissionSettingStore(db),
	}
}

func (a *Aggregate) MessageBus() affiliate.CommandBus {
	b := bus.New()
	return affiliate.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateOrUpdateCommissionSetting(ctx context.Context, args *affiliate.CreateCommissionSettingArgs) (*affiliate.CommissionSetting, error) {
	availableUnit := []string{"vnd", "percent"}
	if args.Unit == "" {
		args.Unit = "percent"
	} else if !cm.StringsContain(availableUnit, args.Unit) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Unit not valid")
	}

	if args.Amount < 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Amount should be greater than 0")
	}
	if args.ProductID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ProductID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}

	shopCommissionSetting, err := a.commissionSetting(ctx).AccountID(args.AccountID).ProductID(args.ProductID).GetShopCommissionSettingDB()
	if err != nil {
		shopCommissionSetting = &model.CommissionSetting{
			AccountID: args.AccountID,
			ProductID: args.ProductID,
			Amount:    args.Amount,
			Unit:      args.Unit,
		}
		err = a.commissionSetting(ctx).CreateShopCommissionSetting(shopCommissionSetting)
	} else {
		shopCommissionSetting.Amount = args.Amount
		shopCommissionSetting.Unit = args.Unit
		err = a.commissionSetting(ctx).UpdateShopCommissionSetting(shopCommissionSetting)
	}

	if err != nil {
		return nil, err
	}

	return a.commissionSetting(ctx).GetShopCommissionSettingByProductID(args.AccountID, args.ProductID)
}
