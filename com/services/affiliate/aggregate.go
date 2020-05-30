package affiliate

import (
	"context"
	"math"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/services/affiliate"
	com "o.o/backend/com/main"
	"o.o/backend/com/services/affiliate/model"
	"o.o/backend/com/services/affiliate/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ affiliate.Aggregate = &Aggregate{}

var ll = l.New()

var AvailableUnit = []string{"vnd", "percent"}
var AvailablePromotionTypes = []string{"cashback", "discount"}
var AvailableDependOnValues = []string{"product", "customer"}
var AvailableDurationTypes = []string{"day", "month", "year"}

const DependOnProduct = "product"
const DependOnCustomer = "customer"

type Aggregate struct {
	commissionSetting       sqlstore.CommissionSettingStoreFactory
	supplyCommissionSetting sqlstore.SupplyCommissionSettingStoreFactory
	productPromotion        sqlstore.ProductPromotionStoreFactory
	affiliateCommission     sqlstore.SellerCommissionStoreFactory
	orderCreatedNotify      sqlstore.OrderCreatedNotifyStoreFactory
	affiliateReferralCode   sqlstore.AffiliateReferralCodeStoreFactory
	userReferral            sqlstore.UserReferralStoreFactory
	orderPromotion          sqlstore.OrderPromotionStoreFactory
	shopCashback            sqlstore.ShopCashbackStoreFactory
	orderCommissionSetting  sqlstore.OrderCommissionSettingStoreFactory
	shopOrderProductHistory sqlstore.ShopOrderProductHistoryStoreFactory
	customerPolicyGroup     sqlstore.CustomerPolicyGroupStoreFactory
	identityQuery           identity.QueryBus
	catalogQuery            catalog.QueryBus
	orderingQuery           ordering.QueryBus

	db *cmsql.Database
}

func NewAggregate(db com.MainDB, idenQuery identity.QueryBus, catalogQuery catalog.QueryBus, orderQuery ordering.QueryBus) *Aggregate {
	return &Aggregate{
		commissionSetting:       sqlstore.NewCommissionSettingStore(db),
		supplyCommissionSetting: sqlstore.NewSupplyCommissionSettingStore(db),
		productPromotion:        sqlstore.NewProductPromotionStore(db),
		affiliateCommission:     sqlstore.NewSellerCommissionSettingStore(db),
		orderCreatedNotify:      sqlstore.NewOrderCreatedNotifyStore(db),
		affiliateReferralCode:   sqlstore.NewAffiliateReferralCodeStore(db),
		userReferral:            sqlstore.NewUserReferralStore(db),
		orderPromotion:          sqlstore.NewOrderPromotionStore(db),
		shopCashback:            sqlstore.NewShopCashbackStore(db),
		orderCommissionSetting:  sqlstore.NewOrderCommissionSettingStoreFactory(db),
		shopOrderProductHistory: sqlstore.NewShopOrderProductHistoryStore(db),
		customerPolicyGroup:     sqlstore.NewCustomerPolicyGroupStore(db),
		identityQuery:           idenQuery,
		catalogQuery:            catalogQuery,
		orderingQuery:           orderQuery,
		db:                      db,
	}
}

func AggregateMessageBus(a *Aggregate) affiliate.CommandBus {
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
	if args.Amount < 0 {
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

	if args.Amount < 0 {
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
		return cm.Errorf(cm.ValidationFailed, nil, "Không thể sử dụng mã giới thiệu của chính bạn")
	}

	return nil
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

	go func() {
		if err := a.FetchOrderInfoToNotify(bus.Ctx(), notify.ID); err != nil {
			ll.Error(err.Error())
			return
		}

		if err := a.CreateOrderPromotions(bus.Ctx(), notify.ID, args.OrderID); err != nil {
			ll.Error(err.Error())
			notify.PromotionSnapshotErr = err.Error()
			notify.PromotionSnapshotStatus = -1
		}
		if err := a.CreateOrderCommissionSettings(bus.Ctx(), notify.ID); err != nil {
			ll.Error(err.Error())
			notify.CommissionSnapshotErr = err.Error()
			notify.CommissionSnapshotStatus = -1
		}
		notifyUpdate, err := a.orderCreatedNotify(bus.Ctx()).ID(notify.ID).GetOrderCreatedNotifyDB()
		if err != nil {
			ll.Error(err.Error())
			return
		}
		notifyUpdate.PromotionSnapshotErr = notify.PromotionSnapshotErr
		notifyUpdate.CommissionSnapshotErr = notify.CommissionSnapshotErr
		_ = a.orderCreatedNotify(bus.Ctx()).UpdateOrderCreatedNotify(notifyUpdate)

		if err = a.ProcessOrderNotify(bus.Ctx(), notify.ID); err != nil {
			ll.Error("ProcessOrderNotify Error", l.Object("err", err))
		}

	}()

	return nil
}

func (a *Aggregate) FetchOrderInfoToNotify(ctx context.Context, orderCreatedNotifyID dot.ID) error {
	notify, err := a.orderCreatedNotify(ctx).ID(orderCreatedNotifyID).GetOrderCreatedNotifyDB()
	if err != nil {
		return err
	}

	orderQ := &ordering.GetOrderByIDQuery{
		ID: notify.OrderID,
	}
	if err := a.orderingQuery.Dispatch(ctx, orderQ); err != nil {
		return err
	}

	notify.ShopID = orderQ.Result.TradingShopID
	notify.SupplyID = orderQ.Result.ShopID

	shopQ := &identity.GetShopByIDQuery{
		ID: orderQ.Result.ShopID,
	}
	if err := a.identityQuery.Dispatch(ctx, shopQ); err == nil {
		notify.ShopUserID = shopQ.Result.OwnerID
	}

	sellerReferralCode, err := a.affiliateReferralCode(ctx).Code(notify.ReferralCode).GetAffiliateReferralCodeDB()

	if err == nil {
		notify.SellerID = sellerReferralCode.AffiliateID
	}

	notify.Status = 2
	return a.orderCreatedNotify(ctx).UpdateOrderCreatedNotify(notify)
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
	if args.DependOn == DependOnProduct {
		args.Group = ""
	}

	if args.Level1LimitDurationType == "" {
		args.Level1LimitDurationType = "day"
		args.Level1LimitDuration = 0
	}

	if args.LifetimeDurationType == "" {
		args.LifetimeDurationType = "day"
		args.LifetimeDuration = 0
	}

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
	if args.Level1LimitCount < 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giá trị giới hạn khách hàng mới tối thiểu là 1")
	}

	if args.DependOn == "customer" && args.Group == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập nhóm khách hàng khi cài đặt chiết khấu theo khách hàng")
	}

	if args.Level1DirectCommission < 0 ||
		args.Level1IndirectCommission < 0 ||
		args.Level2DirectCommission < 0 ||
		args.Level2IndirectCommission < 0 ||
		args.LifetimeDuration < 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giá trị số phải lớn hơn 0")
	}

	customerPolicyGroupID := dot.ID(0)
	customerPolicyGroupName := args.Group
	if customerPolicyGroupName != "" {
		customerPolicyGroup, err := a.customerPolicyGroup(ctx).Name(customerPolicyGroupName).GetCustomerPolicyGroupDB()
		if cm.ErrorCode(err) == cm.NotFound {
			customerPolicyGroup = &model.CustomerPolicyGroup{
				ID:        cm.NewID(),
				SupplyID:  args.ShopID,
				Name:      args.Group,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err = a.customerPolicyGroup(ctx).CreateCustomerPolicyGroup(customerPolicyGroup)
			if err != nil {
				return nil, err
			}
		}
		customerPolicyGroupID = customerPolicyGroup.ID
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

	ll.Info("CHECK ARGS", l.Object("ARGS", args))

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
			Duration: args.Level1LimitDuration,
			Type:     args.Level1LimitDurationType,
		},
		LifetimeDuration: lifetimeDuration,
		MLifetimeDuration: &model.DurationJSON{
			Duration: args.LifetimeDuration,
			Type:     args.LifetimeDurationType,
		},
		CustomerPolicyGroupID: customerPolicyGroupID,
		Group:                 customerPolicyGroupName,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
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

func (a *Aggregate) CreateOrderPromotions(ctx context.Context, orderNotifyID dot.ID, orderID dot.ID) error {
	orderNotify, err := a.orderCreatedNotify(ctx).ID(orderNotifyID).GetOrderCreatedNotifyDB()
	if err != nil {
		return err
	}
	getOrderQ := &ordering.GetOrderByIDQuery{
		ID: orderID,
	}
	if err := a.orderingQuery.Dispatch(ctx, getOrderQ); err != nil {
		return err
	}

	var promotions []*model.OrderPromotion
	for _, line := range getOrderQ.Result.Lines {
		var basePrice = float64(line.TotalPrice)
		tradingPromotion, err := a.productPromotion(ctx).ProductID(line.ProductID).GetProductPromotionDB()
		if err == nil {
			promotions = append(promotions, &model.OrderPromotion{
				ID:                   cm.NewID(),
				ProductID:            line.ProductID,
				OrderID:              getOrderQ.Result.ID,
				BaseValue:            int(basePrice),
				Amount:               tradingPromotion.Amount,
				Unit:                 tradingPromotion.Unit,
				Type:                 tradingPromotion.Type,
				OrderCreatedNotifyID: orderNotify.ID,
				Src:                  "etop",
				Description:          "Khuyến mãi từ eTop Trading",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			})
			switch tradingPromotion.Unit {
			case "vnd":
				basePrice = basePrice - float64(tradingPromotion.Amount)
			case "percent":
				basePrice = math.Round(basePrice - basePrice*(float64(tradingPromotion.Amount)/100/100))
			}
		}

		affReferralCode, err := a.affiliateReferralCode(ctx).Code(orderNotify.ReferralCode).GetAffiliateReferralCodeDB()
		if err != nil {
			continue
		}
		sellerCashback, err := a.commissionSetting(ctx).AccountID(affReferralCode.AffiliateID).ProductID(line.ProductID).GetCommissionSettingDB()
		if err != nil {
			continue
		}
		promotions = append(promotions, &model.OrderPromotion{
			ID:                   cm.NewID(),
			ProductID:            line.ProductID,
			OrderID:              getOrderQ.Result.ID,
			ProductQuantity:      line.Quantity,
			BaseValue:            int(basePrice),
			Amount:               sellerCashback.Amount,
			Unit:                 sellerCashback.Unit,
			Type:                 "cashback",
			OrderCreatedNotifyID: orderNotify.ID,
			Src:                  "seller",
			Description:          "Khuyến mãi từ mã giảm giá",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		})
	}
	if len(promotions) != 0 {
		return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
			for _, promotion := range promotions {
				err := a.orderPromotion(ctx).CreateOrderPromotion(promotion)
				if err != nil {
					return err
				}
			}

			orderNotify.PromotionSnapshotStatus = 1
			return a.orderCreatedNotify(ctx).UpdateOrderCreatedNotify(orderNotify)
		})
	}
	return nil
}

