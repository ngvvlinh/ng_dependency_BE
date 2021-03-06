package admin

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(AccountService), "*"),
	wire.Struct(new(OrderService), "*"),
	wire.Struct(new(FulfillmentService), "*"),
	wire.Struct(new(MoneyTransactionService), "*"),
	wire.Struct(new(ShopService), "*"),
	wire.Struct(new(CreditService), "*"),
	wire.Struct(new(NotificationService), "*"),
	wire.Struct(new(ConnectionService), "*"),
	wire.Struct(new(ShipmentPriceService), "*"),
	wire.Struct(new(LocationService), "*"),
	wire.Struct(new(SubscriptionService), "*"),
	wire.Struct(new(UserService), "*"),
	wire.Struct(new(TicketService), "*"),
	wire.Struct(new(EtelecomService), "*"),
	wire.Struct(new(InvoiceService), "*"),
	wire.Struct(new(TransactionService), "*"),
	wire.Struct(new(AccountUserService), "*"),
)
