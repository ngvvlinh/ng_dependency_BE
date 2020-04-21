package admin

import (
	service "o.o/api/top/int/admin"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/admin
// +gen:wrapper:package=admin

func NewAdminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewAccountServiceServer(WrapAccountService(accountService.Clone)),
		service.NewOrderServiceServer(WrapOrderService(orderService.Clone)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewMoneyTransactionServiceServer(WrapMoneyTransactionService(moneyTransactionService.Clone)),
		service.NewShopServiceServer(WrapShopService(shopService.Clone)),
		service.NewCreditServiceServer(WrapCreditService(creditService.Clone)),
		service.NewNotificationServiceServer(WrapNotificationService(notificationService.Clone)),
		service.NewConnectionServiceServer(WrapConnectionService(connectionService.Clone)),
		service.NewShipmentPriceServiceServer(WrapShipmentPriceService(shipmentPriceService.Clone)),
		service.NewLocationServiceServer(WrapLocationService(locationService.Clone)),
		service.NewSubscriptionServiceServer(WrapSubscriptionService(subscriptionService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
