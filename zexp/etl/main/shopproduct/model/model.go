package model

import (
	"time"

	"o.o/api/top/types/etc/product_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type ShopProduct struct {
	ExternalID         string
	ExternalCode       string
	PartnerID          dot.ID
	ExternalBrandID    string
	ExternalCategoryID string

	ShopID        dot.ID
	ProductID     dot.ID
	CollectionIDs []dot.ID `sq:"-"`

	Code        string
	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string
	Unit        string
	CategoryID  dot.ID

	CostPrice   int
	ListPrice   int
	RetailPrice int
	BrandID     dot.ID

	Status status3.Status `sql_type:"int2"`

	CreatedAt time.Time
	UpdatedAt time.Time

	ProductType product_type.ProductType `sql_type:"text"`

	MetaFields []*MetaField

	Rid dot.ID
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
