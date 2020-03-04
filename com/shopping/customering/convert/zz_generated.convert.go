// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	addressing "etop.vn/api/shopping/addressing"
	customering "etop.vn/api/shopping/customering"
	tradering "etop.vn/api/shopping/tradering"
	status3 "etop.vn/api/top/types/etc/status3"
	addressconvert "etop.vn/backend/com/main/address/convert"
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions:
    CreateShopCustomer         // not use, no conversions between params
    CreateShopTraderAddress    // in use
    ShopTraderAddress          // in use
    ShopTraderAddressDB        // in use
    UpdateShopTraderAddress    // in use
    shopCustomer               // in use
    shopCustomerDB             // in use
    updateShopCustomer         // in use

Ignored functions:
    GenerateCode           // params are not pointer to named types
    ParseCodeNorm          // not recognized
    UpdateCustomerGroup    // not recognized
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*customeringmodel.ShopTraderAddress)(nil), (*addressing.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(arg.(*customeringmodel.ShopTraderAddress), out.(*addressing.ShopTraderAddress))
		return nil
	})
	s.Register(([]*customeringmodel.ShopTraderAddress)(nil), (*[]*addressing.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopTraderAddresses_addressing_ShopTraderAddresses(arg.([]*customeringmodel.ShopTraderAddress))
		*out.(*[]*addressing.ShopTraderAddress) = out0
		return nil
	})
	s.Register((*addressing.ShopTraderAddress)(nil), (*customeringmodel.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		Convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(arg.(*addressing.ShopTraderAddress), out.(*customeringmodel.ShopTraderAddress))
		return nil
	})
	s.Register(([]*addressing.ShopTraderAddress)(nil), (*[]*customeringmodel.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressing_ShopTraderAddresses_customeringmodel_ShopTraderAddresses(arg.([]*addressing.ShopTraderAddress))
		*out.(*[]*customeringmodel.ShopTraderAddress) = out0
		return nil
	})
	s.Register((*addressing.CreateAddressArgs)(nil), (*addressing.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		Apply_addressing_CreateAddressArgs_addressing_ShopTraderAddress(arg.(*addressing.CreateAddressArgs), out.(*addressing.ShopTraderAddress))
		return nil
	})
	s.Register((*addressing.UpdateAddressArgs)(nil), (*addressing.ShopTraderAddress)(nil), func(arg, out interface{}) error {
		Apply_addressing_UpdateAddressArgs_addressing_ShopTraderAddress(arg.(*addressing.UpdateAddressArgs), out.(*addressing.ShopTraderAddress))
		return nil
	})
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
	s.Register((*customeringmodel.ShopCustomerGroup)(nil), (*customering.ShopCustomerGroup)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(arg.(*customeringmodel.ShopCustomerGroup), out.(*customering.ShopCustomerGroup))
		return nil
	})
	s.Register(([]*customeringmodel.ShopCustomerGroup)(nil), (*[]*customering.ShopCustomerGroup)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopCustomerGroups_customering_ShopCustomerGroups(arg.([]*customeringmodel.ShopCustomerGroup))
		*out.(*[]*customering.ShopCustomerGroup) = out0
		return nil
	})
	s.Register((*customering.ShopCustomerGroup)(nil), (*customeringmodel.ShopCustomerGroup)(nil), func(arg, out interface{}) error {
		Convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(arg.(*customering.ShopCustomerGroup), out.(*customeringmodel.ShopCustomerGroup))
		return nil
	})
	s.Register(([]*customering.ShopCustomerGroup)(nil), (*[]*customeringmodel.ShopCustomerGroup)(nil), func(arg, out interface{}) error {
		out0 := Convert_customering_ShopCustomerGroups_customeringmodel_ShopCustomerGroups(arg.([]*customering.ShopCustomerGroup))
		*out.(*[]*customeringmodel.ShopCustomerGroup) = out0
		return nil
	})
	s.Register((*customeringmodel.ShopCustomerGroupCustomer)(nil), (*customering.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(arg.(*customeringmodel.ShopCustomerGroupCustomer), out.(*customering.ShopCustomerGroupCustomer))
		return nil
	})
	s.Register(([]*customeringmodel.ShopCustomerGroupCustomer)(nil), (*[]*customering.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopCustomerGroupCustomers_customering_ShopCustomerGroupCustomers(arg.([]*customeringmodel.ShopCustomerGroupCustomer))
		*out.(*[]*customering.ShopCustomerGroupCustomer) = out0
		return nil
	})
	s.Register((*customering.ShopCustomerGroupCustomer)(nil), (*customeringmodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		Convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg.(*customering.ShopCustomerGroupCustomer), out.(*customeringmodel.ShopCustomerGroupCustomer))
		return nil
	})
	s.Register(([]*customering.ShopCustomerGroupCustomer)(nil), (*[]*customeringmodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customering_ShopCustomerGroupCustomers_customeringmodel_ShopCustomerGroupCustomers(arg.([]*customering.ShopCustomerGroupCustomer))
		*out.(*[]*customeringmodel.ShopCustomerGroupCustomer) = out0
		return nil
	})
	s.Register((*customeringmodel.ShopTrader)(nil), (*tradering.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopTrader_tradering_ShopTrader(arg.(*customeringmodel.ShopTrader), out.(*tradering.ShopTrader))
		return nil
	})
	s.Register(([]*customeringmodel.ShopTrader)(nil), (*[]*tradering.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopTraders_tradering_ShopTraders(arg.([]*customeringmodel.ShopTrader))
		*out.(*[]*tradering.ShopTrader) = out0
		return nil
	})
	s.Register((*tradering.ShopTrader)(nil), (*customeringmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_tradering_ShopTrader_customeringmodel_ShopTrader(arg.(*tradering.ShopTrader), out.(*customeringmodel.ShopTrader))
		return nil
	})
	s.Register(([]*tradering.ShopTrader)(nil), (*[]*customeringmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_tradering_ShopTraders_customeringmodel_ShopTraders(arg.([]*tradering.ShopTrader))
		*out.(*[]*customeringmodel.ShopTrader) = out0
		return nil
	})
}

