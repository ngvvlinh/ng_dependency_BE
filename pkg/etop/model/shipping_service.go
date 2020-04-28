package model

import "o.o/api/top/types/etc/shipping_provider"

type ShippingServiceRegistry struct {
	getShippingServiceNamesFuncs map[shipping_provider.ShippingProvider]func(code string) (name string, ok bool)

	// the map is only read after init, no need to lock
	init bool
}

var shippingServiceRegistry = &ShippingServiceRegistry{
	getShippingServiceNamesFuncs: make(map[shipping_provider.ShippingProvider]func(code string) (name string, ok bool)),
}

func GetShippingServiceRegistry() *ShippingServiceRegistry {
	return shippingServiceRegistry
}

func (s *ShippingServiceRegistry) Initialize() {
	if s.init {
		panic("already init")
	}
	if len(s.getShippingServiceNamesFuncs) != 4 { // ghn, ghtk, vtpost, etop
		panic("unexpected number of shipping service providers")
	}
	s.init = true
}

func (s *ShippingServiceRegistry) RegisterNameFunc(
	provider shipping_provider.ShippingProvider,
	fn func(code string) (name string, ok bool),
) {
	if s.init {
		panic("already init")
	}
	if s.getShippingServiceNamesFuncs[provider] != nil {
		panic("duplicated")
	}
	s.getShippingServiceNamesFuncs[provider] = fn
}

func (s *ShippingServiceRegistry) GetName(provider shipping_provider.ShippingProvider, code string) (name string, ok bool) {
	fn := s.getShippingServiceNamesFuncs[provider]
	if fn == nil {
		return "", false
	}
	return fn(code)
}
