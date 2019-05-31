package types

import (
	"errors"
	"strings"

	location "etop.vn/api/main/location"
	v1types "etop.vn/api/main/ordering/v1/types"
)

type ItemLine = v1types.ItemLine
type ProductInfo = v1types.ProductInfo
type Address = v1types.Address
type Location = v1types.Location
type Coordinates = v1types.Coordinates

type Fulfill = v1types.Fulfill

const (
	FulfillNone     = v1types.Fulfill_none
	FulfillManual   = v1types.Fulfill_manual
	FulfillShipment = v1types.Fulfill_shipment
	FulfillShipnow  = v1types.Fulfill_shipnow
)

func FulfillFromInt(s int32) (Fulfill, error) {
	_, ok := v1types.Fulfill_name[s]
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
