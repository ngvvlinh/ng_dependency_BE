package product_type

import (
	"etop.vn/common/jsonx"
)

// +enum
type ProductType int

const (
	// +enum=unknown
	ProductType_unknown ProductType = 0

	// +enum=services
	ProductType_services ProductType = 1

	// +enum=goods
	ProductType_goods ProductType = 2
)

var ProductType_name = map[int]string{
	0: "unknown",
	1: "services",
	2: "goods",
}

var ProductType_value = map[string]int{
	"unknown":  0,
	"services": 1,
	"goods":    2,
}

func (s ProductType) Enum() *ProductType {
	p := new(ProductType)
	*p = s
	return p
}

func (s ProductType) String() string {
	return jsonx.EnumName(ProductType_name, int(s))
}

func (s *ProductType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ProductType_value, data, "ProductType")
	if err != nil {
		return err
	}
	*s = ProductType(value)
	return nil
}
