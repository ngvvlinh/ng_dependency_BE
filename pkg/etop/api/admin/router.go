package admin

import (
	service "etop.vn/api/top/int/admin"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/int/admin
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
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
