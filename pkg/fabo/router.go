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
) FaboServers {
	servers := httprpc.MustNewServers(
		hooks,
		pageService.Clone,
		conversationService.Clone,
	)
	return servers
}
