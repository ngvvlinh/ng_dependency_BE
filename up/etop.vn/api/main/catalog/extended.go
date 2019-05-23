package catalog

type ShopVariantExtended struct {
	*ShopVariant
	*ShopProduct
	*Variant
	*Product
}

type ShopProductExtended struct {
	*ShopProduct
	*Product
}

type ShopVariants []ShopVariantExtended
type ShopProducts []ShopProductExtended

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
func (v ShopVariantExtended) GetProduct() *Product         { return v.Product }
func (v ShopVariantExtended) GetVariant() *Variant         { return v.Variant }
func (v ShopVariantExtended) GetShopProduct() *ShopProduct { return v.ShopProduct }
func (v ShopVariantExtended) GetShopVariant() *ShopVariant { return v.ShopVariant }
