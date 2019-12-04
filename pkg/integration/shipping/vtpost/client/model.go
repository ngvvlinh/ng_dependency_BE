package client

import (
	"encoding/json"
	"math"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
)

type VTPostOrderServiceCode string
type VTPostShippingFeeType string
type VTPostProductType string

const (
	// Nhanh - SCOD Giao hàng thu tiền
	OrderServiceCodeSCOD VTPostOrderServiceCode = "SCOD"
	// Nhanh - VCN Chuyển phát nhanh - Express dilivery
	OrderServiceCodeVCN VTPostOrderServiceCode = "VCN"
	// Chậm - VTK - VTK Tiết kiệm - Express Saver
	OrderServiceCodeVTK VTPostOrderServiceCode = "VTK"
	// Chậm - PHS Phát hôm sau nội tỉnh
	OrderServiceCodePHS VTPostOrderServiceCode = "PHS"
	// Không dùng - V60 Dịch vụ Nhanh 60h
	OrderServiceCodeV60 VTPostOrderServiceCode = "V60"
	// Chậm - VVT Dịch vụ vận tải
	OrderServiceCodeVVT VTPostOrderServiceCode = "VVT"
	// Nhanh - VHT Phát Hỏa tốc
	OrderServiceCodeVHT VTPostOrderServiceCode = "VHT"
	// Nhanh - PTN Phát trong ngày nội tỉnh
	OrderServiceCodePTN VTPostOrderServiceCode = "PTN"
	// Nhanh - PHT Phát hỏa tốc nội tỉnh
	OrderServiceCodePHT VTPostOrderServiceCode = "PHT"
	// Nhanh - VBS Nhanh theo hộp
	OrderServiceCodeVBS VTPostOrderServiceCode = "VBS"
	// Chậm - VBE Tiết kiệm theo hộp
	OrderServiceCodeVBE VTPostOrderServiceCode = "VBE"

	OrderServiceCodeInsurance VTPostOrderServiceCode = "GBH"

	ShippingFeeTypeAll VTPostShippingFeeType = "ALL"
	ShippingFeeTypeVAT VTPostShippingFeeType = "VAT"
	// Thu hộ SCOD
	ShippingFeeTypeCODS VTPostShippingFeeType = "CODS"
	// Phí xăng dầu
	ShippingFeeTypePXD VTPostShippingFeeType = "PXD"
	// SCOD Giao hàng thu tiền
	ShippingFeeTypeSCOD VTPostShippingFeeType = "SCOD"
	// COD Phát hàng thu tiền
	ShippingFeeTypeCOD VTPostShippingFeeType = "COD"
	// GHB: Phí bảo hiểm
	ShippingFeeTypeGHB VTPostShippingFeeType = "GBH"
	// VTK Tiết kiệm
	ShippingFeeTypeVTK VTPostShippingFeeType = "VTK"
	// VCN Chuyển phát nhanh
	ShippingFeeTypeVCN VTPostShippingFeeType = "VCN"
	// PHS Phát hôm sau nội tỉnh
	ShippingFeeTypePHS VTPostShippingFeeType = "PHS"
	// Type: Hang hoa
	ProductTypeHH VTPostProductType = "HH"
	// Type: Thu
	ProductTypeTH VTPostProductType = "TH"

	IncludeVAT = 1.1
)

func ToShippingType(s string) VTPostShippingFeeType {
	return VTPostShippingFeeType(s)
}

func (s VTPostOrderServiceCode) Name() string {
	switch s {
	case OrderServiceCodeVTK, OrderServiceCodePHS, OrderServiceCodeVBE, OrderServiceCodeVVT:
		return model.ShippingServiceNameStandard
	case OrderServiceCodeSCOD, OrderServiceCodeVCN, OrderServiceCodeVHT, OrderServiceCodePTN,
		OrderServiceCodePHT, OrderServiceCodeVBS, OrderServiceCodeV60:
		return model.ShippingServiceNameFaster
	}
	return string(s)
}

type StateShipping string

