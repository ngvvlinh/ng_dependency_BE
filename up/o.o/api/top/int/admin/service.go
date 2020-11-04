package admin

import (
	"context"

	etop "o.o/api/top/int/etop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/int/types"
	cm "o.o/api/top/types/common"
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
	CreateAdminUser(context.Context, *CreateAdminUserRequest) (*CreateAdminUserResponse, error)
	UpdateAdminUser(context.Context, *UpdateAdminUserRequest) (*UpdateAdminUserResponse, error)
	GetAdminUsers(context.Context, *GetAdminUsersRequest) (*GetAdminUserResponse, error)
	DeleteAdminUser(context.Context, *DeleteAdminUserRequest) (*DeleteAdminUserResponse, error)
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
	UpdateFulfillmentInfo(context.Context, *UpdateFulfillmentInfoRequest) (*cm.UpdatedResponse, error)
	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateRequest) (*cm.UpdatedResponse, error)
	UpdateFulfillmentShippingFees(context.Context, *UpdateFulfillmentShippingFeesRequest) (*cm.UpdatedResponse, error)
	UpdateFulfillmentCODAmount(context.Context, *UpdateFulfillmentCODAmountRequest) (*cm.UpdatedResponse, error)
	AddShippingFee(context.Context, *AddShippingFeeRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.MoneyTransaction
type MoneyTransactionService interface {
	CreateMoneyTransaction(context.Context, *CreateMoneyTransactionRequest) (*types.MoneyTransaction, error)
	GetMoneyTransaction(context.Context, *cm.IDRequest) (*types.MoneyTransaction, error)
	GetMoneyTransactions(context.Context, *GetMoneyTransactionsRequest) (*types.MoneyTransactionsResponse, error)
	ConfirmMoneyTransaction(context.Context, *ConfirmMoneyTransactionRequest) (*cm.UpdatedResponse, error)
	UpdateMoneyTransaction(context.Context, *UpdateMoneyTransactionRequest) (*types.MoneyTransaction, error)

	GetMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*types.MoneyTransactionShippingExternal, error)
	GetMoneyTransactionShippingExternals(context.Context, *GetMoneyTransactionShippingExternalsRequest) (*types.MoneyTransactionShippingExternalsResponse, error)
	RemoveMoneyTransactionShippingExternalLines(context.Context, *RemoveMoneyTransactionShippingExternalLinesRequest) (*types.MoneyTransactionShippingExternal, error)
	DeleteMoneyTransactionShippingExternal(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
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
	GetShopsByIDs(context.Context, *cm.IDsRequest) (*GetShopsResponse, error)
	UpdateShopInfo(context.Context, *UpdateShopInfoRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.User
type UserService interface {
	GetUser(context.Context, *cm.IDRequest) (*etop.User, error)
	GetUsersByIDs(context.Context, *cm.IDsRequest) (*UserResponse, error)
	GetUsers(context.Context, *GetUsersRequest) (*UserResponse, error)

	BlockUser(context.Context, *BlockUserRequest) (*etop.User, error)
	UnblockUser(context.Context, *UnblockUserRequest) (*etop.User, error)

	UpdateUserRef(context.Context, *UpdateUserRefRequest) (*cm.Empty, error)
}

// +apix:path=/admin.Credit
type CreditService interface {
	GetCredit(context.Context, *GetCreditRequest) (*etop.Credit, error)
	GetCredits(context.Context, *GetCreditsRequest) (*etop.CreditsResponse, error)
	CreateCredit(context.Context, *CreateCreditRequest) (*etop.Credit, error)
	ConfirmCredit(context.Context, *ConfirmCreditRequest) (*cm.UpdatedResponse, error)
	DeleteCredit(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/admin.Notification
type NotificationService interface {
	CreateNotifications(context.Context, *CreateNotificationsRequest) (*CreateNotificationsResponse, error)
}

// +apix:path=/admin.Connection
type ConnectionService interface {
	GetConnections(context.Context, *types.GetConnectionsRequest) (*types.GetConnectionsResponse, error)
	ConfirmConnection(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	DisableConnection(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	GetConnectionServices(context.Context, *cm.IDRequest) (*types.GetConnectionServicesResponse, error)

	CreateBuiltinConnection(context.Context, *types.CreateBuiltinConnectionRequest) (*types.Connection, error)
	GetBuiltinShopConnections(context.Context, *cm.Empty) (*types.GetShopConnectionsResponse, error)
	UpdateBuiltinShopConnection(context.Context, *types.UpdateShopConnectionRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/admin.ShipmentPrice
type ShipmentPriceService interface {
	GetShippingServices(context.Context, *GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)

	GetShipmentServices(context.Context, *GetShipmentServicesRequest) (*GetShipmentServicesResponse, error)
	GetShipmentService(context.Context, *cm.IDRequest) (*ShipmentService, error)
	CreateShipmentService(context.Context, *CreateShipmentServiceRequest) (*ShipmentService, error)
	UpdateShipmentService(context.Context, *UpdateShipmentServiceRequest) (*cm.UpdatedResponse, error)
	DeleteShipmentService(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	UpdateShipmentServicesAvailableLocations(context.Context, *UpdateShipmentServicesAvailableLocationsRequest) (*cm.UpdatedResponse, error)
	UpdateShipmentServicesBlacklistLocations(context.Context, *UpdateShipmentServicesBlacklistLocationsRequest) (*cm.UpdatedResponse, error)

	GetShipmentPriceLists(context.Context, *GetShipmentPriceListsRequest) (*GetShipmentPriceListsResponse, error)
	GetShipmentPriceList(context.Context, *cm.IDRequest) (*ShipmentPriceList, error)
	CreateShipmentPriceList(context.Context, *CreateShipmentPriceListRequest) (*ShipmentPriceList, error)
	UpdateShipmentPriceList(context.Context, *UpdateShipmentPriceListRequest) (*cm.UpdatedResponse, error)
	SetDefaultShipmentPriceList(context.Context, *ActiveShipmentPriceListRequest) (*cm.UpdatedResponse, error)
	DeleteShipmentPriceList(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)

	GetShipmentPrices(context.Context, *GetShipmentPricesRequest) (*GetShipmentPricesResponse, error)
	GetShipmentPrice(context.Context, *cm.IDRequest) (*ShipmentPrice, error)
	CreateShipmentPrice(context.Context, *CreateShipmentPriceRequest) (*ShipmentPrice, error)
	UpdateShipmentPrice(context.Context, *UpdateShipmentPriceRequest) (*ShipmentPrice, error)
	DeleteShipmentPrice(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	UpdateShipmentPricesPriorityPoint(context.Context, *UpdateShipmentPricesPriorityPointRequest) (*cm.UpdatedResponse, error)

	GetShopShipmentPriceLists(context.Context, *GetShopShipmentPriceListsRequest) (*GetShopShipmentPriceListsResponse, error)
	GetShopShipmentPriceList(context.Context, *GetShopShipmentPriceListRequest) (*ShopShipmentPriceList, error)
	CreateShopShipmentPriceList(context.Context, *CreateShopShipmentPriceList) (*ShopShipmentPriceList, error)
	UpdateShopShipmentPriceList(context.Context, *UpdateShopShipmentPriceListRequest) (*cm.UpdatedResponse, error)
	DeleteShopShipmentPriceList(context.Context, *GetShopShipmentPriceListRequest) (*cm.DeletedResponse, error)

	GetShipmentPriceListPromotions(context.Context, *GetShipmentPriceListPromotionsRequest) (*GetShipmentPriceListPromotionsResponse, error)
	GetShipmentPriceListPromotion(context.Context, *cm.IDRequest) (*ShipmentPriceListPromotion, error)
	CreateShipmentPriceListPromotion(context.Context, *CreateShipmentPriceListPromotionRequest) (*ShipmentPriceListPromotion, error)
	UpdateShipmentPriceListPromotion(context.Context, *UpdateShipmentPriceListPromotionRequest) (*cm.UpdatedResponse, error)
	DeleteShipmentPriceListPromotion(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/admin.Location
type LocationService interface {
	GetCustomRegion(context.Context, *cm.IDRequest) (*CustomRegion, error)
	GetCustomRegions(context.Context, *cm.Empty) (*GetCustomRegionsResponse, error)
	CreateCustomRegion(context.Context, *CreateCustomRegionRequest) (*CustomRegion, error)
	UpdateCustomRegion(context.Context, *CustomRegion) (*cm.UpdatedResponse, error)
	DeleteCustomRegion(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/admin.Subscription
type SubscriptionService interface {
	CreateSubscriptionProduct(context.Context, *types.CreateSubrProductRequest) (*types.SubscriptionProduct, error)
	GetSubscriptionProducts(context.Context, *cm.Empty) (*types.GetSubrProductsResponse, error)
	DeleteSubscriptionProduct(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)

	CreateSubscriptionPlan(context.Context, *types.CreateSubrPlanRequest) (*types.SubscriptionPlan, error)
	UpdateSubscriptionPlan(context.Context, *types.UpdateSubrPlanRequest) (*cm.UpdatedResponse, error)
	GetSubscriptionPlans(context.Context, *cm.Empty) (*types.GetSubrPlansResponse, error)
	DeleteSubscriptionPlan(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)

	GetSubscription(context.Context, *types.SubscriptionIDRequest) (*types.Subscription, error)
	GetSubscriptions(context.Context, *types.GetSubscriptionsRequest) (*types.GetSubscriptionsResponse, error)
	CreateSubscription(context.Context, *types.CreateSubscriptionRequest) (*types.Subscription, error)
	UpdateSubscriptionInfo(context.Context, *types.UpdateSubscriptionInfoRequest) (*cm.UpdatedResponse, error)
	CancelSubscription(context.Context, *types.SubscriptionIDRequest) (*cm.UpdatedResponse, error)
	ActivateSubscription(context.Context, *types.SubscriptionIDRequest) (*cm.UpdatedResponse, error)
	DeleteSubscription(context.Context, *types.SubscriptionIDRequest) (*cm.DeletedResponse, error)

	GetSubscriptionBills(context.Context, *types.GetSubscriptionBillsRequest) (*types.GetSubscriptionBillsResponse, error)
	CreateSubscriptionBill(context.Context, *types.CreateSubscriptionBillRequest) (*types.SubscriptionBill, error)
	ManualPaymentSubscriptionBill(context.Context, *types.ManualPaymentSubscriptionBillRequest) (*cm.UpdatedResponse, error)
	DeleteSubscriptionBill(context.Context, *types.SubscriptionIDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/admin.Ticket
type TicketService interface {
	// ticket
	CreateTicket(context.Context, *CreateTicketRequest) (*shoptypes.Ticket, error)
	GetTickets(context.Context, *GetTicketsRequest) (*GetTicketsResponse, error)
	GetTicket(context.Context, *GetTicketRequest) (*shoptypes.Ticket, error)
	GetTicketsByRefTicketID(context.Context, *shoptypes.GetTicketsByRefTicketIDRequest) (*shoptypes.GetTicketsByRefTicketIDResponse, error)
	AssignTicket(context.Context, *AssignTicketRequest) (*shoptypes.Ticket, error)
	UnassignTicket(context.Context, *AssignTicketRequest) (*shoptypes.Ticket, error)
	ConfirmTicket(context.Context, *ConfirmTicketRequest) (*shoptypes.Ticket, error)
	CloseTicket(context.Context, *CloseTicketRequest) (*shoptypes.Ticket, error)
	ReopenTicket(context.Context, *ReopenTicketRequest) (*shoptypes.Ticket, error)
	UpdateTicketRefTicketID(context.Context, *UpdateTicketRefTicketIDRequest) (*cm.UpdatedResponse, error)

	// ticket comment
	CreateTicketComment(context.Context, *CreateTicketCommentRequest) (*shoptypes.TicketComment, error)
	UpdateTicketComment(context.Context, *UpdateTicketCommentRequest) (*shoptypes.TicketComment, error)
	DeleteTicketComment(context.Context, *DeleteTicketCommentRequest) (*DeleteTicketCommentResponse, error)
	GetTicketComments(context.Context, *GetTicketCommentsRequest) (*GetTicketCommentsResponse, error)

	// ticket label
	CreateTicketLabel(context.Context, *CreateTicketLabelRequest) (*shoptypes.TicketLabel, error)
	UpdateTicketLabel(context.Context, *UpdateTicketLabelRequest) (*shoptypes.TicketLabel, error)
	DeleteTicketLabel(context.Context, *DeleteTicketLabelRequest) (*DeleteTicketLabelResponse, error)
}
