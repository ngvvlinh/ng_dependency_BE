// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	customermodel "etop.vn/backend/zexp/etl/main/customer/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*customermodel.ShopCustomer)(nil), (*customeringmodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_customermodel_ShopCustomer_customeringmodel_ShopCustomer(arg.(*customermodel.ShopCustomer), out.(*customeringmodel.ShopCustomer))
		return nil
	})
	s.Register(([]*customermodel.ShopCustomer)(nil), (*[]*customeringmodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customermodel_ShopCustomers_customeringmodel_ShopCustomers(arg.([]*customermodel.ShopCustomer))
		*out.(*[]*customeringmodel.ShopCustomer) = out0
		return nil
	})
	s.Register((*customeringmodel.ShopCustomer)(nil), (*customermodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopCustomer_customermodel_ShopCustomer(arg.(*customeringmodel.ShopCustomer), out.(*customermodel.ShopCustomer))
		return nil
	})
	s.Register(([]*customeringmodel.ShopCustomer)(nil), (*[]*customermodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopCustomers_customermodel_ShopCustomers(arg.([]*customeringmodel.ShopCustomer))
		*out.(*[]*customermodel.ShopCustomer) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/shopping/customering/model.ShopCustomer --//

func Convert_customermodel_ShopCustomer_customeringmodel_ShopCustomer(arg *customermodel.ShopCustomer, out *customeringmodel.ShopCustomer) *customeringmodel.ShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopCustomer{}
	}
	convert_customermodel_ShopCustomer_customeringmodel_ShopCustomer(arg, out)
	return out
}

func convert_customermodel_ShopCustomer_customeringmodel_ShopCustomer(arg *customermodel.ShopCustomer, out *customeringmodel.ShopCustomer) {
	out.ExternalID = arg.ExternalID     // simple assign
	out.ExternalCode = arg.ExternalCode // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.ID = arg.ID                     // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.Code = arg.Code                 // simple assign
	out.CodeNorm = 0                    // zero value
	out.FullName = arg.FullName         // simple assign
	out.Gender = arg.Gender             // simple assign
	out.Type = arg.Type                 // simple assign
	out.Birthday = arg.Birthday         // simple assign
	out.Note = arg.Note                 // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Status = arg.Status             // simple assign
	out.FullNameNorm = ""               // zero value
	out.PhoneNorm = ""                  // zero value
	out.GroupIDs = arg.GroupIDs         // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.DeletedAt = time.Time{}         // zero value
	out.Rid = arg.Rid                   // simple assign
}

func Convert_customermodel_ShopCustomers_customeringmodel_ShopCustomers(args []*customermodel.ShopCustomer) (outs []*customeringmodel.ShopCustomer) {
	if args == nil {
		return nil
	}
	tmps := make([]customeringmodel.ShopCustomer, len(args))
	outs = make([]*customeringmodel.ShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customermodel_ShopCustomer_customeringmodel_ShopCustomer(args[i], &tmps[i])
	}
	return outs
}

func Convert_customeringmodel_ShopCustomer_customermodel_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customermodel.ShopCustomer) *customermodel.ShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customermodel.ShopCustomer{}
	}
	convert_customeringmodel_ShopCustomer_customermodel_ShopCustomer(arg, out)
	return out
}

func convert_customeringmodel_ShopCustomer_customermodel_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customermodel.ShopCustomer) {
	out.ExternalID = arg.ExternalID     // simple assign
	out.ExternalCode = arg.ExternalCode // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.ID = arg.ID                     // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.Code = arg.Code                 // simple assign
	out.FullName = arg.FullName         // simple assign
	out.Gender = arg.Gender             // simple assign
	out.Type = arg.Type                 // simple assign
	out.Birthday = arg.Birthday         // simple assign
	out.Note = arg.Note                 // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Status = arg.Status             // simple assign
	out.GroupIDs = arg.GroupIDs         // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Rid = arg.Rid                   // simple assign
}

func Convert_customeringmodel_ShopCustomers_customermodel_ShopCustomers(args []*customeringmodel.ShopCustomer) (outs []*customermodel.ShopCustomer) {
	if args == nil {
		return nil
	}
	tmps := make([]customermodel.ShopCustomer, len(args))
	outs = make([]*customermodel.ShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopCustomer_customermodel_ShopCustomer(args[i], &tmps[i])
	}
	return outs
}
