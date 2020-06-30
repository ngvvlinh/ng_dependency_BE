package admin

import (
	"time"

	etop "o.o/api/top/int/etop"
	"o.o/api/top/int/types"
	common "o.o/api/top/types/common"
	credit_type "o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/filter_type"
	"o.o/api/top/types/etc/location_type"
	notifier_entity "o.o/api/top/types/etc/notifier_entity"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/route_type"
	shipping "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	status3 "o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
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

type GetShopsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetShopsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UsersFilter struct {
	Name      string      `json:"name"`
	Phone     string      `json:"phone"`
	Email     string      `json:"email"`
	CreatedAt filter.Date `json:"created_at"`
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
}

func (m *UpdateFulfillmentShippingStateRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentShippingFeesRequest struct {
	ID               dot.ID                   `json:"id"`
	ShippingCode     string                   `json:"shipping_code"`
	ShippingFeeLines []*types.ShippingFeeLine `json:"shipping_fee_lines"`
	TotalCODAmount   dot.NullInt              `json:"total_cod_amount"`
}

func (m *UpdateFulfillmentShippingFeesRequest) String() string { return jsonx.MustMarshalToString(m) }

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
	FeeType shipping_fee_type.ShippingFeeType `json:"fee_type"`
	Rules   []*AdditionalFeeRule              `json:"rules"`
}

func (m *AdditionalFee) String() string { return jsonx.MustMarshalToString(m) }

type AdditionalFeeRule struct {
	MinValue          int                                   `json:"min_value"`
	MaxValue          int                                   `json:"max_value"`
	PriceModifierType price_modifier_type.PriceModifierType `json:"price_modifier_type"`
	Amount            float64                               `json:"amount"`
	MinPrice          int                                   `json:"min_price"`
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

type UserResponse struct {
	Users  []*etop.User           `json:"users"`
	Paging *common.CursorPageInfo `json:"paging"`
}

func (m *UserResponse) String() string { return jsonx.MustMarshalToString(m) }