//-- convert etop.vn/api/shopping/addressing.ShopTraderAddress --//

func Convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(arg *customeringmodel.ShopTraderAddress, out *addressing.ShopTraderAddress) *addressing.ShopTraderAddress {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressing.ShopTraderAddress{}
	}
	ShopTraderAddress(arg, out)
	return out
}

func convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(arg *customeringmodel.ShopTraderAddress, out *addressing.ShopTraderAddress) {
	out.ID = arg.ID                     // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.TraderID = arg.TraderID         // simple assign
	out.FullName = arg.FullName         // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.DistrictCode = arg.DistrictCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.IsDefault = arg.IsDefault       // simple assign
	out.Position = arg.Position         // simple assign
	out.Coordinates = addressconvert.Convert_addressmodel_Coordinates_orderingtypes_Coordinates(arg.Coordinates, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.Deleted = false           // zero value
}

func Convert_customeringmodel_ShopTraderAddresses_addressing_ShopTraderAddresses(args []*customeringmodel.ShopTraderAddress) (outs []*addressing.ShopTraderAddress) {
	tmps := make([]addressing.ShopTraderAddress, len(args))
	outs = make([]*addressing.ShopTraderAddress, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopTraderAddress_addressing_ShopTraderAddress(args[i], &tmps[i])
	}
	return outs
}

func Convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(arg *addressing.ShopTraderAddress, out *customeringmodel.ShopTraderAddress) *customeringmodel.ShopTraderAddress {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopTraderAddress{}
	}
	ShopTraderAddressDB(arg, out)
	return out
}

func convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(arg *addressing.ShopTraderAddress, out *customeringmodel.ShopTraderAddress) {
	out.ID = arg.ID                     // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.TraderID = arg.TraderID         // simple assign
	out.FullName = arg.FullName         // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.DistrictCode = arg.DistrictCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.Position = arg.Position         // simple assign
	out.IsDefault = arg.IsDefault       // simple assign
	out.Coordinates = addressconvert.Convert_orderingtypes_Coordinates_addressmodel_Coordinates(arg.Coordinates, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.DeletedAt = time.Time{}   // zero value
	out.Status = 0                // zero value
}

func Convert_addressing_ShopTraderAddresses_customeringmodel_ShopTraderAddresses(args []*addressing.ShopTraderAddress) (outs []*customeringmodel.ShopTraderAddress) {
	tmps := make([]customeringmodel.ShopTraderAddress, len(args))
	outs = make([]*customeringmodel.ShopTraderAddress, len(args))
	for i := range tmps {
		outs[i] = Convert_addressing_ShopTraderAddress_customeringmodel_ShopTraderAddress(args[i], &tmps[i])
	}
	return outs
}

func Apply_addressing_CreateAddressArgs_addressing_ShopTraderAddress(arg *addressing.CreateAddressArgs, out *addressing.ShopTraderAddress) *addressing.ShopTraderAddress {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressing.ShopTraderAddress{}
	}
	CreateShopTraderAddress(arg, out)
	return out
}

