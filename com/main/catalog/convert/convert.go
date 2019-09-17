package convert

import (
	"time"

	"etop.vn/api/main/catalog"
	catalogtypes "etop.vn/api/main/catalog/types"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	. "etop.vn/capi/dot"
)

// +gen:convert: etop.vn/backend/com/main/catalog/model->etop.vn/api/main/catalog,etop.vn/api/main/catalog/types
// +gen:convert: etop.vn/api/main/catalog

func AttributeDB(in *catalogtypes.Attribute) (out catalogmodel.ProductAttribute) {
	if in == nil {
		return catalogmodel.ProductAttribute{}
	}
	return catalogmodel.ProductAttribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func AttributesDB(ins []*catalogtypes.Attribute) (outs []catalogmodel.ProductAttribute) {
	outs = make([]catalogmodel.ProductAttribute, len(ins))
	for i, in := range ins {
		outs[i] = AttributeDB(in)
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
			ListPrice:   in.ListPrice,
			CostPrice:   in.CostPrice,
			RetailPrice: in.RetailPrice,
		},
		Status:      int32(in.Status),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		ProductType: catalog.ProductType(in.ProductType),
	}
	return out
}

func ShopProductUpdate(in *catalogmodel.ShopProduct) (out *catalog.UpdateShopProductInfoArgs) {
	if in == nil {
		return nil
	}
	out = &catalog.UpdateShopProductInfoArgs{
		ShopID:      in.ShopID,
		ProductID:   in.ProductID,
		Code:        PString(&in.Code),
		Name:        PString(&in.Name),
		Unit:        PString(&in.Unit),
		Note:        PString(&in.Note),
		ShortDesc:   PString(&in.ShortDesc),
		Description: PString(&in.Description),
		DescHTML:    PString(&in.DescHTML),
		CostPrice:   PInt32(&in.CostPrice),
		ListPrice:   PInt32(&in.ListPrice),
		RetailPrice: PInt32(&in.RetailPrice),
	}
	return out
}

func ShopProductDB(in *catalog.ShopProduct) (out *catalogmodel.ShopProduct) {
	if in == nil {
		return nil
	}
	out = &catalogmodel.ShopProduct{
		ShopID:        in.ShopID,
		ProductID:     in.ProductID,
		CollectionIDs: in.CollectionIDs,
		Code:          in.Code,
		Name:          in.Name,
		Description:   in.Description,
		DescHTML:      in.DescHTML,
		ShortDesc:     in.ShortDesc,
		ImageURLs:     in.ImageURLs,
		Note:          in.Note,
		Tags:          in.Tags,
		Unit:          in.Unit,
		CategoryID:    in.CategoryID,
		CostPrice:     in.CostPrice,
		ListPrice:     in.ListPrice,
		RetailPrice:   in.RetailPrice,
		Status:        model.Status3(in.Status),
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
		DeletedAt:     time.Time{},
		NameNorm:      validate.NormalizeSearch(in.Name),
		NameNormUa:    validate.NormalizeUnaccent(in.Name),
		ProductType:   string(in.ProductType),
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

func ShopVariantDB(in *catalog.ShopVariant) (out *catalogmodel.ShopVariant) {
	if in == nil {
		return nil
	}
	out = &catalogmodel.ShopVariant{
		ShopID:      in.ShopID,
		VariantID:   in.VariantID,
		ProductID:   in.ProductID,
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		DescHTML:    in.DescHTML,
		ShortDesc:   in.ShortDesc,
		ImageURLs:   in.ImageURLs,
		Note:        in.Note,
		Tags:        nil,
		CostPrice:   in.CostPrice,
		ListPrice:   in.ListPrice,
		RetailPrice: in.RetailPrice,
		Status:      0,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		DeletedAt:   time.Time{},
		NameNorm:    validate.NormalizeSearch(in.Name),
	}
	out.Attributes, out.AttrNormKv = catalogmodel.NormalizeAttributes(AttributesDB(in.Attributes))
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

func UpdateShopProduct(in *catalogmodel.ShopProduct, args *catalog.UpdateShopProductInfoArgs) (out *catalogmodel.ShopProduct) {
	if in == nil {
		return nil
	}
	shopProduct := &catalogmodel.ShopProduct{
		ShopID:      args.ShopID,
		ProductID:   args.ProductID,
		Code:        args.Code.Apply(in.Code),
		Name:        args.Name.Apply(in.Name),
		Description: args.Description.Apply(in.Description),
		DescHTML:    args.DescHTML.Apply(in.DescHTML),
		ShortDesc:   args.ShortDesc.Apply(in.ShortDesc),
		Note:        args.Note.Apply(in.Note),
		Unit:        args.Unit.Apply(in.Unit),
		CostPrice:   args.CostPrice.Apply(in.CostPrice),
		ListPrice:   args.ListPrice.Apply(in.ListPrice),
		RetailPrice: args.RetailPrice.Apply(in.RetailPrice),

		CategoryID: in.CategoryID,
		Status:     in.Status,
		DeletedAt:  in.DeletedAt,
		NameNorm:   in.NameNorm,
		NameNormUa: in.NameNormUa,
	}
	return shopProduct
}

func UpdateShopVariant(in *catalogmodel.ShopVariant, args *catalog.UpdateShopVariantInfoArgs) (out *catalogmodel.ShopVariant) {
	if in == nil {
		return nil
	}
	shopVariant := &catalogmodel.ShopVariant{
		ShopID:      args.ShopID,
		VariantID:   args.VariantID,
		Code:        args.Code.Apply(in.Code),
		Name:        args.Name.Apply(in.Name),
		Description: args.Descripttion.Apply(in.Description),
		DescHTML:    args.DescHTML.Apply(in.DescHTML),
		ShortDesc:   args.ShortDesc.Apply(in.ShortDesc),
		Note:        args.Note.Apply(in.Note),
		CostPrice:   args.CostPrice.Apply(in.CostPrice),
		ListPrice:   args.ListPrice.Apply(in.ListPrice),
		RetailPrice: args.RetailPrice.Apply(in.RetailPrice),

		ImageURLs: in.ImageURLs,
		ProductID: in.ProductID,
		Tags:      in.Tags,
		DeletedAt: in.DeletedAt,
	}
	// var attributes catalogmodel.ProductAttributes
	// for _, value := range args.Attributes {
	// 	attributes = append(attributes, catalogmodel.ProductAttribute{
	// 		Name:  value.Name,
	// 		Value: value.Value,
	// 	})
	// // }
	// shopVariant.Attributes = attributes
	return shopVariant
}
