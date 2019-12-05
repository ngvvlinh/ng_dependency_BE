package affiliate

import (
	"context"

	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=services/affiliate

// +apix:path=/affiliate.User
type UserService interface {
	UpdateReferral(context.Context, *UpdateReferralRequest) (*UserReferral, error)
}

// +apix:path=/affiliate.Trading
type TradingService interface {
	TradingGetProducts(context.Context, *cm.CommonListRequest) (*SupplyGetProductsResponse, error)
	CreateOrUpdateTradingCommissionSetting(context.Context, *CreateOrUpdateTradingCommissionSettingRequest) (*SupplyCommissionSetting, error)
	GetTradingProductPromotionByProductIDs(context.Context, *GetTradingProductPromotionByIDsRequest) (*GetTradingProductPromotionByIDsResponse, error)
	GetTradingProductPromotions(context.Context, *cm.CommonListRequest) (*GetProductPromotionsResponse, error)
	CreateTradingProductPromotion(context.Context, *CreateOrUpdateProductPromotionRequest) (*ProductPromotion, error)
	UpdateTradingProductPromotion(context.Context, *CreateOrUpdateProductPromotionRequest) (*ProductPromotion, error)
}

// +apix:path=/affiliate.Shop
type ShopService interface {
	GetProductPromotion(context.Context, *GetProductPromotionRequest) (*GetProductPromotionResponse, error)
	ShopGetProducts(context.Context, *cm.CommonListRequest) (*ShopGetProductsResponse, error)
	CheckReferralCodeValid(context.Context, *CheckReferralCodeValidRequest) (*GetProductPromotionResponse, error)
}

// +apix:path=/affiliate.Affiliate
type AffiliateService interface {
	GetCommissions(context.Context, *cm.CommonListRequest) (*GetCommissionsResponse, error)
	NotifyNewShopPurchase(context.Context, *NotifyNewShopPurchaseRequest) (*NotifyNewShopPurchaseResponse, error)
	GetTransactions(context.Context, *cm.CommonListRequest) (*GetTransactionsResponse, error)
	CreateOrUpdateAffiliateCommissionSetting(context.Context, *CreateOrUpdateCommissionSettingRequest) (*CommissionSetting, error)
	GetProductPromotionByProductID(context.Context, *GetProductPromotionByProductIDRequest) (*GetProductPromotionByProductIDResponse, error)
	AffiliateGetProducts(context.Context, *cm.CommonListRequest) (*AffiliateGetProductsResponse, error)
	CreateReferralCode(context.Context, *CreateReferralCodeRequest) (*ReferralCode, error)
	GetReferralCodes(context.Context, *cm.CommonListRequest) (*GetReferralCodesResponse, error)
	GetReferrals(context.Context, *cm.CommonListRequest) (*GetReferralsResponse, error)
}