func apply_addressing_CreateAddressArgs_addressing_ShopTraderAddress(arg *addressing.CreateAddressArgs, out *addressing.ShopTraderAddress) {
	out.ID = 0                          // zero value
	out.ShopID = arg.ShopID             // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.TraderID = arg.TraderID         // simple assign
	out.FullName = arg.FullName         // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.DistrictCode = arg.DistrictCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.IsDefault = arg.IsDefault       // simple assign
	out.Position = arg.Position         // simple assign
	out.Coordinates = arg.Coordinates   // simple assign
	out.CreatedAt = time.Time{}         // zero value
	out.UpdatedAt = time.Time{}         // zero value
	out.Deleted = false                 // zero value
}

func Apply_addressing_UpdateAddressArgs_addressing_ShopTraderAddress(arg *addressing.UpdateAddressArgs, out *addressing.ShopTraderAddress) *addressing.ShopTraderAddress {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressing.ShopTraderAddress{}
	}
	UpdateShopTraderAddress(arg, out)
	return out
}

func apply_addressing_UpdateAddressArgs_addressing_ShopTraderAddress(arg *addressing.UpdateAddressArgs, out *addressing.ShopTraderAddress) {
	out.ID = out.ID                                             // identifier
	out.ShopID = out.ShopID                                     // no change
	out.PartnerID = out.PartnerID                               // no change
	out.TraderID = out.TraderID                                 // no change
	out.FullName = arg.FullName.Apply(out.FullName)             // apply change
	out.Phone = arg.Phone.Apply(out.Phone)                      // apply change
	out.Email = arg.Email.Apply(out.Email)                      // apply change
	out.Company = arg.Company.Apply(out.Company)                // apply change
	out.Address1 = arg.Address1.Apply(out.Address1)             // apply change
	out.Address2 = arg.Address2.Apply(out.Address2)             // apply change
	out.DistrictCode = arg.DistrictCode.Apply(out.DistrictCode) // apply change
	out.WardCode = arg.WardCode.Apply(out.WardCode)             // apply change
	out.IsDefault = arg.IsDefault.Apply(out.IsDefault)          // apply change
	out.Position = arg.Position.Apply(out.Position)             // apply change
	out.Coordinates = arg.Coordinates                           // simple assign
	out.CreatedAt = out.CreatedAt                               // no change
	out.UpdatedAt = out.UpdatedAt                               // no change
	out.Deleted = out.Deleted                                   // no change
}

//-- convert etop.vn/api/shopping/customering.ShopCustomer --//

func Convert_customeringmodel_ShopCustomer_customering_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customering.ShopCustomer) *customering.ShopCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopCustomer{}
	}
	shopCustomer(arg, out)
	return out
}

func convert_customeringmodel_ShopCustomer_customering_ShopCustomer(arg *customeringmodel.ShopCustomer, out *customering.ShopCustomer) {
	out.ExternalID = arg.ExternalID         // simple assign
	out.ExternalCode = arg.ExternalCode     // simple assign
	out.PartnerID = arg.PartnerID           // simple assign
	out.ID = arg.ID                         // simple assign
	out.ShopID = arg.ShopID                 // simple assign
	out.GroupIDs = arg.GroupIDs             // simple assign
	out.Code = arg.Code                     // simple assign
	out.CodeNorm = arg.CodeNorm             // simple assign
	out.FullName = arg.FullName             // simple assign
	out.Gender = arg.Gender                 // simple assign
	out.Type = arg.Type                     // simple assign
	out.Birthday = arg.Birthday             // simple assign
	out.Note = arg.Note                     // simple assign
	out.Phone = arg.Phone                   // simple assign
	out.Email = arg.Email                   // simple assign
	out.Status = status3.Status(arg.Status) // simple conversion
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
	out.Deleted = false                     // zero value
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
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopCustomer{}
	}
	shopCustomerDB(arg, out)
	return out
}

func convert_customering_ShopCustomer_customeringmodel_ShopCustomer(arg *customering.ShopCustomer, out *customeringmodel.ShopCustomer) {
	out.ExternalID = arg.ExternalID     // simple assign
	out.ExternalCode = arg.ExternalCode // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.ID = arg.ID                     // simple assign
	out.ShopID = arg.ShopID             // simple assign
	out.Code = arg.Code                 // simple assign
	out.CodeNorm = arg.CodeNorm         // simple assign
	out.FullName = arg.FullName         // simple assign
	out.Gender = arg.Gender             // simple assign
	out.Type = arg.Type                 // simple assign
	out.Birthday = arg.Birthday         // simple assign
	out.Note = arg.Note                 // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Status = int(arg.Status)        // simple conversion
	out.FullNameNorm = ""               // zero value
	out.PhoneNorm = ""                  // zero value
	out.GroupIDs = arg.GroupIDs         // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.DeletedAt = time.Time{}         // zero value
	out.Rid = 0                         // zero value
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
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopCustomer{}
	}
	apply_customering_CreateCustomerArgs_customering_ShopCustomer(arg, out)
	return out
}

