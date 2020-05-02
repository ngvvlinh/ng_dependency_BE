package api

import (
	service "o.o/api/top/services/affiliate"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/affiliate
// +gen:wrapper:package=affiliate

func NewAffiliateServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewUserServiceServer(WrapUserService(userService.Clone)),
		service.NewTradingServiceServer(WrapTradingService(tradingService.Clone)),
		service.NewShopServiceServer(WrapShopService(shopService.Clone)),
		service.NewAffiliateServiceServer(WrapAffiliateService(affiliateService.Clone, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
