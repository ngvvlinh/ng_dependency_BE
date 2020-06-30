package admin_min

import (
	"github.com/google/wire"

	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	NewServers,
)

func NewServers(
	rd redis.Store,
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
	admin.InitIdemp(rd)
	servers := httprpc.MustNewServers(
		miscService.Clone,
		accountService.Clone,
		orderService.Clone,
		fulfillmentService.Clone,
		moneyTransactionService.Clone,
		shopService.Clone,
		creditService.Clone,
		notificationService.Clone,
		connectionService.Clone,
		shipmentPriceService.Clone,
		locationService.Clone,
	)
	return servers
}
