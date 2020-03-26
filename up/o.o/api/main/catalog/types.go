package catalog

import (
	"time"

	"o.o/api/main/catalog/types"
	"o.o/api/top/types/etc/product_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	cmutil "o.o/capi/util"
)

// +gen:event:topic=event/catalog

type ShopProduct struct {
	ExternalID         string `json:"external_id"`
	ExternalCode       string `json:"external_code"`
	ExternalBrandID    string `json:"external_brand_id"`
	ExternalCategoryID string `json:"external_category_id"`
	PartnerID          dot.ID `json:"partner_id"`

	ShopID    dot.ID `json:"shop_id"`
	ProductID dot.ID `json:"product_id"`
	Code      string `json:"code"`
	CodeNorm  int    `json:"code_norm"`
	Name      string `json:"name"`
	NameNorm  string `json:"name_norm"`
	Unit      string `json:"unit"`

	ImageURLs   []string `json:"image_urls"`
	Note        string   `json:"note"`
	ShortDesc   string   `json:"short_desc"`
	Description string   `json:"description"`
	DescHTML    string   `json:"desc_html"`

	CostPrice   int `json:"cost_price"`
	ListPrice   int `json:"list_price"`
	RetailPrice int `json:"retail_price"`

	CategoryID    dot.ID         `json:"category_id"`
	CollectionIDs []dot.ID       `json:"collectionIDs"`
	Tags          []string       `json:"tags"`
	Status        status3.Status `json:"status"`

	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Deleted   bool      `json:"deleted"`

	ProductType product_type.ProductType `json:"product_type"`
	MetaFields  []*MetaField             `json:"metafeilds"`
	BrandID     dot.ID                   `json:"brand_id"`
}

type ShopVariant struct {
	ExternalID        string `json:"external_id"`
	ExternalCode      string `json:"external_code"`
	ExternalProductID string `json:"external_product_id"`
	PartnerID         dot.ID `json:"partner_id"`

	ShopID    dot.ID `json:"shop_id"`
	ProductID dot.ID `json:"product_id"`
	VariantID dot.ID `json:"variant_id"`
	// variant.code is also known as sku
	Code string `json:"code"`

	CodeNorm    int              `json:"code_norm"`
	Name        string           `json:"name"`
	ShortDesc   string           `json:"short_desc"`
	Description string           `json:"description"`
	DescHTML    string           `json:"descHTML"`
	ImageURLs   []string         `json:"image_urls"`
	Status      status3.Status   `json:"status"`
	Attributes  types.Attributes `json:"attributes"`

	CostPrice   int    `json:"cost_price"`
	ListPrice   int    `json:"list_price"`
	RetailPrice int    `json:"retail_price"`
	Note        string `json:"note"` // only in ShopProduct and ShopVariant

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Deleted   bool      `json:"deleted"`
}

type ShopCategory struct {
	ID        dot.ID
	PartnerID dot.ID

	ParentID dot.ID
	ShopID   dot.ID

	ExternalID       string
	ExternalParentID string

	Name string

	Status int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ShopCollection struct {
	ID        dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	ExternalID string

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time
	UpdatedAt time.Time

	Deleted bool
}

func (v *ShopVariant) GetName() string {
	if len(v.Attributes) == 0 {
		return ""
	}
	return v.Attributes.ShortLabel()
}

type DescriptionInfo struct {
	ShortDesc string

	Description string

	DescHTML string
}

type PriceInfo struct {
	CostPrice int

	ListPrice int

	RetailPrice int
}

//-- extended --//

type ShopVariantWithProduct struct {
	*ShopVariant

	ShopProduct *ShopProduct
}

func (v *ShopVariantWithProduct) GetFullName() string {
	if v.ShopProduct.Name != "" {
		return v.ShopProduct.Name + " - " + v.ShopVariant.GetName()
	}
	return v.ShopVariant.GetName()
}

func (v ShopVariantWithProduct) GetListPrice() int {
	return cmutil.CoalesceInt(
		v.ShopVariant.ListPrice,
		v.ShopProduct.ListPrice,
	)
}

func (v ShopVariantWithProduct) GetRetailPrice() int {
	return cmutil.CoalesceInt(
		v.ShopVariant.RetailPrice, v.ShopVariant.ListPrice,
		v.ShopProduct.RetailPrice, v.ShopProduct.ListPrice,
	)
}

func (v ShopVariantWithProduct) ProductWithVariantName() string {
	productName := v.ShopProduct.Name
	variantLabel := v.ShopVariant.Attributes.Label()
	if variantLabel == "" {
		return productName
	}
	return productName + " - " + variantLabel
}

type ShopProductWithVariants struct {
	*ShopProduct
	Variants []*ShopVariant `json:"variants"`
}

type ShopCategories struct {
	Categories []*ShopCategory
}

type ShopProductCollection struct {
	PartnerID            dot.ID
	ExternalCollectionID string
	ExternalProductID    string

	ProductID    dot.ID
	CollectionID dot.ID
	ShopID       dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time

	Deleted bool
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ShopVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ShopProductDeletedEvent struct {
	ShopID     dot.ID
	ProductIDs []dot.ID
}
