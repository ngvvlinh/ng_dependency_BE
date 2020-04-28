package client

import (
	"strconv"

	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	cm "o.o/backend/pkg/common"
)

type (
	PaymentStatus string
)

const (
	PaymentStatusNotCreate PaymentStatus = "-1"
	PaymentStatusPending   PaymentStatus = "0"
	PaymentStatusSuccess   PaymentStatus = "1"
	PaymentStatusFailed    PaymentStatus = "2"
	PaymentStatusUnknown   PaymentStatus = "3"
	PaymentStatusCancelled PaymentStatus = "5"

	TransactionSuccessCode = "00"
)

var PaymentStatusMap = map[PaymentStatus]string{
	PaymentStatusNotCreate: "Chưa phát sinh giao dịch",
	PaymentStatusPending:   "Giao dịch đang chờ xử lý",
	PaymentStatusSuccess:   "Giao dịch thành công",
	PaymentStatusFailed:    "Giao dịch thất bại",
	PaymentStatusUnknown:   "Giao dịch chưa rõ kết quả",
	PaymentStatusCancelled: "Khách hàng hủy không thanh toán",
}

func (s PaymentStatus) ToState() payment_state.PaymentState {
	switch s {
	case PaymentStatusNotCreate:
		return payment_state.Default
	case PaymentStatusPending:
		return payment_state.Pending
	case PaymentStatusSuccess:
		return payment_state.Success
	case PaymentStatusFailed:
		return payment_state.Failed
	case PaymentStatusUnknown:
		return payment_state.Unknown
	case PaymentStatusCancelled:
		return payment_state.Cancelled
	default:
		return payment_state.Unknown
	}
}

func (s PaymentStatus) ToStatus() status4.Status {
	switch s {
	case PaymentStatusNotCreate:
		return status4.Z
	case PaymentStatusPending, PaymentStatusUnknown:
		return status4.S
	case PaymentStatusSuccess:
		return status4.P
	case PaymentStatusFailed:
		return status4.N
	default:
		return status4.S
	}
}

var TransactionResultCodeMap = map[string]string{
	TransactionSuccessCode: "Giao dịch thành công",
	"22":                   "KH nhập sai OTP tại CTT",
	"V01":                  "Sai check_sum",
	"V02":                  "KH nhập sai OTP tại CTT",
	"V03":                  "OTP hết hạn",
	"21":                   "KH nhập sai mật khẩu (mã PIN)",
	"685":                  "KH nhập sai mật khẩu (mã PIN)",
	"16":                   "KH không đủ số dư để thanh toán",
	"W04":                  "Kết nối timeout",
	"V04":                  "Có lỗi khi truy vấn hệ thống tại VIETTEL",
	"V05":                  "Không xác nhận được giao dịch (Gọi sang API xác nhận giao dịch của đối tác thất bại)",
	"V06":                  "Khách hàng hủy thanh toán",
	"S_MAINTAIN":           "CTT bảo trì",
	"99":                   "Lỗi không xác định",
	"M01":                  "Mã đối tác chưa được đăng ký (liên hệ kỹ thuật Viettel để kiểm tra)",
	"M02":                  "Chưa thiết lập tài khoản nhận tiền cho đối tác (liên hệ kỹ thuật Viettel)",
	"M03":                  "Hình thức thanh toán không phù hợp (liên hệ kỹ thuật Viettel)",
	"M04":                  "Ảnh QR bị lỗi hoặc không đọc được giá trị cần thiết từ ảnh",
	"813":                  "Lỗi kết nối tới CTT",
}

var ReturnTransactionResultCodeMap = map[string]string{
	TransactionSuccessCode: "Hoàn hủy thành công",
	"M01":                  "Không tìm thấy thông tin của đối tác. Cần rà soát lại cấu hình",
	"M05":                  "Cấu hình hoàn tiền không phù hợp. Cần rà soát lại thông tin",
	"218":                  "Lỗi không truyền sang mã giao dịch của đối tác",
	"KG3":                  "Không tìm được giao dịch thanh toán tương ứng với thông số đối tác truyền sang",
	"27":                   "Giao dịch không phải của đối tác",
	"KG8":                  "Số tiền hoàn không hợp lệ",
	"176":                  "Trạng thái giao dịch thanh toán này không phù hợp để tiếp tục thực hiện hoàn hủy. Có thể đã hoàn hủy thành công hoặc đang không xác định hoặc đã quá hạn",
	"485":                  "Giao dịch thanh toán không thành công nên không được phép hoàn hủy",
	"655":                  "Không đủ dữ liệu để thực hiện hoàn tiền cho Khách hàng",
	"457":                  "Không đủ dữ liệu để thực hiện hoàn tiền cho Khách hàng",
	"17":                   "Tài khoản ViettelPay/BankPlus của Khách hàng không đủ điều kiện để nhận tiền hoàn. Liên hệ với Viettel",
	"159":                  "Không hỗ trợ hoàn tiền vào loại tài khoản BankPlus này của Khách hàng",
	"813":                  "Lỗi trong quá trình hoàn hủy tại Viettel. Hãy thử lại sau",
	"927":                  "Lỗi trong quá trình hoàn hủy tại Viettel. Hãy thử lại sau",
}

