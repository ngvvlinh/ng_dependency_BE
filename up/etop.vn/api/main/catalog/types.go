package catalog

import (
	"time"

	"etop.vn/api/main/catalog/types"
	"etop.vn/api/top/types/etc/product_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	cmutil "etop.vn/capi/util"
)

// +gen:event:topic=event/catalog

type ShopProduct struct {
	ExternalID string

	ExternalCode string

	ExternalBrandID string

	ExternalCategoryID string

	PartnerID dot.ID

	ShopID dot.ID

	ProductID dot.ID

	Code string

	CodeNorm int

	Name string

	Unit string

	ImageURLs []string

	Note string

	ShortDesc string

	Description string

	DescHTML string

	CostPrice int

	ListPrice int

	RetailPrice int

	CategoryID dot.ID

	CollectionIDs []dot.ID

	Tags []string

	Status status3.Status

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time

	ProductType product_type.ProductType

	MetaFields []*MetaField

	BrandID dot.ID

	Deleted bool
}

type ShopVariant struct {
	ExternalID string

	ExternalCode string

	ExternalProductID string

	PartnerID dot.ID

	ShopID dot.ID

	ProductID dot.ID

	VariantID dot.ID

	// variant.code is also known as sku
	Code string

	CodeNorm int

	Name string

	ShortDesc string

	Description string

	DescHTML string

	ImageURLs []string

	Status status3.Status

	Attributes types.Attributes

	CostPrice int

	ListPrice int

	RetailPrice int

	Note string // only in ShopProduct and ShopVariant

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time

	Deleted bool
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
	Variants []*ShopVariant
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
	Key   string
	Value string
}

type ShopVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
