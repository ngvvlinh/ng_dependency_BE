package credit_type

import (
	"etop.vn/common/jsonx"
)

// +enum
type CreditType int

const (
	// +enum=shop
	CreditType_shop CreditType = 1
)

var CreditType_name = map[int]string{
	1: "shop",
}

var CreditType_value = map[string]int{
	"shop": 1,
}

func (x CreditType) Enum() *CreditType {
	p := new(CreditType)
	*p = x
	return p
}

func (x CreditType) String() string {
	return jsonx.EnumName(CreditType_name, int(x))
}

func (x *CreditType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(CreditType_value, data, "CreditType")
	if err != nil {
		return err
	}
	*x = CreditType(value)
	return nil
}
