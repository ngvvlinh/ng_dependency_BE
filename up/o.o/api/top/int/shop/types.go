package shop

import (
	"time"

	catalogtypes "o.o/api/main/catalog/types"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/top/int/etop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/int/types"
	"o.o/api/top/int/types/spreadsheet"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/ledger_type"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/product_type"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/ref_action"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/shop_user_role"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/stocktake_type"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/api/top/types/etc/ticket/ticket_type"
	"o.o/api/top/types/etc/try_on"
	"o.o/api/top/types/etc/ws_banner_show_style"
	"o.o/api/top/types/etc/ws_list_product_show_style"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
	"o.o/common/xerrors"
)

type GetTicketCommentsResponse struct {
	TicketComments []*shoptypes.TicketComment `json:"ticket_comments"`
	Paging         *common.PageInfo           `json:"paging"`
}

func (m *GetTicketCommentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketLabelsResponse struct {
	TicketLabels []*shoptypes.TicketLabel `json:"ticket_labels"`
}

func (m *GetTicketLabelsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketCommentsRequest struct {
	TicketID dot.ID                  `json:"ticket_id"`
	Filter   *FilterGetTicketComment `json:"filter"`
	Paging   *common.Paging          `json:"paging"`
}

type FilterGetTicketComment struct {
	IDs       []dot.ID `json:"ids"`
	Title     string   `json:"title"`
	CreatedBy dot.ID   `json:"created_by"`
	ParentID  dot.ID   `json:"parent_id"`
}

func (m *FilterGetTicketComment) String() string { return jsonx.MustMarshalToString(m) }

func (m *GetTicketCommentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateTicketCommentRequest struct {
	ID        dot.ID   `json:"id"`
	Message   string   `json:"message"`
	ImageUrl  string   `json:"image_url"`
	ImageUrls []string `json:"image_urls"`
}

func (m *UpdateTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateTicketCommentRequest struct {
	TicketID  dot.ID   `json:"ticket_id"`
	Message   string   `json:"message"`
	ImageUrl  string   `json:"image_url"`
	ImageUrls []string `json:"image_urls"`
	ParentID  dot.ID   `json:"parent_id"`
}

func (m *CreateTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketCommentRequest struct {
	ID dot.ID `json:"id"`
}

func (m *DeleteTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketCommentResponse struct {
	Count int `json:"count"`
}

func (m *DeleteTicketCommentResponse) String() string { return jsonx.MustMarshalToString(m) }

type FilterShopGetTicket struct {
	IDs             []dot.ID                      `json:"ids"`
	CreatedBy       dot.ID                        `json:"created_by"`
	ConfirmedBy     dot.ID                        `json:"confirmed_by"`
	ClosedBy        dot.ID                        `json:"closed_by"`
	AccountID       dot.ID                        `json:"account_id"`
	LabelIDs        []dot.ID                      `json:"label_ids"`
	AssignedUserIDs []dot.ID                      `json:"assigned_user_ids"`
	RefID           dot.ID                        `json:"ref_id"`
	RefType         ticket_ref_type.TicketRefType `json:"ref_type"`
	RefCode         string                        `json:"ref_code"`
	State           ticket_state.TicketState      `json:"state"`
	Title           filter.FullTextSearch         `json:"title"`
	Code            string                        `json:"code"`
	Types           []ticket_type.TicketType      `json:"types"`
}

func (m *FilterShopGetTicket) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketsRequest struct {
	Paging *common.Paging       `json:"paging"`
	Filter *FilterShopGetTicket `json:"filter"`
}

func (m *GetTicketsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketsResponse struct {
	Paging  *common.PageInfo    `json:"paging"`
	Tickets []*shoptypes.Ticket `json:"tickets"`
}

func (m *GetTicketsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
}

func (m *ConfirmTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type ReopenTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
}

func (m *ReopenTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type CloseTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
	// @required
	State ticket_state.TicketState `json:"state"`
}

func (m *CloseTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type UnassignTicketRequest struct {
	AssignedUserIDs []dot.ID `json:"assigned_user_ids"`
	TicketID        dot.ID   `json:"ticket_id"`
}

func (m *UnassignTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type AssignTicketRequest struct {
	AssignedUserIDs []dot.ID `json:"assigned_user_ids"`
	TicketID        dot.ID   `json:"ticket_id"`
}

func (m *AssignTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateTicketRequest struct {
	LabelIDs []dot.ID `json:"label_ids"`

	Title       string `json:"title"`
	Description string `json:"description"`

	// user note
	Note string `json:"note"`

	RefID   dot.ID                        `json:"ref_id"`
	RefType ticket_ref_type.TicketRefType `json:"ref_type"`
	RefCode string                        `json:"ref_code"`
	Source  ticket_source.TicketSource    `json:"source"`

	Type ticket_type.NullTicketType `json:"type"`
}

func (m *CreateTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateTicketRequest struct {
	ID          dot.ID                        `json:"id"`
	LabelIDs    []dot.ID                      `json:"label_ids"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	RefID       dot.ID                        `json:"ref_id"`
	RefType     ticket_ref_type.TicketRefType `json:"ref_type"`
}

func (m *UpdateTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *UpdateTicketRequest) Validate() error {
	if m.ID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing ID")
	}
	return nil
}

type UpdateTicketRefTicketIDRequest struct {
	ID dot.ID `json:"id"`
	// Truy???n l??n 0 ????? x??a ref_ticket_id
	RefTicketID dot.NullID `json:"ref_ticket_id"`
}

func (m *UpdateTicketRefTicketIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketLabelRequest struct {
	ID          dot.ID `json:"id"`
	DeleteChild bool   `json:"delete_child"`
}

func (m *DeleteTicketLabelRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetTicketLabelsRequest struct {
	Filter *FilterGetTicketLabels `json:"filter"`
	Tree   bool                   `json:"tree"`
}

func (m *GetTicketLabelsRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetTicketLabels struct {
	Type ticket_type.NullTicketType `json:"type"`
}

func (m *FilterGetTicketLabels) String() string { return jsonx.MustMarshalToString(m) }

type CreateTicketLabelRequest struct {
	ParentID dot.ID `json:"parent_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Code     string `json:"code"`
}

func (m *CreateTicketLabelRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateTicketLabelRequest struct {
	ID       dot.ID         `json:"id"`
	Name     dot.NullString `json:"name"`
	Color    string         `json:"color"`
	Code     dot.NullString `json:"code"`
	ParentID dot.NullID     `json:"parent_id"`
}

func (m *UpdateTicketLabelRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketLabelResponse struct {
	Count int `json:"deleted"`
}

func (m *DeleteTicketLabelResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsCategoriesByIDsResponse struct {
	WsCategories []*WsCategory `json:"ws_categorys"`
}

func (m *GetWsCategoriesByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsCategoriesResponse struct {
	WsCategories []*WsCategory    `json:"ws_categorys"`
	Paging       *common.PageInfo `json:"paging"`
}

func (m *GetWsCategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsCategoriesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetWsCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsCategoriesByIDsRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *GetWsCategoriesByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsCategoryRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetWsCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrUpdateWsCategoryRequest struct {
	CategoryID dot.ID         `json:"category_id"`
	Slug       dot.NullString `json:"slug"`
	SEOConfig  *WsSEOConfig   `json:"seo_config"`
	Image      dot.NullString `json:"image"`
	Appear     dot.NullBool   `json:"appear"`
}

func (m *CreateOrUpdateWsCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type WsCategory struct {
	ID           dot.ID        `json:"id"`
	ShopID       dot.ID        `json:"shop_id"`
	Slug         string        `json:"slug"`
	SEOConfig    *WsSEOConfig  `json:"seo_config"`
	Image        string        `json:"image"`
	Appear       bool          `json:"appear"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Category     *ShopCategory `json:"category"`
	ProductCount int           `json:"product_count"`
}

func (m *WsCategory) String() string { return jsonx.MustMarshalToString(m) }

type DeteleWsPageRequest struct {
	ID dot.ID `json:"id"`
}

func (m *DeteleWsPageRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeteleWsPageResponse struct {
	Count int `json:"count"`
}

func (m *DeteleWsPageResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateWsPageRequest struct {
	ID        dot.ID         `json:"id"`
	SEOConfig *WsSEOConfig   `json:"seo_config"`
	Name      dot.NullString `json:"name"`
	Slug      dot.NullString `json:"slug"`
	DescHTML  dot.NullString `json:"desc_html"`
	Image     dot.NullString `json:"image"`
	Appear    dot.NullBool   `json:"appear"`
}

func (m *UpdateWsPageRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateWsPageRequest struct {
	SEOConfig *WsSEOConfig `json:"seo_config"`
	Name      string       `json:"name"`
	Slug      string       `json:"slug"`
	DescHTML  string       `json:"desc_html"`
	Image     string       `json:"image"`
	Appear    bool         `json:"appear"`
}

func (m *CreateWsPageRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsPagesByIDsResponse struct {
	WsPages []*WsPage `json:"ws_pages"`
}

func (m *GetWsPagesByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsPagesResponse struct {
	WsPages []*WsPage        `json:"ws_pages"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *GetWsPagesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsPagesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetWsPagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsPagesByIDsRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *GetWsPagesByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsPageRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetWsPageRequest) String() string { return jsonx.MustMarshalToString(m) }

type WsPage struct {
	Name      string       `json:"name"`
	ShopID    dot.ID       `json:"shop_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	ID        dot.ID       `json:"id"`
	SEOConfig *WsSEOConfig `json:"seo_config"`
	Slug      string       `json:"slug"`
	Appear    bool         `json:"appear"`
	DescHTML  string       `json:"desc_html"`
}

func (m *WsPage) String() string { return jsonx.MustMarshalToString(m) }

type GetWsProductsByIDsResponse struct {
	WsProducts []*WsProduct `json:"ws_products"`
}

func (m *GetWsProductsByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsProductsResponse struct {
	WsProducts []*WsProduct     `json:"ws_products"`
	Paging     *common.PageInfo `json:"paging"`
}

func (m *GetWsProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsProductsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetWsProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsProductsByIDsRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *GetWsProductsByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsProductRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetWsProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type WsProduct struct {
	ShopID       dot.ID          `json:"shop_id"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ID           dot.ID          `json:"id"`
	SEOConfig    *WsSEOConfig    `json:"seo_config"`
	Slug         string          `json:"slug"`
	Appear       bool            `json:"appear"`
	ComparePrice []*ComparePrice `json:"compare_prices"`
	DescHTML     string          `json:"desc_html"`
	Product      *ShopProduct    `json:"shop_product"`
	Sale         bool            `json:"Sale"`
}

func (m *WsProduct) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrUpdateWsProductRequest struct {
	ProductID     dot.ID          `json:"product_id"`
	SEOConfig     *WsSEOConfig    `json:"seo_config"`
	Slug          dot.NullString  `json:"slug"`
	Appear        dot.NullBool    `json:"appear"`
	ComparePrices []*ComparePrice `json:"compare_prices"`
	DescHTML      dot.NullString  `json:"desc_html"`
}

type ComparePrice struct {
	VariantID    dot.ID `json:"variant_id"`
	ComparePrice int    `json:"compare_price"`
}

func (m *CreateOrUpdateWsProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type WsSEOConfig struct {
	Content     string `json:"content"`
	Keyword     string `json:"keyword"`
	Description string `json:"description"`
}

func (m *WsSEOConfig) String() string { return jsonx.MustMarshalToString(m) }

type GetWsWebsitesByIDsResponse struct {
	WsWebsites []*WsWebsite `json:"ws_websites"`
}

func (m *GetWsWebsitesByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsWebsitesResponse struct {
	WsWebsites []*WsWebsite     `json:"ws_websites"`
	Paging     *common.PageInfo `json:"paging"`
}

func (m *GetWsWebsitesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWsWebsitesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetWsWebsitesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsWebsitesByIDsRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *GetWsWebsitesByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetWsWebsiteRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetWsWebsiteRequest) String() string { return jsonx.MustMarshalToString(m) }

type WsWebsite struct {
	ShopID             dot.ID          `json:"shop_id"`
	ID                 dot.ID          `json:"id"`
	MainColor          string          `json:"main_color"`
	Banner             *Banner         `json:"banner"`
	OutstandingProduct *SpecialProduct `json:"outstanding_product"`
	NewProduct         *SpecialProduct `json:"new_product"`
	SEOConfig          *WsGeneralSEO   `json:"seo_config"`
	Facebook           *Facebook       `json:"facebook"`
	GoogleAnalyticsID  string          `json:"google_analytics_id"`
	DomainName         string          `json:"domain_name"`
	OverStock          bool            `json:"over_stock"`
	ShopInfo           *ShopInfo       `json:"shop_info"`
	Description        string          `json:"description"`
	LogoImage          string          `json:"logo_image"`
	FaviconImage       string          `json:"favicon_image"`
	DeletedAt          time.Time       `json:"deleted_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CreatedAt          time.Time       `json:"created_at"`
	SiteSubdomain      string          `json:"site_subdomain"`
}

func (m *WsWebsite) String() string { return jsonx.MustMarshalToString(m) }

type CreateCreditRequest struct {
	Amount   int                        `json:"amount"`
	Type     credit_type.CreditType     `json:"type"`
	Classify credit_type.CreditClassify `json:"classify"`
}

func (m *CreateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateWsWebsiteRequest struct {
	MainColor          string          `json:"main_color"`
	Banner             *Banner         `json:"banner"`
	OutstandingProduct *SpecialProduct `json:"outstanding_product"`
	NewProduct         *SpecialProduct `json:"new_product"`
	SEOConfig          *WsGeneralSEO   `json:"seo_config"`
	Facebook           *Facebook       `json:"facebook"`
	GoogleAnalyticsID  string          `json:"google_analytics_id"`
	DomainName         string          `json:"domain_name"`
	OverStock          bool            `json:"over_stock"`
	ShopInfo           *ShopInfo       `json:"shop_info"`
	Description        string          `json:"description"`
	LogoImage          string          `json:"logo_image"`
	FaviconImage       string          `json:"favicon_image"`
	SiteSubdomain      string          `json:"site_subdomain"`
}

func (m *CreateWsWebsiteRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateWsWebsiteRequest struct {
	ID                 dot.ID          `json:"id"`
	MainColor          dot.NullString  `json:"main_color"`
	Banner             *Banner         `json:"banner"`
	OutstandingProduct *SpecialProduct `json:"outstanding_product"`
	NewProduct         *SpecialProduct `json:"new_product"`
	SEOConfig          *WsGeneralSEO   `json:"seo_config"`
	Facebook           *Facebook       `json:"facebook"`
	GoogleAnalyticsID  dot.NullString  `json:"google_analytics_id"`
	DomainName         dot.NullString  `json:"domain_name"`
	OverStock          dot.NullBool    `json:"over_stock"`
	ShopInfo           *ShopInfo       `json:"shop_info"`
	Description        dot.NullString  `json:"description"`
	LogoImage          dot.NullString  `json:"logo_image"`
	FaviconImage       dot.NullString  `json:"favicon_image"`
	SiteSubdomain      dot.NullString  `json:"site_subdomain"`
}

func (m *UpdateWsWebsiteRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopInfo struct {
	Email           string           `json:"email"`
	Name            string           `json:"name"`
	Phone           string           `json:"phone"`
	Address         *AddressShopInfo `json:"address"`
	FacebookFanpage string           `json:"facebook_fanpage"`
}

type AddressShopInfo struct {
	Province     string `json:"province"`
	ProvinceCode string `json:"province_code"`
	District     string `json:"district"`
	DistrictCode string `json:"district_code"`
	Ward         string `json:"ward"`
	WardCode     string `json:"ward_code"`
	Address      string `json:"address"`
}

func (m *AddressShopInfo) String() string { return jsonx.MustMarshalToString(m) }

func (m *ShopInfo) String() string { return jsonx.MustMarshalToString(m) }

type Facebook struct {
	FacebookID     string `json:"facebook_id"`
	WelcomeMessage string `json:"welcome_message"`
}

func (m *Facebook) String() string { return jsonx.MustMarshalToString(m) }

type WsGeneralSEO struct {
	Title               string `json:"title"`
	SiteContent         string `json:"site_content"`
	SiteMetaKeyword     string `json:"site_meta_keyword"`
	SiteMetaDescription string `json:"site_meta_description"`
}

func (m *WsGeneralSEO) String() string { return jsonx.MustMarshalToString(m) }

type Banner struct {
	BannerItems []*BannerItem                          `json:"banner_items"`
	Style       ws_banner_show_style.WsBannerShowStyle `json:"style"`
}

func (m *Banner) String() string { return jsonx.MustMarshalToString(m) }

type BannerItem struct {
	Alt   string `json:"alt"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

func (m *BannerItem) String() string { return jsonx.MustMarshalToString(m) }

type SpecialProduct struct {
	ProductIDs []dot.ID                                          `json:"product_ids"`
	Style      ws_list_product_show_style.WsListProductShowStyle `json:"style"`
	Products   []*ShopProduct                                    `json:"products"`
}

func (m *SpecialProduct) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrder struct {
	Id                      dot.ID                  `json:"id"`
	ShopId                  dot.ID                  `json:"shop_id"`
	SupplierId              dot.ID                  `json:"supplier_id"`
	Supplier                *PurchaseOrderSupplier  `json:"supplier"`
	BasketValue             int                     `json:"basket_value"`
	DiscountLines           []*types.DiscountLine   `json:"discount_lines"`
	TotalDiscount           int                     `json:"total_discount"`
	FeeLines                []*types.FeeLine        `json:"fee_lines"`
	TotalFee                int                     `json:"total_fee"`
	PurchaseOrderAdjustment []*types.AdjustmentLine `json:"purchase_order_adjustment"`
	TotalAmount             int                     `json:"total_amount"`
	Code                    string                  `json:"code"`
	Note                    string                  `json:"note"`
	Status                  status3.Status          `json:"status"`
	Lines                   []*PurchaseOrderLine    `json:"lines"`
	CreatedBy               dot.ID                  `json:"created_by"`
	CancelReason            string                  `json:"cancel_reason"`
	ConfirmedAt             dot.Time                `json:"confirmed_at"`
	CancelledAt             dot.Time                `json:"cancelled_at"`
	CreatedAt               dot.Time                `json:"created_at"`
	UpdatedAt               dot.Time                `json:"updated_at"`
	DeletedAt               dot.Time                `json:"deleted_at"`
	InventoryVoucher        *InventoryVoucher       `json:"inventory_voucher"`
	PaidAmount              int                     `json:"paid_amount"`
}

func (m *PurchaseOrder) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrderSupplier struct {
	FullName           string `json:"full_name"`
	Phone              string `json:"phone"`
	Email              string `json:"email"`
	CompanyName        string `json:"company_name"`
	TaxNumber          string `json:"tax_number"`
	HeadquarterAddress string `json:"headquarter_address"`
	Deleted            bool   `json:"deleted"`
}

func (m *PurchaseOrderSupplier) String() string { return jsonx.MustMarshalToString(m) }

type GetPurchaseOrdersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetPurchaseOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrdersResponse struct {
	PurchaseOrders []*PurchaseOrder `json:"purchase_orders"`
	Paging         *common.PageInfo `json:"paging"`
}

func (m *PurchaseOrdersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreatePurchaseOrderRequest struct {
	Lines         []*PurchaseOrderLine  `json:"lines"`
	BasketValue   int                   `json:"basket_value"`
	FeeLines      []*types.FeeLine      `json:"fee_lines"`
	TotalFee      int                   `json:"total_fee"`
	DiscountLines []*types.DiscountLine `json:"discount_lines"`
	TotalDiscount int                   `json:"total_discount"`
	TotalAmount   int                   `json:"total_amount"`
	Note          string                `json:"note"`
	SupplierId    dot.ID                `json:"supplier_id"`
}

func (m *CreatePurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePurchaseOrderRequest struct {
	Id            dot.ID                `json:"id"`
	TotalAmount   dot.NullInt           `json:"total_amount"`
	Note          dot.NullString        `json:"note"`
	Lines         []*PurchaseOrderLine  `json:"lines"`
	BasketValue   dot.NullInt           `json:"basket_value"`
	FeeLines      []*types.FeeLine      `json:"fee_lines"`
	TotalFee      dot.NullInt           `json:"total_fee"`
	DiscountLines []*types.DiscountLine `json:"discount_lines"`
	TotalDiscount dot.NullInt           `json:"total_discount"`
}

func (m *UpdatePurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrderLine struct {
	VariantId    dot.ID                    `json:"variant_id"`
	Quantity     int                       `json:"quantity"`
	PaymentPrice int                       `json:"payment_price"`
	ProductId    dot.ID                    `json:"product_id"`
	ProductName  string                    `json:"product_name"`
	ImageUrl     string                    `json:"image_url"`
	Code         string                    `json:"code"`
	Attributes   []*catalogtypes.Attribute `json:"attributes"`
	Discount     int                       `json:"discount"`
}

func (m *PurchaseOrderLine) String() string { return jsonx.MustMarshalToString(m) }

type CancelPurchaseOrderRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
	// @deprecated use cancel_reason instead
	Reason               string                              `json:"reason"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *CancelPurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmPurchaseOrderRequest struct {
	Id dot.ID `json:"id"`
	// enum create, confirm
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *ConfirmPurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetLedgersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetLedgersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateLedgerRequest struct {
	Name        string            `json:"name"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        string            `json:"note"`
}

func (m *CreateLedgerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateLedgerRequest struct {
	Id          dot.ID            `json:"id"`
	Name        dot.NullString    `json:"name"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        dot.NullString    `json:"note"`
}

func (m *UpdateLedgerRequest) String() string { return jsonx.MustMarshalToString(m) }

type LedgersResponse struct {
	Ledgers []*Ledger        `json:"ledgers"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *LedgersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Ledger struct {
	Id          dot.ID            `json:"id"`
	Name        string            `json:"name"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        string            `json:"note"`
	// enum: cash, bank
	Type      ledger_type.LedgerType `json:"type"`
	CreatedBy dot.ID                 `json:"created_by"`
	CreatedAt dot.Time               `json:"created_at"`
	UpdatedAt dot.Time               `json:"updated_at"`
}

func (m *Ledger) String() string { return jsonx.MustMarshalToString(m) }

type RegisterShopRequest struct {
	// @required
	Name        string            `json:"name"`
	Address     *etop.Address     `json:"address"`
	Phone       string            `json:"phone"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	WebsiteUrl  dot.NullString    `json:"website_url"`
	ImageUrl    string            `json:"image_url"`
	Email       string            `json:"email"`
	UrlSlug     string            `json:"url_slug"`
	CompanyInfo *etop.CompanyInfo `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
}

func (m *RegisterShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterShopResponse struct {
	// @required
	Shop *etop.Shop `json:"shop"`
}

func (m *RegisterShopResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShopRequest struct {
	InventoryOverstock dot.NullBool      `json:"inventory_overstock"`
	Name               string            `json:"name"`
	Address            *etop.Address     `json:"address"`
	Phone              string            `json:"phone"`
	BankAccount        *etop.BankAccount `json:"bank_account"`
	WebsiteUrl         dot.NullString    `json:"website_url"`
	ImageUrl           string            `json:"image_url"`
	Email              string            `json:"email"`
	AutoCreateFfm      dot.NullBool      `json:"auto_create_ffm"`
	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.NullGHNNoteCode `json:"ghn_note_code"`
	TryOn       try_on.NullTryOnCode          `json:"try_on"`
	CompanyInfo *etop.CompanyInfo             `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
}

func (m *UpdateShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShopResponse struct {
	Shop *etop.Shop `json:"shop"`
}

func (m *UpdateShopResponse) String() string { return jsonx.MustMarshalToString(m) }

type Collection struct {
	// @required
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
	// @required
	Name        string `json:"name"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`
	DescHtml    string `json:"desc_html"`
	// @required
	CreatedAt dot.Time `json:"created_at"`
	// @required
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *Collection) String() string { return jsonx.MustMarshalToString(m) }

type CreateCollectionRequest struct {
	// @required
	Name        string `json:"name"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`
	DescHtml    string `json:"desc_html"`
}

func (m *CreateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductCategoryRequest struct {
	ProductId  dot.ID `json:"product_id"`
	CategoryId dot.ID `json:"category_id"`
}

func (m *UpdateProductCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type CollectionsResponse struct {
	Collections []*ShopCollection `json:"collections"`
}

func (m *CollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCollectionRequest struct {
	// @required
	Id dot.ID `json:"id"`
	// @required
	Name        dot.NullString `json:"name"`
	Description dot.NullString `json:"description"`
	ShortDesc   dot.NullString `json:"short_desc"`
	DescHtml    dot.NullString `json:"desc_html"`
}

func (m *UpdateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsCollectionRequest struct {
	// @required
	CollectionId dot.ID   `json:"collection_id"`
	ProductIds   []dot.ID `json:"product_ids"`
}

func (m *UpdateProductsCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveProductsCollectionRequest struct {
	// @required
	CollectionId dot.ID   `json:"collection_id"`
	ProductIds   []dot.ID `json:"product_ids"`
}

func (m *RemoveProductsCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type EtopVariant struct {
	Id          dot.ID                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	ShortDesc   string                    `json:"short_desc"`
	DescHtml    string                    `json:"desc_html"`
	ImageUrls   []string                  `json:"image_urls"`
	ListPrice   int                       `json:"list_price"`
	CostPrice   int                       `json:"cost_price"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
}

func (m *EtopVariant) String() string { return jsonx.MustMarshalToString(m) }

type EtopProduct struct {
	Id          dot.ID   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	DescHtml    string   `json:"desc_html"`
	Unit        string   `json:"unit"`
	ImageUrls   []string `json:"image_urls"`
	ListPrice   int      `json:"list_price"`
	CostPrice   int      `json:"cost_price"`
	CategoryId  dot.ID   `json:"category_id"`
	// @deprecated
	ProductSourceCategoryId dot.ID `json:"product_source_category_id"`
}

func (m *EtopProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariant struct {
	// @required
	Id   dot.ID       `json:"id"`
	Info *EtopVariant `json:"info"`

	Code string `json:"code"`
	// @deprecated use code instead
	EdCode string `json:"ed_code"`

	Name        string         `json:"name"`
	Description string         `json:"description"`
	ShortDesc   string         `json:"short_desc"`
	DescHtml    string         `json:"desc_html"`
	ImageUrls   []string       `json:"image_urls"`
	ListPrice   int            `json:"list_price"`
	RetailPrice int            `json:"retail_price"`
	Note        string         `json:"note"`
	Status      status3.Status `json:"status"`
	IsAvailable bool           `json:"is_available"`

	QuantityOnHand int `json:"quantity_on_hand"`
	QuantityPicked int `json:"quantity_picked"`
	CostPrice      int `json:"cost_price"`
	Quantity       int `json:"quantity"`
	// @deprecated
	InventoryVariant *InventoryVariantShopVariant `json:"inventory_variant"`

	// @deprecated use stags instead
	Tags  []string `json:"tags"`
	Stags []*Tag   `json:"stags"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`
	Product    *ShopShortProduct         `json:"product"`
	ProductId  dot.ID                    `json:"product_id"`
}

func (m *ShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVariantShopVariant struct {
	QuantityOnHand int `json:"quantity_on_hand"`
	QuantityPicked int `json:"quantity_picked"`
	CostPrice      int `json:"cost_price"`
	Quantity       int `json:"quantity"`
}

func (m *InventoryVariantShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type ShopProduct struct {
	// @required
	Id   dot.ID       `json:"id"`
	Info *EtopProduct `json:"info"`
	Code string       `json:"code"`
	// @deprecated use code instead
	EdCode      string   `json:"ed_code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	DescHtml    string   `json:"desc_html"`
	ImageUrls   []string `json:"image_urls"`
	CategoryId  dot.ID   `json:"category_id"`
	// @deprecated use stags instead
	Tags            []string                     `json:"tags"`
	Stags           []*Tag                       `json:"stags"`
	Note            string                       `json:"note"`
	Status          status3.Status               `json:"status"`
	IsAvailable     bool                         `json:"is_available"`
	ListPrice       int                          `json:"list_price"`
	RetailPrice     int                          `json:"retail_price"`
	CollectionIds   []dot.ID                     `json:"collection_ids"`
	Variants        []*ShopVariant               `json:"variants"`
	ProductSourceId dot.ID                       `json:"product_source_id"`
	CreatedAt       dot.Time                     `json:"created_at"`
	UpdatedAt       dot.Time                     `json:"updated_at"`
	ProductType     product_type.NullProductType `json:"product_type"`
	MetaFields      []*common.MetaField          `json:"meta_fields"`
	BrandId         dot.ID                       `json:"brand_id"`
}

func (m *ShopProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopShortProduct struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *ShopShortProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopCollection struct {
	Id          dot.ID `json:"id"`
	ShopId      dot.ID `json:"shop_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DescHtml    string `json:"desc_html"`
	ShortDesc   string `json:"short_desc"`
}

func (m *ShopCollection) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetVariantsRequest struct {
	Name filter.FullTextSearch `json:"name"`
}

func (m *FilterGetVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantsRequest struct {
	Paging  *common.Paging            `json:"paging"`
	Filters []*common.Filter          `json:"filters"`
	Filter  *FilterGetVariantsRequest `json:"filter"`
}

func (m *GetVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCategoriesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariantsResponse struct {
	Paging   *common.PageInfo `json:"paging"`
	Variants []*ShopVariant   `json:"variants"`
}

func (m *ShopVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopProductsResponse struct {
	Paging   *common.PageInfo `json:"paging"`
	Products []*ShopProduct   `json:"products"`
}

func (m *ShopProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopCategoriesResponse struct {
	Paging     *common.PageInfo `json:"paging"`
	Categories []*ShopCategory  `json:"categories"`
}

func (m *ShopCategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantRequest struct {
	// @required
	Id          dot.ID                    `json:"id"`
	Name        dot.NullString            `json:"name"`
	Note        dot.NullString            `json:"note"`
	Code        dot.NullString            `json:"code"`
	CostPrice   dot.NullInt               `json:"cost_price"`
	ListPrice   dot.NullInt               `json:"list_price"`
	RetailPrice dot.NullInt               `json:"retail_price"`
	Description dot.NullString            `json:"description"`
	ShortDesc   dot.NullString            `json:"short_desc"`
	DescHtml    dot.NullString            `json:"desc_html"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
	// deprecated
	Sku string `json:"sku"`
}

func (m *UpdateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductRequest struct {
	// @required
	Id          dot.ID                       `json:"id"`
	Name        dot.NullString               `json:"name"`
	Code        dot.NullString               `json:"code"`
	Note        dot.NullString               `json:"note"`
	Unit        dot.NullString               `json:"unit"`
	Description dot.NullString               `json:"description"`
	ShortDesc   dot.NullString               `json:"short_desc"`
	DescHtml    dot.NullString               `json:"desc_html"`
	CostPrice   dot.NullInt                  `json:"cost_price"`
	ListPrice   dot.NullInt                  `json:"list_price"`
	RetailPrice dot.NullInt                  `json:"retail_price"`
	CategoryID  dot.NullID                   `json:"category_id"`
	ProductType product_type.NullProductType `json:"product_type"`
	MetaFields  *common.MetaField            `json:"meta_fields"`
	BrandId     dot.NullID                   `json:"brand_id"`
}

func (m *UpdateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCategoryRequest struct {
	Id       dot.ID         `json:"id"`
	Name     dot.NullString `json:"name"`
	ParentId dot.ID         `json:"parent_id"`
}

func (m *UpdateCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantsRequest struct {
	Updates []*UpdateVariantRequest `json:"updates"`
}

func (m *UpdateVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsTagsRequest struct {
	// @required
	Ids        []dot.ID `json:"ids"`
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`
}

func (m *UpdateProductsTagsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantsResponse struct {
	Variants []*ShopVariant  `json:"variants"`
	Errors   []*common.Error `json:"errors"`
}

func (m *UpdateVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type AddVariantsRequest struct {
	Ids          []dot.ID `json:"ids"`
	Tags         []string `json:"tags"`
	CollectionId dot.ID   `json:"collection_id"`
}

func (m *AddVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddVariantsResponse struct {
	Variants []*ShopVariant  `json:"variants"`
	Errors   []*common.Error `json:"errors"`
}

func (m *AddVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type RemoveVariantsRequest struct {
	Ids []dot.ID `json:"ids"`
}

func (m *RemoveVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetOrdersRequest struct {
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
	Mixed   *etop.MixedAccount `json:"mixed"`
}

func (m *GetOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateOrdersStatusRequest struct {
	// @required
	Ids []dot.ID `json:"ids"`
	// @required
	Confirm      status3.NullStatus `json:"confirm"`
	CancelReason string             `json:"cancel_reason"`
	Status       status5.Status     `json:"status"`
}

func (m *UpdateOrdersStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmOrderRequest struct {
	OrderId dot.ID `json:"order_id"`
	// enum ('create', 'create')
	AutoInventoryVoucher  inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
	AutoCreateFulfillment bool                                `json:"auto_create_fulfillment"`
}

func (m *ConfirmOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDRequest struct {
	OrderId dot.ID `json:"order_id"`
}

func (m *OrderIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDsRequest struct {
	OrderIds []dot.ID `json:"order_ids"`
}

func (m *OrderIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateFulfillmentsForOrderRequest struct {
	OrderId    dot.ID   `json:"order_id"`
	VariantIds []dot.ID `json:"variant_ids"`
}

func (m *CreateFulfillmentsForOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelOrderRequest struct {
	OrderId              dot.ID                              `json:"order_id"`
	CancelReason         string                              `json:"cancel_reason"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *CancelOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelOrdersRequest struct {
	Ids    []dot.ID `json:"ids"`
	Reason string   `json:"reason"`
}

func (m *CancelOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductSource struct {
	Id        dot.ID         `json:"id"`
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Status    status3.Status `json:"status"`
	UpdatedAt dot.Time       `json:"updated_at"`
	CreatedAt dot.Time       `json:"created_at"`
}

func (m *ProductSource) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type CreateProductSourceRequest struct {
}

func (m *CreateProductSourceRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductSourcesResponse struct {
	ProductSources []*ProductSource `json:"product_sources"`
}

func (m *ProductSourcesResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
}

func (m *CreateCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateProductRequest struct {
	Code        string                       `json:"code"`
	Name        string                       `json:"name"`
	Unit        string                       `json:"unit"`
	Note        string                       `json:"note"`
	Description string                       `json:"description"`
	ShortDesc   string                       `json:"short_desc"`
	DescHtml    string                       `json:"desc_html"`
	ImageUrls   []string                     `json:"image_urls"`
	CostPrice   int                          `json:"cost_price"`
	ListPrice   int                          `json:"list_price"`
	RetailPrice int                          `json:"retail_price"`
	ProductType product_type.NullProductType `json:"product_type"`
	BrandId     dot.ID                       `json:"brand_id"`
	MetaFields  []*common.MetaField          `json:"meta_fields"`
}

func (m *CreateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateVariantRequest struct {
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	ProductId   dot.ID                    `json:"product_id"`
	Note        string                    `json:"note"`
	Description string                    `json:"description"`
	ShortDesc   string                    `json:"short_desc"`
	DescHtml    string                    `json:"desc_html"`
	ImageUrls   []string                  `json:"image_urls"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
	CostPrice   int                       `json:"cost_price"`
	ListPrice   int                       `json:"list_price"`
	RetailPrice int                       `json:"retail_price"`
}

func (m *CreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeprecatedCreateVariantRequest struct {
	// required
	ProductSourceId dot.ID `json:"product_source_id"`
	ProductId       dot.ID `json:"product_id"`
	// In `D??p Adidas Adilette Slides - Full ?????`, product_name is "D??p Adidas Adilette Slides"
	ProductName string `json:"product_name"`
	// In `D??p Adidas Adilette Slides - Full ?????`, name is "Full ?????"
	Name              string                    `json:"name"`
	Description       string                    `json:"description"`
	ShortDesc         string                    `json:"short_desc"`
	DescHtml          string                    `json:"desc_html"`
	ImageUrls         []string                  `json:"image_urls"`
	Tags              []string                  `json:"tags"`
	Status            status3.Status            `json:"status"`
	CostPrice         int                       `json:"cost_price"`
	ListPrice         int                       `json:"list_price"`
	RetailPrice       int                       `json:"retail_price"`
	Code              string                    `json:"code"`
	QuantityAvailable int                       `json:"quantity_available"`
	QuantityOnHand    int                       `json:"quantity_on_hand"`
	QuantityReserved  int                       `json:"quantity_reserved"`
	Attributes        []*catalogtypes.Attribute `json:"attributes"`
	Unit              string                    `json:"unit"`
	// deprecated: use code instead
	Sku string `json:"sku"`
}

func (m *DeprecatedCreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConnectProductSourceResquest struct {
	ProductSourceId dot.ID `json:"product_source_id"`
}

func (m *ConnectProductSourceResquest) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type CreatePSCategoryRequest struct {
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
}

func (m *CreatePSCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsPSCategoryRequest struct {
	CategoryId dot.ID   `json:"category_id"`
	ProductIds []dot.ID `json:"product_ids"`
}

func (m *UpdateProductsPSCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsCollectionResponse struct {
	Updated int             `json:"updated"`
	Errors  []*common.Error `json:"errors"`
}

func (m *UpdateProductsCollectionResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductSourceCategoryRequest struct {
	Id       dot.ID `json:"id"`
	ParentId dot.ID `json:"parent_id"`
	Name     string `json:"name"`
}

func (m *UpdateProductSourceCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type GetProductSourceCategoriesRequest struct {
}

func (m *GetProductSourceCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsRequest struct {
	Paging        *common.Paging     `json:"paging"`
	Filters       []*common.Filter   `json:"filters"`
	Mixed         *etop.MixedAccount `json:"mixed"`
	OrderId       dot.ID             `json:"order_id"`
	Status        status3.NullStatus `json:"status"`
	ConnectionIDs []dot.ID           `json:"connection_ids"`
}

func (m *GetFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsByIDsRequest struct {
	IDs []dot.ID `json:"ids"`
}

func (m *GetFulfillmentsByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentHistoryRequest struct {
	Paging  *common.Paging `json:"paging"`
	All     bool           `json:"all"`
	Id      dot.ID         `json:"id"`
	OrderId dot.ID         `json:"order_id"`
}

func (m *GetFulfillmentHistoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetPublicFulfillmentRequest struct {
	// @Required
	Code string `json:"code"`
}

func (m *GetPublicFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentsShippingStateRequest struct {
	// Only support for manual order
	Ids []dot.ID `json:"ids"`
	// @required
	ShippingState shipping.State `json:"shipping_state"`
}

func (m *UpdateFulfillmentsShippingStateRequest) Reset() {
	*m = UpdateFulfillmentsShippingStateRequest{}
}
func (m *UpdateFulfillmentsShippingStateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateOrderPaymentStatusRequest struct {
	OrderId dot.ID             `json:"order_id"`
	Status  status4.NullStatus `json:"status"`
}

func (m *UpdateOrderPaymentStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeFulfillmentsRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type GetExternalPaymenUrlRequest struct {
	Type            subject_referral.SubjectReferral `json:"type"`
	RefID           dot.ID                           `json:"ref_id"`
	ReturnURL       string                           `json:"return_url"`
	CancelURL       string                           `json:"cancel_url"`
	PaymentProvider payment_provider.PaymentProvider `json:"payment_provider"`
}

func (m *GetExternalPaymenUrlRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePaymentStatusRequest struct {
	ID     dot.ID                     `json:"id"`
	State  payment_state.PaymentState `json:"state"`
	Status status4.Status             `json:"status"`
}

func (m *UpdatePaymentStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExternalPaymentUrlResponse struct {
	PaymentUrl string `json:"payment_url"`
}

func (m *GetExternalPaymentUrlResponse) String() string { return jsonx.MustMarshalToString(m) }

func (m *SummarizeFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeFulfillmentsResponse struct {
	Tables []*SummaryTable `json:"tables"`
}

func (m *SummarizeFulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeTopShipRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (m *SummarizeTopShipRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeTopShipResponse struct {
	Tables []*SummaryTable `json:"tables"`
}

func (m *SummarizeTopShipResponse) String() string { return jsonx.MustMarshalToString(m) }

type SummarizePOSResponse struct {
	Tables []*SummaryTable `json:"tables"`
}

func (m *SummarizePOSResponse) String() string { return jsonx.MustMarshalToString(m) }

type SummarizePOSRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (m *SummarizePOSRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummaryTable struct {
	Label   string          `json:"label"`
	Tags    []string        `json:"tags"`
	Columns []SummaryColRow `json:"columns"`
	Rows    []SummaryColRow `json:"rows"`
	Data    []SummaryItem   `json:"data"`
}

func (m *SummaryTable) String() string { return jsonx.MustMarshalToString(m) }

type SummaryColRow struct {
	Label  string `json:"label"`
	Spec   string `json:"spec"`
	Unit   string `json:"unit"`
	Indent int    `json:"indent"`
}

func (m *SummaryColRow) String() string { return jsonx.MustMarshalToString(m) }

type SummaryItem struct {
	Label     string   `json:"label"`
	Spec      string   `json:"spec"`
	Value     int64    `json:"value"`
	Unit      string   `json:"unit"`
	ImageUrls []string `json:"image_urls"`
}

func (m *SummaryItem) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductsResponse struct {
	Data         *spreadsheet.SpreadsheetData `json:"data"`
	ImportErrors []*common.Error              `json:"import_errors"`
	CellErrors   []*common.Error              `json:"cell_errors"`
	ImportId     dot.ID                       `json:"import_id"`
	StocktakeID  dot.ID                       `json:"stocktake_id"`
}

func (m *ImportProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CalcBalanceUserRequest struct {
	// @deprecated: use service_classify instead
	CreditClassify credit_type.NullCreditClassify `json:"credit_classify"`

	ServiceClassify service_classify.NullServiceClassify `json:"service_classify"`
}

func (m *CalcBalanceUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type CalcBalanceUserResponse struct {
	// S??? d???ng cho v???n chuy???n
	AvailableBalance int `json:"available_balance"`
	// S??? d???ng cho v???n chuy???n
	ActualBalance int `json:"actual_balance"`
	// S??? d???ng cho Telecom
	TelecomBalance int `json:"telecom_balance"`
}

func (m *CalcBalanceUserResponse) String() string { return jsonx.MustMarshalToString(m) }

type RequestExportRequest struct {
	ExportType string           `json:"export_type"`
	Filters    []*common.Filter `json:"filters"`
	DateFrom   string           `json:"date_from"`
	DateTo     string           `json:"date_to"`
	// Accept '\t', ',' or ';'. Default to ','.
	Delimiter string `json:"delimiter"`
	// For exporting csv compatible with Excel
	ExcelCompatibleMode bool `json:"excel_compatible_mode"`
	// Export specific ids
	Ids []dot.ID `json:"ids"`
}

func (m *RequestExportRequest) String() string { return jsonx.MustMarshalToString(m) }

type RequestExportResponse struct {
	Id         string         `json:"id"`
	Filename   string         `json:"filename"`
	ExportType string         `json:"export_type"`
	Status     status4.Status `json:"status"`
}

func (m *RequestExportResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsRequest struct {
}

func (m *GetExportsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsResponse struct {
	ExportItems []*ExportItem `json:"export_items"`
}

func (m *GetExportsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ExportItem struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	// example: shop/fulfillments, admin/orders
	ExportType   string          `json:"export_type"`
	DownloadUrl  string          `json:"download_url"`
	AccountId    dot.ID          `json:"account_id"`
	UserId       dot.ID          `json:"user_id"`
	CreatedAt    dot.Time        `json:"created_at"`
	DeletedAt    dot.Time        `json:"deleted_at"`
	RequestQuery string          `json:"request_query"`
	MimeType     string          `json:"mime_type"`
	Status       status4.Status  `json:"status"`
	ExportErrors []*common.Error `json:"export_errors"`
	Error        *common.Error   `json:"error"`
}

func (m *ExportItem) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsStatusRequest struct {
}

func (m *GetExportsStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type ExportStatusItem struct {
	Id            string        `json:"id"`
	ProgressMax   int           `json:"progress_max"`
	ProgressValue int           `json:"progress_value"`
	ProgressError int           `json:"progress_error"`
	Error         *common.Error `json:"error"`
}

func (m *ExportStatusItem) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizePartnerRequest struct {
	PartnerId dot.ID `json:"partner_id"`
}

func (m *AuthorizePartnerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetPartnersResponse struct {
	Partners []*etop.PublicAccountInfo `json:"partners"`
}

func (m *GetPartnersResponse) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizedPartnerResponse struct {
	Partner     *etop.PublicAccountInfo `json:"partner"`
	RedirectUrl string                  `json:"redirect_url"`
}

func (m *AuthorizedPartnerResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetAuthorizedPartnersResponse struct {
	Partners []*AuthorizedPartnerResponse `json:"partners"`
}

func (m *GetAuthorizedPartnersResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantImagesRequest struct {
	// @required
	Id         dot.ID   `json:"id"`
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`
}

func (m *UpdateVariantImagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductMetaFieldsRequest struct {
	// @required
	Id         dot.ID              `json:"id"`
	MetaFields []*common.MetaField `json:"meta_fields"`
}

func (m *UpdateProductMetaFieldsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (m *CategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Category struct {
	Id              dot.ID `json:"id"`
	Name            string `json:"name"`
	ProductSourceId dot.ID `json:"product_source_id"`
	ParentId        dot.ID `json:"parent_id"`
	ShopId          dot.ID `json:"shop_id"`
}

func (m *Category) String() string { return jsonx.MustMarshalToString(m) }

type Tag struct {
	Id    dot.ID `json:"id"`
	Label string `json:"label"`
}

func (m *Tag) String() string { return jsonx.MustMarshalToString(m) }

type ExternalAccountAhamove struct {
	Id               dot.ID `json:"id"`
	OwnerID          dot.ID `json:"owner_id"`
	Phone            string `json:"phone"`
	Name             string `json:"name"`
	ExternalVerified bool   `json:"external_verified"`
	//    optional string external_token = 5 [(gogoproto.nullable) = false];
	ExternalCreatedAt   dot.Time `json:"external_created_at"`
	CreatedAt           dot.Time `json:"created_at"`
	UpdatedAt           dot.Time `json:"updated_at"`
	LastSendVerifyAt    dot.Time `json:"last_send_verify_at"`
	ExternalTicketId    string   `json:"external_ticket_id"`
	IdCardFrontImg      string   `json:"id_card_front_img"`
	IdCardBackImg       string   `json:"id_card_back_img"`
	PortraitImg         string   `json:"portrait_img"`
	UploadedAt          dot.Time `json:"uploaded_at"`
	FanpageUrl          string   `json:"fanpage_url"`
	WebsiteUrl          string   `json:"website_url"`
	CompanyImgs         []string `json:"company_imgs"`
	BusinessLicenseImgs []string `json:"business_license_imgs"`
	ConnectionID        dot.ID   `json:"connection_id"`
}

func (m *ExternalAccountAhamove) String() string { return jsonx.MustMarshalToString(m) }

type UpdateXAccountAhamoveVerificationRequest struct {
	Id                  dot.ID   `json:"id"`
	IdCardFrontImg      string   `json:"id_card_front_img"`
	IdCardBackImg       string   `json:"id_card_back_img"`
	PortraitImg         string   `json:"portrait_img"`
	FanpageUrl          string   `json:"fanpage_url"`
	WebsiteUrl          string   `json:"website_url"`
	CompanyImgs         []string `json:"company_imgs"`
	BusinessLicenseImgs []string `json:"business_license_imgs"`
}

func (m *UpdateXAccountAhamoveVerificationRequest) Reset() {
	*m = UpdateXAccountAhamoveVerificationRequest{}
}
func (m *UpdateXAccountAhamoveVerificationRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CustomerLiability struct {
	TotalOrders    int `json:"total_orders"`
	TotalAmount    int `json:"total_amount"`
	ReceivedAmount int `json:"received_amount"`
	Liability      int `json:"liability"`
}

func (m *CustomerLiability) String() string { return jsonx.MustMarshalToString(m) }

type Customer struct {
	Id        dot.ID                     `json:"id"`
	ShopId    dot.ID                     `json:"shop_id"`
	FullName  string                     `json:"full_name"`
	Code      string                     `json:"code"`
	Note      string                     `json:"note"`
	Phone     string                     `json:"phone"`
	Email     string                     `json:"email"`
	Gender    gender.Gender              `json:"gender"`
	Type      customer_type.CustomerType `json:"type"`
	Birthday  string                     `json:"birthday"`
	CreatedAt dot.Time                   `json:"created_at"`
	UpdatedAt dot.Time                   `json:"updated_at"`
	Status    status3.Status             `json:"status"`
	Deleted   bool                       `json:"deleted"`
	GroupIds  []dot.ID                   `json:"group_ids"`
	Liability *CustomerLiability         `json:"liability"`
}

func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerRequest struct {
	// @required
	FullName string        `json:"full_name"`
	Gender   gender.Gender `json:"gender"`
	Birthday string        `json:"birthday"`
	// enum ('individual', 'organization')
	Type customer_type.CustomerType `json:"type"`
	Note string                     `json:"note"`
	// @required
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *CreateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerRequest struct {
	Id       dot.ID            `json:"id"`
	FullName dot.NullString    `json:"full_name"`
	Gender   gender.NullGender `json:"gender"`
	Birthday dot.NullString    `json:"birthday"`
	// enum ('individual', 'organization','independent')
	Type  customer_type.NullCustomerType `json:"type"`
	Note  dot.NullString                 `json:"note"`
	Phone dot.NullString                 `json:"phone"`
	Email dot.NullString                 `json:"email"`
}

func (m *UpdateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetCustomersRequest struct {
	FullName filter.FullTextSearch `json:"full_name"`
}

func (m *FilterGetCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomersRequest struct {
	Paging  *common.Paging             `json:"paging"`
	Filters []*common.Filter           `json:"filters"`
	GetAll  dot.NullBool               `json:"get_all"`
	Filter  *FilterGetCustomersRequest `json:"filter"`
}

func (m *GetCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomersResponse struct {
	Customers []*Customer      `json:"customers"`
	Paging    *common.PageInfo `json:"paging"`
}

func (m *CustomersResponse) String() string { return jsonx.MustMarshalToString(m) }

type SetCustomersStatusRequest struct {
	Ids    []dot.ID       `json:"ids"`
	Status status3.Status `json:"status"`
}

func (m *SetCustomersStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerDetailsResponse struct {
	Customer     *Customer                 `json:"customer"`
	SummaryItems []*IndependentSummaryItem `json:"summary_items"`
}

func (m *CustomerDetailsResponse) String() string { return jsonx.MustMarshalToString(m) }

type IndependentSummaryItem struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Spec  string `json:"spec"`
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

func (m *IndependentSummaryItem) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetCustomerAddresses struct {
	Phone string `json:"phone"`
}

func (m *FilterGetCustomerAddresses) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerAddressesRequest struct {
	CustomerId dot.ID                      `json:"customer_id"`
	Filter     *FilterGetCustomerAddresses `json:"filter"`
}

func (m *GetCustomerAddressesRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddress struct {
	Id           dot.ID            `json:"id"`
	Province     string            `json:"province"`
	ProvinceCode string            `json:"province_code"`
	District     string            `json:"district"`
	DistrictCode string            `json:"district_code"`
	Ward         string            `json:"ward"`
	WardCode     string            `json:"ward_code"`
	Address1     string            `json:"address1"`
	Address2     string            `json:"address2"`
	Country      string            `json:"country"`
	FullName     string            `json:"full_name"`
	Company      string            `json:"company"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Position     string            `json:"position"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
	CustomerID   dot.ID            `json:"customer_id"`
}

func (m *CustomerAddress) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerAddressRequest struct {
	CustomerId   dot.ID            `json:"customer_id"`
	ProvinceCode string            `json:"province_code"`
	DistrictCode string            `json:"district_code"`
	WardCode     string            `json:"ward_code"`
	Address1     string            `json:"address1"`
	Address2     string            `json:"address2"`
	Country      string            `json:"country"`
	FullName     string            `json:"full_name"`
	Company      string            `json:"company"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Position     string            `json:"position"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *CreateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerAddressRequest struct {
	Id           dot.ID            `json:"id"`
	ProvinceCode dot.NullString    `json:"province_code"`
	DistrictCode dot.NullString    `json:"district_code"`
	WardCode     dot.NullString    `json:"ward_code"`
	Address1     dot.NullString    `json:"address1"`
	Address2     dot.NullString    `json:"address2"`
	Country      dot.NullString    `json:"country"`
	FullName     dot.NullString    `json:"full_name"`
	Phone        dot.NullString    `json:"phone"`
	Email        dot.NullString    `json:"email"`
	Position     dot.NullString    `json:"position"`
	Company      dot.NullString    `json:"company"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *UpdateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddressesResponse struct {
	Addresses []*CustomerAddress `json:"addresses"`
}

func (m *CustomerAddressesResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductStatusRequest struct {
	Ids    []dot.ID       `json:"ids"`
	Status status3.Status `json:"status"`
}

func (m *UpdateProductStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductStatusResponse struct {
	Updated int `json:"updated"`
}

func (m *UpdateProductStatusResponse) String() string { return jsonx.MustMarshalToString(m) }

type PaymentTradingOrderRequest struct {
	OrderId         dot.ID                           `json:"order_id"`
	Desc            string                           `json:"desc"`
	ReturnUrl       string                           `json:"return_url"`
	Amount          int                              `json:"amount"`
	PaymentProvider payment_provider.PaymentProvider `json:"payment_provider"`
}

func (m *PaymentTradingOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type PaymentTradingOrderResponse struct {
	Url string `json:"url"`
}

func (m *PaymentTradingOrderResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantAttributesRequest struct {
	// @required
	VariantId  dot.ID                    `json:"variant_id"`
	Attributes []*catalogtypes.Attribute `json:"attributes"`
}

func (m *UpdateVariantAttributesRequest) String() string { return jsonx.MustMarshalToString(m) }

type PaymentCheckReturnDataRequest struct {
	Id                    string                           `json:"id"`
	Code                  string                           `json:"code"`
	PaymentStatus         string                           `json:"payment_status"`
	Amount                int                              `json:"amount"`
	ExternalTransactionId string                           `json:"external_transaction_id"`
	PaymentProvider       payment_provider.PaymentProvider `json:"payment_provider"`
}

func (m *PaymentCheckReturnDataRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopCategory struct {
	Id       dot.ID `json:"id"`
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
	ShopId   dot.ID `json:"shop_id"`
	Status   dot.ID `json:"status"`
}

func (m *ShopCategory) String() string { return jsonx.MustMarshalToString(m) }

type GetCollectionsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCollectionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopCollectionsResponse struct {
	Paging      *common.PageInfo  `json:"paging"`
	Collections []*ShopCollection `json:"collections"`
}

func (m *ShopCollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type AddShopProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *AddShopProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveShopProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *RemoveShopProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddCustomerToGroupRequest struct {
	CustomerIds []dot.ID `json:"customer_ids"`
	GroupId     dot.ID   `json:"group_id"`
}

func (m *AddCustomerToGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveCustomerOutOfGroupRequest struct {
	CustomerIds []dot.ID `json:"customer_ids"`
	GroupId     dot.ID   `json:"group_id"`
}

func (m *RemoveCustomerOutOfGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type SupplierLiability struct {
	TotalPurchaseOrders int `json:"total_purchase_orders"`
	TotalAmount         int `json:"total_amount"`
	PaidAmount          int `json:"paid_amount"`
	Liability           int `json:"liability"`
}

func (m *SupplierLiability) String() string { return jsonx.MustMarshalToString(m) }

type Supplier struct {
	Id                dot.ID             `json:"id"`
	ShopId            dot.ID             `json:"shop_id"`
	FullName          string             `json:"full_name"`
	Note              string             `json:"note"`
	Phone             string             `json:"phone"`
	Email             string             `json:"email"`
	CompanyName       string             `json:"company_name"`
	TaxNumber         string             `json:"tax_number"`
	HeadquaterAddress string             `json:"headquater_address"`
	Code              string             `json:"code"`
	Status            status3.Status     `json:"status"`
	CreatedAt         dot.Time           `json:"created_at"`
	UpdatedAt         dot.Time           `json:"updated_at"`
	Liability         *SupplierLiability `json:"liability"`
}

func (m *Supplier) String() string { return jsonx.MustMarshalToString(m) }

type CreateSupplierRequest struct {
	FullName          string `json:"full_name"`
	Note              string `json:"note"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	CompanyName       string `json:"company_name"`
	TaxNumber         string `json:"tax_number"`
	HeadquaterAddress string `json:"headquater_address"`
}

func (m *CreateSupplierRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateSupplierRequest struct {
	Id                dot.ID         `json:"id"`
	FullName          dot.NullString `json:"full_name"`
	Note              dot.NullString `json:"note"`
	Phone             dot.NullString `json:"phone"`
	Email             dot.NullString `json:"email"`
	CompanyName       dot.NullString `json:"company_name"`
	TaxNumber         dot.NullString `json:"tax_number"`
	HeadquaterAddress dot.NullString `json:"headquater_address"`
}

func (m *UpdateSupplierRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSuppliersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetSuppliersRequest) String() string { return jsonx.MustMarshalToString(m) }

type SuppliersResponse struct {
	Suppliers []*Supplier      `json:"suppliers"`
	Paging    *common.PageInfo `json:"paging"`
}

func (m *SuppliersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Carrier struct {
	Id        dot.ID         `json:"id"`
	ShopId    dot.ID         `json:"shop_id"`
	FullName  string         `json:"full_name"`
	Note      string         `json:"note"`
	Status    status3.Status `json:"status"`
	CreatedAt dot.Time       `json:"created_at"`
	UpdatedAt dot.Time       `json:"updated_at"`
}

func (m *Carrier) String() string { return jsonx.MustMarshalToString(m) }

type CreateCarrierRequest struct {
	FullName string `json:"full_name"`
	Note     string `json:"note"`
}

func (m *CreateCarrierRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCarrierRequest struct {
	Id       dot.ID         `json:"id"`
	FullName dot.NullString `json:"full_name"`
	Note     dot.NullString `json:"note"`
}

func (m *UpdateCarrierRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCarriersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCarriersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CarriersResponse struct {
	Carriers []*Carrier       `json:"carriers"`
	Paging   *common.PageInfo `json:"paging"`
}

func (m *CarriersResponse) String() string { return jsonx.MustMarshalToString(m) }

type ReceiptLine struct {
	RefId  dot.ID `json:"ref_id"`
	Title  string `json:"title"`
	Amount int    `json:"amount"`
}

func (m *ReceiptLine) String() string { return jsonx.MustMarshalToString(m) }

type Trader struct {
	Id       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Deleted  bool   `json:"deleted"`
}

func (m *Trader) String() string { return jsonx.MustMarshalToString(m) }

type Receipt struct {
	Id           dot.ID                   `json:"id"`
	ShopId       dot.ID                   `json:"shop_id"`
	TraderId     dot.ID                   `json:"trader_id"`
	CreatedBy    dot.ID                   `json:"created_by"`
	Mode         receipt_mode.ReceiptMode `json:"mode"`
	Code         string                   `json:"code"`
	Title        string                   `json:"title"`
	Type         receipt_type.ReceiptType `json:"type"`
	Description  string                   `json:"description"`
	Amount       int                      `json:"amount"`
	LedgerId     dot.ID                   `json:"ledger_id"`
	RefType      receipt_ref.ReceiptRef   `json:"ref_type"`
	Lines        []*ReceiptLine           `json:"lines"`
	Status       status3.Status           `json:"status"`
	CancelReason string                   `json:"cancel_reason"`
	PaidAt       dot.Time                 `json:"paid_at"`
	CreatedAt    dot.Time                 `json:"created_at"`
	UpdatedAt    dot.Time                 `json:"updated_at"`
	ConfirmedAt  dot.Time                 `json:"confirmed_at"`
	CancelledAt  dot.Time                 `json:"cancelled_at"`
	User         *etop.User               `json:"user"`
	Trader       *Trader                  `json:"trader"`
	Ledger       *Ledger                  `json:"ledger"`
	Note         string                   `json:"note"`

	// deprecated: use mode
	CreatedType receipt_mode.ReceiptMode `json:"created_type"`
}

func (m *Receipt) String() string { return jsonx.MustMarshalToString(m) }

type CreateReceiptRequest struct {
	TraderId    dot.ID                   `json:"trader_id"`
	Title       string                   `json:"title"`
	Type        receipt_type.ReceiptType `json:"type"`
	Description string                   `json:"description"`
	Amount      int                      `json:"amount"`
	LedgerId    dot.ID                   `json:"ledger_id"`
	RefType     receipt_ref.ReceiptRef   `json:"ref_type"`
	PaidAt      dot.Time                 `json:"paid_at"`
	Lines       []*ReceiptLine           `json:"lines"`
	Note        string                   `json:"note"`
}

func (m *CreateReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReceiptRequest struct {
	Id          dot.ID                     `json:"id"`
	TraderId    dot.NullID                 `json:"trader_id"`
	Title       dot.NullString             `json:"title"`
	Description dot.NullString             `json:"description"`
	Amount      dot.NullInt                `json:"amount"`
	LedgerId    dot.NullID                 `json:"ledger_id"`
	RefType     receipt_ref.NullReceiptRef `json:"ref_type"`
	PaidAt      dot.Time                   `json:"paid_at"`
	Lines       []*ReceiptLine             `json:"lines"`
	Note        dot.NullString             `json:"note"`
}

func (m *UpdateReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelReceiptRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
	// @deprecated use cancel_reason instead
	Reason string `json:"reson"`
}

func (m *CancelReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetReceiptsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetReceiptsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetReceiptsByLedgerTypeRequest struct {
	Type    ledger_type.LedgerType `json:"type"`
	Paging  *common.Paging         `json:"paging"`
	Filters []*common.Filter       `json:"filters"`
}

func (m *GetReceiptsByLedgerTypeRequest) String() string { return jsonx.MustMarshalToString(m) }

type ReceiptsResponse struct {
	TotalAmountConfirmedReceipt int              `json:"total_amount_confirmed_receipt"`
	TotalAmountConfirmedPayment int              `json:"total_amount_confirmed_payment"`
	Receipts                    []*Receipt       `json:"receipts"`
	Paging                      *common.PageInfo `json:"paging"`
}

func (m *ReceiptsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetShopCollectionsByProductIDRequest struct {
	ProductId dot.ID `json:"product_id"`
}

func (m *GetShopCollectionsByProductIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateInventoryVoucherRequest struct {
	RefId   dot.ID                                    `json:"ref_id"`
	RefType inventory_voucher_ref.InventoryVoucherRef `json:"ref_type"`
	//enum "in" or "out" only for ref_type = "order"
	Type inventory_type.InventoryVoucherType `json:"type"`
}

func (m *CreateInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVoucherLine struct {
	VariantId   dot.ID                    `json:"variant_id"`
	Code        string                    `json:"code"`
	VariantName string                    `json:"variant_name"`
	ProductId   dot.ID                    `json:"product_id"`
	ProductName string                    `json:"product_name"`
	ImageUrl    string                    `json:"image_url"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
	Price       int                       `json:"price"`
	Quantity    int                       `json:"quantity"`
}

func (m *InventoryVoucherLine) String() string { return jsonx.MustMarshalToString(m) }

type CreateInventoryVoucherResponse struct {
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *CreateInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmInventoryVoucherRequest struct {
	Id dot.ID `json:"id"`
}

func (m *ConfirmInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmInventoryVoucherResponse struct {
	InventoryVoucher *InventoryVoucher `json:"inventory_voucher"`
}

func (m *ConfirmInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type CancelInventoryVoucherRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelInventoryVoucherResponse struct {
	Inventory *InventoryVoucher `json:"inventory"`
}

func (m *CancelInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateInventoryVoucherRequest struct {
	Id          dot.ID                 `json:"id"`
	TraderId    dot.NullID             `json:"trader_id"`
	Lines       []InventoryVoucherLine `json:"lines"`
	Note        dot.NullString         `json:"note"`
	Type        string                 `json:"type"`
	Title       dot.NullString         `json:"title"`
	TotalAmount int                    `json:"total_amount"`
}

func (m *UpdateInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateInventoryVoucherResponse struct {
	InventoryVoucher *InventoryVoucher `json:"inventory_voucher"`
}

func (m *UpdateInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type AdjustInventoryQuantityRequest struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
	Note              string              `json:"note"`
}

func (m *AdjustInventoryQuantityRequest) String() string { return jsonx.MustMarshalToString(m) }

type AdjustInventoryQuantityResponse struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *AdjustInventoryQuantityResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsRequest struct {
	Paging common.Paging `json:"paging"`
}

func (m *GetInventoryVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsResponse struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
}

func (m *GetInventoryVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsByVariantIDsRequest struct {
	VariantIds []dot.ID `json:"variant_ids"`
}

func (m *GetInventoryVariantsByVariantIDsRequest) Reset() {
	*m = GetInventoryVariantsByVariantIDsRequest{}
}
func (m *GetInventoryVariantsByVariantIDsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type InventoryVariant struct {
	ShopId         dot.ID   `json:"shop_id"`
	VariantId      dot.ID   `json:"variant_id"`
	QuantityOnHand int      `json:"quantity_on_hand"`
	QuantityPicked int      `json:"quantity_picked"`
	CostPrice      int      `json:"cost_price"`
	Quantity       int      `json:"quantity"`
	CreatedAt      dot.Time `json:"created_at"`
	UpdatedAt      dot.Time `json:"updated_at"`
}

func (m *InventoryVariant) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVariantQuantity struct {
	QuantityOnHand int `json:"quantity_on_hand"`
	QuantityPicked int `json:"quantity_picked"`
	Quantity       int `json:"quantity"`
}

func (m *InventoryVariantQuantity) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVoucher struct {
	TotalAmount  int                     `json:"total_amount"`
	CreatedBy    dot.ID                  `json:"created_by"`
	UpdatedBy    dot.ID                  `json:"updated_by"`
	Lines        []*InventoryVoucherLine `json:"lines"`
	TraderId     dot.ID                  `json:"trader_id"`
	Note         string                  `json:"note"`
	Type         string                  `json:"type"`
	Id           dot.ID                  `json:"id"`
	ShopId       dot.ID                  `json:"shop_id"`
	Title        string                  `json:"title"`
	RefId        dot.ID                  `json:"ref_id"`
	RefType      string                  `json:"ref_type"`
	RefName      string                  `json:"ref_name"`
	RefAction    ref_action.RefAction    `json:"ref_action"`
	RefCode      string                  `json:"ref_code"`
	Code         string                  `json:"code"`
	CreatedAt    dot.Time                `json:"created_at"`
	UpdatedAt    dot.Time                `json:"updated_at"`
	CancelledAt  dot.Time                `json:"cancelled_at"`
	ConfirmedAt  dot.Time                `json:"confirmed_at"`
	CancelReason string                  `json:"cancel_reason"`
	Trader       *Trader                 `json:"trader"`
	Status       status3.Status          `json:"status"`
}

func (m *InventoryVoucher) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetInventoryVouchersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersByIDsRequest struct {
	Ids []dot.ID `json:"ids"`
}

func (m *GetInventoryVouchersByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersResponse struct {
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *GetInventoryVouchersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroup struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *CustomerGroup) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerGroupRequest struct {
	Name string `json:"name"`
}

func (m *CreateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerGroupRequest struct {
	GroupId dot.ID `json:"group_id"`
	Name    string `json:"name"`
}

func (m *UpdateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerGroupsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCustomerGroupsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupsResponse struct {
	CustomerGroups []*CustomerGroup `json:"customer_groups"`
	Paging         *common.PageInfo `json:"paging"`
}

func (m *CustomerGroupsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetOrdersByReceiptIDRequest struct {
	ReceiptId dot.ID `json:"receipt_id"`
}

func (m *GetOrdersByReceiptIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantRequest struct {
	VariantId dot.ID `json:"variant_id"`
}

func (m *GetInventoryVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateBrandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *CreateBrandRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateBrandRequest struct {
	Id          dot.ID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *UpdateBrandRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteBrandResponse struct {
	Count int `json:"count"`
}

func (m *DeleteBrandResponse) String() string { return jsonx.MustMarshalToString(m) }

type Brand struct {
	ShopId      dot.ID   `json:"shop_id"`
	Id          dot.ID   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
}

func (m *Brand) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsByIDsResponse struct {
	Brands []*Brand `json:"brands"`
}

func (m *GetBrandsByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsRequest struct {
	Paging common.Paging `json:"paging"`
}

func (m *GetBrandsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsResponse struct {
	Brands []*Brand         `json:"brands"`
	Paging *common.PageInfo `json:"paging"`
}

func (m *GetBrandsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersByReferenceRequest struct {
	RefId dot.ID `json:"ref_id"`
	// enum ('order', 'purchase_order', 'return', 'purchase_order')
	RefType inventory_voucher_ref.InventoryVoucherRef `json:"ref_type"`
}

func (m *GetInventoryVouchersByReferenceRequest) Reset() {
	*m = GetInventoryVouchersByReferenceRequest{}
}
func (m *GetInventoryVouchersByReferenceRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersByReferenceResponse struct {
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *GetInventoryVouchersByReferenceResponse) Reset() {
	*m = GetInventoryVouchersByReferenceResponse{}
}
func (m *GetInventoryVouchersByReferenceResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateOrderShippingInfoRequest struct {
	OrderId         dot.ID               `json:"order_id"`
	Shipping        *types.OrderShipping `json:"shipping"`
	ShippingAddress *types.OrderAddress  `json:"shipping_address"`
}

func (m *UpdateOrderShippingInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesByIDsResponse struct {
	Stocktakes []*Stocktake `json:"stocktakes"`
}

func (m *GetStocktakesByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateStocktakeRequest struct {
	TotalQuantity int    `json:"total_quantity"`
	Note          string `json:"note"`
	//  length more than one
	Lines []*StocktakeLine             `json:"lines"`
	Type  stocktake_type.StocktakeType `json:"type"`
}

func (m *CreateStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateStocktakeRequest struct {
	Id            dot.ID `json:"id"`
	TotalQuantity int    `json:"total_quantity"`
	Note          string `json:"note"`
	//  length more than one
	Lines []*StocktakeLine `json:"lines"`
}

func (m *UpdateStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type Stocktake struct {
	Id            dot.ID           `json:"id"`
	ShopId        dot.ID           `json:"shop_id"`
	TotalQuantity int              `json:"total_quantity"`
	Note          string           `json:"note"`
	Code          string           `json:"code"`
	CancelReason  string           `json:"cancel_reason"`
	CreatedBy     dot.ID           `json:"created_by"`
	UpdatedBy     dot.ID           `json:"updated_by"`
	CreatedAt     dot.Time         `json:"created_at"`
	UpdatedAt     dot.Time         `json:"updated_at"`
	ConfirmedAt   dot.Time         `json:"confirmed_at"`
	CancelledAt   dot.Time         `json:"cancelled_at"`
	Status        status3.Status   `json:"status"`
	Lines         []*StocktakeLine `json:"lines"`
	Type          string           `json:"type"`
}

func (m *Stocktake) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetStocktakesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesResponse struct {
	Stocktakes []*Stocktake     `json:"stocktakes"`
	Paging     *common.PageInfo `json:"paging"`
}

func (m *GetStocktakesResponse) String() string { return jsonx.MustMarshalToString(m) }

type StocktakeLine struct {
	ProductId   dot.ID                    `json:"product_id"`
	ProductName string                    `json:"product_name"`
	VariantName string                    `json:"variant_name"`
	VariantId   dot.ID                    `json:"variant_id"`
	OldQuantity int                       `json:"old_quantity"`
	NewQuantity int                       `json:"new_quantity"`
	Code        string                    `json:"code"`
	ImageUrl    string                    `json:"image_url"`
	CostPrice   int                       `json:"cost_price"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
}

func (m *StocktakeLine) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmStocktakeRequest struct {
	Id                   dot.ID                              `json:"id"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *ConfirmStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantsBySupplierIDRequest struct {
	SupplierId dot.ID `json:"supplier_id"`
}

func (m *GetVariantsBySupplierIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSuppliersByVariantIDRequest struct {
	VariantId dot.ID `json:"variant_id"`
}

func (m *GetSuppliersByVariantIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelStocktakeRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetRefundsByIDsResponse struct {
	Refund []*Refund `json:"refunds"`
}

func (m *GetRefundsByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetRefundsResponse struct {
	Refunds []*Refund        `json:"refunds"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *GetRefundsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetRefundsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetRefundsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelRefundRequest struct {
	ID                   dot.ID                              `json:"id"`
	CancelReason         string                              `json:"cancel_reason"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *CancelRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmRefundRequest struct {
	ID                   dot.ID                              `json:"id"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *ConfirmRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateRefundRequest struct {
	Lines           []*RefundLine           `json:"lines"`
	Note            dot.NullString          `json:"note"`
	ID              dot.ID                  `json:"id"`
	TotalAjustment  dot.NullInt             `json:"total_adjustment"`
	AdjustmentLines []*types.AdjustmentLine `json:"adjustment_lines"`
	TotalAmount     dot.NullInt             `json:"total_amount"`
	BasketValue     dot.NullInt             `json:"basket_value"`
}

func (m *UpdateRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateRefundRequest struct {
	Lines           []*RefundLine           `json:"lines"`
	OrderID         dot.ID                  `json:"order_id"`
	Note            string                  `json:"note"`
	TotalAjustment  int                     `json:"total_adjustment"`
	AdjustmentLines []*types.AdjustmentLine `json:"adjustment_lines"`
	TotalAmount     int                     `json:"total_amount"`
	BasketValue     int                     `json:"basket_value"`
}

func (m *CreateRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type Refund struct {
	ID               dot.ID                  `json:"id"`
	ShopID           dot.ID                  `json:"shop_id"`
	OrderID          dot.ID                  `json:"order_id"`
	Note             string                  `json:"note"`
	Code             string                  `json:"code"`
	AdjustmentLines  []*types.AdjustmentLine `json:"adjustment_lines"`
	TotalAdjustment  int                     `json:"total_adjustment"`
	Lines            []*RefundLine           `json:"lines"`
	CreatedAt        dot.Time                `json:"created_at"`
	UpdatedAt        dot.Time                `json:"updated_at"`
	CancelledAt      dot.Time                `json:"cancelled_at"`
	ConfirmedAt      dot.Time                `json:"confirmed_at"`
	CreatedBy        dot.ID                  `json:"created_by"`
	UpdatedBy        dot.ID                  `json:"updated_by"`
	CancelReason     string                  `json:"cancel_reason"`
	Customer         *types.OrderCustomer    `json:"customer"`
	CustomerID       dot.ID                  `json:"customer_id"`
	Status           status3.Status          `json:"status"`
	TotalAmount      int                     `json:"total_amount"`
	BasketValue      int                     `json:"basket_value"`
	PaidAmount       int                     `json:"paid_amount"`
	InventoryVoucher *InventoryVoucher       `json:"inventory_voucher"`
}

func (m *Refund) String() string { return jsonx.MustMarshalToString(m) }

type RefundLine struct {
	VariantID   dot.ID                    `json:"variant_id"`
	ProductID   dot.ID                    `json:"product_id"`
	Quantity    int                       `json:"quantity"`
	Code        string                    `json:"code"`
	ImageURL    string                    `json:"image_url"`
	Name        string                    `json:"name"`
	RetailPrice int                       `json:"retail_price"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
	Adjustment  int                       `json:"adjustment"`
}

func (m *RefundLine) String() string { return jsonx.MustMarshalToString(m) }

type GetPurchaseRefundsByIDsResponse struct {
	PurchaseRefund []*PurchaseRefund `json:"purchase_refunds"`
}

func (m *GetPurchaseRefundsByIDsResponse) Reset()         { *m = GetPurchaseRefundsByIDsResponse{} }
func (m *GetPurchaseRefundsByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetPurchaseRefundsResponse struct {
	PurchaseRefunds []*PurchaseRefund `json:"purchase_refunds"`
	Paging          *common.PageInfo  `json:"paging"`
}

func (m *GetPurchaseRefundsResponse) Reset()         { *m = GetPurchaseRefundsResponse{} }
func (m *GetPurchaseRefundsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetPurchaseRefundsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetPurchaseRefundsRequest) Reset()         { *m = GetPurchaseRefundsRequest{} }
func (m *GetPurchaseRefundsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelPurchaseRefundRequest struct {
	ID                   dot.ID                              `json:"id"`
	CancelReason         string                              `json:"cancel_reason"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *CancelPurchaseRefundRequest) Reset()         { *m = CancelPurchaseRefundRequest{} }
func (m *CancelPurchaseRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmPurchaseRefundRequest struct {
	ID                   dot.ID                              `json:"id"`
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher `json:"auto_inventory_voucher"`
}

func (m *ConfirmPurchaseRefundRequest) Reset()         { *m = ConfirmPurchaseRefundRequest{} }
func (m *ConfirmPurchaseRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePurchaseRefundRequest struct {
	Lines           []*PurchaseRefundLine   `json:"lines"`
	Note            dot.NullString          `json:"note"`
	ID              dot.ID                  `json:"id"`
	TotalAdjustment dot.NullInt             `json:"total_adjustment"`
	AdjustmentLines []*types.AdjustmentLine `json:"adjustment_lines"`
	BasketValue     dot.NullInt             `json:"basket_value"`
	TotalAmount     dot.NullInt             `json:"total_amount"`
}

func (m *UpdatePurchaseRefundRequest) Reset()         { *m = UpdatePurchaseRefundRequest{} }
func (m *UpdatePurchaseRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreatePurchaseRefundRequest struct {
	Lines           []*PurchaseRefundLine   `json:"lines"`
	PurchaseOrderID dot.ID                  `json:"purchase_order_id"`
	Note            string                  `json:"note"`
	Discount        int                     `json:"discount"`
	TotalAdjustment int                     `json:"total_adjustment"`
	AdjustmentLines []*types.AdjustmentLine `json:"adjustment_lines"`
	BasketValue     int                     `json:"basket_value"`
	TotalAmount     int                     `json:"total_amount"`
}

func (m *CreatePurchaseRefundRequest) Reset()         { *m = CreatePurchaseRefundRequest{} }
func (m *CreatePurchaseRefundRequest) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseRefund struct {
	ID               dot.ID                  `json:"id"`
	ShopID           dot.ID                  `json:"shop_id"`
	PurchaseOrderID  dot.ID                  `json:"purchase_order_id"`
	Note             string                  `json:"note"`
	Code             string                  `json:"code"`
	TotalAdjustment  int                     `json:"total_adjustment"`
	AdjustmentLines  []*types.AdjustmentLine `json:"adjustment_lines"`
	Lines            []*PurchaseRefundLine   `json:"lines"`
	CreatedAt        dot.Time                `json:"created_at"`
	UpdatedAt        dot.Time                `json:"updated_at"`
	CancelledAt      dot.Time                `json:"cancelled_at"`
	ConfirmedAt      dot.Time                `json:"confirmed_at"`
	CreatedBy        dot.ID                  `json:"created_by"`
	UpdatedBy        dot.ID                  `json:"updated_by"`
	CancelReason     string                  `json:"cancel_reason"`
	Supplier         *PurchaseOrderSupplier  `json:"supplier"`
	SupplierID       dot.ID                  `json:"supplier_id"`
	Status           status3.Status          `json:"status"`
	TotalAmount      int                     `json:"total_amount"`
	BasketValue      int                     `json:"basket_value"`
	InventoryVoucher *InventoryVoucher       `json:"inventory"`
}

func (m *PurchaseRefund) Reset()         { *m = PurchaseRefund{} }
func (m *PurchaseRefund) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseRefundLine struct {
	VariantID    dot.ID                    `json:"variant_id"`
	ProductID    dot.ID                    `json:"product_id"`
	Quantity     int                       `json:"quantity"`
	Code         string                    `json:"code"`
	ImageURL     string                    `json:"image_url"`
	Name         string                    `json:"name"`
	PaymentPrice int                       `json:"payment_price"`
	Attributes   []*catalogtypes.Attribute `json:"attributes"`
	Adjustment   int                       `json:"adjustment"`
}

func (m *PurchaseRefundLine) Reset()         { *m = PurchaseRefundLine{} }
func (m *PurchaseRefundLine) String() string { return jsonx.MustMarshalToString(m) }

type UpdateInventoryVariantCostPriceResponse struct {
	InventoryVariant *InventoryVariant `json:"inventory_variant"`
}

func (m *UpdateInventoryVariantCostPriceResponse) Reset() {
	*m = UpdateInventoryVariantCostPriceResponse{}
}
func (m *UpdateInventoryVariantCostPriceResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateInventoryVariantCostPriceRequest struct {
	VariantId dot.ID `json:"variant_id"`
	CostPrice int    `json:"cost_price"`
}

func (m *UpdateInventoryVariantCostPriceRequest) Reset() {
	*m = UpdateInventoryVariantCostPriceRequest{}
}
func (m *UpdateInventoryVariantCostPriceRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummaryTableJSON struct {
	Label   string          `json:"label"`
	Tags    []string        `json:"tags"`
	Columns []SummaryColRow `json:"columns"`
	Rows    []SummaryColRow `json:"rows"`
	Data    [][]SummaryItem `json:"data"`
}

// MarshalJSON implements JSONMarshaler
func (m *SummaryTable) MarshalJSON() ([]byte, error) {
	ncol := len(m.Columns)
	data := make([][]SummaryItem, len(m.Rows))
	for r := range m.Rows {
		data[r] = m.Data[r*ncol : (r+1)*ncol]
	}
	res := SummaryTableJSON{
		Label:   m.Label,
		Tags:    m.Tags,
		Columns: m.Columns,
		Rows:    m.Rows,
		Data:    data,
	}
	return jsonx.Marshal(res)
}

// UnmarshalJSON implements JSONUnmarshaler
func (m *SummaryTable) UnmarshalJSON(data []byte) error {
	var tmp SummaryTableJSON
	if err := jsonx.Unmarshal(data, &tmp); err != nil {
		return err
	}
	ncol := len(tmp.Columns)
	mdata := make([]SummaryItem, len(tmp.Rows)*ncol)
	for r := range tmp.Rows {
		copy(mdata[r*ncol:], tmp.Data[r])
	}
	*m = SummaryTable{
		Label:   tmp.Label,
		Tags:    tmp.Tags,
		Columns: tmp.Columns,
		Rows:    tmp.Rows,
		Data:    mdata,
	}
	return nil
}

func (m *ImportProductsResponse) HasErrors() []*common.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}

type GetVariantRequest struct {
	Code string `json:"code"`
	ID   dot.ID `json:"id"`
}

func (m *GetVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateFulfillmentsRequest struct {
	OrderID             dot.ID                                        `json:"order_id"`
	ShippingType        ordertypes.ShippingType                       `json:"shipping_type"`
	ShippingServiceCode string                                        `json:"shipping_service_code"`
	ShippingServiceFee  int                                           `json:"shipping_service_fee"`
	ShippingServiceName string                                        `json:"shipping_service_name"`
	ShippingNote        string                                        `json:"shipping_note"`
	PickupAddress       *types.OrderAddress                           `json:"pickup_address"`
	ReturnAddress       *types.OrderAddress                           `json:"return_address"`
	ShippingAddress     *types.OrderAddress                           `json:"shipping_address"`
	TryOn               try_on.TryOnCode                              `json:"try_on"`
	ShippingPaymentType shipping_payment_type.NullShippingPaymentType `json:"shipping_payment_type"`
	ChargeableWeight    int                                           `json:"chargeable_weight"`
	GrossWeight         int                                           `json:"gross_weight"`
	Height              int                                           `json:"height"`
	Width               int                                           `json:"width"`
	Length              int                                           `json:"length"`
	CODAmount           int                                           `json:"cod_amount"`
	IncludeInsurance    bool                                          `json:"include_insurance"`
	InsuranceValue      dot.NullInt                                   `json:"insurance_value"`
	Coupon              string                                        `json:"coupon"`

	ConnectionID  dot.ID `json:"connection_id"`
	ShopCarrierID dot.ID `json:"shop_carrier_id"`
}

func (m *CreateFulfillmentsRequest) Reset() {
	*m = CreateFulfillmentsRequest{}
}
func (m *CreateFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateFulfillmentsResponse struct {
	Fulfillments []*types.Fulfillment `json:"fulfillment"`
	Errors       []*common.Error      `json:"errors"`
}

func (m *CreateFulfillmentsResponse) Reset() {
	*m = CreateFulfillmentsResponse{}
}
func (m *CreateFulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentCODRequest struct {
	FulfillmentID dot.ID `json:"fulfillment_id"`
	CODAmount     int    `json:"cod_amount"`
}

func (m *UpdateFulfillmentCODRequest) Reset() {
	*m = UpdateFulfillmentCODRequest{}
}

func (m *UpdateFulfillmentCODRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelFulfillmentRequest struct {
	FulfillmentID dot.ID `json:"fulfillment_id"`
	CancelReason  string `json:"cancel_reason"`
}

func (m *CancelFulfillmentRequest) Reset() {
	*m = CancelFulfillmentRequest{}
}
func (m *CancelFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentInfoRequest struct {
	FulfillmentID       dot.ID                                        `json:"fulfillment_id"`
	PickupAddress       *types.OrderAddress                           `json:"pickup_address"`
	ShippingAddress     *types.OrderAddress                           `json:"shipping_address"`
	IncludeInsurance    dot.NullBool                                  `json:"include_insurance"`
	InsuranceValue      dot.NullInt                                   `json:"insurance_value"`
	GrossWeight         dot.NullInt                                   `json:"gross_weight"`
	TryOn               try_on.TryOnCode                              `json:"try_on"`
	ShippingNote        dot.NullString                                `json:"shipping_note"`
	ShippingPaymentType shipping_payment_type.NullShippingPaymentType `json:"shipping_payment_type"`
}

func (m *UpdateFulfillmentInfoRequest) Reset() { *m = UpdateFulfillmentInfoRequest{} }

func (m *UpdateFulfillmentInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateFulfillmentFromImportArgs struct {
	EdCode              string              `json:"ed_code"`
	PickupAddress       *types.OrderAddress `json:"pickup_address"`
	ShippingAddress     *types.OrderAddress `json:"shipping_address"`
	ProductDescription  string              `json:"product_description"`
	TotalWeight         int                 `json:"total_weight"` // gram
	BasketValue         int                 `json:"basket_value"`
	IncludeInsurance    bool                `json:"include_insurance"`
	CODAmount           int                 `json:"cod_amount"`
	ShippingNote        string              `json:"shipping_note"`
	TryOn               try_on.TryOnCode    `json:"try_on"`
	ConnectionID        dot.ID              `json:"connection_id"`
	ShippingServiceCode string              `json:"shipping_service_code"`
	ShippingServiceFee  int                 `json:"shipping_service_fee"`
	ShippingServiceName string              `json:"shipping_service_name"`
}

func (m *CreateFulfillmentFromImportArgs) String() string { return jsonx.MustMarshalToString(m) }

func (m *CreateFulfillmentFromImportArgs) Reset() {
	*m = CreateFulfillmentFromImportArgs{}
}

type CreateFulfillmentsFromImportRequest struct {
	Fulfillments []*CreateFulfillmentFromImportArgs `json:"fulfillments"`
}

func (m *CreateFulfillmentsFromImportRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *CreateFulfillmentsFromImportRequest) Reset() {
	*m = CreateFulfillmentsFromImportRequest{}
}

type CreateFulfillmentsFromImportResponse struct {
	Fulfillments []*types.Fulfillment `json:"fulfillments"`
	Errors       []*common.Error      `json:"errors"`
}

func (m *CreateFulfillmentsFromImportResponse) String() string { return jsonx.MustMarshalToString(m) }

func (m *CreateFulfillmentsFromImportResponse) Reset() {
	*m = CreateFulfillmentsFromImportResponse{}
}

type GetAccountShipnowRequest struct {
	Identity     string `json:"identity"`
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *GetAccountShipnowRequest) String() string { return jsonx.MustMarshalToString(m) }

type Contact struct {
	ID        dot.ID   `json:"id"`
	ShopID    dot.ID   `json:"shop_id"`
	FullName  string   `json:"full_name"`
	Phone     string   `json:"phone"`
	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *Contact) String() string { return jsonx.MustMarshalToString(m) }

type CreateContactRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

func (m *CreateContactRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateContactRequest struct {
	ID       dot.ID         `json:"id"`
	FullName dot.NullString `json:"full_name"`
	Phone    dot.NullString `json:"phone"`
}

func (m *UpdateContactRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetContactRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetContactRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetContactsRequest struct {
	Paging *common.CursorPaging `json:"paging"`
	Filter *FilterGetContacts   `json:"filter"`
}

func (m *GetContactsRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetContacts struct {
	IDs   []dot.ID `json:"ids"`
	Phone string   `json:"phone"`
	Name  string   `json:"name"`
}

func (m *FilterGetContacts) String() string { return jsonx.MustMarshalToString(m) }

type GetContactsResponse struct {
	Contacts []*Contact             `json:"contacts"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetContactsResponse) String() string { return jsonx.MustMarshalToString(m) }

type DeleteContactResponse struct {
	Count int `json:"count"`
}

func (m *DeleteContactResponse) String() string { return jsonx.MustMarshalToString(m) }

type DeleteContactRequest struct {
	ID dot.ID `json:"id"`
}

func (m *DeleteContactRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopSetting struct {
	ShopID                     dot.ID                                        `json:"shop_id"`
	PaymentTypeID              shipping_payment_type.NullShippingPaymentType `json:"payment_type_id"`
	ReturnAddressID            dot.ID                                        `json:"return_address_id"`
	ReturnAddress              *etop.Address                                 `json:"return_address"`
	TryOn                      try_on.NullTryOnCode                          `json:"try_on"`
	ShippingNote               string                                        `json:"shipping_note"`
	Weight                     int                                           `json:"weight"`
	HideAllComments            dot.NullBool                                  `json:"hide_all_comments"`
	CreatedAt                  dot.Time                                      `json:"created_at"`
	UpdatedAt                  dot.Time                                      `json:"updated_at"`
	AllowConnectDirectShipment bool                                          `json:"allow_connect_direct_shipment"`
}

func (m *ShopSetting) String() string { return jsonx.MustMarshalToString(m) }

type UpdateSettingRequest struct {
	PaymentTypeID   shipping_payment_type.NullShippingPaymentType `json:"payment_type_id"`
	ReturnAddress   *etop.Address                                 `json:"return_address"`
	TryOn           try_on.NullTryOnCode                          `json:"try_on"`
	ShippingNote    dot.NullString                                `json:"shipping_note"`
	Weight          dot.NullInt                                   `json:"weight"`
	HideAllComments dot.NullBool                                  `json:"hide_all_comments"`
}

func (m *UpdateSettingRequest) String() string { return jsonx.MustMarshalToString(m) }

type JiraCustomField struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (m *JiraCustomField) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomFieldsResponse struct {
	CustomFields []*JiraCustomField `json:"custom_fields"`
}

func (m *GetCustomFieldsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateJiraIssueRequest struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	// CustomField l???y t??? api: GetJiraCustomFields
	CustomFields []CreateCustomField `json:"custom_fields"`
}

func (m *CreateJiraIssueRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *CreateCustomField) String() string { return jsonx.MustMarshalToString(m) }

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (m *CreateIssueResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateAccountUserRequest struct {
	FullName     string                    `json:"full_name"`
	Phone        string                    `json:"phone"`
	Roles        []shop_user_role.UserRole `json:"roles"`
	Password     string                    `json:"password"`
	DepartmentID dot.ID                    `json:"department_id"`
}

func (m *CreateAccountUserRequest) String() string { return jsonx.MustMarshalToString(m) }

func (r *CreateAccountUserRequest) GetAccountUserRoles() []string {
	var res = make([]string, 0, len(r.Roles))
	for _, role := range r.Roles {
		res = append(res, role.String())
	}
	return res
}

func (r *CreateAccountUserRequest) Validate() error {
	if r.FullName == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui l??ng ??i???n h??? t??n")
	}
	if r.Password == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui l??ng ??i???n password")
	}
	if len(r.Roles) == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui l??ng ch???n roles")
	}
	return nil
}

type GetAccountUsersRequest struct {
	Paging *common.CursorPaging          `json:"paging"`
	Filter *FilterGetAccountUsersRequest `json:"filter"`
}

func (m *GetAccountUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetAccountUsersRequest struct {
	Name            filter.FullTextSearch     `json:"name"`
	Phone           filter.FullTextSearch     `json:"phone"`
	ExtensionNumber filter.FullTextSearch     `json:"extension_number"`
	Roles           []shop_user_role.UserRole `json:"roles"`
	ExactRoles      []shop_user_role.UserRole `json:"exact_roles"`
	UserIDs         []dot.ID                  `json:"user_ids"`
	HasExtension    dot.NullBool              `json:"has_extension"`
	HasDepartment   dot.NullBool              `json:"has_department"`
	DepartmentID    dot.ID                    `json:"department_id"`
}

func (m *FilterGetAccountUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetAccountUsersResponse struct {
	AccountUsers []*shoptypes.AccountUserExtended `json:"account_users"`
	Paging       *common.CursorPageInfo           `json:"paging"`
}

func (m *GetAccountUsersResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAccountUserRequest struct {
	UserID   dot.ID                    `json:"user_id"`
	FullName string                    `json:"full_name"`
	Roles    []shop_user_role.UserRole `json:"roles"`
	// C?? th??? ?????i m???t kh???u cho nh??n vi??n, n???u nh??n vi??n ???? ch??a ?????i m???t kh???u l???n n??o.
	// T???c m???t kh???u v???n ??ang ??? d???ng kh???i t???o khi v???a m???i t???o t??i kho???n.
	Password     string `json:"password"`
	DepartmentID dot.ID `json:"department_id"`
}

func (m *UpdateAccountUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteAccountUserRequest struct {
	UserID dot.ID `json:"user_id"`
}

func (m *DeleteAccountUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateDepartmentRequest struct {
	// @required
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *CreateDepartmentRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *CreateDepartmentRequest) Validate() error {
	if m.Name == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "missing name")
	}
	return nil
}

type GetDepartmentsRequest struct {
	Filter *GetDepartmentsFilter `json:"filter"`
	Paging *common.CursorPaging  `json:"paging"`
}

func (m *GetDepartmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetDepartmentsFilter struct {
	Name filter.FullTextSearch `json:"name"`
}

func (m *GetDepartmentsFilter) String() string { return jsonx.MustMarshalToString(m) }

type GetDepartmentsResponse struct {
	Departments []*shoptypes.Department `json:"departments"`
	Paging      *common.CursorPageInfo  `json:"paging"`
}

func (m *GetDepartmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateDepartmentRequest struct {
	ID          dot.ID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *UpdateDepartmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveUserOutOfDepartmentRequest struct {
	UserID       dot.ID `json:"user_id"`
	DepartmentID dot.ID `json:"department_id"`
}

func (m *RemoveUserOutOfDepartmentRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *RemoveUserOutOfDepartmentRequest) Validate() error {
	if m.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "missing user id")
	}
	if m.DepartmentID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "missing department id")
	}
	return nil
}
