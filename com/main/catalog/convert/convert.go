package convert

import (
	"time"

	"etop.vn/api/main/catalog"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/catalog/model->etop.vn/api/main/catalog,etop.vn/api/main/catalog/types
// +gen:convert: etop.vn/api/main/catalog

func shopProduct(in *catalogmodel.ShopProduct, out *catalog.ShopProduct) {
	metaFields := []*catalog.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}

	convert_catalogmodel_ShopProduct_catalog_ShopProduct(in, out)
	out.MetaFields = metaFields
}

func shopProductDB(in *catalog.ShopProduct, out *catalogmodel.ShopProduct) {

	metaFields := []*catalogmodel.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalogmodel.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	convert_catalog_ShopProduct_catalogmodel_ShopProduct(in, out)
	out.MetaFields = metaFields
}

func ShopProductWithVariants(in *catalogmodel.ShopProductWithVariants) (out *catalog.ShopProductWithVariants) {
	if in == nil {
		return nil
	}
	shopVariants := Convert_catalogmodel_ShopVariants_catalog_ShopVariants(in.Variants)
	out = &catalog.ShopProductWithVariants{
		ShopProduct: Convert_catalogmodel_ShopProduct_catalog_ShopProduct(in.ShopProduct, nil),
		Variants:    shopVariants,
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

func shopVariantDB(in *catalog.ShopVariant, out *catalogmodel.ShopVariant) {
	convert_catalog_ShopVariant_catalogmodel_ShopVariant(in, out)
	attributes, attrNormKv := catalogmodel.NormalizeAttributes(in.Attributes)
	out.Attributes = Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(attributes)
	out.AttrNormKv = attrNormKv
}

func ShopVariantWithProduct(in *catalogmodel.ShopVariantWithProduct) (out *catalog.ShopVariantWithProduct) {
	if in == nil {
		return nil
	}
	var shopVariant *catalog.ShopVariant
	convert_catalogmodel_ShopVariant_catalog_ShopVariant(in.ShopVariant, shopVariant)
	out = &catalog.ShopVariantWithProduct{
		ShopProduct: Convert_catalogmodel_ShopProduct_catalog_ShopProduct(in.ShopProduct, nil),
		ShopVariant: shopVariant,
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

func createShopBrand(args *catalog.CreateBrandArgs, out *catalog.ShopBrand) {
	apply_catalog_CreateBrandArgs_catalog_ShopBrand(args, out)
	out.ID = cm.NewID()
}

func updateShopProduct(args *catalog.UpdateShopProductInfoArgs, in *catalog.ShopProduct) (out *catalog.ShopProduct) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopProductInfoArgs_catalog_ShopProduct(args, in)
	in.UpdatedAt = time.Now()
	return in
}

func updateShopVariant(args *catalog.UpdateShopVariantInfoArgs, in *catalog.ShopVariant) (out *catalog.ShopVariant) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopVariantInfoArgs_catalog_ShopVariant(args, in)
	in.UpdatedAt = time.Now()
	return in
}

func updateShopCollection(args *catalog.UpdateShopCollectionArgs, in *catalog.ShopCollection) (out *catalog.ShopCollection) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopCollectionArgs_catalog_ShopCollection(args, in)
	in.UpdatedAt = time.Now()
	return in
}
func updateShopCategory(args *catalog.UpdateShopCategoryArgs, in *catalog.ShopCategory) (out *catalog.ShopCategory) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopCategoryArgs_catalog_ShopCategory(args, in)
	in.UpdatedAt = time.Now()
	return in
}
func updateShopProductCategory(args *catalog.UpdateShopProductCategoryArgs, in *catalog.ShopProduct) (out *catalog.ShopProduct) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopProductCategoryArgs_catalog_ShopProduct(args, in)
	in.UpdatedAt = time.Now()
	return in
}
