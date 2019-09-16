package model

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

// Normalize attributes, do not sort them. Empty attributes is '_'.
func NormalizeAttributes(attrs []ProductAttribute) ([]ProductAttribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}

	normAttrs := make([]ProductAttribute, 0, len(attrs))
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

		normAttrs = append(normAttrs, ProductAttribute{Name: attr.Name, Value: attr.Value})
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

type ShopVariant struct {
	ShopID    int64
	VariantID int64
	ProductID int64

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

var _ = sqlgenShopProduct(&ShopProduct{})

type ShopProduct struct {
	ShopID        int64
	ProductID     int64
	CollectionIDs []int64 `sq:"-"`

	Code        string
	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string
	Unit        string
	CategoryID  int64

	CostPrice   int32
	ListPrice   int32
	RetailPrice int32

	Status model.Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	NameNorm   string
	NameNormUa string // unaccent normalization
}

type ShopProductWithVariants struct {
	*ShopProduct
	Variants []*ShopVariant
}

var _ = sqlgenShopCollection(&ShopCollection{})

type ShopCollection struct {
	ID     int64
	ShopID int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenProductShopCollection(&ProductShopCollection{})

type ProductShopCollection struct {
	CollectionID int64
	ProductID    int64
	ShopID       int64
	Status       int
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

type ProductAttributes []ProductAttribute

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

type ProductAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var _ = sqlgenShopCategory(&ShopCategory{})

type ShopCategory struct {
	ID int64

	ParentID int64
	ShopID   int64

	Name string

	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}
