package admin

import (
	"context"

	etop "etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/admin

// +apix:path=/admin.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	AdminLoginAsAccount(context.Context, *LoginAsAccountRequest) (*etop.LoginResponse, error)
}

// +apix:path=/admin.Account
type AccountService interface {
	CreatePartner(context.Context, *CreatePartnerRequest) (*etop.Partner, error)
	GenerateAPIKey(context.Context, *GenerateAPIKeyRequest) (*GenerateAPIKeyResponse, error)
}

// +apix:path=/admin.Order
type OrderService interface {
	GetOrder(context.Context, *cm.IDRequest) (*types.Order, error)
	GetOrders(context.Context, *GetOrdersRequest) (*types.OrdersResponse, error)
	GetOrdersByIDs(context.Context, *cm.IDsRequest) (*types.OrdersResponse, error)
}

// +apix:path=/admin.Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *cm.IDRequest) (*types.Fulfillment, error)
	GetFulfillments(context.Context, *GetFulfillmentsRequest) (*types.FulfillmentsResponse, error)

	// UpdateFulfillment
	//
	// `shipping_state`
	//
	// Only update from any state to `undeliverable`
	// Or update from `undeliverable`to any state
	UpdateFulfillment(context.Context, *UpdateFulfillmentRequest) (*cm.UpdatedResponse, error)

	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateRequest) (*cm.UpdatedResponse, error)
	UpdateFulfillmentShippingFee(context.Context, *UpdateFulfillmentShippingFeeRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.MoneyTransaction
type MoneyTransactionService interface {
	GetMoneyTransaction(context.Context, *cm.IDRequest) (*types.MoneyTransaction, error)
	GetMoneyTransactions(context.Context, *GetMoneyTransactionsRequest) (*types.MoneyTransactionsResponse, error)
	ConfirmMoneyTransaction(context.Context, *ConfirmMoneyTransactionRequest) (*cm.UpdatedResponse, error)
	UpdateMoneyTransaction(context.Context, *UpdateMoneyTransactionRequest) (*types.MoneyTransaction, error)

	GetMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*types.MoneyTransactionShippingExternal, error)
	GetMoneyTransactionShippingExternals(context.Context, *GetMoneyTransactionShippingExternalsRequest) (*types.MoneyTransactionShippingExternalsResponse, error)
	RemoveMoneyTransactionShippingExternalLines(context.Context, *RemoveMoneyTransactionShippingExternalLinesRequest) (*types.MoneyTransactionShippingExternal, error)
	DeleteMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
	ConfirmMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	ConfirmMoneyTransactionShippingExternals(context.Context, *cm.IDsRequest) (*cm.UpdatedResponse, error)
	UpdateMoneyTransactionShippingExternal(context.Context, *UpdateMoneyTransactionShippingExternalRequest) (*types.MoneyTransactionShippingExternal, error)

	GetMoneyTransactionShippingEtop(context.Context, *cm.IDRequest) (*types.MoneyTransactionShippingEtop, error)
	GetMoneyTransactionShippingEtops(context.Context, *GetMoneyTransactionShippingEtopsRequest) (*types.MoneyTransactionShippingEtopsResponse, error)
	CreateMoneyTransactionShippingEtop(context.Context, *cm.IDsRequest) (*types.MoneyTransactionShippingEtop, error)
	UpdateMoneyTransactionShippingEtop(context.Context, *UpdateMoneyTransactionShippingEtopRequest) (*types.MoneyTransactionShippingEtop, error)
	DeleteMoneyTransactionShippingEtop(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	ConfirmMoneyTransactionShippingEtop(context.Context, *ConfirmMoneyTransactionShippingEtopRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.Shop
type ShopService interface {
	GetShop(context.Context, *cm.IDRequest) (*etop.Shop, error)
	GetShops(context.Context, *GetShopsRequest) (*GetShopsResponse, error)
}

// +apix:path=/admin.Credit
type CreditService interface {
	GetCredit(context.Context, *GetCreditRequest) (*etop.Credit, error)
	GetCredits(context.Context, *GetCreditsRequest) (*etop.CreditsResponse, error)
	CreateCredit(context.Context, *CreateCreditRequest) (*etop.Credit, error)
	UpdateCredit(context.Context, *UpdateCreditRequest) (*etop.Credit, error)
	ConfirmCredit(context.Context, *ConfirmCreditRequest) (*cm.UpdatedResponse, error)
	DeleteCredit(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/admin.Notification
type NotificationService interface {
	CreateNotifications(context.Context, *CreateNotificationsRequest) (*CreateNotificationsResponse, error)
}
