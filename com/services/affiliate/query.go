package affiliate

import (
	"context"

	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ affiliate.QueryService = &QueryService{}

type QueryService struct {
	commissionSetting     sqlstore.CommissionSettingStoreFactory
	productPromotion      sqlstore.ProductPromotionStoreFactory
	affiliateReferralCode sqlstore.AffiliateReferralCodeStoreFactory
	userReferral          sqlstore.UserReferralStoreFactory
}

func NewQuery(db cmsql.Database) *QueryService {
	return &QueryService{
		commissionSetting:     sqlstore.NewCommissionSettingStore(db),
		productPromotion:      sqlstore.NewProductPromotionStore(db),
		affiliateReferralCode: sqlstore.NewAffiliateReferralCodeStore(db),
		userReferral:          sqlstore.NewUserReferralStore(db),
	}
}

func (a *QueryService) MessageBus() affiliate.QueryBus {
	b := bus.New()
	return affiliate.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (a *QueryService) GetCommissionByProductIDs(ctx context.Context, args *affiliate.GetCommissionByProductIDsArgs) ([]*affiliate.CommissionSetting, error) {
	return a.commissionSetting(ctx).AccountID(args.AccountID).ProductIDs(args.ProductIDs).GetCommissionSettings()
}

func (a *QueryService) GetCommissionByProductID(ctx context.Context, args *affiliate.GetCommissionByProductIDArgs) (*affiliate.CommissionSetting, error) {
	return a.commissionSetting(ctx).AccountID(args.AccountID).ProductID(args.ProductID).GetCommissionSetting()
}

func (a *QueryService) ListShopProductPromotions(ctx context.Context, args *affiliate.ListShopProductPromotionsArgs) (*affiliate.ListShopProductPromotionsResponse, error) {
	query := a.productPromotion(ctx).ShopID(args.ShopID).Filters(args.Filters)
	promotions, err := query.Paging(args.Paging).GetProductPromotions()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &affiliate.ListShopProductPromotionsResponse{
		Promotions: promotions,
		Count:      int32(count),
		Paging:     query.GetPaging(),
	}, nil
}

func (a *QueryService) GetShopProductPromotion(ctx context.Context, args *affiliate.GetProductPromotionArgs) (*affiliate.ProductPromotion, error) {
	return a.productPromotion(ctx).ShopID(args.ShopID).ProductID(args.ProductID).GetProductPromotion()
}

func (a *QueryService) GetShopProductPromotionByProductIDs(ctx context.Context, args *affiliate.GetShopProductPromotionByProductIDs) ([]*affiliate.ProductPromotion, error) {
	return a.productPromotion(ctx).ShopID(args.ShopID).ProductIDs(args.ProductIDs...).GetProductPromotions()
}

func (a *QueryService) GetAffiliateAccountReferralCodes(ctx context.Context, args *affiliate.GetAffiliateAccountReferralCodesArgs) ([]*affiliate.AffiliateReferralCode, error) {
	return a.affiliateReferralCode(ctx).AffiliateID(args.AffiliateAccountID).GetAffiliateReferralCodes()
}

func (a *QueryService) GetReferralsByReferralID(ctx context.Context, args *affiliate.GetReferralsByReferralIDArgs) ([]*affiliate.UserReferral, error) {
	return a.userReferral(ctx).ReferralID(args.ID).GetUserReferrals()
}

func (q *QueryService) GetAffiliateAccountReferralByCode(ctx context.Context, args *affiliate.GetAffiliateAccountReferralByCodeArgs) (*affiliate.AffiliateReferralCode, error) {
	return q.affiliateReferralCode(ctx).Code(args.Code).GetAffiliateReferralCode()
}
