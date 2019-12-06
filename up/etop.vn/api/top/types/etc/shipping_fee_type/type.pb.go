package shipping_fee_type

import (
	"etop.vn/common/jsonx"
)

// +enum
type ShippingFeeType int

const (
	// +enum=main
	ShippingFeeType_main ShippingFeeType = 0

	// +enum=return
	ShippingFeeType_return ShippingFeeType = 1

	// +enum=adjustment
	ShippingFeeType_adjustment ShippingFeeType = 2

	// +enum=insurance
	ShippingFeeType_insurance ShippingFeeType = 3

	// +enum=tax
	ShippingFeeType_tax ShippingFeeType = 4

	// +enum=other
	ShippingFeeType_other ShippingFeeType = 5

	// +enum=cods
	ShippingFeeType_cods ShippingFeeType = 6

	// +enum=address_change
	ShippingFeeType_address_change ShippingFeeType = 7

	// +enum=discount
	ShippingFeeType_discount ShippingFeeType = 8

	// +enum=unknown
	ShippingFeeType_unknown ShippingFeeType = 127
)

var ShippingFeeType_name = map[int]string{
	0:   "main",
	1:   "return",
	2:   "adjustment",
	3:   "insurance",
	4:   "tax",
	5:   "other",
	6:   "cods",
	7:   "address_change",
	8:   "discount",
	127: "unknown",
}

var ShippingFeeType_value = map[string]int{
	"main":           0,
	"return":         1,
	"adjustment":     2,
	"insurance":      3,
	"tax":            4,
	"other":          5,
	"cods":           6,
	"address_change": 7,
	"discount":       8,
	"unknown":        127,
}

func (x ShippingFeeType) Enum() *ShippingFeeType {
	p := new(ShippingFeeType)
	*p = x
	return p
}

func (x ShippingFeeType) String() string {
	return jsonx.EnumName(ShippingFeeType_name, int(x))
}

func (x *ShippingFeeType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ShippingFeeType_value, data, "ShippingFeeType")
	if err != nil {
		return err
	}
	*x = ShippingFeeType(value)
	return nil
}
