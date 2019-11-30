package admin

import (
	"context"

	cm "etop.vn/api/pb/common"
	"etop.vn/api/pb/etop"
	admin "etop.vn/api/pb/etop/admin"
	"etop.vn/api/pb/etop/order"
)

// +gen:apix
// +gen:swagger:doc-path=etop/admin

// +apix:path=/admin.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	AdminLoginAsAccount(context.Context, *admin.LoginAsAccountRequest) (*etop.LoginResponse, error)
}

// +apix:path=/admin.Account
type AccountService interface {
	CreatePartner(context.Context, *admin.CreatePartnerRequest) (*etop.Partner, error)
	GenerateAPIKey(context.Context, *admin.GenerateAPIKeyRequest) (*admin.GenerateAPIKeyResponse, error)
}

// +apix:path=/admin.Order
type OrderService interface {
	GetOrder(context.Context, *cm.IDRequest) (*order.Order, error)
	GetOrders(context.Context, *admin.GetOrdersRequest) (*order.OrdersResponse, error)
	GetOrdersByIDs(context.Context, *cm.IDsRequest) (*order.OrdersResponse, error)
}

// +apix:path=/admin.Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *cm.IDRequest) (*order.Fulfillment, error)
	GetFulfillments(context.Context, *admin.GetFulfillmentsRequest) (*order.FulfillmentsResponse, error)

	// UpdateFulfillment
	//
	// `shipping_state`
	//
	// Only update from any state to `undeliverable`
	// Or update from `undeliverable`to any state
	UpdateFulfillment(context.Context, *admin.UpdateFulfillmentRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.MoneyTransaction
type MoneyTransactionService interface {
	GetMoneyTransaction(context.Context, *cm.IDRequest) (*order.MoneyTransaction, error)
	GetMoneyTransactions(context.Context, *admin.GetMoneyTransactionsRequest) (*order.MoneyTransactionsResponse, error)
	ConfirmMoneyTransaction(context.Context, *admin.ConfirmMoneyTransactionRequest) (*cm.UpdatedResponse, error)
	UpdateMoneyTransaction(context.Context, *admin.UpdateMoneyTransactionRequest) (*order.MoneyTransaction, error)

	GetMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*order.MoneyTransactionShippingExternal, error)
	GetMoneyTransactionShippingExternals(context.Context, *admin.GetMoneyTransactionShippingExternalsRequest) (*order.MoneyTransactionShippingExternalsResponse, error)
	RemoveMoneyTransactionShippingExternalLines(context.Context, *admin.RemoveMoneyTransactionShippingExternalLinesRequest) (*order.MoneyTransactionShippingExternal, error)
	DeleteMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
	ConfirmMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	ConfirmMoneyTransactionShippingExternals(context.Context, *cm.IDsRequest) (*cm.UpdatedResponse, error)
	UpdateMoneyTransactionShippingExternal(context.Context, *admin.UpdateMoneyTransactionShippingExternalRequest) (*order.MoneyTransactionShippingExternal, error)

	GetMoneyTransactionShippingEtop(context.Context, *cm.IDRequest) (*order.MoneyTransactionShippingEtop, error)
	GetMoneyTransactionShippingEtops(context.Context, *admin.GetMoneyTransactionShippingEtopsRequest) (*order.MoneyTransactionShippingEtopsResponse, error)
	CreateMoneyTransactionShippingEtop(context.Context, *cm.IDsRequest) (*order.MoneyTransactionShippingEtop, error)
	UpdateMoneyTransactionShippingEtop(context.Context, *admin.UpdateMoneyTransactionShippingEtopRequest) (*order.MoneyTransactionShippingEtop, error)
	DeleteMoneyTransactionShippingEtop(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	ConfirmMoneyTransactionShippingEtop(context.Context, *admin.ConfirmMoneyTransactionShippingEtopRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.Shop
type ShopService interface {
	GetShop(context.Context, *cm.IDRequest) (*etop.Shop, error)
	GetShops(context.Context, *admin.GetShopsRequest) (*admin.GetShopsResponse, error)
}

// +apix:path=/admin.Credit
type CreditService interface {
	GetCredit(context.Context, *admin.GetCreditRequest) (*etop.Credit, error)
	GetCredits(context.Context, *admin.GetCreditsRequest) (*etop.CreditsResponse, error)
	CreateCredit(context.Context, *admin.CreateCreditRequest) (*etop.Credit, error)
	UpdateCredit(context.Context, *admin.UpdateCreditRequest) (*etop.Credit, error)
	ConfirmCredit(context.Context, *admin.ConfirmCreditRequest) (*cm.UpdatedResponse, error)
	DeleteCredit(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/admin.Notification
type NotificationService interface {
	CreateNotifications(context.Context, *admin.CreateNotificationsRequest) (*admin.CreateNotificationsResponse, error)
}
