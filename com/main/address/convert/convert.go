package convert

import (
	"etop.vn/api/main/address"
	orderingv1types "etop.vn/api/main/ordering/types"
	"etop.vn/backend/pkg/etop/model"
)

func AddressToModel(in *address.Address) (out *model.Address) {
	if in == nil {
		return nil
	}
	out = &model.Address{
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
		Coordinates:  CoordinatesToModel(in.Coordinates),
	}
	return
}

func Address(in *model.Address) (out *address.Address) {
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

func CoordinatesToModel(in *orderingv1types.Coordinates) (out *model.Coordinates) {
	if in == nil {
		return nil
	}
	out = &model.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
	return
}

func Coordinates(in *model.Coordinates) (out *orderingv1types.Coordinates) {
	if in == nil {
		return nil
	}
	out = &orderingv1types.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
	return
}
