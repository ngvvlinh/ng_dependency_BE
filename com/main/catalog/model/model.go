package model

import (
	"sort"
	"strings"
	"time"

	"o.o/api/main/catalog/types"
	"o.o/api/top/types/etc/product_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

// Normalize attributes. Empty attributes is '_'.
func NormalizeAttributes(attrs []*types.Attribute) ([]*types.Attribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}
	normAttrs := make([]*types.Attribute, 0, len(attrs))
	normAttrKv := make([]*types.Attribute, 0, len(attrs))
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
		normAttrKv = append(normAttrKv, &types.Attribute{Name: nameNorm, Value: valueNorm})
	}
	sort.Slice(normAttrKv, func(i, j int) bool {
		return normAttrKv[i].Name > normAttrKv[j].Name
	})
	for _, v := range normAttrKv {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(v.Name)
		b.WriteByte('=')
		b.WriteString(v.Value)
	}

	s := b.String()
	if s == "" {
		s = "_"
	}
	return normAttrs, s
}

// +sqlgen:           ShopVariant as sv
// +sqlgen:left-join: ShopProduct as sp on sp.product_id = sv.product_id
type ShopVariantWithProduct struct {
	*ShopVariant
	*ShopProduct
}

// +convert:type=catalog.ShopVariant
// +sqlgen
type ShopVariant struct {
	ExternalID        string
	ExternalCode      string
	ExternalProductID string
	PartnerID         dot.ID

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

	Rid dot.ID
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// +convert:type=catalog.ShopProduct
// +sqlgen
type ShopProduct struct {
	ExternalID         string
	ExternalCode       string
	PartnerID          dot.ID
	ExternalBrandID    string
	ExternalCategoryID string

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

	Rid dot.ID
}

type ShopProductWithVariants struct {
	*ShopProduct
	Variants []*ShopVariant
}

// +sqlgen
type ProductShopCollection struct {
	CollectionID dot.ID
	ProductID    dot.ID
	ShopID       dot.ID
	Status       int
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	Rid          dot.ID
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

// +convert:type=catalog.ShopCategory
// +sqlgen
type ShopCategory struct {
	ID        dot.ID
	PartnerID dot.ID
	ShopID    dot.ID

	ExternalID       string
	ExternalParentID string

	ParentID dot.ID

	Name string

	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Rid dot.ID
}

// +convert:type=catalog.ShopCollection
// +sqlgen
type ShopCollection struct {
	ID        dot.ID `paging:"id"`
	ShopID    dot.ID
	PartnerID dot.ID

	ExternalID string

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"  paging:"updated_at"`
	DeletedAt time.Time

	Rid dot.ID
}

// +convert:type=catalog.ShopProductCollection
// +sqlgen
type ShopProductCollection struct {
	PartnerID            dot.ID
	ShopID               dot.ID
	ExternalCollectionID string
	ExternalProductID    string

	ProductID    dot.ID
	CollectionID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	Rid dot.ID
}

// +convert:type=catalog.ShopBrand
// +sqlgen
type ShopBrand struct {
	ID         dot.ID
	ShopID     dot.ID
	ExternalID string
	PartnerID  dot.ID

	BrandName   string
	Description string

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	Rid dot.ID
}

// +convert:type=catalog.ShopVariantSupplier
// +sqlgen
type ShopVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`

	Rid dot.ID
}
