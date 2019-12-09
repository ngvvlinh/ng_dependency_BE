package types

import (
	catalogtypes "etop.vn/api/main/catalog/types"
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	common "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/account_type"
	gender "etop.vn/api/top/types/etc/gender"
	shipping "etop.vn/api/top/types/etc/shipping"
	shipping_provider "etop.vn/api/top/types/etc/shipping_provider"
	status3 "etop.vn/api/top/types/etc/status3"
	status4 "etop.vn/api/top/types/etc/status4"
	status5 "etop.vn/api/top/types/etc/status5"
	try_on "etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type Partner struct {
	Id         dot.ID                   `json:"id"`
	Name       string                   `json:"name"`
	PublicName string                   `json:"public_name"`
	Type       account_type.AccountType `json:"type"`
	Phone      string                   `json:"phone"`
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
	Entity   dot.NullString        `json:"entity"`
	EntityId dot.NullString        `json:"entity_id"`
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

type ListFulfillmentsRequest struct {
	OrderID dot.ID               `json:"order_id"`
	Paging  *common.CursorPaging `json:"paging"`
}

func (m *ListFulfillmentsRequest) Reset()         { *m = ListFulfillmentsRequest{} }
func (m *ListFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

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
	Code                      dot.NullString        `json:"code"`
	ExternalId                dot.NullString        `json:"external_id"`
	ExternalCode              dot.NullString        `json:"external_code"`
	ExternalUrl               dot.NullString        `json:"external_url"`
	SelfUrl                   dot.NullString        `json:"self_url"`
	CustomerAddress           *OrderAddress         `json:"customer_address"`
	ShippingAddress           *OrderAddress         `json:"shipping_address"`
	CreatedAt                 dot.Time              `json:"created_at"`
	ProcessedAt               dot.Time              `json:"processed_at"`
	UpdatedAt                 dot.Time              `json:"updated_at"`
	ClosedAt                  dot.Time              `json:"closed_at"`
	ConfirmedAt               dot.Time              `json:"confirmed_at"`
	CancelledAt               dot.Time              `json:"cancelled_at"`
	CancelReason              dot.NullString        `json:"cancel_reason"`
	ConfirmStatus             status3.NullStatus    `json:"confirm_status"`
	Status                    status5.NullStatus    `json:"status"`
	FulfillmentShippingStatus status5.NullStatus    `json:"fulfillment_shipping_status"`
	EtopPaymentStatus         status4.NullStatus    `json:"etop_payment_status"`
	Lines                     []*OrderLine          `json:"lines"`
	TotalItems                dot.NullInt           `json:"total_items"`
	BasketValue               dot.NullInt           `json:"basket_value"`
	OrderDiscount             dot.NullInt           `json:"order_discount"`
	TotalDiscount             dot.NullInt           `json:"total_discount"`
	TotalFee                  dot.NullInt           `json:"total_fee"`
	FeeLines                  []*types.OrderFeeLine `json:"fee_lines"`
	TotalAmount               dot.NullInt           `json:"total_amount"`
	OrderNote                 dot.NullString        `json:"order_note"`
	Shipping                  *OrderShipping        `json:"shipping"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return jsonx.MustMarshalToString(m) }

type OrderShipping struct {
	PickupAddress       *OrderAddress  `json:"pickup_address"`
	ReturnAddress       *OrderAddress  `json:"return_address"`
	ShippingServiceName dot.NullString `json:"shipping_service_name"`
	ShippingServiceCode dot.NullString `json:"shipping_service_code"`
	ShippingServiceFee  dot.NullInt    `json:"shipping_service_fee"`
	// @Deprecated use connection_id instead
	Carrier          shipping_provider.ShippingProvider `json:"carrier"`
	IncludeInsurance dot.NullBool                       `json:"include_insurance"`
	TryOn            try_on.TryOnCode                   `json:"try_on"`
	ShippingNote     dot.NullString                     `json:"shipping_note"`
	CodAmount        dot.NullInt                        `json:"cod_amount"`
	GrossWeight      dot.NullInt                        `json:"gross_weight"`
	Length           dot.NullInt                        `json:"length"`
	Width            dot.NullInt                        `json:"width"`
	Height           dot.NullInt                        `json:"height"`
	ChargeableWeight dot.NullInt                        `json:"chargeable_weight"`
	ConnectionID     dot.ID                             `json:"connection_id"`
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
	TotalFee      dot.NullInt           `json:"total_fee"`
	FeeLines      []*types.OrderFeeLine `json:"fee_lines"`
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
	PaymentPrice dot.NullInt               `json:"payment_price"`
	ImageUrl     string                    `json:"image_url"`
	Attributes   []*catalogtypes.Attribute `json:"attributes"`
}

func (m *OrderLine) Reset()         { *m = OrderLine{} }
func (m *OrderLine) String() string { return jsonx.MustMarshalToString(m) }

type Fulfillment struct {
	Id                       dot.ID                             `json:"id"`
	OrderId                  dot.ID                             `json:"order_id"`
	ShopId                   dot.ID                             `json:"shop_id"`
	SelfUrl                  dot.NullString                     `json:"self_url"`
	TotalItems               dot.NullInt                        `json:"total_items"`
	BasketValue              dot.NullInt                        `json:"basket_value"`
	CreatedAt                dot.Time                           `json:"created_at"`
	UpdatedAt                dot.Time                           `json:"updated_at"`
	ClosedAt                 dot.Time                           `json:"closed_at"`
	CancelledAt              dot.Time                           `json:"cancelled_at"`
	CancelReason             dot.NullString                     `json:"cancel_reason"`
	Carrier                  shipping_provider.ShippingProvider `json:"carrier"`
	ShippingServiceName      dot.NullString                     `json:"shipping_service_name"`
	ShippingServiceFee       dot.NullInt                        `json:"shipping_service_fee"`
	ActualShippingServiceFee dot.NullInt                        `json:"actual_shipping_service_fee"`
	ShippingServiceCode      dot.NullString                     `json:"shipping_service_code"`
	ShippingCode             dot.NullString                     `json:"shipping_code"`
	ShippingNote             dot.NullString                     `json:"shipping_note"`
	TryOn                    try_on.TryOnCode                   `json:"try_on"`
	IncludeInsurance         dot.NullBool                       `json:"include_insurance"`
	ConfirmStatus            status3.NullStatus                 `json:"confirm_status"`
	ShippingState            shipping.NullState                 `json:"shipping_state"`
	ShippingStatus           status5.NullStatus                 `json:"shipping_status"`
	Status                   status5.NullStatus                 `json:"status"`
	CodAmount                dot.NullInt                        `json:"cod_amount"`
	ActualCodAmount          dot.NullInt                        `json:"actual_cod_amount"`
	ChargeableWeight         dot.NullInt                        `json:"chargeable_weight"`
	PickupAddress            *OrderAddress                      `json:"pickup_address"`
	ReturnAddress            *OrderAddress                      `json:"return_address"`
	ShippingAddress          *OrderAddress                      `json:"shipping_address"`
	EtopPaymentStatus        status4.NullStatus                 `json:"etop_payment_status"`
	EstimatedDeliveryAt      dot.Time                           `json:"estimated_delivery_at"`
	EstimatedPickupAt        dot.Time                           `json:"estimated_pickup_at"`
}

func (m *Fulfillment) Reset()         { *m = Fulfillment{} }
func (m *Fulfillment) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentsResponse struct {
	Fulfillments []*Fulfillment         `json:"fulfillments"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (m *FulfillmentsResponse) Reset()         { *m = FulfillmentsResponse{} }
func (m *FulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

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
	ConnectionIDs   []dot.ID         `json:"connection_ids"`
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
	Height           int          `json:"height"`
	BasketValue      int          `json:"basket_value"`
	CodAmount        int          `json:"cod_amount"`
	IncludeInsurance dot.NullBool `json:"include_insurance"`
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
	Code string `json:"code"`
	// @deprecated use carrier info instead
	Name string `json:"name"`
	Fee  int    `json:"fee"`
	// @deprecated
	Carrier             shipping_provider.ShippingProvider `json:"carrier"`
	EstimatedPickupAt   dot.Time                           `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time                           `json:"estimated_delivery_at"`
	CarrierInfo         *CarrierInfo                       `json:"carrier_info"`
}

func (m *ShippingService) Reset()         { *m = ShippingService{} }
func (m *ShippingService) String() string { return jsonx.MustMarshalToString(m) }

type CarrierInfo struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func (m *CarrierInfo) Reset()         { *m = CarrierInfo{} }
func (m *CarrierInfo) String() string { return jsonx.MustMarshalToString(m) }

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

func (m *Order) HasChanged() bool {
	return m.Status.Valid ||
		m.ConfirmStatus.Valid ||
		m.FulfillmentShippingStatus.Valid ||
		m.EtopPaymentStatus.Valid ||
		m.BasketValue.Valid ||
		m.TotalAmount.Valid ||
		m.Shipping != nil ||
		m.CustomerAddress != nil || m.ShippingAddress != nil
}

func (m *Fulfillment) HasChanged() bool {
	return m.Status.Valid ||
		m.ShippingState.Valid ||
		m.EtopPaymentStatus.Valid ||
		m.ActualShippingServiceFee.Valid ||
		m.CodAmount.Valid ||
		m.ActualCodAmount.Valid ||
		m.ShippingNote.Valid ||
		m.ChargeableWeight.Valid
}

type Customer struct {
	Id           dot.ID         `json:"id"`
	ShopId       dot.ID         `json:"shop_id"`
	ExternalId   string         `json:"external_id"`
	ExternalCode string         `json:"external_code"`
	FullName     string         `json:"full_name"`
	Code         string         `json:"code"`
	Note         string         `json:"note"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Gender       string         `json:"gender"`
	Type         string         `json:"type"`
	Birthday     string         `json:"birthday"`
	CreatedAt    dot.Time       `json:"created_at"`
	UpdatedAt    dot.Time       `json:"updated_at"`
	Status       status3.Status `json:"status"`
	GroupIds     []dot.ID       `json:"group_ids"`
}

func (m *Customer) Reset()         { *m = Customer{} }
func (m *Customer) String() string { return jsonx.MustMarshalToString(m) }

type CustomersResponse struct {
	Customers []*Customer            `json:"customers"`
	Paging    *common.CursorPageInfo `json:"paging"`
}

func (m *CustomersResponse) Reset()         { *m = CustomersResponse{} }
func (m *CustomersResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	// @required
	FullName string        `json:"full_name"`
	Gender   gender.Gender `json:"gender"`
	Birthday string        `json:"birthday"`
	// enum ('independent', 'individual', 'organization')
	Type customer_type.CustomerType `json:"type"`
	Note string                     `json:"note"`
	// @required
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (m *CreateCustomerRequest) Reset()         { *m = CreateCustomerRequest{} }
func (m *CreateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerRequest struct {
	Id       dot.ID            `json:"id"`
	FullName dot.NullString    `json:"full_name"`
	Gender   gender.NullGender `json:"gender"`
	Birthday dot.NullString    `json:"birthday"`
	// enum ('individual', 'organization','independent')
	Type  dot.NullString `json:"type"`
	Note  dot.NullString `json:"note"`
	Phone dot.NullString `json:"phone"`
	Email dot.NullString `json:"email"`
}

func (m *UpdateCustomerRequest) Reset()         { *m = UpdateCustomerRequest{} }
func (m *UpdateCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteCustomerRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *DeleteCustomerRequest) Reset()         { *m = DeleteCustomerRequest{} }
func (m *DeleteCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetCustomerRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetCustomerRequest) Reset()         { *m = GetCustomerRequest{} }
func (m *GetCustomerRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *CustomerFilter) Reset()         { *m = CustomerFilter{} }
func (m *CustomerFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomersRequest struct {
	Filter CustomerFilter       `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCustomersRequest) Reset()         { *m = ListCustomersRequest{} }
func (m *ListCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddress struct {
	Id           dot.ID            `json:"id"`
	Province     string            `json:"province"`
	ProvinceCode string            `json:"province_code"`
	District     string            `json:"district"`
	DistrictCode string            `json:"district_code"`
	Ward         string            `json:"ward"`
	WardCode     string            `json:"ward_code"`
	Address1     string            `json:"address1"`
	Address2     string            `json:"address2"`
	FullName     string            `json:"full_name"`
	Company      string            `json:"company"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Position     string            `json:"position"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *CustomerAddress) Reset()         { *m = CustomerAddress{} }
func (m *CustomerAddress) String() string { return jsonx.MustMarshalToString(m) }

type CustomerAddressesResponse struct {
	CustomerAddresses []*CustomerAddress `json:"addresses"`
}

func (m *CustomerAddressesResponse) Reset()         { *m = CustomerAddressesResponse{} }
func (m *CustomerAddressesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomerAddressesRequest struct {
	CustomerId dot.ID `json:"customer_id"`
}

func (m *ListCustomerAddressesRequest) Reset()         { *m = ListCustomerAddressesRequest{} }
func (m *ListCustomerAddressesRequest) String() string { return jsonx.MustMarshalToString(m) }

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Coordinates) Reset()         { *m = Coordinates{} }
func (m *Coordinates) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerAddressRequest struct {
	CustomerId   dot.ID       `json:"customer_id"`
	ProvinceCode string       `json:"province_code"`
	DistrictCode string       `json:"district_code"`
	WardCode     string       `json:"ward_code"`
	Address1     string       `json:"address1"`
	Address2     string       `json:"address2"`
	Country      string       `json:"country"`
	FullName     string       `json:"full_name"`
	Company      string       `json:"company"`
	Phone        string       `json:"phone"`
	Email        string       `json:"email"`
	Position     string       `json:"position"`
	Coordinates  *Coordinates `json:"coordinates"`
}

func (m *CreateCustomerAddressRequest) Reset()         { *m = CreateCustomerAddressRequest{} }
func (m *CreateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerAddressRequest struct {
	Id           dot.ID            `json:"id"`
	ProvinceCode dot.NullString    `json:"province_code"`
	DistrictCode dot.NullString    `json:"district_code"`
	WardCode     dot.NullString    `json:"ward_code"`
	Address1     dot.NullString    `json:"address1"`
	Address2     dot.NullString    `json:"address2"`
	Country      dot.NullString    `json:"country"`
	FullName     dot.NullString    `json:"full_name"`
	Phone        dot.NullString    `json:"phone"`
	Email        dot.NullString    `json:"email"`
	Position     dot.NullString    `json:"position"`
	Company      dot.NullString    `json:"company"`
	Coordinates  *etop.Coordinates `json:"coordinates"`
}

func (m *UpdateCustomerAddressRequest) Reset()         { *m = UpdateCustomerAddressRequest{} }
func (m *UpdateCustomerAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddCustomersRequest struct {
	GroupId     dot.ID   `json:"id"`
	CustomerIds []dot.ID `json:"customer_ids"`
}

func (m *AddCustomersRequest) Reset()         { *m = AddCustomersRequest{} }
func (m *AddCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveCustomersRequest struct {
	GroupId     dot.ID   `json:"id"`
	CustomerIds []dot.ID `json:"customer_ids"`
}

func (m *RemoveCustomersRequest) Reset()         { *m = RemoveCustomersRequest{} }
func (m *RemoveCustomersRequest) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroup struct {
	Id   dot.ID `json:"id"`
	Name string `json:"name"`
}

func (m *CustomerGroup) Reset()         { *m = CustomerGroup{} }
func (m *CustomerGroup) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupsResponse struct {
	CustomerGroups []*CustomerGroup       `json:"customer_groups"`
	Paging         *common.CursorPageInfo `json:"paging"`
}

func (m *CustomerGroupsResponse) Reset()         { *m = CustomerGroupsResponse{} }
func (m *CustomerGroupsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CustomerGroupFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *CustomerGroupFilter) Reset()         { *m = CustomerGroupFilter{} }
func (m *CustomerGroupFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCustomerGroupsRequest struct {
	Filter CustomerGroupFilter  `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListCustomerGroupsRequest) Reset()         { *m = ListCustomerGroupsRequest{} }
func (m *ListCustomerGroupsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateCustomerGroupRequest struct {
	Name string `json:"name"`
}

func (m *CreateCustomerGroupRequest) Reset()         { *m = CreateCustomerGroupRequest{} }
func (m *CreateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCustomerGroupRequest struct {
	GroupId dot.ID         `json:"id"`
	Name    dot.NullString `json:"name"`
}

func (m *UpdateCustomerGroupRequest) Reset()         { *m = UpdateCustomerGroupRequest{} }
func (m *UpdateCustomerGroupRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListInventoryLevelsRequest struct {
	VariantIds []dot.ID `json:"variant_ids"`
}

func (m *ListInventoryLevelsRequest) Reset()         { *m = ListInventoryLevelsRequest{} }
func (m *ListInventoryLevelsRequest) String() string { return jsonx.MustMarshalToString(m) }

type InventoryLevel struct {
	VariantId         dot.ID   `json:"variant_id"`
	AvailableQuantity int      `json:"available_quantity"`
	ReservedQuantity  int      `json:"reserved_quantity"`
	PickedQuantity    int      `json:"picked_quantity"`
	UpdatedAt         dot.Time `json:"updated_at"`
}

func (m *InventoryLevel) Reset()         { *m = InventoryLevel{} }
func (m *InventoryLevel) String() string { return jsonx.MustMarshalToString(m) }

type InventoryLevelsResponse struct {
	InventoryLevels []*InventoryLevel `json:"inventory_levels"`
}

func (m *InventoryLevelsResponse) Reset()         { *m = InventoryLevelsResponse{} }
func (m *InventoryLevelsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmOrderRequest struct {
	OrderId dot.ID `json:"order_id"`
	// enum ('create', 'confirm')
	AutoInventoryVoucher dot.NullString `json:"auto_inventory_voucher"`
	// enum ('obey', 'ignore')
	InventoryPolicy bool `json:"inventory_policy"`
}

func (m *ConfirmOrderRequest) Reset()         { *m = ConfirmOrderRequest{} }
func (m *ConfirmOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type EtopProduct struct {
	Id          dot.ID   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	DescHtml    string   `json:"desc_html"`
	Unit        string   `json:"unit"`
	ImageUrls   []string `json:"image_urls"`
	ListPrice   int      `json:"list_price"`
	CostPrice   int      `json:"cost_price"`
	CategoryId  dot.ID   `json:"category_id"`
}

type Tag struct {
	Id    dot.ID `json:"id"`
	Label string `json:"label"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return jsonx.MustMarshalToString(m) }

type EtopVariant struct {
	Id          dot.ID                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	ShortDesc   string                    `json:"short_desc"`
	DescHtml    string                    `json:"desc_html"`
	ImageUrls   []string                  `json:"image_urls"`
	ListPrice   int                       `json:"list_price"`
	CostPrice   int                       `json:"cost_price"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
}

func (m *EtopVariant) Reset()         { *m = EtopVariant{} }
func (m *EtopVariant) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariant struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	// @required
	Id   dot.ID       `json:"id"`
	Info *EtopVariant `json:"info"`

	Code string `json:"code"`

	Name        string         `json:"name"`
	Description string         `json:"description"`
	ShortDesc   string         `json:"short_desc"`
	DescHtml    string         `json:"desc_html"`
	ImageUrls   []string       `json:"image_urls"`
	ListPrice   int            `json:"list_price"`
	RetailPrice int            `json:"retail_price"`
	Note        string         `json:"note"`
	Status      status3.Status `json:"status"`

	CostPrice int `json:"cost_price"`

	Tags []string `json:"tags"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`
}

func (m *ShopVariant) Reset()         { *m = ShopVariant{} }
func (m *ShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariantsResponse struct {
	ShopVariants []*ShopVariant         `json:"shop_variants"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (m *ShopVariantsResponse) Reset()         { *m = ShopVariantsResponse{} }
func (m *ShopVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShopProduct struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	// @required
	Id dot.ID `json:"id"`

	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ShortDesc     string         `json:"short_desc"`
	DescHtml      string         `json:"desc_html"`
	ImageUrls     []string       `json:"image_urls"`
	CategoryId    dot.ID         `json:"category_id"`
	Tags          []string       `json:"tags"`
	Note          string         `json:"note"`
	Status        status3.Status `json:"status"`
	ListPrice     int            `json:"list_price"`
	RetailPrice   int            `json:"retail_price"`
	CollectionIds []dot.ID       `json:"collection_ids"`
	Variants      []*ShopVariant `json:"variants"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
	BrandId   dot.ID   `json:"brand_id"`
}

func (m *ShopProduct) Reset()         { *m = ShopProduct{} }
func (m *ShopProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopProductsResponse struct {
	Products []*ShopProduct         `json:"products"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *ShopProductsResponse) Reset()         { *m = ShopProductsResponse{} }
func (m *ShopProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProductRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetProductRequest) Reset()         { *m = GetProductRequest{} }
func (m *GetProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *ProductFilter) Reset()         { *m = ProductFilter{} }
func (m *ProductFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListProductsRequest struct {
	Filter ProductFilter        `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListProductsRequest) Reset()         { *m = ListProductsRequest{} }
func (m *ListProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateProductRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Unit        string   `json:"unit"`
	Note        string   `json:"note"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	DescHtml    string   `json:"desc_html"`
	ImageUrls   []string `json:"image_urls"`
	CostPrice   int      `json:"cost_price"`
	ListPrice   int      `json:"list_price"`
	RetailPrice int      `json:"retail_price"`
	BrandId     dot.ID   `json:"brand_id"`
}

func (m *CreateProductRequest) Reset()         { *m = CreateProductRequest{} }
func (m *CreateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductRequest struct {
	// @required
	Id          dot.ID         `json:"id"`
	Name        dot.NullString `json:"name"`
	Note        dot.NullString `json:"note"`
	Unit        dot.NullString `json:"unit"`
	Description dot.NullString `json:"description"`
	ShortDesc   dot.NullString `json:"short_desc"`
	DescHtml    dot.NullString `json:"desc_html"`
	CostPrice   dot.NullInt    `json:"cost_price"`
	ListPrice   dot.NullInt    `json:"list_price"`
	RetailPrice dot.NullInt    `json:"retail_price"`
	BrandId     dot.NullID     `json:"brand_id"`
}

func (m *UpdateProductRequest) Reset()         { *m = UpdateProductRequest{} }
func (m *UpdateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type BatchUpdateProductsRequest struct {
	ProductIds   []dot.ID `json:"product_ids"`
	CollectionId dot.ID   `json:"collection_id"`

	// tags
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`

	Status status3.Status `json:"status"`
}

func (m *BatchUpdateProductsRequest) Reset()         { *m = BatchUpdateProductsRequest{} }
func (m *BatchUpdateProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveProductsRequest struct {
	Ids []dot.ID `json:"ids"`
}

func (m *RemoveProductsRequest) Reset()         { *m = RemoveProductsRequest{} }
func (m *RemoveProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *AddProductCollectionRequest) Reset()         { *m = AddProductCollectionRequest{} }
func (m *AddProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveProductCollectionRequest struct {
	ProductId     dot.ID   `json:"product_id"`
	CollectionIds []dot.ID `json:"collection_ids"`
}

func (m *RemoveProductCollectionRequest) Reset()         { *m = RemoveProductCollectionRequest{} }
func (m *RemoveProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetVariantRequest) Reset()         { *m = GetVariantRequest{} }
func (m *GetVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type VariantFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *VariantFilter) Reset()         { *m = VariantFilter{} }
func (m *VariantFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListVariantsRequest struct {
	Filter VariantFilter        `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListVariantsRequest) Reset()         { *m = ListVariantsRequest{} }
func (m *ListVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateVariantRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	ProductId   dot.ID                    `json:"product_id"`
	Note        string                    `json:"note"`
	Description string                    `json:"description"`
	ShortDesc   string                    `json:"short_desc"`
	DescHtml    string                    `json:"desc_html"`
	ImageUrls   []string                  `json:"image_urls"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`
	CostPrice   int                       `json:"cost_price"`
	ListPrice   int                       `json:"list_price"`
	RetailPrice int                       `json:"retail_price"`
}

func (m *CreateVariantRequest) Reset()         { *m = CreateVariantRequest{} }
func (m *CreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantRequest struct {
	// @required
	Id          dot.ID                    `json:"id"`
	Name        dot.NullString            `json:"name"`
	Note        dot.NullString            `json:"note"`
	Code        dot.NullString            `json:"code"`
	CostPrice   dot.NullInt               `json:"cost_price"`
	ListPrice   dot.NullInt               `json:"list_price"`
	RetailPrice dot.NullInt               `json:"retail_price"`
	Description dot.NullString            `json:"description"`
	ShortDesc   dot.NullString            `json:"short_desc"`
	DescHtml    dot.NullString            `json:"desc_html"`
	Attributes  []*catalogtypes.Attribute `json:"attributes"`

	// images
	Adds       []string `json:"adds"`
	Deletes    []string `json:"deletes"`
	ReplaceAll []string `json:"replace_all"`
	DeleteAll  bool     `json:"delete_all"`
}

func (m *UpdateVariantRequest) Reset()         { *m = UpdateVariantRequest{} }
func (m *UpdateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type BatchUpdateVariantsRequest struct {
	Ids    []dot.ID       `json:"ids"`
	Status status3.Status `json:"status"`
}

func (m *BatchUpdateVariantsRequest) Reset()         { *m = BatchUpdateVariantsRequest{} }
func (m *BatchUpdateVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }
