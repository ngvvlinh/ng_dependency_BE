package api

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/services/handler"
)

// +gen:wrapper=etop.vn/backend/zexp/api/root/services/handler
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
