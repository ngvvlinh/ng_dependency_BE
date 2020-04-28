package types

import (
	"errors"
	"strings"

	"o.o/api/main/catalog/types"
	"o.o/api/main/location"
	"o.o/capi/dot"
)

type ItemLine struct {
	OrderID     dot.ID
	Quantity    int
	ProductID   dot.ID
	VariantID   dot.ID
	IsOutSide   bool
	ProductInfo ProductInfo
	TotalPrice  int
}

type ProductInfo struct {
	ProductName  string
	ImageURL     string
	Attributes   []*types.Attribute
	ListPrice    int
	RetailPrice  int
	PaymentPrice int
}

type Address struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	Location
}

type Location struct {
	ProvinceCode string       `json:"province_code"`
	Province     string       `json:"province"`
	DistrictCode string       `json:"district_code"`
	District     string       `json:"district"`
	WardCode     string       `json:"ward_code"`
	Ward         string       `json:"ward"`
	Coordinates  *Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// +enum
// +enum:sql=int
// +enum:zero=null
type ShippingType int

const (
	// +enum=none
	ShippingTypeNone ShippingType = 0

	// +enum=manual
	ShippingTypeManual ShippingType = 1

	// +enum=shipment
	ShippingTypeShipment ShippingType = 10

	// +enum=shipnow
	ShippingTypeShipnow ShippingType = 11
)

func ShippingTypeFromInt(s int) (ShippingType, error) {
	_, ok := enumShippingTypeName[s]
	if !ok {
		return 0, errors.New("invalid fulfill code")
	}
	return ShippingType(s), nil
}

func GetFullAddress(a *Address, location *location.LocationQueryResult) string {
	b := strings.Builder{}
	if a.Address1 != "" {
		b.WriteString(a.Address1)
		b.WriteByte('\n')
	}
	if a.Address2 != "" {
		b.WriteString(a.Address2)
		b.WriteByte('\n')
	}
	if a.Company != "" {
		b.WriteString(a.Company)
		b.WriteByte('\n')
	}
	flag := false
	if location.Ward != nil && location.Ward.Name != "" {
		b.WriteString(location.Ward.Name)
		flag = true
	}
	if location.District != nil && location.District.Name != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(location.District.Name)
		flag = true
	}
	if location.Province != nil && location.Province.Name != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(location.Province.Name)
	}
	return b.String()
}
