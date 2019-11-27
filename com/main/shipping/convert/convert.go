package convert

import (
	ordertypes "etop.vn/api/main/ordering/types"
	etopmodel "etop.vn/backend/pkg/etop/model"
)

func ModelAddress(in *ordertypes.Address) *etopmodel.Address {
	if in == nil {
		return nil
	}
	res := &etopmodel.Address{
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
		res.Coordinates = &etopmodel.Coordinates{
			Latitude:  in.Coordinates.Latitude,
			Longitude: in.Coordinates.Longitude,
		}
	}
	return res
}
