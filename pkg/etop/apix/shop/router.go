package xshop

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/external_shop"
)

// +gen:wrapper=etop.vn/backend/pb/external/shop
// +gen:wrapper:package=shop
// +gen:wrapper:prefix=ext

func NewShopServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewWebhookServiceServer(NewWebhookService(webhookService)),
		service.NewHistoryServiceServer(NewHistoryService(historyService)),
		service.NewShippingServiceServer(NewShippingService(shippingService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