var (
	// reference link: https://docs.google.com/document/d/1FhQE45QC0kkzkJFaoom9Oz5VJfiC-ZzhCzomqw84bjQ/edit#
	StateCanceled   StateShipping = "Hủy đơn hàng"  // []int{107, 201, 503, 101}
	StateNotConfirm StateShipping = "Chưa duyệt"    // []int{-100}
	StateConfirmed  StateShipping = "Đã duyệt"      // []int{100, -108}
	StatePicking    StateShipping = "Đang lấy hàng" // []int{102, 103, 104}
	// -109, -110: Đã gửi tại cửa hàng tiện ích
	// 200 ,202 ,300 ,320, 400: Đang vận chuyển
	StateStored     StateShipping = "Đang vận chuyển" // []int{-109, -110, 200, 202, 300, 320, 400, 301}
	StatePicked     StateShipping = "Đã lấy hàng"     // []int{105}
	StateDelivering StateShipping = "Đang giao hàng"  // []int{500, 506, 570, 508, 509, 550}
	// StateUndeliverable             StateShipping = "Giao hàng thất bại"   // []int{507}
	StateReturning StateShipping = "Duyệt hoàn"      // []int{505, 502, 515}
	StateReturned  StateShipping = "Hoàn thành công" // []int{504}
	// StateWaitingToConfirmReturning StateShipping = "Chờ duyệt hoàn"       // []int{505}
	StateDelivered StateShipping = "Giao hàng thành công" // []int{501}
)

var StateCodeMap = map[StateShipping][]int{
	StateCanceled:   {107, 201, 503, 101},
	StateNotConfirm: {-100},
	StateConfirmed:  {100, -108},
	StatePicking:    {102, 103, 104},
	StateStored:     {-109, -110, 200, 202, 300, 320, 400, 301},
	StatePicked:     {105},
	StateDelivering: {500, 505, 506, 507, 570, 508, 509, 550},
	StateReturning:  {502, 515},
	StateReturned:   {504},
	// StateWaitingToConfirmReturning: {505},
	StateDelivered: {501},
}

var SubStateMap = map[int]string{
	100:  "Tiếp nhận đơn hàng từ đối tác",
	101:  "ViettelPost yêu cầu hủy đơn hàng",
	102:  "Giải trình không liên hệ được người gửi hoặc không có hàng",
	103:  "Giao bưu cục đi lấy hàng",
	104:  "Giao bưu tá đi lấy hàng",
	-108: "Đơn hàng người gửi mang tới bưu cục",
	105:  "Buu Tá đã nhận hàng",
	106:  "Đối tác yêu cầu lấy lại hàng",
	107:  "Đối tác yêu cầu hủy qua API",
	200:  "Nhận từ bưu tá - Bưu cục gốc",
	201:  "Hủy nhập phiếu gửi",
	202:  "Sửa phiếu gửi",
	300:  "Đóng bảng kê đi",
	301:  "Ðóng túi gói",
	302:  "Đóng Chuyến thư",
	303:  "Đóng tuyến xe",
	400:  "Nhận bảng kê đến",
	401:  "Nhận Túi gói",
	402:  "Nhận chuyến thư",
	403:  "Nhận chuyến xe",
	500:  "Giao bưu tá đi phát",
	501:  "Thành công - Phát thành công",
	502:  "Chuyển hoàn bưu cục gốc",
	503:  "Hủy - Theo yêu cầu khách hàng",
	504:  "Thành công - Chuyển trả người gửi",
	// 505:  "Tồn - Thông báo chuyển hoàn bưu cục gốc",
	505: "Giao hàng không thành công vì không liên lạc được người nhận (Cần xác nhận lại số điện thoại, địa chỉ, tên ....)",
	// 506:  "Tồn - Khách hàng nghỉ, không có nhà",
	506: "Người nhận hẹn giao lại (Đi vắng, không tiện nhận hàng ...)",
	507: "Tồn - Khách hàng đến bưu cục nhận",
	508: "Phát tiếp",
	509: "Chuyển tiếp bưu cục khác",
	510: "Hủy phân công phát",
	515: "Bưu cục phát duyệt hoàn",
	550: "Đơn Vị Yêu Cầu Phát Tiếp",
}

func ToVTPostShippingState(code int) StateShipping {
	var stateName StateShipping
	for state, stateIDs := range StateCodeMap {
		if cm.IntsContain(stateIDs, code) {
			stateName = state
			break
		}
	}
	return stateName
}

func (s StateShipping) ToStatus5() model.Status5 {
	switch s {
	case StateCanceled:
		return model.S5Negative
	case StateReturned:
		return model.S5NegSuper
	case StateDelivered:
		return model.S5Positive
	}
	return model.S5SuperPos
}

