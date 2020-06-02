// +build wireinject

package shop

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	wire.Struct(new(AccountService), "*"),
	wire.Struct(new(AuthorizeService), "*"),
	wire.Struct(new(BrandService), "*"),
	wire.Struct(new(CarrierService), "*"),
	wire.Struct(new(CategoryService), "*"),
	wire.Struct(new(CollectionService), "*"),
	wire.Struct(new(ConnectionService), "*"),
	wire.Struct(new(CustomerGroupService), "*"),
	wire.Struct(new(CustomerService), "*"),
	wire.Struct(new(ExportService), "*"),
	wire.Struct(new(FulfillmentService), "*"),
	wire.Struct(new(HistoryService), "*"),
	wire.Struct(new(InventoryService), "*"),
	wire.Struct(new(LedgerService), "*"),
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(MoneyTransactionService), "*"),
	wire.Struct(new(NotificationService), "*"),
	wire.Struct(new(OrderService), "*"),
	wire.Struct(new(PaymentService), "*"),
	wire.Struct(new(ProductService), "*"),
	wire.Struct(new(ProductSourceService), "*"),
	wire.Struct(new(PurchaseOrderService), "*"),
	wire.Struct(new(PurchaseRefundService), "*"),
	wire.Struct(new(ReceiptService), "*"),
	wire.Struct(new(RefundService), "*"),
	wire.Struct(new(ShipmentService), "*"),
	wire.Struct(new(ShipnowService), "*"),
	wire.Struct(new(StocktakeService), "*"),
	wire.Struct(new(SubscriptionService), "*"),
	wire.Struct(new(SummaryService), "*"),
	wire.Struct(new(SupplierService), "*"),
	wire.Struct(new(TradingService), "*"),
	wire.Struct(new(WebServerService), "*"),
)
