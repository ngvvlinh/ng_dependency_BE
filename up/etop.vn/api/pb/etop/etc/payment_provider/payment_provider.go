package payment_provider

import "etop.vn/api/external/payment"

func (s *PaymentProvider) ToPaymentProvider() payment.PaymentProvider {
	if s == nil || *s == 0 {
		return ""
	}
	return payment.PaymentProvider(PaymentProvider_name[int32(*s)])
}

func (x PaymentProvider) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