type ConnectPaymentGatewayRequest struct {
	// Mã hóa đơn. Trường hợp đối tác không có nhu cầu truyền mã hóa đơn thì đặt giá trị bằng với order_id. Lưu ý không dùng tiếng Việt.
	BillCode string `url:"billcode"` // required
	// Mã lệnh. Giá trị cố định là: PAYMENT
	Command string `url:"command"` // required
	// Nội dung giao dịch / đơn hàng
	Desc string `url:"desc"`
	// Mã vị trí. Giá trị mặc định là: Vi
	Locale string `url:"locale"`
	// Mã đối tác mà Viettel đã cung cấp.
	MerchantCode string `url:"merchant_code"` // required
	// Mã giao dịch duy nhất bên phía đối tác. Lưu ý không dùng tiếng Việt.
	OrderID string `url:"order_id"` // required
	// Địa chỉ để chuyển sau khi KH thanh toán.
	// Url này sử dụng trong trường hợp đối tác không truyền trong quá trình thanh toán.
	ReturnURL string `url:"return_url"` // required
	// Địa chỉ để chuyển sau nếu KH hủy giao dịch
	CancelURL string `url:"cancel_url"`
	// Số tiền của giao dịch Đơn vị tính là VND
	TransAmount int `url:"trans_amount"` // required
	// Phiên bản kết nối. Giá trị cố định là: 2.0
	Version string `url:"version" ` // required
	// Chuỗi mã hóa tạo ra dựa trên các trường dữ liệu truyền sang và được encode UTF-8: access_code + billcode + command + merchant_code + order_id + trans_amount + version
	CheckSum string `url:"check_sum"` // required
}

func (r *ConnectPaymentGatewayRequest) DataCheckSum() string {
	amountStr := strconv.Itoa(r.TransAmount)
	return r.BillCode + r.Command + r.MerchantCode + r.OrderID + amountStr + r.Version
}

func (r *ConnectPaymentGatewayRequest) Validate() error {
	if r.BillCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu mã hóa đơn.")
	}
	if r.OrderID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu mã đơn hàng")
	}
	if r.ReturnURL == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu return url")
	}
	if r.TransAmount == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu giá trị thanh toán")
	}
	return nil
}

type GetTransactionRequest struct {
	// Mã lệnh. Giá trị cố định là: TRANS_INQUIRY
	Cmd          string `url:"cmd"`
	MerchantCode string `url:"merchant_code"`
	OrderID      string `url:"order_id"`
	// Phiên bản kết nối.
	// Giá trị cố định là nếu là cổng thanh toán: 2.0 Nếu là TMĐT: WEBVIEW
	Version string `url:"version"`
	// Chuỗi mã hóa tạo ra dựa trên các trường dữ liệu truyền sang:
	// access_code + cmd + merchant_code + order_id + version
	Checksum string `url:"check_sum"`
}

func (r GetTransactionRequest) DataCheckSum() string {
	return r.Cmd + r.MerchantCode + r.OrderID + r.Version
}

type GetTransactionResponse struct {
	MerchantCode    string `json:"merchant_code"`
	OrderID         string `json:"order_id"`
	ErrorCode       string `json:"error_code"`
	VTTransactionID string `json:"vt_transaction_id"`
	// Trạng thái thanh toán
	// -1: Chưa phát sinh giao dịch
	// 0: giao dịch đang chờ xử lý 1: giao dịch thành công
	// 2: giao dịch thất bại
	// 3: giao dịch chưa rõ kết quả
	PaymentStatus string `json:"payment_status"`
	Version       string `json:"version"`
	CheckSum      string `json:"check_sum"`
}

type CancelTransactionRequest struct {
	// Mã lệnh. Giá trị cố định là: REFUND_PAYMENT
	Cmd          string `url:"cmd"`
	MerchantCode string `url:"merchant_code"`
	OrderID      string `url:"order_id"`
	// Mã giao dịch thanh toán tương ứng bên Viettel
	OriginalRequestID string `url:"originalRequestId"`
	// Hình thức hoàn tiền. 0: hoàn toàn phần | 1: hoàn một phần. 	// Hiện tại Viettel chỉ hỗ trợ hoàn toàn phần các giao dịch.
	RefundType int `url:"refundType"`
	// Số tiền hoàn (chính là số tiền thanh toán). Đối tác cần truyền sang để kiểm chứng với giao dịch thanh toán.
	TransAmount int `url:"trans_amount"`
	// Lý do hoàn hủy. Truyền sang tiếng Việt không dấu. Giới hạn 100 ký tự.
	TransContent string `url:"trans_content"`
	// Phiên bản kết nối. Giá trị cố định là: 2.0
	Version string `url:"version"`
	// Chuỗi mã hóa tạo ra dựa trên các trường dữ liệu truyền sang:
	//  access_code + cmd + merchant_code + order_id + originalRequestId + refundType + trans_amount + version
	CheckSum string `url:"check_sum"`
}

func (r CancelTransactionRequest) DataCheckSum() string {
	return r.Cmd + r.MerchantCode + r.OrderID + r.OriginalRequestID + strconv.Itoa(r.RefundType) + strconv.Itoa(r.TransAmount) + r.Version
}

type CancelTransactionResponse struct {
	MerchantCode string `json:"merchant_code"`
	OrderID      string `json:"order_id"`
	ErrorCode    string `json:"error_code"`
	ErrorMsg     string `json:"error_msg"`
	// Mã giao dịch hủy bên Viettel (nếu có)
	MerchantRequestID string `json:"merchant_request_id"`
	// Mã giao dịch hoàn tiền cho KH của Viettel (nếu có)
	RefundRequestID string `json:"refund_request_id"`
	// Phiên bản kết nối. Cố định là: 2.0
	Version string `json:"version"`
	// Chuỗi mã hóa tạo ra dựa trên các tham số dưới:
	// access_code + error_code + merchant_code + order_id + version
	CheckSum string `json:"check_sum"`
}
