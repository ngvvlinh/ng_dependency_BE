// +build wireinject

package shop

import (
	"github.com/google/wire"
	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/main/address"
	"o.o/api/main/catalog"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/inventory"
	"o.o/api/main/ledgering"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/purchaserefund"
	"o.o/api/main/receipting"
	"o.o/api/main/refund"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipping"
	st "o.o/api/main/stocktaking"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/carrying"
	"o.o/api/shopping/customering"
	"o.o/api/shopping/suppliering"
	"o.o/api/shopping/tradering"
	"o.o/api/subscripting/subscription"
	"o.o/api/summary"
	"o.o/api/webserver"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/capi"
)

var WireSet = wire.NewSet(
	wire.Struct(new(AccountService)),
	wire.Struct(new(AuthorizeService)),
	wire.Struct(new(BrandService)),
	wire.Struct(new(CarrierService)),
	wire.Struct(new(CategoryService)),
	wire.Struct(new(CollectionService)),
	wire.Struct(new(ConnectionService)),
	wire.Struct(new(CustomerGroupService)),
	wire.Struct(new(CustomerService)),
	wire.Struct(new(ExportService)),
	wire.Struct(new(FulfillmentService)),
	wire.Struct(new(HistoryService)),
	wire.Struct(new(InventoryService)),
	wire.Struct(new(LedgerService)),
	wire.Struct(new(MiscService)),
	wire.Struct(new(MoneyTransactionService)),
	wire.Struct(new(NotificationService)),
	wire.Struct(new(OrderService)),
	wire.Struct(new(PaymentService)),
	wire.Struct(new(ProductService)),
	wire.Struct(new(ProductSourceService)),
	wire.Struct(new(PurchaseOrderService)),
	wire.Struct(new(PurchaseRefundService)),
	wire.Struct(new(ReceiptService)),
	wire.Struct(new(RefundService)),
	wire.Struct(new(ShipmentService)),
	wire.Struct(new(ShipnowService)),
	wire.Struct(new(StocktakeService)),
	wire.Struct(new(SubscriptionService)),
	wire.Struct(new(SummaryService)),
	wire.Struct(new(SupplierService)),
	wire.Struct(new(TradingService)),
	wire.Struct(new(WebServerService)),
	NewServers,
)

func BuildServers(
	locationQ location.QueryBus,
	catalogQueryBus catalog.QueryBus,
	catalogCommandBus catalog.CommandBus,
	shipnow shipnow.CommandBus,
	shipnowQS shipnow.QueryBus,
	identity identity.CommandBus,
	identityQS identity.QueryBus,
	addressQS address.QueryBus,
	providerManager *shipping_provider.ProviderManager,
	customerA customering.CommandBus,
	customerQS customering.QueryBus,
	traderAddressA addressing.CommandBus,
	traderAddressQ addressing.QueryBus,
	orderA ordering.CommandBus,
	orderQ ordering.QueryBus,
	paymentManager paymentmanager.CommandBus,
	supplierA suppliering.CommandBus,
	supplierQ suppliering.QueryBus,
	carrierA carrying.CommandBus,
	carrierQ carrying.QueryBus,
	traderQ tradering.QueryBus,
	eventB capi.EventBus,
	receiptA receipting.CommandBus,
	receiptQS receipting.QueryBus,
	sd cmservice.Shutdowner,
	rd redis.Store,
	inventoryA inventory.CommandBus,
	inventoryQ inventory.QueryBus,
	ledgerA ledgering.CommandBus,
	ledgerQ ledgering.QueryBus,
	purchaseOrderA purchaseorder.CommandBus,
	purchaseOrderQ purchaseorder.QueryBus,
	summary summary.QueryBus,
	StocktakeQ st.QueryBus,
	StocktakeA st.CommandBus,
	shipmentM *shippingcarrier.ShipmentManager,
	shippingA shipping.CommandBus,
	refundA refund.CommandBus,
	refundQ refund.QueryBus,
	purchaseRefundA purchaserefund.CommandBus,
	purchaseRefundQ purchaserefund.QueryBus,
	connectionQ connectioning.QueryBus,
	connectionA connectioning.CommandBus,
	shippingQ shipping.QueryBus,
	webserverA webserver.CommandBus,
	webserverQ webserver.QueryBus,
	subscriptionQ subscription.QueryBus,
) Servers {
	panic(wire.Build(WireSet))
}
