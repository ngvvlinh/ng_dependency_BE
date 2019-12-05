package affiliate

import (
	etop "etop.vn/api/top/int/etop"
	shop "etop.vn/api/top/int/shop"
	"etop.vn/api/top/int/types"
	common "etop.vn/api/top/types/common"
	status4 "etop.vn/api/top/types/etc/status4"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type UpdateReferralRequest struct {
	ReferralCode     string `json:"referral_code"`
	SaleReferralCode string `json:"sale_referral_code"`
}

func (m *UpdateReferralRequest) Reset()         { *m = UpdateReferralRequest{} }
func (m *UpdateReferralRequest) String() string { return jsonx.MustMarshalToString(m) }

type UserReferral struct {
	UserId           dot.ID `json:"user_id"`
	ReferralCode     string `json:"referral_code"`
	SaleReferralCode string `json:"sale_referral_code"`
}

func (m *UserReferral) Reset()         { *m = UserReferral{} }
func (m *UserReferral) String() string { return jsonx.MustMarshalToString(m) }

type SellerCommission struct {
	Id          dot.ID            `json:"id"`
	Value       int               `json:"value"`
	Description string            `json:"description"`
	Note        string            `json:"note"`
	Status      status4.Status    `json:"status"`
	Type        string            `json:"type"`
	OValue      int               `json:"o_value"`
	OBaseValue  int               `json:"o_base_value"`
	Product     *shop.ShopProduct `json:"product"`
	Order       *types.Order      `json:"order"`
	FromSeller  *etop.Affiliate   `json:"from_seller"`
	ValidAt     dot.Time          `json:"valid_at"`
	CreatedAt   dot.Time          `json:"created_at"`
	UpdatedAt   dot.Time          `json:"updated_at"`
}

func (m *SellerCommission) Reset()         { *m = SellerCommission{} }
func (m *SellerCommission) String() string { return jsonx.MustMarshalToString(m) }

type GetCommissionsResponse struct {
	Commissions []*SellerCommission `json:"commissions"`
}

func (m *GetCommissionsResponse) Reset()         { *m = GetCommissionsResponse{} }
func (m *GetCommissionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type NotifyNewShopPurchaseRequest struct {
	OrderId dot.ID `json:"order_id"`
}

func (m *NotifyNewShopPurchaseRequest) Reset()         { *m = NotifyNewShopPurchaseRequest{} }
func (m *NotifyNewShopPurchaseRequest) String() string { return jsonx.MustMarshalToString(m) }

type NotifyNewShopPurchaseResponse struct {
	Message string `json:"message"`
}

func (m *NotifyNewShopPurchaseResponse) Reset()         { *m = NotifyNewShopPurchaseResponse{} }
func (m *NotifyNewShopPurchaseResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetCouponsResponse struct {
	Coupons []*Coupon `json:"coupons"`
}

func (m *GetCouponsResponse) Reset()         { *m = GetCouponsResponse{} }
func (m *GetCouponsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateCouponRequest struct {
	Value       int            `json:"value"`
	Unit        dot.NullString `json:"unit"`
	Description dot.NullString `json:"description"`
	ProductId   dot.ID         `json:"product_id"`
}

func (m *CreateCouponRequest) Reset()         { *m = CreateCouponRequest{} }
func (m *CreateCouponRequest) String() string { return jsonx.MustMarshalToString(m) }

type Coupon struct {
	Id          dot.ID         `json:"id"`
	Code        string         `json:"code"`
	Value       int            `json:"value"`
	Unit        string         `json:"unit"`
	Description dot.NullString `json:"description"`
	UserId      dot.ID         `json:"user_id"`
	StartDate   dot.Time       `json:"start_date"`
	EndDate     dot.Time       `json:"end_date"`
	ProductId   dot.ID         `json:"product_id"`
	CreatedAt   dot.Time       `json:"created_at"`
	UpdatedAt   dot.Time       `json:"updated_at"`
}

func (m *Coupon) Reset()         { *m = Coupon{} }
func (m *Coupon) String() string { return jsonx.MustMarshalToString(m) }

type GetTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions"`
}

func (m *GetTransactionsResponse) Reset()         { *m = GetTransactionsResponse{} }
func (m *GetTransactionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type Transaction struct {
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return jsonx.MustMarshalToString(m) }

type CommissionSetting struct {
	ProductId dot.ID   `json:"product_id"`
	Amount    int      `json:"amount"`
	Unit      string   `json:"unit"`
	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *CommissionSetting) Reset()         { *m = CommissionSetting{} }
func (m *CommissionSetting) String() string { return jsonx.MustMarshalToString(m) }

type SupplyCommissionSetting struct {
	ProductId                dot.ID                                 `json:"product_id"`
	Level1DirectCommission   int                                    `json:"level1_direct_commission"`
	Level1IndirectCommission int                                    `json:"level1_indirect_commission"`
	Level2DirectCommission   int                                    `json:"level2_direct_commission"`
	Level2IndirectCommission int                                    `json:"level2_indirect_commission"`
	DependOn                 string                                 `json:"depend_on"`
	Level1LimitCount         int                                    `json:"level1_limit_count"`
	MLifetimeDuration        *SupplyCommissionSettingDurationObject `json:"m_lifetime_duration"`
	MLevel1LimitDuration     *SupplyCommissionSettingDurationObject `json:"m_level1_limit_duration"`
	CreatedAt                dot.Time                               `json:"created_at"`
	UpdatedAt                dot.Time                               `json:"updated_at"`
	Group                    string                                 `json:"group"`
}

func (m *SupplyCommissionSetting) Reset()         { *m = SupplyCommissionSetting{} }
func (m *SupplyCommissionSetting) String() string { return jsonx.MustMarshalToString(m) }

type SupplyCommissionSettingDurationObject struct {
	Duration int    `json:"duration"`
	Type     string `json:"type"`
}

func (m *SupplyCommissionSettingDurationObject) Reset()         { *m = SupplyCommissionSettingDurationObject{} }
func (m *SupplyCommissionSettingDurationObject) String() string { return jsonx.MustMarshalToString(m) }

type GetEtopCommissionSettingByProductIDsRequest struct {
	ProductIds []dot.ID `json:"product_ids"`
}

func (m *GetEtopCommissionSettingByProductIDsRequest) Reset() {
	*m = GetEtopCommissionSettingByProductIDsRequest{}
}
func (m *GetEtopCommissionSettingByProductIDsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetEtopCommissionSettingByProductIDsResponse struct {
	EtopCommissionSettings []*CommissionSetting `json:"etop_commission_settings"`
}

func (m *GetEtopCommissionSettingByProductIDsResponse) Reset() {
	*m = GetEtopCommissionSettingByProductIDsResponse{}
}
func (m *GetEtopCommissionSettingByProductIDsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetCommissionSettingByProductIDsRequest struct {
	ProductIds []dot.ID `json:"product_ids"`
}

func (m *GetCommissionSettingByProductIDsRequest) Reset() {
	*m = GetCommissionSettingByProductIDsRequest{}
}
func (m *GetCommissionSettingByProductIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCommissionSettingByProductIDsResponse struct {
	CommissionSettings []*CommissionSetting `json:"commission_settings"`
}

func (m *GetCommissionSettingByProductIDsResponse) Reset() {
	*m = GetCommissionSettingByProductIDsResponse{}
}
func (m *GetCommissionSettingByProductIDsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateOrUpdateCommissionSettingRequest struct {
	ProductId dot.ID         `json:"product_id"`
	Amount    int            `json:"amount"`
	Unit      dot.NullString `json:"unit"`
}

func (m *CreateOrUpdateCommissionSettingRequest) Reset() {
	*m = CreateOrUpdateCommissionSettingRequest{}
}
func (m *CreateOrUpdateCommissionSettingRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrUpdateTradingCommissionSettingRequest struct {
	ProductId                dot.ID `json:"product_id"`
	Level1DirectCommission   int    `json:"level1_direct_commission"`
	Level1IndirectCommission int    `json:"level1_indirect_commission"`
	Level2DirectCommission   int    `json:"level2_direct_commission"`
	Level2IndirectCommission int    `json:"level2_indirect_commission"`
	// product, customer
	DependOn         string `json:"depend_on"`
	Level1LimitCount int    `json:"level1_limit_count"`
	// day, month
	Level1LimitDurationType string `json:"level1_limit_duration_type"`
	Level1LimitDuration     int    `json:"level1_limit_duration"`
	// day, month
	LifetimeDurationType string `json:"lifetime_duration_type"`
	LifetimeDuration     int    `json:"lifetime_duration"`
	Group                string `json:"group"`
}

func (m *CreateOrUpdateTradingCommissionSettingRequest) Reset() {
	*m = CreateOrUpdateTradingCommissionSettingRequest{}
}
func (m *CreateOrUpdateTradingCommissionSettingRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type ProductPromotion struct {
	Product   *shop.ShopProduct `json:"product"`
	Id        dot.ID            `json:"id"`
	ProductId dot.ID            `json:"product_id"`
	Amount    int               `json:"amount"`
	Unit      string            `json:"unit"`
	Type      string            `json:"type"`
}

func (m *ProductPromotion) Reset()         { *m = ProductPromotion{} }
func (m *ProductPromotion) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrUpdateProductPromotionRequest struct {
	Id          dot.ID `json:"id"`
	ProductId   dot.ID `json:"product_id"`
	Amount      int    `json:"amount"`
	Unit        string `json:"unit"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Type        string `json:"type"`
}

func (m *CreateOrUpdateProductPromotionRequest) Reset()         { *m = CreateOrUpdateProductPromotionRequest{} }
func (m *CreateOrUpdateProductPromotionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetProductPromotionByProductIDRequest struct {
}

func (m *GetProductPromotionByProductIDRequest) Reset()         { *m = GetProductPromotionByProductIDRequest{} }
func (m *GetProductPromotionByProductIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetProductPromotionByProductIDResponse struct {
}

func (m *GetProductPromotionByProductIDResponse) Reset() {
	*m = GetProductPromotionByProductIDResponse{}
}
func (m *GetProductPromotionByProductIDResponse) String() string { return jsonx.MustMarshalToString(m) }

type SupplyProductResponse struct {
	Product                 *shop.ShopProduct        `json:"product"`
	SupplyCommissionSetting *SupplyCommissionSetting `json:"supply_commission_setting"`
	Promotion               *ProductPromotion        `json:"promotion"`
}

func (m *SupplyProductResponse) Reset()         { *m = SupplyProductResponse{} }
func (m *SupplyProductResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopProductResponse struct {
	Product   *shop.ShopProduct `json:"product"`
	Promotion *ProductPromotion `json:"promotion"`
}

func (m *ShopProductResponse) Reset()         { *m = ShopProductResponse{} }
func (m *ShopProductResponse) String() string { return jsonx.MustMarshalToString(m) }

type AffiliateProductResponse struct {
	Product                    *shop.ShopProduct  `json:"product"`
	ShopCommissionSetting      *CommissionSetting `json:"shop_commission_setting"`
	AffiliateCommissionSetting *CommissionSetting `json:"affiliate_commission_setting"`
	Promotion                  *ProductPromotion  `json:"promotion"`
}

func (m *AffiliateProductResponse) Reset()         { *m = AffiliateProductResponse{} }
func (m *AffiliateProductResponse) String() string { return jsonx.MustMarshalToString(m) }

type SupplyGetProductsResponse struct {
	Paging   *common.PageInfo         `json:"paging"`
	Products []*SupplyProductResponse `json:"products"`
}

func (m *SupplyGetProductsResponse) Reset()         { *m = SupplyGetProductsResponse{} }
func (m *SupplyGetProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopGetProductsResponse struct {
	Paging   *common.PageInfo       `json:"paging"`
	Products []*ShopProductResponse `json:"products"`
}

func (m *ShopGetProductsResponse) Reset()         { *m = ShopGetProductsResponse{} }
func (m *ShopGetProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CheckReferralCodeValidRequest struct {
	ProductId    dot.ID `json:"product_id"`
	ReferralCode string `json:"referral_code"`
}

func (m *CheckReferralCodeValidRequest) Reset()         { *m = CheckReferralCodeValidRequest{} }
func (m *CheckReferralCodeValidRequest) String() string { return jsonx.MustMarshalToString(m) }

type AffiliateGetProductsResponse struct {
	Paging   *common.PageInfo            `json:"paging"`
	Products []*AffiliateProductResponse `json:"products"`
}

func (m *AffiliateGetProductsResponse) Reset()         { *m = AffiliateGetProductsResponse{} }
func (m *AffiliateGetProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProductPromotionRequest struct {
	ProductId    dot.ID         `json:"product_id"`
	ReferralCode dot.NullString `json:"referral_code"`
}

func (m *GetProductPromotionRequest) Reset()         { *m = GetProductPromotionRequest{} }
func (m *GetProductPromotionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetProductPromotionResponse struct {
	Promotion        *ProductPromotion  `json:"promotion"`
	ReferralDiscount *CommissionSetting `json:"referral_discount"`
}

func (m *GetProductPromotionResponse) Reset()         { *m = GetProductPromotionResponse{} }
func (m *GetProductPromotionResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProductPromotionsResponse struct {
	Paging     *common.PageInfo    `json:"paging"`
	Promotions []*ProductPromotion `json:"promotions"`
}

func (m *GetProductPromotionsResponse) Reset()         { *m = GetProductPromotionsResponse{} }
func (m *GetProductPromotionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTradingProductPromotionByIDsRequest struct {
	ProductIds []dot.ID `json:"product_ids"`
}

func (m *GetTradingProductPromotionByIDsRequest) Reset() {
	*m = GetTradingProductPromotionByIDsRequest{}
}
func (m *GetTradingProductPromotionByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTradingProductPromotionByIDsResponse struct {
	Promotions []*ProductPromotion `json:"promotions"`
}

func (m *GetTradingProductPromotionByIDsResponse) Reset() {
	*m = GetTradingProductPromotionByIDsResponse{}
}
func (m *GetTradingProductPromotionByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateReferralCodeRequest struct {
	Code string `json:"code"`
}

func (m *CreateReferralCodeRequest) Reset()         { *m = CreateReferralCodeRequest{} }
func (m *CreateReferralCodeRequest) String() string { return jsonx.MustMarshalToString(m) }

type ReferralCode struct {
	Code string `json:"code"`
}

func (m *ReferralCode) Reset()         { *m = ReferralCode{} }
func (m *ReferralCode) String() string { return jsonx.MustMarshalToString(m) }

type GetReferralCodesResponse struct {
	ReferralCodes []*ReferralCode `json:"referral_codes"`
}

func (m *GetReferralCodesResponse) Reset()         { *m = GetReferralCodesResponse{} }
func (m *GetReferralCodesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Referral struct {
	Name            string   `json:"name"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	OrderCount      int      `json:"order_count"`
	TotalRevenue    int      `json:"total_revenue"`
	TotalCommission int      `json:"total_commission"`
	CreatedAt       dot.Time `json:"created_at"`
}

func (m *Referral) Reset()         { *m = Referral{} }
func (m *Referral) String() string { return jsonx.MustMarshalToString(m) }

type GetReferralsResponse struct {
	Referrals []*Referral `json:"referrals"`
}

func (m *GetReferralsResponse) Reset()         { *m = GetReferralsResponse{} }
func (m *GetReferralsResponse) String() string { return jsonx.MustMarshalToString(m) }
