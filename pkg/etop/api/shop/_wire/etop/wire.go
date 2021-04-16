// +build wireinject

package etop

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/account"
	"o.o/backend/pkg/etop/api/shop/accountshipnow"
	"o.o/backend/pkg/etop/api/shop/authorize"
	"o.o/backend/pkg/etop/api/shop/brand"
	"o.o/backend/pkg/etop/api/shop/carrier"
	"o.o/backend/pkg/etop/api/shop/category"
	"o.o/backend/pkg/etop/api/shop/collection"
	"o.o/backend/pkg/etop/api/shop/connection"
	"o.o/backend/pkg/etop/api/shop/contact"
	"o.o/backend/pkg/etop/api/shop/credit"
	"o.o/backend/pkg/etop/api/shop/customer"
	"o.o/backend/pkg/etop/api/shop/customergroup"
	"o.o/backend/pkg/etop/api/shop/etelecom"
	"o.o/backend/pkg/etop/api/shop/etelecom/etelecomuser"
	"o.o/backend/pkg/etop/api/shop/export"
	"o.o/backend/pkg/etop/api/shop/fulfillment"
	"o.o/backend/pkg/etop/api/shop/history"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/api/shop/invoice"
	"o.o/backend/pkg/etop/api/shop/ledger"
	"o.o/backend/pkg/etop/api/shop/money_transaction"
	"o.o/backend/pkg/etop/api/shop/notification"
	"o.o/backend/pkg/etop/api/shop/order"
	"o.o/backend/pkg/etop/api/shop/payment"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/product_source"
	"o.o/backend/pkg/etop/api/shop/purchase_order"
	"o.o/backend/pkg/etop/api/shop/purchase_refund"
	"o.o/backend/pkg/etop/api/shop/receipt"
	"o.o/backend/pkg/etop/api/shop/refund"
	"o.o/backend/pkg/etop/api/shop/setting"
	"o.o/backend/pkg/etop/api/shop/shipment"
	"o.o/backend/pkg/etop/api/shop/shipnow"
	"o.o/backend/pkg/etop/api/shop/stocktake"
	"o.o/backend/pkg/etop/api/shop/subscription"
	"o.o/backend/pkg/etop/api/shop/summary"
	"o.o/backend/pkg/etop/api/shop/supplier"
	"o.o/backend/pkg/etop/api/shop/ticket"
	"o.o/backend/pkg/etop/api/shop/trading"
	"o.o/backend/pkg/etop/api/shop/transaction"
	"o.o/backend/pkg/etop/api/shop/ws"
)

var WireSet = wire.NewSet(
	wire.Struct(new(account.AccountService), "*"),
	wire.Struct(new(authorize.AuthorizeService), "*"),
	wire.Struct(new(brand.BrandService), "*"),
	wire.Struct(new(carrier.CarrierService), "*"),
	wire.Struct(new(category.CategoryService), "*"),
	wire.Struct(new(collection.CollectionService), "*"),
	wire.Struct(new(connection.ConnectionService), "*"),
	wire.Struct(new(customergroup.CustomerGroupService), "*"),
	wire.Struct(new(customer.CustomerService), "*"),
	wire.Struct(new(credit.CreditService), "*"),
	wire.Struct(new(export.ExportService), "*"),
	wire.Struct(new(fulfillment.FulfillmentService), "*"),
	wire.Struct(new(history.HistoryService), "*"),
	wire.Struct(new(inventory.InventoryService), "*"),
	wire.Struct(new(ledger.LedgerService), "*"),
	wire.Struct(new(shop.MiscService), "*"),
	wire.Struct(new(money_transaction.MoneyTransactionService), "*"),
	wire.Struct(new(notification.NotificationService), "*"),
	wire.Struct(new(order.OrderService), "*"),
	wire.Struct(new(payment.PaymentService), "*"),
	wire.Struct(new(product.ProductService), "*"),
	wire.Struct(new(product_source.ProductSourceService), "*"),
	wire.Struct(new(purchase_order.PurchaseOrderService), "*"),
	wire.Struct(new(purchase_refund.PurchaseRefundService), "*"),
	wire.Struct(new(receipt.ReceiptService), "*"),
	wire.Struct(new(refund.RefundService), "*"),
	wire.Struct(new(shipment.ShipmentService), "*"),
	wire.Struct(new(shipnow.ShipnowService), "*"),
	wire.Struct(new(stocktake.StocktakeService), "*"),
	wire.Struct(new(subscription.SubscriptionService), "*"),
	wire.Struct(new(summary.SummaryService), "*"),
	wire.Struct(new(supplier.SupplierService), "*"),
	wire.Struct(new(trading.TradingService), "*"),
	wire.Struct(new(ws.WebServerService), "*"),
	wire.Struct(new(ticket.TicketService), "*"),
	wire.Struct(new(accountshipnow.AccountShipnowService), "*"),
	wire.Struct(new(contact.ContactService), "*"),
	wire.Struct(new(setting.SettingService), "*"),
	wire.Struct(new(etelecom.EtelecomService), "*"),
	wire.Struct(new(etelecomuser.EtelecomUserService), "*"),
	wire.Struct(new(transaction.TransactionService), "*"),
	wire.Struct(new(invoice.InvoiceService), "*"),
)
