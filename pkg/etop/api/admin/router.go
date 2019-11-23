package admin

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/admin"
)

// +gen:wrapper=etop.vn/backend/zexp/api/root/int/admin
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
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
