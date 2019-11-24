package affiliate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateOrUpdateCommissionSetting(context.Context, *CreateCommissionSettingArgs) (*CommissionSetting, error)
	CreateOrUpdateSupplyCommissionSetting(context.Context, *CreateOrUpdateSupplyCommissionSettingArgs) (*SupplyCommissionSetting, error)
	CreateProductPromotion(context.Context, *CreateProductPromotionArgs) (*ProductPromotion, error)
	UpdateProductPromotion(context.Context, *UpdateProductPromotionArgs) (*ProductPromotion, error)
	OnTradingOrderCreated(context.Context, *OnTradingOrderCreatedArgs) error
	TradingOrderCreating(context.Context, *TradingOrderCreating) error
	CreateAffiliateReferralCode(context.Context, *CreateReferralCodeArgs) (*AffiliateReferralCode, error)
	CreateOrUpdateUserReferral(context.Context, *CreateOrUpdateReferralArgs) (*UserReferral, error)
	OrderPaymentSuccess(context.Context, *OrderPaymentSuccessEvent) error
}

type CreateCommissionSettingArgs struct {
	ProductID   dot.ID
	AccountID   dot.ID
	Amount      int32
	Unit        string
	Type        string
	Description string
	Note        string
}

type CreateOrUpdateSupplyCommissionSettingArgs struct {
	ShopID                   dot.ID
	ProductID                dot.ID
	Level1DirectCommission   int32
	Level1IndirectCommission int32
	Level2DirectCommission   int32
	Level2IndirectCommission int32
	DependOn                 string
	Level1LimitCount         int32
	Level1LimitDuration      int32
	Level1LimitDurationType  string
	LifetimeDuration         int32
	LifetimeDurationType     string
	Group                    string
}

type CreateProductPromotionArgs struct {
	ShopID      dot.ID
	ProductID   dot.ID
	Amount      int32
	Code        string
	Description string
	Unit        string
	Note        string
	Type        string
}

type UpdateProductPromotionArgs struct {
	ID          dot.ID
	Amount      int32
	Unit        string
	Code        string
	Description string
	Note        string
	Type        string
}

type OnTradingOrderCreatedArgs struct {
	OrderID      dot.ID
	ReferralCode string
}

type TradingOrderCreating struct {
	ProductIDs   []dot.ID
	ReferralCode string
	UserID       dot.ID
}

type CreateReferralCodeArgs struct {
	AffiliateAccountID dot.ID
	Code               string
}

type CreateOrUpdateReferralArgs struct {
	UserID           dot.ID
	ReferralCode     string
	SaleReferralCode string
}

type QueryService interface {
	GetCommissionByProductIDs(context.Context, *GetCommissionByProductIDsArgs) ([]*CommissionSetting, error)
	GetCommissionByProductID(context.Context, *GetCommissionByProductIDArgs) (*CommissionSetting, error)
	ListShopProductPromotions(context.Context, *ListShopProductPromotionsArgs) (*ListShopProductPromotionsResponse, error)
	GetShopProductPromotion(context.Context, *GetProductPromotionArgs) (*ProductPromotion, error)
	GetShopProductPromotionByProductIDs(context.Context, *GetShopProductPromotionByProductIDs) ([]*ProductPromotion, error)
	GetAffiliateAccountReferralCodes(context.Context, *GetAffiliateAccountReferralCodesArgs) ([]*AffiliateReferralCode, error)
	GetReferralsByReferralID(context.Context, *GetReferralsByReferralIDArgs) ([]*UserReferral, error)
	GetAffiliateAccountReferralByCode(context.Context, *GetAffiliateAccountReferralByCodeArgs) (*AffiliateReferralCode, error)
	GetSupplyCommissionSettingsByProductIDs(context.Context, *GetSupplyCommissionSettingsByProductIDsArgs) ([]*SupplyCommissionSetting, error)
	GetSellerCommissions(context.Context, *GetSellerCommissionsArgs) ([]*SellerCommission, error)
}

type GetCommissionByProductIDsArgs struct {
	AccountID  dot.ID
	ProductIDs []dot.ID
}

type GetCommissionByProductIDArgs struct {
	AccountID dot.ID
	ProductID dot.ID
}

type ListShopProductPromotionsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListShopProductPromotionsResponse struct {
	Promotions []*ProductPromotion
	Count      int32
	Paging     meta.PageInfo
}

type GetProductPromotionArgs struct {
	ShopID    dot.ID
	ProductID dot.ID
}

type GetShopProductPromotionByProductIDs struct {
	ShopID     dot.ID
	ProductIDs []dot.ID
}

type GetAffiliateAccountReferralCodesArgs struct {
	AffiliateAccountID dot.ID
}

type GetReferralsByReferralIDArgs struct {
	ID dot.ID
}

type GetAffiliateReferralsArgs struct {
	Paging meta.Paging
}

type GetAffiliateAccountReferralByCodeArgs struct {
	Code string
}

type GetSupplyCommissionSettingsByProductIDsArgs struct {
	ShopID     dot.ID
	ProductIDs []dot.ID
}

type GetSellerCommissionsArgs struct {
	SellerID dot.ID
	Paging   meta.Paging
	Filters  meta.Filters
}
