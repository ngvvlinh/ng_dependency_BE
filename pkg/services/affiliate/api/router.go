package api

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/services_affiliate"
)

// +gen:wrapper=etop.vn/backend/pb/services/affiliate
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
