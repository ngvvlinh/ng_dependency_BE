// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	catalog "etop.vn/api/main/catalog"
	catalogtypes "etop.vn/api/main/catalog/types"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	scheme "etop.vn/backend/pkg/common/scheme"
	etopmodel "etop.vn/backend/pkg/etop/model"
)

/*
Custom conversions:
    Attribute                  // not use, no conversions between params
    AttributeDB                // not use, no conversions between params
    ShopCategory               // in use
    ShopCategoryDB             // in use
    ShopCollection             // in use
    ShopCollectionDB           // in use
    ShopProducCollection       // in use
    ShopProduct                // in use
    ShopProductCollectionDB    // in use
    ShopProductDB              // in use
    ShopProductUpdate          // not use, no conversions between params
    ShopProductWithVariants    // not use, no conversions between params
    ShopVariant                // in use
    ShopVariantDB              // in use
    ShopVariantWithProduct     // not use, no conversions between params

Ignored functions:
    Attributes                   // params are not pointer to named types
    AttributesDB                 // params are not pointer to named types
    ShopCategories               // params are not pointer to named types
    ShopCollections              // params are not pointer to named types
    ShopProductCollections       // params are not pointer to named types
    ShopProducts                 // params are not pointer to named types
    ShopProductsWithVariants     // params are not pointer to named types
    ShopVariants                 // params are not pointer to named types
    ShopVariantsWithProduct      // params are not pointer to named types
    UpdateShopCategory           // not recognized
    UpdateShopCollection         // not recognized
    UpdateShopProduct            // not recognized
    UpdateShopProductCategory    // not recognized
    UpdateShopVariant            // not recognized
*/

func init() {
	registerConversionFunctions(scheme.Global)
}

