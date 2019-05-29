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
		Code:      in.EdCode,
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
		Status:      int16(in.Status),
		Code:        in.EdCode,
		ListPrice:   in.ListPrice,
		CostPrice:   in.CostPrice,
		RetailPrice: 0,
		Attributes:  Attributes(in.Attributes),
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
