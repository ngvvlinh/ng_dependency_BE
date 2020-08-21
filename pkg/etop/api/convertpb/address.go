package convertpb

import (
	"o.o/api/main/address"
	orderv1types "o.o/api/main/ordering/types"
	api "o.o/api/top/int/etop"
	"o.o/capi/dot"
)

func Convert_core_Address_To_api_Address(in *address.Address) *api.Address {
	if in == nil {
		return nil
	}
	return &api.Address{
		Province:     in.Province,
		ProvinceCode: in.ProvinceCode,
		District:     in.District,
		DistrictCode: in.DistrictCode,
		Ward:         in.Ward,
		WardCode:     in.WardCode,
		Address1:     in.Address1,
		Address2:     in.Address2,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Type:         in.Type,
		Country:      in.Country,
		FirstName:    in.FirstName,
		Id:           in.ID,
		LastName:     in.LastName,
		Position:     in.Position,
		Zip:          in.Zip,
		Coordinates:  Convert_core_Coordinates_To_api_Coordinates(in.Coordinates),
		Notes:        Convert_core_AddressNote_To_api_AddressNote(in.Notes),
	}
}

func Convert_core_AddressNote_To_api_AddressNote(in *address.AddressNote) *api.AddressNote {
	if in == nil {
		return nil
	}
	return &api.AddressNote{
		LunchBreak: in.LunchBreak,
		Note:       in.Note,
		OpenTime:   in.OpenTime,
		Other:      in.Other,
	}
}

func Convert_core_Addresses_To_api_Addresses(items []*address.Address) []*api.Address {
	result := make([]*api.Address, len(items))
	for i, item := range items {
		result[i] = Convert_core_Address_To_api_Address(item)
	}
	return result
}

func SetCreateAddressArgs(in *api.CreateAddressRequest, accountID dot.ID) *address.CreateAddressCommand {
	var coordinates orderv1types.Coordinates
	var notes address.AddressNote
	if in.Coordinates != nil {
		coordinates.Latitude = in.Coordinates.Latitude
		coordinates.Longitude = in.Coordinates.Longitude
	}

	if in.Notes != nil {
		notes.LunchBreak = in.Notes.LunchBreak
		notes.Note = in.Notes.Note
		notes.OpenTime = in.Notes.OpenTime
		notes.Other = in.Notes.Other
	}
	return &address.CreateAddressCommand{
		Address2:     in.Address2,
		Address1:     in.Address1,
		Country:      in.Country,
		District:     in.District,
		DistrictCode: in.DistrictCode,
		Email:        in.Email,
		FirstName:    in.FirstName,
		FullName:     in.FullName,
		LastName:     in.LastName,
		Phone:        in.Phone,
		Province:     in.Province,
		Position:     in.Position,
		ProvinceCode: in.ProvinceCode,
		Ward:         in.Ward,
		WardCode:     in.WardCode,
		Zip:          in.Zip,
		AccountID:    accountID,
		Coordinates:  &coordinates,
		Notes:        &notes,
		Type:         in.Type,
	}
}

func SetUpdateAddressArgs(in *api.UpdateAddressRequest, accountID dot.ID) *address.UpdateAddressCommand {
	var coordinates orderv1types.Coordinates
	var notes address.AddressNote
	if in.Coordinates != nil {
		coordinates.Latitude = in.Coordinates.Latitude
		coordinates.Longitude = in.Coordinates.Longitude
	}

	if in.Notes != nil {
		notes.LunchBreak = in.Notes.LunchBreak
		notes.Note = in.Notes.Note
		notes.OpenTime = in.Notes.OpenTime
		notes.Other = in.Notes.Other
	}
	return &address.UpdateAddressCommand{
		AccountID:    accountID,
		Address1:     in.Address1,
		Address2:     in.Address2,
		Country:      in.Country,
		District:     in.District,
		DistrictCode: in.DistrictCode,
		Email:        in.Email,
		Zip:          in.Zip,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		FullName:     in.FullName,
		ID:           in.Id,
		Phone:        in.Phone,
		Position:     in.Position,
		Province:     in.Province,
		ProvinceCode: in.ProvinceCode,
		Ward:         in.Ward,
		WardCode:     in.WardCode,
		Type:         in.Type,
		Coordinates:  &coordinates,
		Notes:        &notes,
	}
}
