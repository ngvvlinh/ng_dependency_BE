package client

import (
	typesshipping "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/pkg/common/apifw/httpreq"
	cc "o.o/backend/pkg/common/config"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type Config struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	PartnerID int    `yaml:"partner_id"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_NTX"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_USERNAME":  &c.Username,
		p + "_PASSWORD":  &c.Password,
		p + "_PARTNERID": &c.PartnerID,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Username:  "etop",
		Password:  "hm73Tyrk9EzNnhSJ",
		PartnerID: 350,
	}
}

type State string

const (
	// ID: 1
	// Mã trạng thái : OPS
	// Vận đơn được khởi tạo qua hệ thống NTX
	StateStaging State = "Vận đơn mới"

	// ID: 2
	// Mã trạng thái : PRI
	// Vận đơn được NTX xác nhận lấy hàng
	StateConfirmed State = "Đã xác nhận"

	// ID: 3
	// Mã trạng thái : LPC
	// Vận đơn được NTX xác nhận lấy hàng
	StateSuccessfulPickup State = "Đã lấy hàng"

	// ID: 4
	// Mã trạng thái : EOD
	// Vận đơn được xuất kho trả hàng và bắt đầu giao đến Người nhận
	StateInProgressForDelivery State = "Đang giao hàng"

	// ID: 5
	// Mã trạng thái : FBC
	// Vận đơn đã được giao đến Người nhận thành công
	StateSuccessfulDelivery State = "Giao thành công"

	// ID: 6
	// Mã trạng thái : GBV
	// Vận đã đã được Hủy từ người gửi hoặc Bộ phận vận hành NTX khi không liên lạc được với Người gửi
	StateCancelled State = "Hủy"

	// ID: 7
	// Mã trạng thái : MRC
	// Vận đơn đi giao nhưng không liên hệ hoặc sự cố giao hàng. NTX lập đơn chuyển hoàn và Người nhận hàng hoàn
	StateReturnedToSender State = "Đã chuyển hoàn"

	// ID: 8
	// Mã trạng thái : NRT
	// Vận đơn đi giao nhưng không liên hệ hoặc sự cố giao hàng. NTX lập đơn chuyển hoàn và Người nhận hàng hoàn
	StateReturnToSender State = "Đang chuyển hoàn"

	// ID: 9
	// Mã trạng thái : QIU
	// Vận đơn khi đang giao hàng không liên hệ được Người nhận hoặc gặp sự cố giao quá 3 lần
	StateDeliveryFailed State = "Không giao được"

	// ID: 10
	// Mã trạng thái : DIT
	// Vận đơn đang luân chuyển giữa các kho NTX
	StateTranshipment State = "Đang luân chuyển"

	// ID: 11
	// Mã trạng thái : OEP
	// Vận đơn gặp sự cố lấy hàng quá 3 lần
	StateUnsuccessfulPickup State = "Không lấy được"

	// ID: 12
	// Mã trạng thái : OEP
	// Vận đơn mới khi Nhân viên NTX xác nhận và đến lấy hàng không được
	StatePickupError State = "Lỗi lấy hàng"

	// ID: 13
	// Mã trạng thái : FUD
	// Vận đơn được giao đến người nhận nhưng gặp sự cố giao
	StateDeliveryDelay = "Delay giao hàng"

	// ID: 14
	// Mã trạng thái : CIW
	// Vận đơn đã được lấy và nhập về kho
	StateEnterTheWarehouse = "Nhập kho"

	// ID: 15
	// Mã trạng thái : WFA
	// Vận đơn chờ duyệt hoàn từ khách. Đối với trạng thái này, sau 48h Khách không có thao tác gì, hệ thống sẽ tự động chuyển sang Đã duyệt hoàn
	StateWaitForApproval = "Chờ duyệt hoàn"

	// ID: 16
	// Mã trạng thái : QIU
	// Vận đơn bị sự cố giao, sẽ giao lại trong phiên giao hàng sau
	StateWaitingForDelivery = "Chờ giao lại"

	// ID: 17
	// Mã trạng thái : RFA
	// Vận đơn được duyệt hoàn từ khách. Tiền hành lập hoàn
	StateFullyApproved = "Đã duyệt hoàn"
)

func (s State) ToModel() typesshipping.State {
	switch s {
	case StateStaging:
		return typesshipping.Created
	case StateConfirmed:
		return typesshipping.Confirmed
	case StateSuccessfulPickup, StateTranshipment, StateEnterTheWarehouse, StateWaitForApproval, StateFullyApproved:
		return typesshipping.Holding
	case StateInProgressForDelivery, StateDeliveryFailed, StateDeliveryDelay, StateWaitingForDelivery:
		return typesshipping.Delivering
	case StateSuccessfulDelivery:
		return typesshipping.Delivered
	case StateCancelled:
		return typesshipping.Cancelled
	case StateReturnToSender:
		return typesshipping.Returning
	case StateReturnedToSender:
		return typesshipping.Returned
	case StateUnsuccessfulPickup, StatePickupError:
		return typesshipping.Picking
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

type ResponseInterface interface {
	GetCommonResponse() CommonResponse
}

type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c CommonResponse) Error() string {
	if !c.Success {
		return c.Message
	}
	return "Lỗi không xác định: " + c.Message
}

func (c CommonResponse) GetCommonResponse() CommonResponse {
	return c
}

type CalcShippingFeeRequest struct {
	PartnerID     int     `json:"partner_id"`
	CodAmount     int     `json:"cod_amount"`
	CargoValue    int     `json:"cargo_value"`
	Weight        float64 `json:"weight"`
	PaymentMethod int     `json:"payment_method"`
	SProvinceID   int     `json:"s_province_id"`
	SDistrictID   int     `json:"s_district_id"`
	RProvinceID   int     `json:"r_province_id"`
	RDistrictID   int     `json:"r_district_id"`
	PackageNo     int     `json:"package_no"`
	UtmSource     string  `json:"utm_source"`
}

type CalcShippingFeeResponse struct {
	CommonResponse
	Data []*ShippingFeeData `json:"data"`
}

type ShippingFeeData struct {
	Weight      float64 `json:"weight"`
	TotalFee    int     `json:"total_fee"`
	MainFee     int     `json:"main_fee"`
	ServiceID   int     `json:"service_id"`
	ServiceCode string  `json:"service_code"`
	ServiceName string  `json:"service_name"`
	LeadTime    string  `json:"lead_time"`
}

type FindAvailableServicesResponse struct {
	AvailableServices []*AvailableService
}

type AvailableService struct {
	Name        String
	ServiceCode String
}

type CreateOrderRequest struct {
	PartnerID      int     `json:"partner_id"`
	SName          string  `json:"s_name"`
	SPhone         string  `json:"s_phone"`
	SAddress       string  `json:"s_address"`
	SProvinceID    int     `json:"s_province_id"`
	SDistrictID    int     `json:"s_district_id"`
	SWardID        int     `json:"s_ward_id"`
	RName          string  `json:"r_name"`
	RPhone         string  `json:"r_phone"`
	RAddress       string  `json:"r_address"`
	RProvinceID    int     `json:"r_province_id"`
	RDistrictID    int     `json:"r_district_id"`
	RWardID        int     `json:"r_ward_id"`
	CodAmount      int     `json:"cod_amount"`
	ServiceID      int     `json:"service_id"`
	PaymentMethod  int     `json:"payment_method"`
	Weight         float64 `json:"weight"`
	CargoContentID int     `json:"cargo_content_id"`
	CargoContent   string  `json:"cargo_content"`
	CargoValue     int     `json:"cargo_value"`
	Note           string  `json:"note"`
	UtmSource      string  `json:"utm_source"`
	RefCode        string  `json:"ref_code"`
	PackageNo      int     `json:"package_no"`
	IsDeductCod    int     `json:"is_deduct_cod"`
}

type CreateOrderResponse struct {
	CommonResponse
	Order OrderResponse `json:"order"`
}

type OrderResponse struct {
	BillID            int     `json:"bill_id"`
	BillCode          string  `json:"bill_code"`
	RefCode           string  `json:"ref_code"`
	StatusID          int     `json:"status_id"`
	StatusName        string  `json:"status_name"`
	CodAmount         int     `json:"cod_amount"`
	ServiceName       string  `json:"service_name"`
	PaymentMethodCode string  `json:"payment_method_code"`
	PaymentMethod     string  `json:"payment_method"`
	CreatedAt         string  `json:"created_at"`
	TotalFee          int     `json:"total_fee"`
	CodFee            int     `json:"cod_fee"`
	SName             string  `json:"s_name"`
	SPhone            string  `json:"s_phone"`
	SAddress          string  `json:"s_address"`
	SProvinceName     string  `json:"s_province_name"`
	SDistrictName     string  `json:"s_district_name"`
	SWardName         string  `json:"s_ward_name"`
	SPostName         string  `json:"s_post_name"`
	RName             string  `json:"r_name"`
	RPhone            string  `json:"r_phone"`
	RAddress          string  `json:"r_address"`
	RProvinceName     string  `json:"r_province_name"`
	RDistrictName     string  `json:"r_district_name"`
	RWardName         string  `json:"r_ward_name"`
	RPostName         string  `json:"r_post_name"`
	PackageNo         int     `json:"package_no"`
	Weight            float64 `json:"weight"`
	CargoContent      string  `json:"cargo_content"`
	CargoName         string  `json:"cargo_name"`
	CargoValue        int     `json:"cargo_value"`
	Note              string  `json:"note"`
}
