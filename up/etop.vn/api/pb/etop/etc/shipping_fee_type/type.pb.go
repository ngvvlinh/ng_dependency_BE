package shipping_fee_type

import (
	"etop.vn/common/jsonx"
)

type ShippingFeeType int32

const (
	ShippingFeeType_main           ShippingFeeType = 0
	ShippingFeeType_return         ShippingFeeType = 1
	ShippingFeeType_adjustment     ShippingFeeType = 2
	ShippingFeeType_insurance      ShippingFeeType = 3
	ShippingFeeType_tax            ShippingFeeType = 4
	ShippingFeeType_other          ShippingFeeType = 5
	ShippingFeeType_cods           ShippingFeeType = 6
	ShippingFeeType_address_change ShippingFeeType = 7
	ShippingFeeType_discount       ShippingFeeType = 8
	ShippingFeeType_unknown        ShippingFeeType = 127
)

var ShippingFeeType_name = map[int32]string{
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

var ShippingFeeType_value = map[string]int32{
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
	return jsonx.EnumName(ShippingFeeType_name, int32(x))
}

func (x *ShippingFeeType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ShippingFeeType_value, data, "ShippingFeeType")
	if err != nil {
		return err
	}
	*x = ShippingFeeType(value)
	return nil
}
