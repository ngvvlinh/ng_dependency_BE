package catalog

import (
	"time"

	"etop.vn/api/main/catalog/types"
	cmutil "etop.vn/capi/util"
)

type ShopProduct struct {
	ShopID int64

	ProductID int64

	Code string

	Name string

	Unit string

	ImageURLs []string

	Note string

	DescriptionInfo

	PriceInfo

	CategoryID int64

	CollectionIDs []int64

	Tags []string

	Status int32

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time
}

type ShopVariant struct {
	ShopID int64

	ProductID int64

	VariantID int64

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
