package types

import (
	"errors"

	v1types "etop.vn/api/main/ordering/v1/types"
)

type ItemLine = v1types.ItemLine
type ProductInfo = v1types.ProductInfo
type Address = v1types.Address
type Location = v1types.Location
type Coordinates = v1types.Coordinates

type Fulfill = v1types.Fulfill

const (
	FulfillNone               = v1types.Fulfill_none
	FulfillManual             = v1types.Fulfill_manual
	FulfillFulfillment        = v1types.Fulfill_fulfillment
	FulfillShipnowFulfillment = v1types.Fulfill_shipnow_fulfillment
)

func FulfillFromInt(s int32) (Fulfill, error) {
	_, ok := v1types.Fulfill_name[s]
	if !ok {
		return 0, errors.New("invalid fulfill code")
	}
	return Fulfill(s), nil
}
