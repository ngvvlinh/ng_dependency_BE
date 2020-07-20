package clientv2

import (
	"encoding/json"
	"strings"
	"time"

	typesshipping "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status5"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/integration/shipping"
	"o.o/common/strs"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type State string

type GHNAccountCfg struct {
	ClientID    int    `yaml:"client_id"`
	ShopID      int    `json:"shop_id"`
	Token       string `yaml:"token"`
	AffiliateID int    `yaml:"affiliate_id"`
}

const (
	StateReadyToPick        State = "ready_to_pick"       // -> picking
	StateCancel             State = "cancel"              // -> cancel
	StatePicking            State = "picking"             // -> picking
	StatePicked             State = "picked"              // -> holding
	StateStoring            State = "storing"             // -> holding
	StateTransporting       State = "transporting"        // -> delivering
	StateDelivering         State = "delivering"          // -> delivering
	StateDeliveryFail       State = "delivery_fail"       // -> delivering
	StateDelivered          State = "delivered"           // -> delivered
	StateWaitingToReturn    State = "waiting_to_return"   // -> delivering
	StateReturn             State = "return"              // -> delivering
	StateReturnTransporting State = "return_transporting" // -> delivering
	StateReturning          State = "returning"           // -> delivering
	StateReturnFail         State = "return_fail"         // -> delivering
	StateReturned           State = "returned"            // -> returned
	StateDamage             State = "damage"              // -> undeliverable
	StateLost               State = "lost"                // -> undeliverable
)

func (s State) ToModel() typesshipping.State {
	switch s {
	case StateReadyToPick:
		return typesshipping.Created
	case StateCancel:
		return typesshipping.Cancelled
	case StatePicking:
		return typesshipping.Picking
	case StatePicked:
		return typesshipping.Holding
	case StateStoring:
		return typesshipping.Holding
	case StateTransporting:
		return typesshipping.Delivering
	case StateDelivering:
		return typesshipping.Delivering
	case StateDeliveryFail:
		return typesshipping.Returning
	case StateWaitingToReturn:
		return typesshipping.Returning
	case StateReturn:
		return typesshipping.Returning
	case StateReturnTransporting:
		return typesshipping.Returning
	case StateReturning:
		return typesshipping.Returning
	case StateReturnFail:
		return typesshipping.Returning
	case StateReturned:
		return typesshipping.Returned
	case StateDelivered:
		return typesshipping.Delivered
	case StateDamage:
		return typesshipping.Undeliverable
	case StateLost:
		return typesshipping.Undeliverable
	default:
		return typesshipping.Unknown
	}
}

func (s State) ToStatus5() status5.Status {
	switch s.ToModel() {
	case typesshipping.Cancelled:
		return status5.N
	case typesshipping.Returned, typesshipping.Returning, typesshipping.Undeliverable:
		return status5.NS
	case typesshipping.Delivered:
		return status5.P
	}

	return status5.S
}

type ErrorResponse struct {
	Code        Int             `json:"code"`
	Message     String          `json:"message"`
	Data        json.RawMessage `json:"data"`
	CodeMessage String          `json:"code_message"`

	ErrorData map[string]string `json:"-"`
}

func (e *ErrorResponse) Error() (s string) {
	defer func() {
		s = strs.TrimLastPunctuation(s)
	}()

	if len(e.ErrorData) == 0 {
		return e.Message.String()
	}
	b := strings.Builder{}
	b.WriteString(e.Message.String())
	b.WriteString(" (")
	for _, v := range e.ErrorData {
		b.WriteString(v)
		break
	}
	b.WriteString(")")
	return b.String()
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}

func (o *OrderFee) ToFeeLines() []*shippingsharemodel.ShippingFeeLine {
	var res []*shippingsharemodel.ShippingFeeLine
	if o.Coupon != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Discount,
			Cost:            o.Coupon.Int(),
		})
	}
	if o.Insurance != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Insurance,
			Cost:            o.Insurance.Int(),
		})
	}
	if o.MainService != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Main,
			Cost:            o.MainService.Int(),
		})
	}
	if o.R2S != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Redelivery,
			Cost:            o.R2S.Int(),
		})
	}
	if o.Return != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Return,
			Cost:            o.Return.Int(),
		})
	}
	if o.StationPu != 0 || o.StationDo != 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType: shipping_fee_type.Other,
			Cost:            o.StationPu.Int() + o.StationDo.Int(),
		})
	}
	return res
}

