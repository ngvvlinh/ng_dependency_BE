package partner

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/partner"
)

// +gen:wrapper=etop.vn/backend/pb/external/partner
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

func NewPartnerServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewShopServiceServer(NewShopService(shopService)),
		service.NewWebhookServiceServer(NewWebhookService(webhookService)),
		service.NewHistoryServiceServer(NewHistoryService(historyService)),
		service.NewShippingServiceServer(NewShippingService(shippingService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
