package api

import (
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(
	miscService *MiscService,
	webhookService *WebhookService,
) Servers {
	servers := httprpc.MustNewServers(
		miscService.Clone,
		webhookService.Clone,
	)
	return servers
}
