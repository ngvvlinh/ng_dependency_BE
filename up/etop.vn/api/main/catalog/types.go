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
	ShopID dot.ID

	ProductID dot.ID

	Code string

	Name string

	Unit string

	ImageURLs []string

	Note string

	DescriptionInfo

	PriceInfo

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
}

type ShopVariant struct {
	ShopID dot.ID

	ProductID dot.ID

	VariantID dot.ID

	// variant.code is also known as sku
	Code string

	Name string

	DescriptionInfo

	ImageURLs []string

	Status status3.Status

	Attributes types.Attributes

	PriceInfo

	Note string // only in ShopProduct and ShopVariant

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time
}

type ShopCategory struct {
	ID dot.ID

	ParentID dot.ID
	ShopID   dot.ID

	Name string

	Status int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ShopCollection struct {
	ID     dot.ID
	ShopID dot.ID

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time
	UpdatedAt time.Time
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

type PriceDeclareInfo struct {
	ListPrice int

	CostPrice int

	RetailPrice int
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
	ProductID    dot.ID
	CollectionID dot.ID
	ShopID       dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time
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
