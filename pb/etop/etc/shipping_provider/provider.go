package shipping_provider

import "etop.vn/backend/pkg/etop/model"

func (s *ShippingProvider) ToModel() model.ShippingProvider {
	if s == nil || *s == 0 {
		return ""
	}
	return model.ShippingProvider(ShippingProvider_name[int32(*s)])
}

func PbShippingProviderType(sp model.ShippingProvider) ShippingProvider {
	spString := string(sp)
	return ShippingProvider(ShippingProvider_value[spString])
}

func PbPtrShippingProvider(sp model.ShippingProvider) *ShippingProvider {
	res := PbShippingProviderType(sp)
	return &res
}

func PbShippingProviderPtr(s *string) *ShippingProvider {
	if s == nil || *s == "" {
		return nil
	}
	sp := PbShippingProviderType(model.ShippingProvider(*s))
	return &sp
}
