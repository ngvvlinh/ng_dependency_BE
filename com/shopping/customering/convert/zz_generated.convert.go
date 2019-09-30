// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	customering "etop.vn/api/shopping/customering"
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	scheme "etop.vn/backend/pkg/common/scheme"
)

/*
Custom conversions:
    CreateShopCustomer         // in use
    CreateShopTraderAddress    // not use, no conversions between params
    ShopCustomer               // in use
    ShopCustomerDB             // in use
    ShopTraderAddress          // not use, no conversions between params
    ShopTraderAddressDB        // not use, no conversions between params

Ignored functions:
    Addresses                  // params are not pointer to named types
    ShopCustomers              // params are not pointer to named types
    UpdateShopCustomer         // not recognized
    UpdateShopTraderAddress    // not recognized
*/

func init() {
	registerConversionFunctions(scheme.Global)
}

func registerConversionFunctions(s *scheme.Scheme) {
	s.Register((*customeringmodel.ShopCustomer)(nil), (*customering.ShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopCustomer_customering_ShopCustomer(arg.(*customeringmodel.ShopCustomer), out.(*customering.ShopCustomer))
		return nil
	})
	s.Register(([]*customeringmodel.ShopCustomer)(nil), (*[]*customering.ShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopCustomers_customering_ShopCustomers(arg.([]*customeringmodel.ShopCustomer))
		*out.(*[]*customering.ShopCustomer) = out0
		return nil
	})
	s.Register((*customering.ShopCustomer)(nil), (*customeringmodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		Convert_customering_ShopCustomer_customeringmodel_ShopCustomer(arg.(*customering.ShopCustomer), out.(*customeringmodel.ShopCustomer))
		return nil
	})
	s.Register(([]*customering.ShopCustomer)(nil), (*[]*customeringmodel.ShopCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customering_ShopCustomers_customeringmodel_ShopCustomers(arg.([]*customering.ShopCustomer))
		*out.(*[]*customeringmodel.ShopCustomer) = out0
		return nil
	})
	s.Register((*customering.CreateCustomerArgs)(nil), (*customering.ShopCustomer)(nil), func(arg, out interface{}) error {
		Apply_customering_CreateCustomerArgs_customering_ShopCustomer(arg.(*customering.CreateCustomerArgs), out.(*customering.ShopCustomer))
		return nil
	})
	s.Register((*customering.UpdateCustomerArgs)(nil), (*customering.ShopCustomer)(nil), func(arg, out interface{}) error {
		Apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg.(*customering.UpdateCustomerArgs), out.(*customering.ShopCustomer))
		return nil
	})
	s.Register((*customeringmodel.ShopTrader)(nil), (*customering.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopTrader_customering_ShopTrader(arg.(*customeringmodel.ShopTrader), out.(*customering.ShopTrader))
		return nil
	})
	s.Register(([]*customeringmodel.ShopTrader)(nil), (*[]*customering.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopTraders_customering_ShopTraders(arg.([]*customeringmodel.ShopTrader))
		*out.(*[]*customering.ShopTrader) = out0
		return nil
	})
	s.Register((*customering.ShopTrader)(nil), (*customeringmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_customering_ShopTrader_customeringmodel_ShopTrader(arg.(*customering.ShopTrader), out.(*customeringmodel.ShopTrader))
		return nil
	})
	s.Register(([]*customering.ShopTrader)(nil), (*[]*customeringmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_customering_ShopTraders_customeringmodel_ShopTraders(arg.([]*customering.ShopTrader))
		*out.(*[]*customeringmodel.ShopTrader) = out0
		return nil
	})
}

//-- convert etop.vn/api/shopping/customering.ShopCustomer --//

func Convert_customeringmodel_ShopCustomer_customering_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customering.ShopCustomer) *customering.ShopCustomer {
	return ShopCustomer(arg)
}

func convert_customeringmodel_ShopCustomer_customering_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customering.ShopCustomer) {
	out.ID = arg.ID               // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.Code = arg.Code           // simple assign
	out.FullName = arg.FullName   // simple assign
	out.Gender = arg.Gender       // simple assign
	out.Type = arg.Type           // simple assign
	out.Birthday = arg.Birthday   // simple assign
	out.Note = arg.Note           // simple assign
	out.Phone = arg.Phone         // simple assign
	out.Email = arg.Email         // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
}

func Convert_customeringmodel_ShopCustomers_customering_ShopCustomers(args []*customeringmodel.ShopCustomer) (outs []*customering.ShopCustomer) {
	tmps := make([]customering.ShopCustomer, len(args))
	outs = make([]*customering.ShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopCustomer_customering_ShopCustomer(args[i], &tmps[i])
	}
	return outs
}

func Convert_customering_ShopCustomer_customeringmodel_ShopCustomer(arg *customering.ShopCustomer, out *customeringmodel.ShopCustomer) *customeringmodel.ShopCustomer {
	return ShopCustomerDB(arg)
}

