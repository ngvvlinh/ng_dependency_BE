package whitelabel

import (
	catalogtypes "o.o/api/main/catalog/types"
	"o.o/api/shopping/customering/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type ImportProductRequest struct {
	ExternalID   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	ExternalBrandID    string `json:"external_brand_id"`
	ExternalCategoryID string `json:"external_category_id"`

	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	ImageUrls   []string `json:"image_urls"`

	CodePrice   int      `json:"code_price"`
	ListPrice   int      `json:"list_price"`
	Note        string   `json:"note"`
	RetailPrice int      `json:"retail_price"`
	Unit        string   `json:"unit"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
	DeletedAt   dot.Time `json:"deleted_at"`
}

func (m *ImportProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductsRequest struct {
	Products []*ImportProductRequest `json:"products"`
}

func (m *ImportProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type Product struct {
	ExternalId   dot.NullString `json:"external_id"`
	ExternalCode dot.NullString `json:"external_code"`

	ExternalBrandID    string `json:"external_brand_id"`
	ExternalCategoryID string `json:"external_category_id"`

	// @required
	Id        dot.ID `json:"id"`
	PartnerID dot.ID `json:"partner_id"`
	ShopID    dot.ID `json:"shop_id"`

	Name        dot.NullString     `json:"name"`
	Description dot.NullString     `json:"description"`
	ShortDesc   dot.NullString     `json:"short_desc"`
	ImageUrls   []string           `json:"image_urls"`
	CategoryId  dot.NullID         `json:"category_id"`
	Note        dot.NullString     `json:"note"`
	Status      status3.NullStatus `json:"status"`
	ListPrice   dot.NullInt        `json:"list_price"`
	RetailPrice dot.NullInt        `json:"retail_price"`

	CreatedAt dot.Time   `json:"created_at"`
	UpdatedAt dot.Time   `json:"updated_at"`
	BrandId   dot.NullID `json:"brand_id"`
}

func (m *Product) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductsResponse struct {
	Products []*Product `json:"products"`
}

func (m *ImportProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportBrandRequest struct {
	ExternalID  string   `json:"external_id"`
	BrandName   string   `json:"brand_name"`
	Description string   `json:"description"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
	DeletedAt   dot.Time `json:"deleted_at"`
}

