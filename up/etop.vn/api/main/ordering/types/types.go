package types

import (
	"errors"
	"strings"

	"etop.vn/api/main/catalog/types"
	"etop.vn/api/main/location"
	"etop.vn/capi/dot"
)

type ItemLine struct {
	OrderId     dot.ID
	Quantity    int
	ProductId   dot.ID
	VariantId   dot.ID
	IsOutside   bool
	ProductInfo ProductInfo
	TotalPrice  int
}

type ProductInfo struct {
	ProductName  string
	ImageUrl     string
	Attributes   []*types.Attribute
	ListPrice    int
	RetailPrice  int
	PaymentPrice int
}

type Address struct {
	FullName string
	Phone    string
	Email    string
	Company  string
	Address1 string
	Address2 string
	Location
}

type Location struct {
	ProvinceCode string
	Province     string
	DistrictCode string
	District     string
	WardCode     string
	Ward         string
	Coordinates  *Coordinates
}

type Coordinates struct {
	Latitude  float32
	Longitude float32
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