func registerConversionFunctions(s *scheme.Scheme) {
	s.Register((*catalogmodel.ProductAttribute)(nil), (*catalogtypes.Attribute)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ProductAttribute_catalogtypes_Attribute(arg.(*catalogmodel.ProductAttribute), out.(*catalogtypes.Attribute))
		return nil
	})
	s.Register(([]*catalogmodel.ProductAttribute)(nil), (*[]*catalogtypes.Attribute)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(arg.([]*catalogmodel.ProductAttribute))
		*out.(*[]*catalogtypes.Attribute) = out0
		return nil
	})
	s.Register((*catalogtypes.Attribute)(nil), (*catalogmodel.ProductAttribute)(nil), func(arg, out interface{}) error {
		Convert_catalogtypes_Attribute_catalogmodel_ProductAttribute(arg.(*catalogtypes.Attribute), out.(*catalogmodel.ProductAttribute))
		return nil
	})
	s.Register(([]*catalogtypes.Attribute)(nil), (*[]*catalogmodel.ProductAttribute)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(arg.([]*catalogtypes.Attribute))
		*out.(*[]*catalogmodel.ProductAttribute) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopCategory)(nil), (*catalog.ShopCategory)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopCategory_catalog_ShopCategory(arg.(*catalogmodel.ShopCategory), out.(*catalog.ShopCategory))
		return nil
	})
	s.Register(([]*catalogmodel.ShopCategory)(nil), (*[]*catalog.ShopCategory)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopCategories_catalog_ShopCategories(arg.([]*catalogmodel.ShopCategory))
		*out.(*[]*catalog.ShopCategory) = out0
		return nil
	})
	s.Register((*catalog.ShopCategory)(nil), (*catalogmodel.ShopCategory)(nil), func(arg, out interface{}) error {
		Convert_catalog_ShopCategory_catalogmodel_ShopCategory(arg.(*catalog.ShopCategory), out.(*catalogmodel.ShopCategory))
		return nil
	})
	s.Register(([]*catalog.ShopCategory)(nil), (*[]*catalogmodel.ShopCategory)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalog_ShopCategories_catalogmodel_ShopCategories(arg.([]*catalog.ShopCategory))
		*out.(*[]*catalogmodel.ShopCategory) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopCollection)(nil), (*catalog.ShopCollection)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopCollection_catalog_ShopCollection(arg.(*catalogmodel.ShopCollection), out.(*catalog.ShopCollection))
		return nil
	})
	s.Register(([]*catalogmodel.ShopCollection)(nil), (*[]*catalog.ShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopCollections_catalog_ShopCollections(arg.([]*catalogmodel.ShopCollection))
		*out.(*[]*catalog.ShopCollection) = out0
		return nil
	})
	s.Register((*catalog.ShopCollection)(nil), (*catalogmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		Convert_catalog_ShopCollection_catalogmodel_ShopCollection(arg.(*catalog.ShopCollection), out.(*catalogmodel.ShopCollection))
		return nil
	})
	s.Register(([]*catalog.ShopCollection)(nil), (*[]*catalogmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalog_ShopCollections_catalogmodel_ShopCollections(arg.([]*catalog.ShopCollection))
		*out.(*[]*catalogmodel.ShopCollection) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopProduct)(nil), (*catalog.ShopProduct)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopProduct_catalog_ShopProduct(arg.(*catalogmodel.ShopProduct), out.(*catalog.ShopProduct))
		return nil
	})
	s.Register(([]*catalogmodel.ShopProduct)(nil), (*[]*catalog.ShopProduct)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopProducts_catalog_ShopProducts(arg.([]*catalogmodel.ShopProduct))
		*out.(*[]*catalog.ShopProduct) = out0
		return nil
	})
	s.Register((*catalog.ShopProduct)(nil), (*catalogmodel.ShopProduct)(nil), func(arg, out interface{}) error {
		Convert_catalog_ShopProduct_catalogmodel_ShopProduct(arg.(*catalog.ShopProduct), out.(*catalogmodel.ShopProduct))
		return nil
	})
	s.Register(([]*catalog.ShopProduct)(nil), (*[]*catalogmodel.ShopProduct)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalog_ShopProducts_catalogmodel_ShopProducts(arg.([]*catalog.ShopProduct))
		*out.(*[]*catalogmodel.ShopProduct) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopProductCollection)(nil), (*catalog.ShopProductCollection)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopProductCollection_catalog_ShopProductCollection(arg.(*catalogmodel.ShopProductCollection), out.(*catalog.ShopProductCollection))
		return nil
	})
	s.Register(([]*catalogmodel.ShopProductCollection)(nil), (*[]*catalog.ShopProductCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopProductCollections_catalog_ShopProductCollections(arg.([]*catalogmodel.ShopProductCollection))
		*out.(*[]*catalog.ShopProductCollection) = out0
		return nil
	})
	s.Register((*catalog.ShopProductCollection)(nil), (*catalogmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		Convert_catalog_ShopProductCollection_catalogmodel_ShopProductCollection(arg.(*catalog.ShopProductCollection), out.(*catalogmodel.ShopProductCollection))
		return nil
	})
	s.Register(([]*catalog.ShopProductCollection)(nil), (*[]*catalogmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalog_ShopProductCollections_catalogmodel_ShopProductCollections(arg.([]*catalog.ShopProductCollection))
		*out.(*[]*catalogmodel.ShopProductCollection) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopVariant)(nil), (*catalog.ShopVariant)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopVariant_catalog_ShopVariant(arg.(*catalogmodel.ShopVariant), out.(*catalog.ShopVariant))
		return nil
	})
	s.Register(([]*catalogmodel.ShopVariant)(nil), (*[]*catalog.ShopVariant)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopVariants_catalog_ShopVariants(arg.([]*catalogmodel.ShopVariant))
		*out.(*[]*catalog.ShopVariant) = out0
		return nil
	})
	s.Register((*catalog.ShopVariant)(nil), (*catalogmodel.ShopVariant)(nil), func(arg, out interface{}) error {
		Convert_catalog_ShopVariant_catalogmodel_ShopVariant(arg.(*catalog.ShopVariant), out.(*catalogmodel.ShopVariant))
		return nil
	})
	s.Register(([]*catalog.ShopVariant)(nil), (*[]*catalogmodel.ShopVariant)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalog_ShopVariants_catalogmodel_ShopVariants(arg.([]*catalog.ShopVariant))
		*out.(*[]*catalogmodel.ShopVariant) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/catalog.Attribute --//

func Convert_catalogmodel_ProductAttribute_catalogtypes_Attribute(arg *catalogmodel.ProductAttribute, out *catalogtypes.Attribute) *catalogtypes.Attribute {
	return Attribute(arg)
}

func convert_catalogmodel_ProductAttribute_catalogtypes_Attribute(arg *catalogmodel.ProductAttribute, out *catalogtypes.Attribute) {
	out.Name = arg.Name   // simple assign
	out.Value = arg.Value // simple assign
}

func Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(args []*catalogmodel.ProductAttribute) (outs []*catalogtypes.Attribute) {
	tmps := make([]catalogtypes.Attribute, len(args))
	outs = make([]*catalogtypes.Attribute, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ProductAttribute_catalogtypes_Attribute(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalogtypes_Attribute_catalogmodel_ProductAttribute(arg *catalogtypes.Attribute, out *catalogmodel.ProductAttribute) *catalogmodel.ProductAttribute {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ProductAttribute{}
	}
	AttributeDB(arg, out)
	return out
}

func convert_catalogtypes_Attribute_catalogmodel_ProductAttribute(arg *catalogtypes.Attribute, out *catalogmodel.ProductAttribute) {
	out.Name = arg.Name   // simple assign
	out.Value = arg.Value // simple assign
}

func Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(args []*catalogtypes.Attribute) (outs []*catalogmodel.ProductAttribute) {
	tmps := make([]catalogmodel.ProductAttribute, len(args))
	outs = make([]*catalogmodel.ProductAttribute, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogtypes_Attribute_catalogmodel_ProductAttribute(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/catalog.ShopCategory --//

func Convert_catalogmodel_ShopCategory_catalog_ShopCategory(arg *catalogmodel.ShopCategory, out *catalog.ShopCategory) *catalog.ShopCategory {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalog.ShopCategory{}
	}
	ShopCategory(arg, out)
	return out
}

func convert_catalogmodel_ShopCategory_catalog_ShopCategory(arg *catalogmodel.ShopCategory, out *catalog.ShopCategory) {
	out.ID = arg.ID               // simple assign
	out.ParentID = arg.ParentID   // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.Name = arg.Name           // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
}

func Convert_catalogmodel_ShopCategories_catalog_ShopCategories(args []*catalogmodel.ShopCategory) (outs []*catalog.ShopCategory) {
	tmps := make([]catalog.ShopCategory, len(args))
	outs = make([]*catalog.ShopCategory, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopCategory_catalog_ShopCategory(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalog_ShopCategory_catalogmodel_ShopCategory(arg *catalog.ShopCategory, out *catalogmodel.ShopCategory) *catalogmodel.ShopCategory {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ShopCategory{}
	}
	ShopCategoryDB(arg, out)
	return out
}

func convert_catalog_ShopCategory_catalogmodel_ShopCategory(arg *catalog.ShopCategory, out *catalogmodel.ShopCategory) {
	out.ID = arg.ID               // simple assign
	out.ParentID = arg.ParentID   // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.Name = arg.Name           // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
}

func Convert_catalog_ShopCategories_catalogmodel_ShopCategories(args []*catalog.ShopCategory) (outs []*catalogmodel.ShopCategory) {
	tmps := make([]catalogmodel.ShopCategory, len(args))
	outs = make([]*catalogmodel.ShopCategory, len(args))
	for i := range tmps {
		outs[i] = Convert_catalog_ShopCategory_catalogmodel_ShopCategory(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/catalog.ShopCollection --//

func Convert_catalogmodel_ShopCollection_catalog_ShopCollection(arg *catalogmodel.ShopCollection, out *catalog.ShopCollection) *catalog.ShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalog.ShopCollection{}
	}
	ShopCollection(arg, out)
	return out
}

func convert_catalogmodel_ShopCollection_catalog_ShopCollection(arg *catalogmodel.ShopCollection, out *catalog.ShopCollection) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.Name = arg.Name               // simple assign
	out.Description = arg.Description // simple assign
	out.DescHTML = arg.DescHTML       // simple assign
	out.ShortDesc = arg.ShortDesc     // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
}

func Convert_catalogmodel_ShopCollections_catalog_ShopCollections(args []*catalogmodel.ShopCollection) (outs []*catalog.ShopCollection) {
	tmps := make([]catalog.ShopCollection, len(args))
	outs = make([]*catalog.ShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopCollection_catalog_ShopCollection(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalog_ShopCollection_catalogmodel_ShopCollection(arg *catalog.ShopCollection, out *catalogmodel.ShopCollection) *catalogmodel.ShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ShopCollection{}
	}
	ShopCollectionDB(arg, out)
	return out
}

func convert_catalog_ShopCollection_catalogmodel_ShopCollection(arg *catalog.ShopCollection, out *catalogmodel.ShopCollection) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.Name = arg.Name               // simple assign
	out.Description = arg.Description // simple assign
	out.DescHTML = arg.DescHTML       // simple assign
	out.ShortDesc = arg.ShortDesc     // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
}

func Convert_catalog_ShopCollections_catalogmodel_ShopCollections(args []*catalog.ShopCollection) (outs []*catalogmodel.ShopCollection) {
	tmps := make([]catalogmodel.ShopCollection, len(args))
	outs = make([]*catalogmodel.ShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalog_ShopCollection_catalogmodel_ShopCollection(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/catalog.ShopProduct --//

func Convert_catalogmodel_ShopProduct_catalog_ShopProduct(arg *catalogmodel.ShopProduct, out *catalog.ShopProduct) *catalog.ShopProduct {
	return ShopProduct(arg)
}

func convert_catalogmodel_ShopProduct_catalog_ShopProduct(arg *catalogmodel.ShopProduct, out *catalog.ShopProduct) {
	out.ShopID = arg.ShopID                         // simple assign
	out.ProductID = arg.ProductID                   // simple assign
	out.Code = arg.Code                             // simple assign
	out.Name = arg.Name                             // simple assign
	out.Unit = arg.Unit                             // simple assign
	out.ImageURLs = arg.ImageURLs                   // simple assign
	out.Note = arg.Note                             // simple assign
	out.DescriptionInfo = catalog.DescriptionInfo{} // zero value
	out.PriceInfo = catalog.PriceInfo{}             // zero value
	out.CategoryID = arg.CategoryID                 // simple assign
	out.VendorID = arg.VendorID                     // simple assign
	out.CollectionIDs = arg.CollectionIDs           // simple assign
	out.Tags = arg.Tags                             // simple assign
	out.Status = int32(arg.Status)                  // simple conversion
	out.CreatedAt = arg.CreatedAt                   // simple assign
	out.UpdatedAt = arg.UpdatedAt                   // simple assign
	out.DeletedAt = arg.DeletedAt                   // simple assign
	out.ProductType = arg.ProductType               // simple assign
	out.MetaFields = nil                            // types do not match
}

func Convert_catalogmodel_ShopProducts_catalog_ShopProducts(args []*catalogmodel.ShopProduct) (outs []*catalog.ShopProduct) {
	tmps := make([]catalog.ShopProduct, len(args))
	outs = make([]*catalog.ShopProduct, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopProduct_catalog_ShopProduct(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalog_ShopProduct_catalogmodel_ShopProduct(arg *catalog.ShopProduct, out *catalogmodel.ShopProduct) *catalogmodel.ShopProduct {
	return ShopProductDB(arg)
}

func convert_catalog_ShopProduct_catalogmodel_ShopProduct(arg *catalog.ShopProduct, out *catalogmodel.ShopProduct) {
	out.ShopID = arg.ShopID                    // simple assign
	out.ProductID = arg.ProductID              // simple assign
	out.CollectionIDs = arg.CollectionIDs      // simple assign
	out.Code = arg.Code                        // simple assign
	out.Name = arg.Name                        // simple assign
	out.Description = ""                       // zero value
	out.DescHTML = ""                          // zero value
	out.ShortDesc = ""                         // zero value
	out.ImageURLs = arg.ImageURLs              // simple assign
	out.Note = arg.Note                        // simple assign
	out.Tags = arg.Tags                        // simple assign
	out.Unit = arg.Unit                        // simple assign
	out.CategoryID = arg.CategoryID            // simple assign
	out.VendorID = arg.VendorID                // simple assign
	out.CostPrice = 0                          // zero value
	out.ListPrice = 0                          // zero value
	out.RetailPrice = 0                        // zero value
	out.Status = etopmodel.Status3(arg.Status) // simple conversion
	out.CreatedAt = arg.CreatedAt              // simple assign
	out.UpdatedAt = arg.UpdatedAt              // simple assign
	out.DeletedAt = arg.DeletedAt              // simple assign
	out.NameNorm = ""                          // zero value
	out.NameNormUa = ""                        // zero value
	out.ProductType = arg.ProductType          // simple assign
	out.MetaFields = nil                       // types do not match
}

func Convert_catalog_ShopProducts_catalogmodel_ShopProducts(args []*catalog.ShopProduct) (outs []*catalogmodel.ShopProduct) {
	tmps := make([]catalogmodel.ShopProduct, len(args))
	outs = make([]*catalogmodel.ShopProduct, len(args))
	for i := range tmps {
		outs[i] = Convert_catalog_ShopProduct_catalogmodel_ShopProduct(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/catalog.ShopProductCollection --//

func Convert_catalogmodel_ShopProductCollection_catalog_ShopProductCollection(arg *catalogmodel.ShopProductCollection, out *catalog.ShopProductCollection) *catalog.ShopProductCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalog.ShopProductCollection{}
	}
	ShopProducCollection(arg, out)
	return out
}

func convert_catalogmodel_ShopProductCollection_catalog_ShopProductCollection(arg *catalogmodel.ShopProductCollection, out *catalog.ShopProductCollection) {
	out.ProductID = arg.ProductID       // simple assign
	out.CollectionID = arg.CollectionID // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
}

func Convert_catalogmodel_ShopProductCollections_catalog_ShopProductCollections(args []*catalogmodel.ShopProductCollection) (outs []*catalog.ShopProductCollection) {
	tmps := make([]catalog.ShopProductCollection, len(args))
	outs = make([]*catalog.ShopProductCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopProductCollection_catalog_ShopProductCollection(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalog_ShopProductCollection_catalogmodel_ShopProductCollection(arg *catalog.ShopProductCollection, out *catalogmodel.ShopProductCollection) *catalogmodel.ShopProductCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ShopProductCollection{}
	}
	ShopProductCollectionDB(arg, out)
	return out
}

func convert_catalog_ShopProductCollection_catalogmodel_ShopProductCollection(arg *catalog.ShopProductCollection, out *catalogmodel.ShopProductCollection) {
	out.ProductID = arg.ProductID       // simple assign
	out.CollectionID = arg.CollectionID // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
}

func Convert_catalog_ShopProductCollections_catalogmodel_ShopProductCollections(args []*catalog.ShopProductCollection) (outs []*catalogmodel.ShopProductCollection) {
	tmps := make([]catalogmodel.ShopProductCollection, len(args))
	outs = make([]*catalogmodel.ShopProductCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalog_ShopProductCollection_catalogmodel_ShopProductCollection(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/catalog.ShopVariant --//

func Convert_catalogmodel_ShopVariant_catalog_ShopVariant(arg *catalogmodel.ShopVariant, out *catalog.ShopVariant) *catalog.ShopVariant {
	return ShopVariant(arg)
}

func convert_catalogmodel_ShopVariant_catalog_ShopVariant(arg *catalogmodel.ShopVariant, out *catalog.ShopVariant) {
	out.ShopID = arg.ShopID                         // simple assign
	out.ProductID = arg.ProductID                   // simple assign
	out.VariantID = arg.VariantID                   // simple assign
	out.Code = arg.Code                             // simple assign
	out.Name = arg.Name                             // simple assign
	out.DescriptionInfo = catalog.DescriptionInfo{} // zero value
	out.ImageURLs = arg.ImageURLs                   // simple assign
	out.Status = int16(arg.Status)                  // simple conversion
	out.Attributes = Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(arg.Attributes)
	out.PriceInfo = catalog.PriceInfo{} // zero value
	out.Note = arg.Note                 // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.DeletedAt = arg.DeletedAt       // simple assign
}

func Convert_catalogmodel_ShopVariants_catalog_ShopVariants(args []*catalogmodel.ShopVariant) (outs []*catalog.ShopVariant) {
	tmps := make([]catalog.ShopVariant, len(args))
	outs = make([]*catalog.ShopVariant, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopVariant_catalog_ShopVariant(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalog_ShopVariant_catalogmodel_ShopVariant(arg *catalog.ShopVariant, out *catalogmodel.ShopVariant) *catalogmodel.ShopVariant {
	return ShopVariantDB(arg)
}

func convert_catalog_ShopVariant_catalogmodel_ShopVariant(arg *catalog.ShopVariant, out *catalogmodel.ShopVariant) {
	out.ShopID = arg.ShopID                    // simple assign
	out.VariantID = arg.VariantID              // simple assign
	out.ProductID = arg.ProductID              // simple assign
	out.Code = arg.Code                        // simple assign
	out.Name = arg.Name                        // simple assign
	out.Description = ""                       // zero value
	out.DescHTML = ""                          // zero value
	out.ShortDesc = ""                         // zero value
	out.ImageURLs = arg.ImageURLs              // simple assign
	out.Note = arg.Note                        // simple assign
	out.Tags = nil                             // zero value
	out.CostPrice = 0                          // zero value
	out.ListPrice = 0                          // zero value
	out.RetailPrice = 0                        // zero value
	out.Status = etopmodel.Status3(arg.Status) // simple conversion
	out.Attributes = Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(arg.Attributes)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = arg.DeletedAt // simple assign
	out.NameNorm = ""             // zero value
	out.AttrNormKv = ""           // zero value
}

func Convert_catalog_ShopVariants_catalogmodel_ShopVariants(args []*catalog.ShopVariant) (outs []*catalogmodel.ShopVariant) {
	tmps := make([]catalogmodel.ShopVariant, len(args))
	outs = make([]*catalogmodel.ShopVariant, len(args))
	for i := range tmps {
		outs[i] = Convert_catalog_ShopVariant_catalogmodel_ShopVariant(args[i], &tmps[i])
	}
	return outs
}
