// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	shopproductcollectionmodel "etop.vn/backend/zexp/etl/main/shopproductcollection/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopproductcollectionmodel.ShopProductCollection)(nil), (*catalogmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		Convert_shopproductcollectionmodel_ShopProductCollection_catalogmodel_ShopProductCollection(arg.(*shopproductcollectionmodel.ShopProductCollection), out.(*catalogmodel.ShopProductCollection))
		return nil
	})
	s.Register(([]*shopproductcollectionmodel.ShopProductCollection)(nil), (*[]*catalogmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopproductcollectionmodel_ShopProductCollections_catalogmodel_ShopProductCollections(arg.([]*shopproductcollectionmodel.ShopProductCollection))
		*out.(*[]*catalogmodel.ShopProductCollection) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopProductCollection)(nil), (*shopproductcollectionmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopProductCollection_shopproductcollectionmodel_ShopProductCollection(arg.(*catalogmodel.ShopProductCollection), out.(*shopproductcollectionmodel.ShopProductCollection))
		return nil
	})
	s.Register(([]*catalogmodel.ShopProductCollection)(nil), (*[]*shopproductcollectionmodel.ShopProductCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopProductCollections_shopproductcollectionmodel_ShopProductCollections(arg.([]*catalogmodel.ShopProductCollection))
		*out.(*[]*shopproductcollectionmodel.ShopProductCollection) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/main/catalog/model.ShopProductCollection --//

func Convert_shopproductcollectionmodel_ShopProductCollection_catalogmodel_ShopProductCollection(arg *shopproductcollectionmodel.ShopProductCollection, out *catalogmodel.ShopProductCollection) *catalogmodel.ShopProductCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ShopProductCollection{}
	}
	convert_shopproductcollectionmodel_ShopProductCollection_catalogmodel_ShopProductCollection(arg, out)
	return out
}

func convert_shopproductcollectionmodel_ShopProductCollection_catalogmodel_ShopProductCollection(arg *shopproductcollectionmodel.ShopProductCollection, out *catalogmodel.ShopProductCollection) {
	out.PartnerID = 0                   // zero value
	out.ShopID = arg.ShopID             // simple assign
	out.ExternalCollectionID = ""       // zero value
	out.ExternalProductID = ""          // zero value
	out.ProductID = arg.ProductID       // simple assign
	out.CollectionID = arg.CollectionID // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Rid = arg.Rid                   // simple assign
}

func Convert_shopproductcollectionmodel_ShopProductCollections_catalogmodel_ShopProductCollections(args []*shopproductcollectionmodel.ShopProductCollection) (outs []*catalogmodel.ShopProductCollection) {
	tmps := make([]catalogmodel.ShopProductCollection, len(args))
	outs = make([]*catalogmodel.ShopProductCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_shopproductcollectionmodel_ShopProductCollection_catalogmodel_ShopProductCollection(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalogmodel_ShopProductCollection_shopproductcollectionmodel_ShopProductCollection(arg *catalogmodel.ShopProductCollection, out *shopproductcollectionmodel.ShopProductCollection) *shopproductcollectionmodel.ShopProductCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopproductcollectionmodel.ShopProductCollection{}
	}
	convert_catalogmodel_ShopProductCollection_shopproductcollectionmodel_ShopProductCollection(arg, out)
	return out
}

func convert_catalogmodel_ShopProductCollection_shopproductcollectionmodel_ShopProductCollection(arg *catalogmodel.ShopProductCollection, out *shopproductcollectionmodel.ShopProductCollection) {
	out.ShopID = arg.ShopID             // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.CollectionID = arg.CollectionID // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Rid = arg.Rid                   // simple assign
}

func Convert_catalogmodel_ShopProductCollections_shopproductcollectionmodel_ShopProductCollections(args []*catalogmodel.ShopProductCollection) (outs []*shopproductcollectionmodel.ShopProductCollection) {
	tmps := make([]shopproductcollectionmodel.ShopProductCollection, len(args))
	outs = make([]*shopproductcollectionmodel.ShopProductCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopProductCollection_shopproductcollectionmodel_ShopProductCollection(args[i], &tmps[i])
	}
	return outs
}
