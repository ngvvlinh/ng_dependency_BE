package admin_min

import (
	"github.com/google/wire"

	service "o.o/api/top/int/admin"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	NewServers,
)

func NewServers(
	miscService admin.MiscService,
	accountService admin.AccountService,
	orderService admin.OrderService,
	fulfillmentService admin.FulfillmentService,
	moneyTransactionService admin.MoneyTransactionService,
	shopService admin.ShopService,
	creditService admin.CreditService,
	notificationService admin.NotificationService,
	connectionService admin.ConnectionService,
	shipmentPriceService admin.ShipmentPriceService,
	locationService admin.LocationService,
) admin.Servers {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(admin.WrapMiscService(miscService.Clone)),
		service.NewAccountServiceServer(admin.WrapAccountService(accountService.Clone)),
		service.NewOrderServiceServer(admin.WrapOrderService(orderService.Clone)),
		service.NewFulfillmentServiceServer(admin.WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewMoneyTransactionServiceServer(admin.WrapMoneyTransactionService(moneyTransactionService.Clone)),
		service.NewShopServiceServer(admin.WrapShopService(shopService.Clone)),
		service.NewCreditServiceServer(admin.WrapCreditService(creditService.Clone)),
		service.NewNotificationServiceServer(admin.WrapNotificationService(notificationService.Clone)),
		service.NewConnectionServiceServer(admin.WrapConnectionService(connectionService.Clone)),
		service.NewShipmentPriceServiceServer(admin.WrapShipmentPriceService(shipmentPriceService.Clone)),
		service.NewLocationServiceServer(admin.WrapLocationService(locationService.Clone)),
	}
	return servers
}
