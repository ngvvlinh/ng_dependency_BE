package types

import (
	"errors"
	"strings"

	"etop.vn/api/main/catalog/types"
	"etop.vn/api/main/location"
)

type ItemLine struct {
	OrderId     int64
	Quantity    int32
	ProductId   int64
	VariantId   int64
	IsOutside   bool
	ProductInfo ProductInfo
	TotalPrice  int32
}

type ProductInfo struct {
	ProductName  string
	ImageUrl     string
	Attributes   []*types.Attribute
	ListPrice    int32
	RetailPrice  int32
	PaymentPrice int32
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

type Fulfill int32

const (
	FulfillNone     Fulfill = 0
	FulfillManual   Fulfill = 1
	FulfillShipment Fulfill = 10
	FulfillShipnow  Fulfill = 11
)

var Fulfill_name = map[int32]string{
	0:  "none",
	1:  "manual",
	10: "shipment",
	11: "shipnow",
}

var Fulfill_value = map[string]int32{
	"none":     0,
	"manual":   1,
	"shipment": 10,
	"shipnow":  11,
}

func FulfillFromInt(s int32) (Fulfill, error) {
	_, ok := Fulfill_name[s]
	if !ok {
		return 0, errors.New("invalid fulfill code")
	}
	return Fulfill(s), nil
}

func (f Fulfill) String() string {
	return Fulfill_name[int32(f)]
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
