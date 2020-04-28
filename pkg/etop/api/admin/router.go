package admin

import (
	service "o.o/api/top/int/admin"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/admin
// +gen:wrapper:package=admin

func NewAdminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewAccountServiceServer(WrapAccountService(accountService)),
		service.NewOrderServiceServer(WrapOrderService(orderService)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService)),
		service.NewMoneyTransactionServiceServer(WrapMoneyTransactionService(moneyTransactionService)),
		service.NewShopServiceServer(WrapShopService(shopService)),
		service.NewCreditServiceServer(WrapCreditService(creditService)),
		service.NewNotificationServiceServer(WrapNotificationService(notificationService)),
		service.NewConnectionServiceServer(WrapConnectionService(connectionService)),
		service.NewShipmentPriceServiceServer(WrapShipmentPriceService(shipmentPriceService)),
		service.NewLocationServiceServer(WrapLocationService(locationService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
