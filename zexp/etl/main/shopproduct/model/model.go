package model

import (
	"time"

	"etop.vn/api/top/types/etc/product_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
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
	CodeNorm    int
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

	Status status3.Status

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	NameNorm    string
	NameNormUa  string // unaccent normalization
	ProductType product_type.ProductType

	MetaFields []*MetaField

	Rid dot.ID
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
