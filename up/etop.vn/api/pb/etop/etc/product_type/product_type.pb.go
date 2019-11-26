package product_type

import (
	"etop.vn/common/jsonx"
)

type ProductType int32

const (
	ProductType_unknown  ProductType = 0
	ProductType_services ProductType = 1
	ProductType_goods    ProductType = 2
)

var ProductType_name = map[int32]string{
	0: "unknown",
	1: "services",
	2: "goods",
}

var ProductType_value = map[string]int32{
	"unknown":  0,
	"services": 1,
	"goods":    2,
}

func (x ProductType) Enum() *ProductType {
	p := new(ProductType)
	*p = x
	return p
}

func (x ProductType) String() string {
	return jsonx.EnumName(ProductType_name, int32(x))
}

func (x *ProductType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ProductType_value, data, "ProductType")
	if err != nil {
		return err
	}
	*x = ProductType(value)
	return nil
}
