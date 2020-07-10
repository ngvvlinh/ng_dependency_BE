package api

import (
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

type Servers []httprpc.Server

func NewServers(
	userService *UserService,
	tradingService *TradingService,
	shopService *ShopService,
	affiliateService *AffiliateService,
) Servers {
	servers := httprpc.MustNewServers(
		userService.Clone,
		tradingService.Clone,
		shopService.Clone,
		affiliateService.Clone,
	)
	return servers
}
