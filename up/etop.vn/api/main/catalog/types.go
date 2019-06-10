package catalog

import (
	"time"

	"etop.vn/api/main/catalog/types"
)

type ProductSource struct {
	ID   int64
	Type string
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID int64

	ProductSourceID int64

	Code string

	Name string

	DescriptionInfo

	ImageURLs []string

	Tags []string

	Unit string

	Status int16

	CreatedAt time.Time

	UpdatedAt time.Time

	PriceDeclareInfo
}

type ShopProduct struct {
	ShopID int64

	ProductID int64

	Code string

	Name string

	DescriptionInfo

	ImageURLs []string

	Tags []string

	Status int32

	PriceInfo

	Note string // only in ShopProduct and ShopVariant

	CreatedAt time.Time

	UpdatedAt time.Time
}

type Variant struct {
	ID int64

	ProductID int64

	// Code is also known as sku
	Code string

	Name string

	DescriptionInfo

	ImageURLs []string

	Status int16

	Attributes // only in Variant, not in ShopVariant

	// TODO: price shoule be managed in pricing service
	PriceDeclareInfo

	CreatedAt time.Time

	UpdatedAt time.Time
}

type ShopVariant struct {
	ShopID int64

	ProductID int64

	VariantID int64

	// Code is also known as sku
	Code string

	Name string

	DescriptionInfo

	ImageURLs []string

	Status int16

	PriceInfo

	Note string // only in ShopProduct and ShopVariant

	CreatedAt time.Time

	UpdatedAt time.Time
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
	ListPrice int32

	CostPrice int32

	RetailPrice int32
}

type Attribute = types.Attribute
type Attributes = types.Attributes
