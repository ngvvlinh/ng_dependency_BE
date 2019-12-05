package payment_provider

import "etop.vn/api/external/payment"

func (x *PaymentProvider) ToPaymentProvider() payment.PaymentProvider {
	if x == nil || *x == 0 {
		return ""
	}
	return payment.PaymentProvider(PaymentProvider_name[int(*x)])
}

func (x PaymentProvider) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
