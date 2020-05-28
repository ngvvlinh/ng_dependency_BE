package admin

import (
	service "o.o/api/top/int/admin"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

// +gen:wrapper=o.o/api/top/int/admin
// +gen:wrapper:package=admin

var ll = l.New()

type Servers []httprpc.Server

func NewServers(
	miscService *MiscService,
	accountService *AccountService,
	orderService *OrderService,
	fulfillmentService *FulfillmentService,
	moneyTransactionService *MoneyTransactionService,
	shopService *ShopService,
	creditService *CreditService,
	notificationService *NotificationService,
	connectionService *ConnectionService,
	shipmentPriceService *ShipmentPriceService,
	locationService *LocationService,
	subscriptionService *SubscriptionService,
) Servers {
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
	return servers
}
