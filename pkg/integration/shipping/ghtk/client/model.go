package client

import (
	"time"

	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status5"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpreq"
	"etop.vn/backend/pkg/etop/model"
)

type TransportType string

const (
	TransportRoad TransportType = "road"
	TransportFly  TransportType = "fly"
)

func (t TransportType) Name() string {
	switch t {
	case TransportRoad:
		return model.ShippingServiceNameStandard
	case TransportFly:
		return model.ShippingServiceNameFaster
	}
	return string(t)
}

type GhtkAccount struct {
	AccountID   string `yaml:"account_id"`
	Token       string `yaml:"token"`
	AffiliateID string `yaml:"-"`
	B2CToken    string `yaml:"-"`
}

type StateID int

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

const (
	// reference link: https://docs.giaohangtietkiem.vn/?http#tr-ng-th-i-n-h-ng22
	StateIDCanceled             StateID = -1 // Hủy đơn hàng
	StateIDNotConfirm           StateID = 1  // Chưa tiếp nhận
	StateIDConfirmed            StateID = 2  // Đã tiếp nhận
	StateIDStored               StateID = 3  // Đã lấy hàng/Đã nhập kho
	StateIDDelivering           StateID = 4  // Đã điều phối giao hàng/Đang giao hàng
	StateIDDelivered            StateID = 5  // Đã giao hàng/Chưa đối soát
	StateIDCrossChecked         StateID = 6  // Đã đối soát
	StateIDNotTake              StateID = 7  // Không lấy được hàng
	StateIDDelayPicking         StateID = 8  // Hoãn lấy hàng
	StateIDNotDeliver           StateID = 9  // Không giao được hàng
	StateIDDelayDelivering      StateID = 10 // Delay giao hàng
	StateIDCrossCheckedReturned StateID = 11 // Đã đối soát công nợ trả hàng
	StateIDPicking              StateID = 12 // Đã điều phối lấy hàng/Đang lấy hàng
	StateIDBoiHoan              StateID = 13 // Đơn hàng bồi hoàn
	StateIDReturning            StateID = 20 // Đang trả hàng (COD cầm hàng đi trả)
	StateIDReturned             StateID = 21 // Đã trả hàng (COD đã trả xong hàng)

	// Trạng thái phụ (dành cho Shipper) - trạng thái này chỉ được trả về trong webhook, khi lấy chi tiết đơn hàng nó sẽ trả về trạng thái chính
	StateIDShipperPicked          StateID = 123 // Shipper báo đã lấy hàng
	StateIDShipperCanNotPick      StateID = 127 // Shipper (nhân viên lấy/giao hàng) báo không lấy được hàng
	StateIDShipperDelayPicking    StateID = 128 // Shipper báo delay lấy hàng
	StateIDShipperDelivered       StateID = 45  // Shipper báo đã giao hàng
	StateIDShipperCanNotDelivery  StateID = 49  // Shipper báo không giao được giao hàng
	StateIDShipperDelayDelivering StateID = 410 // Shipper báo delay giao hàng
)

var StateMapping = map[StateID]string{
	StateIDCanceled:             "Hủy đơn hàng",
	StateIDNotConfirm:           "Chưa tiếp nhận",
	StateIDConfirmed:            "Đã tiếp nhận",
	StateIDStored:               "Đã lấy hàng/Đã nhập kho",
	StateIDDelivering:           "Đã điều phối giao hàng/Đang giao hàng",
	StateIDDelivered:            "Đã giao hàng/Chưa đối soát",
	StateIDCrossChecked:         "Đã đối soát",
	StateIDNotTake:              "Không lấy được hàng",
	StateIDDelayPicking:         "Hoãn lấy hàng",
	StateIDNotDeliver:           "Không giao được hàng",
	StateIDDelayDelivering:      "Delay giao hàng",
	StateIDCrossCheckedReturned: "Đã đối soát công nợ trả hàng",
	StateIDPicking:              "Đã điều phối lấy hàng/Đang lấy hàng",
	StateIDBoiHoan:              "Đơn hàng bồi hoàn",
	StateIDReturning:            "Đang trả hàng (COD cầm hàng đi trả)",
	StateIDReturned:             "Đã trả hàng (COD đã trả xong hàng)",
}

var SubStateMapping = map[StateID]string{
	StateIDShipperPicked:          "Shipper báo đã lấy hàng",
	StateIDShipperCanNotPick:      "Shipper báo không lấy được hàng",
	StateIDShipperDelayPicking:    "Shipper báo delay lấy hàng",
	StateIDShipperDelivered:       "Shipper báo đã giao hàng",
	StateIDShipperCanNotDelivery:  "Shipper báo không giao được giao hàng",
	StateIDShipperDelayDelivering: "Shipper báo delay giao hàng",
}

