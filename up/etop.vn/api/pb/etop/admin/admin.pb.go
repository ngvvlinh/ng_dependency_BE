package admin

import (
	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	credit_type "etop.vn/api/pb/etop/etc/credit_type"
	notifier_entity "etop.vn/api/pb/etop/etc/notifier_entity"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	status3 "etop.vn/api/pb/etop/etc/status3"
	_ "etop.vn/api/pb/etop/order"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type GetOrdersRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetOrdersRequest) Reset()         { *m = GetOrdersRequest{} }
func (m *GetOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetFulfillmentsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	ShopId  dot.ID           `json:"shop_id"`
	OrderId dot.ID           `json:"order_id"`
	Status  *status3.Status  `json:"status"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetFulfillmentsRequest) Reset()         { *m = GetFulfillmentsRequest{} }
func (m *GetFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginAsAccountRequest struct {
	UserId    dot.ID `json:"user_id"`
	AccountId dot.ID `json:"account_id"`
	// This is a sensitive API, so admin must provide password before processing!
	Password string `json:"password"`
}

func (m *LoginAsAccountRequest) Reset()         { *m = LoginAsAccountRequest{} }
func (m *LoginAsAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetMoneyTransactionsRequest struct {
	Ids                                []dot.ID         `json:"ids"`
	ShopId                             dot.ID           `json:"shop_id"`
	MoneyTransactionShippingExternalId dot.ID           `json:"money_transaction_shipping_external_id"`
	Paging                             *common.Paging   `json:"paging"`
	Filters                            []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionsRequest) Reset()         { *m = GetMoneyTransactionsRequest{} }
func (m *GetMoneyTransactionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveFfmsMoneyTransactionRequest struct {
	FulfillmentIds     []dot.ID `json:"fulfillment_ids"`
	MoneyTransactionId dot.ID   `json:"money_transaction_id"`
	ShopId             dot.ID   `json:"shop_id"`
}

func (m *RemoveFfmsMoneyTransactionRequest) Reset()         { *m = RemoveFfmsMoneyTransactionRequest{} }
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
	TotalCod           dot.ID `json:"total_cod"`
	TotalAmount        dot.ID `json:"total_amount"`
	TotalOrders        dot.ID `json:"total_orders"`
}

func (m *ConfirmMoneyTransactionRequest) Reset()         { *m = ConfirmMoneyTransactionRequest{} }
func (m *ConfirmMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteMoneyTransactionRequest struct {
	MoneyTransactionId dot.ID `json:"money_transaction_id"`
	ShopId             dot.ID `json:"shop_id"`
}

func (m *DeleteMoneyTransactionRequest) Reset()         { *m = DeleteMoneyTransactionRequest{} }
func (m *DeleteMoneyTransactionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopsRequest struct {
	Paging *common.Paging `json:"paging"`
}

func (m *GetShopsRequest) Reset()         { *m = GetShopsRequest{} }
func (m *GetShopsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopsResponse struct {
	Paging *common.PageInfo `json:"paging"`
	Shops  []*etop.Shop     `json:"shops"`
}

func (m *GetShopsResponse) Reset()         { *m = GetShopsResponse{} }
func (m *GetShopsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetCreditRequest struct {
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
}

func (m *GetCreditRequest) Reset()         { *m = GetCreditRequest{} }
func (m *GetCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCreditsRequest struct {
	ShopId dot.ID         `json:"shop_id"`
	Paging *common.Paging `json:"paging"`
}

func (m *GetCreditsRequest) Reset()         { *m = GetCreditsRequest{} }
func (m *GetCreditsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCreditRequest struct {
	Amount dot.ID                  `json:"amount"`
	ShopId dot.ID                  `json:"shop_id"`
	Type   *credit_type.CreditType `json:"type"`
	PaidAt dot.Time                `json:"paid_at"`
}

func (m *CreateCreditRequest) Reset()         { *m = CreateCreditRequest{} }
func (m *CreateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCreditRequest struct {
	Id     dot.ID                  `json:"id"`
	Amount dot.ID                  `json:"amount"`
	ShopId dot.ID                  `json:"shop_id"`
	Type   *credit_type.CreditType `json:"type"`
	PaidAt dot.Time                `json:"paid_at"`
}

func (m *UpdateCreditRequest) Reset()         { *m = UpdateCreditRequest{} }
func (m *UpdateCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmCreditRequest struct {
	Id     dot.ID `json:"id"`
	ShopId dot.ID `json:"shop_id"`
}

func (m *ConfirmCreditRequest) Reset()         { *m = ConfirmCreditRequest{} }
func (m *ConfirmCreditRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreatePartnerRequest struct {
	Partner etop.Partner `json:"partner"`
}

func (m *CreatePartnerRequest) Reset()         { *m = CreatePartnerRequest{} }
func (m *CreatePartnerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentRequest struct {
	Id                       dot.ID          `json:"id"`
	FullName                 string          `json:"full_name"`
	Phone                    string          `json:"phone"`
	TotalCodAmount           *int32          `json:"total_cod_amount"`
	IsPartialDelivery        bool            `json:"is_partial_delivery"`
	AdminNote                string          `json:"admin_note"`
	ActualCompensationAmount int32           `json:"actual_compensation_amount"`
	ShippingState            *shipping.State `json:"shipping_state"`
}

func (m *UpdateFulfillmentRequest) Reset()         { *m = UpdateFulfillmentRequest{} }
func (m *UpdateFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type GenerateAPIKeyRequest struct {
	AccountId dot.ID `json:"account_id"`
}

func (m *GenerateAPIKeyRequest) Reset()         { *m = GenerateAPIKeyRequest{} }
func (m *GenerateAPIKeyRequest) String() string { return jsonx.MustMarshalToString(m) }

type GenerateAPIKeyResponse struct {
	AccountId dot.ID `json:"account_id"`
	ApiKey    string `json:"api_key"`
}

func (m *GenerateAPIKeyResponse) Reset()         { *m = GenerateAPIKeyResponse{} }
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
	Ids     []dot.ID         `json:"ids"`
	Status  *status3.Status  `json:"status"`
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetMoneyTransactionShippingEtopsRequest) Reset() {
	*m = GetMoneyTransactionShippingEtopsRequest{}
}
func (m *GetMoneyTransactionShippingEtopsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmMoneyTransactionShippingEtopRequest struct {
	Id          dot.ID `json:"id"`
	TotalCod    dot.ID `json:"total_cod"`
	TotalAmount dot.ID `json:"total_amount"`
	TotalOrders dot.ID `json:"total_orders"`
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

func (m *UpdateMoneyTransactionRequest) Reset()         { *m = UpdateMoneyTransactionRequest{} }
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

func (m *CreateNotificationsRequest) Reset()         { *m = CreateNotificationsRequest{} }
func (m *CreateNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateNotificationsResponse struct {
	Created int32 `json:"created"`
	Errored int32 `json:"errored"`
}

func (m *CreateNotificationsResponse) Reset()         { *m = CreateNotificationsResponse{} }
func (m *CreateNotificationsResponse) String() string { return jsonx.MustMarshalToString(m) }
