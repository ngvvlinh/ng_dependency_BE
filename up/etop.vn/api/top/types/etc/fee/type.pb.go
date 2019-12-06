package fee

import (
	"etop.vn/common/jsonx"
)

// +enum
type FeeType int

const (
	// +enum=other
	FeeType_other FeeType = 0

	// +enum=shipping
	FeeType_shipping FeeType = 1

	// +enum=tax
	FeeType_tax FeeType = 2
)

var FeeType_name = map[int]string{
	0: "other",
	1: "shipping",
	2: "tax",
}

var FeeType_value = map[string]int{
	"other":    0,
	"shipping": 1,
	"tax":      2,
}

func (x FeeType) Enum() *FeeType {
	p := new(FeeType)
	*p = x
	return p
}

func (x FeeType) String() string {
	return jsonx.EnumName(FeeType_name, int(x))
}

func (x *FeeType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(FeeType_value, data, "FeeType")
	if err != nil {
		return err
	}
	*x = FeeType(value)
	return nil
}
