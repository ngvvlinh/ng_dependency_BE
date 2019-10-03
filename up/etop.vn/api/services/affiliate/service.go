package affiliate

import (
	"context"

	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateOrUpdateCommissionSetting(context.Context, *CreateCommissionSettingArgs) (*CommissionSetting, error)
	CreateProductPromotion(context.Context, *CreateProductPromotionArgs) (*ProductPromotion, error)
	UpdateProductPromotion(context.Context, *UpdateProductPromotionArgs) (*ProductPromotion, error)
	OnTradingOrderCreated(context.Context, *OnTradingOrderCreatedArgs) error
	TradingOrderCreating(context.Context, *TradingOrderCreating) error
	CreateAffiliateReferralCode(context.Context, *CreateReferralCodeArgs) (*AffiliateReferralCode, error)
	CreateOrUpdateUserReferral(context.Context, *CreateOrUpdateReferralArgs) (*UserReferral, error)
}

type CreateCommissionSettingArgs struct {
	ProductID   int64
	AccountID   int64
	Amount      int32
	Unit        string
	Type        string
	Description string
	Note        string
}

type CreateProductPromotionArgs struct {
	ShopID      int64
	ProductID   int64
	Amount      int32
	Code        string
	Description string
	Unit        string
	Note        string
	Type        string
}

type UpdateProductPromotionArgs struct {
	ID          int64
	Amount      int32
	Unit        string
	Code        string
	Description string
	Note        string
	Type        string
}

type OnTradingOrderCreatedArgs struct {
	OrderID      int64
	ReferralCode string
}

type TradingOrderCreating struct {
	ProductIDs   []int64
	ReferralCode string
	UserID       int64
}

type CreateReferralCodeArgs struct {
	AffiliateAccountID int64
	Code               string
}

type CreateOrUpdateReferralArgs struct {
	UserID           int64
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
}

type GetCommissionByProductIDsArgs struct {
	AccountID  int64
	ProductIDs []int64
}

type GetCommissionByProductIDArgs struct {
	AccountID int64
	ProductID int64
}

type ListShopProductPromotionsArgs struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters
}

type ListShopProductPromotionsResponse struct {
	Promotions []*ProductPromotion
	Count      int32
	Paging     meta.PageInfo
}

type GetProductPromotionArgs struct {
	ShopID    int64
	ProductID int64
}

type GetShopProductPromotionByProductIDs struct {
	ShopID     int64
	ProductIDs []int64
}

type GetAffiliateAccountReferralCodesArgs struct {
	AffiliateAccountID int64
}

type GetReferralsByReferralIDArgs struct {
	ID int64
}

type GetAffiliateReferralsArgs struct {
	Paging meta.Paging
}

type GetAffiliateAccountReferralByCodeArgs struct {
	Code string
}
