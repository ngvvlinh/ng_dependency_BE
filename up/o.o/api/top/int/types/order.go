package types

import (
	catalogtypes "o.o/api/main/catalog/types"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/api/top/int/etop"
	"o.o/api/top/int/types/spreadsheet"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/fee"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/ghn_note_code"
	source "o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func (m *OrderWithErrorsResponse) HasErrors() []*common.Error {
	return m.FulfillmentErrors
}

func (m *ImportOrdersResponse) HasErrors() []*common.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}

type OrdersResponse struct {
	Paging *common.PageInfo `json:"paging"`
	Orders []*Order         `json:"orders"`
}

func (m *OrdersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Order struct {
	ExportedFields []string `json:"exported_fields"`
	Id             dot.ID   `json:"id"`
	ShopId         dot.ID   `json:"shop_id"`
	ShopName       string   `json:"shop_name"`
	Code           string   `json:"code"`
	// the same as external_code
	EdCode          string                       `json:"ed_code"`
	ExternalCode    string                       `json:"external_code"`
	Source          source.Source                `json:"source"`
	PartnerId       dot.ID                       `json:"partner_id"`
	ExternalId      string                       `json:"external_id"`
	ExternalUrl     string                       `json:"external_url"`
	SelfUrl         string                       `json:"self_url"`
	PaymentMethod   payment_method.PaymentMethod `json:"payment_method"`
	Customer        *OrderCustomer               `json:"customer"`
	CustomerAddress *OrderAddress                `json:"customer_address"`
	BillingAddress  *OrderAddress                `json:"billing_address"`
	ShippingAddress *OrderAddress                `json:"shipping_address"`
	CreatedAt       dot.Time                     `json:"created_at"`
	ProcessedAt     dot.Time                     `json:"processed_at"`
	UpdatedAt       dot.Time                     `json:"updated_at"`
	ClosedAt        dot.Time                     `json:"closed_at"`
	ConfirmedAt     dot.Time                     `json:"confirmed_at"`
	CancelledAt     dot.Time                     `json:"cancelled_at"`
	CancelReason    string                       `json:"cancel_reason"`
	ShConfirm       status3.Status               `json:"sh_confirm"`
	// @deprecated replaced by confirm_status
	Confirm       status3.Status `json:"confirm"`
	ConfirmStatus status3.Status `json:"confirm_status"`
	Status        status5.Status `json:"status"`
	// @deprecated replaced by fulfillment_shipping_status
	FulfillmentStatus         status5.Status   `json:"fulfillment_status"`
	FulfillmentShippingStatus status5.Status   `json:"fulfillment_shipping_status"`
	CustomerPaymentStatus     status3.Status   `json:"customer_payment_status"`
	EtopPaymentStatus         status4.Status   `json:"etop_payment_status"`
	Lines                     []*OrderLine     `json:"lines"`
	Discounts                 []*OrderDiscount `json:"discounts"`
	TotalItems                int              `json:"total_items"`
	BasketValue               int              `json:"basket_value"`
	TotalWeight               int              `json:"total_weight"`
	OrderDiscount             int              `json:"order_discount"`
	TotalDiscount             int              `json:"total_discount"`
	TotalAmount               int              `json:"total_amount"`
	OrderNote                 string           `json:"order_note"`
	// @deprecated use fee_lines.shipping
	ShippingFee int             `json:"shipping_fee"`
	TotalFee    int             `json:"total_fee"`
	FeeLines    []*OrderFeeLine `json:"fee_lines"`
	// @deperecated use fee_lines.shipping instead
	ShopShippingFee int `json:"shop_shipping_fee"`
	// @deprecated use shop_shipping.shipping_note instead
	ShippingNote string `json:"shipping_note"`
	// @deperecated use shop_shipping.cod_amount instead
	ShopCod      int             `json:"shop_cod"`
	ReferenceUrl string          `json:"reference_url"`
	Fulfillments []*XFulfillment `json:"fulfillments"`
	// @deprecated use shipping instead
	ShopShipping *OrderShipping `json:"shop_shipping"`
	Shipping     *OrderShipping `json:"shipping"`
	// @deprecated use try_on_code instead
	GhnNoteCode     ghn_note_code.GHNNoteCode `json:"ghn_note_code,omitempty"`
	FulfillmentType string                    `json:"fulfillment_type"`
	FulfillmentIds  []dot.ID                  `json:"fulfillment_ids"`
	// received_amount get from receipt
	ReceivedAmount int            `json:"received_amount"`
	CustomerId     dot.ID         `json:"customer_id"`
	PaymentStatus  status4.Status `json:"payment_status"`
	CreatedBy      dot.ID         `json:"created_by"`
	PreOrder       bool           `json:"pre_order"`
}

func (m *Order) String() string { return jsonx.MustMarshalToString(m) }

type OrderLineMetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

func (m *OrderLineMetaField) String() string { return jsonx.MustMarshalToString(m) }

type OrderLine struct {
	ExportedFields []string                  `json:"exported_fields"`
	OrderId        dot.ID                    `json:"order_id"`
	VariantId      dot.ID                    `json:"variant_id"`
	ProductName    string                    `json:"product_name"`
	IsOutsideEtop  bool                      `json:"is_outside_etop"`
	Quantity       int                       `json:"quantity"`
	ListPrice      int                       `json:"list_price"`
	RetailPrice    int                       `json:"retail_price"`
	PaymentPrice   int                       `json:"payment_price"`
	ImageUrl       string                    `json:"image_url"`
	Attributes     []*catalogtypes.Attribute `json:"attributes"`
	ProductId      dot.ID                    `json:"product_id"`
	TotalDiscount  int                       `json:"total_discount"`
	MetaFields     []*OrderLineMetaField     `json:"meta_fields"`
	Code           string                    `json:"code"`
}

func (m *OrderLine) String() string { return jsonx.MustMarshalToString(m) }

type OrderFeeLine struct {
	Type fee.FeeType `json:"type"`
	// @required
	Name string `json:"name"`
	Code string `json:"code"`
	Desc string `json:"desc"`
	// @required
	Amount Int `json:"amount"`
}

func (m *OrderFeeLine) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrderLine struct {
	VariantId    dot.ID                    `json:"variant_id"`
	ProductName  string                    `json:"product_name"`
	Quantity     int                       `json:"quantity"`
	ListPrice    int                       `json:"list_price"`
	RetailPrice  int                       `json:"retail_price"`
	PaymentPrice int                       `json:"payment_price"`
	ImageUrl     string                    `json:"image_url"`
	Attributes   []*catalogtypes.Attribute `json:"attributes"`
	MetaFields   []*OrderLineMetaField     `json:"meta_fields"`
}

func (m *CreateOrderLine) String() string { return jsonx.MustMarshalToString(m) }

type OrderCustomer struct {
	ExportedFields []string                   `json:"exported_fields"`
	FirstName      string                     `json:"first_name"`
	LastName       string                     `json:"last_name"`
	FullName       string                     `json:"full_name"`
	Email          string                     `json:"email"`
	Phone          string                     `json:"phone"`
	Gender         gender.Gender              `json:"gender"`
	Type           customer_type.CustomerType `json:"type"`
	Deleted        bool                       `json:"deleted"`
}

func (m *OrderCustomer) String() string { return jsonx.MustMarshalToString(m) }

type OrderAddress struct {
	ExportedFields []string          `json:"exported_fields"`
	FullName       string            `json:"full_name"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	Phone          string            `json:"phone"`
	Email          string            `json:"email"`
	Country        string            `json:"country"`
	City           string            `json:"city"`
	Province       string            `json:"province"`
	District       string            `json:"district"`
	Ward           string            `json:"ward"`
	Zip            string            `json:"zip"`
	Company        string            `json:"company"`
	Address1       string            `json:"address1"`
	Address2       string            `json:"address2"`
	ProvinceCode   string            `json:"province_code"`
	DistrictCode   string            `json:"district_code"`
	WardCode       string            `json:"ward_code"`
	Coordinates    *etop.Coordinates `json:"coordinates"`
}

func (m *OrderAddress) String() string { return jsonx.MustMarshalToString(m) }

type OrderDiscount struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

func (m *OrderDiscount) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrderRequest struct {
	Source        source.Source                `json:"source"`
	ExternalId    string                       `json:"external_id"`
	ExternalCode  string                       `json:"external_code"`
	ExternalUrl   string                       `json:"external_url"`
	PaymentMethod payment_method.PaymentMethod `json:"payment_method"`
	// If order_source is self, customer must be shop information
	// and customer_address, shipping_address must be shop address.
	Customer        *OrderCustomer `json:"customer"`
	CustomerAddress *OrderAddress  `json:"customer_address"`
	BillingAddress  *OrderAddress  `json:"billing_address"`
	ShippingAddress *OrderAddress  `json:"shipping_address"`
	// If there are products from shop, this field should be set.
	// Otherwise, use shop default address.
	ShopAddress *OrderAddress      `json:"shop_address"`
	ShConfirm   status3.NullStatus `json:"sh_confirm"`
	Lines       []*CreateOrderLine `json:"lines"`
	Discounts   []*OrderDiscount   `json:"discounts"`
	TotalItems  int                `json:"total_items"`
	BasketValue int                `json:"basket_value"`
	// @deprecated use shipping.gross_weight, shipping.chargeable_weight
	TotalWeight   int             `json:"total_weight"`
	OrderDiscount int             `json:"order_discount"`
	TotalFee      int             `json:"total_fee"`
	FeeLines      []*OrderFeeLine `json:"fee_lines"`
	TotalDiscount dot.NullInt     `json:"total_discount"`
	TotalAmount   int             `json:"total_amount"`
	OrderNote     string          `json:"order_note"`
	ShippingNote  string          `json:"shipping_note"`
	// @deprecated use fee_lines.shipping instead
	ShopShippingFee int `json:"shop_shipping_fee"`
	// @deprecated use shipping.cod_amount instead
	ShopCod      int    `json:"shop_cod"`
	ReferenceUrl string `json:"reference_url"`
	// @deprecated use shipping instead
	ShopShipping *OrderShipping `json:"shop_shipping"`
	Shipping     *OrderShipping `json:"shipping"`
	// @deprecated use shop_shipping.try_on instead
	GhnNoteCode  ghn_note_code.GHNNoteCode `json:"ghn_note_code"`
	ExternalMeta map[string]string         `json:"external_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ReferralMeta map[string]string         `json:"referral_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CustomerId   dot.ID                    `json:"customer_id"`
	PreOrder     bool                      `json:"pre_order"`
}

func (m *CreateOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateOrderRequest struct {
	// @required
	Id              dot.ID         `json:"id"`
	Customer        *OrderCustomer `json:"customer"`
	CustomerAddress *OrderAddress  `json:"customer_address"`
	BillingAddress  *OrderAddress  `json:"billing_address"`
	ShippingAddress *OrderAddress  `json:"shipping_address"`
	ShopAddress     *OrderAddress  `json:"shop_address"`
	OrderNote       string         `json:"order_note"`
	ShippingNote    string         `json:"shipping_note"`
	// @deprecated use fee_lines instead
	ShopShippingFee dot.NullInt `json:"shop_shipping_fee"`
	// @deprecated use shipping.cod_amount instead
	ShopCod dot.NullInt `json:"shop_cod"`
	// @deprecated use shipping instead
	ShopShipping  *OrderShipping  `json:"shop_shipping"`
	Shipping      *OrderShipping  `json:"shipping"`
	FeeLines      []*OrderFeeLine `json:"fee_lines"`
	OrderDiscount dot.NullInt     `json:"order_discount"`
	// @deprecated
	TotalWeight      int                `json:"total_weight"`
	ChargeableWeight int                `json:"chargeable_weight"`
	Lines            []*CreateOrderLine `json:"lines"`
	BasketValue      int                `json:"basket_value"`
	TotalAmount      int                `json:"total_amount"`
	TotalItems       int                `json:"total_items"`
	TotalFee         dot.NullInt        `json:"total_fee"`
	CustomerId       dot.ID             `json:"customer_id"`
}

func (m *UpdateOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

// OrderShipping provides shipping information for an order to be shipped.
type OrderShipping struct {
	ExportedFields []string `json:"exported_fields"`
	// @deprecated use pickup_address
	ShAddress *OrderAddress `json:"sh_address"`
	// @deprecated
	XServiceId string `json:"x_service_id"`
	// @deprecated
	XShippingFee int `json:"x_shipping_fee"`
	// @deprecated
	XServiceName        string        `json:"x_service_name"`
	PickupAddress       *OrderAddress `json:"pickup_address"`
	ReturnAddress       *OrderAddress `json:"return_address"`
	ShippingServiceName string        `json:"shipping_service_name"`
	ShippingServiceCode string        `json:"shipping_service_code"`
	ShippingServiceFee  int           `json:"shipping_service_fee"`
	// @deprecated use carrier
	ShippingProvider shipping_provider.ShippingProvider `json:"shipping_provider"`
	Carrier          shipping_provider.ShippingProvider `json:"carrier"`
	IncludeInsurance bool                               `json:"include_insurance"`
	TryOn            try_on.TryOnCode                   `json:"try_on"`
	ShippingNote     string                             `json:"shipping_note"`
	CodAmount        dot.NullInt                        `json:"cod_amount"`
	// @deprecated
	Weight           dot.NullInt `json:"weight"`
	GrossWeight      dot.NullInt `json:"gross_weight"`
	Length           dot.NullInt `json:"length"`
	Width            dot.NullInt `json:"width"`
	Height           dot.NullInt `json:"height"`
	ChargeableWeight dot.NullInt `json:"chargeable_weight"`
}

func (m *OrderShipping) String() string { return jsonx.MustMarshalToString(m) }

type Fulfillment struct {
	ExportedFields []string     `json:"exported_fields"`
	Id             dot.ID       `json:"id"`
	OrderId        dot.ID       `json:"order_id"`
	ShopId         dot.ID       `json:"shop_id"`
	PartnerId      dot.ID       `json:"partner_id"`
	SelfUrl        string       `json:"self_url"`
	Lines          []*OrderLine `json:"lines"`
	TotalItems     int          `json:"total_items"`
	// @deprecated use chargeable_weight
	TotalWeight int `json:"total_weight"`
	BasketValue int `json:"basket_value"`
	// @deprecated use cod_amount
	TotalCodAmount int `json:"total_cod_amount"`
	CodAmount      int `json:"cod_amount"`
	// @deprecated
	TotalAmount      int      `json:"total_amount"`
	ChargeableWeight int      `json:"chargeable_weight"`
	CreatedAt        dot.Time `json:"created_at"`
	UpdatedAt        dot.Time `json:"updated_at"`
	ClosedAt         dot.Time `json:"closed_at"`
	CancelledAt      dot.Time `json:"cancelled_at"`
	CancelReason     string   `json:"cancel_reason"`
	// deprecated: use carrier instead
	ShippingProvider     shipping_provider.ShippingProvider `json:"shipping_provider"`
	Carrier              shipping_provider.ShippingProvider `json:"carrier"`
	ShippingServiceName  string                             `json:"shipping_service_name"`
	ShippingServiceFee   int                                `json:"shipping_service_fee"`
	ShippingServiceCode  string                             `json:"shipping_service_code"`
	ShippingCode         string                             `json:"shipping_code"`
	ShippingNote         string                             `json:"shipping_note"`
	TryOn                try_on.TryOnCode                   `json:"try_on"`
	IncludeInsurance     bool                               `json:"include_insurance"`
	ShConfirm            status3.Status                     `json:"sh_confirm"`
	ShippingState        shipping.State                     `json:"shipping_state"`
	Status               status5.Status                     `json:"status"`
	ShippingStatus       status5.Status                     `json:"shipping_status"`
	EtopPaymentStatus    status4.Status                     `json:"etop_payment_status"`
	ShippingFeeCustomer  int                                `json:"shipping_fee_customer"`
	ShippingFeeShop      int                                `json:"shipping_fee_shop"`
	XShippingFee         int                                `json:"x_shipping_fee"`
	XShippingId          string                             `json:"x_shipping_id"`
	XShippingCode        string                             `json:"x_shipping_code"`
	XShippingCreatedAt   dot.Time                           `json:"x_shipping_created_at"`
	XShippingUpdatedAt   dot.Time                           `json:"x_shipping_updated_at"`
	XShippingCancelledAt dot.Time                           `json:"x_shipping_cancelled_at"`
	XShippingDeliveredAt dot.Time                           `json:"x_shipping_delivered_at"`
	XShippingReturnedAt  dot.Time                           `json:"x_shipping_returned_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt dot.Time `json:"expected_delivery_at"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt              dot.Time               `json:"expected_pick_at"`
	EstimatedDeliveryAt         dot.Time               `json:"estimated_delivery_at"`
	EstimatedPickupAt           dot.Time               `json:"estimated_pickup_at"`
	CodEtopTransferedAt         dot.Time               `json:"cod_etop_transfered_at"`
	ShippingFeeShopTransferedAt dot.Time               `json:"shipping_fee_shop_transfered_at"`
	XShippingState              string                 `json:"x_shipping_state"`
	XShippingStatus             status5.Status         `json:"x_shipping_status"`
	XSyncStatus                 status4.Status         `json:"x_sync_status"`
	XSyncStates                 *FulfillmentSyncStates `json:"x_sync_states"`
	// @deprecated use shipping_address instead
	AddressTo *etop.Address `json:"address_to"`
	// @deprecated use pickup_address instead
	AddressFrom                        *etop.Address          `json:"address_from"`
	PickupAddress                      *OrderAddress          `json:"pickup_address"`
	ReturnAddress                      *OrderAddress          `json:"return_address"`
	ShippingAddress                    *OrderAddress          `json:"shipping_address"`
	Shop                               *etop.Shop             `json:"shop"`
	Order                              *Order                 `json:"order"`
	ProviderShippingFeeLines           []*ShippingFeeLine     `json:"provider_shipping_fee_lines"`
	ShippingFeeShopLines               []*ShippingFeeLine     `json:"shipping_fee_shop_lines"`
	EtopDiscount                       int                    `json:"etop_discount"`
	MoneyTransactionShippingId         dot.ID                 `json:"money_transaction_shipping_id"`
	MoneyTransactionShippingExternalId dot.ID                 `json:"money_transaction_shipping_external_id"`
	XShippingLogs                      []*ExternalShippingLog `json:"x_shipping_logs"`
	XShippingNote                      string                 `json:"x_shipping_note"`
	XShippingSubState                  string                 `json:"x_shipping_sub_state"`
	Code                               string                 `json:"code"`
	ActualCompensationAmount           int                    `json:"actual_compensation_amount"`
	AdminNote                          string                 `json:"admin_note"`

	ConnectionID  dot.ID `json:"connection_id"`
	ShopCarrierID dot.ID `json:"shop_carrier_id"`
}

func (m *Fulfillment) String() string { return jsonx.MustMarshalToString(m) }

type ShippingFeeLine struct {
	ShippingFeeType          shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
	Cost                     int                               `json:"cost"`
	ExternalServiceId        string                            `json:"external_service_id"`
	ExternalServiceName      string                            `json:"external_service_name"`
	ExternalServiceType      string                            `json:"external_service_type"`
	ExternalShippingOrderId  string                            `json:"external_shipping_order_id"`
	ExternalPaymentChannelId string                            `json:"external_payment_channel_id"`
}

func (m *ShippingFeeLine) String() string { return jsonx.MustMarshalToString(m) }

type ShippingFeeShortLine struct {
	ShippingFeeType shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
	Cost            int                               `json:"cost"`
}

func (m *ShippingFeeShortLine) String() string { return jsonx.MustMarshalToString(m) }

type ExternalShippingLog struct {
	StateText string `json:"state_text"`
	Time      string `json:"time"`
	Message   string `json:"message"`
}

func (m *ExternalShippingLog) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentsResponse struct {
	Paging       *common.PageInfo `json:"paging"`
	Fulfillments []*Fulfillment   `json:"fulfillments"`
}

func (m *FulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type OrderWithErrorsResponse struct {
	// @deprecated
	Errors            []*common.Error `json:"errors"`
	Order             *Order          `json:"order"`
	FulfillmentErrors []*common.Error `json:"fulfillment_errors"`
}

func (m *OrderWithErrorsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetExternalShippingServicesRequest struct {
	// @deprecated use carrier instead
	Provider         shipping_provider.ShippingProvider `json:"provider"`
	Carrier          shipping_provider.ShippingProvider `json:"carrier"`
	FromDistrictCode string                             `json:"from_district_code"`
	FromProvinceCode string                             `json:"from_province_code"`
	ToDistrictCode   string                             `json:"to_district_code"`
	ToProvinceCode   string                             `json:"to_province_code"`
	FromProvince     string                             `json:"from_province"`
	FromDistrict     string                             `json:"from_district"`
	ToProvince       string                             `json:"to_province"`
	ToDistrict       string                             `json:"to_district"`
	// @deprecated use gross_weight instead
	Weight           int `json:"weight"`
	GrossWeight      int `json:"gross_weight"`
	ChargeableWeight int `json:"chargeable_weight"`
	Length           int `json:"length"`
	Width            int `json:"width"`
	Height           int `json:"height"`
	// @deprecated use basket_value instead
	Value int `json:"value"`
	// @deprecated use cod_amount instead
	TotalCodAmount   int          `json:"total_cod_amount"`
	CodAmount        int          `json:"cod_amount"`
	BasketValue      int          `json:"basket_value"`
	IncludeInsurance dot.NullBool `json:"include_insurance"`
}

func (m *GetExternalShippingServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetExternalShippingServicesResponse struct {
	Services []*ExternalShippingService `json:"services"`
}

func (m *GetExternalShippingServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ExternalShippingService struct {
	ExportedFields []string `json:"exported_fields"`
	// @deprecated use code
	ExternalId string `json:"external_id"`
	// @deprecated use fee
	ServiceFee int `json:"service_fee"`
	// @deprecated use carier
	Provider shipping_provider.ShippingProvider `json:"provider"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt dot.Time `json:"expected_pick_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt dot.Time `json:"expected_delivery_at"`
	Name               string   `json:"name"`
	Code               string   `json:"code"`
	Fee                int      `json:"fee"`
	// @deprecated use connection instead
	Carrier             shipping_provider.ShippingProvider `json:"carrier"`
	EstimatedPickupAt   dot.Time                           `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time                           `json:"estimated_delivery_at"`
	ConnectionInfo      *ConnectionInfo                    `json:"connection_info"`
	ShipmentServiceInfo *ShipmentServiceInfo               `json:"shipment_service_info"`
	ShipmentPriceInfo   *ShipmentPriceInfo                 `json:"shipment_price_info"`
}

func (m *ExternalShippingService) String() string { return jsonx.MustMarshalToString(m) }

type ConnectionInfo struct {
	ID       dot.ID `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func (m *ConnectionInfo) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentServiceInfo struct {
	ID           dot.ID `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	IsAvailable  bool   `json:"is_available"`
	ErrorMessage string `json:"error_message"`
}

func (m *ShipmentServiceInfo) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentPriceInfo struct {
	ID            dot.ID `json:"id"`
	OriginFee     int    `json:"origin_fee"`
	OriginMainFee int    `json:"origin_main_fee"`
	MakeupMainFee int    `json:"makeup_main_fee"`
}

func (m *ShipmentPriceInfo) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentSyncStates struct {
	SyncAt            dot.Time       `json:"sync_at"`
	NextShippingState shipping.State `json:"next_shipping_state"`
	Error             *common.Error  `json:"error"`
}

func (m *FulfillmentSyncStates) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransaction struct {
	Id                                 dot.ID            `json:"id"`
	ShopId                             dot.ID            `json:"shop_id"`
	Status                             status3.Status    `json:"status"`
	TotalCod                           int               `json:"total_cod"`
	TotalOrders                        int               `json:"total_orders"`
	Code                               string            `json:"code"`
	Provider                           string            `json:"provider"`
	MoneyTransactionShippingExternalId dot.ID            `json:"money_transaction_shipping_external_id"`
	MoneyTransactionShippingEtopId     dot.ID            `json:"money_transaction_shipping_etop_id"`
	TotalAmount                        int               `json:"total_amount"`
	CreatedAt                          dot.Time          `json:"created_at"`
	UpdatedAt                          dot.Time          `json:"updated_at"`
	ClosedAt                           dot.Time          `json:"closed_at"`
	ConfirmedAt                        dot.Time          `json:"confirmed_at"`
	EtopTransferedAt                   dot.Time          `json:"etop_transfered_at"`
	Note                               string            `json:"note"`
	InvoiceNumber                      string            `json:"invoice_number"`
	BankAccount                        *etop.BankAccount `json:"bank_account"`
}

func (m *MoneyTransaction) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransactionsResponse struct {
	Paging            *common.PageInfo    `json:"paging"`
	MoneyTransactions []*MoneyTransaction `json:"money_transactions"`
}

func (m *MoneyTransactionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransactionShippingExternalLine struct {
	Id                                 dot.ID        `json:"id"`
	ExternalCode                       string        `json:"external_code"`
	ExternalCustomer                   string        `json:"external_customer"`
	ExternalAddress                    string        `json:"external_address"`
	ExternalTotalCod                   int           `json:"external_total_cod"`
	ExternalTotalShippingFee           int           `json:"external_total_shipping_fee"`
	EtopFulfillmentId                  dot.ID        `json:"etop_fulfillment_id"`
	EtopFulfillmentIdRaw               string        `json:"etop_fulfillment_id_raw"`
	Note                               string        `json:"note"`
	MoneyTransactionShippingExternalId dot.ID        `json:"money_transaction_shipping_external_id"`
	ImportError                        *common.Error `json:"import_error"`
	CreatedAt                          dot.Time      `json:"created_at"`
	UpdatedAt                          dot.Time      `json:"updated_at"`
	ExternalCreatedAt                  dot.Time      `json:"external_created_at"`
	ExternalClosedAt                   dot.Time      `json:"external_closed_at"`
	Fulfillment                        *Fulfillment  `json:"fulfillment"`
}

func (m *MoneyTransactionShippingExternalLine) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransactionShippingExternal struct {
	Id             dot.ID                                  `json:"id"`
	Code           string                                  `json:"code"`
	TotalCod       int                                     `json:"total_cod"`
	TotalOrders    int                                     `json:"total_orders"`
	Status         status3.Status                          `json:"status"`
	Provider       string                                  `json:"provider"`
	Lines          []*MoneyTransactionShippingExternalLine `json:"lines"`
	CreatedAt      dot.Time                                `json:"created_at"`
	UpdatedAt      dot.Time                                `json:"updated_at"`
	ExternalPaidAt dot.Time                                `json:"external_paid_at"`
	Note           string                                  `json:"note"`
	InvoiceNumber  string                                  `json:"invoice_number"`
	BankAccount    *etop.BankAccount                       `json:"bank_account"`
}

func (m *MoneyTransactionShippingExternal) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransactionShippingExternalsResponse struct {
	Paging            *common.PageInfo                    `json:"paging"`
	MoneyTransactions []*MoneyTransactionShippingExternal `json:"money_transactions"`
}

func (m *MoneyTransactionShippingExternalsResponse) Reset() {
	*m = MoneyTransactionShippingExternalsResponse{}
}
func (m *MoneyTransactionShippingExternalsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type MoneyTransactionShippingEtop struct {
	Id                    dot.ID         `json:"id"`
	Code                  string         `json:"code"`
	TotalCod              int            `json:"total_cod"`
	TotalOrders           int            `json:"total_orders"`
	TotalAmount           int            `json:"total_amount"`
	TotalFee              int            `json:"total_fee"`
	TotalMoneyTxShippings int            `json:"total_money_tx_shippings"`
	Status                status3.Status `json:"status"`
	// @deprecated
	MoneyTransactions []*MoneyTransaction `json:"money_transactions"`
	CreatedAt         dot.Time            `json:"created_at"`
	UpdatedAt         dot.Time            `json:"updated_at"`
	ConfirmedAt       dot.Time            `json:"confirmed_at"`
	Note              string              `json:"note"`
	InvoiceNumber     string              `json:"invoice_number"`
	BankAccount       *etop.BankAccount   `json:"bank_account"`
}

func (m *MoneyTransactionShippingEtop) String() string { return jsonx.MustMarshalToString(m) }

type MoneyTransactionShippingEtopsResponse struct {
	Paging                        *common.PageInfo                `json:"paging"`
	MoneyTransactionShippingEtops []*MoneyTransactionShippingEtop `json:"money_transaction_shipping_etops"`
}

func (m *MoneyTransactionShippingEtopsResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportOrdersResponse struct {
	Data         *spreadsheet.SpreadsheetData `json:"data"`
	Orders       []*Order                     `json:"orders"`
	ImportErrors []*common.Error              `json:"import_errors"`
	CellErrors   []*common.Error              `json:"cell_errors"`
	ImportId     dot.ID                       `json:"import_id"`
}

func (m *ImportOrdersResponse) String() string { return jsonx.MustMarshalToString(m) }

// Public API for using with ManyChat
type PublicFulfillment struct {
	Id            dot.ID         `json:"id"`
	ShippingState shipping.State `json:"shipping_state"`
	Status        status5.Status `json:"status"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt  dot.Time `json:"expected_delivery_at"`
	DeliveredAt         dot.Time `json:"delivered_at"`
	EstimatedPickupAt   dot.Time `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time `json:"estimated_delivery_at"`
	ShippingCode        string   `json:"shipping_code"`
	OrderId             dot.ID   `json:"order_id"`
	// For using with ManyChat
	DeliveredAtText   string `json:"delivered_at_text"`
	ShippingStateText string `json:"shipping_state_text"`
}

func (m *PublicFulfillment) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowFulfillment struct {
	Id                         dot.ID           `json:"id"`
	ShopId                     dot.ID           `json:"shop_id"`
	PartnerId                  dot.ID           `json:"partner_id"`
	PickupAddress              *OrderAddress    `json:"pickup_address"`
	DeliveryPoints             []*DeliveryPoint `json:"delivery_points"`
	Carrier                    string           `json:"carrier"`
	ShippingServiceCode        string           `json:"shipping_service_code"`
	ShippingServiceFee         int              `json:"shipping_service_fee"`
	ShippingServiceName        string           `json:"shipping_service_name"`
	ShippingServiceDescription string           `json:"shipping_service_description"`
	WeightInfo                 `json:"weight_info"`
	ValueInfo                  ValueInfo           `json:"value_info"`
	ShippingNote               string              `json:"shipping_note"`
	RequestPickupAt            dot.Time            `json:"request_pickup_at"`
	CreatedAt                  dot.Time            `json:"created_at"`
	UpdatedAt                  dot.Time            `json:"updated_at"`
	Status                     status5.Status      `json:"status"`
	ShippingStatus             status5.Status      `json:"shipping_status"`
	ShippingState              shipnow_state.State `json:"shipping_state"`
	ConfirmStatus              status3.Status      `json:"confirm_status"`
	OrderIds                   []dot.ID            `json:"order_ids"`
	ShippingCreatedAt          dot.Time            `json:"shipping_created_at"`
	ShippingCode               string              `json:"shipping_code"`
	EtopPaymentStatus          status4.Status      `json:"etop_payment_status"`
	CodEtopTransferedAt        dot.Time            `json:"cod_etop_transfered_at"`
	ShippingPickingAt          dot.Time            `json:"shipping_picking_at"`
	ShippingDeliveringAt       dot.Time            `json:"shipping_delivering_at"`
	ShippingDeliveredAt        dot.Time            `json:"shipping_delivered_at"`
	ShippingCancelledAt        dot.Time            `json:"shipping_cancelled_at"`
	ShippingSharedLink         string              `json:"shipping_shared_link"`
	CancelReason               string              `json:"cancel_reason"`
}

func (m *ShipnowFulfillment) String() string { return jsonx.MustMarshalToString(m) }

type GetShipnowFulfillmentsRequest struct {
	Paging  *common.Paging     `json:"paging"`
	Filters []*common.Filter   `json:"filters"`
	Mixed   *etop.MixedAccount `json:"mixed"`
}

func (m *GetShipnowFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowFulfillments struct {
	Paging              *common.PageInfo      `json:"paging"`
	ShipnowFulfillments []*ShipnowFulfillment `json:"shipnow_fulfillments"`
}

func (m *ShipnowFulfillments) String() string { return jsonx.MustMarshalToString(m) }

type DeliveryPoint struct {
	ShippingAddress *OrderAddress `json:"shipping_address"`
	Lines           []*OrderLine  `json:"lines"`
	ShippingNote    string        `json:"shipping_note"`
	OrderId         dot.ID        `json:"order_id"`
	WeightInfo      `json:"weight_info"`
	ValueInfo       `json:"value_info"`
	TryOn           try_on.TryOnCode `json:"try_on"`
}

func (m *DeliveryPoint) String() string { return jsonx.MustMarshalToString(m) }

type WeightInfo struct {
	GrossWeight      int `json:"gross_weight"`
	ChargeableWeight int `json:"chargeable_weight"`
	Length           int `json:"length"`
	Width            int `json:"width"`
	Height           int `json:"height"`
}

func (m *WeightInfo) String() string { return jsonx.MustMarshalToString(m) }

type ValueInfo struct {
	BasketValue      int  `json:"basket_value"`
	CodAmount        int  `json:"cod_amount"`
	IncludeInsurance bool `json:"include_insurance"`
}

func (m *ValueInfo) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipnowFulfillmentRequest struct {
	OrderIds            []dot.ID      `json:"order_ids"`
	Carrier             string        `json:"carrier"`
	ShippingServiceCode string        `json:"shipping_service_code"`
	ShippingServiceFee  int           `json:"shipping_service_fee"`
	ShippingNote        string        `json:"shipping_note"`
	RequestPickupAt     dot.Time      `json:"request_pickup_at"`
	PickupAddress       *OrderAddress `json:"pickup_address"`
}

func (m *CreateShipnowFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipnowFulfillmentV2Request struct {
	DeliveryPoints      []*ShipnowDeliveryPoint     `json:"delivery_points"`
	Carrier             shipnowcarriertypes.Carrier `json:"carrier"`
	ShippingServiceCode string                      `json:"shipping_service_code"`
	ShippingServiceFee  int                         `json:"shipping_service_fee"`
	ShippingNote        string                      `json:"shipping_note"`
	PickupAddress       *OrderAddress               `json:"pickup_address"`
}

func (m *CreateShipnowFulfillmentV2Request) Reset()         { *m = CreateShipnowFulfillmentV2Request{} }
func (m *CreateShipnowFulfillmentV2Request) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowDeliveryPoint struct {
	OrderID          dot.ID        `json:"order_id"`
	ShippingAddress  *OrderAddress `json:"shipping_address"`
	ShippingNote     string        `json:"shipping_note"`
	ChargeableWeight int           `json:"chargeable_weight"`
	GrossWeight      int           `json:"gross_weight"`
	CODAmount        int           `json:"cod_amount"`
}

func (m *ShipnowDeliveryPoint) Reset()         { *m = ShipnowDeliveryPoint{} }
func (m *ShipnowDeliveryPoint) String() string { return jsonx.MustMarshalToString(m) }

type UpdateShipnowFulfillmentRequest struct {
	Id                  dot.ID        `json:"id"`
	OrderIds            []dot.ID      `json:"order_ids"`
	Carrier             string        `json:"carrier"`
	ShippingServiceCode string        `json:"shipping_service_code"`
	ShippingServiceFee  int           `json:"shipping_service_fee"`
	ShippingNote        string        `json:"shipping_note"`
	RequestPickupAt     dot.Time      `json:"request_pickup_at"`
	PickupAddress       *OrderAddress `json:"pickup_address"`
}

func (m *UpdateShipnowFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelShipnowFulfillmentRequest struct {
	Id           dot.ID `json:"id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelShipnowFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipnowServicesRequest struct {
	OrderIds       []dot.ID                `json:"order_ids"`
	PickupAddress  *OrderAddress           `json:"pickup_address"`
	DeliveryPoints []*DeliveryPointRequest `json:"delivery_points"`
}

func (m *GetShipnowServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeliveryPointRequest struct {
	ShippingAddress *OrderAddress `json:"shipping_address"`
	CodAmount       int           `json:"cod_amount"`
}

func (m *DeliveryPointRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShipnowServicesResponse struct {
	Services []*ShippnowService `json:"services"`
}

func (m *GetShipnowServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShippnowService struct {
	Carrier     string `json:"carrier"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Fee         int    `json:"fee"`
	Description string `json:"description"`
}

func (m *ShippnowService) String() string { return jsonx.MustMarshalToString(m) }

type XFulfillment struct {
	Shipnow  *ShipnowFulfillment `json:"shipnow"`
	Shipment *Fulfillment        `json:"shipment"`
	// backward-compatible fields from shipment
	Id         dot.ID       `json:"id"`
	OrderId    dot.ID       `json:"order_id"`
	ShopId     dot.ID       `json:"shop_id"`
	PartnerId  dot.ID       `json:"partner_id"`
	SelfUrl    string       `json:"self_url"`
	Lines      []*OrderLine `json:"lines"`
	TotalItems int          `json:"total_items"`
	// @deprecated use chargeable_weight
	TotalWeight int `json:"total_weight"`
	BasketValue int `json:"basket_value"`
	// @deprecated use cod_amount
	TotalCodAmount int `json:"total_cod_amount"`
	CodAmount      int `json:"cod_amount"`
	// @deprecated
	TotalAmount      int      `json:"total_amount"`
	ChargeableWeight int      `json:"chargeable_weight"`
	CreatedAt        dot.Time `json:"created_at"`
	UpdatedAt        dot.Time `json:"updated_at"`
	ClosedAt         dot.Time `json:"closed_at"`
	CancelledAt      dot.Time `json:"cancelled_at"`
	CancelReason     string   `json:"cancel_reason"`
	// @deprecated use carrier instead
	ShippingProvider     shipping_provider.ShippingProvider `json:"shipping_provider"`
	Carrier              shipping_provider.ShippingProvider `json:"carrier"`
	ShippingServiceName  string                             `json:"shipping_service_name"`
	ShippingServiceFee   int                                `json:"shipping_service_fee"`
	ShippingServiceCode  string                             `json:"shipping_service_code"`
	ShippingCode         string                             `json:"shipping_code"`
	ShippingNote         string                             `json:"shipping_note"`
	TryOn                try_on.TryOnCode                   `json:"try_on"`
	IncludeInsurance     bool                               `json:"include_insurance"`
	ShConfirm            status3.Status                     `json:"sh_confirm"`
	ShippingState        shipping.State                     `json:"shipping_state"`
	Status               status5.Status                     `json:"status"`
	ShippingStatus       status5.Status                     `json:"shipping_status"`
	EtopPaymentStatus    status4.Status                     `json:"etop_payment_status"`
	ShippingFeeCustomer  int                                `json:"shipping_fee_customer"`
	ShippingFeeShop      int                                `json:"shipping_fee_shop"`
	XShippingFee         int                                `json:"x_shipping_fee"`
	XShippingId          string                             `json:"x_shipping_id"`
	XShippingCode        string                             `json:"x_shipping_code"`
	XShippingCreatedAt   dot.Time                           `json:"x_shipping_created_at"`
	XShippingUpdatedAt   dot.Time                           `json:"x_shipping_updated_at"`
	XShippingCancelledAt dot.Time                           `json:"x_shipping_cancelled_at"`
	XShippingDeliveredAt dot.Time                           `json:"x_shipping_delivered_at"`
	XShippingReturnedAt  dot.Time                           `json:"x_shipping_returned_at"`
	// @deprecated use estimated_delivery_at
	ExpectedDeliveryAt dot.Time `json:"expected_delivery_at"`
	// @deprecated use estimated_pickup_at
	ExpectedPickAt              dot.Time               `json:"expected_pick_at"`
	EstimatedDeliveryAt         dot.Time               `json:"estimated_delivery_at"`
	EstimatedPickupAt           dot.Time               `json:"estimated_pickup_at"`
	CodEtopTransferedAt         dot.Time               `json:"cod_etop_transfered_at"`
	ShippingFeeShopTransferedAt dot.Time               `json:"shipping_fee_shop_transfered_at"`
	XShippingState              string                 `json:"x_shipping_state"`
	XShippingStatus             status5.Status         `json:"x_shipping_status"`
	XSyncStatus                 status4.Status         `json:"x_sync_status"`
	XSyncStates                 *FulfillmentSyncStates `json:"x_sync_states"`
	// @deprecated use shipping_address instead
	AddressTo *etop.Address `json:"address_to"`
	// @deprecated use pickup_address instead
	AddressFrom                        *etop.Address          `json:"address_from"`
	PickupAddress                      *OrderAddress          `json:"pickup_address"`
	ReturnAddress                      *OrderAddress          `json:"return_address"`
	ShippingAddress                    *OrderAddress          `json:"shipping_address"`
	Shop                               *etop.Shop             `json:"shop"`
	Order                              *Order                 `json:"order"`
	ProviderShippingFeeLines           []*ShippingFeeLine     `json:"provider_shipping_fee_lines"`
	ShippingFeeShopLines               []*ShippingFeeLine     `json:"shipping_fee_shop_lines"`
	EtopDiscount                       int                    `json:"etop_discount"`
	MoneyTransactionShippingId         dot.ID                 `json:"money_transaction_shipping_id"`
	MoneyTransactionShippingExternalId dot.ID                 `json:"money_transaction_shipping_external_id"`
	XShippingLogs                      []*ExternalShippingLog `json:"x_shipping_logs"`
	XShippingNote                      string                 `json:"x_shipping_note"`
	XShippingSubState                  string                 `json:"x_shipping_sub_state"`
	Code                               string                 `json:"code"`
	ActualCompensationAmount           int                    `json:"actual_compensation_amount"`
}

func (m *XFulfillment) String() string { return jsonx.MustMarshalToString(m) }

type TradingCreateOrderRequest struct {
	// Customer should be shop's information
	// and customer_address, shipping_address must be shop address.
	Customer        *OrderCustomer               `json:"customer"`
	CustomerAddress *OrderAddress                `json:"customer_address"`
	BillingAddress  *OrderAddress                `json:"billing_address"`
	ShippingAddress *OrderAddress                `json:"shipping_address"`
	Lines           []*CreateOrderLine           `json:"lines"`
	Discounts       []*OrderDiscount             `json:"discounts"`
	TotalItems      int                          `json:"total_items"`
	BasketValue     int                          `json:"basket_value"`
	OrderDiscount   int                          `json:"order_discount"`
	TotalFee        int                          `json:"total_fee"`
	FeeLines        []*OrderFeeLine              `json:"fee_lines"`
	TotalDiscount   dot.NullInt                  `json:"total_discount"`
	TotalAmount     int                          `json:"total_amount"`
	OrderNote       string                       `json:"order_note"`
	PaymentMethod   payment_method.PaymentMethod `json:"payment_method"`
	ReferralMeta    map[string]string            `json:"referral_meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *TradingCreateOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShippingServicesRequest struct {
	ConnectionIDs    []dot.ID     `json:"connection_ids"`
	FromDistrictCode string       `json:"from_district_code"`
	FromProvinceCode string       `json:"from_province_code"`
	ToDistrictCode   string       `json:"to_district_code"`
	ToProvinceCode   string       `json:"to_province_code"`
	FromProvince     string       `json:"from_province"`
	FromDistrict     string       `json:"from_district"`
	ToProvince       string       `json:"to_province"`
	ToDistrict       string       `json:"to_district"`
	GrossWeight      int          `json:"gross_weight"`
	ChargeableWeight int          `json:"chargeable_weight"`
	Length           int          `json:"length"`
	Width            int          `json:"width"`
	Height           int          `json:"height"`
	TotalCodAmount   int          `json:"total_cod_amount"`
	BasketValue      int          `json:"basket_value"`
	IncludeInsurance dot.NullBool `json:"include_insurance"`
}

func (m *GetShippingServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShippingServicesResponse struct {
	Services []*ExternalShippingService `json:"services"`
}

func (m *GetShippingServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShippingService struct {
	Name                string   `json:"name"`
	Code                string   `json:"code"`
	Fee                 int      `json:"fee"`
	ConnectionID        dot.ID   `json:"connection_id"`
	EstimatedPickupAt   dot.Time `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time `json:"estimated_delivery_at"`
}

func (m *ShippingService) String() string { return jsonx.MustMarshalToString(m) }
