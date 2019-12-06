package payment_provider

import (
	"etop.vn/common/jsonx"
)

// +enum
type PaymentProvider int

const (
	// +enum=unknown
	PaymentProvider_unknown PaymentProvider = 0

	// +enum=vtpay
	PaymentProvider_vtpay PaymentProvider = 1
)

var PaymentProvider_name = map[int]string{
	0: "unknown",
	1: "vtpay",
}

var PaymentProvider_value = map[string]int{
	"unknown": 0,
	"vtpay":   1,
}

func (x PaymentProvider) Enum() *PaymentProvider {
	p := new(PaymentProvider)
	*p = x
	return p
}

func (x PaymentProvider) String() string {
	return jsonx.EnumName(PaymentProvider_name, int(x))
}

func (x *PaymentProvider) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(PaymentProvider_value, data, "PaymentProvider")
	if err != nil {
		return err
	}
	*x = PaymentProvider(value)
	return nil
}
