package _fabo

import (
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/capi/httprpc"
)

func NewServers(
	webhookService *sadmin.WebhookService,
) sadmin.Servers {
	servers := httprpc.MustNewServers(
		webhookService.Clone,
	)
	return servers
}
