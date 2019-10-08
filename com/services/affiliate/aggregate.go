package affiliate

import (
	"context"
	"time"

	"etop.vn/api/main/identity"

	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/com/services/affiliate/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/l"

	"etop.vn/api/services/affiliate"
)

var _ affiliate.Aggregate = &Aggregate{}

var ll = l.New()

var AvailableUnit = []string{"vnd", "percent"}
var AvailablePromotionTypes = []string{"cashback", "discount"}
var AvailableDependOnValues = []string{"product", "customer"}
var AvailableDurationTypes = []string{"day", "month", "year"}

type Aggregate struct {
	commissionSetting       sqlstore.CommissionSettingStoreFactory
	supplyCommissionSetting sqlstore.SupplyCommissionSettingStoreFactory
	productPromotion        sqlstore.ProductPromotionStoreFactory
	affiliateCommission     sqlstore.AffiliateCommissionStoreFactory
	orderCreatedNotify      sqlstore.OrderCreatedNotifyStoreFactory
	affiliateReferralCode   sqlstore.AffiliateReferralCodeStoreFactory
	userReferral            sqlstore.UserReferralStoreFactory
	identityQuery           identity.QueryBus
}

func NewAggregate(db cmsql.Database, idenQuery identity.QueryBus) *Aggregate {
	return &Aggregate{
		commissionSetting:       sqlstore.NewCommissionSettingStore(db),
		supplyCommissionSetting: sqlstore.NewSupplyCommissionSettingStore(db),
		productPromotion:        sqlstore.NewProductPromotionStore(db),
		affiliateCommission:     sqlstore.NewAffiliateCommissionSettingStore(db),
		orderCreatedNotify:      sqlstore.NewOrderCreatedNotifyStore(db),
		affiliateReferralCode:   sqlstore.NewAffiliateReferralCodeStore(db),
		userReferral:            sqlstore.NewUserReferralStore(db),
		identityQuery:           idenQuery,
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

func (a *Aggregate) OnTradingOrderCreated(ctx context.Context, args *affiliate.OnTradingOrderCreatedArgs) error {
	if args.OrderID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}

	notify, err := a.orderCreatedNotify(ctx).OrderID(args.OrderID).GetOrderCreatedNotifyDB()
	if err == nil {
		return cm.Errorf(cm.AlreadyExists, nil, "Notify does existed")
	}

	notify = &model.OrderCreatedNotify{
		ID:           cm.NewID(),
		OrderID:      args.OrderID,
		ReferralCode: args.ReferralCode,
		Status:       0,
		CompletedAt:  time.Time{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := a.orderCreatedNotify(ctx).CreateOrderCreatedNotify(notify); err != nil {
		return err
	}
	return nil
}

func (a *Aggregate) TradingOrderCreating(ctx context.Context, args *affiliate.TradingOrderCreating) error {
	if args.ReferralCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ReferralCode")
	}
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}
	userQ := &identity.GetUserByIDQuery{UserID: args.UserID}
	if err := a.identityQuery.Dispatch(ctx, userQ); err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "User not found")
	}

	// TODO: Get affiliate account
	affiliateReferralCode, err := a.affiliateReferralCode(ctx).Code(args.ReferralCode).GetAffiliateReferralCodeDB()
	if err != nil {
		return err
	}

	if userQ.Result.ID == affiliateReferralCode.UserID {
		return cm.Errorf(cm.ValidationFailed, nil, "Mã giới thiệu không hợp lệ")
	}

	return nil
}

func (a *Aggregate) CreateAffiliateReferralCode(ctx context.Context, args *affiliate.CreateReferralCodeArgs) (*affiliate.AffiliateReferralCode, error) {
	if args.AffiliateAccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "AffiliateAccountID missing")
	}
	if args.Code == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Code is missing")
	}
	affiliateQ := &identity.GetAffiliateByIDQuery{
		ID: args.AffiliateAccountID,
	}
	if err := a.identityQuery.Dispatch(ctx, affiliateQ); err != nil {
		return nil, err
	}
	_, err := a.affiliateReferralCode(ctx).Code(args.Code).GetAffiliateReferralCodeDB()
	if err == nil {
		return nil, cm.Errorf(cm.ValidationFailed, nil, "Referral Code does exist")
	}
	affiliateReferralCode := &model.AffiliateReferralCode{
		ID:          cm.NewID(),
		Code:        args.Code,
		UserID:      affiliateQ.Result.OwnerID,
		AffiliateID: args.AffiliateAccountID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := a.affiliateReferralCode(ctx).CreateAffiliateReferralCode(affiliateReferralCode); err != nil {
		return nil, err
	}
	return a.affiliateReferralCode(ctx).ID(affiliateReferralCode.ID).GetAffiliateReferralCode()
}

