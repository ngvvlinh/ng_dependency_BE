package api

import (
	service "o.o/api/top/services/affiliate"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

// +gen:wrapper=o.o/api/top/services/affiliate
// +gen:wrapper:package=affiliate

var ll = l.New()

type Secret string
type Servers []httprpc.Server

func NewServers(
	secret Secret,
	userService *UserService,
	tradingService *TradingService,
	shopService *ShopService,
	affiliateService *AffiliateService,
) Servers {
	servers := []httprpc.Server{
		service.NewUserServiceServer(WrapUserService(userService.Clone)),
		service.NewTradingServiceServer(WrapTradingService(tradingService.Clone)),
		service.NewShopServiceServer(WrapShopService(shopService.Clone)),
		service.NewAffiliateServiceServer(WrapAffiliateService(affiliateService.Clone, string(secret))),
	}
	return servers
}