func (a *Aggregate) CreateOrderCommissionSettings(ctx context.Context, orderCreatedNotifyID dot.ID) error {
	orderNotify, err := a.orderCreatedNotify(ctx).ID(orderCreatedNotifyID).GetOrderCreatedNotifyDB()
	if err != nil {
		return err
	}
	getOrderQ := &ordering.GetOrderByIDQuery{
		ID: orderNotify.OrderID,
	}
	if err := a.orderingQuery.Dispatch(ctx, getOrderQ); err != nil {
		return err
	}

	for _, line := range getOrderQ.Result.Lines {
		supplyCommissionSetting, err := a.supplyCommissionSetting(ctx).ShopID(orderNotify.SupplyID).ProductID(line.ProductID).GetSupplyCommissionSettingDB()
		if err != nil {
			return err
		}

		if err := a.orderCommissionSetting(ctx).CreateOrderCommissionSetting(&model.OrderCommissionSetting{
			OrderID:                  orderNotify.OrderID,
			SupplyID:                 orderNotify.SupplyID,
			ProductID:                line.ProductID,
			ProductQuantity:          line.Quantity,
			Level1DirectCommission:   supplyCommissionSetting.Level1DirectCommission,
			Level1IndirectCommission: supplyCommissionSetting.Level1IndirectCommission,
			Level2DirectCommission:   supplyCommissionSetting.Level2DirectCommission,
			Level2IndirectCommission: supplyCommissionSetting.Level2IndirectCommission,
			DependOn:                 supplyCommissionSetting.DependOn,
			Level1LimitCount:         supplyCommissionSetting.Level1LimitCount,
			Level1LimitDuration:      supplyCommissionSetting.Level1LimitDuration,
			LifetimeDuration:         supplyCommissionSetting.LifetimeDuration,
			Group:                    supplyCommissionSetting.Group,
			CustomerPolicyGroupID:    supplyCommissionSetting.CustomerPolicyGroupID,
			CreatedAt:                time.Now(),
			UpdatedAt:                time.Now(),
		}); err != nil {
			return err
		}

	}

	orderNotify.CommissionSnapshotStatus = 1
	return a.orderCreatedNotify(ctx).UpdateOrderCreatedNotify(orderNotify)
}

