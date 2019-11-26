package address_type

import (
	"etop.vn/common/jsonx"
)

type AddressType int32

const (
	AddressType_unknown   AddressType = 0
	AddressType_general   AddressType = 1
	AddressType_warehouse AddressType = 2
	AddressType_shipto    AddressType = 3
	AddressType_shipfrom  AddressType = 4
)

var AddressType_name = map[int32]string{
	0: "unknown",
	1: "general",
	2: "warehouse",
	3: "shipto",
	4: "shipfrom",
}

var AddressType_value = map[string]int32{
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
	return jsonx.EnumName(AddressType_name, int32(x))
}

func (x *AddressType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(AddressType_value, data, "AddressType")
	if err != nil {
		return err
	}
	*x = AddressType(value)
	return nil
}
