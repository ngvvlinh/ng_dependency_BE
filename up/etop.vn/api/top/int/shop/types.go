package shop

import (
	etop "etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	spreadsheet "etop.vn/api/top/int/types/spreadsheet"
	common "etop.vn/api/top/types/common"
	ghn_note_code "etop.vn/api/top/types/etc/ghn_note_code"
	payment_provider "etop.vn/api/top/types/etc/payment_provider"
	product_type "etop.vn/api/top/types/etc/product_type"
	shipping "etop.vn/api/top/types/etc/shipping"
	status3 "etop.vn/api/top/types/etc/status3"
	status4 "etop.vn/api/top/types/etc/status4"
	try_on "etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type PurchaseOrder struct {
	Id               dot.ID                 `json:"id"`
	ShopId           dot.ID                 `json:"shop_id"`
	SupplierId       dot.ID                 `json:"supplier_id"`
	Supplier         *PurchaseOrderSupplier `json:"supplier"`
	BasketValue      int                    `json:"basket_value"`
	TotalDiscount    int                    `json:"total_discount"`
	TotalAmount      int                    `json:"total_amount"`
	Code             string                 `json:"code"`
	Note             string                 `json:"note"`
	Status           status3.Status         `json:"status"`
	Lines            []*PurchaseOrderLine   `json:"lines"`
	CreatedBy        dot.ID                 `json:"created_by"`
	CancelledReason  string                 `json:"cancelled_reason"`
	ConfirmedAt      dot.Time               `json:"confirmed_at"`
	CancelledAt      dot.Time               `json:"cancelled_at"`
	CreatedAt        dot.Time               `json:"created_at"`
	UpdatedAt        dot.Time               `json:"updated_at"`
	DeletedAt        dot.Time               `json:"deleted_at"`
	InventoryVoucher *InventoryVoucher      `json:"inventory_voucher"`
	PaidAmount       int                    `json:"paid_amount"`
}

func (m *PurchaseOrder) Reset()         { *m = PurchaseOrder{} }
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

func (m *PurchaseOrderSupplier) Reset()         { *m = PurchaseOrderSupplier{} }
func (m *PurchaseOrderSupplier) String() string { return jsonx.MustMarshalToString(m) }

type GetPurchaseOrdersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetPurchaseOrdersRequest) Reset()         { *m = GetPurchaseOrdersRequest{} }
func (m *GetPurchaseOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrdersResponse struct {
	PurchaseOrders []*PurchaseOrder `json:"purchase_orders"`
	Paging         *common.PageInfo `json:"paging"`
}

func (m *PurchaseOrdersResponse) Reset()         { *m = PurchaseOrdersResponse{} }
func (m *PurchaseOrdersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreatePurchaseOrderRequest struct {
	SupplierId    dot.ID               `json:"supplier_id"`
	BasketValue   int                  `json:"basket_value"`
	TotalDiscount int                  `json:"total_discount"`
	TotalAmount   int                  `json:"total_amount"`
	Note          string               `json:"note"`
	Lines         []*PurchaseOrderLine `json:"lines"`
}

func (m *CreatePurchaseOrderRequest) Reset()         { *m = CreatePurchaseOrderRequest{} }
func (m *CreatePurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePurchaseOrderRequest struct {
	Id            dot.ID               `json:"id"`
	SupplierId    dot.NullID           `json:"supplier_id"`
	BasketValue   dot.NullInt          `json:"basket_value"`
	TotalDiscount dot.NullInt          `json:"total_discount"`
	TotalAmount   dot.NullInt          `json:"total_amount"`
	Note          dot.NullString       `json:"note"`
	Lines         []*PurchaseOrderLine `json:"lines"`
}

func (m *UpdatePurchaseOrderRequest) Reset()         { *m = UpdatePurchaseOrderRequest{} }
func (m *UpdatePurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type PurchaseOrderLine struct {
	VariantId    dot.ID       `json:"variant_id"`
	Quantity     int          `json:"quantity"`
	PaymentPrice int          `json:"payment_price"`
	ProductId    dot.ID       `json:"product_id"`
	ProductName  string       `json:"product_name"`
	ImageUrl     string       `json:"image_url"`
	Code         string       `json:"code"`
	Attributes   []*Attribute `json:"attributes"`
}

func (m *PurchaseOrderLine) Reset()         { *m = PurchaseOrderLine{} }
func (m *PurchaseOrderLine) String() string { return jsonx.MustMarshalToString(m) }

type CancelPurchaseOrderRequest struct {
	Id     dot.ID `json:"id"`
	Reason string `json:"reason"`
}

func (m *CancelPurchaseOrderRequest) Reset()         { *m = CancelPurchaseOrderRequest{} }
func (m *CancelPurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmPurchaseOrderRequest struct {
	Id dot.ID `json:"id"`
	// enum create, confirm
	AutoInventoryVoucher string `json:"auto_inventory_voucher"`
}

func (m *ConfirmPurchaseOrderRequest) Reset()         { *m = ConfirmPurchaseOrderRequest{} }
func (m *ConfirmPurchaseOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetLedgersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetLedgersRequest) Reset()         { *m = GetLedgersRequest{} }
func (m *GetLedgersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateLedgerRequest struct {
	Name string `json:"name"`
	// required name, account_number, account_name
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        string            `json:"note"`
	CreatedBy   string            `json:"created_by"`
}

func (m *CreateLedgerRequest) Reset()         { *m = CreateLedgerRequest{} }
func (m *CreateLedgerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateLedgerRequest struct {
	Id   dot.ID         `json:"id"`
	Name dot.NullString `json:"name"`
	// required name, account_number, account_name
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        dot.NullString    `json:"note"`
}

func (m *UpdateLedgerRequest) Reset()         { *m = UpdateLedgerRequest{} }
func (m *UpdateLedgerRequest) String() string { return jsonx.MustMarshalToString(m) }

type LedgersResponse struct {
	Ledgers []*Ledger        `json:"ledgers"`
	Paging  *common.PageInfo `json:"paging"`
}

func (m *LedgersResponse) Reset()         { *m = LedgersResponse{} }
func (m *LedgersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Ledger struct {
	Id          dot.ID            `json:"id"`
	Name        string            `json:"name"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	Note        string            `json:"note"`
	// enum: cash, bank
	Type      string   `json:"type"`
	CreatedBy dot.ID   `json:"created_by"`
	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *Ledger) Reset()         { *m = Ledger{} }
func (m *Ledger) String() string { return jsonx.MustMarshalToString(m) }

type RegisterShopRequest struct {
	// @required
	Name        string            `json:"name"`
	Address     *etop.Address     `json:"address"`
	Phone       string            `json:"phone"`
	BankAccount *etop.BankAccount `json:"bank_account"`
	WebsiteUrl  string            `json:"website_url"`
	ImageUrl    string            `json:"image_url"`
	Email       string            `json:"email"`
	UrlSlug     string            `json:"url_slug"`
	CompanyInfo *etop.CompanyInfo `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
}

func (m *RegisterShopRequest) Reset()         { *m = RegisterShopRequest{} }
func (m *RegisterShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterShopResponse struct {
	// @required
	Shop *etop.Shop `json:"shop"`
}

func (m *RegisterShopResponse) Reset()         { *m = RegisterShopResponse{} }
func (m *RegisterShopResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShopRequest struct {
	InventoryOverstock dot.NullBool      `json:"inventory_overstock"`
	Name               string            `json:"name"`
	Address            *etop.Address     `json:"address"`
	Phone              string            `json:"phone"`
	BankAccount        *etop.BankAccount `json:"bank_account"`
	WebsiteUrl         string            `json:"website_url"`
	ImageUrl           string            `json:"image_url"`
	Email              string            `json:"email"`
	AutoCreateFfm      dot.NullBool      `json:"auto_create_ffm"`
	// @deprecated use try_on instead
	GhnNoteCode *ghn_note_code.GHNNoteCode `json:"ghn_note_code"`
	TryOn       *try_on.TryOnCode          `json:"try_on"`
	CompanyInfo *etop.CompanyInfo          `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                                    `json:"money_transaction_rrule"`
	SurveyInfo                    []*etop.SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*etop.ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
}

func (m *UpdateShopRequest) Reset()         { *m = UpdateShopRequest{} }
func (m *UpdateShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShopResponse struct {
	Shop *etop.Shop `json:"shop"`
}

func (m *UpdateShopResponse) Reset()         { *m = UpdateShopResponse{} }
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

func (m *Collection) Reset()         { *m = Collection{} }
func (m *Collection) String() string { return jsonx.MustMarshalToString(m) }

type CreateCollectionRequest struct {
	// @required
	Name        string `json:"name"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`
	DescHtml    string `json:"desc_html"`
}

func (m *CreateCollectionRequest) Reset()         { *m = CreateCollectionRequest{} }
func (m *CreateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductCategoryRequest struct {
	ProductId  dot.ID `json:"product_id"`
	CategoryId dot.ID `json:"category_id"`
}

func (m *UpdateProductCategoryRequest) Reset()         { *m = UpdateProductCategoryRequest{} }
func (m *UpdateProductCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type CollectionsResponse struct {
	Collections []*ShopCollection `json:"collections"`
}

func (m *CollectionsResponse) Reset()         { *m = CollectionsResponse{} }
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

func (m *UpdateCollectionRequest) Reset()         { *m = UpdateCollectionRequest{} }
func (m *UpdateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsCollectionRequest struct {
	// @required
	CollectionId dot.ID   `json:"collection_id"`
	ProductIds   []dot.ID `json:"product_ids"`
}

func (m *UpdateProductsCollectionRequest) Reset()         { *m = UpdateProductsCollectionRequest{} }
func (m *UpdateProductsCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveProductsCollectionRequest struct {
	// @required
	CollectionId dot.ID   `json:"collection_id"`
	ProductIds   []dot.ID `json:"product_ids"`
}

func (m *RemoveProductsCollectionRequest) Reset()         { *m = RemoveProductsCollectionRequest{} }
func (m *RemoveProductsCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type EtopVariant struct {
	Id          dot.ID       `json:"id"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ShortDesc   string       `json:"short_desc"`
	DescHtml    string       `json:"desc_html"`
	ImageUrls   []string     `json:"image_urls"`
	ListPrice   int          `json:"list_price"`
	CostPrice   int          `json:"cost_price"`
	Attributes  []*Attribute `json:"attributes"`
}

func (m *EtopVariant) Reset()         { *m = EtopVariant{} }
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

func (m *EtopProduct) Reset()         { *m = EtopProduct{} }
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

	Attributes []*Attribute      `json:"attributes"`
	Product    *ShopShortProduct `json:"product"`
}

func (m *ShopVariant) Reset()         { *m = ShopVariant{} }
func (m *ShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVariantShopVariant struct {
	QuantityOnHand int `json:"quantity_on_hand"`
	QuantityPicked int `json:"quantity_picked"`
	CostPrice      int `json:"cost_price"`
	Quantity       int `json:"quantity"`
}

func (m *InventoryVariantShopVariant) Reset()         { *m = InventoryVariantShopVariant{} }
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
	Tags            []string                  `json:"tags"`
	Stags           []*Tag                    `json:"stags"`
	Note            string                    `json:"note"`
	Status          status3.Status            `json:"status"`
	IsAvailable     bool                      `json:"is_available"`
	ListPrice       int                       `json:"list_price"`
	RetailPrice     int                       `json:"retail_price"`
	CollectionIds   []dot.ID                  `json:"collection_ids"`
	Variants        []*ShopVariant            `json:"variants"`
	ProductSourceId dot.ID                    `json:"product_source_id"`
	CreatedAt       dot.Time                  `json:"created_at"`
	UpdatedAt       dot.Time                  `json:"updated_at"`
	ProductType     *product_type.ProductType `json:"product_type"`
	MetaFields      []*common.MetaField       `json:"meta_fields"`
	BrandId         dot.ID                    `json:"brand_id"`
}

func (m *ShopProduct) Reset()         { *m = ShopProduct{} }
func (m *ShopProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopShortProduct struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *ShopShortProduct) Reset()         { *m = ShopShortProduct{} }
func (m *ShopShortProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopCollection struct {
	Id          dot.ID `json:"id"`
	ShopId      dot.ID `json:"shop_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DescHtml    string `json:"desc_html"`
	ShortDesc   string `json:"short_desc"`
}

func (m *ShopCollection) Reset()         { *m = ShopCollection{} }
func (m *ShopCollection) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetVariantsRequest) Reset()         { *m = GetVariantsRequest{} }
func (m *GetVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCategoriesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCategoriesRequest) Reset()         { *m = GetCategoriesRequest{} }
func (m *GetCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariantsResponse struct {
	Paging   *common.PageInfo `json:"paging"`
	Variants []*ShopVariant   `json:"variants"`
}

func (m *ShopVariantsResponse) Reset()         { *m = ShopVariantsResponse{} }
func (m *ShopVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopProductsResponse struct {
	Paging   *common.PageInfo `json:"paging"`
	Products []*ShopProduct   `json:"products"`
}

func (m *ShopProductsResponse) Reset()         { *m = ShopProductsResponse{} }
func (m *ShopProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopCategoriesResponse struct {
	Paging     *common.PageInfo `json:"paging"`
	Categories []*ShopCategory  `json:"categories"`
}

func (m *ShopCategoriesResponse) Reset()         { *m = ShopCategoriesResponse{} }
func (m *ShopCategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantRequest struct {
	// @required
	Id          dot.ID         `json:"id"`
	Name        dot.NullString `json:"name"`
	Note        dot.NullString `json:"note"`
	Code        dot.NullString `json:"code"`
	CostPrice   dot.NullInt    `json:"cost_price"`
	ListPrice   dot.NullInt    `json:"list_price"`
	RetailPrice dot.NullInt    `json:"retail_price"`
	Description dot.NullString `json:"description"`
	ShortDesc   dot.NullString `json:"short_desc"`
	DescHtml    dot.NullString `json:"desc_html"`
	Attributes  []*Attribute   `json:"attributes"`
	// deprecated
	Sku string `json:"sku"`
}

func (m *UpdateVariantRequest) Reset()         { *m = UpdateVariantRequest{} }
func (m *UpdateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductRequest struct {
	// @required
	Id          dot.ID                    `json:"id"`
	Name        dot.NullString            `json:"name"`
	Code        dot.NullString            `json:"code"`
	Note        dot.NullString            `json:"note"`
	Unit        dot.NullString            `json:"unit"`
	Description dot.NullString            `json:"description"`
	ShortDesc   dot.NullString            `json:"short_desc"`
	DescHtml    dot.NullString            `json:"desc_html"`
	CostPrice   dot.NullInt               `json:"cost_price"`
	ListPrice   dot.NullInt               `json:"list_price"`
	RetailPrice dot.NullInt               `json:"retail_price"`
	ProductType *product_type.ProductType `json:"product_type"`
	MetaFields  *common.MetaField         `json:"meta_fields"`
	BrandId     dot.NullID                `json:"brand_id"`
}

func (m *UpdateProductRequest) Reset()         { *m = UpdateProductRequest{} }
func (m *UpdateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCategoryRequest struct {
	Id       dot.ID         `json:"id"`
	Name     dot.NullString `json:"name"`
	ParentId dot.ID         `json:"parent_id"`
}

func (m *UpdateCategoryRequest) Reset()         { *m = UpdateCategoryRequest{} }
func (m *UpdateCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantsRequest struct {
	Updates []*UpdateVariantRequest `json:"updates"`
}

func (m *UpdateVariantsRequest) Reset()         { *m = UpdateVariantsRequest{} }
func (m *UpdateVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsTagsRequest struct {
	// @required
	Ids        []dot.ID `json:"ids"`
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`
}

func (m *UpdateProductsTagsRequest) Reset()         { *m = UpdateProductsTagsRequest{} }
func (m *UpdateProductsTagsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantsResponse struct {
	Variants []*ShopVariant  `json:"variants"`
	Errors   []*common.Error `json:"errors"`
}

func (m *UpdateVariantsResponse) Reset()         { *m = UpdateVariantsResponse{} }
func (m *UpdateVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type AddVariantsRequest struct {
	Ids          []dot.ID `json:"ids"`
	Tags         []string `json:"tags"`
	CollectionId dot.ID   `json:"collection_id"`
}

func (m *AddVariantsRequest) Reset()         { *m = AddVariantsRequest{} }
func (m *AddVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddVariantsResponse struct {
	Variants []*ShopVariant  `json:"variants"`
	Errors   []*common.Error `json:"errors"`
}

func (m *AddVariantsResponse) Reset()         { *m = AddVariantsResponse{} }
func (m *AddVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type RemoveVariantsRequest struct {
	Ids []dot.ID `json:"ids"`
}

func (m *RemoveVariantsRequest) Reset()         { *m = RemoveVariantsRequest{} }
func (m *RemoveVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetOrdersRequest struct {
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
	Mixed   *etop.MixedAccount `json:"mixed"`
}

func (m *GetOrdersRequest) Reset()         { *m = GetOrdersRequest{} }
func (m *GetOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateOrdersStatusRequest struct {
	// @required
	Ids []dot.ID `json:"ids"`
	// @required
	Confirm      *status3.Status `json:"confirm"`
	CancelReason string          `json:"cancel_reason"`
	Status       status4.Status  `json:"status"`
}

func (m *UpdateOrdersStatusRequest) Reset()         { *m = UpdateOrdersStatusRequest{} }
func (m *UpdateOrdersStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmOrderRequest struct {
	OrderId dot.ID `json:"order_id"`
	// enum ('create', 'create')
	AutoInventoryVoucher  dot.NullString `json:"auto_inventory_voucher"`
	AutoCreateFulfillment bool           `json:"auto_create_fulfillment"`
}

func (m *ConfirmOrderRequest) Reset()         { *m = ConfirmOrderRequest{} }
func (m *ConfirmOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDRequest struct {
	OrderId dot.ID `json:"order_id"`
}

func (m *OrderIDRequest) Reset()         { *m = OrderIDRequest{} }
func (m *OrderIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDsRequest struct {
	OrderIds []dot.ID `json:"order_ids"`
}

func (m *OrderIDsRequest) Reset()         { *m = OrderIDsRequest{} }
func (m *OrderIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateFulfillmentsForOrderRequest struct {
	OrderId    dot.ID   `json:"order_id"`
	VariantIds []dot.ID `json:"variant_ids"`
}

func (m *CreateFulfillmentsForOrderRequest) Reset()         { *m = CreateFulfillmentsForOrderRequest{} }
func (m *CreateFulfillmentsForOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelOrderRequest struct {
	OrderId              dot.ID `json:"order_id"`
	CancelReason         string `json:"cancel_reason"`
	AutoInventoryVoucher string `json:"auto_inventory_voucher"`
}

func (m *CancelOrderRequest) Reset()         { *m = CancelOrderRequest{} }
func (m *CancelOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelOrdersRequest struct {
	Ids    []dot.ID `json:"ids"`
	Reason string   `json:"reason"`
}

func (m *CancelOrdersRequest) Reset()         { *m = CancelOrdersRequest{} }
func (m *CancelOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductSource struct {
	Id        dot.ID         `json:"id"`
	Type      string         `json:"type"`
	Name      string         `json:"name"`
	Status    status3.Status `json:"status"`
	UpdatedAt dot.Time       `json:"updated_at"`
	CreatedAt dot.Time       `json:"created_at"`
}

func (m *ProductSource) Reset()         { *m = ProductSource{} }
func (m *ProductSource) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type CreateProductSourceRequest struct {
}

func (m *CreateProductSourceRequest) Reset()         { *m = CreateProductSourceRequest{} }
func (m *CreateProductSourceRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductSourcesResponse struct {
	ProductSources []*ProductSource `json:"product_sources"`
}

func (m *ProductSourcesResponse) Reset()         { *m = ProductSourcesResponse{} }
func (m *ProductSourcesResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
}

func (m *CreateCategoryRequest) Reset()         { *m = CreateCategoryRequest{} }
func (m *CreateCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateProductRequest struct {
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Unit        string                    `json:"unit"`
	Note        string                    `json:"note"`
	Description string                    `json:"description"`
	ShortDesc   string                    `json:"short_desc"`
	DescHtml    string                    `json:"desc_html"`
	ImageUrls   []string                  `json:"image_urls"`
	CostPrice   int                       `json:"cost_price"`
	ListPrice   int                       `json:"list_price"`
	RetailPrice int                       `json:"retail_price"`
	ProductType *product_type.ProductType `json:"product_type"`
	BrandId     dot.ID                    `json:"brand_id"`
	MetaFields  []*common.MetaField       `json:"meta_fields"`
}

func (m *CreateProductRequest) Reset()         { *m = CreateProductRequest{} }
func (m *CreateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateVariantRequest struct {
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	ProductId   dot.ID       `json:"product_id"`
	Note        string       `json:"note"`
	Description string       `json:"description"`
	ShortDesc   string       `json:"short_desc"`
	DescHtml    string       `json:"desc_html"`
	ImageUrls   []string     `json:"image_urls"`
	Attributes  []*Attribute `json:"attributes"`
	CostPrice   int          `json:"cost_price"`
	ListPrice   int          `json:"list_price"`
	RetailPrice int          `json:"retail_price"`
}

func (m *CreateVariantRequest) Reset()         { *m = CreateVariantRequest{} }
func (m *CreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeprecatedCreateVariantRequest struct {
	// required
	ProductSourceId dot.ID `json:"product_source_id"`
	ProductId       dot.ID `json:"product_id"`
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string `json:"product_name"`
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	ShortDesc         string         `json:"short_desc"`
	DescHtml          string         `json:"desc_html"`
	ImageUrls         []string       `json:"image_urls"`
	Tags              []string       `json:"tags"`
	Status            status3.Status `json:"status"`
	CostPrice         int            `json:"cost_price"`
	ListPrice         int            `json:"list_price"`
	RetailPrice       int            `json:"retail_price"`
	Code              string         `json:"code"`
	QuantityAvailable int            `json:"quantity_available"`
	QuantityOnHand    int            `json:"quantity_on_hand"`
	QuantityReserved  int            `json:"quantity_reserved"`
	Attributes        []*Attribute   `json:"attributes"`
	Unit              string         `json:"unit"`
	// deprecated: use code instead
	Sku string `json:"sku"`
}

func (m *DeprecatedCreateVariantRequest) Reset()         { *m = DeprecatedCreateVariantRequest{} }
func (m *DeprecatedCreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConnectProductSourceResquest struct {
	ProductSourceId dot.ID `json:"product_source_id"`
}

func (m *ConnectProductSourceResquest) Reset()         { *m = ConnectProductSourceResquest{} }
func (m *ConnectProductSourceResquest) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type CreatePSCategoryRequest struct {
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
}

func (m *CreatePSCategoryRequest) Reset()         { *m = CreatePSCategoryRequest{} }
func (m *CreatePSCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsPSCategoryRequest struct {
	CategoryId dot.ID   `json:"category_id"`
	ProductIds []dot.ID `json:"product_ids"`
}

func (m *UpdateProductsPSCategoryRequest) Reset()         { *m = UpdateProductsPSCategoryRequest{} }
func (m *UpdateProductsPSCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductsCollectionResponse struct {
	Updated int             `json:"updated"`
	Errors  []*common.Error `json:"errors"`
}

func (m *UpdateProductsCollectionResponse) Reset()         { *m = UpdateProductsCollectionResponse{} }
func (m *UpdateProductsCollectionResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductSourceCategoryRequest struct {
	Id       dot.ID `json:"id"`
	ParentId dot.ID `json:"parent_id"`
	Name     string `json:"name"`
}

func (m *UpdateProductSourceCategoryRequest) Reset()         { *m = UpdateProductSourceCategoryRequest{} }
func (m *UpdateProductSourceCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

// deprecated
type GetProductSourceCategoriesRequest struct {
}

func (m *GetProductSourceCategoriesRequest) Reset()         { *m = GetProductSourceCategoriesRequest{} }
func (m *GetProductSourceCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsRequest struct {
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
	Mixed   *etop.MixedAccount `json:"mixed"`
	OrderId dot.ID             `json:"order_id"`
	Status  *status3.Status    `json:"status"`
}

func (m *GetFulfillmentsRequest) Reset()         { *m = GetFulfillmentsRequest{} }
func (m *GetFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentHistoryRequest struct {
	Paging  *common.Paging `json:"paging"`
	All     bool           `json:"all"`
	Id      dot.ID         `json:"id"`
	OrderId dot.ID         `json:"order_id"`
}

func (m *GetFulfillmentHistoryRequest) Reset()         { *m = GetFulfillmentHistoryRequest{} }
func (m *GetFulfillmentHistoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetBalanceShopResponse struct {
	Amount int `json:"amount"`
}

func (m *GetBalanceShopResponse) Reset()         { *m = GetBalanceShopResponse{} }
func (m *GetBalanceShopResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionsRequest) Reset()         { *m = GetMoneyTransactionsRequest{} }
func (m *GetMoneyTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetPublicFulfillmentRequest struct {
	// @Required
	Code string `json:"code"`
}

func (m *GetPublicFulfillmentRequest) Reset()         { *m = GetPublicFulfillmentRequest{} }
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
	OrderId dot.ID          `json:"order_id"`
	Status  *status3.Status `json:"status"`
}

func (m *UpdateOrderPaymentStatusRequest) Reset()         { *m = UpdateOrderPaymentStatusRequest{} }
func (m *UpdateOrderPaymentStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeFulfillmentsRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (m *SummarizeFulfillmentsRequest) Reset()         { *m = SummarizeFulfillmentsRequest{} }
func (m *SummarizeFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummarizeFulfillmentsResponse struct {
	Tables []*SummaryTable `json:"tables"`
}

func (m *SummarizeFulfillmentsResponse) Reset()         { *m = SummarizeFulfillmentsResponse{} }
func (m *SummarizeFulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type SummarizePOSResponse struct {
	Tables []*SummaryTable `json:"tables"`
}

func (m *SummarizePOSResponse) Reset()         { *m = SummarizePOSResponse{} }
func (m *SummarizePOSResponse) String() string { return jsonx.MustMarshalToString(m) }

type SummarizePOSRequest struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

func (m *SummarizePOSRequest) Reset()         { *m = SummarizePOSRequest{} }
func (m *SummarizePOSRequest) String() string { return jsonx.MustMarshalToString(m) }

type SummaryTable struct {
	Label   string          `json:"label"`
	Tags    []string        `json:"tags"`
	Columns []SummaryColRow `json:"columns"`
	Rows    []SummaryColRow `json:"rows"`
	Data    []SummaryItem   `json:"data"`
}

func (m *SummaryTable) Reset()         { *m = SummaryTable{} }
func (m *SummaryTable) String() string { return jsonx.MustMarshalToString(m) }

type SummaryColRow struct {
	Label  string `json:"label"`
	Spec   string `json:"spec"`
	Unit   string `json:"unit"`
	Indent int    `json:"indent"`
}

func (m *SummaryColRow) Reset()         { *m = SummaryColRow{} }
func (m *SummaryColRow) String() string { return jsonx.MustMarshalToString(m) }

type SummaryItem struct {
	Label     string   `json:"label"`
	Spec      string   `json:"spec"`
	Value     int      `json:"value"`
	Unit      string   `json:"unit"`
	ImageUrls []string `json:"image_urls"`
}

func (m *SummaryItem) Reset()         { *m = SummaryItem{} }
func (m *SummaryItem) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductsResponse struct {
	Data         *spreadsheet.SpreadsheetData `json:"data"`
	ImportErrors []*common.Error              `json:"import_errors"`
	CellErrors   []*common.Error              `json:"cell_errors"`
	ImportId     dot.ID                       `json:"import_id"`
	StocktakeID  dot.ID                       `json:"stocktake_id"`
}

func (m *ImportProductsResponse) Reset()         { *m = ImportProductsResponse{} }
func (m *ImportProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CalcBalanceShopResponse struct {
	Balance int `json:"balance"`
}

func (m *CalcBalanceShopResponse) Reset()         { *m = CalcBalanceShopResponse{} }
func (m *CalcBalanceShopResponse) String() string { return jsonx.MustMarshalToString(m) }

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

func (m *RequestExportRequest) Reset()         { *m = RequestExportRequest{} }
func (m *RequestExportRequest) String() string { return jsonx.MustMarshalToString(m) }

type RequestExportResponse struct {
	Id         string         `json:"id"`
	Filename   string         `json:"filename"`
	ExportType string         `json:"export_type"`
	Status     status4.Status `json:"status"`
}

func (m *RequestExportResponse) Reset()         { *m = RequestExportResponse{} }
func (m *RequestExportResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsRequest struct {
}

func (m *GetExportsRequest) Reset()         { *m = GetExportsRequest{} }
func (m *GetExportsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsResponse struct {
	ExportItems []*ExportItem `json:"export_items"`
}

func (m *GetExportsResponse) Reset()         { *m = GetExportsResponse{} }
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

func (m *ExportItem) Reset()         { *m = ExportItem{} }
func (m *ExportItem) String() string { return jsonx.MustMarshalToString(m) }

type GetExportsStatusRequest struct {
}

func (m *GetExportsStatusRequest) Reset()         { *m = GetExportsStatusRequest{} }
func (m *GetExportsStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type ExportStatusItem struct {
	Id            string        `json:"id"`
	ProgressMax   int           `json:"progress_max"`
	ProgressValue int           `json:"progress_value"`
	ProgressError int           `json:"progress_error"`
	Error         *common.Error `json:"error"`
}

func (m *ExportStatusItem) Reset()         { *m = ExportStatusItem{} }
func (m *ExportStatusItem) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizePartnerRequest struct {
	PartnerId dot.ID `json:"partner_id"`
}

func (m *AuthorizePartnerRequest) Reset()         { *m = AuthorizePartnerRequest{} }
func (m *AuthorizePartnerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetPartnersResponse struct {
	Partners []*etop.PublicAccountInfo `json:"partners"`
}

func (m *GetPartnersResponse) Reset()         { *m = GetPartnersResponse{} }
func (m *GetPartnersResponse) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizedPartnerResponse struct {
	Partner     *etop.PublicAccountInfo `json:"partner"`
	RedirectUrl string                  `json:"redirect_url"`
}

func (m *AuthorizedPartnerResponse) Reset()         { *m = AuthorizedPartnerResponse{} }
func (m *AuthorizedPartnerResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetAuthorizedPartnersResponse struct {
	Partners []*AuthorizedPartnerResponse `json:"partners"`
}

func (m *GetAuthorizedPartnersResponse) Reset()         { *m = GetAuthorizedPartnersResponse{} }
func (m *GetAuthorizedPartnersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (m *Attribute) Reset()         { *m = Attribute{} }
func (m *Attribute) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantImagesRequest struct {
	// @required
	Id         dot.ID   `json:"id"`
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`
}

func (m *UpdateVariantImagesRequest) Reset()         { *m = UpdateVariantImagesRequest{} }
func (m *UpdateVariantImagesRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductMetaFieldsRequest struct {
	// @required
	Id         dot.ID              `json:"id"`
	MetaFields []*common.MetaField `json:"meta_fields"`
}

func (m *UpdateProductMetaFieldsRequest) Reset()         { *m = UpdateProductMetaFieldsRequest{} }
func (m *UpdateProductMetaFieldsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (m *CategoriesResponse) Reset()         { *m = CategoriesResponse{} }
func (m *CategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type Category struct {
	Id              dot.ID `json:"id"`
	Name            string `json:"name"`
	ProductSourceId dot.ID `json:"product_source_id"`
	ParentId        dot.ID `json:"parent_id"`
	ShopId          dot.ID `json:"shop_id"`
}

func (m *Category) Reset()         { *m = Category{} }
func (m *Category) String() string { return jsonx.MustMarshalToString(m) }

type Tag struct {
	Id    dot.ID `json:"id"`
	Label string `json:"label"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return jsonx.MustMarshalToString(m) }

type ExternalAccountAhamove struct {
	Id               dot.ID `json:"id"`
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
}

func (m *ExternalAccountAhamove) Reset()         { *m = ExternalAccountAhamove{} }
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

type ExternalAccountHaravan struct {
	Id                                dot.ID   `json:"id"`
	ShopId                            dot.ID   `json:"shop_id"`
	Subdomain                         string   `json:"subdomain"`
	ExternalCarrierServiceId          int      `json:"external_carrier_service_id"`
	ExternalConnectedCarrierServiceAt dot.Time `json:"external_connected_carrier_service_at"`
	ExpiresAt                         dot.Time `json:"expires_at"`
	CreatedAt                         dot.Time `json:"created_at"`
	UpdatedAt                         dot.Time `json:"updated_at"`
}

func (m *ExternalAccountHaravan) Reset()         { *m = ExternalAccountHaravan{} }
func (m *ExternalAccountHaravan) String() string { return jsonx.MustMarshalToString(m) }

type ExternalAccountHaravanRequest struct {
	// @required
	Subdomain string `json:"subdomain"`
	// @required OAuth code
	Code string `json:"code"`
	// @required
	RedirectUri string `json:"redirect_uri"`
}

func (m *ExternalAccountHaravanRequest) Reset()         { *m = ExternalAccountHaravanRequest{} }
func (m *ExternalAccountHaravanRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerLiability struct {
	TotalOrders    int `json:"total_orders"`
	TotalAmount    int `json:"total_amount"`
	ReceivedAmount int `json:"received_amount"`
	Liability      int `json:"liability"`
}

func (m *CustomerLiability) Reset()         { *m = CustomerLiability{} }
func (m *CustomerLiability) String() string { return jsonx.MustMarshalToString(m) }

type Customer struct {
	Id        dot.ID             `json:"id"`
	ShopId    dot.ID             `json:"shop_id"`
	FullName  string             `json:"full_name"`
	Code      string             `json:"code"`
	Note      string             `json:"note"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email"`
	Gender    string             `json:"gender"`
	Type      string             `json:"type"`
	Birthday  string             `json:"birthday"`
	CreatedAt dot.Time           `json:"created_at"`
	UpdatedAt dot.Time           `json:"updated_at"`
	Status    status3.Status     `json:"status"`
	GroupIds  []dot.ID           `json:"group_ids"`
	Liability *CustomerLiability `json:"liability"`
}

func (m *Customer) Reset()         { *m = Customer{} }
func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerRequest struct {
	// @required
	FullName string `json:"full_name"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
	// enum ('individual', 'organization')
	Type string `json:"type"`
	Note string `json:"note"`
	// @required
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *CreateCustomerRequest) Reset()         { *m = CreateCustomerRequest{} }
func (m *CreateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerRequest struct {
	Id       dot.ID         `json:"id"`
	FullName dot.NullString `json:"full_name"`
	Gender   dot.NullString `json:"gender"`
	Birthday dot.NullString `json:"birthday"`
	// enum ('individual', 'organization','independent')
	Type  dot.NullString `json:"type"`
	Note  dot.NullString `json:"note"`
	Phone dot.NullString `json:"phone"`
	Email dot.NullString `json:"email"`
}

func (m *UpdateCustomerRequest) Reset()         { *m = UpdateCustomerRequest{} }
func (m *UpdateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCustomersRequest) Reset()         { *m = GetCustomersRequest{} }
func (m *GetCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomersResponse struct {
	Customers []*Customer      `json:"customers"`
	Paging    *common.PageInfo `json:"paging"`
}

func (m *CustomersResponse) Reset()         { *m = CustomersResponse{} }
func (m *CustomersResponse) String() string { return jsonx.MustMarshalToString(m) }

type SetCustomersStatusRequest struct {
	Ids    []dot.ID       `json:"ids"`
	Status status3.Status `json:"status"`
}

func (m *SetCustomersStatusRequest) Reset()         { *m = SetCustomersStatusRequest{} }
func (m *SetCustomersStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerDetailsResponse struct {
	Customer     *Customer                 `json:"customer"`
	SummaryItems []*IndependentSummaryItem `json:"summary_items"`
}

func (m *CustomerDetailsResponse) Reset()         { *m = CustomerDetailsResponse{} }
func (m *CustomerDetailsResponse) String() string { return jsonx.MustMarshalToString(m) }

type IndependentSummaryItem struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Spec  string `json:"spec"`
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

func (m *IndependentSummaryItem) Reset()         { *m = IndependentSummaryItem{} }
func (m *IndependentSummaryItem) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerAddressesRequest struct {
	CustomerId dot.ID `json:"customer_id"`
}

func (m *GetCustomerAddressesRequest) Reset()         { *m = GetCustomerAddressesRequest{} }
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
}

func (m *CustomerAddress) Reset()         { *m = CustomerAddress{} }
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

func (m *CreateCustomerAddressRequest) Reset()         { *m = CreateCustomerAddressRequest{} }
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

func (m *UpdateCustomerAddressRequest) Reset()         { *m = UpdateCustomerAddressRequest{} }
func (m *UpdateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddressesResponse struct {
	Addresses []*CustomerAddress `json:"addresses"`
}

func (m *CustomerAddressesResponse) Reset()         { *m = CustomerAddressesResponse{} }
func (m *CustomerAddressesResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductStatusRequest struct {
	Ids    []dot.ID       `json:"ids"`
	Status status3.Status `json:"status"`
}

func (m *UpdateProductStatusRequest) Reset()         { *m = UpdateProductStatusRequest{} }
func (m *UpdateProductStatusRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductStatusResponse struct {
	Updated int `json:"updated"`
}

func (m *UpdateProductStatusResponse) Reset()         { *m = UpdateProductStatusResponse{} }
func (m *UpdateProductStatusResponse) String() string { return jsonx.MustMarshalToString(m) }

type PaymentTradingOrderRequest struct {
	OrderId         dot.ID                           `json:"order_id"`
	Desc            string                           `json:"desc"`
	ReturnUrl       string                           `json:"return_url"`
	Amount          int                              `json:"amount"`
	PaymentProvider payment_provider.PaymentProvider `json:"payment_provider"`
}

func (m *PaymentTradingOrderRequest) Reset()         { *m = PaymentTradingOrderRequest{} }
func (m *PaymentTradingOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type PaymentTradingOrderResponse struct {
	Url string `json:"url"`
}

func (m *PaymentTradingOrderResponse) Reset()         { *m = PaymentTradingOrderResponse{} }
func (m *PaymentTradingOrderResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantAttributesRequest struct {
	// @required
	VariantId  dot.ID       `json:"variant_id"`
	Attributes []*Attribute `json:"attributes"`
}

func (m *UpdateVariantAttributesRequest) Reset()         { *m = UpdateVariantAttributesRequest{} }
func (m *UpdateVariantAttributesRequest) String() string { return jsonx.MustMarshalToString(m) }

type PaymentCheckReturnDataRequest struct {
	Id                    string                           `json:"id"`
	Code                  string                           `json:"code"`
	PaymentStatus         string                           `json:"payment_status"`
	Amount                int                              `json:"amount"`
	ExternalTransactionId string                           `json:"external_transaction_id"`
	PaymentProvider       payment_provider.PaymentProvider `json:"payment_provider"`
}

func (m *PaymentCheckReturnDataRequest) Reset()         { *m = PaymentCheckReturnDataRequest{} }
func (m *PaymentCheckReturnDataRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopCategory struct {
	Id       dot.ID `json:"id"`
	Name     string `json:"name"`
	ParentId dot.ID `json:"parent_id"`
	ShopId   dot.ID `json:"shop_id"`
	Status   dot.ID `json:"status"`
}

func (m *ShopCategory) Reset()         { *m = ShopCategory{} }
func (m *ShopCategory) String() string { return jsonx.MustMarshalToString(m) }

type GetCollectionsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCollectionsRequest) Reset()         { *m = GetCollectionsRequest{} }
func (m *GetCollectionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopCollectionsResponse struct {
	Paging      *common.PageInfo  `json:"paging"`
	Collections []*ShopCollection `json:"collections"`
}

func (m *ShopCollectionsResponse) Reset()         { *m = ShopCollectionsResponse{} }
func (m *ShopCollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type AddShopProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *AddShopProductCollectionRequest) Reset()         { *m = AddShopProductCollectionRequest{} }
func (m *AddShopProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveShopProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *RemoveShopProductCollectionRequest) Reset()         { *m = RemoveShopProductCollectionRequest{} }
func (m *RemoveShopProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddCustomerToGroupRequest struct {
	CustomerIds []dot.ID `json:"customer_ids"`
	GroupId     dot.ID   `json:"group_id"`
}

func (m *AddCustomerToGroupRequest) Reset()         { *m = AddCustomerToGroupRequest{} }
func (m *AddCustomerToGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveCustomerOutOfGroupRequest struct {
	CustomerIds []dot.ID `json:"customer_ids"`
	GroupId     dot.ID   `json:"group_id"`
}

func (m *RemoveCustomerOutOfGroupRequest) Reset()         { *m = RemoveCustomerOutOfGroupRequest{} }
func (m *RemoveCustomerOutOfGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type SupplierLiability struct {
	TotalPurchaseOrders int `json:"total_purchase_orders"`
	TotalAmount         int `json:"total_amount"`
	PaidAmount          int `json:"paid_amount"`
	Liability           int `json:"liability"`
}

func (m *SupplierLiability) Reset()         { *m = SupplierLiability{} }
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

func (m *Supplier) Reset()         { *m = Supplier{} }
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

func (m *CreateSupplierRequest) Reset()         { *m = CreateSupplierRequest{} }
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

func (m *UpdateSupplierRequest) Reset()         { *m = UpdateSupplierRequest{} }
func (m *UpdateSupplierRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSuppliersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetSuppliersRequest) Reset()         { *m = GetSuppliersRequest{} }
func (m *GetSuppliersRequest) String() string { return jsonx.MustMarshalToString(m) }

type SuppliersResponse struct {
	Suppliers []*Supplier      `json:"suppliers"`
	Paging    *common.PageInfo `json:"paging"`
}

func (m *SuppliersResponse) Reset()         { *m = SuppliersResponse{} }
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

func (m *Carrier) Reset()         { *m = Carrier{} }
func (m *Carrier) String() string { return jsonx.MustMarshalToString(m) }

type CreateCarrierRequest struct {
	FullName string `json:"full_name"`
	Note     string `json:"note"`
}

func (m *CreateCarrierRequest) Reset()         { *m = CreateCarrierRequest{} }
func (m *CreateCarrierRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCarrierRequest struct {
	Id       dot.ID         `json:"id"`
	FullName dot.NullString `json:"full_name"`
	Note     dot.NullString `json:"note"`
}

func (m *UpdateCarrierRequest) Reset()         { *m = UpdateCarrierRequest{} }
func (m *UpdateCarrierRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCarriersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCarriersRequest) Reset()         { *m = GetCarriersRequest{} }
func (m *GetCarriersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CarriersResponse struct {
	Carriers []*Carrier       `json:"carriers"`
	Paging   *common.PageInfo `json:"paging"`
}

func (m *CarriersResponse) Reset()         { *m = CarriersResponse{} }
func (m *CarriersResponse) String() string { return jsonx.MustMarshalToString(m) }

type ReceiptLine struct {
	RefId  dot.ID `json:"ref_id"`
	Title  string `json:"title"`
	Amount int    `json:"amount"`
}

func (m *ReceiptLine) Reset()         { *m = ReceiptLine{} }
func (m *ReceiptLine) String() string { return jsonx.MustMarshalToString(m) }

type Trader struct {
	Id       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Deleted  bool   `json:"deleted"`
}

func (m *Trader) Reset()         { *m = Trader{} }
func (m *Trader) String() string { return jsonx.MustMarshalToString(m) }

type Receipt struct {
	Id          dot.ID         `json:"id"`
	ShopId      dot.ID         `json:"shop_id"`
	TraderId    dot.ID         `json:"trader_id"`
	CreatedBy   dot.ID         `json:"created_by"`
	CreatedType string         `json:"created_type"`
	Code        string         `json:"code"`
	Title       string         `json:"title"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Amount      int            `json:"amount"`
	LedgerId    dot.ID         `json:"ledger_id"`
	RefType     string         `json:"ref_type"`
	Lines       []*ReceiptLine `json:"lines"`
	Status      status3.Status `json:"status"`
	Reason      string         `json:"reason"`
	PaidAt      dot.Time       `json:"paid_at"`
	CreatedAt   dot.Time       `json:"created_at"`
	UpdatedAt   dot.Time       `json:"updated_at"`
	ConfirmedAt dot.Time       `json:"confirmed_at"`
	CancelledAt dot.Time       `json:"cancelled_at"`
	User        *etop.User     `json:"user"`
	Trader      *Trader        `json:"trader"`
	Ledger      *Ledger        `json:"ledger"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return jsonx.MustMarshalToString(m) }

type CreateReceiptRequest struct {
	TraderId dot.ID `json:"trader_id"`
	Title    string `json:"title"`
	// enum('receipt', 'payment')
	Type        string `json:"type"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	LedgerId    dot.ID `json:"ledger_id"`
	// enum ('order', 'fulfillment', 'inventory voucher'
	RefType string         `json:"ref_type"`
	PaidAt  dot.Time       `json:"paid_at"`
	Lines   []*ReceiptLine `json:"lines"`
}

func (m *CreateReceiptRequest) Reset()         { *m = CreateReceiptRequest{} }
func (m *CreateReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReceiptRequest struct {
	Id          dot.ID         `json:"id"`
	TraderId    dot.NullID     `json:"trader_id"`
	Title       dot.NullString `json:"title"`
	Description dot.NullString `json:"description"`
	Amount      dot.NullInt    `json:"amount"`
	LedgerId    dot.NullID     `json:"ledger_id"`
	// enum ('order', 'fulfillment', 'inventory voucher'
	RefType dot.NullString `json:"ref_type"`
	PaidAt  dot.Time       `json:"paid_at"`
	Lines   []*ReceiptLine `json:"lines"`
}

func (m *UpdateReceiptRequest) Reset()         { *m = UpdateReceiptRequest{} }
func (m *UpdateReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelReceiptRequest struct {
	Id     dot.ID `json:"id"`
	Reason string `json:"reason"`
}

func (m *CancelReceiptRequest) Reset()         { *m = CancelReceiptRequest{} }
func (m *CancelReceiptRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetReceiptsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetReceiptsRequest) Reset()         { *m = GetReceiptsRequest{} }
func (m *GetReceiptsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetReceiptsByLedgerTypeRequest struct {
	Type    string           `json:"type"`
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetReceiptsByLedgerTypeRequest) Reset()         { *m = GetReceiptsByLedgerTypeRequest{} }
func (m *GetReceiptsByLedgerTypeRequest) String() string { return jsonx.MustMarshalToString(m) }

type ReceiptsResponse struct {
	TotalAmountConfirmedReceipt int              `json:"total_amount_confirmed_receipt"`
	TotalAmountConfirmedPayment int              `json:"total_amount_confirmed_payment"`
	Receipts                    []*Receipt       `json:"receipts"`
	Paging                      *common.PageInfo `json:"paging"`
}

func (m *ReceiptsResponse) Reset()         { *m = ReceiptsResponse{} }
func (m *ReceiptsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetShopCollectionsByProductIDRequest struct {
	ProductId dot.ID `json:"product_id"`
}

func (m *GetShopCollectionsByProductIDRequest) Reset()         { *m = GetShopCollectionsByProductIDRequest{} }
func (m *GetShopCollectionsByProductIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateInventoryVoucherRequest struct {
	RefId   dot.ID `json:"ref_id"`
	RefType string `json:"ref_type"`
	//enum "in" or "out" only for ref_type = "order"
	Type string `json:"type"`
}

func (m *CreateInventoryVoucherRequest) Reset()         { *m = CreateInventoryVoucherRequest{} }
func (m *CreateInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVoucherLine struct {
	VariantId   dot.ID      `json:"variant_id"`
	Code        string      `json:"code"`
	VariantName string      `json:"variant_name"`
	ProductId   dot.ID      `json:"product_id"`
	ProductName string      `json:"product_name"`
	ImageUrl    string      `json:"image_url"`
	Attributes  []Attribute `json:"attributes"`
	Price       int         `json:"price"`
	Quantity    int         `json:"quantity"`
}

func (m *InventoryVoucherLine) Reset()         { *m = InventoryVoucherLine{} }
func (m *InventoryVoucherLine) String() string { return jsonx.MustMarshalToString(m) }

type CreateInventoryVoucherResponse struct {
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *CreateInventoryVoucherResponse) Reset()         { *m = CreateInventoryVoucherResponse{} }
func (m *CreateInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmInventoryVoucherRequest struct {
	Id dot.ID `json:"id"`
}

func (m *ConfirmInventoryVoucherRequest) Reset()         { *m = ConfirmInventoryVoucherRequest{} }
func (m *ConfirmInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmInventoryVoucherResponse struct {
	InventoryVoucher *InventoryVoucher `json:"inventory_voucher"`
}

func (m *ConfirmInventoryVoucherResponse) Reset()         { *m = ConfirmInventoryVoucherResponse{} }
func (m *ConfirmInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type CancelInventoryVoucherRequest struct {
	Id     dot.ID `json:"id"`
	Reason string `json:"reason"`
}

func (m *CancelInventoryVoucherRequest) Reset()         { *m = CancelInventoryVoucherRequest{} }
func (m *CancelInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelInventoryVoucherResponse struct {
	Inventory *InventoryVoucher `json:"inventory"`
}

func (m *CancelInventoryVoucherResponse) Reset()         { *m = CancelInventoryVoucherResponse{} }
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

func (m *UpdateInventoryVoucherRequest) Reset()         { *m = UpdateInventoryVoucherRequest{} }
func (m *UpdateInventoryVoucherRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateInventoryVoucherResponse struct {
	InventoryVoucher *InventoryVoucher `json:"inventory_voucher"`
}

func (m *UpdateInventoryVoucherResponse) Reset()         { *m = UpdateInventoryVoucherResponse{} }
func (m *UpdateInventoryVoucherResponse) String() string { return jsonx.MustMarshalToString(m) }

type AdjustInventoryQuantityRequest struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
	Note              string              `json:"note"`
}

func (m *AdjustInventoryQuantityRequest) Reset()         { *m = AdjustInventoryQuantityRequest{} }
func (m *AdjustInventoryQuantityRequest) String() string { return jsonx.MustMarshalToString(m) }

type AdjustInventoryQuantityResponse struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *AdjustInventoryQuantityResponse) Reset()         { *m = AdjustInventoryQuantityResponse{} }
func (m *AdjustInventoryQuantityResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsRequest struct {
	Paging common.Paging `json:"paging"`
}

func (m *GetInventoryVariantsRequest) Reset()         { *m = GetInventoryVariantsRequest{} }
func (m *GetInventoryVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsResponse struct {
	InventoryVariants []*InventoryVariant `json:"inventory_variants"`
}

func (m *GetInventoryVariantsResponse) Reset()         { *m = GetInventoryVariantsResponse{} }
func (m *GetInventoryVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantsByVariantIDsRequest struct {
	VariantIds []dot.ID `json:"variant_ids"`
}

func (m *GetInventoryVariantsByVariantIDsRequest) Reset() {
	*m = GetInventoryVariantsByVariantIDsRequest{}
}
func (m *GetInventoryVariantsByVariantIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

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

func (m *InventoryVariant) Reset()         { *m = InventoryVariant{} }
func (m *InventoryVariant) String() string { return jsonx.MustMarshalToString(m) }

type InventoryVariantQuantity struct {
	QuantityOnHand int `json:"quantity_on_hand"`
	QuantityPicked int `json:"quantity_picked"`
	Quantity       int `json:"quantity"`
}

func (m *InventoryVariantQuantity) Reset()         { *m = InventoryVariantQuantity{} }
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

func (m *InventoryVoucher) Reset()         { *m = InventoryVoucher{} }
func (m *InventoryVoucher) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetInventoryVouchersRequest) Reset()         { *m = GetInventoryVouchersRequest{} }
func (m *GetInventoryVouchersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersByIDsRequest struct {
	Ids []dot.ID `json:"ids"`
}

func (m *GetInventoryVouchersByIDsRequest) Reset()         { *m = GetInventoryVouchersByIDsRequest{} }
func (m *GetInventoryVouchersByIDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersResponse struct {
	InventoryVouchers []*InventoryVoucher `json:"inventory_vouchers"`
}

func (m *GetInventoryVouchersResponse) Reset()         { *m = GetInventoryVouchersResponse{} }
func (m *GetInventoryVouchersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroup struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *CustomerGroup) Reset()         { *m = CustomerGroup{} }
func (m *CustomerGroup) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerGroupRequest struct {
	Name string `json:"name"`
}

func (m *CreateCustomerGroupRequest) Reset()         { *m = CreateCustomerGroupRequest{} }
func (m *CreateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerGroupRequest struct {
	GroupId dot.ID `json:"group_id"`
	Name    string `json:"name"`
}

func (m *UpdateCustomerGroupRequest) Reset()         { *m = UpdateCustomerGroupRequest{} }
func (m *UpdateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerGroupsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetCustomerGroupsRequest) Reset()         { *m = GetCustomerGroupsRequest{} }
func (m *GetCustomerGroupsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupsResponse struct {
	CustomerGroups []*CustomerGroup `json:"customer_groups"`
	Paging         *common.PageInfo `json:"paging"`
}

func (m *CustomerGroupsResponse) Reset()         { *m = CustomerGroupsResponse{} }
func (m *CustomerGroupsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetOrdersByReceiptIDRequest struct {
	ReceiptId dot.ID `json:"receipt_id"`
}

func (m *GetOrdersByReceiptIDRequest) Reset()         { *m = GetOrdersByReceiptIDRequest{} }
func (m *GetOrdersByReceiptIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVariantRequest struct {
	VariantId dot.ID `json:"variant_id"`
}

func (m *GetInventoryVariantRequest) Reset()         { *m = GetInventoryVariantRequest{} }
func (m *GetInventoryVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateBrandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *CreateBrandRequest) Reset()         { *m = CreateBrandRequest{} }
func (m *CreateBrandRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateBrandRequest struct {
	Id          dot.ID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *UpdateBrandRequest) Reset()         { *m = UpdateBrandRequest{} }
func (m *UpdateBrandRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteBrandResponse struct {
	Count int `json:"count"`
}

func (m *DeleteBrandResponse) Reset()         { *m = DeleteBrandResponse{} }
func (m *DeleteBrandResponse) String() string { return jsonx.MustMarshalToString(m) }

type Brand struct {
	ShopId      dot.ID   `json:"shop_id"`
	Id          dot.ID   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
}

func (m *Brand) Reset()         { *m = Brand{} }
func (m *Brand) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsByIDsResponse struct {
	Brands []*Brand `json:"brands"`
}

func (m *GetBrandsByIDsResponse) Reset()         { *m = GetBrandsByIDsResponse{} }
func (m *GetBrandsByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsRequest struct {
	Paging common.Paging `json:"paging"`
}

func (m *GetBrandsRequest) Reset()         { *m = GetBrandsRequest{} }
func (m *GetBrandsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetBrandsResponse struct {
	Brands []*Brand         `json:"brands"`
	Paging *common.PageInfo `json:"paging"`
}

func (m *GetBrandsResponse) Reset()         { *m = GetBrandsResponse{} }
func (m *GetBrandsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInventoryVouchersByReferenceRequest struct {
	RefId dot.ID `json:"ref_id"`
	// enum ('order', 'purchase_order', 'return', 'purchase_order')
	RefType string `json:"ref_type"`
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
func (m *GetInventoryVouchersByReferenceResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateOrderShippingInfoRequest struct {
	OrderId         dot.ID               `json:"order_id"`
	Shipping        *types.OrderShipping `json:"shipping"`
	ShippingAddress *types.OrderAddress  `json:"shipping_address"`
}

func (m *UpdateOrderShippingInfoRequest) Reset()         { *m = UpdateOrderShippingInfoRequest{} }
func (m *UpdateOrderShippingInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesByIDsResponse struct {
	Stocktakes []*Stocktake `json:"stocktakes"`
}

func (m *GetStocktakesByIDsResponse) Reset()         { *m = GetStocktakesByIDsResponse{} }
func (m *GetStocktakesByIDsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateStocktakeRequest struct {
	TotalQuantity int    `json:"total_quantity"`
	Note          string `json:"note"`
	//  length more than one
	Lines []*StocktakeLine `json:"lines"`
}

func (m *CreateStocktakeRequest) Reset()         { *m = CreateStocktakeRequest{} }
func (m *CreateStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateStocktakeRequest struct {
	Id            dot.ID `json:"id"`
	TotalQuantity int    `json:"total_quantity"`
	Note          string `json:"note"`
	//  length more than one
	Lines []*StocktakeLine `json:"lines"`
}

func (m *UpdateStocktakeRequest) Reset()         { *m = UpdateStocktakeRequest{} }
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
}

func (m *Stocktake) Reset()         { *m = Stocktake{} }
func (m *Stocktake) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetStocktakesRequest) Reset()         { *m = GetStocktakesRequest{} }
func (m *GetStocktakesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetStocktakesResponse struct {
	Stocktakes []*Stocktake     `json:"stocktakes"`
	Paging     *common.PageInfo `json:"paging"`
}

func (m *GetStocktakesResponse) Reset()         { *m = GetStocktakesResponse{} }
func (m *GetStocktakesResponse) String() string { return jsonx.MustMarshalToString(m) }

type StocktakeLine struct {
	ProductId   dot.ID       `json:"product_id"`
	ProductName string       `json:"product_name"`
	VariantName string       `json:"variant_name"`
	VariantId   dot.ID       `json:"variant_id"`
	OldQuantity int          `json:"old_quantity"`
	NewQuantity int          `json:"new_quantity"`
	Code        string       `json:"code"`
	ImageUrl    string       `json:"image_url"`
	CostPrice   int          `json:"cost_price"`
	Attributes  []*Attribute `json:"attributes"`
}

func (m *StocktakeLine) Reset()         { *m = StocktakeLine{} }
func (m *StocktakeLine) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmStocktakeRequest struct {
	Id                   dot.ID `json:"id"`
	AutoInventoryVoucher string `json:"auto_inventory_voucher"`
}

func (m *ConfirmStocktakeRequest) Reset()         { *m = ConfirmStocktakeRequest{} }
func (m *ConfirmStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantsBySupplierIDRequest struct {
	SupplierId dot.ID `json:"supplier_id"`
}

func (m *GetVariantsBySupplierIDRequest) Reset()         { *m = GetVariantsBySupplierIDRequest{} }
func (m *GetVariantsBySupplierIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSuppliersByVariantIDRequest struct {
	VariantId dot.ID `json:"variant_id"`
}

func (m *GetSuppliersByVariantIDRequest) Reset()         { *m = GetSuppliersByVariantIDRequest{} }
func (m *GetSuppliersByVariantIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelStocktakeRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelStocktakeRequest) Reset()         { *m = CancelStocktakeRequest{} }
func (m *CancelStocktakeRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateInventoryVariantCostPriceResponse struct {
	InventoryVariant *InventoryVariant `json:"inventory_variant"`
}

func (m *UpdateInventoryVariantCostPriceResponse) Reset() {
	*m = UpdateInventoryVariantCostPriceResponse{}
}
func (m *UpdateInventoryVariantCostPriceResponse) String() string { return jsonx.MustMarshalToString(m) }

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
