package model

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

// Normalize attributes, do not sort them. Empty attributes is '_'.
func NormalizeAttributes(attrs []*ProductAttribute) ([]*ProductAttribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}

	normAttrs := make([]*ProductAttribute, 0, len(attrs))
	b := make([]byte, 0, 256)
	for _, attr := range attrs {
		attr.Name, _ = validate.NormalizeName(attr.Name)
		attr.Value, _ = validate.NormalizeName(attr.Value)
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		nameNorm := validate.NormalizeUnderscore(attr.Name)
		valueNorm := validate.NormalizeUnderscore(attr.Value)
		if nameNorm == "" || valueNorm == "" {
			continue
		}

		normAttrs = append(normAttrs, &ProductAttribute{Name: attr.Name, Value: attr.Value})
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, nameNorm...)
		b = append(b, '=')
		b = append(b, valueNorm...)
	}
	s := string(b)
	if s == "" {
		s = "_"
	}
	return normAttrs, s
}

var _ = sqlgenShopVariantWithProduct(
	&ShopVariantWithProduct{}, &ShopProduct{}, sq.AS("sp"),
	sq.LEFT_JOIN, &ShopVariant{}, sq.AS("sv"), "sp.product_id = sv.product_id",
)

type ShopVariantWithProduct struct {
	*ShopVariant
	*ShopProduct
}

var _ = sqlgenShopVariant(&ShopVariant{})

// +convert:type=catalog.ShopVariant
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

	CostPrice   int32
	ListPrice   int32
	RetailPrice int32

	Status     model.Status3
	Attributes ProductAttributes

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	NameNorm string

	// key-value normalization, must be non-null. Empty attributes is '_'
	AttrNormKv string
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var _ = sqlgenShopProduct(&ShopProduct{})

// +convert:type=catalog.ShopProduct
type ShopProduct struct {
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

	CostPrice   int32
	ListPrice   int32
	RetailPrice int32
	BrandID     dot.ID

	Status model.Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	NameNorm    string
	NameNormUa  string // unaccent normalization
	ProductType string

	MetaFields []*MetaField
}

type ShopProductWithVariants struct {
	*ShopProduct
	Variants []*ShopVariant
}

var _ = sqlgenProductShopCollection(&ProductShopCollection{})

type ProductShopCollection struct {
	CollectionID dot.ID
	ProductID    dot.ID
	ShopID       dot.ID
	Status       int
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

type ProductAttributes []*ProductAttribute

func (attrs ProductAttributes) Name() string {
	if len(attrs) == 0 {
		return ""
	}
	return attrs.ShortLabel()
}

func (attrs ProductAttributes) Label() string {
	if len(attrs) == 0 {
		return "Mặc định"
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

func (attrs ProductAttributes) ShortLabel() string {
	if len(attrs) == 0 {
		return "Mặc định"
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

// +convert:type=catalog.Attribute
type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var _ = sqlgenShopCategory(&ShopCategory{})

// +convert:type=catalog.ShopCategory
type ShopCategory struct {
	ID dot.ID

	ParentID dot.ID
	ShopID   dot.ID

	Name string

	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

var _ = sqlgenShopCollection(&ShopCollection{})

// +convert:type=catalog.ShopCollection
type ShopCollection struct {
	ID     dot.ID
	ShopID dot.ID

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenShopProductCollection(&ShopProductCollection{})

// +convert:type=catalog.ShopProductCollection
type ShopProductCollection struct {
	ProductID    dot.ID
	CollectionID dot.ID
	ShopID       dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenShopBrand(&ShopBrand{})

// +convert:type=catalog.ShopBrand
type ShopBrand struct {
	ID     dot.ID
	ShopID dot.ID

	BrandName   string
	Description string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

var _ = sqlgenShopSupplierVariant(&ShopVariantSupplier{})

// +convert:type=catalog.ShopVariantSupplier
type ShopVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
}
