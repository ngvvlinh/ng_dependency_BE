// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	addressmodel "o.o/backend/com/main/address/model"
	conversion "o.o/backend/pkg/common/conversion"
	addressmodel1 "o.o/backend/zexp/etl/main/address/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*addressmodel1.Address)(nil), (*addressmodel.Address)(nil), func(arg, out interface{}) error {
		Convert_addressmodel1_Address_addressmodel_Address(arg.(*addressmodel1.Address), out.(*addressmodel.Address))
		return nil
	})
	s.Register(([]*addressmodel1.Address)(nil), (*[]*addressmodel.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel1_Addresses_addressmodel_Addresses(arg.([]*addressmodel1.Address))
		*out.(*[]*addressmodel.Address) = out0
		return nil
	})
	s.Register((*addressmodel.Address)(nil), (*addressmodel1.Address)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_Address_addressmodel1_Address(arg.(*addressmodel.Address), out.(*addressmodel1.Address))
		return nil
	})
	s.Register(([]*addressmodel.Address)(nil), (*[]*addressmodel1.Address)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_Addresses_addressmodel1_Addresses(arg.([]*addressmodel.Address))
		*out.(*[]*addressmodel1.Address) = out0
		return nil
	})
	s.Register((*addressmodel1.AddressNote)(nil), (*addressmodel.AddressNote)(nil), func(arg, out interface{}) error {
		Convert_addressmodel1_AddressNote_addressmodel_AddressNote(arg.(*addressmodel1.AddressNote), out.(*addressmodel.AddressNote))
		return nil
	})
	s.Register(([]*addressmodel1.AddressNote)(nil), (*[]*addressmodel.AddressNote)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel1_AddressNotes_addressmodel_AddressNotes(arg.([]*addressmodel1.AddressNote))
		*out.(*[]*addressmodel.AddressNote) = out0
		return nil
	})
	s.Register((*addressmodel.AddressNote)(nil), (*addressmodel1.AddressNote)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_AddressNote_addressmodel1_AddressNote(arg.(*addressmodel.AddressNote), out.(*addressmodel1.AddressNote))
		return nil
	})
	s.Register(([]*addressmodel.AddressNote)(nil), (*[]*addressmodel1.AddressNote)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_AddressNotes_addressmodel1_AddressNotes(arg.([]*addressmodel.AddressNote))
		*out.(*[]*addressmodel1.AddressNote) = out0
		return nil
	})
	s.Register((*addressmodel1.Coordinates)(nil), (*addressmodel.Coordinates)(nil), func(arg, out interface{}) error {
		Convert_addressmodel1_Coordinates_addressmodel_Coordinates(arg.(*addressmodel1.Coordinates), out.(*addressmodel.Coordinates))
		return nil
	})
	s.Register(([]*addressmodel1.Coordinates)(nil), (*[]*addressmodel.Coordinates)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel1_Coordinateses_addressmodel_Coordinateses(arg.([]*addressmodel1.Coordinates))
		*out.(*[]*addressmodel.Coordinates) = out0
		return nil
	})
	s.Register((*addressmodel.Coordinates)(nil), (*addressmodel1.Coordinates)(nil), func(arg, out interface{}) error {
		Convert_addressmodel_Coordinates_addressmodel1_Coordinates(arg.(*addressmodel.Coordinates), out.(*addressmodel1.Coordinates))
		return nil
	})
	s.Register(([]*addressmodel.Coordinates)(nil), (*[]*addressmodel1.Coordinates)(nil), func(arg, out interface{}) error {
		out0 := Convert_addressmodel_Coordinateses_addressmodel1_Coordinateses(arg.([]*addressmodel.Coordinates))
		*out.(*[]*addressmodel1.Coordinates) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/address/model.Address --//

func Convert_addressmodel1_Address_addressmodel_Address(arg *addressmodel1.Address, out *addressmodel.Address) *addressmodel.Address {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel.Address{}
	}
	convert_addressmodel1_Address_addressmodel_Address(arg, out)
	return out
}

func convert_addressmodel1_Address_addressmodel_Address(arg *addressmodel1.Address, out *addressmodel.Address) {
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
	out.IsDefault = false               // zero value
	out.DistrictCode = arg.DistrictCode // simple assign
	out.ProvinceCode = arg.ProvinceCode // simple assign
	out.WardCode = arg.WardCode         // simple assign
	out.Company = arg.Company           // simple assign
	out.Address1 = arg.Address1         // simple assign
	out.Address2 = arg.Address2         // simple assign
	out.Type = arg.Type                 // simple assign
	out.AccountID = arg.AccountID       // simple assign
	out.Notes = Convert_addressmodel1_AddressNote_addressmodel_AddressNote(arg.Notes, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.Coordinates = Convert_addressmodel1_Coordinates_addressmodel_Coordinates(arg.Coordinates, nil)
	out.Rid = arg.Rid // simple assign
}

func Convert_addressmodel1_Addresses_addressmodel_Addresses(args []*addressmodel1.Address) (outs []*addressmodel.Address) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel.Address, len(args))
	outs = make([]*addressmodel.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel1_Address_addressmodel_Address(args[i], &tmps[i])
	}
	return outs
}