func apply_customering_CreateCustomerArgs_customering_ShopCustomer(arg *customering.CreateCustomerArgs, out *customering.ShopCustomer) {
	out.ExternalID = arg.ExternalID     // simple assign
	out.ExternalCode = arg.ExternalCode // simple assign
	out.PartnerID = arg.PartnerID       // simple assign
	out.ID = 0                          // zero value
	out.ShopID = arg.ShopID             // simple assign
	out.GroupIDs = nil                  // zero value
	out.Code = ""                       // zero value
	out.CodeNorm = 0                    // zero value
	out.FullName = arg.FullName         // simple assign
	out.Gender = arg.Gender             // simple assign
	out.Type = arg.Type                 // simple assign
	out.Birthday = arg.Birthday         // simple assign
	out.Note = arg.Note                 // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Email = arg.Email               // simple assign
	out.Status = 0                      // zero value
	out.CreatedAt = time.Time{}         // zero value
	out.UpdatedAt = time.Time{}         // zero value
	out.Deleted = false                 // zero value
}

func Apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg *customering.UpdateCustomerArgs, out *customering.ShopCustomer) *customering.ShopCustomer {
	return updateShopCustomer(arg, out)
}

func apply_customering_UpdateCustomerArgs_customering_ShopCustomer(arg *customering.UpdateCustomerArgs, out *customering.ShopCustomer) {
	out.ExternalID = out.ExternalID                 // no change
	out.ExternalCode = out.ExternalCode             // no change
	out.PartnerID = out.PartnerID                   // no change
	out.ID = out.ID                                 // identifier
	out.ShopID = out.ShopID                         // identifier
	out.GroupIDs = out.GroupIDs                     // no change
	out.Code = out.Code                             // no change
	out.CodeNorm = out.CodeNorm                     // no change
	out.FullName = arg.FullName.Apply(out.FullName) // apply change
	out.Gender = arg.Gender.Apply(out.Gender)       // apply change
	out.Type = arg.Type                             // simple assign
	out.Birthday = arg.Birthday.Apply(out.Birthday) // apply change
	out.Note = arg.Note.Apply(out.Note)             // apply change
	out.Phone = arg.Phone.Apply(out.Phone)          // apply change
	out.Email = arg.Email.Apply(out.Email)          // apply change
	out.Status = out.Status                         // no change
	out.CreatedAt = out.CreatedAt                   // no change
	out.UpdatedAt = out.UpdatedAt                   // no change
	out.Deleted = out.Deleted                       // no change
}

//-- convert etop.vn/api/shopping/customering.ShopCustomerGroup --//

func Convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(arg *customeringmodel.ShopCustomerGroup, out *customering.ShopCustomerGroup) *customering.ShopCustomerGroup {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopCustomerGroup{}
	}
	convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(arg, out)
	return out
}

func convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(arg *customeringmodel.ShopCustomerGroup, out *customering.ShopCustomerGroup) {
	out.ID = arg.ID               // simple assign
	out.PartnerID = arg.PartnerID // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.Name = arg.Name           // simple assign
	out.Deleted = false           // zero value
}

func Convert_customeringmodel_ShopCustomerGroups_customering_ShopCustomerGroups(args []*customeringmodel.ShopCustomerGroup) (outs []*customering.ShopCustomerGroup) {
	tmps := make([]customering.ShopCustomerGroup, len(args))
	outs = make([]*customering.ShopCustomerGroup, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopCustomerGroup_customering_ShopCustomerGroup(args[i], &tmps[i])
	}
	return outs
}

func Convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(arg *customering.ShopCustomerGroup, out *customeringmodel.ShopCustomerGroup) *customeringmodel.ShopCustomerGroup {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopCustomerGroup{}
	}
	convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(arg, out)
	return out
}

func convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(arg *customering.ShopCustomerGroup, out *customeringmodel.ShopCustomerGroup) {
	out.ID = arg.ID               // simple assign
	out.PartnerID = arg.PartnerID // simple assign
	out.Name = arg.Name           // simple assign
	out.ShopID = arg.ShopID       // simple assign
	out.CreatedAt = time.Time{}   // zero value
	out.UpdatedAt = time.Time{}   // zero value
	out.DeletedAt = time.Time{}   // zero value
}

