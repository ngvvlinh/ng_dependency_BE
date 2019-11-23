package api

import (
	service "etop.vn/api/root/services/affiliate"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/root/services/affiliate
// +gen:wrapper:package=affiliate

func NewAffiliateServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewUserServiceServer(WrapUserService(userService)),
		service.NewTradingServiceServer(WrapTradingService(tradingService)),
		service.NewShopServiceServer(WrapShopService(shopService)),
		service.NewAffiliateServiceServer(WrapAffiliateService(affiliateService, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
