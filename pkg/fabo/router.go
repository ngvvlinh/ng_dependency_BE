package fabo

import (
	"o.o/capi/httprpc"
)

type FaboServer struct {
	pageService                 *PageService
	customerConversationService *CustomerConversationService
}

type Servers []httprpc.Server

func NewServer(
	pageService *PageService,
	conversationService *CustomerConversationService,
	customerService *CustomerService,
) Servers {
	servers := httprpc.MustNewServers(
		pageService.Clone,
		conversationService.Clone,
		customerService.Clone,
	)
	return servers
}