func (o *OrderFeeCallback) ToOrderFee() *OrderFee {
	return &OrderFee{
		MainService: o.MainService,
		Insurance:   o.Insurance,
		StationDo:   o.StationDO,
		StationPu:   o.StationPU,
		Return:      o.Return,
		R2S:         o.R2S,
		Coupon:      o.Coupon,
	}
}

// Send OTP affiliate

type SendOTPShopAffiliateRequest struct {
	Phone string `json:"phone"`
}

type SendOTPShopAffiliateResponse struct {
	TTL Int `json:"TTL"`
}

// Create shop affiliate

type CreateShopAffiliateRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type CreateShopAffiliateResponse struct {
	ShopID   Int `json:"shop_id"`
	ClientID Int `json:"client_id"`
}

// API Add Staff to Store by OTP (Affiliate)

type AffiliateCreateWithShopRequest struct {
	Phone  string `json:"phone"`
	OTP    string `json:"otp"`
	ShopID int    `json:"shop_id"`
}

// Get shops

type ShopsRequest struct {
	Offset Int `json:"offset"`
	Limit  Int `json:"limit"`
}

type ShopsResponse struct {
	LastOffset Int         `json:"last_offset"`
	Shops      []*ShopInfo `json:"shops"`
}

type ShopInfo struct {
	ID              Int    `json:"_id"`
	Name            String `json:"name"`
	Phone           String `json:"phone"`
	Address         String `json:"address"`
	WardCode        String `json:"ward_code"`
	DistrictID      Int    `json:"district_id"`
	ClientID        Int    `json:"client_id"`
	BankAccountID   Int    `json:"bank_account_id"`
	Status          Int    `json:"status"`
	VersionNo       String `json:"version_no"`
	UpdatedIP       String `json:"updated_ip"`
	UpdatedEmployee Int    `json:"update_employee"`
	UpdatedClient   Int    `json:"updated_client"`
	UpdatedSource   String `json:"updated_source"`
	UpdatedDate     Time   `json:"updated_date"`
	CreatedIP       String `json:"created_ip"`
	CreatedEmployee Int    `json:"created_employee"`
	CreatedClient   Int    `json:"created_client"`
	CreatedSource   String `json:"created_source"`
	CreatedDate     Time   `json:"created_date"`
}

// Get provinces

type ProvincesResponse struct {
	ProvinceID   Int    `json:"ProvinceID"`
	ProvinceName String `json:"ProvinceName"`
	Code         String `json:"Code"`
}

// Get districts

type DistrictsRequest struct {
	ProvinceID Int `json:"province_id"`
}

type DistrictsResponse struct {
	DistrictID   Int    `json:"DistrictID"`
	ProvinceID   Int    `json:"ProvinceID"`
	DistrictName String `json:"DistrictName"`
	Code         String `json:"Code"`
	Type         Int    `json:"Type"`
	SupportType  Int    `json:"SupportType"`
}

// Get wards

type GetWardsRequest struct {
	DistrictID int `json:"district_id"`
}

type GetWardsResponse []*WardInfo

type WardInfo struct {
	WardCode   Int    `json:"WardCode"`
	DistrictID Int    `json:"DistrictID"`
	WardName   String `json:"WardName"`
}

// Services

type FindAvailableServicesRequest struct {
	Weight int `json:"Weight"` // 10000,
	Length int `json:"Length"` // 10,
	Width  int `json:"Width"`  // 110,
	Height int `json:"Height"` // 20,

	FromDistrictID int    `json:"FromDistrictID"` // 1455,
	FromWardCode   string `json:"FromWardCode"`
	ToDistrictID   int    `json:"ToDistrictID"` // 1462
	ToWardCode     string `json:"ToWardCode"`
	InsuranceFee   int    `json:"InsuranceFee"`
	Coupon         string `json:"Coupon"`
}

type FindAvailableServicesResponse struct {
	AvailableServices []*AvailableService
}

