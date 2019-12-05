package shipping_provider

import (
	"etop.vn/common/jsonx"
)

type ShippingProvider int

const (
	ShippingProvider_unknown ShippingProvider = 0
	ShippingProvider_all     ShippingProvider = 22
	ShippingProvider_manual  ShippingProvider = 20
	ShippingProvider_ghn     ShippingProvider = 19
	ShippingProvider_ghtk    ShippingProvider = 21
	ShippingProvider_vtpost  ShippingProvider = 23
)

var ShippingProvider_name = map[int]string{
	0:  "unknown",
	22: "all",
	20: "manual",
	19: "ghn",
	21: "ghtk",
	23: "vtpost",
}

var ShippingProvider_value = map[string]int{
	"unknown": 0,
	"all":     22,
	"manual":  20,
	"ghn":     19,
	"ghtk":    21,
	"vtpost":  23,
}

func (x ShippingProvider) Enum() *ShippingProvider {
	p := new(ShippingProvider)
	*p = x
	return p
}

func (x ShippingProvider) String() string {
	return jsonx.EnumName(ShippingProvider_name, int(x))
}

func (x *ShippingProvider) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(ShippingProvider_value, data, "ShippingProvider")
	if err != nil {
		return err
	}
	*x = ShippingProvider(value)
	return nil
}
