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

func ShopProduct(in *catalogmodel.ShopProduct) (out *catalog.ShopProduct) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopProduct{
		ShopID:        in.ShopID,
		ProductID:     in.ProductID,
		CollectionIDs: in.CollectionIDs,
		Name:          in.Name,
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

func ShopProducts(ins []*catalogmodel.ShopProduct) (outs []*catalog.ShopProduct) {
	outs = make([]*catalog.ShopProduct, len(ins))
	for i, in := range ins {
		outs[i] = ShopProduct(in)
	}
	return outs
}

func ShopProductWithVariants(in *catalogmodel.ShopProductWithVariants) (out *catalog.ShopProductWithVariants) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopProductWithVariants{
		ShopProduct: ShopProduct(in.ShopProduct),
		Variants:    ShopVariants(in.Variants),
	}
	return out
}

func ShopProductsWithVariants(ins []*catalogmodel.ShopProductWithVariants) (outs []*catalog.ShopProductWithVariants) {
	outs = make([]*catalog.ShopProductWithVariants, len(ins))
	for i, in := range ins {
		outs[i] = ShopProductWithVariants(in)
	}
	return outs
}

func ShopVariant(in *catalogmodel.ShopVariant) (out *catalog.ShopVariant) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopVariant{
		ShopID:    in.ShopID,
		ProductID: in.ProductID,
		VariantID: in.VariantID,
		Code:      in.Code,
		Name:      in.Name,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   in.ShortDesc,
			Description: in.Description,
			DescHTML:    in.DescHTML,
		},
		ImageURLs:  in.ImageURLs,
		Status:     int16(in.Status),
		Attributes: Attributes(in.Attributes),
		PriceInfo: catalog.PriceInfo{
			ListPrice:   in.ListPrice,
			CostPrice:   in.CostPrice,
			RetailPrice: in.RetailPrice,
		},
		Note:      in.Note,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
	return out
}

func ShopVariants(ins []*catalogmodel.ShopVariant) (outs []*catalog.ShopVariant) {
	outs = make([]*catalog.ShopVariant, len(ins))
	for i, in := range ins {
		outs[i] = ShopVariant(in)
	}
	return outs
}

func ShopVariantWithProduct(in *catalogmodel.ShopVariantWithProduct) (out *catalog.ShopVariantWithProduct) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopVariantWithProduct{
		ShopProduct: ShopProduct(in.ShopProduct),
		ShopVariant: ShopVariant(in.ShopVariant),
	}
	return out
}

func ShopVariantsWithProduct(ins []*catalogmodel.ShopVariantWithProduct) (outs []*catalog.ShopVariantWithProduct) {
	outs = make([]*catalog.ShopVariantWithProduct, len(ins))
	for i, in := range ins {
		outs[i] = ShopVariantWithProduct(in)
	}
	return outs
}