func Convert_customering_ShopCustomerGroups_customeringmodel_ShopCustomerGroups(args []*customering.ShopCustomerGroup) (outs []*customeringmodel.ShopCustomerGroup) {
	tmps := make([]customeringmodel.ShopCustomerGroup, len(args))
	outs = make([]*customeringmodel.ShopCustomerGroup, len(args))
	for i := range tmps {
		outs[i] = Convert_customering_ShopCustomerGroup_customeringmodel_ShopCustomerGroup(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/shopping/customering.ShopCustomerGroupCustomer --//

func Convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(arg *customeringmodel.ShopCustomerGroupCustomer, out *customering.ShopCustomerGroupCustomer) *customering.ShopCustomerGroupCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customering.ShopCustomerGroupCustomer{}
	}
	convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(arg, out)
	return out
}

func convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(arg *customeringmodel.ShopCustomerGroupCustomer, out *customering.ShopCustomerGroupCustomer) {
	out.GroupID = arg.GroupID       // simple assign
	out.CustomerID = arg.CustomerID // simple assign
	out.CreatedAt = arg.CreatedAt   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
}

func Convert_customeringmodel_ShopCustomerGroupCustomers_customering_ShopCustomerGroupCustomers(args []*customeringmodel.ShopCustomerGroupCustomer) (outs []*customering.ShopCustomerGroupCustomer) {
	tmps := make([]customering.ShopCustomerGroupCustomer, len(args))
	outs = make([]*customering.ShopCustomerGroupCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopCustomerGroupCustomer_customering_ShopCustomerGroupCustomer(args[i], &tmps[i])
	}
	return outs
}

func Convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg *customering.ShopCustomerGroupCustomer, out *customeringmodel.ShopCustomerGroupCustomer) *customeringmodel.ShopCustomerGroupCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopCustomerGroupCustomer{}
	}
	convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg, out)
	return out
}

func convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg *customering.ShopCustomerGroupCustomer, out *customeringmodel.ShopCustomerGroupCustomer) {
	out.GroupID = arg.GroupID       // simple assign
	out.CustomerID = arg.CustomerID // simple assign
	out.CreatedAt = arg.CreatedAt   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
}

func Convert_customering_ShopCustomerGroupCustomers_customeringmodel_ShopCustomerGroupCustomers(args []*customering.ShopCustomerGroupCustomer) (outs []*customeringmodel.ShopCustomerGroupCustomer) {
	tmps := make([]customeringmodel.ShopCustomerGroupCustomer, len(args))
	outs = make([]*customeringmodel.ShopCustomerGroupCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customering_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/shopping/tradering.ShopTrader --//

func Convert_customeringmodel_ShopTrader_tradering_ShopTrader(arg *customeringmodel.ShopTrader, out *tradering.ShopTrader) *tradering.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &tradering.ShopTrader{}
	}
	convert_customeringmodel_ShopTrader_tradering_ShopTrader(arg, out)
	return out
}

func convert_customeringmodel_ShopTrader_tradering_ShopTrader(arg *customeringmodel.ShopTrader, out *tradering.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
	out.FullName = ""       // zero value
	out.Phone = ""          // zero value
}

func Convert_customeringmodel_ShopTraders_tradering_ShopTraders(args []*customeringmodel.ShopTrader) (outs []*tradering.ShopTrader) {
	tmps := make([]tradering.ShopTrader, len(args))
	outs = make([]*tradering.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopTrader_tradering_ShopTrader(args[i], &tmps[i])
	}
	return outs
}

func Convert_tradering_ShopTrader_customeringmodel_ShopTrader(arg *tradering.ShopTrader, out *customeringmodel.ShopTrader) *customeringmodel.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopTrader{}
	}
	convert_tradering_ShopTrader_customeringmodel_ShopTrader(arg, out)
	return out
}

func convert_tradering_ShopTrader_customeringmodel_ShopTrader(arg *tradering.ShopTrader, out *customeringmodel.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
}

func Convert_tradering_ShopTraders_customeringmodel_ShopTraders(args []*tradering.ShopTrader) (outs []*customeringmodel.ShopTrader) {
	tmps := make([]customeringmodel.ShopTrader, len(args))
	outs = make([]*customeringmodel.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_tradering_ShopTrader_customeringmodel_ShopTrader(args[i], &tmps[i])
	}
	return outs
}
