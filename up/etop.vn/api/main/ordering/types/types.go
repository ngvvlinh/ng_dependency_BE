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
	DistrictCode string
	WardCode     string
	Coordinates  *Coordinates
}

type Coordinates struct {
	Latitude  float32
	Longitude float32
}

// +enum
type Fulfill int

const (
	// +enum=none
	FulfillNone Fulfill = 0

	// +enum=manual
	FulfillManual Fulfill = 1

	// +enum=shipment
	FulfillShipment Fulfill = 10

	// +enum=shipnow
	FulfillShipnow Fulfill = 11
)

func FulfillFromInt(s int) (Fulfill, error) {
	_, ok := enumFulfillName[s]
	if !ok {
		return 0, errors.New("invalid fulfill code")
	}
	return Fulfill(s), nil
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