type AvailableService struct {
	ExpectedDeliveryTime Time `json:"ExpectedDeliveryTime"`

	Name           String       `json:"Name"`
	ServiceFee     Int          `json:"ServiceFee"`
	ServiceFeeMain Int          `json:"-"`
	ServiceID      Int          `json:"ServiceID"`
	Extras         ExtraService `json:"Extra"`

	IsPromotion bool `json:"-"`
}

type ExtraService struct {
	MaxValue   Int    `json:"MaxValue"`   // 0,
	Name       String `json:"Name"`       // "Gửi hàng tại điểm",
	ServiceFee Int    `json:"ServiceFee"` // -2000,
	ServiceID  Int    `json:"ServiceID"`  // 53337
}

// Get services

type GetServicesRequest struct {
	FromDistrict Int `json:"from_district"`
	ToDistrict   Int `json:"to_district"`
}

type GetServicesResponse []*ServiceInfo

type ServiceInfo struct {
	ServiceID     Int    `json:"service_id"`
	ShortName     String `json:"short_name"`
	ServiceTypeID Int    `json:"service_type_id"`
}

// Get shipping fee

type CalculateFeeRequest struct {
	ShopID         int    `json:"shop_id"`
	ServiceID      int    `json:"service_id"`
	ServiceTypeID  int    `json:"service_type_id"`
	InsuranceValue int    `json:"insurance_value"`
	Coupon         string `json:"coupon"`
	FromDistrictID int    `json:"from_district_id"`
	FromWardCode   string `json:"from_ward_code"`
	ToDistrictID   int    `json:"to_district_id"`
	ToWardCode     string `json:"to_ward_code"`
	Height         int    `json:"height"`
	Length         int    `json:"length"`
	Weight         int    `json:"weight"`
	Width          int    `json:"width"`
}

type CalculateFeeResponse struct {
	Total          Int `json:"total"`
	ServiceFee     Int `json:"service_fee"`
	InsuranceFee   Int `json:"insurance_fee"`
	PickStationFee Int `json:"pick_station_fee"`
	CouponValue    Int `json:"coupon_value"`
	R2SFee         Int `json:"r2s_fee"`
}

// Get leadtime

type GetLeadTimeRequest struct {
	FromDistrictID int    `json:"from_district_id"`
	FromWardCode   string `json:"from_ward_code"`
	ToDistrictID   int    `json:"to_district_id"`
	ToWardCode     string `json:"to_ward_code"`
	ServiceID      int    `json:"service_id"`
}

type GetLeadTimeResponse struct {
	LeadTime  Time `json:"leadtime"`
	OrderDate Time `json:"order_date"`
}

// Create order

