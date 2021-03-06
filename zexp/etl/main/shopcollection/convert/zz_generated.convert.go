// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	catalogmodel "o.o/backend/com/main/catalog/model"
	conversion "o.o/backend/pkg/common/conversion"
	shopcollectionmodel "o.o/backend/zexp/etl/main/shopcollection/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopcollectionmodel.ShopCollection)(nil), (*catalogmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		Convert_shopcollectionmodel_ShopCollection_catalogmodel_ShopCollection(arg.(*shopcollectionmodel.ShopCollection), out.(*catalogmodel.ShopCollection))
		return nil
	})
	s.Register(([]*shopcollectionmodel.ShopCollection)(nil), (*[]*catalogmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopcollectionmodel_ShopCollections_catalogmodel_ShopCollections(arg.([]*shopcollectionmodel.ShopCollection))
		*out.(*[]*catalogmodel.ShopCollection) = out0
		return nil
	})
	s.Register((*catalogmodel.ShopCollection)(nil), (*shopcollectionmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		Convert_catalogmodel_ShopCollection_shopcollectionmodel_ShopCollection(arg.(*catalogmodel.ShopCollection), out.(*shopcollectionmodel.ShopCollection))
		return nil
	})
	s.Register(([]*catalogmodel.ShopCollection)(nil), (*[]*shopcollectionmodel.ShopCollection)(nil), func(arg, out interface{}) error {
		out0 := Convert_catalogmodel_ShopCollections_shopcollectionmodel_ShopCollections(arg.([]*catalogmodel.ShopCollection))
		*out.(*[]*shopcollectionmodel.ShopCollection) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/catalog/model.ShopCollection --//

func Convert_shopcollectionmodel_ShopCollection_catalogmodel_ShopCollection(arg *shopcollectionmodel.ShopCollection, out *catalogmodel.ShopCollection) *catalogmodel.ShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &catalogmodel.ShopCollection{}
	}
	convert_shopcollectionmodel_ShopCollection_catalogmodel_ShopCollection(arg, out)
	return out
}

func convert_shopcollectionmodel_ShopCollection_catalogmodel_ShopCollection(arg *shopcollectionmodel.ShopCollection, out *catalogmodel.ShopCollection) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.PartnerID = arg.PartnerID     // simple assign
	out.ExternalID = arg.ExternalID   // simple assign
	out.Name = arg.Name               // simple assign
	out.Description = arg.Description // simple assign
	out.DescHTML = arg.DescHTML       // simple assign
	out.ShortDesc = arg.ShortDesc     // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
	out.DeletedAt = time.Time{}       // zero value
	out.Rid = arg.Rid                 // simple assign
}

func Convert_shopcollectionmodel_ShopCollections_catalogmodel_ShopCollections(args []*shopcollectionmodel.ShopCollection) (outs []*catalogmodel.ShopCollection) {
	if args == nil {
		return nil
	}
	tmps := make([]catalogmodel.ShopCollection, len(args))
	outs = make([]*catalogmodel.ShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_shopcollectionmodel_ShopCollection_catalogmodel_ShopCollection(args[i], &tmps[i])
	}
	return outs
}

func Convert_catalogmodel_ShopCollection_shopcollectionmodel_ShopCollection(arg *catalogmodel.ShopCollection, out *shopcollectionmodel.ShopCollection) *shopcollectionmodel.ShopCollection {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopcollectionmodel.ShopCollection{}
	}
	convert_catalogmodel_ShopCollection_shopcollectionmodel_ShopCollection(arg, out)
	return out
}

func convert_catalogmodel_ShopCollection_shopcollectionmodel_ShopCollection(arg *catalogmodel.ShopCollection, out *shopcollectionmodel.ShopCollection) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.PartnerID = arg.PartnerID     // simple assign
	out.ExternalID = arg.ExternalID   // simple assign
	out.Name = arg.Name               // simple assign
	out.Description = arg.Description // simple assign
	out.DescHTML = arg.DescHTML       // simple assign
	out.ShortDesc = arg.ShortDesc     // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
	out.Rid = arg.Rid                 // simple assign
}

func Convert_catalogmodel_ShopCollections_shopcollectionmodel_ShopCollections(args []*catalogmodel.ShopCollection) (outs []*shopcollectionmodel.ShopCollection) {
	if args == nil {
		return nil
	}
	tmps := make([]shopcollectionmodel.ShopCollection, len(args))
	outs = make([]*shopcollectionmodel.ShopCollection, len(args))
	for i := range tmps {
		outs[i] = Convert_catalogmodel_ShopCollection_shopcollectionmodel_ShopCollection(args[i], &tmps[i])
	}
	return outs
}
