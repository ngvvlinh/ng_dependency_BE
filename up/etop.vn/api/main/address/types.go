package address

import (
	"strings"
	"time"

	orderv1types "etop.vn/api/main/ordering/v1/types"
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
	Coordinates *orderv1types.Coordinates
}

func (a *Address) ToOrderAddress() *orderv1types.Address {
	if a == nil {
		return nil
	}
	return &orderv1types.Address{
		FullName: a.FullName,
		Phone:    a.Phone,
		Email:    a.Email,
		Company:  a.Company,
		Address1: a.Address1,
		Address2: a.Address2,
		Location: orderv1types.Location{
			ProvinceCode: a.ProvinceCode,
			DistrictCode: a.DistrictCode,
			WardCode:     a.WardCode,
			Coordinates:  a.Coordinates,
		},
	}
}

func (a *Address) GetFullAddress() string {
	b := strings.Builder{}
	if a.Address1 != "" {
		b.WriteString(a.Address1)
		b.WriteByte('\n')
	}
	if a.Address2 != "" {
		b.WriteString(a.Address2)
		b.WriteByte('\n')
	}
	flag := false
	if a.Ward != "" {
		b.WriteString(a.Ward)
		flag = true
	}
	if a.District != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(a.District)
		flag = true
	}
	if a.Province != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(a.Province)
	}
	return b.String()
}
