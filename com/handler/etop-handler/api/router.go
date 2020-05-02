package api

import (
	service "o.o/api/top/services/handler"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/handler
// +gen:wrapper:package=handler

func NewHandlerServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone, secret)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService.Clone, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