func (sID StateID) ToModel() shipping.State {
	switch sID {
	case StateIDCanceled, StateIDNotTake:
		return shipping.Cancelled
	case StateIDNotConfirm, StateIDConfirmed:
		return shipping.Created
	case StateIDPicking, StateIDDelayPicking, StateIDShipperCanNotPick, StateIDShipperDelayPicking:
		return shipping.Picking
	case StateIDStored, StateIDShipperPicked:
		return shipping.Holding
	case StateIDDelivering, StateIDDelayDelivering, StateIDShipperCanNotDelivery, StateIDShipperDelayDelivering:
		return shipping.Delivering
	case StateIDDelivered, StateIDCrossChecked, StateIDShipperDelivered:
		return shipping.Delivered
	case StateIDReturning, StateIDNotDeliver, StateIDCrossCheckedReturned:
		return shipping.Returning
	case StateIDReturned:
		return shipping.Returned
	case StateIDBoiHoan:
		return shipping.Undeliverable
	default:
		return shipping.Unknown
	}
}

func (sID StateID) ToStatus5() status5.Status {
	switch sID {
	case StateIDCanceled, StateIDNotTake:
		return status5.N
	case StateIDReturned, StateIDBoiHoan:
		return status5.NS
	case StateIDCrossChecked, StateIDDelivered:
		return status5.P
	default:
		return status5.S
	}
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

type ShippingFee struct {
	Name         String `json:"name"`
	Fee          Int    `json:"fee"`
	InsuranceFee Int    `json:"insurance_fee"`
	// Delivery - Boolean - Hỗ trợ giao ở địa chỉ này chưa, nếu điểm giao đã được GHTK hỗ trợ giao trả về true, nếu GTHK chưa hỗ trợ giao đến khu vực này thì trả về false
	Delivery     Bool   `json:"delivery"`
	DeliveryType String `json:"delivery_type"`
}

func (s *ShippingFee) ToShippingService(providerServiceID string, transport TransportType, expectedPickTime, expectedDeliveryTime time.Time) *model.AvailableShippingService {
	if s == nil {
		return nil
	}
	return &model.AvailableShippingService{
		// GHTK always returns the same name, so we replace it with our name
		// (Chuẩn, Nhanh)
		Name:              transport.Name(),
		ServiceFee:        int(s.Fee),
		ShippingFeeMain:   int(s.Fee - s.InsuranceFee),
		Provider:          shipping_provider.GHTK,
		ProviderServiceID: providerServiceID,

		ExpectedPickAt:     expectedPickTime,
		ExpectedDeliveryAt: expectedDeliveryTime,
	}
}

type CalcShippingFeeRequest struct {
	Weight int `url:"weight"` // Integer - Cân nặng của gói hàng, đơn vị sử dụng Gram
	Value  int `url:"value"`  // Integer - Giá trị thực của đơn hàng áp dụng để tính phí bảo hiểm, đơn vị sử dụng VNĐ

	PickingAddress  string `url:"pick_address"`  // Địa chỉ ngắn gọn để lấy nhận hàng hóa. Ví dụ: nhà số 5, tổ 3, ngách 11, ngõ 45
	PickingProvince string `url:"pick_province"` // required - Tên tỉnh/thành phố nơi lấy hàng hóa
	PickingDistrict string `url:"pick_district"` // required - Tên quận/huyện nơi lấy hàng hóa
	PickingWard     string `url:"picking_ward"`  // Tên phường/xã nơi lấy hàng hóa
	Address         string `url:"address"`       // Địa chỉ chi tiết của người nhận hàng
	Province        string `url:"province"`      // required - Tên tỉnh/thành phố của người nhận hàng hóa
	District        string `url:"district"`      // required - Tên quận/huyện của người nhận hàng hóa
	Ward            string `url:"ward"`          // Tên phường/xã của người nhận hàng hóa

	Transport TransportType `url:"transport"` // fly or road (road is default)
}

type CalcShippingFeeResponse struct {
	CommonResponse
	Fee *ShippingFee `json:"fee"`

	ErrorData map[string]string `json:"-"`
}

type CreateOrderRequest struct {
	Products []*ProductRequest `json:"products"`
	Order    *OrderRequest     `json:"order"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.Order.PickName == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên người liên hệ lấy hàng")
	}
	if r.Order.PickAddress == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu địa chỉ để lấy hàng")
	}
	if r.Order.PickProvince == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên tỉnh/thành phố nơi lấy hàng")
	}
	if r.Order.PickDistrict == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên quận/huyện nơi lấy hàng")
	}
	if r.Order.PickTel == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu số điện thoại liên hệ nơi lấy hàng")
	}
	if r.Order.Name == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên người nhận hàng")
	}
	if r.Order.Address == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu địa chỉ nhận hàng")
	}
	if r.Order.Province == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên tỉnh/thành phố nơi nhận hàng")
	}
	if r.Order.District == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu tên quận/huyện nơi nhận hàng")
	}
	if r.Order.Tel == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với GHTK: Thiếu số điện thoại liên hệ nơi nhận hàng")
	}
	return nil
}

type ProductRequest struct {
	Name     string  `json:"name"`     // required - Tên hàng hóa
	Weight   float32 `json:"weight"`   // required - Khối lượng hàng hóa Tính theo đơn vị KG
	Price    int     `json:"price"`    // Giá trị hàng hóa
	Quantity int     `json:"quantity"` // Số lượng hàng hóa
}

type OrderRequest struct {
	ID string `json:"id"` // required - mã đơn hàng thuộc hệ thống của đối tác
	// Thông tin điểm lấy hàng
	PickName     string `json:"pick_name"`     // required - Tên người liên hệ lấy hàng hóa
	PickMoney    int    `json:"pick_money"`    // required - Số tiền cần thu hộ. Nếu bằng 0 thì không thu hộ tiền. Tính theo VNĐ
	PickAddress  string `json:"pick_address"`  // required - Địa chỉ ngắn gọn để lấy nhận hàng hóa
	PickProvince string `json:"pick_province"` // required - Tên tỉnh/thành phố nơi lấy hàng hóa
	PickDistrict string `json:"pick_district"` // required - Tên quận/huyện nơi lấy hàng hóa
	PickWard     string `json:"pick_ward"`     // Tên phường/xã nơi lấy hàng hóa
	PickStreet   string `json:"pick_street"`   // Tên đường/phố nơi lấy hàng hóa
	PickTel      string `json:"pick_tel"`      // required - Số điện thoại liên hệ nơi lấy hàng hóa
	PickEmail    string `json:"pick_email"`
	// Thông tin điểm giao hàng
	Name     string `json:"name"`     // required - tên người nhận hàng
	Address  string `json:"address"`  // required - Địa chỉ chi tiết của người nhận hàng, vd: Chung cư CT1, ngõ 58, đường Trần Bình
	Province string `json:"province"` // required - Tên tỉnh/thành phố của người nhận hàng hóa
	District string `json:"district"` // required - Tên quận/huyện của người nhận hàng hóa
	Ward     string `json:"ward"`     // Tên phường/xã của người nhận hàng hóa
	Street   string `json:"street"`   // Tên đường/phố của người nhận hàng hóa
	Tel      string `json:"tel"`      // required - Số điện thoại người nhận hàng hóa
	Note     string `json:"note"`
	Email    string `json:"string"` // required - Email người nhận hàng hóa
	// Thông tin điểm trả hàng
	UseReturnAddress int `json:"use_return_address"` // mặc định là 0. Field này có thể truyền vào một trong hai giá trị 0 hoặc 1.
	// Bằng 0 nghĩa là địa chỉ trả hàng giống địa chỉ lấy hàng nên các field địa chỉ trả hàng không cần truyền qua.
	// Bằng 1 nghĩa là sử dụng địa chỉ trả hàng khác địa chỉ lấy hàng và cần truyền vào giá trị cho các field tiếp theo
	ReturnName     string `json:"return_name"`     // required
	ReturnAddress  string `json:"return_address"`  // required
	ReturnProvince string `json:"return_province"` // required
	ReturnDistrict string `json:"return_district"` // required
	ReturnWard     string `json:"return_ward"`
	ReturnStreet   string `json:"return_street"`
	ReturnTel      string `json:"return_tel"`   // required - Số điện thoại người nhận hàng hóa
	ReturnEmail    string `json:"return_email"` // required - Email người nhận hàng hóa

	// Thông tin thêm
	IsFreeship       int     `json:"is_freeship"`        // Freeship cho người nhận hàng. Nếu bằng 1 COD sẽ chỉ thu người nhận hàng số tiền bằng pick_money, nếu bằng 0 COD sẽ thu tiền người nhận số tiền bằng pick_money + phí ship của đơn hàng, giá trị mặc định bằng 0
	WeightOption     string  `json:"weight_option"`      // nhận một trong hai giá trị gram và kilogram, mặc định là kilogram, đơn vị khối lượng của các sản phẩm có trong gói hàng
	DeliverWorkShift int     `json:"deliver_work_shift"` // Nếu set bằng 3 đơn hàng sẽ được giao vào buổi tối. Giá trị mặc định bằng 0
	Value            int     `json:"value"`              // Giá trị đóng bảo hiểm, là căn cứ để tính phí bảo hiểm và bồi thường khi có sự cố.
	PickOption       string  `json:"pick_option"`        // Nhận một trong hai giá trị cod và post, mặc định là cod, biểu thị lấy hàng bởi COD hoặc Shop sẽ gửi tại bưu cục
	TotalWeight      float32 `json:"total_weight"`       // Tổng khối lượng của đơn hàng, mặc định sẽ tính theo products.weight nếu không truyền giá trị này.

	Transport TransportType `json:"transport"` // fly or road (road is default)
}

type CreateOrderResponse struct {
	CommonResponse
	Order OrderResponse `json:"order"`
}

type OrderResponse struct {
	PartnerID            String `json:"partner_id"`
	Label                String `json:"label"`
	Area                 String `json:"area"`
	Fee                  Int    `json:"fee"`
	InsuranceFee         Int    `json:"insurance_fee"`
	EstimatedPickTime    String `json:"estimated_pick_time"` // "Sáng 2017-07-01"
	EstimatedDeliverTime String `json:"estimated_deliver_time"`
	StatusID             Int    `json:"status_id"`
	TrackingID           Int    `json:"tracking_id"`
}

type OrderInfo struct {
	LabelID          String `json:"label_id"`   // Mã đơn hàng của hệ thống GHTK
	PartnerID        String `json:"partner_id"` // Mã đơn hàng thuộc hệ thống của đối tác
	Status           Int    `json:"status"`
	StatusText       String `json:"status_text"`
	Created          String `json:"created"` // Thời gian tạo đơn hàng, định dạng YY-MM-DD hh:mm:ss
	Modified         String `json:"modified"`
	Message          String `json:"message"`      // Ghi chú của đơn hàng
	PickDate         String `json:"pick_date"`    //  Ngày hẹn lấy hàng của đơn hàng nếu có, nếu đơn hàng đã được lấy thành công thì là ngày lấy hàng
	DeliverDate      String `json:"deliver_date"` // Ngày hẹn giao đơn hàng nếu có, nếu đơn hàng đã được giao hàng thì là ngày giao hàng thành công
	CustomerFullname String `json:"customer_fullname"`
	CustomerTel      String `json:"customer_tel"`
	Address          String `json:"address"`
	StorageDay       Int    `json:"storage_day"` // Số ngày giữ đơn hàng tại kho GHTK trước khi trả hàng
	ShipMoney        Int    `json:"ship_money"`  // Phí giao hàng
	Insurance        Int    `json:"insurance"`   // Phí bảo hiểm
	Value            Int    `json:"value"`       //  Giá trị đóng bảo hiểm - căn cứ để bồi thường cho người gửi khi có sự cố xảy ra
	Weight           Int    `json:"weight"`      //  Khối lượng đơn hàng tính theo gram
	PickMoney        Int    `json:"pick_money"`  // Số tiền cần thu hộ
	IsFreeship       Int    `json:"is_freeship"` // Freeship cho người nhận hàng
}

type GetOrderResponse struct {
	CommonResponse
	Order OrderInfo `json:"order"`
}

type CallbackOrder struct {
	PartnerID  String `json:"partner_id"`
	LabelID    String `json:"label_id"`
	StatusID   Int    `json:"status_id"`
	ActionTime String `json:"action_time"`
	ReasonCode String `json:"reason_code"` // mã lý do cập nhật
	Reason     String `json:"reason"`      // Lý do chi tiết cập nhật
	Weight     Float  `json:"weight"`      // khối lượng đơn hàng tính theo kilogram
	Fee        Int    `json:"fee"`
}

type SignUpRequest struct {
	// required
	Name string `json:"name"`
	// required: Địa chỉ chi tiết của tài khoản (số nhà, ngõ, đường, phường,…)
	FirstAddress string `json:"first_address"`
	// required
	Province string `json:"province"`
	// required
	District string `json:"district"`
	// required
	Tel string `json:"tel"`
	// required
	Email string `json:"email"`
}

type SignUpResponse struct {
	CommonResponse
	Data *AccountData `json:"data"`
}

type AccountData struct {
	// Mã tài khoản trên hệ thống GHTK
	Code  string `json:"code"`
	Token string `json:"token"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	CommonResponse
	Data *AccountData `json:"data"`
}
