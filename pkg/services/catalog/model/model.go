package model

import (
	"encoding/json"
	"fmt"
	"time"

	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenProduct(&Product{})

type Product struct {
	ID                      int64
	ProductSourceID         int64
	SupplierID              int64
	ProductSourceCategoryID string
	EtopCategoryID          int64

	Name          string
	ShortDesc     string
	Description   string
	DescHTML      string `sq:"'desc_html'"`
	EdName        string
	EdShortDesc   string
	EdDescription string
	EdDescHTML    string `sq:"'ed_desc_html'"`
	EdTags        []string
	Unit          string

	Status model.Status3
	Code   string
	EdCode string

	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int

	ImageURLs []string `sq:"'image_urls'"`

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

func (p *Product) IsAvailable() bool {
	return p.QuantityAvailable > 0
}

func (p *Product) GetFullName() string {
	return coalesce(p.Name, p.EdName)
}

var _ = sqlgenProductExternal(&ProductExternal{})

type ProductExternal struct {
	ID int64

	ProductExternalCommon `sq:"inline"`

	ExternalUnits []*model.Unit
}

type ProductExternalCommon struct {
	ProductSourceID   int64
	ProductSourceType string

	ExternalID          string
	ExternalName        string
	ExternalCode        string
	ExternalCategoryID  string
	ExternalDescription string
	ExternalImageURLs   []string
	ExternalUnit        string

	ExternalData      json.RawMessage
	ExternalStatus    model.Status3
	ExternalCreatedAt time.Time
	ExternalUpdatedAt time.Time
	ExternalDeletedAt time.Time
	LastSyncAt        time.Time
}

var _ = sqlgenProductExtended(
	&ProductExtended{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &ProductExternal{}, sq.AS("px"), "p.id = px.id",
	sq.LEFT_JOIN, &model.ProductSource{}, sq.AS("ps"), "p.product_source_id = ps.id",
)

type ProductExtended struct {
	*Product
	*ProductExternal
	*model.ProductSource
}

type ProductFtVariant struct {
	ProductExtended
	Variants []*VariantExternalExtended
}

var _ = sqlgenVariant(&Variant{})

type Variant struct {
	ID              int64
	ProductID       int64
	ProductSourceID int64
	// ProductSourceType string
	SupplierID int64

	ProductSourceCategoryID int64
	EtopCategoryID          int64

	// Name          string
	ShortDesc     string
	Description   string
	DescHTML      string `sq:"'desc_html'"`
	EdName        string
	EdShortDesc   string
	EdDescription string
	EdDescHTML    string `sq:"'ed_desc_html'"`

	DescNorm string

	// key-value normalization, must be non-null. Empty attributes is '_'.
	AttrNormKv string

	Status     model.Status3
	EtopStatus model.Status3
	EdStatus   model.Status3
	Code       string
	EdCode     string

	WholesalePrice0 int `sq:"'wholesale_price_0'"`
	WholesalePrice  int
	ListPrice       int
	RetailPriceMin  int
	RetailPriceMax  int

	EdWholesalePrice0 int `sq:"'ed_wholesale_price_0'"`
	EdWholesalePrice  int
	EdListPrice       int
	EdRetailPriceMin  int
	EdRetailPriceMax  int

	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
	ImageURLs         []string `sq:"'image_urls'"`
	SupplierMeta      json.RawMessage

	CostPrice  int
	Attributes model.ProductAttributes

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

func (v *Variant) GetName() string {
	if len(v.Attributes) == 0 {
		return ""
	}
	return model.ProductAttributes(v.Attributes).ShortLabel()
}

func (v *Variant) IsAvailable() bool {
	return v.QuantityAvailable > 0
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
func NormalizeAttributes(attrs []model.ProductAttribute) ([]model.ProductAttribute, string) {
	if len(attrs) == 0 {
		return nil, "_"
	}
	const maxAttrs = 8
	if len(attrs) > maxAttrs {
		attrs = attrs[:maxAttrs]
	}

	normAttrs := make([]model.ProductAttribute, 0, len(attrs))
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

		normAttrs = append(normAttrs, model.ProductAttribute{Name: attr.Name, Value: attr.Value})
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
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "v.id = vx.id",
)

type VariantExtended struct {
	*Variant
	Product         *Product
	VariantExternal *VariantExternal
}

func (v *VariantExtended) GetFullName() string {
	if v.Product.Name != "" {
		return v.Product.Name + " - " + v.GetName()
	}
	return v.GetName()
}

var _ = sqlgenVariantExternal(&VariantExternal{})

type VariantExternal struct {
	ID int64

	ProductExternalCommon `sq:"inline"`

	ExternalProductID  string
	ExternalPrice      int
	ExternalBaseUnitID string
	ExternalUnitConv   float64
	ExternalAttributes []model.ProductAttribute
}

var _ = sqlgenVariantExternalExtended(
	&VariantExternalExtended{}, &Variant{}, sq.AS("v"),
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "v.id = vx.id",
)

type VariantExternalExtended struct {
	*Variant
	*VariantExternal
}

var _ = sqlgenPriceDef(&PriceDef{}, &Variant{})

type PriceDef struct {
	WholesalePrice0 int `sq:"'wholesale_price_0'"`
	WholesalePrice  int
	ListPrice       int
	RetailPriceMin  int
	RetailPriceMax  int
}

func (p *PriceDef) IsValid() bool {
	return p.ListPrice > 0 &&
		p.WholesalePrice > 0 && p.WholesalePrice0 > 0 &&
		p.RetailPriceMin > 0 && p.RetailPriceMax > 0
}

func (p *PriceDef) ApplyTo(v *Variant) {
	v.ListPrice = p.ListPrice
	v.WholesalePrice0 = p.WholesalePrice0
	v.WholesalePrice = p.WholesalePrice
	v.RetailPriceMin = p.RetailPriceMin
	v.RetailPriceMax = p.RetailPriceMax
}

var _ = sqlgenVariantQuantity(&VariantQuantity{}, &Variant{})

type VariantQuantity struct {
	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
}

var _ = substructEtopProduct(&EtopProduct{}, &Product{})

type EtopProduct struct {
	ID         int64
	SupplierID int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
	ImageURLs   []string `sq:"'image_urls'"`
	Status      model.Status3
	Code        string

	// Only default branch
	QuantityAvailable int // Available = OnHand - Reserved - 3
	QuantityOnHand    int
	QuantityReserved  int
}

func (p *EtopProduct) IsAvailable() bool {
	return p.QuantityAvailable > 0
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
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
	sq.LEFT_JOIN, &ShopVariant{}, sq.AS("sv"), "sp.shop_id = sv.shop_id and sv.variant_id = v.id",
)

type ShopProductFtProductFtVariantFtShopVariant struct {
	*ShopProduct
	*Product
	*Variant
	*VariantExternal
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
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
)

type ShopVariantExt struct {
	*ShopVariant
	*Variant
	*VariantExternal
}

var _ = sqlgenProductFtVariantFtShopProduct(
	&ProductFtVariantFtShopProduct{}, &Product{}, sq.AS("p"),
	sq.LEFT_JOIN, &Variant{}, sq.AS("v"), "v.product_id = p.id",
	sq.LEFT_JOIN, &VariantExternal{}, sq.AS("vx"), "vx.id = v.id",
	sq.LEFT_JOIN, &ShopProduct{}, sq.AS("sp"), "sp.product_id = p.id",
)

type ProductFtVariantFtShopProduct struct {
	*Product
	*Variant
	*VariantExternal
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
