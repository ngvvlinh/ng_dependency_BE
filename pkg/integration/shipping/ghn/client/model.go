package client

import (
	"encoding/json"
	"strconv"
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
type ServiceFeeID string
type ServiceFeeType string

type GHNAccountCfg struct {
	ClientID    int    `yaml:"client_id"`
	AffiliateID int    `yaml:"affiliate_id"`
	Token       string `yaml:"token"`
}

const (
	StateReadyToPick     State = "ReadyToPick"
	StatePicking         State = "Picking"
	StateStoring         State = "Storing"
	StateDelivering      State = "Delivering"
	StateDelivered       State = "Delivered"
	StateReturn          State = "Return"
	StateReturned        State = "Returned"
	StateWaitingToFinish State = "WaitingToFinish"
	StateFinish          State = "Finish"
	StateCancel          State = "Cancel"
	StateLostOrder       State = "LostOrder"

	ServiceFeeInsurance       ServiceFeeID = "16"
	ServiceFeeCostDeclaration ServiceFeeID = "53332"
	ServiceFeeReturnType1     ServiceFeeID = "100026" // Trả hàng nội tỉnh ( giống địa chỉ lấy )
	ServiceFeeReturnType2     ServiceFeeID = "100027" // Trả hàng nội tỉnh ( khác địa chỉ lấy )
	ServiceFeeReturnType3     ServiceFeeID = "100028" // Trả hàng liên tỉnh ( giống địa chỉ lấy )
	ServiceFeeReturnType4     ServiceFeeID = "100029" // Trả hàng liên tỉnh ( khác địa chỉ lấy )
	ServiceFeeAddressChange   ServiceFeeID = "100030" // Phí thay đổi địa chỉ giao
	ServiceFeeAdjustment      ServiceFeeID = "100032" // Phí chênh lệch thay đổi đơn hàng

	// https://api.ghn.vn/home/docs/detail?id=23
	ServiceFee6Hours   ServiceFeeID   = "53319"
	ServiceFee1Day     ServiceFeeID   = "53320"
	ServiceFee2Days    ServiceFeeID   = "53321"
	ServiceFee3Days    ServiceFeeID   = "53322"
	ServiceFee4Days    ServiceFeeID   = "53323"
	ServiceFee5Days    ServiceFeeID   = "53324"
	ServiceFee6Days    ServiceFeeID   = "53327"
	ServiceFeeExtend   ServiceFeeID   = "53337"
	ServiceFeeTypeMain ServiceFeeType = "1"
)

func (s State) ToModel(old typesshipping.State, isReturnOrder bool) typesshipping.State {
	switch s {
	case StateReadyToPick:
		return typesshipping.Created
	case StatePicking:
		return typesshipping.Picking
	case StateStoring:
		return typesshipping.Holding
	case StateDelivering:
		return typesshipping.Delivering
	case StateDelivered:
		return typesshipping.Delivered
	case StateReturn:
		return typesshipping.Returning
	case StateReturned:
		return typesshipping.Returned
	case StateCancel:
		return typesshipping.Cancelled
	case StateLostOrder:
		return typesshipping.Undeliverable
	case StateWaitingToFinish, StateFinish:
		switch old {
		case typesshipping.Returning, typesshipping.Returned:
			return typesshipping.Returned
		case typesshipping.Cancelled, typesshipping.Delivered, typesshipping.Undeliverable:
			return old
		}
		if isReturnOrder {
			return typesshipping.Returned
		}
		return typesshipping.Delivered
	default:
		return typesshipping.Unknown
	}
}

func (s State) ToStatus5(old typesshipping.State) status5.Status {
	switch s {
	case StateCancel:
		return status5.N

	case StateReturned:
		return status5.NS

	case StateFinish:
		switch old {
		case typesshipping.Cancelled:
			return status5.N
		case typesshipping.Returned, typesshipping.Returning:
			return status5.NS
		default:
			return status5.P
		}
	}

	return status5.S
}

func (s State) ToShippingStatus5(old typesshipping.State) status5.Status {
	switch s {
	case StateCancel:
		return status5.N

	case StateReturned:
		return status5.NS

	case StateWaitingToFinish, StateFinish:
		switch old {
		case typesshipping.Cancelled:
			return status5.N
		case typesshipping.Returned, typesshipping.Returning:
			return status5.NS
		default:
			return status5.P
		}
	}

	return status5.S
}

func (id ServiceFeeID) ToModel() shipping_fee_type.ShippingFeeType {
	switch id {
	// case ServiceFee6Hours, ServiceFee1Day, ServiceFee2Days,
	// 	ServiceFee3Days, ServiceFee4Days, ServiceFee5Days,
	// 	ServiceFee6Days, ServiceFeeExtend:
	// 	return model.Main
	case ServiceFeeInsurance, ServiceFeeCostDeclaration:
		return shipping_fee_type.Insurance
	case ServiceFeeReturnType1, ServiceFeeReturnType2, ServiceFeeReturnType3, ServiceFeeReturnType4:
		return shipping_fee_type.Return
	case ServiceFeeAdjustment:
		return shipping_fee_type.Adjustment
	case ServiceFeeAddressChange:
		return shipping_fee_type.AddressChange
	default:
		return shipping_fee_type.Other
	}
}

func (s *ShippingOrderCost) CalcAndConvertShippingFee() *shippingsharemodel.ShippingFeeLine {
	if s == nil {
		return nil
	}

	cost, err := strconv.Atoi(s.Cost.String())
	if err != nil {
		cost = 0
	}

	var shippingType shipping_fee_type.ShippingFeeType
	sType := ServiceFeeType(s.ServiceType)
	sID := ServiceFeeID(s.ServiceID)
	if sType == ServiceFeeTypeMain {
		shippingType = shipping_fee_type.Main
	} else if cost < 0 {
		shippingType = shipping_fee_type.Discount
	} else {
		shippingType = sID.ToModel()
	}

	return &shippingsharemodel.ShippingFeeLine{
		ShippingFeeType:          shippingType,
		Cost:                     cost,
		ExternalServiceID:        s.ServiceID.String(),
		ExternalServiceName:      s.ServiceName.String(),
		ExternalServiceType:      s.ServiceType.String(),
		ExternalPaymentChannelID: s.PaymentChannelID.String(),
		ExternalShippingOrderID:  s.ShippingOrderID.String(),
	}
}

func CalcAndConvertShippingFeeLines(items []*ShippingOrderCost) []*shippingsharemodel.ShippingFeeLine {
	res := make([]*shippingsharemodel.ShippingFeeLine, len(items))
	for i, item := range items {
		res[i] = item.CalcAndConvertShippingFee()
	}
	return res
}

type ErrorResponse struct {
	Code Int             `json:"code"`
	Msg  String          `json:"msg"`
	Data json.RawMessage `json:"data"`

	ErrorData map[string]string `json:"-"`
}

func (e *ErrorResponse) Error() (s string) {
	defer func() {
		s = strs.TrimLastPunctuation(s)
	}()

	if len(e.ErrorData) == 0 {
		return e.Msg.String()
	}
	b := strings.Builder{}
	b.WriteString(e.Msg.String())
	b.WriteString(" (")
	for _, v := range e.ErrorData {
		b.WriteString(v)
		break
	}
	b.WriteString(")")
	return b.String()
}

type Connection struct {
	// The token field name is lowercase!
	Token string `json:"token"`
}

//-- Requests & Responses

type CreateOrderRequest struct {
	Connection

	PaymentTypeID        int     `json:"PaymentTypeID"`        // 1,
	FromDistrictID       int     `json:"FromDistrictID"`       // 1455,
	FromWardCode         string  `json:"FromWardCode"`         // "21402",
	ToDistrictID         int     `json:"ToDistrictID"`         // 1462,
	ToWardCode           string  `json:"ToWardCode"`           // "21609",
	Note                 string  `json:"Note"`                 // "Tạo ĐH qua API",
	SealCode             string  `json:"SealCode"`             // "tem niêm phong",
	ExternalCode         string  `json:"ExternalCode"`         // "",
	ClientContactName    string  `json:"ClientContactName"`    // "client name",
	ClientContactPhone   string  `json:"ClientContactPhone"`   // "0987654321",
	ClientAddress        string  `json:"ClientAddress"`        // "140 Lê Trọng Tấn",
	CustomerName         string  `json:"CustomerName"`         // "Nguyễn Văn A",
	CustomerPhone        string  `json:"CustomerPhone"`        // "01666666666",
	ShippingAddress      string  `json:"ShippingAddress"`      // "137 Lê Quang Định",
	CoDAmount            int     `json:"CoDAmount"`            // 1500000,
	NoteCode             string  `json:"NoteCode"`             // "CHOXEMHANGKHONGTHU",
	InsuranceFee         int     `json:"InsuranceFee"`         // 0,
	ClientHubID          int     `json:"ClientHubID"`          // 0,
	ServiceID            int     `json:"ServiceID"`            // 53319,
	ToLatitude           float64 `json:"ToLatitude"`           // 1.2343322,
	ToLongitude          float64 `json:"ToLongitude"`          // 10.54324322,
	FromLat              float64 `json:"FromLat"`              // 1.2343322,
	FromLng              float64 `json:"FromLng"`              // 10.54324322,
	Content              string  `json:"Content"`              // "Test nội dung",
	CouponCode           string  `json:"CouponCode"`           // "",
	Weight               int     `json:"Weight"`               // 10200,
	Length               int     `json:"Length"`               // 10,
	Width                int     `json:"Width"`                // 10,
	Height               int     `json:"Height"`               // 10,
	CheckMainBankAccount bool    `json:"CheckMainBankAccount"` // false
	ReturnContactName    string  `json:"ReturnContactName"`    // "",
	ReturnContactPhone   string  `json:"ReturnContactPhone"`   // "",
	ReturnAddress        string  `json:"ReturnAddress"`        // "",
	ReturnDistrictID     int     `json:"ReturnDistrictID"`     // "",
	ExternalReturnCode   string  `json:"ExternalReturnCode"`   // "",
	IsCreditCreate       bool    `json:"IsCreditCreate"`       // true,

	ShippingOrderCosts []ShippingOrderCostRequest `json:"ShippingOrderCosts"`
	// AffiliateID: This field is ClientID, use for tracking commission between GHN and Client.
	AffiliateID int `json:"AffiliateID"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.FromDistrictID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu thông tin địa chỉ nơi gửi hàng")
	}
	if r.ToDistrictID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu thông tin địa chỉ nơi nhận hàng")
	}
	if r.ClientContactName == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu tên người gửi hàng")
	}
	if r.ClientContactPhone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu số điện thoại nơi gửi hàng")
	}
	if r.ClientAddress == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu địa chỉ nơi gửi hàng")
	}
	if r.CustomerName == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu tên người nhận hàng")
	}
	if r.CustomerPhone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu số điện thoại người nhận hàng")
	}
	if r.ShippingAddress == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu địa chỉ nơi nhận hàng")
	}
	if r.ServiceID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với Giao Hàng Nhanh: Thiếu mã dịch vụ")
	}
	return nil
}

type CreateOrderResponse struct {
	ErrorMessage         String `json:"ErrorMessage"`         // "",
	OrderID              Int    `json:"OrderID"`              // 268916,
	PaymentTypeID        Int    `json:"PaymentTypeID"`        // 4,
	OrderCode            String `json:"OrderCode"`            // "236697NF",
	ExtraFee             Int    `json:"ExtraFee"`             // 0,
	TotalServiceFee      Int    `json:"TotalServiceFee"`      // 81400,
	ExpectedDeliveryTime Time   `json:"ExpectedDeliveryTime"` // "2017-09-22T23:00:00+07:00",
	ClientHubID          Int    `json:"ClientHubID"`          // 0,
	SortCode             String `json:"SortCode"`             // "N/A"
}

type ShippingOrderCostRequest struct {
	ServiceID int `json:"ServiceID"`
}

type OrderCodeRequest struct {
	Connection
	OrderCode string `json:"OrderCode"`
}

type Order struct {
	CODTransferDate      String `json:"CODTransferDate"`      // "",
	CSLostPackageID      Int    `json:"CSLostPackageID"`      // 0,
	CheckMainBankAccount Bool   `json:"CheckMainBankAccount"` // false,
	ClientHubID          Int    `json:"ClientHubID"`          // 2,
	CoDAmount            Int    `json:"CoDAmount"`            // 1500000,
	Content              String `json:"Content"`              // "Test nội dung",
	CouponCode           String `json:"CouponCode"`           // "",
	CurrentStatus        String `json:"CurrentStatus"`        // "Cancel",
	CurrentWarehouseName String `json:"CurrentWarehouseName"` // "",
	CustomerID           Int    `json:"CustomerID"`           // 0,
	CustomerName         String `json:"CustomerName"`         // "Nguyễn Văn UpdateInfo",
	CustomerPhone        String `json:"CustomerPhone"`        // "0934035687",
	DeliverWarehouseName String `json:"DeliverWarehouseName"` // "",
	EndDeliveryTime      String `json:"EndDeliveryTime"`      // "",
	EndPickTime          String `json:"EndPickTime"`          // "2017-09-14T15:23:56.643",
	EndReturnTime        String `json:"EndReturnTime"`        // "",
	ErrorMessage         String `json:"ErrorMessage"`         // "",
	ExpectedDeliveryTime String `json:"ExpectedDeliveryTime"` // "2017-09-18T23:00:00",
	ExternalCode         String `json:"ExternalCode"`         // "",
	ExternalReturnCode   String `json:"ExternalReturnCode"`   // "",
	FirstDeliveredTime   String `json:"FirstDeliveredTime"`   // "",
	FromDistrictID       Int    `json:"FromDistrictID"`       // 0,
	FromLat              Int    `json:"FromLat"`              // 0,
	FromLng              Int    `json:"FromLng"`              // 0,
	FromWardCode         String `json:"FromWardCode"`         // "",
	Height               Int    `json:"Height"`               // 10,
	InsuranceFee         Int    `json:"InsuranceFee"`         // 0,
	Length               Int    `json:"Length"`               // 10,
	Note                 String `json:"Note"`                 // "Lưu ĐH qua API ",
	NoteCode             String `json:"NoteCode"`             // "",
	NumDeliver           Int    `json:"NumDeliver"`           // 0,
	NumPick              Int    `json:"NumPick"`              // 0,
	OrderCode            String `json:"OrderCode"`            // "23HY9557",
	OriginServiceName    String `json:"OriginServiceName"`    // "1 Ngày",
	PaymentTypeID        Int    `json:"PaymentTypeID"`        // 1,
	PickWarehouseName    String `json:"PickWarehouseName"`    // "",
	ReturnInfo           String `json:"ReturnInfo"`           // "",
	SealCode             String `json:"SealCode"`             // "",
	ServiceID            Int    `json:"ServiceID"`            // 53320,
	ServiceName          String `json:"ServiceName"`          // "1 Ngày",
	ShippingAddress      String `json:"ShippingAddress"`      // "137 Lê Quang Định",
	ShippingOrderID      Int    `json:"ShippingOrderID"`      // 268263,
	StartDeliveryTime    String `json:"StartDeliveryTime"`    // "",
	StartPickTime        String `json:"StartPickTime"`        // "",
	StartReturnTime      String `json:"StartReturnTime"`      // "",
	ToDistrictID         Int    `json:"ToDistrictID"`         // 1462,
	ToLatitude           Int    `json:"ToLatitude"`           // 0,
	ToLongitude          Int    `json:"ToLongitude"`          // 0,
	ToWardCode           String `json:"ToWardCode"`           // "",
	TotalServiceCost     Int    `json:"TotalServiceCost"`     // 81400,
	Weight               Int    `json:"Weight"`               // 10200,
	Width                Int    `json:"Width"`                // 10

	ShippingOrderCosts []*ShippingOrderCost `json:"ShippingOrderCosts"` //

	//ExtraFees            String `json:"ExtraFees"`           // null,
}

type OrderLog struct {
	OrderCode       String       `json:"OrderCode"`       // "DHYDFYHF"
	ShippingOrderID Int          `json:"ShippingOrderID"` // 268263,
	CurrentStatus   String       `json:"CurrentStatus"`   // "Cancel",
	CustomerID      Int          `json:"CustomerID"`      // 0,
	IsPushed        Bool         `json:"IsPushed"`        // true
	OrderInfo       OrderLogInfo `json:"OrderInfo"`
	StatusCode      Int          `json:"StatusCode"` // 0
	CreateTime      Time         `json:"CreateTime"` // "2018-08-07T14:42:23.31+07:00"
	UpdateTime      Time         `json:"UpdateTime"`
}

type OrderLogInfo struct {
	CoDAmount            Int                  `json:"CoDAmount"`            // 1500000,
	CurrentStatus        String               `json:"CurrentStatus"`        // "Cancel",
	CurrentWarehouseName String               `json:"CurrentWarehouseName"` // "",
	CustomerID           Int                  `json:"CustomerID"`           // 0,
	CustomerName         String               `json:"CustomerName"`         // "Nguyễn Văn UpdateInfo",
	CustomerPhone        String               `json:"CustomerPhone"`        // "0934035687",
	ExternalCode         String               `json:"ExternalCode"`         // "",
	Note                 String               `json:"Note"`                 // "Lưu ĐH qua API ",
	OrderCode            String               `json:"OrderCode"`            // "23HY9557",
	ReturnInfo           String               `json:"ReturnInfo"`           // "",
	ServiceName          String               `json:"ServiceName"`          // "1 Ngày",
	ShippingOrderCosts   []*ShippingOrderCost `json:"ShippingOrderCosts"`   //
	Weight               Int                  `json:"Weight"`               // 10200,
}

// Use for webhook
type CallbackOrder struct {
	CoDAmount            Int                  `json:"CoDAmount"`            // 0,
	CurrentStatus        String               `json:"CurrentStatus"`        // "ReadyToPick",
	CurrentWarehouseName String               `json:"CurrentWarehouseName"` // "Kho giao nhận Đống Đa_Hà Nội",
	CustomerID           Int                  `json:"CustomerID"`           // 252905,
	CustomerName         String               `json:"CustomerName"`         // "Hà Anh",
	CustomerPhone        String               `json:"CustomerPhone"`        // "0973636049",
	ExternalCode         String               `json:"ExternalCode"`         // "",
	Note                 String               `json:"Note"`                 // "Gửi hàng",
	OrderCode            String               `json:"OrderCode"`            // "DC5D4NFUH",
	ReturnInfo           String               `json:"ReturnInfo"`           // "",
	ServiceName          String               `json:"ServiceName"`          // "Nhanh",
	Weight               Int                  `json:"Weight"`               // 800
	ShippingOrderCosts   []*ShippingOrderCost `json:"ShippingOrderCosts"`
}
type ShippingOrderCost struct {
	Cost             String `json:"Cost"`             // 0,
	PaymentChannelID String `json:"PaymentChannelID"` // 4,
	ServiceID        String `json:"ServiceID"`        // 53332,
	ServiceName      String `json:"ServiceName"`      // "Dịch vụ khai giá",
	ServiceType      String `json:"ServiceType"`      // 5,
	ShippingOrderID  String `json:"ShippingOrderID"`  // 268263
}

type FindAvailableServicesRequest struct {
	Connection

	Weight int `json:"Weight"` // 10000,
	Length int `json:"Length"` // 10,
	Width  int `json:"Width"`  // 110,
	Height int `json:"Height"` // 20,

	FromDistrictID int `json:"FromDistrictID"` // 1455,
	ToDistrictID   int `json:"ToDistrictID"`   // 1462
	InsuranceFee   int `json:"InsuranceFee"`
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

type OrderLogsRequest struct {
	Token     string              `json:"token"`
	FromTime  int64               `json:"FromTime"`
	ToTime    int64               `json:"ToTime"`
	Skip      int                 `json:"Skip"`
	Condition *OrderLogsCondition `json:"Condition"`
}

type OrderLogsResponse struct {
	Logs []*OrderLog `json:"Logs"`
}

type OrderLogsCondition struct {
	OrderCode  string `json:"OrderCode"`
	CustomerID int    `json:"CustomerID"`
}

func (info *OrderLogInfo) ToGHNOrder() *Order {
	return &Order{
		CoDAmount:            info.CoDAmount,
		CurrentStatus:        info.CurrentStatus,
		CurrentWarehouseName: info.CurrentWarehouseName,
		CustomerID:           info.CustomerID,
		CustomerName:         info.CustomerName,
		CustomerPhone:        info.CustomerName,
		ExternalCode:         info.ExternalCode,
		Note:                 info.Note,
		OrderCode:            info.OrderCode,
		ReturnInfo:           info.ReturnInfo,
		ServiceName:          info.ServiceName,
		Weight:               info.Weight,
		ShippingOrderCosts:   info.ShippingOrderCosts,
	}
}

func (ghnOrder *Order) ToFakeCallbackOrder() *CallbackOrder {
	return &CallbackOrder{
		CoDAmount:            ghnOrder.CoDAmount,
		CurrentStatus:        ghnOrder.CurrentStatus,
		CurrentWarehouseName: ghnOrder.CurrentWarehouseName,
		CustomerID:           ghnOrder.CustomerID,
		CustomerName:         ghnOrder.CustomerName,
		CustomerPhone:        ghnOrder.CustomerName,
		ExternalCode:         ghnOrder.ExternalCode,
		Note:                 ghnOrder.Note,
		OrderCode:            ghnOrder.OrderCode,
		ReturnInfo:           ghnOrder.ReturnInfo,
		ServiceName:          ghnOrder.ServiceName,
		Weight:               ghnOrder.Weight,
	}
}

func (info *OrderLogInfo) ToFakeCallbackOrder() *CallbackOrder {
	return &CallbackOrder{
		CoDAmount:            info.CoDAmount,
		CurrentStatus:        info.CurrentStatus,
		CurrentWarehouseName: info.CurrentWarehouseName,
		CustomerID:           info.CustomerID,
		CustomerName:         info.CustomerName,
		CustomerPhone:        info.CustomerName,
		ExternalCode:         info.ExternalCode,
		Note:                 info.Note,
		OrderCode:            info.OrderCode,
		ReturnInfo:           info.ReturnInfo,
		ServiceName:          info.ServiceName,
		Weight:               info.Weight,
	}
}

type FindAvailableServicesResponse struct {
	AvailableServices []*AvailableService
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}

type CalculateFeeRequest struct {
	Connection

	Weight int `json:"Weight"` // 10000,
	Length int `json:"Length"` // 10,
	Width  int `json:"Width"`  // 110,
	Height int `json:"Height"` // 20,

	FromDistrictID int `json:"FromDistrictID"` // 1455,
	ToDistrictID   int `json:"ToDistrictID"`   // 1462
	InsuranceFee   int `json:"InsuranceFee"`
	ServiceID      int `json:"ServiceID"`
}

type CalculateFeeResponse struct {
	ErrorMessage    string       `json:"ErrorMessage"`
	CalculatedFee   int          `json:"CalculatedFee"`
	ServiceFee      int          `json:"ServiceFee"`
	CODFee          int          `json:"CoDFee"`
	OrderCosts      []*OrderCost `json:"OrderCosts"`
	DiscountFee     int          `json:"DiscountFee"`
	WeightDimension int          `json:"WeightDimension"`
}

func (r *CalculateFeeResponse) ToShippingService(providerServiceID string) *shippingsharemodel.AvailableShippingService {
	if r == nil {
		return nil
	}
	return &shippingsharemodel.AvailableShippingService{
		Name:                "",
		ServiceFee:          r.CalculatedFee,
		ShippingFeeMain:     r.ServiceFee,
		Provider:            shipping_provider.GHN,
		ProviderServiceID:   providerServiceID,
		ExpectedPickAt:      time.Time{},
		ExpectedDeliveryAt:  time.Time{},
		Source:              "",
		ConnectionInfo:      nil,
		ShipmentServiceInfo: nil,
		ShipmentPriceInfo:   nil,
	}
}

type OrderCost struct {
	Cost      int    `json:"Cost"`
	Name      string `json:"Name"`
	ServiceID int    `json:"ServiceID"`
}

func GetInsuranceFee(orderCosts []*OrderCost) int {
	if orderCosts == nil {
		return 0
	}
	insuranceFee := 0
	for _, orderCost := range orderCosts {
		sID := ServiceFeeID(strconv.Itoa(orderCost.ServiceID))
		shippingType := sID.ToModel()
		if shippingType == shipping_fee_type.Insurance {
			insuranceFee = orderCost.Cost
			break
		}
	}
	return insuranceFee
}

type SignInRequest struct {
	Connection
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type SignInResponse struct {
	ClientID   Int    `json:"ClientID"`
	ClientName string `json:"ClientName"`
	Token      string `json:"Token"`
}

type SignUpRequest struct {
	Connection
	Email        string `json:"Email"`
	Password     string `json:"Password"`
	ContactPhone string `json:"ContactPhone"`
	ContactName  string `json:"ContactName"`
}

// document: https://api.ghn.vn/home/docs/detail?id=41
type RegisterWebhookForClientRequest struct {
	Connection
	TokenClient []string `json:"TokenClient"`
	// ConfigCod: Receive COD callback
	ConfigCOD        bool `json:"ConfigCod"`
	ConfigReturnData bool `json:"ConfigReturnData"`
	// URLCallback: the endpoint you receive data.
	URLCallback  string              `json:"URLCallback"`
	ConfigField  WebhookConfigField  `json:"ConfigField"`
	ConfigStatus WebhookConfigStatus `json:"ConfigStatus"`
}

type WebhookConfigField struct {
	CODAmount            bool `json:"CoDAmount"`
	CurrentWarehouseName bool `json:"CurrentWarehouseName"`
	CustomerID           bool `json:"CustomerID"`
	CustomerName         bool `json:"CustomerName"`
	CustomerPhone        bool `json:"CustomerPhone"`
	Note                 bool `json:"Note"`
	OrderCode            bool `json:"OrderCode"`
	ServiceName          bool `json:"ServiceName"`
	ShippingOrderCosts   bool `json:"ShippingOrderCosts"`
	Weight               bool `json:"Weight"`
	ExternalCode         bool `json:"ExternalCode"`
	ReturnInfo           bool `json:"ReturnInfo"`
}

type WebhookConfigStatus struct {
	ReadyToPick     bool `json:"ReadyToPick"`
	Picking         bool `json:"Picking"`
	Storing         bool `json:"Storing"`
	Delivering      bool `json:"Delivering"`
	Delivered       bool `json:"Delivered"`
	WaitingToFinish bool `json:"WaitingToFinish"`
	Return          bool `json:"Return"`
	Returned        bool `json:"Returned"`
	Finish          bool `json:"Finish"`
	LostOrder       bool `json:"LostOrder"`
	Cancel          bool `json:"Cancel"`
}
