package catalog

type VariantWithProduct struct {
	*Variant
	*Product
}

type ShopVariantWithProduct struct {
	*Product
	*Variant
	*ShopProduct
	*ShopVariant
}

type ShopVariantExtended struct {
	*ShopVariant
	*Variant
}

type ShopProductExtended struct {
	*ShopProduct
	*Product
}

type ProductWithVariants struct {
	*Product
	Variants []*Variant
}

type ShopProductWithVariants struct {
	*Product
	*ShopProduct
	Variants []*ShopVariantExtended
}

type ProductInterface interface{ GetProduct() *Product }
type VariantInterface interface{ GetVariant() *Variant }
type ShopProductInterface interface{ GetShopProduct() *ShopProduct }
type ShopVariantInterface interface{ GetShopVariant() *ShopVariant }

type VariantExtendedInterface interface {
	ProductInterface
	VariantInterface
	ShopProductInterface
	ShopVariantInterface
}

func (p *Product) GetProduct() *Product                    { return p }
func (p *ShopProduct) GetShopProduct() *ShopProduct        { return p }
func (v *Variant) GetVariant() *Variant                    { return v }
func (v *ShopVariant) GetShopVariant() *ShopVariant        { return v }
func (p ShopProductExtended) GetProduct() *Product         { return p.Product }
func (p ShopProductExtended) GetShopProduct() *ShopProduct { return p.ShopProduct }
func (v ShopVariantExtended) GetVariant() *Variant         { return v.Variant }
func (v ShopVariantExtended) GetShopVariant() *ShopVariant { return v.ShopVariant }
