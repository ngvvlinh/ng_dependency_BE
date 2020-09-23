package vnp

import (
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(
	shipnowService *VNPostService,
	webhookService *VNPostWebhookService,
) Servers {
	servers := httprpc.MustNewServers(
		shipnowService.Clone,
		webhookService.Clone,
	)
	return servers
}