func (a *Aggregate) ProcessOrderNotify(ctx context.Context, orderCreatedNotifyID dot.ID) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		orderCreatedNotify, err := a.orderCreatedNotify(ctx).ID(orderCreatedNotifyID).GetOrderCreatedNotifyDB()
		if err != nil {
			return err
		}

		orderPromotions, err := a.orderPromotion(ctx).OrderCreatedNotifyID(orderCreatedNotify.ID).GetOrderPromotionsDB()
		if err != nil {
			return err
		}
		getOrderQ := &ordering.GetOrderByIDQuery{ID: orderCreatedNotify.OrderID}
		if err := a.orderingQuery.Dispatch(ctx, getOrderQ); err != nil {
			return err
		}

		shopQ := &identity.GetShopByIDQuery{
			ID: orderCreatedNotify.ShopID,
		}
		if err := a.identityQuery.Dispatch(ctx, shopQ); err != nil {
			return err
		}

		// process cashback for shop
		var shopCashbacks []*model.ShopCashback
		shopCashbackMap := map[dot.ID]*model.ShopCashback{}
		sellerPromotionByProductID := map[dot.ID]int{}
		for _, promotion := range orderPromotions {
			var cashback = float64(0)
			if promotion.Type != "cashback" {
				continue
			}
			switch promotion.Unit {
			case "vnd":
				cashback = cashback + float64(promotion.Amount)
			case "percent":
				cashback = cashback + (float64(promotion.BaseValue) * float64(promotion.Amount) / 100 / 100)
			}
			productQ := &catalog.GetShopProductByIDQuery{
				ProductID: promotion.ProductID,
				ShopID:    getOrderQ.Result.ShopID,
			}
			if err := a.catalogQuery.Dispatch(ctx, productQ); err != nil {
				return err
			}

			if promotion.Src == "seller" {
				sellerCashback := int(cashback)
				if promotion.Unit == "vnd" {
					sellerCashback = sellerCashback * promotion.ProductQuantity
				}

				if sellerPromotionByProductID[promotion.ProductID] == 0 {
					sellerPromotionByProductID[promotion.ProductID] = sellerCashback
				} else {
					sellerPromotionByProductID[promotion.ProductID] = sellerPromotionByProductID[promotion.ProductID] + sellerCashback
				}
			}

			if shopCashbackMap[productQ.ProductID] != nil {
				shopCashbackMap[productQ.ProductID].Amount = shopCashbackMap[productQ.ProductID].Amount + int(cashback)
			} else {
				now := time.Now()
				valid := now.Add(time.Hour * 24 * 3)
				shopCashback := &model.ShopCashback{
					ID:                   cm.NewID(),
					ShopID:               getOrderQ.Result.TradingShopID,
					OrderID:              getOrderQ.Result.ID,
					Amount:               int(cashback),
					OrderCreatedNotifyID: orderCreatedNotifyID,
					Description:          "Hoàn tiền khi mua sản phẩm " + productQ.Result.Name + " từ eTop Trading",
					Status:               0,
					ValidAt:              valid,
					CreatedAt:            now,
					UpdatedAt:            now,
				}
				shopCashbacks = append(shopCashbacks, shopCashback)
				shopCashbackMap[productQ.ProductID] = shopCashback
			}
		}
		for _, shopCashback := range shopCashbacks {
			err := a.shopCashback(ctx).CreateShopCashback(shopCashback)
			if err != nil {
				return err
			}
		}

		// process commission for seller
		affReferralCode, err := a.affiliateReferralCode(ctx).Code(orderCreatedNotify.ReferralCode).GetAffiliateReferralCodeDB()
		if err != nil {
			return err
		}
		userReferral, _ := a.userReferral(ctx).UserID(affReferralCode.UserID).GetUserReferralDB()

		for _, line := range getOrderQ.Result.Lines {
			orderCommissionSetting, err := a.orderCommissionSetting(ctx).SupplyID(getOrderQ.Result.ShopID).OrderID(getOrderQ.Result.ID).ProductID(line.ProductID).GetOrderCommissionSettingDB()
			if err != nil {
				continue
			}

			basePrice := float64(line.TotalPrice)
			tradingPromotion, err := a.productPromotion(ctx).ProductID(line.ProductID).GetProductPromotionDB()
			ll.Info("TRADING", l.Object("tradingPromotion", tradingPromotion), l.Object("err", err))
			if err == nil && tradingPromotion != nil {
				switch tradingPromotion.Unit {
				case "vnd":
					basePrice = basePrice - float64(tradingPromotion.Amount)
				case "percent":
					basePrice = math.Round(basePrice - basePrice*(float64(tradingPromotion.Amount)/100/100))
				}
			}
			var countByUser = math.MaxInt64
			var countByProduct = math.MaxInt64
			if orderCommissionSetting.DependOn == DependOnCustomer {
				countByUser, err = a.shopOrderProductHistory(ctx).UserID(shopQ.Result.OwnerID).CustomerPolicyGroup(orderCommissionSetting.CustomerPolicyGroupID).Count()
				if err != nil {
					return err
				}
			}
			if orderCommissionSetting.DependOn == DependOnProduct {
				countByProduct, err = a.shopOrderProductHistory(ctx).UserID(shopQ.Result.OwnerID).ProductID(line.ProductID).Count()
				if err != nil {
					return err
				}
			}

			var directCommission float64
			var indirectCommission float64
			if (orderCommissionSetting.DependOn == DependOnProduct && countByProduct < orderCommissionSetting.Level1LimitCount) ||
				(orderCommissionSetting.DependOn == DependOnCustomer && countByUser < orderCommissionSetting.Level1LimitCount) {
				directCommission = float64(orderCommissionSetting.Level1DirectCommission)
				indirectCommission = float64(orderCommissionSetting.Level1IndirectCommission)
			} else {
				directCommission = float64(orderCommissionSetting.Level2DirectCommission)
				indirectCommission = float64(orderCommissionSetting.Level2IndirectCommission)
			}

			directCommissionValue := math.Round(basePrice * (float64(directCommission) / 100 / 100))

			if sellerPromotionByProductID[line.ProductID] != 0 {
				directCommissionValue = directCommissionValue - float64(sellerPromotionByProductID[line.ProductID])
			}

			if err := a.affiliateCommission(ctx).CreateAffiliateCommission(&model.SellerCommission{
				ID:           cm.NewID(),
				SellerID:     affReferralCode.AffiliateID,
				FromSellerID: 0,
				ProductID:    line.ProductID,
				OrderId:      getOrderQ.Result.ID,
				ShopID:       getOrderQ.Result.TradingShopID,
				SupplyID:     getOrderQ.Result.ShopID,
				Amount:       int(directCommissionValue),
				Description:  "",
				Note:         "",
				Type:         "direct",
				Status:       0,
				OValue:       int(directCommission),
				OBaseValue:   int(basePrice),
				ValidAt:      time.Now().Add(time.Hour * 24 * 10),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}); err != nil {
				return err
			}

			if userReferral != nil && userReferral.ReferralID != 0 {
				indirectCommissionValue := math.Round(basePrice * (float64(indirectCommission) / 100 / 100))
				if err := a.affiliateCommission(ctx).CreateAffiliateCommission(&model.SellerCommission{
					ID:           cm.NewID(),
					SellerID:     userReferral.ReferralID,
					FromSellerID: affReferralCode.AffiliateID,
					ProductID:    line.ProductID,
					OrderId:      getOrderQ.Result.ID,
					ShopID:       getOrderQ.Result.TradingShopID,
					SupplyID:     getOrderQ.Result.ShopID,
					Amount:       int(indirectCommissionValue),
					Description:  "",
					Note:         "",
					Type:         "indirect",
					Status:       0,
					OValue:       int(indirectCommission),
					OBaseValue:   int(basePrice),
					ValidAt:      time.Now().Add(time.Hour * 24 * 10),
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}); err != nil {
					return err
				}
			}
		}

		// store history
		for _, line := range getOrderQ.Result.Lines {
			orderCommissionSetting, _ := a.orderCommissionSetting(ctx).SupplyID(getOrderQ.Result.ShopID).OrderID(getOrderQ.Result.ID).ProductID(line.ProductID).GetOrderCommissionSettingDB()

			if err := a.shopOrderProductHistory(ctx).CreateShopOrderProductHistory(&model.ShopOrderProductHistory{
				UserID:                shopQ.Result.OwnerID,
				ShopID:                getOrderQ.Result.TradingShopID,
				OrderID:               getOrderQ.Result.ID,
				SupplyID:              getOrderQ.Result.ShopID,
				ProductID:             line.ProductID,
				ProductQuantity:       line.Quantity,
				CustomerPolicyGroupID: orderCommissionSetting.CustomerPolicyGroupID,
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			}); err != nil {
				return err
			}
		}

		orderCreatedNotify.CashbackProcessStatus = 1
		orderCreatedNotify.CommissionProcessStatus = 1
		orderCreatedNotify.Status = 1
		if err := a.orderCreatedNotify(ctx).UpdateOrderCreatedNotify(orderCreatedNotify); err != nil {
			return err
		}

		return nil
	})
}