func Convert_addressmodel_Address_addressmodel1_Address(arg *addressmodel.Address, out *addressmodel1.Address) *addressmodel1.Address {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel1.Address{}
	}
	convert_addressmodel_Address_addressmodel1_Address(arg, out)
	return out
}

func convert_addressmodel_Address_addressmodel1_Address(arg *addressmodel.Address, out *addressmodel1.Address) {
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
	out.Notes = Convert_addressmodel_AddressNote_addressmodel1_AddressNote(arg.Notes, nil)
	out.CreatedAt = arg.CreatedAt // simple assign
	out.UpdatedAt = arg.UpdatedAt // simple assign
	out.Coordinates = Convert_addressmodel_Coordinates_addressmodel1_Coordinates(arg.Coordinates, nil)
	out.Rid = arg.Rid // simple assign
}

func Convert_addressmodel_Addresses_addressmodel1_Addresses(args []*addressmodel.Address) (outs []*addressmodel1.Address) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel1.Address, len(args))
	outs = make([]*addressmodel1.Address, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_Address_addressmodel1_Address(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/address/model.AddressNote --//

func Convert_addressmodel1_AddressNote_addressmodel_AddressNote(arg *addressmodel1.AddressNote, out *addressmodel.AddressNote) *addressmodel.AddressNote {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel.AddressNote{}
	}
	convert_addressmodel1_AddressNote_addressmodel_AddressNote(arg, out)
	return out
}

func convert_addressmodel1_AddressNote_addressmodel_AddressNote(arg *addressmodel1.AddressNote, out *addressmodel.AddressNote) {
	out.Note = arg.Note             // simple assign
	out.OpenTime = arg.OpenTime     // simple assign
	out.LunchBreak = arg.LunchBreak // simple assign
	out.Other = arg.Other           // simple assign
}

func Convert_addressmodel1_AddressNotes_addressmodel_AddressNotes(args []*addressmodel1.AddressNote) (outs []*addressmodel.AddressNote) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel.AddressNote, len(args))
	outs = make([]*addressmodel.AddressNote, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel1_AddressNote_addressmodel_AddressNote(args[i], &tmps[i])
	}
	return outs
}

func Convert_addressmodel_AddressNote_addressmodel1_AddressNote(arg *addressmodel.AddressNote, out *addressmodel1.AddressNote) *addressmodel1.AddressNote {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel1.AddressNote{}
	}
	convert_addressmodel_AddressNote_addressmodel1_AddressNote(arg, out)
	return out
}

func convert_addressmodel_AddressNote_addressmodel1_AddressNote(arg *addressmodel.AddressNote, out *addressmodel1.AddressNote) {
	out.Note = arg.Note             // simple assign
	out.OpenTime = arg.OpenTime     // simple assign
	out.LunchBreak = arg.LunchBreak // simple assign
	out.Other = arg.Other           // simple assign
}

func Convert_addressmodel_AddressNotes_addressmodel1_AddressNotes(args []*addressmodel.AddressNote) (outs []*addressmodel1.AddressNote) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel1.AddressNote, len(args))
	outs = make([]*addressmodel1.AddressNote, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_AddressNote_addressmodel1_AddressNote(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/address/model.Coordinates --//

func Convert_addressmodel1_Coordinates_addressmodel_Coordinates(arg *addressmodel1.Coordinates, out *addressmodel.Coordinates) *addressmodel.Coordinates {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel.Coordinates{}
	}
	convert_addressmodel1_Coordinates_addressmodel_Coordinates(arg, out)
	return out
}

func convert_addressmodel1_Coordinates_addressmodel_Coordinates(arg *addressmodel1.Coordinates, out *addressmodel.Coordinates) {
	out.Latitude = arg.Latitude   // simple assign
	out.Longitude = arg.Longitude // simple assign
}

func Convert_addressmodel1_Coordinateses_addressmodel_Coordinateses(args []*addressmodel1.Coordinates) (outs []*addressmodel.Coordinates) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel.Coordinates, len(args))
	outs = make([]*addressmodel.Coordinates, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel1_Coordinates_addressmodel_Coordinates(args[i], &tmps[i])
	}
	return outs
}

func Convert_addressmodel_Coordinates_addressmodel1_Coordinates(arg *addressmodel.Coordinates, out *addressmodel1.Coordinates) *addressmodel1.Coordinates {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &addressmodel1.Coordinates{}
	}
	convert_addressmodel_Coordinates_addressmodel1_Coordinates(arg, out)
	return out
}

func convert_addressmodel_Coordinates_addressmodel1_Coordinates(arg *addressmodel.Coordinates, out *addressmodel1.Coordinates) {
	out.Latitude = arg.Latitude   // simple assign
	out.Longitude = arg.Longitude // simple assign
}

func Convert_addressmodel_Coordinateses_addressmodel1_Coordinateses(args []*addressmodel.Coordinates) (outs []*addressmodel1.Coordinates) {
	if args == nil {
		return nil
	}
	tmps := make([]addressmodel1.Coordinates, len(args))
	outs = make([]*addressmodel1.Coordinates, len(args))
	for i := range tmps {
		outs[i] = Convert_addressmodel_Coordinates_addressmodel1_Coordinates(args[i], &tmps[i])
	}
	return outs
}
