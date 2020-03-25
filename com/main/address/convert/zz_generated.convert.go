// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	address "etop.vn/api/main/address"
	orderingtypes "etop.vn/api/main/ordering/types"
	addressmodel "etop.vn/backend/com/main/address/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions:
    Address                // in use
    AddressToModel         // in use
    Coordinates            // in use
    CoordinatesDB          // in use
    OrderAddress           // in use
    OrderAddressToModel    // in use

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*addressmodel.Address)(nil), (*address.Address)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_Address_address_Address(arg.(*addressmodel.Address), out.(*address.Address))
		return nil
	})
	s.Register(([]*addressmodel.Address)(nil), (*[]*address.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_Addresses_address_Addresses(arg.([]*addressmodel.Address))
		*out.(*[]*address.Address) = out0
		return nil
	})
	s.Register((*address.Address)(nil), (*addressmodel.Address)(nil), func(arg, out interface{}) error {
		Convert_address_Address_addressmodel_Address(arg.(*address.Address), out.(*addressmodel.Address))
		return nil
	})
	s.Register(([]*address.Address)(nil), (*[]*addressmodel.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_address_Addresses_addressmodel_Addresses(arg.([]*address.Address))
		*out.(*[]*addressmodel.Address) = out0
		return nil
	})
	s.Register((*addressmodel.Address)(nil), (*orderingtypes.Address)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_Address_orderingtypes_Address(arg.(*addressmodel.Address), out.(*orderingtypes.Address))
		return nil
	})
	s.Register(([]*addressmodel.Address)(nil), (*[]*orderingtypes.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_Addresses_orderingtypes_Addresses(arg.([]*addressmodel.Address))
		*out.(*[]*orderingtypes.Address) = out0
		return nil
	})
	s.Register((*orderingtypes.Address)(nil), (*addressmodel.Address)(nil), func(arg, out interface{}) error {
		Convert_orderingtypes_Address_addressmodel_Address(arg.(*orderingtypes.Address), out.(*addressmodel.Address))
		return nil
	})
	s.Register(([]*orderingtypes.Address)(nil), (*[]*addressmodel.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_orderingtypes_Addresses_addressmodel_Addresses(arg.([]*orderingtypes.Address))
		*out.(*[]*addressmodel.Address) = out0
		return nil
	})
	s.Register((*addressmodel.Coordinates)(nil), (*orderingtypes.Coordinates)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_Coordinates_orderingtypes_Coordinates(arg.(*addressmodel.Coordinates), out.(*orderingtypes.Coordinates))
		return nil
	})
	s.Register(([]*addressmodel.Coordinates)(nil), (*[]*orderingtypes.Coordinates)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_Coordinateses_orderingtypes_Coordinateses(arg.([]*addressmodel.Coordinates))
		*out.(*[]*orderingtypes.Coordinates) = out0
		return nil
	})
	s.Register((*orderingtypes.Coordinates)(nil), (*addressmodel.Coordinates)(nil), func(arg, out interface{}) error {
		Convert_orderingtypes_Coordinates_addressmodel_Coordinates(arg.(*orderingtypes.Coordinates), out.(*addressmodel.Coordinates))
		return nil
	})
	s.Register(([]*orderingtypes.Coordinates)(nil), (*[]*addressmodel.Coordinates)(nil), func(arg, out interface{}) error {
		out0 := Convert_orderingtypes_Coordinateses_addressmodel_Coordinateses(arg.([]*orderingtypes.Coordinates))
		*out.(*[]*addressmodel.Coordinates) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/address.Address --//

func Convert_addressmodel_Address_address_Address(arg *addressmodel.Address, out *address.Address) *address.Address {
	return Address(arg)
}

func convert_addressmodel_Address_address_Address(arg *addressmodel.Address, out *address.Address) {
	out.ID = arg.ID                     // simple assign
	out.FullName = arg.FullName         // simple assign
	out.FirstName = arg.FirstName       // simple assign
	out.LastName = arg.LastName         // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Position = arg.Position         // simple assign
	out.Email = arg.Email               // simple assign
	out.Country = arg.Country           // simple assign
	out.City = arg.City                 // simple assign
	out.Province = arg.Province         // simple assign
	out.District = arg.District         // simple assign
	out.Ward = arg.Ward                 // simple assign
	out.Zip = arg.Zip                   // simple assign
	out.DistrictCode = arg.DistrictCode // simple assign
	out.ProvinceCode = arg.ProvinceCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.Type = arg.Type                 // simple assign
	out.AccountID = arg.AccountID       // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Coordinates = Convert_addressmodel_Coordinates_orderingtypes_Coordinates(arg.Coordinates, nil)
}

func Convert_addressmodel_Addresses_address_Addresses(args []*addressmodel.Address) (outs []*address.Address) {
	tmps := make([]address.Address, len(args))
	outs = make([]*address.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_Address_address_Address(args[i], &tmps[i])
	}
	return outs
}

func Convert_address_Address_addressmodel_Address(arg *address.Address, out *addressmodel.Address) *addressmodel.Address {
	return AddressToModel(arg)
}

func convert_address_Address_addressmodel_Address(arg *address.Address, out *addressmodel.Address) {
	out.ID = arg.ID                     // simple assign
	out.FullName = arg.FullName         // simple assign
	out.FirstName = arg.FirstName       // simple assign
	out.LastName = arg.LastName         // simple assign
	out.Phone = arg.Phone               // simple assign
	out.Position = arg.Position         // simple assign
	out.Email = arg.Email               // simple assign
	out.Country = arg.Country           // simple assign
	out.City = arg.City                 // simple assign
	out.Province = arg.Province         // simple assign
	out.District = arg.District         // simple assign
	out.Ward = arg.Ward                 // simple assign
	out.Zip = arg.Zip                   // simple assign
	out.DistrictCode = arg.DistrictCode // simple assign
	out.ProvinceCode = arg.ProvinceCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.Type = arg.Type                 // simple assign
	out.AccountID = arg.AccountID       // simple assign
	out.Notes = nil                     // zero value
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.Coordinates = Convert_orderingtypes_Coordinates_addressmodel_Coordinates(arg.Coordinates, nil)
	out.Rid = 0 // zero value
}

func Convert_address_Addresses_addressmodel_Addresses(args []*address.Address) (outs []*addressmodel.Address) {
	tmps := make([]addressmodel.Address, len(args))
	outs = make([]*addressmodel.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_address_Address_addressmodel_Address(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/ordering/types.Address --//

func Convert_addressmodel_Address_orderingtypes_Address(arg *addressmodel.Address, out *orderingtypes.Address) *orderingtypes.Address {
	return OrderAddress(arg)
}

func convert_addressmodel_Address_orderingtypes_Address(arg *addressmodel.Address, out *orderingtypes.Address) {
	out.FullName = arg.FullName             // simple assign
	out.Phone = arg.Phone                   // simple assign
	out.Email = arg.Email                   // simple assign
	out.Company = arg.Company               // simple assign
	out.Address1 = arg.Address1             // simple assign
	out.Address2 = arg.Address2             // simple assign
	out.Location = orderingtypes.Location{} // zero value
}

func Convert_addressmodel_Addresses_orderingtypes_Addresses(args []*addressmodel.Address) (outs []*orderingtypes.Address) {
	tmps := make([]orderingtypes.Address, len(args))
	outs = make([]*orderingtypes.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_Address_orderingtypes_Address(args[i], &tmps[i])
	}
	return outs
}

func Convert_orderingtypes_Address_addressmodel_Address(arg *orderingtypes.Address, out *addressmodel.Address) *addressmodel.Address {
	return OrderAddressToModel(arg)
}

func convert_orderingtypes_Address_addressmodel_Address(arg *orderingtypes.Address, out *addressmodel.Address) {
	out.ID = 0                  // zero value
	out.FullName = arg.FullName // simple assign
	out.FirstName = ""          // zero value
	out.LastName = ""           // zero value
	out.Phone = arg.Phone       // simple assign
	out.Position = ""           // zero value
	out.Email = arg.Email       // simple assign
	out.Country = ""            // zero value
	out.City = ""               // zero value
	out.Province = ""           // zero value
	out.District = ""           // zero value
	out.Ward = ""               // zero value
	out.Zip = ""                // zero value
	out.DistrictCode = ""       // zero value
	out.ProvinceCode = ""       // zero value
	out.WardCode = ""           // zero value
	out.Company = arg.Company   // simple assign
	out.Address1 = arg.Address1 // simple assign
	out.Address2 = arg.Address2 // simple assign
	out.Type = ""               // zero value
	out.AccountID = 0           // zero value
	out.Notes = nil             // zero value
	out.CreatedAt = time.Time{} // zero value
	out.UpdatedAt = time.Time{} // zero value
	out.Coordinates = nil       // zero value
	out.Rid = 0                 // zero value
}

func Convert_orderingtypes_Addresses_addressmodel_Addresses(args []*orderingtypes.Address) (outs []*addressmodel.Address) {
	tmps := make([]addressmodel.Address, len(args))
	outs = make([]*addressmodel.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_orderingtypes_Address_addressmodel_Address(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/ordering/types.Coordinates --//

func Convert_addressmodel_Coordinates_orderingtypes_Coordinates(arg *addressmodel.Coordinates, out *orderingtypes.Coordinates) *orderingtypes.Coordinates {
	return Coordinates(arg)
}

func convert_addressmodel_Coordinates_orderingtypes_Coordinates(arg *addressmodel.Coordinates, out *orderingtypes.Coordinates) {
	out.Latitude = arg.Latitude   // simple assign
	out.Longitude = arg.Longitude // simple assign
}

func Convert_addressmodel_Coordinateses_orderingtypes_Coordinateses(args []*addressmodel.Coordinates) (outs []*orderingtypes.Coordinates) {
	tmps := make([]orderingtypes.Coordinates, len(args))
	outs = make([]*orderingtypes.Coordinates, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_Coordinates_orderingtypes_Coordinates(args[i], &tmps[i])
	}
	return outs
}

func Convert_orderingtypes_Coordinates_addressmodel_Coordinates(arg *orderingtypes.Coordinates, out *addressmodel.Coordinates) *addressmodel.Coordinates {
	return CoordinatesDB(arg)
}

func convert_orderingtypes_Coordinates_addressmodel_Coordinates(arg *orderingtypes.Coordinates, out *addressmodel.Coordinates) {
	out.Latitude = arg.Latitude   // simple assign
	out.Longitude = arg.Longitude // simple assign
}

func Convert_orderingtypes_Coordinateses_addressmodel_Coordinateses(args []*orderingtypes.Coordinates) (outs []*addressmodel.Coordinates) {
	tmps := make([]addressmodel.Coordinates, len(args))
	outs = make([]*addressmodel.Coordinates, len(args))
	for i := range tmps {
		outs[i] = Convert_orderingtypes_Coordinates_addressmodel_Coordinates(args[i], &tmps[i])
	}
	return outs
}