func (a *Aggregate) OrderPaymentSuccess(ctx context.Context, event *affiliate.OrderPaymentSuccessEvent) error {
	orderCreatedNotify, err := a.orderCreatedNotify(ctx).OrderID(event.OrderID).GetOrderCreatedNotifyDB()
	if err != nil {
		return err
	}

	if orderCreatedNotify.Status == 1 {
		return nil
	}

	orderCreatedNotify.PaymentStatus = 1
	if err := a.orderCreatedNotify(ctx).UpdateOrderCreatedNotify(orderCreatedNotify); err != nil {
		return err
	}

	go func() {
		if err := a.ProcessOrderNotify(bus.Ctx(), orderCreatedNotify.ID); err != nil {
			ll.Error("Processing Order created notify failed", l.Object("err", err))
			return
		}
		err := a.db.InTransaction(bus.Ctx(), func(tx cmsql.QueryInterface) error {
			commissions, err := a.affiliateCommission(ctx).OrderID(event.OrderID).GetAffiliateCommissionsDB()
			if err != nil {
				return err
			}
			for _, commission := range commissions {
				commission.Status = 2
				commission.ValidAt = time.Now().Add(time.Hour * 24 * 10)

				if err = a.affiliateCommission(ctx).UpdateAffiliateCommission(commission); err != nil {
					return err
				}
			}

			shopCashBacks, err := a.shopCashback(ctx).OrderID(event.OrderID).GetShopCashbacksDB()
			if err != nil {
				return err
			}
			for _, cashback := range shopCashBacks {
				cashback.Status = 2
				cashback.ValidAt = time.Now().Add(time.Hour * 24 * 3)
				if err = a.shopCashback(ctx).UpdateShopCashback(cashback); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			ll.Error("Processing Order created notify failed", l.Object("err", err))
			return
		}
		return
	}()

	return nil
}
