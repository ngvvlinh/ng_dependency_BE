package fabo

import (
	"o.o/capi/httprpc"
)

type FaboServer struct {
	pageService                 *PageService
	customerConversationService *CustomerConversationService
}

type FaboServers []httprpc.Server

func NewServer(
	hooks httprpc.HooksBuilder,
	pageService *PageService,
	conversationService *CustomerConversationService,
	customerService *CustomerService,
) FaboServers {
	servers := httprpc.MustNewServers(
		hooks,
		pageService.Clone,
		conversationService.Clone,
		customerService.Clone,
	)
	return servers
}
