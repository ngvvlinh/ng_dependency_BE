package payment_provider

import (
	"etop.vn/common/jsonx"
)

type PaymentProvider int32

const (
	PaymentProvider_unknown PaymentProvider = 0
	PaymentProvider_vtpay   PaymentProvider = 1
)

var PaymentProvider_name = map[int32]string{
	0: "unknown",
	1: "vtpay",
}

var PaymentProvider_value = map[string]int32{
	"unknown": 0,
	"vtpay":   1,
}

func (x PaymentProvider) Enum() *PaymentProvider {
	p := new(PaymentProvider)
	*p = x
	return p
}

func (x PaymentProvider) String() string {
	return jsonx.EnumName(PaymentProvider_name, int32(x))
}

func (x *PaymentProvider) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(PaymentProvider_value, data, "PaymentProvider")
	if err != nil {
		return err
	}
	*x = PaymentProvider(value)
	return nil
}
