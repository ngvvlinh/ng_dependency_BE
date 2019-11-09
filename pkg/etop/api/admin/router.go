package admin

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/admin"
)

// +gen:wrapper=etop.vn/backend/pb/etop/admin
// +gen:wrapper:package=admin

func NewAdminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewAccountServiceServer(NewAccountService(accountService)),
		service.NewOrderServiceServer(NewOrderService(orderService)),
		service.NewFulfillmentServiceServer(NewFulfillmentService(fulfillmentService)),
		service.NewMoneyTransactionServiceServer(NewMoneyTransactionService(moneyTransactionService)),
		service.NewShopServiceServer(NewShopService(shopService)),
		service.NewCreditServiceServer(NewCreditService(creditService)),
		service.NewNotificationServiceServer(NewNotificationService(notificationService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
