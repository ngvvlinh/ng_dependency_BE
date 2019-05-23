package catalog

import (
	"time"

	v1types "etop.vn/api/main/catalog/v1/types"
)

type ProductSource struct {
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

type Attribute = v1types.Attribute
type Attributes []Attribute

func (attrs Attributes) Name() string {
	if len(attrs) == 0 {
		return ""
	}
	return attrs.ShortLabel()
}

func (attrs Attributes) Label() string {
	if len(attrs) == 0 {
		return ""
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ", "...)
		}
		b = append(b, attr.Name...)
		b = append(b, ": "...)
		b = append(b, attr.Value...)
	}
	return string(b)
}

func (attrs Attributes) LabelWithDefault(s string) string {
	if len(attrs) == 0 {
		return s
	}
	return attrs.Label()
}

func (attrs Attributes) ShortLabel() string {
	if len(attrs) == 0 {
		return ""
	}
	b := make([]byte, 0, 64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, attr.Value...)
	}
	return string(b)
}

func (attrs Attributes) ShortLabelWithDefault(s string) string {
	if len(attrs) == 0 {
		return s
	}
	return attrs.ShortLabel()
}
