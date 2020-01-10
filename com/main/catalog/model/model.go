package model

import (
	"strings"
	"time"

	"etop.vn/api/main/catalog/types"
	"etop.vn/api/top/types/etc/product_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

// Normalize attributes, do not sort them. Empty attributes is '_'.
func NormalizeAttributes(attrs []*types.Attribute) ([]*types.Attribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}

	normAttrs := make([]*types.Attribute, 0, len(attrs))
	var b strings.Builder
	b.Grow(256)
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

		normAttrs = append(normAttrs, &types.Attribute{Name: attr.Name, Value: attr.Value})
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(nameNorm)
		b.WriteByte('=')
		b.WriteString(valueNorm)
	}
	s := b.String()
	if s == "" {
		s = "_"
	}
	return normAttrs, s
}

var _ = sqlgenShopVariantWithProduct(
	&ShopVariantWithProduct{}, &ShopVariant{}, "sv",
	sq.LEFT_JOIN, &ShopProduct{}, "sp",
	"sp.product_id = sv.product_id",
)

type ShopVariantWithProduct struct {
	*ShopVariant
	*ShopProduct
}

var _ = sqlgenShopVariant(&ShopVariant{})

// +convert:type=catalog.ShopVariant
type ShopVariant struct {
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

	ShopID    dot.ID
	VariantID dot.ID `paging:"id"`
	ProductID dot.ID

	Code        string
	CodeNorm    int
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

	Status     status3.Status
	Attributes ProductAttributes

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update" paging:"updated_at"`
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
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

	ShopID        dot.ID
	ProductID     dot.ID   `paging:"id"`
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

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update" paging:"updated_at"`
	DeletedAt time.Time

	NameNorm    string
	NameNormUa  string // unaccent normalization
	ProductType product_type.ProductType

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
	var b strings.Builder
	b.Grow(64)
	for _, attr := range attrs {
		if attr.Name == "" || attr.Value == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(attr.Value)
	}
	return b.String()
}

// +convert:type=catalog/types.Attribute
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