func (m *ImportBrandRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportBrandsRequest struct {
	Brands []*ImportBrandRequest `json:"brands"`
}

func (m *ImportBrandsRequest) String() string { return jsonx.MustMarshalToString(m) }

type Brand struct {
	ID          dot.ID   `json:"id"`
	PartnerID   dot.ID   `json:"partner_id"`
	ShopID      dot.ID   `json:"shop_id"`
	ExternalID  string   `json:"external_id"`
	BrandName   string   `json:"brand_name"`
	Description string   `json:"description"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
	DeletedAt   dot.Time `json:"deleted_at"`
}

func (m *Brand) String() string { return jsonx.MustMarshalToString(m) }

type ImportBrandsResponse struct {
	Brands []*Brand `json:"brands"`
}

func (m *ImportBrandsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportCategoryRequest struct {
	ExternalID       string   `json:"external_id"`
	ExternalParentID string   `json:"external_parent_id"`
	Name             string   `json:"name"`
	CreatedAt        dot.Time `json:"created_at"`
	UpdatedAt        dot.Time `json:"updated_at"`
	DeletedAt        dot.Time `json:"deleted_at"`
}

func (m *ImportCategoryRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportCategoriesRequest struct {
	Categories []*ImportCategoryRequest `json:"categories"`
}

func (m *ImportCategoriesRequest) String() string { return jsonx.MustMarshalToString(m) }

type Category struct {
	ID               dot.ID   `json:"id"`
	ShopID           dot.ID   `json:"shop_id"`
	PartnerID        dot.ID   `json:"partner_id"`
	ExternalID       string   `json:"external_id"`
	ExternalParentID string   `json:"external_parent_id"`
	ParentID         dot.ID   `json:"parent_id"`
	Name             string   `json:"name"`
	CreatedAt        dot.Time `json:"created_at"`
	UpdatedAt        dot.Time `json:"updated_at"`
	DeletedAt        dot.Time `json:"deleted_at"`
}

func (m *Category) String() string { return jsonx.MustMarshalToString(m) }

type ImportCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (m *ImportCategoriesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportCustomerRequest struct {
	ExternalID   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

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

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
	DeletedAt dot.Time `json:"deleted_at"`
}

func (m *ImportCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportCustomersRequest struct {
	Customers []*ImportCustomerRequest `json:"customers"`
}

func (m *ImportCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type Customer struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	ID        dot.ID `json:"id"`
	PartnerID dot.ID `json:"partner_id"`
	ShopID    dot.ID `json:"shop_id"`

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

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
	DeletedAt dot.Time `json:"deleted_at"`
}

func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

type ImportCustomersResponse struct {
	Customers []*Customer `json:"customers"`
}

func (m *ImportCustomersResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariant struct {
	ExternalId        string `json:"external_id"`
	ExternalCode      string `json:"external_code"`
	ExternalProductId string `json:"external_product_id"`

	Id        dot.ID `json:"id"`
	ProductID dot.ID `json:"product_id"`
	ShopID    dot.ID `json:"shop_id"`
	PartnerID dot.ID `json:"partner_id"`

	Code string `json:"code"`

	Name        string         `json:"name"`
	Description string         `json:"description"`
	ShortDesc   string         `json:"short_desc"`
	ImageUrls   []string       `json:"image_urls"`
	ListPrice   int            `json:"list_price"`
	RetailPrice int            `json:"retail_price"`
	Note        string         `json:"note"`
	Status      status3.Status `json:"status"`

	CostPrice int `json:"cost_price"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
	DeletedAt dot.Time `json:"deleted_at"`
}

func (m *ShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type ImportShopVariantRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	Name string `json:"name"`

	ExternalProductId string `json:"external_product_id"`

	Note string `json:"note"`

	Description string `json:"description"`

	ShortDesc string `json:"short_desc"`

	CostPrice int `json:"cost_price"`

	ListPrice int `json:"list_price"`

	RetailPrice int `json:"retail_price"`

	ImageUrls []string `json:"image_urls"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`

	CreatedAt dot.Time `json:"created_at"`

	UpdatedAt dot.Time `json:"updated_at"`

	DeletedAt dot.Time `json:"deleted_at"`
}

func (m *ImportShopVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportShopVariantsRequest struct {
	Variants []*ImportShopVariantRequest `json:"variants"`
}

func (m *ImportShopVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportShopVariantsResponse struct {
	Variants []*ShopVariant `json:"variants"`
}

func (m *ImportShopVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportCollectionRequest struct {
	ExternalID  string `json:"external_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DescHTML    string `json:"desc_html"`
	ShortDesc   string `json:"short_desc"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *ImportCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportCollectionsRequest struct {
	Collections []*ImportCollectionRequest `json:"collections"`
}

func (m *ImportCollectionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportCollectionsResponse struct {
	Collections []*Collection `json:"collections"`
}

func (m *ImportCollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type Collection struct {
	ID        dot.ID `json:"id"`
	ShopID    dot.ID `json:"shop_id"`
	PartnerID dot.ID `json:"partner_id"`

	ExternalID string `json:"external_id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	DescHTML    string `json:"desc_html"`
	ShortDesc   string `json:"short_desc"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *Collection) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductCollectionRequest struct {
	ExternalProductID    string `json:"external_product_id"`
	ExternalCollectionID string `json:"external_collection_id"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *ImportProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductCollectionsRequest struct {
	ProductCollections []*ImportProductCollectionRequest `json:"product_collections"`
}

func (m *ImportProductCollectionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ImportProductCollectionsResponse struct {
	ProductCollections []*ProductCollection `json:"product_collections"`
}

func (m *ImportProductCollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ProductCollection struct {
	PartnerID dot.ID `json:"partner_id"`
	ShopID    dot.ID `json:"shop_id"`

	ExternalProductID    string `json:"external_product_id"`
	ExternalCollectionID string `json:"external_collection_id"`
	ProductID            dot.ID `json:"product_id"`
	CollectionID         dot.ID `json:"collection_id"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *ProductCollection) String() string { return jsonx.MustMarshalToString(m) }
