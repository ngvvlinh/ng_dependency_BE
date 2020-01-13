package admin

import (
	etop "etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	common "etop.vn/api/top/types/common"
	credit_type "etop.vn/api/top/types/etc/credit_type"
	notifier_entity "etop.vn/api/top/types/etc/notifier_entity"
	shipping "etop.vn/api/top/types/etc/shipping"
	status3 "etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type GetOrdersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsRequest struct {
	Paging  *common.Paging     `json:"paging"`
	ShopId  dot.ID             `json:"shop_id"`
	OrderId dot.ID             `json:"order_id"`
	Status  status3.NullStatus `json:"status"`
	Filters []*common.Filter   `json:"filters"`
}

func (m *GetFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginAsAccountRequest struct {
	UserId    dot.ID `json:"user_id"`
	AccountId dot.ID `json:"account_id"`
	// This is a sensitive API, so admin must provide password before processing!
	Password string `json:"password"`
}

func (m *LoginAsAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionsRequest struct {
	Ids                                []dot.ID         `json:"ids"`
	ShopId                             dot.ID           `json:"shop_id"`
	MoneyTransactionShippingExternalId dot.ID           `json:"money_transaction_shipping_external_id"`
	Paging                             *common.Paging   `json:"paging"`
	Filters                            []*common.Filter `json:"filters"`
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

type GetShopsRequest struct {
	Paging *common.Paging `json:"paging"`
}

func (m *GetShopsRequest) String() string { return jsonx.MustMarshalToString(m) }

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
	ShopId dot.ID         `json:"shop_id"`
	Paging *common.Paging `json:"paging"`
}

func (m *GetCreditsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCreditRequest struct {
	Amount int                    `json:"amount"`
	ShopId dot.ID                 `json:"shop_id"`
	Type   credit_type.CreditType `json:"type"`
	PaidAt dot.Time               `json:"paid_at"`
}

func (m *CreateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCreditRequest struct {
	Id     dot.ID                  `json:"id"`
	Amount int                     `json:"amount"`
	ShopId dot.ID                  `json:"shop_id"`
	Type   *credit_type.CreditType `json:"type"`
	PaidAt dot.Time                `json:"paid_at"`
}

func (m *UpdateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmCreditRequest struct {
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
}

func (m *ConfirmCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreatePartnerRequest struct {
	Partner etop.Partner `json:"partner"`
}

func (m *CreatePartnerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentRequest struct {
	Id                       dot.ID             `json:"id"`
	FullName                 string             `json:"full_name"`
	Phone                    string             `json:"phone"`
	TotalCodAmount           dot.NullInt        `json:"total_cod_amount"`
	IsPartialDelivery        bool               `json:"is_partial_delivery"`
	AdminNote                string             `json:"admin_note"`
	ActualCompensationAmount int                `json:"actual_compensation_amount"`
	ShippingState            shipping.NullState `json:"shipping_state"`
}

func (m *UpdateFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

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
func (m *GetMoneyTransactionShippingEtopsRequest) String() string { return jsonx.MustMarshalToString(m) }

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
}

func (m *UpdateFulfillmentShippingStateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentShippingFeeRequest struct {
	ID           dot.ID                        `json:"id"`
	ShippingCode string                        `json:"shipping_code"`
	ShippingFees []*types.ShippingFeeShortLine `json:"shipping_fees"`
}

func (m *UpdateFulfillmentShippingFeeRequest) String() string { return jsonx.MustMarshalToString(m) }
