package address_type

import (
	"etop.vn/common/jsonx"
)

// +enum
type AddressType int

const (
	// +enum=unknown
	AddressType_unknown AddressType = 0

	// +enum=general
	AddressType_general AddressType = 1

	// +enum=warehouse
	AddressType_warehouse AddressType = 2

	// +enum=shipto
	AddressType_shipto AddressType = 3

	// +enum=shipfrom
	AddressType_shipfrom AddressType = 4
)

var AddressType_name = map[int]string{
	0: "unknown",
	1: "general",
	2: "warehouse",
	3: "shipto",
	4: "shipfrom",
}

var AddressType_value = map[string]int{
	"unknown":   0,
	"general":   1,
	"warehouse": 2,
	"shipto":    3,
	"shipfrom":  4,
}

func (x AddressType) Enum() *AddressType {
	p := new(AddressType)
	*p = x
	return p
}

func (x AddressType) String() string {
	return jsonx.EnumName(AddressType_name, int(x))
}

func (x *AddressType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(AddressType_value, data, "AddressType")
	if err != nil {
		return err
	}
	*x = AddressType(value)
	return nil
}
