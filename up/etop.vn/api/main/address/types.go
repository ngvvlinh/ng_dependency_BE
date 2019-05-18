package address

import (
	"time"

	orderingv1types "etop.vn/api/main/ordering/v1/types"
)

type Address struct {
	ID        int64
	FullName  string
	FirstName string
	LastName  string
	Phone     string
	Position  string
	Email     string

	Country  string
	City     string
	Province string
	District string
	Ward     string
	Zip      string

	DistrictCode string
	ProvinceCode string
	WardCode     string

	Company     string
	Address1    string
	Address2    string
	Type        string
	AccountID   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Coordinates *orderingv1types.Coordinates
}

func (a *Address) ToOrderAddress() *orderingv1types.Address {
	if a == nil {
		return nil
	}
	return &orderingv1types.Address{
		FullName: a.FullName,
		Phone:    a.Phone,
		Email:    a.Email,
		Company:  a.Company,
		Address1: a.Address1,
		Address2: a.Address2,
		Location: orderingv1types.Location{
			ProvinceCode: a.ProvinceCode,
			DistrictCode: a.DistrictCode,
			WardCode:     a.WardCode,
			Coordinates:  a.Coordinates,
		},
	}
}
