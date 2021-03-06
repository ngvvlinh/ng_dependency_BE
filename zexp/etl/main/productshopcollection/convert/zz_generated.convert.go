// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	catalogmodel "o.o/backend/com/main/catalog/model"
	conversion "o.o/backend/pkg/common/conversion"
	productshopcollectionmodel "o.o/backend/zexp/etl/main/productshopcollection/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*productshopcollectionmodel.ProductShopCollection)(nil), (*catalogmodel.ProductShopCollection)(nil), func(arg, out interface{}) error {
		Convert_productshopcollectionmodel_ProductShopCollection_catalogmodel_ProductShopCollection(arg.(*productshopcollectionmodel.ProductShopCollection), out.(*catalogmodel.ProductShopCollection))
		return nil
	})
	s.Register(([]*productshopcollectionmodel.ProductShopCollection)(nil), (*[]*catalogmodel.ProductShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_productshopcollectionmodel_ProductShopCollections_catalogmodel_ProductShopCollections(arg.([]*productshopcollectionmodel.ProductShopCollection))
		*out.(*[]*catalogmodel.ProductShopCollection) = out0
		return nil
	})
	s.Register((*catalogmodel.ProductShopCollection)(nil), (*productshopcollectionmodel.ProductShopCollection)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ProductShopCollection_productshopcollectionmodel_ProductShopCollection(arg.(*catalogmodel.ProductShopCollection), out.(*productshopcollectionmodel.ProductShopCollection))
		return nil
	})
	s.Register(([]*catalogmodel.ProductShopCollection)(nil), (*[]*productshopcollectionmodel.ProductShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ProductShopCollections_productshopcollectionmodel_ProductShopCollections(arg.([]*catalogmodel.ProductShopCollection))
		*out.(*[]*productshopcollectionmodel.ProductShopCollection) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/catalog/model.ProductShopCollection --//

func Convert_productshopcollectionmodel_ProductShopCollection_catalogmodel_ProductShopCollection(arg *productshopcollectionmodel.ProductShopCollection, out *catalogmodel.ProductShopCollection) *catalogmodel.ProductShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ProductShopCollection{}
	}
	convert_productshopcollectionmodel_ProductShopCollection_catalogmodel_ProductShopCollection(arg, out)
	return out
}

func convert_productshopcollectionmodel_ProductShopCollection_catalogmodel_ProductShopCollection(arg *productshopcollectionmodel.ProductShopCollection, out *catalogmodel.ProductShopCollection) {
	out.CollectionID = arg.CollectionID // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.Status = arg.Status             // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Rid = arg.Rid                   // simple assign
}

func Convert_productshopcollectionmodel_ProductShopCollections_catalogmodel_ProductShopCollections(args []*productshopcollectionmodel.ProductShopCollection) (outs []*catalogmodel.ProductShopCollection) {
	if args == nil {
		return nil
	}
	tmps := make([]catalogmodel.ProductShopCollection, len(args))
	outs = make([]*catalogmodel.ProductShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_productshopcollectionmodel_ProductShopCollection_catalogmodel_ProductShopCollection(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalogmodel_ProductShopCollection_productshopcollectionmodel_ProductShopCollection(arg *catalogmodel.ProductShopCollection, out *productshopcollectionmodel.ProductShopCollection) *productshopcollectionmodel.ProductShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &productshopcollectionmodel.ProductShopCollection{}
	}
	convert_catalogmodel_ProductShopCollection_productshopcollectionmodel_ProductShopCollection(arg, out)
	return out
}

func convert_catalogmodel_ProductShopCollection_productshopcollectionmodel_ProductShopCollection(arg *catalogmodel.ProductShopCollection, out *productshopcollectionmodel.ProductShopCollection) {
	out.CollectionID = arg.CollectionID // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.Status = arg.Status             // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Rid = arg.Rid                   // simple assign
}

func Convert_catalogmodel_ProductShopCollections_productshopcollectionmodel_ProductShopCollections(args []*catalogmodel.ProductShopCollection) (outs []*productshopcollectionmodel.ProductShopCollection) {
	if args == nil {
		return nil
	}
	tmps := make([]productshopcollectionmodel.ProductShopCollection, len(args))
	outs = make([]*productshopcollectionmodel.ProductShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ProductShopCollection_productshopcollectionmodel_ProductShopCollection(args[i], &tmps[i])
	}
	return outs
}
