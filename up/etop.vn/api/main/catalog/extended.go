package catalog

import cmutil "etop.vn/common/util"

type VariantWithProduct struct {
	Variant *Variant
	Product *Product
}

type ShopVariantWithProduct struct {
	Product     *Product
	Variant     *Variant
	ShopProduct *ShopProduct
	ShopVariant *ShopVariant
}

func (v ShopVariantWithProduct) GetListPrice() int32 {
	return cmutil.CoalesceInt32(
		v.ShopVariant.ListPrice, v.Variant.ListPrice,
		v.ShopProduct.ListPrice, v.Product.ListPrice,
	)
}

func (v ShopVariantWithProduct) GetRetailPrice() int32 {
	return cmutil.CoalesceInt32(
		v.ShopVariant.RetailPrice, v.ShopVariant.ListPrice, v.Variant.ListPrice,
		v.ShopProduct.RetailPrice, v.ShopProduct.ListPrice, v.Product.ListPrice,
	)
}

func (v ShopVariantWithProduct) ProductWithVariantName() string {
	productName := ShopProductExtended{
		Product:     v.Product,
		ShopProduct: v.ShopProduct,
	}.ProductName()

	variantLabel := v.Variant.Attributes.Label()
	if variantLabel == "" {
		return productName
	}
	return productName + " - " + variantLabel
}

type ShopVariantExtended struct {
	*ShopVariant
	*Variant
}

type ShopProductExtended struct {
	*Product
	*ShopProduct
}

func (p ShopProductExtended) ProductName() string {
	return cmutil.CoalesceString(p.ShopProduct.Name, p.Product.Name)
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
