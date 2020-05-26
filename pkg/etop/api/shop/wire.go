// +build wireinject

package shop

import (
	"github.com/google/wire"

	paymentmanager "o.o/backend/com/external/payment/manager"
	"o.o/backend/com/main/address"
	"o.o/backend/com/main/catalog"
	"o.o/backend/com/main/connectioning"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/main/inventory"
	"o.o/backend/com/main/ledgering"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/ordering"
	"o.o/backend/com/main/purchaseorder"
	"o.o/backend/com/main/purchaserefund"
	"o.o/backend/com/main/receipting"
	"o.o/backend/com/main/refund"
	"o.o/backend/com/main/shipmentpricing"
	"o.o/backend/com/main/shipnow"
	"o.o/backend/com/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/com/main/stocktaking"
	"o.o/backend/com/shopping/carrying"
	"o.o/backend/com/shopping/customering"
	"o.o/backend/com/shopping/suppliering"
	"o.o/backend/com/shopping/tradering"
	"o.o/backend/com/subscripting"
	"o.o/backend/com/summary"
	"o.o/backend/com/web/webserver"
	"o.o/backend/pkg/etop/logic/shipping_provider"
)

var WireDepsSet = wire.NewSet(
	location.WireSet,
	catalog.WireSet,
	shipnow.WireSet,
	identity.WireSet,
	address.WireSet,
	shipping_provider.WireSet,
	customering.WireSet,
	ordering.WireSet,
	paymentmanager.WireSet,
	suppliering.WireSet,
	carrying.WireSet,
	tradering.WireSet,
	receipting.WireSet,
	inventory.WireSet,
	ledgering.WireSet,
	purchaseorder.WireSet,
	shipmentpricing.WireSet,
	summary.WireSet,
	stocktaking.WireSet,
	shippingcarrier.WireSet,
	shipping.WireSet,
	refund.WireSet,
	purchaserefund.WireSet,
	connectioning.WireSet,
	webserver.WireSet,
	subscripting.WireSet,
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
	NewServers,
)