func (s StateShipping) ToShippingStatus5(old model.ShippingState) model.Status5 {
	switch s {
	case StateCanceled:
		return model.S5Negative
	case StateReturned:
		return model.S5NegSuper
	case StateDelivered:
		return model.S5Positive
	case StateStored, StatePicked, StateDelivering:
		switch old {
		case model.StateReturned, model.StateReturning:
			return model.S5NegSuper
		default:
			return model.S5SuperPos
		}
	}
	return model.S5SuperPos
}

func (s StateShipping) ToModel(old model.ShippingState) model.ShippingState {
	switch s {
	case StateCanceled:
		return model.StateCancelled
	case StateNotConfirm, StateConfirmed:
		return model.StateCreated
	case StatePicking:
		return model.StatePicking
	case StateStored, StatePicked:
		switch old {
		case model.StateReturned, model.StateReturning:
			return old
		default:
			return model.StateHolding
		}
	case StateDelivering:
		switch old {
		case model.StateReturned, model.StateReturning:
			return old
		default:
			return model.StateDelivering
		}
	case StateDelivered:
		return model.StateDelivered
	case StateReturning:
		return model.StateReturning
	case StateReturned:
		return model.StateReturned
	default:
		return model.StateUnknown
	}
}

type LoginRequest struct {
	Username string `json:"USERNAME"`
	Password string `json:"PASSWORD"`
}

type LoginResponse struct {
	CommonResponse
	Data struct {
		UserId    int    `json:"userId"`
		Token     string `json:"token"`
		Partner   int    `json:"partner"`
		Phone     string `json:"phone"`
		Expired   int    `json:"expired"`
		Encrypted string `json:"encrypted"`
		Source    int    `json:"source"`
	} `json:"data"`
}

type WarehouseResponse struct {
	CommonResponse
	Data []*Warehouse `json:"data"`
}