type CreateOrderRequest struct {
	FromName         string `json:"from_name"`
	FromPhone        string `json:"from_phone"`
	FromAddress      string `json:"from_address"`
	FromWardCode     string `json:"from_ward_code"`
	FromDistrictID   int    `json:"from_district_id"`
	ToName           string `json:"to_name"`
	ToPhone          string `json:"to_phone"`
	ToAddress        string `json:"to_address"`
	ToWardCode       string `json:"to_ward_code"`
	ToDistrictID     int    `json:"to_district_id"`
	ReturnPhone      string `json:"return_phone"`
	ReturnAddress    string `json:"return_address"`
	ReturnDistrictID int    `json:"return_district_id"`
	ReturnWardCode   string `json:"return_ward_code"`
	ClientOrderCode  string `json:"client_order_code"`
	CODAmount        int    `json:"cod_amount"`
	Content          string `json:"content"`
	Weight           int    `json:"weight"`
	Length           int    `json:"length"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	PickStationID    int    `json:"pick_station_id"`
	InsuranceValue   int    `json:"insurance_value"`
	Coupon           string `json:"coupon"`
	ServiceTypeID    int    `json:"service_type_id"`
	ServiceID        int    `json:"service_id"`
	PaymentTypeID    int    `json:"payment_type_id"`
	Note             string `json:"note"`
	RequiredNote     string `json:"required_note"`
}

type CreateOrderResponse struct {
	OrderCode            String    `json:"order_code"`
	SortCode             String    `json:"sort_code"`
	TransType            String    `json:"trans_type"`
	WardEncode           String    `json:"ward_encode"`
	DistrictEncode       String    `json:"district_encode"`
	Fee                  *OrderFee `json:"fee"`
	TotalFee             Int       `json:"total_fee"`
	ExpectedDeliveryTime Time      `json:"expected_delivery_time"`
}

type OrderFee struct {
	MainService Int `json:"main_service"` // Phí giao hàng
	Insurance   Int `json:"insurance"`    // Phí khai giá
	StationDo   Int `json:"station_do"`   // Phí khách hàng tới station lấy hàng
	StationPu   Int `json:"station_pu"`   // Phí khách hàng mang hàng tới station
	Return      Int `json:"return"`       // Phí trả hàng
	R2S         Int `json:"r2s"`          // Phí giao lại
	Coupon      Int `json:"coupon"`       // Giảm giá
}

// Update order

type UpdateOrderRequest struct {
	OrderCode      string `json:"order_code"`
	FromName       string `json:"from_name"`
	FromPhone      string `json:"from_phone"`
	FromAddress    string `json:"from_address"`
	FromWardCode   string `json:"from_ward_code"`
	FromDistrictID int    `json:"from_district_id"`
	ToName         string `json:"to_name"`
	ToPhone        string `json:"to_phone"`
	ToAddress      string `json:"to_address"`
	ToWardCode     string `json:"to_ward_code"`
	ToDistrictID   int    `json:"to_district_id"`
	Weight         int    `json:"weight,omitempty"`
	InsuranceValue *int   `json:"insurance_value,omitempty"`
	Note           string `json:"note,omitempty"`
	RequiredNote   string `json:"required_note,omitempty"`
}

// UpdateCOD

type UpdateOrderCODRequest struct {
	OrderCode string `json:"order_code"`
	CODAmount int    `json:"cod_amount"`
}

// Get detail order

type GetOrderInfoRequest struct {
	OrderCode String `json:"order_code"`
}

type GetOrderInfoResponse struct {
	ID                 String            `json:"_id"`
	OrderCode          String            `json:"order_code"`
	ShopID             Int               `json:"shop_id"`
	ClientID           Int               `json:"client_id"`
	ReturnName         String            `json:"return_name"`
	ReturnPhone        String            `json:"return_phone"`
	ReturnAddress      String            `json:"return_address"`
	ReturnWardCode     String            `json:"return_ward_code"`
	ReturnDistrictID   Int               `json:"return_district_id"`
	FromName           String            `json:"from_name"`
	FromPhone          String            `json:"from_phone"`
	FromAddress        String            `json:"from_address"`
	FromWardCode       String            `json:"from_ward_code"`
	FromDistrictID     Int               `json:"from_district_id"`
	DeliverStationID   Int               `json:"deliver_station_id"`
	ToName             String            `json:"to_name"`
	ToPhone            String            `json:"to_phone"`
	ToAddress          String            `json:"to_address"`
	ToWardCode         String            `json:"to_ward_code"`
	ToDistrictID       Int               `json:"to_district_id"`
	Weight             Int               `json:"weight"`
	Length             Int               `json:"length"`
	Width              Int               `json:"width"`
	Height             Int               `json:"height"`
	ConvertedWeight    Int               `json:"converted_weight"`
	ServiceTypeID      Int               `json:"service_type_id"`
	ServiceID          Int               `json:"service_id"`
	PaymentTypeID      Int               `json:"payment_type_id"`
	CustomServiceFee   Int               `json:"custom_service_fee"`
	CODAmount          Int               `json:"cod_amount"`
	CODCollectDate     Time              `json:"cod_collect_date"`
	CODTransferDate    Time              `json:"cod_transfer_date"`
	IsCODTransferred   Bool              `json:"is_cod_transferred"`
	IsCODCollected     Bool              `json:"is_cod_collected"`
	InsuranceValue     Int               `json:"insurance_value"`
	OrderValue         Int               `json:"order_value"`
	PickStationID      Int               `json:"pick_station_id"`
	ClientOrderCode    String            `json:"client_order_code"`
	RequiredNote       String            `json:"required_note"`
	Content            String            `json:"content"`
	Note               String            `json:"note"`
	EmployeeNote       String            `json:"employee_note"`
	Coupon             String            `json:"coupon"`
	VersionNo          String            `json:"version_no"`
	UpdatedIP          String            `json:"updated_ip"`
	UpdatedEmployee    Int               `json:"updated_employee"`
	UpdatedClient      Int               `json:"updated_client"`
	UpdatedSource      String            `json:"updated_source"`
	UpdatedDate        Time              `json:"updated_date"`
	UpdatedWarehouse   Int               `json:"updated_warehouse"`
	CreatedIP          String            `json:"created_ip"`
	CreatedEmployee    Int               `json:"created_employee"`
	CreatedClient      Int               `json:"created_client"`
	CreatedSource      String            `json:"created_source"`
	CreatedDate        Time              `json:"created_date"`
	Status             String            `json:"status"`
	PickWareHouseID    Int               `json:"pick_ware_house_id"`
	DeliverWarehouseID Int               `json:"deliver_warehouse_id"`
	CurrentWarehouseID Int               `json:"current_warehouse_id"`
	ReturnWarehouseID  Int               `json:"return_warehouse_id"`
	NextWarehouseID    Int               `json:"next_warehouse_id"`
	Leadtime           Time              `json:"leadtime"`
	OrderDate          Time              `json:"order_date"`
	SOCID              String            `json:"soc_id"`
	FinishDate         Time              `json:"finish_date"`
	Tag                []String          `json:"tag"`
	Log                []*OrderDetailLog `json:"log"`
}

type OrderDetailLog struct {
	Status      String `json:"status"`
	UpdatedDate String `json:"updated_date"`
}

// Cancel Order
type CancelOrderRequest struct {
	OrderCodes []string `json:"order_codes"`
}

type CancelOrderResponse struct {
	OrderCode String `json:"order_code"`
	Result    Bool   `json:"result"`
	Message   String `json:"message"`
}

type OrderCost struct {
	Cost      int    `json:"Cost"`
	Name      string `json:"Name"`
	ServiceID int    `json:"ServiceID"`
}

func (s *AvailableService) ToShippingService(providerServiceID string) *shippingsharemodel.AvailableShippingService {
	if s == nil {
		return nil
	}
	serviceFeeMain := cm.CoalesceInt(int(s.ServiceFeeMain), int(s.ServiceFee))
	return &shippingsharemodel.AvailableShippingService{
		Name:              s.Name.String(),
		ServiceFee:        int(s.ServiceFee),
		ShippingFeeMain:   serviceFeeMain,
		Provider:          shipping_provider.GHN,
		ProviderServiceID: providerServiceID,

		ExpectedPickAt:     shipping.CalcPickTime(shipping_provider.GHN, time.Now()),
		ExpectedDeliveryAt: s.ExpectedDeliveryTime.ToTime(),
	}
}

// Use for webhook
type CallbackOrder struct {
	CODAmount       Int               `json:"CODAmount"`
	CODTransferDate Time              `json:"CODTransferDate"`
	ClientOrderCode String            `json:"ClientOrderCode"`
	Description     String            `json:"Description"`
	Fee             *OrderFeeCallback `json:"Fee"`
	Height          Int               `json:"Height"`
	Length          Int               `json:"Length"`
	OrderCode       String            `json:"OrderCode"`
	Reason          String            `json:"Reason"`
	ReasonCode      String            `json:"ReasonCode"`
	ShipperName     String            `json:"ShipperName"`
	ShipperPhone    String            `json:"ShipperPhone"`
	Status          String            `json:"Status"`
	Time            Time              `json:"Time"`
	TotalFee        Int               `json:"TotalFee"`
	Type            String            `json:"Type"`
	WareHouse       String            `json:"WareHouse"`
	Weight          Int               `json:"Weight"`
	Width           Int               `json:"Width"`
}

type OrderFeeCallback struct {
	Coupon      Int `json:"Coupon"`
	Insurance   Int `json:"Insurance"`
	MainService Int `json:"MainService"`
	R2S         Int `json:"R2S"`
	Return      Int `json:"Return"`
	StationDO   Int `json:"StationDO"`
	StationPU   Int `json:"StationPU"`
}

type GetShopByClientOwnerRequest struct {
	OffsetID int    `json:"offset_id"`
	Phone    string `json:"phone"`
	Limit    int    `json:"limit"`
}

type GetShopByClientOwnerResponse []*ShopInfo
