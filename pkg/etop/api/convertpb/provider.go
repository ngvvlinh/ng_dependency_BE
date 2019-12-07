package convertpb

import (
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/capi/dot"
)

func PbShippingProviderType(sp shipping_provider.ShippingProvider) shipping_provider.ShippingProvider {
	return sp
}

func PbPtrShippingProvider(sp shipping_provider.ShippingProvider) shipping_provider.NullShippingProvider {
	return sp.Wrap()
}

func PbShippingProviderPtr(s dot.NullString) shipping_provider.NullShippingProvider {
	if s.Apply("") == "" {
		return shipping_provider.NullShippingProvider{}
	}
	sp, _ := shipping_provider.ParseShippingProvider(s.String)
	return sp.Wrap()
}
