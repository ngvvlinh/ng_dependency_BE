package convert

import (
	"time"

	"etop.vn/api/main/catalog"
	catalogtypes "etop.vn/api/main/catalog/types"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	. "etop.vn/capi/dot"
)

// +gen:convert: etop.vn/backend/com/main/catalog/model->etop.vn/api/main/catalog,etop.vn/api/main/catalog/types
// +gen:convert: etop.vn/api/main/catalog

func AttributeDB(in *catalogtypes.Attribute, out *catalogmodel.ProductAttribute) {
	convert_catalogtypes_Attribute_catalogmodel_ProductAttribute(in, out)
}

func AttributesDB(ins []*catalogtypes.Attribute) (outs []*catalogmodel.ProductAttribute) {
	return Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(ins)
}

func Attribute(in *catalogmodel.ProductAttribute) (out *catalogtypes.Attribute) {
	return &catalogtypes.Attribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func Attributes(ins []*catalogmodel.ProductAttribute) (outs catalogtypes.Attributes) {
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
	metaFields := []*catalog.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}

	out = &catalog.ShopProduct{
		ShopID:        in.ShopID,
		ProductID:     in.ProductID,
		CollectionIDs: in.CollectionIDs,
		Name:          in.Name,
		Code:          in.Code,
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
		CategoryID:  in.CategoryID,
		VendorID:    in.VendorID,
		ProductType: catalog.ProductType(in.ProductType),
		MetaFields:  metaFields,
		BrandID:     in.BrandID,
	}
	return out
}

func ShopCategory(in *catalogmodel.ShopCategory, out *catalog.ShopCategory) {
	convert_catalogmodel_ShopCategory_catalog_ShopCategory(in, out)
}

func ShopCategoryDB(in *catalog.ShopCategory, out *catalogmodel.ShopCategory) {
	convert_catalog_ShopCategory_catalogmodel_ShopCategory(in, out)
}

func ShopCategories(ins []*catalogmodel.ShopCategory) (outs []*catalog.ShopCategory) {
	return Convert_catalogmodel_ShopCategories_catalog_ShopCategories(ins)
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
	metaFields := []*catalogmodel.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalogmodel.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
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
		VendorID:      in.VendorID,
		CostPrice:     in.CostPrice,
		ListPrice:     in.ListPrice,
		RetailPrice:   in.RetailPrice,
		Status:        model.Status3(in.Status),
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,

		DeletedAt:   time.Time{},
		NameNorm:    validate.NormalizeSearch(in.Name),
		NameNormUa:  validate.NormalizeUnaccent(in.Name),
		ProductType: string(in.ProductType),
		MetaFields:  metaFields,
		BrandID:     in.BrandID,
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
	shopProduct := in
	shopProduct.VendorID = args.VendorID
	shopProduct.Code = args.Code.Apply(in.Code)
	shopProduct.Name = args.Name.Apply(in.Name)
	shopProduct.Description = args.Description.Apply(in.Description)
	shopProduct.DescHTML = args.DescHTML.Apply(in.DescHTML)
	shopProduct.ShortDesc = args.ShortDesc.Apply(in.ShortDesc)
	shopProduct.Note = args.Note.Apply(in.Note)
	shopProduct.Unit = args.Unit.Apply(in.Unit)
	shopProduct.CostPrice = args.CostPrice.Apply(in.CostPrice)
	shopProduct.ListPrice = args.ListPrice.Apply(in.ListPrice)
	shopProduct.RetailPrice = args.RetailPrice.Apply(in.RetailPrice)
	shopProduct.BrandID = args.BrandID.Apply(in.BrandID)
	return shopProduct
}
func UpdateShopCategory(in *catalogmodel.ShopCategory, args *catalog.UpdateShopCategoryArgs) (out *catalogmodel.ShopCategory) {
	if in == nil {
		return nil
	}
	shopCategory := &catalogmodel.ShopCategory{
		ID:       args.ID,
		ShopID:   args.ShopID,
		ParentID: args.ParentID,
		Name:     args.Name.Apply(in.Name),

		Status:    in.Status,
		DeletedAt: in.DeletedAt,
	}
	return shopCategory
}

func UpdateShopVariant(in *catalogmodel.ShopVariant, args *catalog.UpdateShopVariantInfoArgs) (out *catalogmodel.ShopVariant) {
	if in == nil {
		return nil
	}
	shopVariant := in
	shopVariant.Code = args.Code.Apply(in.Code)
	shopVariant.Name = args.Name.Apply(in.Name)
	shopVariant.Description = args.Descripttion.Apply(in.Description)
	shopVariant.DescHTML = args.DescHTML.Apply(in.DescHTML)
	shopVariant.ShortDesc = args.ShortDesc.Apply(in.ShortDesc)
	shopVariant.Note = args.Note.Apply(in.Note)
	shopVariant.CostPrice = args.CostPrice.Apply(in.CostPrice)
	shopVariant.ListPrice = args.ListPrice.Apply(in.ListPrice)
	shopVariant.RetailPrice = args.RetailPrice.Apply(in.RetailPrice)
	return shopVariant
}

func ShopCollectionDB(in *catalog.ShopCollection, out *catalogmodel.ShopCollection) {
	convert_catalog_ShopCollection_catalogmodel_ShopCollection(in, out)
}

func ShopCollection(in *catalogmodel.ShopCollection, out *catalog.ShopCollection) {
	convert_catalogmodel_ShopCollection_catalog_ShopCollection(in, out)
}

func ShopProductCollectionDB(in *catalog.ShopProductCollection, out *catalogmodel.ShopProductCollection) {
	convert_catalog_ShopProductCollection_catalogmodel_ShopProductCollection(in, out)
}

func ShopProducCollection(in *catalogmodel.ShopProductCollection, out *catalog.ShopProductCollection) {
	convert_catalogmodel_ShopProductCollection_catalog_ShopProductCollection(in, out)
}

func ShopCollections(ins []*catalogmodel.ShopCollection) (outs []*catalog.ShopCollection) {
	return Convert_catalogmodel_ShopCollections_catalog_ShopCollections(ins)
}

func ShopProductCollections(ins []*catalogmodel.ShopProductCollection) (outs []*catalog.ShopProductCollection) {
	return Convert_catalogmodel_ShopProductCollections_catalog_ShopProductCollections(ins)
}

func UpdateShopCollection(in *catalogmodel.ShopCollection, args *catalog.UpdateShopCollectionArgs) (out *catalogmodel.ShopCollection) {
	if in == nil {
		return nil
	}
	shopColelction := &catalogmodel.ShopCollection{
		ID:          args.ID,
		ShopID:      args.ShopID,
		Name:        args.Name.Apply(in.Name),
		Description: args.Description.Apply(in.Description),
		DescHTML:    args.DescHTML.Apply(in.DescHTML),
		ShortDesc:   args.ShortDesc.Apply(in.ShortDesc),
	}
	return shopColelction
}

func UpdateShopProductCategory(in *catalogmodel.ShopProduct, args *catalog.UpdateShopProductCategoryArgs) (out *catalogmodel.ShopProduct) {
	if in == nil {
		return nil
	}
	shopProduct := &catalogmodel.ShopProduct{
		ProductID:  args.ProductID,
		CategoryID: args.CategoryID,
	}
	return shopProduct
}

func createShopBrand(args *catalog.CreateBrandArgs, out *catalog.ShopBrand) {
	apply_catalog_CreateBrandArgs_catalog_ShopBrand(args, out)
	out.ID = cm.NewID()
}
