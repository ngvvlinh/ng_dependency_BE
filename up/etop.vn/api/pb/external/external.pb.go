package external

import (
	common "etop.vn/api/pb/common"
	etop "etop.vn/api/pb/etop"
	gender "etop.vn/api/pb/etop/etc/gender"
	shipping "etop.vn/api/pb/etop/etc/shipping"
	shipping_provider "etop.vn/api/pb/etop/etc/shipping_provider"
	status3 "etop.vn/api/pb/etop/etc/status3"
	status4 "etop.vn/api/pb/etop/etc/status4"
	status5 "etop.vn/api/pb/etop/etc/status5"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	order "etop.vn/api/pb/etop/order"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type Partner struct {
	Id         dot.ID           `json:"id"`
	Name       string           `json:"name"`
	PublicName string           `json:"public_name"`
	Type       etop.AccountType `json:"type"`
	Phone      string           `json:"phone"`
	// only domain, no scheme
	Website         string   `json:"website"`
	WebsiteUrl      string   `json:"website_url"`
	ImageUrl        string   `json:"image_url"`
	Email           string   `json:"email"`
	RecognizedHosts []string `json:"recognized_hosts"`
	RedirectUrls    []string `json:"redirect_urls"`
}

func (m *Partner) Reset()         { *m = Partner{} }
func (m *Partner) String() string { return jsonx.MustMarshalToString(m) }

type CreateWebhookRequest struct {
	Entities []string `json:"entities"`
	Fields   []string `json:"fields"`
	Url      string   `json:"url"`
	Metadata string   `json:"metadata"`
}

func (m *CreateWebhookRequest) Reset()         { *m = CreateWebhookRequest{} }
func (m *CreateWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteWebhookRequest struct {
	Id dot.ID `json:"id"`
}

func (m *DeleteWebhookRequest) Reset()         { *m = DeleteWebhookRequest{} }
func (m *DeleteWebhookRequest) String() string { return jsonx.MustMarshalToString(m) }

type WebhooksResponse struct {
	Webhooks []*Webhook `json:"webhooks"`
}

func (m *WebhooksResponse) Reset()         { *m = WebhooksResponse{} }
func (m *WebhooksResponse) String() string { return jsonx.MustMarshalToString(m) }

type Webhook struct {
	Id        dot.ID         `json:"id"`
	Entities  []string       `json:"entities"`
	Fields    []string       `json:"fields"`
	Url       string         `json:"url"`
	Metadata  string         `json:"metadata"`
	CreatedAt dot.Time       `json:"created_at"`
	States    *WebhookStates `json:"states"`
}

func (m *Webhook) Reset()         { *m = Webhook{} }
func (m *Webhook) String() string { return jsonx.MustMarshalToString(m) }

type WebhookStates struct {
	State      string        `json:"state"`
	LastSentAt dot.Time      `json:"last_sent_at"`
	LastError  *WebhookError `json:"last_error"`
}

func (m *WebhookStates) Reset()         { *m = WebhookStates{} }
func (m *WebhookStates) String() string { return jsonx.MustMarshalToString(m) }

type WebhookError struct {
	Error      string   `json:"error"`
	RespStatus int      `json:"resp_status"`
	RespBody   string   `json:"resp_body"`
	Retried    int      `json:"retried"`
	SentAt     dot.Time `json:"sent_at"`
}

func (m *WebhookError) Reset()         { *m = WebhookError{} }
func (m *WebhookError) String() string { return jsonx.MustMarshalToString(m) }

type GetChangesRequest struct {
	Paging   *common.ForwardPaging `json:"paging"`
	Entity   *string               `json:"entity"`
	EntityId *string               `json:"entity_id"`
}

func (m *GetChangesRequest) Reset()         { *m = GetChangesRequest{} }
func (m *GetChangesRequest) String() string { return jsonx.MustMarshalToString(m) }

type Callback struct {
	Id      dot.ID    `json:"id"`
	Changes []*Change `json:"changes"`
}

func (m *Callback) Reset()         { *m = Callback{} }
func (m *Callback) String() string { return jsonx.MustMarshalToString(m) }

// ChangesData serialize changes data for storing in MongoDB
type ChangesData struct {
	// for using with mongodb
	XId       dot.ID    `json:"_id"`
	WebhookId dot.ID    `json:"webhook_id"`
	AccountId dot.ID    `json:"account_id"`
	CreatedAt dot.Time  `json:"created_at"`
	Changes   []*Change `json:"changes"`
}

func (m *ChangesData) Reset()         { *m = ChangesData{} }
func (m *ChangesData) String() string { return jsonx.MustMarshalToString(m) }

type Change struct {
	Time       dot.Time     `json:"time"`
	ChangeType string       `json:"change_type"`
	Entity     string       `json:"entity"`
	Latest     *LatestOneOf `json:"latest"`
	Changed    *ChangeOneOf `json:"changed"`
}

func (m *Change) Reset()         { *m = Change{} }
func (m *Change) String() string { return jsonx.MustMarshalToString(m) }

type LatestOneOf struct {
	Order       *Order       `json:"order"`
	Fulfillment *Fulfillment `json:"fulfillment"`
}

func (m *LatestOneOf) Reset()         { *m = LatestOneOf{} }
func (m *LatestOneOf) String() string { return jsonx.MustMarshalToString(m) }

type ChangeOneOf struct {
	Order       *Order       `json:"order"`
	Fulfillment *Fulfillment `json:"fulfillment"`
}

func (m *ChangeOneOf) Reset()         { *m = ChangeOneOf{} }
func (m *ChangeOneOf) String() string { return jsonx.MustMarshalToString(m) }

type CancelOrderRequest struct {
	Id           dot.ID `json:"id"`
	Code         string `json:"code"`
	ExternalId   string `json:"external_id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelOrderRequest) Reset()         { *m = CancelOrderRequest{} }
func (m *CancelOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *OrderIDRequest) Reset()         { *m = OrderIDRequest{} }
func (m *OrderIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentIDRequest struct {
	Id           dot.ID `json:"id"`
	ShippingCode string `json:"shipping_code"`
}

func (m *FulfillmentIDRequest) Reset()         { *m = FulfillmentIDRequest{} }
func (m *FulfillmentIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderAndFulfillments struct {
	Order             *Order          `json:"order"`
	Fulfillments      []*Fulfillment  `json:"fulfillments"`
	FulfillmentErrors []*common.Error `json:"fulfillment_errors"`
}

func (m *OrderAndFulfillments) Reset()         { *m = OrderAndFulfillments{} }
func (m *OrderAndFulfillments) String() string { return jsonx.MustMarshalToString(m) }

type Order struct {
	Id                        dot.ID                `json:"id"`
	ShopId                    dot.ID                `json:"shop_id"`
	Code                      *string               `json:"code"`
	ExternalId                *string               `json:"external_id"`
	ExternalCode              *string               `json:"external_code"`
	ExternalUrl               *string               `json:"external_url"`
	SelfUrl                   *string               `json:"self_url"`
	CustomerAddress           *OrderAddress         `json:"customer_address"`
	ShippingAddress           *OrderAddress         `json:"shipping_address"`
	CreatedAt                 dot.Time              `json:"created_at"`
	ProcessedAt               dot.Time              `json:"processed_at"`
	UpdatedAt                 dot.Time              `json:"updated_at"`
	ClosedAt                  dot.Time              `json:"closed_at"`
	ConfirmedAt               dot.Time              `json:"confirmed_at"`
	CancelledAt               dot.Time              `json:"cancelled_at"`
	CancelReason              *string               `json:"cancel_reason"`
	ConfirmStatus             *status3.Status       `json:"confirm_status"`
	Status                    *status5.Status       `json:"status"`
	FulfillmentShippingStatus *status5.Status       `json:"fulfillment_shipping_status"`
	EtopPaymentStatus         *status4.Status       `json:"etop_payment_status"`
	Lines                     []*OrderLine          `json:"lines"`
	TotalItems                *int                  `json:"total_items"`
	BasketValue               *int                  `json:"basket_value"`
	OrderDiscount             *int                  `json:"order_discount"`
	TotalDiscount             *int                  `json:"total_discount"`
	TotalFee                  *int                  `json:"total_fee"`
	FeeLines                  []*order.OrderFeeLine `json:"fee_lines"`
	TotalAmount               *int                  `json:"total_amount"`
	OrderNote                 *string               `json:"order_note"`
	Shipping                  *OrderShipping        `json:"shipping"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return jsonx.MustMarshalToString(m) }

type OrderShipping struct {
	PickupAddress       *OrderAddress                       `json:"pickup_address"`
	ReturnAddress       *OrderAddress                       `json:"return_address"`
	ShippingServiceName *string                             `json:"shipping_service_name"`
	ShippingServiceCode *string                             `json:"shipping_service_code"`
	ShippingServiceFee  *int                                `json:"shipping_service_fee"`
	Carrier             *shipping_provider.ShippingProvider `json:"carrier"`
	IncludeInsurance    *bool                               `json:"include_insurance"`
	TryOn               *try_on.TryOnCode                   `json:"try_on"`
	ShippingNote        *string                             `json:"shipping_note"`
	CodAmount           *int                                `json:"cod_amount"`
	GrossWeight         *int                                `json:"gross_weight"`
	Length              *int                                `json:"length"`
	Width               *int                                `json:"width"`
	Height              *int                                `json:"height"`
	ChargeableWeight    *int                                `json:"chargeable_weight"`
}

func (m *OrderShipping) Reset()         { *m = OrderShipping{} }
func (m *OrderShipping) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrderRequest struct {
	ExternalId      string        `json:"external_id"`
	ExternalCode    string        `json:"external_code"`
	ExternalUrl     string        `json:"external_url"`
	CustomerAddress *OrderAddress `json:"customer_address"`
	ShippingAddress *OrderAddress `json:"shipping_address"`
	Lines           []*OrderLine  `json:"lines"`
	TotalItems      int           `json:"total_items"`
	// basket_value = SUM(lines.retail_price)
	BasketValue   int                   `json:"basket_value"`
	OrderDiscount int                   `json:"order_discount"`
	TotalDiscount int                   `json:"total_discount"`
	TotalFee      *int                  `json:"total_fee"`
	FeeLines      []*order.OrderFeeLine `json:"fee_lines"`
	TotalAmount   int                   `json:"total_amount"`
	OrderNote     string                `json:"order_note"`
	Shipping      *OrderShipping        `json:"shipping"`
	ExternalMeta  map[string]string     `json:"external_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *CreateOrderRequest) Reset()         { *m = CreateOrderRequest{} }
func (m *CreateOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderLine struct {
	VariantId   dot.ID `json:"variant_id"`
	ProductId   dot.ID `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	ListPrice   int    `json:"list_price"`
	RetailPrice int    `json:"retail_price"`
	// payment_price = retail_price - discount_per_item
	PaymentPrice *int               `json:"payment_price"`
	ImageUrl     string             `json:"image_url"`
	Attributes   []*order.Attribute `json:"attributes"`
}

func (m *OrderLine) Reset()         { *m = OrderLine{} }
func (m *OrderLine) String() string { return jsonx.MustMarshalToString(m) }

type Fulfillment struct {
	Id                       dot.ID                              `json:"id"`
	OrderId                  dot.ID                              `json:"order_id"`
	ShopId                   dot.ID                              `json:"shop_id"`
	SelfUrl                  *string                             `json:"self_url"`
	TotalItems               *int                                `json:"total_items"`
	BasketValue              *int                                `json:"basket_value"`
	CreatedAt                dot.Time                            `json:"created_at"`
	UpdatedAt                dot.Time                            `json:"updated_at"`
	ClosedAt                 dot.Time                            `json:"closed_at"`
	CancelledAt              dot.Time                            `json:"cancelled_at"`
	CancelReason             *string                             `json:"cancel_reason"`
	Carrier                  *shipping_provider.ShippingProvider `json:"carrier"`
	ShippingServiceName      *string                             `json:"shipping_service_name"`
	ShippingServiceFee       *int                                `json:"shipping_service_fee"`
	ActualShippingServiceFee *int                                `json:"actual_shipping_service_fee"`
	ShippingServiceCode      *string                             `json:"shipping_service_code"`
	ShippingCode             *string                             `json:"shipping_code"`
	ShippingNote             *string                             `json:"shipping_note"`
	TryOn                    *try_on.TryOnCode                   `json:"try_on"`
	IncludeInsurance         *bool                               `json:"include_insurance"`
	ConfirmStatus            *status3.Status                     `json:"confirm_status"`
	ShippingState            *shipping.State                     `json:"shipping_state"`
	ShippingStatus           *status5.Status                     `json:"shipping_status"`
	Status                   *status5.Status                     `json:"status"`
	CodAmount                *int                                `json:"cod_amount"`
	ActualCodAmount          *int                                `json:"actual_cod_amount"`
	ChargeableWeight         *int                                `json:"chargeable_weight"`
	PickupAddress            *OrderAddress                       `json:"pickup_address"`
	ReturnAddress            *OrderAddress                       `json:"return_address"`
	ShippingAddress          *OrderAddress                       `json:"shipping_address"`
	EtopPaymentStatus        *status4.Status                     `json:"etop_payment_status"`
	EstimatedDeliveryAt      dot.Time                            `json:"estimated_delivery_at"`
	EstimatedPickupAt        dot.Time                            `json:"estimated_pickup_at"`
}

func (m *Fulfillment) Reset()         { *m = Fulfillment{} }
func (m *Fulfillment) String() string { return jsonx.MustMarshalToString(m) }

type Ward struct {
	Name string `json:"name"`
}

func (m *Ward) Reset()         { *m = Ward{} }
func (m *Ward) String() string { return jsonx.MustMarshalToString(m) }

type District struct {
	Name  string `json:"name"`
	Wards []Ward `json:"wards"`
}

func (m *District) Reset()         { *m = District{} }
func (m *District) String() string { return jsonx.MustMarshalToString(m) }

type Province struct {
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

func (m *Province) Reset()         { *m = Province{} }
func (m *Province) String() string { return jsonx.MustMarshalToString(m) }

type LocationResponse struct {
	Provinces []Province `json:"provinces"`
}

func (m *LocationResponse) Reset()         { *m = LocationResponse{} }
func (m *LocationResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetShippingServicesRequest struct {
	PickupAddress   *LocationAddress `json:"pickup_address"`
	ShippingAddress *LocationAddress `json:"shipping_address"`
	// in gram (g)
	GrossWeight int `json:"gross_weight"`
	// in gram (g)
	ChargeableWeight int `json:"chargeable_weight"`
	// in centimetre (cm)
	Length int `json:"length"`
	// in centimetre (cm)
	Width int `json:"width"`
	// in centimetre (cm)
	Height           int   `json:"height"`
	BasketValue      int   `json:"basket_value"`
	CodAmount        int   `json:"cod_amount"`
	IncludeInsurance *bool `json:"include_insurance"`
}

func (m *GetShippingServicesRequest) Reset()         { *m = GetShippingServicesRequest{} }
func (m *GetShippingServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type LocationAddress struct {
	Province string `json:"province"`
	District string `json:"district"`
}

func (m *LocationAddress) Reset()         { *m = LocationAddress{} }
func (m *LocationAddress) String() string { return jsonx.MustMarshalToString(m) }

type GetShippingServicesResponse struct {
	Services []*ShippingService `json:"services"`
}

func (m *GetShippingServicesResponse) Reset()         { *m = GetShippingServicesResponse{} }
func (m *GetShippingServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShippingService struct {
	Code                string                             `json:"code"`
	Name                string                             `json:"name"`
	Fee                 int                                `json:"fee"`
	Carrier             shipping_provider.ShippingProvider `json:"carrier"`
	EstimatedPickupAt   dot.Time                           `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time                           `json:"estimated_delivery_at"`
}

func (m *ShippingService) Reset()         { *m = ShippingService{} }
func (m *ShippingService) String() string { return jsonx.MustMarshalToString(m) }

type OrderCustomer struct {
	FullName string        `json:"full_name"`
	Email    string        `json:"email"`
	Phone    string        `json:"phone"`
	Gender   gender.Gender `json:"gender"`
}

func (m *OrderCustomer) Reset()         { *m = OrderCustomer{} }
func (m *OrderCustomer) String() string { return jsonx.MustMarshalToString(m) }

type OrderAddress struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Company  string `json:"company"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
}

func (m *OrderAddress) Reset()         { *m = OrderAddress{} }
func (m *OrderAddress) String() string { return jsonx.MustMarshalToString(m) }
