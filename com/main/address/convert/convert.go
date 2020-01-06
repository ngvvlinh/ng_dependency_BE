package convert

import (
	"etop.vn/api/main/address"
	ordertypes "etop.vn/api/main/ordering/types"
	addressmodel "etop.vn/backend/com/main/address/model"
)

// +gen:convert: etop.vn/backend/com/main/address/model->etop.vn/api/main/ordering/types
// +gen:convert: etop.vn/backend/com/main/address/model->etop.vn/api/main/address

func AddressToModel(in *address.Address) (out *addressmodel.Address) {
	if in == nil {
		return nil
	}
	out = &addressmodel.Address{
		ID:           in.ID,
		FullName:     in.FullName,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Phone:        in.Phone,
		Position:     in.Position,
		Email:        in.Email,
		Country:      in.Country,
		City:         in.City,
		Province:     in.Province,
		District:     in.District,
		Ward:         in.Ward,
		Zip:          in.Zip,
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
		Type:         in.Type,
		AccountID:    in.AccountID,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Coordinates:  CoordinatesDB(in.Coordinates),
	}
	return
}

func Address(in *addressmodel.Address) (out *address.Address) {
	if in == nil {
		return nil
	}
	out = &address.Address{
		ID:           in.ID,
		FullName:     in.FullName,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Phone:        in.Phone,
		Position:     in.Position,
		Email:        in.Email,
		Country:      in.Country,
		City:         in.City,
		Province:     in.Province,
		District:     in.District,
		Ward:         in.Ward,
		Zip:          in.Zip,
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
		Type:         in.Type,
		AccountID:    in.AccountID,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Coordinates:  Coordinates(in.Coordinates),
	}
	return
}

func CoordinatesDB(in *ordertypes.Coordinates) (out *addressmodel.Coordinates) {
	if in == nil {
		return nil
	}
	out = &addressmodel.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
	return
}

func Coordinates(in *addressmodel.Coordinates) (out *ordertypes.Coordinates) {
	if in == nil {
		return nil
	}
	out = &ordertypes.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
	return
}

func OrderAddressToModel(in *ordertypes.Address) *addressmodel.Address {
	if in == nil {
		return nil
	}
	res := &addressmodel.Address{
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Province:     in.Province,
		District:     in.District,
		Ward:         in.Ward,
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      in.Company,
		Address1:     in.Address1,
		Address2:     in.Address2,
	}
	if in.Coordinates != nil {
		res.Coordinates = &addressmodel.Coordinates{
			Latitude:  in.Coordinates.Latitude,
			Longitude: in.Coordinates.Longitude,
		}
	}
	return res
}

func OrderAddress(in *addressmodel.Address) *ordertypes.Address {
	if in == nil {
		return nil
	}
	out := &ordertypes.Address{}
	convert_addressmodel_Address_orderingtypes_Address(in, out)
	out.Location = ordertypes.Location{
		ProvinceCode: in.ProvinceCode,
		Province:     in.Province,
		DistrictCode: in.DistrictCode,
		District:     in.District,
		WardCode:     in.WardCode,
		Ward:         in.Ward,
		Coordinates:  Convert_addressmodel_Coordinates_orderingtypes_Coordinates(in.Coordinates, nil),
	}
	return out
}
