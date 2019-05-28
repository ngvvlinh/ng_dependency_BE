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

	Name string

	DescriptionInfo

	Tags []string

	Unit string

	Status int16

	Code string

	ImageURLs []string

	CreatedAt time.Time

	UpdatedAt time.Time
}

type ShopProduct struct {
	ShopID int64

	ProductID int64

	Name string

	DescriptionInfo

	ImageURLs []string

	Note string

	Tags []string

	RetailPrice int

	Status int32

	CreatedAt time.Time

	DeletedAt time.Time
}

type Variant struct {
	ID int64

	ProductID int64

	DescriptionInfo

	Status int16

	// editable code
	Code string

	// TODO: price shoule be managed in pricing service

	ListPrice int

	CostPrice int

	RetailPrice int

	Attributes
}

type ShopVariant struct {
	ShopID int64

	VariantID int64

	// also known as sku

	Code string

	DescriptionInfo

	Status int16
}

type DescriptionInfo struct {
	ShortDesc string

	Description string

	DescHTML string
}

type Attribute = types.Attribute
type Attributes = types.Attributes
