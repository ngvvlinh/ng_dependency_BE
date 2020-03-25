// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	carryingmodel "etop.vn/backend/com/shopping/carrying/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	shopcarriermodel "etop.vn/backend/zexp/etl/main/shopcarrier/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopcarriermodel.ShopCarrier)(nil), (*carryingmodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		Convert_shopcarriermodel_ShopCarrier_carryingmodel_ShopCarrier(arg.(*shopcarriermodel.ShopCarrier), out.(*carryingmodel.ShopCarrier))
		return nil
	})
	s.Register(([]*shopcarriermodel.ShopCarrier)(nil), (*[]*carryingmodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopcarriermodel_ShopCarriers_carryingmodel_ShopCarriers(arg.([]*shopcarriermodel.ShopCarrier))
		*out.(*[]*carryingmodel.ShopCarrier) = out0
		return nil
	})
	s.Register((*carryingmodel.ShopCarrier)(nil), (*shopcarriermodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		Convert_carryingmodel_ShopCarrier_shopcarriermodel_ShopCarrier(arg.(*carryingmodel.ShopCarrier), out.(*shopcarriermodel.ShopCarrier))
		return nil
	})
	s.Register(([]*carryingmodel.ShopCarrier)(nil), (*[]*shopcarriermodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		out0 := Convert_carryingmodel_ShopCarriers_shopcarriermodel_ShopCarriers(arg.([]*carryingmodel.ShopCarrier))
		*out.(*[]*shopcarriermodel.ShopCarrier) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/shopping/carrying/model.ShopCarrier --//

func Convert_shopcarriermodel_ShopCarrier_carryingmodel_ShopCarrier(arg *shopcarriermodel.ShopCarrier, out *carryingmodel.ShopCarrier) *carryingmodel.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &carryingmodel.ShopCarrier{}
	}
	convert_shopcarriermodel_ShopCarrier_carryingmodel_ShopCarrier(arg, out)
	return out
}

func convert_shopcarriermodel_ShopCarrier_carryingmodel_ShopCarrier(arg *shopcarriermodel.ShopCarrier, out *carryingmodel.ShopCarrier) {
	out.ID = arg.ID               // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.FullName = arg.FullName   // simple assign
	out.Note = arg.Note           // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = time.Time{}   // zero value
	out.Rid = arg.Rid             // simple assign
}

func Convert_shopcarriermodel_ShopCarriers_carryingmodel_ShopCarriers(args []*shopcarriermodel.ShopCarrier) (outs []*carryingmodel.ShopCarrier) {
	tmps := make([]carryingmodel.ShopCarrier, len(args))
	outs = make([]*carryingmodel.ShopCarrier, len(args))
	for i := range tmps {
		outs[i] = Convert_shopcarriermodel_ShopCarrier_carryingmodel_ShopCarrier(args[i], &tmps[i])
	}
	return outs
}

func Convert_carryingmodel_ShopCarrier_shopcarriermodel_ShopCarrier(arg *carryingmodel.ShopCarrier, out *shopcarriermodel.ShopCarrier) *shopcarriermodel.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopcarriermodel.ShopCarrier{}
	}
	convert_carryingmodel_ShopCarrier_shopcarriermodel_ShopCarrier(arg, out)
	return out
}

func convert_carryingmodel_ShopCarrier_shopcarriermodel_ShopCarrier(arg *carryingmodel.ShopCarrier, out *shopcarriermodel.ShopCarrier) {
	out.ID = arg.ID               // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.FullName = arg.FullName   // simple assign
	out.Note = arg.Note           // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.Rid = arg.Rid             // simple assign
}

func Convert_carryingmodel_ShopCarriers_shopcarriermodel_ShopCarriers(args []*carryingmodel.ShopCarrier) (outs []*shopcarriermodel.ShopCarrier) {
	tmps := make([]shopcarriermodel.ShopCarrier, len(args))
	outs = make([]*shopcarriermodel.ShopCarrier, len(args))
	for i := range tmps {
		outs[i] = Convert_carryingmodel_ShopCarrier_shopcarriermodel_ShopCarrier(args[i], &tmps[i])
	}
	return outs
}
