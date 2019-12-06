package convertpb

import (
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func ShippingProviderToModel(s *shipping_provider.ShippingProvider) model.ShippingProvider {
	if s == nil || *s == 0 {
		return ""
	}
	return model.ShippingProvider(s.String())
}

func PbShippingProviderType(sp model.ShippingProvider) shipping_provider.ShippingProvider {
	p, _ := shipping_provider.ParseShippingProvider(string(sp))
	return p
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
