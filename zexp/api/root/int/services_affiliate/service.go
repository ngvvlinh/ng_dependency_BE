package services_affiliate

import (
	"context"

	cm "etop.vn/backend/pb/common"
	aff "etop.vn/backend/pb/services/affiliate"
)

// +gen:apix

// +apix:path=/affiliate.User
type UserAPI interface {
	UpdateReferral(context.Context, *aff.UpdateReferralRequest) (*aff.UserReferral, error)
}

// +apix:path=/affiliate.Trading
type TradingAPI interface {
	TradingGetProducts(context.Context, *cm.CommonListRequest) (*aff.SupplyGetProductsResponse, error)
	CreateOrUpdateTradingCommissionSetting(context.Context, *aff.CreateOrUpdateTradingCommissionSettingRequest) (*aff.SupplyCommissionSetting, error)
	GetTradingProductPromotionByProductIDs(context.Context, *aff.GetTradingProductPromotionByIDsRequest) (*aff.GetTradingProductPromotionByIDsResponse, error)
	GetTradingProductPromotions(context.Context, *cm.CommonListRequest) (*aff.GetProductPromotionsResponse, error)
	CreateTradingProductPromotion(context.Context, *aff.CreateOrUpdateProductPromotionRequest) (*aff.ProductPromotion, error)
	UpdateTradingProductPromotion(context.Context, *aff.CreateOrUpdateProductPromotionRequest) (*aff.ProductPromotion, error)
}

// +apix:path=/affiliate.Shop
type ShopAPI interface {
	GetProductPromotion(context.Context, *aff.GetProductPromotionRequest) (*aff.GetProductPromotionResponse, error)
	ShopGetProducts(context.Context, *cm.CommonListRequest) (*aff.ShopGetProductsResponse, error)
	CheckReferralCodeValid(context.Context, *aff.CheckReferralCodeValidRequest) (*aff.GetProductPromotionResponse, error)
}

// +apix:path=/affiliate.Affiliate
type AffiliateAPI interface {
	GetCommissions(context.Context, *cm.CommonListRequest) (*aff.GetCommissionsResponse, error)
	NotifyNewShopPurchase(context.Context, *aff.NotifyNewShopPurchaseRequest) (*aff.NotifyNewShopPurchaseResponse, error)
	GetTransactions(context.Context, *cm.CommonListRequest) (*aff.GetTransactionsResponse, error)
	CreateOrUpdateAffiliateCommissionSetting(context.Context, *aff.CreateOrUpdateCommissionSettingRequest) (*aff.CommissionSetting, error)
	GetProductPromotionByProductID(context.Context, *aff.GetProductPromotionByProductIDRequest) (*aff.GetProductPromotionByProductIDResponse, error)
	AffiliateGetProducts(context.Context, *cm.CommonListRequest) (*aff.AffiliateGetProductsResponse, error)
	CreateReferralCode(context.Context, *aff.CreateReferralCodeRequest) (*aff.ReferralCode, error)
	GetReferralCodes(context.Context, *cm.CommonListRequest) (*aff.GetReferralCodesResponse, error)
	GetReferrals(context.Context, *cm.CommonListRequest) (*aff.GetReferralsResponse, error)
}
