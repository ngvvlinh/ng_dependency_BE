package model

import (
	"fmt"
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenProduct(&Product{})

type Product struct {
	ID                      int64
	ProductSourceID         int64
	ProductSourceCategoryID int64

	Name        string
	ShortDesc   string
	Description string
	DescHTML    string `sq:"'desc_html'"`
	Unit        string
	ImageURLs   []string `sq:"'image_urls'"`

	Status model.Status3
	Code   string `sq:"'ed_code'"`

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm   string // search normalization
	NameNormUa string // unaccent normalization
}

func (p *Product) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	p.NameNormUa = validate.NormalizeUnaccent(p.Name)
	return nil
}

func (p *Product) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	p.NameNormUa = validate.NormalizeUnaccent(p.Name)
	return nil
}

func (p *Product) GetFullName() string {
	return coalesce(p.Name)
}

var _ = sqlgenProductExtended(
	&ProductExtended{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &ProductSource{}, sq.AS("ps"), "p.product_source_id = ps.id",
)

type ProductExtended struct {
	*Product
	*ProductSource
}

type ProductFtVariant struct {
	*Product
	Variants []*Variant
}

var _ = sqlgenVariant(&Variant{})

type Variant struct {
	ID              int64
	ProductID       int64
	ProductSourceID int64

	Code        string `sq:"'ed_code'"`
	ShortDesc   string
	Description string
	DescHTML    string   `sq:"'desc_html'"`
	ImageURLs   []string `sq:"'image_urls'"`

	Attributes ProductAttributes

	CostPrice int
	ListPrice int

	Status    model.Status3
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	// key-value normalization, must be non-null. Empty attributes is '_'
	AttrNormKv string
}

func (v *Variant) GetName() string {
	if len(v.Attributes) == 0 {
		return ""
	}
	return ProductAttributes(v.Attributes).ShortLabel()
}

func (v *Variant) BeforeInsert() error {
	v.Attributes, v.AttrNormKv = NormalizeAttributes(v.Attributes)
	return nil
}

func (v *Variant) BeforeUpdate() error {
	v.Attributes, v.AttrNormKv = NormalizeAttributes(v.Attributes)
	return nil
}

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
			fmt.Println("1", attr.Name, attr.Value)
			continue
		}
		nameNorm := validate.NormalizeUnderscore(attr.Name)
		valueNorm := validate.NormalizeUnderscore(attr.Value)
		if nameNorm == "" || valueNorm == "" {
			fmt.Println("2")
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

var _ = sqlgenVariantExtended(
	&VariantExtended{}, &Variant{}, sq.AS("v"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "v.product_id = p.id",
)

type VariantExtended struct {
	*Variant
	Product *Product
}

func (v *VariantExtended) GetFullName() string {
	if v.Product.Name != "" {
		return v.Product.Name + " - " + v.GetName()
	}
	return v.GetName()
}

var _ = sqlgenShopVariantExtended(
	&ShopVariantExtended{}, &ShopVariant{}, sq.AS("sv"),
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "sv.variant_id = v.id",
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "v.product_id = p.id",
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ShopVariantExtended struct {
	*ShopVariant
	VariantExtended
	*ShopProduct
}

func (v *ShopVariantExtended) GetFullName() string {
	var productName, variantName string
	if v.ShopProduct != nil && v.ShopProduct.Name != "" {
		productName = v.ShopProduct.Name
	} else {
		productName = v.Product.Name
	}
	if v.ShopVariant != nil && v.ShopVariant.Name != "" {
		variantName = v.ShopVariant.Name
	} else {
		variantName = v.Variant.GetName()
	}
	return productName + " - " + variantName
}

var _ = sqlgenShopVariant(&ShopVariant{})

type ShopVariant struct {
	ShopID       int64
	VariantID    int64
	CollectionID int64
	ProductID    int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string

	RetailPrice int
	Status      model.Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm string
}

func (p *ShopVariant) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

func (p *ShopVariant) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

var _ = sqlgenShopProduct(&ShopProduct{})

type ShopProduct struct {
	ShopID        int64
	ProductID     int64
	CollectionIDs []int64 `sq:"-"`

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Note        string
	Tags        []string

	RetailPrice int
	Status      model.Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	NameNorm string

	ProductSourceID   int64  `sq:"-"`
	ProductSourceName string `sq:"-"`
	ProductSourceType string `sq:"-"`
}

func (p *ShopProduct) BeforeInsert() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

func (p *ShopProduct) BeforeUpdate() error {
	p.NameNorm = validate.NormalizeSearch(p.Name)
	return nil
}

var _ = sqlgenShopProductFtProductFtVariantFtShopVariant(
	&ShopProductFtProductFtVariantFtShopVariant{}, &ShopProduct{}, sq.AS("sp"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "sp.product_id = p.id",
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "sp.product_id = v.product_id",
	sq.LEFT_JOIN, &ShopVariant{}, sq.AS("sv"), "sp.shop_id = sv.shop_id and sv.variant_id = v.id",
)

type ShopProductFtProductFtVariantFtShopVariant struct {
	*ShopProduct
	*Product
	*Variant
	*ShopVariant
}

var _ = sqlgenShopProductExtended(
	&ShopProductExtended{}, &ShopProduct{}, sq.AS("sp"),
	sq.LEFT_JOIN, &Product{}, sq.AS("p"), "sp.product_id = p.id",
)

type ShopProductExtended struct {
	*ShopProduct
	*Product
	Variants []*ShopVariantExt `sq:"-"`
}

var _ = sqlgenShopVariantExt(
	&ShopVariantExt{}, &ShopVariant{}, sq.AS("sv"),
	sq.JOIN, &Variant{}, sq.AS("v"), "sp.product_id = v.product_id",
)

type ShopVariantExt struct {
	*ShopVariant
	*Variant
}

var _ = sqlgenProductFtVariantFtShopProduct(
	&ProductFtVariantFtShopProduct{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "v.product_id = p.id",
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ProductFtVariantFtShopProduct struct {
	*Product
	*Variant
	*ShopProduct
}

var _ = sqlgenProductFtShopProduct(
	&ProductFtShopProduct{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ProductFtShopProduct struct {
	*Product
	*ShopProduct
}

type ShopProductFtVariant struct {
	*ShopProduct
	*Product
	Variants []*ShopVariantExtended
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

var _ = sqlgenProductSource(&ProductSource{})

const ProductSourceCustom = "custom"

type ProductSource struct {
	ID     int64
	Type   string
	Name   string
	Status model.Status3

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenProductSourceCategory(&ProductSourceCategory{})

type ProductSourceCategory struct {
	ID int64

	ProductSourceID   int64
	ProductSourceType string
	ParentID          int64
	ShopID            int64

	Name string

	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}
