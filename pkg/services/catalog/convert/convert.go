package convert

import (
	"etop.vn/api/main/catalog"
	catalogtypes "etop.vn/api/main/catalog/types"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func AttributeToModel(in *catalogtypes.Attribute) (out catalogmodel.ProductAttribute) {
	if in == nil {
		return catalogmodel.ProductAttribute{}
	}
	return catalogmodel.ProductAttribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func AttributesToModel(ins []*catalogtypes.Attribute) (outs []catalogmodel.ProductAttribute) {
	outs = make([]catalogmodel.ProductAttribute, len(ins))
	for i, in := range ins {
		outs[i] = AttributeToModel(in)
	}
	return outs
}

func Attribute(in catalogmodel.ProductAttribute) (out *catalogtypes.Attribute) {
	return &catalogtypes.Attribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func Attributes(ins []catalogmodel.ProductAttribute) (outs catalogtypes.Attributes) {
	outs = make(catalogtypes.Attributes, len(ins))
	for i, in := range ins {
		outs[i] = Attribute(in)
	}
	return outs
}

func Product(in *catalogmodel.Product) (out *catalog.Product) {
	if in == nil {
		return nil
	}
	out = &catalog.Product{
		ID:              in.ID,
		ProductSourceID: in.ProductSourceID,
		Name:            in.Name,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   in.ShortDesc,
			Description: in.Description,
			DescHTML:    in.DescHTML,
		},
		Unit:      in.Unit,
		Status:    int16(in.Status),
		Code:      in.Code,
		ImageURLs: in.ImageURLs,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
	return out
}

func Products(ins []*catalogmodel.Product) (outs []*catalog.Product) {
	outs = make([]*catalog.Product, len(ins))
	for i, in := range ins {
		outs[i] = Product(in)
	}
	return outs
}

func ProductWithVariants(in *catalogmodel.ProductFtVariant) (out *catalog.ProductWithVariants) {
	if in == nil {
		return nil
	}
	out = &catalog.ProductWithVariants{
		Product:  Product(in.Product),
		Variants: Variants(in.Variants),
	}
	return out
}

func ProductsWithVariants(ins []*catalogmodel.ProductFtVariant) (outs []*catalog.ProductWithVariants) {
	outs = make([]*catalog.ProductWithVariants, len(ins))
	for i, in := range ins {
		outs[i] = ProductWithVariants(in)
	}
	return outs
}

func Variant(in *catalogmodel.Variant) (out *catalog.Variant) {
	if in == nil {
		return nil
	}
	out = &catalog.Variant{
		ID:        in.ID,
		ProductID: in.ProductID,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   in.ShortDesc,
			Description: in.Description,
			DescHTML:    in.DescHTML,
		},
		Status: int16(in.Status),
		Code:   in.Code,
		PriceDeclareInfo: catalog.PriceDeclareInfo{
			ListPrice:   in.ListPrice,
			CostPrice:   in.CostPrice,
			RetailPrice: 0,
		},
		Attributes: Attributes(in.Attributes),
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}
	return out
}

func Variants(ins []*catalogmodel.Variant) (outs []*catalog.Variant) {
	outs = make([]*catalog.Variant, len(ins))
	for i, in := range ins {
		outs[i] = Variant(in)
	}
	return outs
}

func VariantWithProduct(in *catalogmodel.VariantExtended) (out *catalog.VariantWithProduct) {
	if in == nil {
		return nil
	}
	out = &catalog.VariantWithProduct{
		Variant: Variant(in.Variant),
		Product: Product(in.Product),
	}
	return out
}

func VariantsWithProduct(ins []*catalogmodel.VariantExtended) (outs []*catalog.VariantWithProduct) {
	outs = make([]*catalog.VariantWithProduct, len(ins))
	for i, in := range ins {
		outs[i] = VariantWithProduct(in)
	}
	return outs
}

func ShopProduct(in *catalogmodel.ShopProduct) (out *catalog.ShopProduct) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopProduct{
		ShopID:    in.ShopID,
		ProductID: in.ProductID,
		Name:      in.Name,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   in.ShortDesc,
			Description: in.Description,
			DescHTML:    in.DescHTML,
		},
		ImageURLs: in.ImageURLs,
		Note:      in.Note,
		Tags:      in.Tags,
		PriceInfo: catalog.PriceInfo{
			ListPrice:   0,
			CostPrice:   0,
			RetailPrice: in.RetailPrice,
		},
		Status:    int32(in.Status),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
	return out
}

func ShopProductExtended(in *catalogmodel.ShopProductExtended) (out *catalog.ShopProductExtended) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopProductExtended{
		ShopProduct: ShopProduct(in.ShopProduct),
		Product:     Product(in.Product),
	}
	return out
}

func ShopProductExtendeds(ins []*catalogmodel.ShopProductExtended) (outs []*catalog.ShopProductExtended) {
	outs = make([]*catalog.ShopProductExtended, len(ins))
	for i, in := range ins {
		outs[i] = ShopProductExtended(in)
	}
	return outs
}

func ShopProductWithVariants(in *catalogmodel.ShopProductFtVariant) (out *catalog.ShopProductWithVariants) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopProductWithVariants{
		Product:     Product(in.Product),
		ShopProduct: ShopProduct(in.ShopProduct),
		Variants:    ShopVariantExtendeds(in.Variants),
	}
	return out
}

func ShopProductsWithVariants(ins []*catalogmodel.ShopProductFtVariant) (outs []*catalog.ShopProductWithVariants) {
	outs = make([]*catalog.ShopProductWithVariants, len(ins))
	for i, in := range ins {
		outs[i] = ShopProductWithVariants(in)
	}
	return outs
}

func ShopVariantExtended(in *catalogmodel.ShopVariantExtended) (out *catalog.ShopVariantExtended) {
	if in == nil {
		return nil
	}
	sv := &catalog.ShopVariant{
		ShopID:    in.ShopVariant.ShopID,
		VariantID: in.ShopVariant.VariantID,
		Code:      in.Variant.Code,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   in.ShopVariant.ShortDesc,
			Description: in.ShopVariant.Description,
			DescHTML:    in.ShopVariant.DescHTML,
		},
		Status:    int16(in.ShopVariant.Status),
		CreatedAt: in.ShopVariant.CreatedAt,
		UpdatedAt: in.ShopVariant.UpdatedAt,
	}
	out = &catalog.ShopVariantExtended{
		ShopVariant: sv,
		Variant:     Variant(in.Variant),
	}
	return out
}

func ShopVariantExtendeds(ins []*catalogmodel.ShopVariantExtended) (outs []*catalog.ShopVariantExtended) {
	outs = make([]*catalog.ShopVariantExtended, len(ins))
	for i, in := range ins {
		outs[i] = ShopVariantExtended(in)
	}
	return outs
}
