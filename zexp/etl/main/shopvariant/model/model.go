package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type ShopVariant struct {
	ShopID    dot.ID
	VariantID dot.ID
	ProductID dot.ID

	Code        string
	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string

	CostPrice   int
	ListPrice   int
	RetailPrice int

	Status     status3.Status `sql_type:"int2"`
	Attributes ProductAttributes

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}

type ProductAttributes []*ProductAttribute

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
