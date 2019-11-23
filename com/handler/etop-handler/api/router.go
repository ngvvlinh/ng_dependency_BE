package api

import (
	service "etop.vn/api/root/services/handler"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/root/services/handler
// +gen:wrapper:package=handler

func NewHandlerServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService, secret)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