func (a *Aggregate) CreateOrUpdateUserReferral(ctx context.Context, args *affiliate.CreateOrUpdateReferralArgs) (*affiliate.UserReferral, error) {
	if args.UserID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UserID is missing")
	}
	userReferral, err := a.userReferral(ctx).UserID(args.UserID).GetUserReferralDB()
	if err != nil {
		userReferral = &model.UserReferral{
			UserID: args.UserID,
		}
		if err := a.userReferral(ctx).CreateUserReferral(userReferral); err != nil {
			return nil, err
		}
	}
	if args.ReferralCode != "" {
		if userReferral.ReferralID != 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tài khoản đã được giới thiệu bởi thành viên khác")
		}
		affiliateReferralCode, err := a.affiliateReferralCode(ctx).Code(args.ReferralCode).GetAffiliateReferralCodeDB()
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "ReferralCode not found")
		}
		userReferral.ReferralCode = args.ReferralCode
		userReferral.ReferralID = affiliateReferralCode.AffiliateID
		userReferral.ReferralAt = time.Now()
	}
	if args.SaleReferralCode != "" {
		if userReferral.SaleReferralID != 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tài khoản đã được giới thiệu bởi thành viên khác")
		}
		affiliateReferralCode, err := a.affiliateReferralCode(ctx).Code(args.SaleReferralCode).GetAffiliateReferralCodeDB()
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "ReferralCode not found")
		}
		userReferral.SaleReferralCode = args.SaleReferralCode
		userReferral.SaleReferralID = affiliateReferralCode.AffiliateID
		userReferral.SaleReferralAt = time.Now()
	}
	err = a.userReferral(ctx).UpdateUserReferral(userReferral)
	if err != nil {
		return nil, err
	}

	return a.userReferral(ctx).UserID(args.UserID).GetUserReferral()
}

func (a *Aggregate) CreateOrUpdateSupplyCommissionSetting(ctx context.Context, args *affiliate.CreateOrUpdateSupplyCommissionSettingArgs) (*affiliate.SupplyCommissionSetting, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ShopID is invalid")
	}
	if args.ProductID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ProductID is invalid")
	}
	if !cm.StringsContain(AvailableDependOnValues, args.DependOn) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Depend On value is not valid")
	}
	if !cm.StringsContain(AvailableDurationTypes, args.Level1LimitDurationType) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Level1LimitDurationType is not valid")
	}
	if !cm.StringsContain(AvailableDurationTypes, args.LifetimeDurationType) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "LifetimeDurationType is not valid")
	}

	var level1LimitDuration int64 = 0
	var lifetimeDuration int64 = 0

	switch args.Level1LimitDurationType {
	case "day":
		level1LimitDuration = int64(args.Level1LimitDuration) * 86400
	case "month":
		level1LimitDuration = int64(args.Level1LimitDuration) * 2592000
	case "year":
		level1LimitDuration = int64(args.Level1LimitDuration) * 31104000
	}

	switch args.LifetimeDurationType {
	case "day":
		lifetimeDuration = int64(args.LifetimeDuration) * 86400
	case "month":
		lifetimeDuration = int64(args.LifetimeDuration) * 2592000
	case "year":
		lifetimeDuration = int64(args.LifetimeDuration) * 31104000
	}

	var createOrUpdateErr error
	supplyCommissionSetting := &model.SupplyCommissionSetting{
		ShopID:                   args.ShopID,
		ProductID:                args.ProductID,
		Level1DirectCommission:   args.Level1DirectCommission,
		Level1IndirectCommission: args.Level1IndirectCommission,
		Level2DirectCommission:   args.Level2DirectCommission,
		Level2IndirectCommission: args.Level2IndirectCommission,
		DependOn:                 args.DependOn,
		Level1LimitCount:         args.Level1LimitCount,
		Level1LimitDuration:      level1LimitDuration,
		MLevel1LimitDuration: &model.DurationJSON{
			Duration: int32(args.Level1LimitDuration),
			Type:     args.Level1LimitDurationType,
		},
		LifetimeDuration: lifetimeDuration,
		MLifetimeDuration: &model.DurationJSON{
			Duration: int32(args.LifetimeDuration),
			Type:     args.LifetimeDurationType,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := a.supplyCommissionSetting(ctx).ProductID(args.ProductID).ShopID(args.ShopID).GetSupplyCommissionSettingDB()

	if cm.ErrorCode(err) == cm.NotFound {
		createOrUpdateErr = a.supplyCommissionSetting(ctx).CreateSupplyCommissionSetting(supplyCommissionSetting)
	} else {
		createOrUpdateErr = a.supplyCommissionSetting(ctx).UpdateSupplyCommissionSetting(supplyCommissionSetting)
	}

	if createOrUpdateErr != nil {
		return nil, createOrUpdateErr
	}

	return a.supplyCommissionSetting(ctx).ShopID(args.ShopID).ProductID(args.ProductID).GetSupplyCommissionSetting()
}
