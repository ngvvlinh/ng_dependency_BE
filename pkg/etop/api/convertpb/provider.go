package convertpb

import (
	"etop.vn/api/pb/etop/etc/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func ShippingProviderToModel(s *shipping_provider.ShippingProvider) model.ShippingProvider {
	if s == nil || *s == 0 {
		return ""
	}
	return model.ShippingProvider(shipping_provider.ShippingProvider_name[int(*s)])
}

func PbShippingProviderType(sp model.ShippingProvider) shipping_provider.ShippingProvider {
	spString := string(sp)
	return shipping_provider.ShippingProvider(shipping_provider.ShippingProvider_value[spString])
}

func PbPtrShippingProvider(sp model.ShippingProvider) *shipping_provider.ShippingProvider {
	res := PbShippingProviderType(sp)
	return &res
}

func PbShippingProviderPtr(s dot.NullString) *shipping_provider.ShippingProvider {
	if s.Apply("") == "" {
		return nil
	}
	sp := PbShippingProviderType(model.ShippingProvider(s.String))
	return &sp
}