func convert_customering_ShopCustomer_customeringmodel_ShopCustomer(arg *customering.ShopCustomer, out *customeringmodel.ShopCustomer) {
	out.ID = arg.ID               // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.Code = arg.Code           // simple assign
	out.FullName = arg.FullName   // simple assign
	out.Gender = arg.Gender       // simple assign
	out.Type = arg.Type           // simple assign
	out.Birthday = arg.Birthday   // simple assign
	out.Note = arg.Note           // simple assign
	out.Phone = arg.Phone         // simple assign
	out.Email = arg.Email         // simple assign
	out.Status = arg.Status       // simple assign
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = time.Time{}   // zero value
}

func Convert_customering_ShopCustomers_customeringmodel_ShopCustomers(args []*customering.ShopCustomer) (outs []*customeringmodel.ShopCustomer) {
	tmps := make([]customeringmodel.ShopCustomer, len(args))
	outs = make([]*customeringmodel.ShopCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customering_ShopCustomer_customeringmodel_ShopCustomer(args[i], &tmps[i])
	}
	return outs
}

func Apply_customering_CreateCustomerArgs_customering_ShopCustomer(arg *customering.CreateCustomerArgs, out *customering.ShopCustomer) *customering.ShopCustomer {
	return CreateShopCustomer(arg)
}

func apply_customering_CreateCustomerArgs_customering_ShopCustomer(arg *customering.CreateCustomerArgs, out *customering.ShopCustomer) {
	out.ID = 0                  // zero value
	out.ShopID = arg.ShopID     // simple assign
	out.Code = arg.Code         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Gender = arg.Gender     // simple assign
	out.Type = arg.Type         // simple assign
	out.Birthday = arg.Birthday // simple assign
	out.Note = arg.Note         // simple assign
	out.Phone = arg.Phone       // simple assign
	out.Email = arg.Email       // simple assign
	out.Status = 0              // zero value
	out.CreatedAt = time.Time{} // zero value
	out.UpdatedAt = time.Time{} // zero value
}

func Apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg *customering.UpdateCustomerArgs, out *customering.ShopCustomer) *customering.ShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopCustomer{}
	}
	apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg, out)
	return out
}

func apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg *customering.UpdateCustomerArgs, out *customering.ShopCustomer) {
	out.ID = out.ID                                 // identifier
	out.ShopID = out.ShopID                         // identifier
	out.Code = arg.Code.Apply(out.Code)             // apply change
	out.FullName = arg.FullName.Apply(out.FullName) // apply change
	out.Gender = arg.Gender.Apply(out.Gender)       // apply change
	out.Type = arg.Type.Apply(out.Type)             // apply change
	out.Birthday = arg.Birthday.Apply(out.Birthday) // apply change
	out.Note = arg.Note.Apply(out.Note)             // apply change
	out.Phone = arg.Phone.Apply(out.Phone)          // apply change
	out.Email = arg.Email.Apply(out.Email)          // apply change
	out.Status = out.Status                         // no change
	out.CreatedAt = out.CreatedAt                   // no change
	out.UpdatedAt = out.UpdatedAt                   // no change
}

//-- convert etop.vn/api/shopping/customering.ShopTrader --//

func Convert_customeringmodel_ShopTrader_customering_ShopTrader(arg *customeringmodel.ShopTrader, out *customering.ShopTrader) *customering.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopTrader{}
	}
	convert_customeringmodel_ShopTrader_customering_ShopTrader(arg, out)
	return out
}

func convert_customeringmodel_ShopTrader_customering_ShopTrader(arg *customeringmodel.ShopTrader, out *customering.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
}

func Convert_customeringmodel_ShopTraders_customering_ShopTraders(args []*customeringmodel.ShopTrader) (outs []*customering.ShopTrader) {
	tmps := make([]customering.ShopTrader, len(args))
	outs = make([]*customering.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopTrader_customering_ShopTrader(args[i], &tmps[i])
	}
	return outs
}

func Convert_customering_ShopTrader_customeringmodel_ShopTrader(arg *customering.ShopTrader, out *customeringmodel.ShopTrader) *customeringmodel.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopTrader{}
	}
	convert_customering_ShopTrader_customeringmodel_ShopTrader(arg, out)
	return out
}

func convert_customering_ShopTrader_customeringmodel_ShopTrader(arg *customering.ShopTrader, out *customeringmodel.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
}

func Convert_customering_ShopTraders_customeringmodel_ShopTraders(args []*customering.ShopTrader) (outs []*customeringmodel.ShopTrader) {
	tmps := make([]customeringmodel.ShopTrader, len(args))
	outs = make([]*customeringmodel.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_customering_ShopTrader_customeringmodel_ShopTrader(args[i], &tmps[i])
	}
	return outs
}
