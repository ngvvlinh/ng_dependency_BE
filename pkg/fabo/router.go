package fabo

import (
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
)

var idempgroup *idemp.RedisGroup

type FaboServer struct {
	pageService                 *PageService
	customerConversationService *CustomerConversationService
}

type Servers []httprpc.Server

func NewServers(
	pageService *PageService,
	demoService *DemoService,
	conversationService *CustomerConversationService,
	customerService *CustomerService,
	shopService *ShopService,
	extraShipmentService *ExtraShipmentService,
	summaryService *SummaryService,
	rd redis.Store,
) Servers {
	idempgroup = idemp.NewRedisGroup(rd, "idemp_fabo", 30)
	servers := httprpc.MustNewServers(
		pageService.Clone,
		demoService.Clone,
		conversationService.Clone,
		customerService.Clone,
		shopService.Clone,
		extraShipmentService.Clone,
		summaryService.Clone,
	)
	return servers
}
