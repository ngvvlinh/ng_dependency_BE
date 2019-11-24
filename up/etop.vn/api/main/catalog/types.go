package catalog

import (
	"time"

	"etop.vn/api/main/catalog/types"
	dot "etop.vn/capi/dot"
	cmutil "etop.vn/capi/util"
)

// +gen:event:topic=event/catalog

type ProductType = string

const (
	ProductTypeService ProductType = "services"
	ProductTypeGoods   ProductType = "goods"
)

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

	Status int32

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time

	ProductType ProductType

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

	Status int16

	Attributes Attributes

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
	ListPrice int32

	CostPrice int32

	RetailPrice int32
}

type PriceInfo struct {
	CostPrice int32

	ListPrice int32

	RetailPrice int32
}

type Attribute = types.Attribute
type Attributes = types.Attributes

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

func (v ShopVariantWithProduct) GetListPrice() int32 {
	return cmutil.CoalesceInt32(
		v.ShopVariant.ListPrice,
		v.ShopProduct.ListPrice,
	)
}

func (v ShopVariantWithProduct) GetRetailPrice() int32 {
	return cmutil.CoalesceInt32(
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
