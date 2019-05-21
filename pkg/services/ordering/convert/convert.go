package convert

import (
	"time"

	"etop.vn/api/main/ordering"
	"etop.vn/api/main/ordering/v1/types"
	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

func AddressToModel(in *types.Address) (out *model.Address) {
	if in == nil {
		return nil
	}
	out = &model.Address{
		ID:           0,
		FullName:     in.FullName,
		FirstName:    "",
		LastName:     "",
		Phone:        in.Phone,
		Position:     "",
		Email:        in.Email,
		Country:      "",
		City:         "",
		Province:     "",
		District:     "",
		Ward:         "",
		Zip:          "",
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      "",
		Address1:     in.Address1,
		Address2:     in.Address2,
		Type:         "",
		AccountID:    0,
		Notes:        nil,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
	return out
}

func Address(in *model.Address) (out *types.Address) {
	if in == nil {
		return nil
	}
	out = &types.Address{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Company:  in.Company,
		Address1: in.Address1,
		Address2: in.Address2,
		Location: types.Location{
			ProvinceCode: in.ProvinceCode,
			DistrictCode: in.DistrictCode,
			WardCode:     in.WardCode,
		},
	}
	return out
}

func Order(in *ordermodel.Order) (out *ordering.Order) {
	if in == nil {
		return nil
	}
	out = &ordering.Order{
		ID: in.ID,
	}
	return out
}
