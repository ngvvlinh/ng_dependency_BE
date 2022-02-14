package admin

import (
	"time"

	etop "o.o/api/top/int/etop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/int/types"
	common "o.o/api/top/types/common"
	"o.o/api/top/types/etc/additional_fee_base_value"
	"o.o/api/top/types/etc/calculation_method"
	credit_type "o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/filter_type"
	"o.o/api/top/types/etc/location_type"
	notifier_entity "o.o/api/top/types/etc/notifier_entity"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/route_type"
	shipping "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shop_user_role"
	status3 "o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type GetTicketCommentsResponse struct {
	TicketComments []*shoptypes.TicketComment `json:"ticket_comments"`
	Paging         *common.PageInfo           `json:"paging"`
}

func (m *GetTicketCommentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketCommentsRequest struct {
	Filter *FilterGetTicketComment `json:"filter"`
	Paging *common.Paging          `json:"paging"`
}

func (m *GetTicketCommentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetTicketComment struct {
	IDs       []dot.ID `json:"ids"`
	Title     string   `json:"title"`
	CreatedBy dot.ID   `json:"created_by"`
	ParentID  dot.ID   `json:"parent_id"`
	AccountID dot.ID   `json:"account_id"`
	TicketID  dot.ID   `json:"ticket_id"`
}

func (m *FilterGetTicketComment) String() string { return jsonx.MustMarshalToString(m) }

type UpdateTicketCommentRequest struct {
	AccountID dot.ID `json:"account_id"`
	ID        dot.ID `json:"id"`
	Message   string `json:"message"`
	// @deprecated
	ImageUrl  string   `json:"image_url"`
	ImageUrls []string `json:"image_urls"`
}

func (m *UpdateTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateTicketCommentRequest struct {
	AccountID dot.ID `json:"account_id"`
	TicketID  dot.ID `json:"ticket_id"`
	Message   string `json:"message"`
	// @deprecated
	ImageUrl  string   `json:"image_url"`
	ImageUrls []string `json:"image_urls"`
	ParentID  dot.ID   `json:"parent_id"`
}

func (m *CreateTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketCommentRequest struct {
	ID dot.ID `json:"id"`
}

func (m *DeleteTicketCommentRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketCommentResponse struct {
	Count int `json:"count"`
}

func (m *DeleteTicketCommentResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
}

func (m *ConfirmTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type ReopenTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
}

func (m *ReopenTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type CloseTicketRequest struct {
	TicketID dot.ID `json:"ticket_id"`
	Note     string `json:"note"`
	// @required
	State ticket_state.TicketState `json:"state"`
}

func (m *CloseTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type UnassignTicketRequest struct {
	AssignedUserIDs []dot.ID `json:"assigned_user_ids"`
	TicketID        dot.ID   `json:"ticket_id"`
}

func (m *UnassignTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type AssignTicketRequest struct {
	AssignedUserIDs []dot.ID `json:"assigned_user_ids"`
	TicketID        dot.ID   `json:"ticket_id"`
}

func (m *AssignTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateTicketRequest struct {
	LabelIDs []dot.ID `json:"label_ids"`

	Title       string `json:"title"`
	Description string `json:"description"`

	// user note
	Note string `json:"note"`

	RefID     dot.ID                        `json:"ref_id"`
	RefType   ticket_ref_type.TicketRefType `json:"ref_type"`
	RefCode   string                        `json:"ref_code"`
	Source    ticket_source.TicketSource    `json:"source"`
	AccountID dot.ID                        `json:"account_id"`
	// Ticket ID liên quan
	RefTicketID dot.NullID `json:"ref_ticket_id"`
}

func (m *CreateTicketRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateTicketRefTicketIDRequest struct {
	ID dot.ID `json:"id"`
	// Truyền lên 0 để xóa ref_ticket_id
	RefTicketID dot.NullID `json:"ref_ticket_id"`
}

func (m *UpdateTicketRefTicketIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterShopGetTicket struct {
	IDs            []dot.ID                      `json:"ids"`
	CreatedBy      dot.ID                        `json:"created_by"`
	ClosedBy       dot.ID                        `json:"closed_by"`
	AccountID      dot.ID                        `json:"account_id"`
	LabelIDs       []dot.ID                      `json:"label_ids"`
	Title          filter.FullTextSearch         `json:"title"`
	AssignedUserID []dot.ID                      `json:"assigned_user_id"`
	Code           string                        `json:"code"`
	State          ticket_state.TicketState      `json:"state"`
	RefID          dot.ID                        `json:"ref_id"`
	RefType        ticket_ref_type.TicketRefType `json:"ref_type"`
	RefCode        string                        `json:"ref_code"`
}

func (m *FilterShopGetTicket) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketsResponse struct {
	Paging  *common.PageInfo    `json:"paging"`
	Tickets []*shoptypes.Ticket `json:"tickets"`
}

func (m *GetTicketsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketRequest struct {
	ID        dot.ID `json:"id"`
	AccountID dot.ID `json:"account_id"`
}

func (m *GetTicketRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetTicketsRequest struct {
	Paging *common.Paging       `json:"paging"`
	Filter *FilterShopGetTicket `json:"filter"`
}

func (m *GetTicketsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetTicketResponse struct {
	Tickets []*shoptypes.Ticket `json:"tickets"`
}

func (m *GetTicketResponse) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketLabelRequest struct {
	ID          dot.ID `json:"id"`
	DeleteChild bool   `json:"delete_child"`
}

func (m *DeleteTicketLabelRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetTicketLabelsRequest struct {
	Tree bool `json:"tree"`
}

func (m *GetTicketLabelsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateTicketLabelRequest struct {
	ParentID dot.ID `json:"parent_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Code     string `json:"code"`
}

func (m *CreateTicketLabelRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateTicketLabelRequest struct {
	ID       dot.ID         `json:"id"`
	Name     dot.NullString `json:"name"`
	Color    string         `json:"color"`
	Code     dot.NullString `json:"code"`
	ParentID dot.NullID     `json:"parent_id"`
}

func (m *UpdateTicketLabelRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteTicketLabelResponse struct {
	Count int `json:"deleted"`
}

func (m *DeleteTicketLabelResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketLabelsResponse struct {
	TicketLabels []*shoptypes.TicketLabel `json:"ticket_labels"`
}

func (m *GetTicketLabelsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketLabelIDRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetTicketLabelIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetOrdersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsRequest struct {
	Paging        *common.Paging     `json:"paging"`
	ShopId        dot.ID             `json:"shop_id"`
	OrderId       dot.ID             `json:"order_id"`
	Status        status3.NullStatus `json:"status"`
	ConnectionIDs []dot.ID           `json:"connection_ids"`
	Filters       []*common.Filter   `json:"filters"`
	DateFrom      time.Time          `json:"date_from"`
	DateTo        time.Time          `json:"date_to"`
}

func (m *GetFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginAsAccountRequest struct {
	UserId    dot.ID `json:"user_id"`
	AccountId dot.ID `json:"account_id"`
	// This is a sensitive API, so admin must provide password before processing!
	Password string `json:"password"`
}

func (m *LoginAsAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateMoneyTransactionRequest struct {
	ShopID         dot.ID   `json:"shop_id"`
	FulfillmentIDs []dot.ID `json:"fulfillment_ids"`
	TotalCOD       int      `json:"total_cod"`
	TotalAmount    int      `json:"total_amount"`
	TotalOrders    int      `json:"total_orders"`
}

func (m *CreateMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionsRequest struct {
	Ids     []dot.ID         `json:"ids"`
	ShopId  dot.ID           `json:"shop_id"`
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveFfmsMoneyTransactionRequest struct {
	FulfillmentIds     []dot.ID `json:"fulfillment_ids"`
	MoneyTransactionId dot.ID   `json:"money_transaction_id"`
	ShopId             dot.ID   `json:"shop_id"`
}

func (m *RemoveFfmsMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionShippingExternalsRequest struct {
	Ids     []dot.ID         `json:"ids"`
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionShippingExternalsRequest) Reset() {
	*m = GetMoneyTransactionShippingExternalsRequest{}
}
func (m *GetMoneyTransactionShippingExternalsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type RemoveMoneyTransactionShippingExternalLinesRequest struct {
	LineIds                            []dot.ID `json:"line_ids"`
	MoneyTransactionShippingExternalId dot.ID   `json:"money_transaction_shipping_external_id"`
}

func (m *RemoveMoneyTransactionShippingExternalLinesRequest) Reset() {
	*m = RemoveMoneyTransactionShippingExternalLinesRequest{}
}
func (m *RemoveMoneyTransactionShippingExternalLinesRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type ConfirmMoneyTransactionRequest struct {
	MoneyTransactionId dot.ID `json:"money_transaction_id"`
	ShopId             dot.ID `json:"shop_id"`
	TotalCod           int    `json:"total_cod"`
	TotalAmount        int    `json:"total_amount"`
	TotalOrders        int    `json:"total_orders"`
}

func (m *ConfirmMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteMoneyTransactionRequest struct {
	MoneyTransactionId dot.ID `json:"money_transaction_id"`
	ShopId             dot.ID `json:"shop_id"`
}

func (m *DeleteMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterShopRequest struct {
	Name     filter.FullTextSearch `json:"name"`
	ShopIDs  []dot.ID              `json:"shop_ids"`
	OwnerID  dot.ID                `json:"owner_id"`
	DateFrom time.Time             `json:"date_from"`
	DateTo   time.Time             `json:"date_to"`
}

func (m *FilterShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopsRequest struct {
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
	Filter  *FilterShopRequest `json:"filter"`
}

func (m *GetShopsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UsersFilter struct {
	Name      filter.FullTextSearch `json:"name"`
	Phone     string                `json:"phone"`
	Email     string                `json:"email"`
	CreatedAt filter.Date           `json:"created_at"`
	RefAff    string                `json:"ref_aff"`
	RefSale   string                `json:"ref_sale"`
}

func (m *UsersFilter) String() string { return jsonx.MustMarshalToString(m) }

type GetUsersRequest struct {
	Paging  *common.Paging `json:"paging"`
	Filters *UsersFilter   `json:"filters"`
}

func (m *GetUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopsResponse struct {
	Paging *common.PageInfo `json:"paging"`
	Shops  []*etop.Shop     `json:"shops"`
}

func (m *GetShopsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetCreditRequest struct {
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
}

func (m *GetCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCreditsRequest struct {
	// @deprecated: use Filter instead
	ShopId dot.ID         `json:"shop_id"`
	Paging *common.Paging `json:"paging"`
	Filter *CreditsFilter `json:"filter"`
}

func (m *GetCreditsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreditsFilter struct {
	ShopID   dot.ID                         `json:"shop_id"`
	Classify credit_type.NullCreditClassify `json:"classify"`
	DateFrom time.Time                      `json:"date_from"`
	DateTo   time.Time                      `json:"date_to"`
}

type CreateCreditRequest struct {
	Amount   int                            `json:"amount"`
	ShopId   dot.ID                         `json:"shop_id"`
	Type     credit_type.CreditType         `json:"type"`
	PaidAt   dot.Time                       `json:"paid_at"`
	Classify credit_type.NullCreditClassify `json:"classify"`
}

func (m *CreateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmCreditRequest struct {
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
}

func (m *ConfirmCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreatePartnerRequest struct {
	Partner etop.Partner `json:"partner"`
}

func (m *CreatePartnerRequest) String() string { return jsonx.MustMarshalToString(m) }

type BlockUserRequest struct {
	UserID      dot.ID `json:"user_id"`
	BlockReason string `json:"block_reason"`
}

func (m *BlockUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type UnblockUserRequest struct {
	UserID dot.ID `json:"user_id"`
}

func (m *UnblockUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserRefRequest struct {
	UserID  dot.ID         `json:"user_id"`
	RefAff  dot.NullString `json:"ref_aff"`
	RefSale dot.NullString `json:"ref_sale"`
}

func (m *UpdateUserRefRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentInfoRequest struct {
	ID           dot.ID         `json:"id"`
	ShippingCode string         `json:"shipping_code"`
	FullName     dot.NullString `json:"full_name"`
	Phone        dot.NullString `json:"phone"`
	AdminNote    string         `json:"admin_note"`
}

func (m *UpdateFulfillmentInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

type GenerateAPIKeyRequest struct {
	AccountId dot.ID `json:"account_id"`
}

func (m *GenerateAPIKeyRequest) String() string { return jsonx.MustMarshalToString(m) }

type GenerateAPIKeyResponse struct {
	AccountId dot.ID `json:"account_id"`
	ApiKey    string `json:"api_key"`
}

func (m *GenerateAPIKeyResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateMoneyTransactionShippingEtopRequest struct {
	Id            dot.ID            `json:"id"`
	Adds          []dot.ID          `json:"adds"`
	Deletes       []dot.ID          `json:"deletes"`
	ReplaceAll    []dot.ID          `json:"replace_all"`
	Note          string            `json:"note"`
	InvoiceNumber string            `json:"invoice_number"`
	BankAccount   *etop.BankAccount `json:"bank_account"`
}

func (m *UpdateMoneyTransactionShippingEtopRequest) Reset() {
	*m = UpdateMoneyTransactionShippingEtopRequest{}
}
func (m *UpdateMoneyTransactionShippingEtopRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetMoneyTransactionShippingEtopsRequest struct {
	Ids     []dot.ID           `json:"ids"`
	Status  status3.NullStatus `json:"status"`
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
}

func (m *GetMoneyTransactionShippingEtopsRequest) Reset() {
	*m = GetMoneyTransactionShippingEtopsRequest{}
}
func (m *GetMoneyTransactionShippingEtopsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type ConfirmMoneyTransactionShippingEtopRequest struct {
	Id          dot.ID `json:"id"`
	TotalCod    int    `json:"total_cod"`
	TotalAmount int    `json:"total_amount"`
	TotalOrders int    `json:"total_orders"`
}

func (m *ConfirmMoneyTransactionShippingEtopRequest) Reset() {
	*m = ConfirmMoneyTransactionShippingEtopRequest{}
}
func (m *ConfirmMoneyTransactionShippingEtopRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateMoneyTransactionRequest struct {
	Id            dot.ID            `json:"id"`
	Note          string            `json:"note"`
	InvoiceNumber string            `json:"invoice_number"`
	BankAccount   *etop.BankAccount `json:"bank_account"`
}

func (m *UpdateMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateMoneyTransactionShippingExternalRequest struct {
	Id            dot.ID            `json:"id"`
	Note          string            `json:"note"`
	InvoiceNumber string            `json:"invoice_number"`
	BankAccount   *etop.BankAccount `json:"bank_account"`
}

func (m *UpdateMoneyTransactionShippingExternalRequest) Reset() {
	*m = UpdateMoneyTransactionShippingExternalRequest{}
}
func (m *UpdateMoneyTransactionShippingExternalRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateNotificationsRequest struct {
	Title      string                         `json:"title"`
	Message    string                         `json:"message"`
	Entity     notifier_entity.NotifierEntity `json:"entity"`
	EntityId   dot.ID                         `json:"entity_id"`
	AccountIds []dot.ID                       `json:"account_ids"`
	// Send to all subscribers
	SendAll bool `json:"send_all"`
}

func (m *CreateNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateNotificationsResponse struct {
	Created int `json:"created"`
	Errored int `json:"errored"`
}

func (m *CreateNotificationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentShippingStateRequest struct {
	ID                       dot.ID         `json:"id"`
	ShippingCode             string         `json:"shipping_code"`
	ShippingState            shipping.State `json:"shipping_state"`
	ActualCompensationAmount dot.NullInt    `json:"actual_compensation_amount"`
	AdminNote                string         `json:"admin_note"`
}

func (m *UpdateFulfillmentShippingStateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentShippingFeesRequest struct {
	ID               dot.ID                   `json:"id"`
	ShippingCode     string                   `json:"shipping_code"`
	ShippingFeeLines []*types.ShippingFeeLine `json:"shipping_fee_lines"`
	// @deprecated TotalCODAmount
	TotalCODAmount dot.NullInt `json:"total_cod_amount"`
	AdminNote      string      `json:"admin_note"`
}

func (m *UpdateFulfillmentShippingFeesRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentCODAmountRequest struct {
	ID                dot.ID       `json:"id"`
	ShippingCode      string       `json:"shipping_code"`
	TotalCODAmount    dot.NullInt  `json:"total_cod_amount"`
	IsPartialDelivery dot.NullBool `json:"is_partial_delivery"`
	AdminNote         string       `json:"admin_note"`
}

func (m *UpdateFulfillmentCODAmountRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddShippingFeeRequest struct {
	ID              dot.ID                            `json:"id"`
	ShippingCode    string                            `json:"shipping_code"`
	ShippingFeeType shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
}

func (m *AddShippingFeeRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentPrice struct {
	ID                  dot.ID                             `json:"id"`
	ShipmentPriceListID dot.ID                             `json:"shipment_price_list_id"`
	ShipmentServiceID   dot.ID                             `json:"shipment_service_id"`
	Name                string                             `json:"name"`
	CustomRegionTypes   []route_type.CustomRegionRouteType `json:"custom_region_types"`
	CustomRegionIDs     []dot.ID                           `json:"custom_region_ids"`
	RegionTypes         []route_type.RegionRouteType       `json:"region_types"`
	ProvinceTypes       []route_type.ProvinceRouteType     `json:"province_types"`
	UrbanTypes          []route_type.UrbanType             `json:"urban_types"`
	PriorityPoint       int                                `json:"priority_point"`
	Details             []*PricingDetail                   `json:"details"`
	CreatedAt           time.Time                          `json:"created_at"`
	UpdatedAt           time.Time                          `json:"updated_at"`
	Status              status3.Status                     `json:"status"`
	AdditionalFees      []*AdditionalFee                   `json:"additional_fees"`
}

func (m *ShipmentPrice) String() string { return jsonx.MustMarshalToString(m) }

type PricingDetail struct {
	Weight     int                        `json:"weight"`
	Price      int                        `json:"price"`
	Overweight []*PricingDetailOverweight `json:"overweight"`
}

func (m *PricingDetail) String() string { return jsonx.MustMarshalToString(m) }

type PricingDetailOverweight struct {
	MinWeight  int `json:"min_weight"`
	MaxWeight  int `json:"max_weight"`
	WeightStep int `json:"weight_step"`
	PriceStep  int `json:"price_step"`
}

func (m *PricingDetailOverweight) String() string { return jsonx.MustMarshalToString(m) }

type AdditionalFee struct {
	FeeType           shipping_fee_type.ShippingFeeType        `json:"fee_type"`
	CalculationMethod calculation_method.CalculationMethodType `json:"calculation_method"`
	BaseValueType     additional_fee_base_value.BaseValueType  `json:"base_value_type"`
	Rules             []*AdditionalFeeRule                     `json:"rules"`
}

func (m *AdditionalFee) String() string { return jsonx.MustMarshalToString(m) }

type AdditionalFeeRule struct {
	MinValue          int                                   `json:"min_value"`
	MaxValue          int                                   `json:"max_value"`
	PriceModifierType price_modifier_type.PriceModifierType `json:"price_modifier_type"`
	Amount            float64                               `json:"amount"`
	MinPrice          int                                   `json:"min_price"`
	StartValue        int                                   `json:"start_value"`
}

func (m *AdditionalFeeRule) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPricesResponse struct {
	ShipmentPrices []*ShipmentPrice `json:"shipment_prices"`
}

func (m *GetShipmentPricesResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipmentPriceRequest struct {
	Name                string                             `json:"name"`
	ShipmentPriceListID dot.ID                             `json:"shipment_price_list_id"`
	ShipmentServiceID   dot.ID                             `json:"shipment_service_id"`
	CustomRegionTypes   []route_type.CustomRegionRouteType `json:"custom_region_types"`
	CustomRegionIDs     []dot.ID                           `json:"custom_region_ids"`
	RegionTypes         []route_type.RegionRouteType       `json:"region_types"`
	ProvinceTypes       []route_type.ProvinceRouteType     `json:"province_types"`
	UrbanTypes          []route_type.UrbanType             `json:"urban_types"`
	PriorityPoint       int                                `json:"priority_point"`
	Details             []*PricingDetail                   `json:"details"`
	AdditionalFees      []*AdditionalFee                   `json:"additional_fees"`
}

func (m *CreateShipmentPriceRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipmentPriceRequest struct {
	ID                  dot.ID                             `json:"id"`
	Name                string                             `json:"name"`
	ShipmentPriceListID dot.ID                             `json:"shipment_price_list_id"`
	ShipmentServiceID   dot.ID                             `json:"shipment_service_id"`
	CustomRegionTypes   []route_type.CustomRegionRouteType `json:"custom_region_types"`
	CustomRegionIDs     []dot.ID                           `json:"custom_region_ids"`
	RegionTypes         []route_type.RegionRouteType       `json:"region_types"`
	ProvinceTypes       []route_type.ProvinceRouteType     `json:"province_types"`
	UrbanTypes          []route_type.UrbanType             `json:"urban_types"`
	PriorityPoint       int                                `json:"priority_point"`
	Details             []*PricingDetail                   `json:"details"`
	Status              status3.Status                     `json:"status"`
	AdditionalFees      []*AdditionalFee                   `json:"additional_fees"`
}

func (m *UpdateShipmentPriceRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPricesRequest struct {
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`
	ShipmentServiceID   dot.ID `json:"shipment_service_id"`
}

func (m *GetShipmentPricesRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipmentPricesPriorityPointRequest struct {
	ShipmentPrices []*UpdateShipmentPricePriorityPointRequest `json:"shipment_prices"`
}

func (m *UpdateShipmentPricesPriorityPointRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateShipmentPricePriorityPointRequest struct {
	ID            dot.ID `json:"id"`
	PriorityPoint int    `json:"priority_point"`
}

func (m *UpdateShipmentPricePriorityPointRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CustomRegion struct {
	ID            dot.ID    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ProvinceCodes []string  `json:"province_codes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (m *CustomRegion) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomRegionRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ProvinceCodes []string `json:"province_codes"`
}

func (m *CreateCustomRegionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomRegionRequest struct {
	ID            dot.ID   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ProvinceCodes []string `json:"province_codes"`
}

func (m *UpdateCustomRegionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomRegionsResponse struct {
	CustomRegions []*CustomRegion `json:"custom_regions"`
}

func (m *GetCustomRegionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentService struct {
	ID                 dot.ID               `json:"id"`
	ConnectionID       dot.ID               `json:"connection_id"`
	Name               string               `json:"name"`
	EdCode             string               `json:"ed_code"`
	ServiceIDs         []string             `json:"service_ids"`
	Description        string               `json:"description"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	Status             status3.Status       `json:"status"`
	ImageURL           string               `json:"image_url"`
	AvailableLocations []*AvailableLocation `json:"available_locations"`
	BlacklistLocations []*BlacklistLocation `json:"blacklist_locations"`
	OtherCondition     *OtherCondition      `json:"other_condition"`
}

func (m *ShipmentService) String() string { return jsonx.MustMarshalToString(m) }

type AvailableLocation struct {
	FilterType           filter_type.FilterType             `json:"filter_type"`
	ShippingLocationType location_type.ShippingLocationType `json:"shipping_location_type"`
	RegionTypes          []location_type.RegionType         `json:"regions"`
	CustomRegionIDs      []dot.ID                           `json:"custom_region_ids"`
	ProvinceCodes        []string                           `json:"province_codes"`
}

func (m *AvailableLocation) String() string { return jsonx.MustMarshalToString(m) }

type BlacklistLocation struct {
	ShippingLocationType location_type.ShippingLocationType `json:"shipping_location_type"`
	ProvinceCodes        []string                           `json:"province_codes"`
	DistrictCodes        []string                           `json:"district_codes"`
	WardCodes            []string                           `json:"ward_codes"`
	Reason               string                             `json:"reason"`
}

func (m *BlacklistLocation) String() string { return jsonx.MustMarshalToString(m) }

type OtherCondition struct {
	MinWeight int `json:"min_weight"`
	MaxWeight int `json:"max_weight"`
}

func (m *OtherCondition) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentServicesRequest struct {
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *GetShipmentServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentServicesResponse struct {
	ShipmentServices []*ShipmentService `json:"shipment_services"`
}

func (m *GetShipmentServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipmentServiceRequest struct {
	ConnectionID       dot.ID               `json:"connection_id"`
	Name               string               `json:"name"`
	EdCode             string               `json:"ed_code"`
	ServiceIDs         []string             `json:"service_ids"`
	Description        string               `json:"description"`
	ImageURL           string               `json:"image_url"`
	OtherCondition     *OtherCondition      `json:"other_condition"`
	AvailableLocations []*AvailableLocation `json:"available_locations"`
	BlacklistLocations []*BlacklistLocation `json:"blacklist_locations"`
}

func (m *CreateShipmentServiceRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipmentServiceRequest struct {
	ID             dot.ID             `json:"id"`
	ConnectionID   dot.ID             `json:"connection_id"`
	Name           string             `json:"name"`
	EdCode         string             `json:"ed_code"`
	ServiceIDs     []string           `json:"service_ids"`
	Description    string             `json:"description"`
	ImageURL       string             `json:"image_url"`
	Status         status3.NullStatus `json:"status"`
	OtherCondition *OtherCondition    `json:"other_condition"`
}

func (m *UpdateShipmentServiceRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipmentServicesAvailableLocationsRequest struct {
	IDs                []dot.ID             `json:"ids"`
	AvailableLocations []*AvailableLocation `json:"available_locations"`
}

func (m *UpdateShipmentServicesAvailableLocationsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateShipmentServicesBlacklistLocationsRequest struct {
	IDs                []dot.ID             `json:"ids"`
	BlacklistLocations []*BlacklistLocation `json:"blacklist_locations"`
}

func (m *UpdateShipmentServicesBlacklistLocationsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type ShipmentPriceList struct {
	ID           dot.ID    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsDefault    bool      `json:"is_default"`
	ConnectionID dot.ID    `json:"connection_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (m *ShipmentPriceList) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPriceListsRequest struct {
	ConnectionID dot.ID       `json:"connection_id"`
	IsDefault    dot.NullBool `json:"is_default"`
}

func (m *GetShipmentPriceListsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPriceListsResponse struct {
	ShipmentPriceLists []*ShipmentPriceList `json:"shipment_price_lists"`
}

func (m *GetShipmentPriceListsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipmentPriceListRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsDefault    bool   `json:"is_default"`
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *CreateShipmentPriceListRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipmentPriceListRequest struct {
	ID          dot.ID `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m *UpdateShipmentPriceListRequest) String() string { return jsonx.MustMarshalToString(m) }

type ActiveShipmentPriceListRequest struct {
	ID           dot.ID `json:"id"`
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *ActiveShipmentPriceListRequest) String() string { return jsonx.MustMarshalToString(m) }

// Use to test shipment prices
type GetShippingServicesRequest struct {
	// AccountID
	//
	// dùng để tính cước phí cho 1 shop (trường hợp gắn bảng giá cho shop)
	AccountID dot.ID `json:"account_id"`
	// ShipmentPriceListID
	//
	// trường hợp tính cước phí cho 1 bảng giá cụ thể, nếu field này rỗng mới quan tâm tới AccountID để lấy bảng giá
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`

	ConnectionIDs    []dot.ID     `json:"connection_ids"`
	FromProvinceCode string       `json:"from_province_code"`
	FromDistrictCode string       `json:"from_district_code"`
	ToProvinceCode   string       `json:"to_province_code"`
	ToDistrictCode   string       `json:"to_district_code"`
	GrossWeight      int          `json:"gross_weight"`
	TotalCodAmount   int          `json:"total_cod_amount"`
	BasketValue      int          `json:"basket_value"`
	IncludeInsurance dot.NullBool `json:"include_insurance"`
}

func (m *GetShippingServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShopShipmentPriceList struct {
	ShopID              dot.ID    `json:"shop_id"`
	ConnectionID        dot.ID    `json:"connection_id"`
	ShipmentPriceListID dot.ID    `json:"shipment_price_list_id"`
	Note                string    `json:"note"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           dot.ID    `json:"updated_by"`
}

func (m *ShopShipmentPriceList) String() string { return jsonx.MustMarshalToString(m) }

type GetShopShipmentPriceListRequest struct {
	ShopID              dot.ID `json:"shop_id"`
	ConnectionID        dot.ID `json:"connection_id"`
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`
}

func (m *GetShopShipmentPriceListRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopShipmentPriceListsRequest struct {
	ShopID              dot.ID         `json:"shop_id"`
	ShipmentPriceListID dot.ID         `json:"shipment_price_list_id"`
	ConnectionID        dot.ID         `json:"connection_id"`
	Paging              *common.Paging `json:"paging"`
}

func (m *GetShopShipmentPriceListsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopShipmentPriceListsResponse struct {
	PriceLists []*ShopShipmentPriceList `json:"price_lists"`
	Paging     *common.PageInfo         `json:"paging"`
}

func (m *GetShopShipmentPriceListsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateShopShipmentPriceList struct {
	ShopID              dot.ID `json:"shop_id"`
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`
	ConnectionID        dot.ID `json:"connection_id"`
	Note                string `json:"note"`
}

func (m *CreateShopShipmentPriceList) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShopShipmentPriceListRequest struct {
	ShopID              dot.ID `json:"shop_id"`
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`
	ConnectionID        dot.ID `json:"connection_id"`
	Note                string `json:"note"`
}

func (m *UpdateShopShipmentPriceListRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentPriceListPromotion struct {
	ID                  dot.ID                          `json:"id"`
	ShipmentPriceListID dot.ID                          `json:"shipment_price_list_id"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	Status              status3.Status                  `json:"status"`
	DateFrom            time.Time                       `json:"date_from"`
	DateTo              time.Time                       `json:"date_to"`
	AppliedRules        *PriceListPromotionAppliedRules `json:"applied_rules"`
	CreatedAt           time.Time                       `json:"created_at"`
	UpdatedAt           time.Time                       `json:"updated_at"`
	ConnectionID        dot.ID                          `json:"connection_id"`
	PriorityPoint       int                             `json:"priority_point"`
}

func (m *ShipmentPriceListPromotion) String() string { return jsonx.MustMarshalToString(m) }

type PriceListPromotionAppliedRules struct {
	// apply cho những đơn có điểm lấy hàng nằm trong vùng tự định nghĩa này
	FromCustomRegionIDs []dot.ID `json:"from_custom_region_ids"`
	// apply cho những shop có ngày tạo trong khoảng này
	ShopCreatedDate filter.Date `json:"shop_created_date"`
	// apply cho những user có ngày tạo trong khoảng này
	UserCreatedDate filter.Date `json:"user_created_date"`
	// apply cho những shop đang xài bảng giá này
	UsingPriceListIDs []dot.ID `json:"using_price_list_ids"`
}

func (m *PriceListPromotionAppliedRules) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPriceListPromotionsRequest struct {
	ShipmentPriceListID dot.ID         `json:"shipment_price_list_id"`
	ConnectionID        dot.ID         `json:"connection_id"`
	Paging              *common.Paging `json:"paging"`
}

func (m *GetShipmentPriceListPromotionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipmentPriceListPromotionsResponse struct {
	ShipmentPriceListPromotions []*ShipmentPriceListPromotion `json:"shipment_price_list_promotions"`
}

func (m *GetShipmentPriceListPromotionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipmentPriceListPromotionRequest struct {
	ShipmentPriceListID dot.ID                          `json:"shipment_price_list_id"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	ConnectionID        dot.ID                          `json:"connection_id"`
	DateFrom            time.Time                       `json:"date_from"`
	DateTo              time.Time                       `json:"date_to"`
	AppliedRules        *PriceListPromotionAppliedRules `json:"applied_rules"`
	PriorityPoint       int                             `json:"priority_point"`
}

func (m *CreateShipmentPriceListPromotionRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateShipmentPriceListPromotionRequest struct {
	ID                  dot.ID                          `json:"id"`
	Name                string                          `json:"name"`
	Description         string                          `json:"description"`
	DateFrom            time.Time                       `json:"date_from"`
	DateTo              time.Time                       `json:"date_to"`
	AppliedRules        *PriceListPromotionAppliedRules `json:"applied_rules"`
	PriorityPoint       int                             `json:"priority_point"`
	Status              status3.NullStatus              `json:"status"`
	ConnectionID        dot.ID                          `json:"connection_id"`
	ShipmentPriceListID dot.ID                          `json:"shipment_price_list_id"`
}

func (m *UpdateShipmentPriceListPromotionRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UserResponse struct {
	Users  []*etop.User           `json:"users"`
	Paging *common.CursorPageInfo `json:"paging"`
}

func (m *UserResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateAdminUserRequest struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

func (m *CreateAdminUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateAdminUserResponse struct {
	UserId dot.ID         `json:"user_id"`
	Roles  []string       `json:"roles"`
	Status status3.Status `json:"status"`
}

func (m *CreateAdminUserResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAdminUserRequest struct {
	UserId dot.ID         `json:"user_id"`
	Roles  []string       `json:"roles"`
	Status status3.Status `json:"status"`
}

func (m *UpdateAdminUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAdminUserResponse struct {
	UserId dot.ID         `json:"user_id"`
	Roles  []string       `json:"roles"`
	Status status3.Status `json:"status"`
}

func (m *UpdateAdminUserResponse) String() string { return jsonx.MustMarshalToString(m) }

type AdminUserFilter struct {
	Roles filter.Strings `json:"roles"`
}

type GetAdminUsersRequest struct {
	Filter AdminUserFilter `json:"filter"`
}

func (g *GetAdminUsersRequest) String() string {
	return jsonx.MustMarshalToString(g)
}

type AdminAccountResponse struct {
	UserId    dot.ID    `json:"user_id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAdminUserResponse struct {
	Admins []*AdminAccountResponse `json:"admins"`
}

func (g *GetAdminUserResponse) String() string {
	return jsonx.MustMarshalToString(g)
}

type DeleteAdminUserRequest struct {
	UserID dot.ID `json:"user_id"`
}

func (d *DeleteAdminUserRequest) String() string {
	return jsonx.MustMarshalToString(d)
}

type DeleteAdminUserResponse struct {
	Updated int `json:"updated"`
}

func (d *DeleteAdminUserResponse) String() string {
	return jsonx.MustMarshalToString(d)
}

type UpdateShopInfoRequest struct {
	ID dot.ID `json:"id"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule   string       `json:"money_transaction_rrule"`
	IsPriorMoneyTransaction dot.NullBool `json:"is_prior_money_transaction"`
}

func (d *UpdateShopInfoRequest) String() string {
	return jsonx.MustMarshalToString(d)
}

type SplitMoneyTxShippingExternalRequest struct {
	ID dot.ID `json:"id"`
	// Tách phiên dựa theo độ ưu tiên đối soát với shop
	IsSplitByShopPriority bool `json:"is_split_by_shop_priority"`
	// Tách phiên dựa theo shop mới (shop có số lượng phiên <= max_money_tx_shipping_count)
	MaxMoneyTxShippingCount int `json:"max_money_tx_shipping_count"`
}

func (d *SplitMoneyTxShippingExternalRequest) String() string {
	return jsonx.MustMarshalToString(d)
}

type GetAPIKeyRequest struct {
	AccountID dot.ID `json:"account_id"`
}

func (m *GetAPIKeyRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetAPIKeyResponse struct {
	ApiKey string `json:"api_key"`
}

func (m *GetAPIKeyResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetAccountUsersRequest struct {
	Paging *common.CursorPaging          `json:"paging"`
	Filter *FilterGetAccountUsersRequest `json:"filter"`
}

func (m *GetAccountUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetAccountUsersRequest struct {
	Name       filter.FullTextSearch     `json:"name"`
	Phone      filter.FullTextSearch     `json:"phone"`
	AccountID  dot.ID                    `json:"account_id"`
	Roles      []shop_user_role.UserRole `json:"roles"`
	ExactRoles []shop_user_role.UserRole `json:"exact_roles"`
}

func (m *FilterGetAccountUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetAccountUsersResponse struct {
	AccountUsers []*shoptypes.AccountUserExtended `json:"account_users"`
	Paging       *common.CursorPageInfo           `json:"paging"`
}

func (m *GetAccountUsersResponse) String() string { return jsonx.MustMarshalToString(m) }