type Warehouse struct {
	CusID          int    `json:"cusId"`
	GroupAddressID int    `json:"groupaddressId"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	ProvinceID     int    `json:"provinceId"`
	DistrictID     int    `json:"districtId"`
	WardsID        int    `json:"wardsId"`
}

type CommonResponse struct {
	Status  int    `json:"status"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type CalcShippingFeeRequest struct {
	SenderProvince   int `json:"SENDER_PROVINCE"`
	SenderDistrict   int `json:"SENDER_DISTRICT"`
	ReceiverProvince int `json:"RECEIVER_PROVINCE"`
	ReceiverDistrict int `json:"RECEIVER_DISTRICT"`
	// TH: Thư/Envelope
	// HH: Hàng hóa
	ProductType VTPostProductType `json:"PRODUCT_TYPE"`
	// VCN - VCN Chuyển phát nhanh - Express dilivery
	// VTK - VTK Tiết kiệm - Express Saver
	// SCOD - SCOD Giao hàng thu tiền - Cash on dilivery plus
	OrderService    VTPostOrderServiceCode `json:"ORDER_SERVICE"`
	OrderServiceAdd VTPostOrderServiceCode `json:"ORDER_SERVICE_ADD"`
	ProductWeight   int                    `json:"PRODUCT_WEIGHT"`
	ProductPrice    int                    `json:"PRODUCT_PRICE"`
	// Tiền thu hộ
	MoneyCollection int `json:"MONEY_COLLECTION"`
	ProductQuantity int `json:"PRODUCT_QUANTITY"`
	// =1 : Trong nước/ inland;
	// =0: Quốc tế/ international
	NATIONAL_TYPE int `json:"NATIONAL_TYPE"`
}

type CalcShippingFeeAllServicesRequest struct {
	SenderProvince   int `json:"SENDER_PROVINCE"`
	SenderDistrict   int `json:"SENDER_DISTRICT"`
	ReceiverProvince int `json:"RECEIVER_PROVINCE"`
	ReceiverDistrict int `json:"RECEIVER_DISTRICT"`
	// Loại hàng hóa/ Product type: Thư/Envelope: TH; Hàng hóa/ Goods: HH
	ProductType   string `json:"PRODUCT_TYPE"`
	ProductWeight int    `json:"PRODUCT_WEIGHT"`
	// Giá trị hàng hóa/ price
	ProductPrice int `json:"PRODUCT_PRICE"`
	// Tiền thu hộ
	MoneyCollection int `json:"MONEY_COLLECTION"`
	// Loại trạng thái:/ Status type
	// 1. Duyệt đơn hàng/ Confirm order
	// 2. Duyệt chuyển hoàn/ Confirm return shipping
	// 3. Phát tiếp/ delivery again
	// 4. Hủy đơn hàng/ delivery again
	// 5. Lấy lại đơn hàng (Gửi lại)/ get back order (re-order)
	// 11. Xóa đơn hàng đã hủy(delete canceled order)
	// set Type = 1
	Type int `json:"TYPE"`
}

type ShippingFeeService struct {
	MaDVChinh string `json:"MA_DV_CHINH"`
	TenDichvu string `json:"TEN_DICHVU"`
	GiaCuoc   int    `json:"GIA_CUOC"`
	// Tổng thời gian giao hàng
	ThoiGian       string `json:"THOI_GIAN"`
	ExchangeWeight int    `json:"EXCHANGE_WEIGHT"`
	// Các dịch vụ thêm, bao gồm các thuộc tính SERVICE_CODE": Mã dịch vụ, "SERVICE_NAME": Tên dịch vụ, "DESCRIPTION": Mô tả
	ExtraService []*ExtraService `json:"EXTRA_SERVICE"`
}

type ExtraService struct {
	ServiceCode string `json:"SERVICE_CODE"`
	ServiceName string `json:"SERVICE_NAME"`
	Description string `json:"DESCRIPTION"`
}

type ShippingFeeLine struct {
	ServiceCode string `json:"SERVICE_CODE"`
	ServiceName string `json:"SERVICE_NAME"`
	Price       string `json:"PRICE"`
}

type _CalcShippingFeeResponse struct {
	CommonResponse
	Data json.RawMessage `json:"data"`
}

type CalcShippingFeeResponse struct {
	CommonResponse
	Data *ShippingFeeData `json:"data"`
}

type ShippingFeeData struct {
	// Tổng tiền ban đầu
	MoneyTotalOld int `json:"MONEY_TOTAL_OLD"`
	// Tổng tiền
	MoneyTotal int `json:"MONEY_TOTAL"`
	// Cước chính/ main charge
	MoneyTotalFee int `json:"MONEY_TOTAL_FEE"`
	// Phụ phí (phụ phí xăng dầu...)
	MoneyFee int `json:"MONEY_FEE"`
	// Tiền thu hộ ( Số tiền khách hàng muốn VTP thu của người nhận hàng)
	MoneyCollectionFee int `json:"MONEY_COLLECTION_FEE"`
	// Phụ phí khác ( phụ phí đóng gói, phát sinh khác….)/
	MoneyOtherFee int `json:"MONEY_OTHER_FEE"`
	// Phí gia tăng ( các dịch vụ cộng thêm khác có phát sinh)
	MoneyVAS int `json:"MONEY_VAS"`
	// Tổng phí VAT/ Sum of VAT
	MoneyVAT int     `json:"MONEY_VAT"`
	KpiHt    float32 `json:"KPI_HT"`
}

func ToShippingService(lines []*model.ShippingFeeLine, providerServiceID string, orderService VTPostOrderServiceCode, expectedPickTime, expectedDeliveryTime time.Time) *model.AvailableShippingService {
	if lines == nil {
		return nil
	}
	shippingFeeShop := 0
	for _, line := range lines {
		shippingFeeShop += line.Cost
	}
	return &model.AvailableShippingService{
		Name:               orderService.Name(),
		ServiceFee:         shippingFeeShop,
		Provider:           model.TypeVTPost,
		ProviderServiceID:  providerServiceID,
		ExpectedPickAt:     expectedPickTime,
		ExpectedDeliveryAt: expectedDeliveryTime,
	}
}

func (s *ShippingFeeService) ToAvailableShippingService(providerServiceID string, expectedPickTime, expectedDeliveryTime time.Time) *model.AvailableShippingService {
	serviceCode := VTPostOrderServiceCode(s.MaDVChinh)
	return &model.AvailableShippingService{
		Name:               serviceCode.Name(),
		ServiceFee:         s.GiaCuoc,
		Provider:           model.TypeVTPost,
		ProviderServiceID:  providerServiceID,
		ExpectedPickAt:     expectedPickTime,
		ExpectedDeliveryAt: expectedDeliveryTime,
	}
}

func (sf *ShippingFeeData) CalcAndConvertShippingFeeLines() ([]*model.ShippingFeeLine, error) {
	var res []*model.ShippingFeeLine
	total := 0

	if sf.MoneyCollectionFee != 0 {
		fee := sf.MoneyCollectionFee
		shippingFeeLine := &model.ShippingFeeLine{
			ShippingFeeType:     model.ShippingFeeTypeCODS,
			Cost:                int(math.Floor(float64(fee) * IncludeVAT)),
			ExternalServiceName: "Phụ phí thu hộ",
		}
		res = append(res, shippingFeeLine)
		total += shippingFeeLine.Cost
	}
	if sf.MoneyOtherFee != 0 {
		// xem như phí bảo hiểm
		fee := sf.MoneyOtherFee
		shippingFeeLine := &model.ShippingFeeLine{
			ShippingFeeType:     model.ShippingFeeTypeInsurance,
			Cost:                int(math.Floor(float64(fee) * IncludeVAT)),
			ExternalServiceName: "Phụ phí bảo hiểm",
		}
		res = append(res, shippingFeeLine)
		total += shippingFeeLine.Cost
	}
	if sf.MoneyFee != 0 {
		fee := sf.MoneyFee
		shippingFeeLine := &model.ShippingFeeLine{
			ShippingFeeType:     model.ShippingFeeTypeOther,
			Cost:                int(math.Floor(float64(fee) * IncludeVAT)),
			ExternalServiceName: "Phụ phí xăng dầu",
		}
		res = append(res, shippingFeeLine)
		total += shippingFeeLine.Cost
	}
	if sf.MoneyTotalFee != 0 {
		fee := sf.MoneyTotal - total
		moneyTotalFee := int(math.Floor(float64(sf.MoneyTotalFee) * IncludeVAT))
		if math.Abs(float64(fee-moneyTotalFee)) > 10 {
			return nil, cm.Error(cm.FailedPrecondition, "VTPOST: Total shipping does not match", nil)
		}
		shippingFeeLine := &model.ShippingFeeLine{
			ShippingFeeType:     model.ShippingFeeTypeMain,
			Cost:                fee,
			ExternalServiceName: "Cước chính",
		}
		res = append(res, shippingFeeLine)
		total += shippingFeeLine.Cost
	}

	if total == 0 {
		return nil, cm.Error(cm.FailedPrecondition, "VTPOST: Total shipping does not found", nil)
	}
	return res, nil
}

type GetProvinceResponse struct {
	CommonResponse
	Data []*Province `json:"data"`
}

type Province struct {
	ProvinceID   int    `json:"PROVINCE_ID"`
	ProvinceCode string `json:"PROVINCE_CODE"`
	ProvinceName string `json:"PROVINCE_NAME"`
}

type GetDistrictResponse struct {
	CommonResponse
	Data []*District `json:"data"`
}

type District struct {
	DistrictID int `json:"DISTRICT_ID"`
	// Mã số quận / huyện/ District code
	DistrictValue string `json:"DISTRICT_VALUE"`
	DistrictName  string `json:"DISTRICT_Name"`
	ProvinceID    int    `json:"PROVINCE_ID"`
}

type GetDistrictsByProvinceRequest struct {
	ProvinceID string
}

type Ward struct {
	WardsID    int    `json:"WARDS_ID"`
	WardsName  string `json:"WARDS_NAME"`
	DistrictID int    `json:"DISTRICT_ID"`
}

type GetWardsByDistrictRequest struct {
	DistrictID int `url:"districtId"`
}

type GetWardsResponse struct {
	CommonResponse
	Data []*Ward `json:"data"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.SenderFullname == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên người liên hệ lấy hàng")
	}
	if r.SenderAddress == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu địa chỉ để lấy hàng")
	}
	if r.SenderProvince == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên tỉnh/thành phố nơi lấy hàng")
	}
	if r.SenderDistrict == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên quận/huyện nơi lấy hàng")
	}
	if r.SenderPhone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu số điện thoại liên hệ nơi lấy hàng")
	}
	if r.ReceiverFullname == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên người nhận hàng")
	}
	if r.ReceiverAddress == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu địa chỉ nhận hàng")
	}
	if r.ReceiverProvince == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên tỉnh/thành phố nơi nhận hàng")
	}
	if r.ReceiverDistrict == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu tên quận/huyện nơi nhận hàng")
	}
	if r.ReceiverPhone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Lỗi tạo đơn hàng với VTPost: Thiếu số điện thoại liên hệ nơi nhận hàng")
	}
	return nil
}

type CreateOrderRequest struct {
	OrderNumber        string `json:"ORDER_NUMBER"`
	GroupAddressID     int    `json:"GROUPADDRESS_ID"`
	CusID              int    `json:"CUS_ID"`
	DeliveryDate       string `json:"DELIVERY_DATE"`
	SenderFullname     string `json:"SENDER_FULLNAME"`
	SenderAddress      string `json:"SENDER_ADDRESS"`
	SenderPhone        string `json:"SENDER_PHONE"`
	SenderEmail        string `json:"SENDER_EMAIL"`
	SenderWard         int    `json:"SENDER_WARD"`
	SenderDistrict     int    `json:"SENDER_DISTRICT"`
	SenderProvince     int    `json:"SENDER_PROVINCE"`
	SenderLatitude     int    `json:"SENDER_LATITUDE"`
	SenderLongtitude   int    `json:"SENDER_LONGITUDE"`
	ReceiverFullname   string `json:"RECEIVER_FULLNAME"`
	ReceiverAddress    string `json:"RECEIVER_ADDRESS"`
	ReceiverPhone      string `json:"RECEIVER_PHONE"`
	ReceiverEmail      string `json:"RECEIVER_EMAIL"`
	ReceiverWard       int    `json:"RECEIVER_WARD"`
	ReceiverDistrict   int    `json:"RECEIVER_DISTRICT"`
	ReceiverProvince   int    `json:"RECEIVER_PROVINCE"`
	ReceiverLatitude   int    `json:"RECEIVER_LATITUDE"`
	ReceiverLongtitude int    `json:"RECEIVER_LONGITUDE"`
	ProductName        string `json:"PRODUCT_NAME"`
	ProductDescription string `json:"PRODUCT_DESCRIPTION"`
	ProductQuantity    int    `json:"PRODUCT_QUANTITY"`
	ProductPrice       int    `json:"PRODUCT_PRICE"`
	ProductWeight      int    `json:"PRODUCT_WEIGHT"`
	ProductLenght      int    `json:"PRODUCT_LENGTH"`
	ProductHeight      int    `json:"PRODUCT_WIDTH"`
	// TH: Thư/ Envelope
	// HH: Hàng hóa / Goods
	ProductType VTPostProductType `json:"PRODUCT_TYPE"`
	// Loại vận đơn/ Oder type
	// 1: Không thu tiền/ Uncollect money
	// 2: Thu hộ tiền cước và tiền hàng/ Collect express fee and price of goods.
	// 3: Thu hộ tiền hàng/ Collect price of goods
	// 4: Thu hộ tiền cước/ Collect express fee.
	Orderpayment    int                    `json:"ORDER_PAYMENT"` // use 3
	OrderService    VTPostOrderServiceCode `json:"ORDER_SERVICE"`
	OrderServiceAdd VTPostOrderServiceCode `json:"ORDER_SERVICE_ADD"`
	OrderVoucher    string                 `json:"ORDER_VOUCHER"`
	OrderNote       string                 `json:"ORDER_NOTE"`
	// Tiền thu hộ
	MoneyCollection int `json:"MONEY_COLLECTION"`
	// Cước chính/ main charge
	MoneyTotalFee int `json:"MONEY_TOTALFEE"`
	// Phí thu hộ ( phí phải trả cho VTP khi VTP thu hộ COD)/ Collection fee (payable fee to VTP when VTP collects money-COD)
	MoneyFeeCOD int `json:"MONEY_FEECOD"`
	// Phí gia tăng ( các dịch vụ cộng thêm khác có phát sinh)/ added fee (fee of addition services)
	MoneyFeeVAS int `json:"MONEY_FEEVAS"`
	// Phí bảo hiểm (= 1% giá trị khai giá) / fee of insurance (= 1% value-based pricing)
	MoneyFeeInsurance int `json:"MONEY_FEEINSURRANCE"`
	// Phụ phí (phụ phí xăng dầu, kết nối huyện...)/ Sub-fee (oil fee, connection fee,…)
	MoneyFee int `json:"MONEY_FEE"`
	// Phụ phí khác ( phụ phí đóng gói, phát sinh khác….)/ Other fee (pack fee,..)
	MoneyFeeOther int        `json:"MONEY_FEEOTHER"`
	MoneyTotalVAT int        `json:"MONEY_TOTALVAT"`
	MoneyTotal    int        `json:"MONEY_TOTAL"`
	ListItem      []*Product `json:"LIST_ITEM"`
}

type Product struct {
	ProductName     string `json:"PRODUCT_NAME"`
	ProductPrice    int    `json:"PRODUCT_PRICE"`
	ProductWeight   int    `json:"PRODUCT_WEIGHT"`
	ProductQuantity int    `json:"PRODUCT_QUANTITY"`
}

type CreateOrderResponse struct {
	CommonResponse
	Data struct {
		OrderNumber   string `json:"ORDER_NUMBER"`
		ProductWeight int    `json:"PRODUCT_WEIGHT"`
		// Tổng tiền
		MoneyTotal int `json:"MONEY_TOTAL"`
		// Cước chính/ main charge
		MoneyTotalFee int `json:"MONEY_TOTAL_FEE"`
		// Phụ phí (phụ phí xăng dầu...)
		MoneyFee int `json:"MONEY_FEE"`
		// Tiền thu hộ ( Số tiền khách hàng muốn VTP thu của người nhận hàng)
		MoneyCollection int `json:"MONEY_COLLECTTION"`
		// Cước phí thu hộ
		MoneyCollectionFee int `json:"MONEY_COLLECTION_FEE"`
		// Phụ phí khác ( phụ phí đóng gói, phát sinh khác….)/
		MoneyOtherFee int `json:"MONEY_OTHER_FEE"`
		// Tổng phí VAT/ Sum of VAT
		MoneyFeeVAT int `json:"MONEY_FEE_VAT"`
		// Cân nặng
		ExchangeWeight int     `json:"EXCHANGE_WEIGHT"`
		KpiHt          float32 `json:"KPI_HT"`
	} `json:"data"`
}

type CancelOrderRequest struct {
	Type        int    `json:"TYPE"`
	OrderNumber string `json:"ORDER_NUMBER"`
	Note        int    `json:"NOTE"`
}

type CallbackOrder struct {
	Data CallbackOrderData `json:"DATA"`
	// ApiKey đối tác/ partner’s token
	Token string `json:"TOKEN"`
}

type CallbackOrderData struct {
	// Mã vận đơn VTP/ VTP Order number
	OrderNumber string `json:"ORDER_NUMBER"`
	// Mã vận đơn đối tác/ Partner Order number
	OrderReference string `json:"ORDER_REFERENCE"`
	// Ngày tháng trạng thái (dd/MM/yyyy H:m:s)/ Date of status
	OrderStatusDate string `json:"ORDER_STATUSDATE"`
	// Mã trạng thái/ status code
	OrderStatus       int    `json:"ORDER_STATUS"`
	StatusName        string `json:"STATUS_NAME"`
	LocationCurrently string `json:"LOCALION_CURRENTLY"`
	Note              string `json:"NOTE"`
	// Tiền thu hộ/ collection money
	MoneyCollection int `json:"MONEY_COLLECTION"`
	// Phí thu hộ/ collection fee
	MoneyFeeCOD int `json:"MONEY_FEECOD"`
	// Tổng tiền/ total  money
	MoneyTotal       int    `json:"MONEY_TOTAL"`
	ExpectedDelivery string `json:"EXPECTED_DELIVERY"`
	ProductWeight    int    `json:"PRODUCT_WEIGHT"`
	OrderService     string `json:"ORDER_SERVICE"`
}

type ResponseInterface interface {
	GetCommonResponse() CommonResponse
}

func (c CommonResponse) GetCommonResponse() CommonResponse {
	return c
}

func GetInsuranceFee(productPrice int) int {
	// BẢO HIỂM (Mã: GBH)
	// 1% Giá trị khai giá
	// Tối thiểu 15.000VNĐ/bưu gửi. (chưa bao gồm VAT)
	insuranceFee := math.Floor(float64(productPrice) * 0.01)
	if insuranceFee < 15000 {
		insuranceFee = 15000
	}
	return int(insuranceFee * 1.1)
}
