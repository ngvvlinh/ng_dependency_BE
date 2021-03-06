// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	carrying "o.o/api/shopping/carrying"
	status3 "o.o/api/top/types/etc/status3"
	carryingmodel "o.o/backend/com/shopping/carrying/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions:
    createShopCarrier    // in use
    updateShopCarrier    // in use

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*carryingmodel.ShopCarrier)(nil), (*carrying.ShopCarrier)(nil), func(arg, out interface{}) error {
		Convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(arg.(*carryingmodel.ShopCarrier), out.(*carrying.ShopCarrier))
		return nil
	})
	s.Register(([]*carryingmodel.ShopCarrier)(nil), (*[]*carrying.ShopCarrier)(nil), func(arg, out interface{}) error {
		out0 := Convert_carryingmodel_ShopCarriers_carrying_ShopCarriers(arg.([]*carryingmodel.ShopCarrier))
		*out.(*[]*carrying.ShopCarrier) = out0
		return nil
	})
	s.Register((*carrying.ShopCarrier)(nil), (*carryingmodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		Convert_carrying_ShopCarrier_carryingmodel_ShopCarrier(arg.(*carrying.ShopCarrier), out.(*carryingmodel.ShopCarrier))
		return nil
	})
	s.Register(([]*carrying.ShopCarrier)(nil), (*[]*carryingmodel.ShopCarrier)(nil), func(arg, out interface{}) error {
		out0 := Convert_carrying_ShopCarriers_carryingmodel_ShopCarriers(arg.([]*carrying.ShopCarrier))
		*out.(*[]*carryingmodel.ShopCarrier) = out0
		return nil
	})
	s.Register((*carrying.CreateCarrierArgs)(nil), (*carrying.ShopCarrier)(nil), func(arg, out interface{}) error {
		Apply_carrying_CreateCarrierArgs_carrying_ShopCarrier(arg.(*carrying.CreateCarrierArgs), out.(*carrying.ShopCarrier))
		return nil
	})
	s.Register((*carrying.UpdateCarrierArgs)(nil), (*carrying.ShopCarrier)(nil), func(arg, out interface{}) error {
		Apply_carrying_UpdateCarrierArgs_carrying_ShopCarrier(arg.(*carrying.UpdateCarrierArgs), out.(*carrying.ShopCarrier))
		return nil
	})
}

//-- convert o.o/api/shopping/carrying.ShopCarrier --//

func Convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(arg *carryingmodel.ShopCarrier, out *carrying.ShopCarrier) *carrying.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &carrying.ShopCarrier{}
	}
	convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(arg, out)
	return out
}

func convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(arg *carryingmodel.ShopCarrier, out *carrying.ShopCarrier) {
	out.ID = arg.ID                         // simple assign
	out.ShopID = arg.ShopID                 // simple assign
	out.FullName = arg.FullName             // simple assign
	out.Note = arg.Note                     // simple assign
	out.Status = status3.Status(arg.Status) // simple conversion
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
}

func Convert_carryingmodel_ShopCarriers_carrying_ShopCarriers(args []*carryingmodel.ShopCarrier) (outs []*carrying.ShopCarrier) {
	if args == nil {
		return nil
	}
	tmps := make([]carrying.ShopCarrier, len(args))
	outs = make([]*carrying.ShopCarrier, len(args))
	for i := range tmps {
		outs[i] = Convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(args[i], &tmps[i])
	}
	return outs
}

func Convert_carrying_ShopCarrier_carryingmodel_ShopCarrier(arg *carrying.ShopCarrier, out *carryingmodel.ShopCarrier) *carryingmodel.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &carryingmodel.ShopCarrier{}
	}
	convert_carrying_ShopCarrier_carryingmodel_ShopCarrier(arg, out)
	return out
}

func convert_carrying_ShopCarrier_carryingmodel_ShopCarrier(arg *carrying.ShopCarrier, out *carryingmodel.ShopCarrier) {
	out.ID = arg.ID               // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.FullName = arg.FullName   // simple assign
	out.Note = arg.Note           // simple assign
	out.Status = int(arg.Status)  // simple conversion
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = time.Time{}   // zero value
	out.Rid = 0                   // zero value
}

func Convert_carrying_ShopCarriers_carryingmodel_ShopCarriers(args []*carrying.ShopCarrier) (outs []*carryingmodel.ShopCarrier) {
	if args == nil {
		return nil
	}
	tmps := make([]carryingmodel.ShopCarrier, len(args))
	outs = make([]*carryingmodel.ShopCarrier, len(args))
	for i := range tmps {
		outs[i] = Convert_carrying_ShopCarrier_carryingmodel_ShopCarrier(args[i], &tmps[i])
	}
	return outs
}

func Apply_carrying_CreateCarrierArgs_carrying_ShopCarrier(arg *carrying.CreateCarrierArgs, out *carrying.ShopCarrier) *carrying.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &carrying.ShopCarrier{}
	}
	createShopCarrier(arg, out)
	return out
}

func apply_carrying_CreateCarrierArgs_carrying_ShopCarrier(arg *carrying.CreateCarrierArgs, out *carrying.ShopCarrier) {
	out.ID = 0                  // zero value
	out.ShopID = arg.ShopID     // simple assign
	out.FullName = arg.FullName // simple assign
	out.Note = arg.Note         // simple assign
	out.Status = 0              // zero value
	out.CreatedAt = time.Time{} // zero value
	out.UpdatedAt = time.Time{} // zero value
}

func Apply_carrying_UpdateCarrierArgs_carrying_ShopCarrier(arg *carrying.UpdateCarrierArgs, out *carrying.ShopCarrier) *carrying.ShopCarrier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &carrying.ShopCarrier{}
	}
	updateShopCarrier(arg, out)
	return out
}

func apply_carrying_UpdateCarrierArgs_carrying_ShopCarrier(arg *carrying.UpdateCarrierArgs, out *carrying.ShopCarrier) {
	out.ID = out.ID                                 // identifier
	out.ShopID = out.ShopID                         // identifier
	out.FullName = arg.FullName.Apply(out.FullName) // apply change
	out.Note = arg.Note.Apply(out.Note)             // apply change
	out.Status = out.Status                         // no change
	out.CreatedAt = out.CreatedAt                   // no change
	out.UpdatedAt = out.UpdatedAt                   // no change
}
