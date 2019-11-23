package xshop

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/external/shop"
)

// +gen:wrapper=etop.vn/backend/zexp/api/root/external/shop
// +gen:wrapper:package=shop
// +gen:wrapper:prefix=ext

func NewShopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService)),
		service.NewShippingServiceServer(WrapShippingService(shippingService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
