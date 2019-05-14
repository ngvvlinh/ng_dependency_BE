package model

type ShippingServiceRegistry struct {
	getShippingServiceNamesFuncs map[ShippingProvider]func(code string) (name string, ok bool)

	// the map is only read after init, no need to lock
	init bool
}

var shippingServiceRegistry = &ShippingServiceRegistry{
	getShippingServiceNamesFuncs: make(map[ShippingProvider]func(code string) (name string, ok bool)),
}

func GetShippingServiceRegistry() *ShippingServiceRegistry {
	return shippingServiceRegistry
}

func (s *ShippingServiceRegistry) Initialize() {
	if s.init {
		panic("already init")
	}
	if len(s.getShippingServiceNamesFuncs) != 5 { // ghn, ghtk, vtpost, etop, ahamove
		panic("unexpected number of shipping service providers")
	}
	s.init = true
}

func (s *ShippingServiceRegistry) RegisterNameFunc(
	provider ShippingProvider,
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

func (s *ShippingServiceRegistry) GetName(provider ShippingProvider, code string) (name string, ok bool) {
	fn := s.getShippingServiceNamesFuncs[provider]
	if fn == nil {
		return "", false
	}
	return fn(code)
}
