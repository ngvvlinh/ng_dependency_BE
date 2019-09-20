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

var AvailableUnit = []string{"vnd", "percent"}
var AvailablePromotionTypes = []string{"cashback", "discount"}

type Aggregate struct {
	commissionSetting sqlstore.CommissionSettingStoreFactory
	productPromotion  sqlstore.ProductPromotionStoreFactory
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{
		commissionSetting: sqlstore.NewCommissionSettingStore(db),
		productPromotion:  sqlstore.NewProductPromotionStore(db),
	}
}

func (a *Aggregate) MessageBus() affiliate.CommandBus {
	b := bus.New()
	return affiliate.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateOrUpdateCommissionSetting(ctx context.Context, args *affiliate.CreateCommissionSettingArgs) (*affiliate.CommissionSetting, error) {
	if args.Unit == "" {
		args.Unit = "percent"
	} else if !cm.StringsContain(AvailableUnit, args.Unit) {
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

	shopCommissionSetting, err := a.commissionSetting(ctx).AccountID(args.AccountID).ProductID(args.ProductID).GetCommissionSettingDB()
	if err != nil {
		shopCommissionSetting = &model.CommissionSetting{
			AccountID: args.AccountID,
			ProductID: args.ProductID,
			Amount:    args.Amount,
			Unit:      args.Unit,
			Type:      args.Type,
		}
		err = a.commissionSetting(ctx).CreateCommissionSetting(shopCommissionSetting)
	} else {
		shopCommissionSetting.Amount = args.Amount
		shopCommissionSetting.Unit = args.Unit
		shopCommissionSetting.Type = args.Type
		err = a.commissionSetting(ctx).UpdateCommissionSetting(shopCommissionSetting)
	}

	if err != nil {
		return nil, err
	}

	return a.commissionSetting(ctx).AccountID(args.AccountID).ProductID(args.ProductID).GetCommissionSetting()
}

func (a *Aggregate) CreateProductPromotion(ctx context.Context, args *affiliate.CreateProductPromotionArgs) (*affiliate.ProductPromotion, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ShopID is missing")
	}
	if args.ProductID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ProductID is missing")
	}
	if !cm.StringsContain(AvailablePromotionTypes, args.Type) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Type not valid")
	}
	if args.Amount <= 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Amount must be greater than 0")
	}
	if args.Unit == "" {
		args.Unit = "percent"
	} else if !cm.StringsContain(AvailableUnit, args.Unit) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Unit not valid")
	}

	productPromotion := &model.ProductPromotion{
		ID:          cm.NewID(),
		ProductID:   args.ProductID,
		ShopID:      args.ShopID,
		Amount:      args.Amount,
		Unit:        args.Unit,
		Code:        args.Code,
		Description: args.Description,
		Note:        args.Note,
		Type:        args.Type,
	}
	if err := a.productPromotion(ctx).CreateProductPromotion(productPromotion); err != nil {
		return nil, err
	}
	return a.productPromotion(ctx).ID(productPromotion.ID).GetProductPromotion()
}

func (a *Aggregate) UpdateProductPromotion(ctx context.Context, args *affiliate.UpdateProductPromotionArgs) (*affiliate.ProductPromotion, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ID is missing")
	}

	if args.Amount <= 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Amount must be greater than 0")
	}
	if !cm.StringsContain(AvailablePromotionTypes, args.Type) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Type not valid")
	}

	if !cm.StringsContain(AvailableUnit, args.Unit) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Unit not valid")
	}

	productPromotion, err := a.productPromotion(ctx).ID(args.ID).GetProductPromotionDB()
	if err != nil {
		return nil, err
	}

	productPromotion.Amount = args.Amount
	productPromotion.Unit = args.Unit
	productPromotion.Type = args.Type

	if err := a.productPromotion(ctx).UpdateProductPromotion(productPromotion); err != nil {
		return nil, err
	}
	return a.productPromotion(ctx).ID(productPromotion.ID).GetProductPromotion()
}
