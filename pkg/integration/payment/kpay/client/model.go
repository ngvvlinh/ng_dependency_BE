package client

import (
	"encoding/json"

	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/backend/pkg/common/apifw/httpreq"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type BodyResponse struct {
	XApiMessage string `json:"x-api-message"`
}

type DecryptedBodyResponse struct {
	Code ErrorCode       `json:"code"`
	Data json.RawMessage `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type BodyRequest struct {
	XAPIMessage string `json:"x-api-message"`
}

type ErrorCode int

const (
	SystemError       ErrorCode = 500
	ServerMaintenance ErrorCode = 501
	InvalidParams     ErrorCode = 400
	InvalidToken      ErrorCode = 401
	RequestSuccess    ErrorCode = 1000
	RequestFail       ErrorCode = 1001
	RequestRefused    ErrorCode = 1002
)

var mapErrorCode = map[ErrorCode]string{
	SystemError:       "Lỗi dịch vụ, liên hệ bộ phận liên quan để xử lý",
	ServerMaintenance: "Server đang bảo trì",
	InvalidParams:     "Dữ liệu yêu cầu không hợp lệ",
	InvalidToken:      "AccessToken chứng thực yêu cầu không hợp lệ",
	RequestSuccess:    "Yêu cầu thành công",
	RequestFail:       "Yêu cầu thất bại, kèm theo message trong data để thông báo lý do lỗi thất bại",
	RequestRefused:    "Yêu cầu bị từ chối vì 1 số lý do bảo mật",
}

type PaymentStatus string

const (
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED" // Thanh toán đã thành công
	PaymentStatusFailed    PaymentStatus = "FAILED"    // Thanh toán thất bại link đã vô hiệu lực
	PaymentStatusPending   PaymentStatus = "PENDING"   // Thanh toán đang chờ khách hàng thanh toán
	PaymentStatusExpired   PaymentStatus = "EXPIRED"   // Thanh toán đã hết hạn không còn hiệu lực
	PaymentStatusCanceled  PaymentStatus = "CANCELED"  // Thanh toán đã bị huỷ từ Merchant
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"  // Thanh toán đã được hoàn lại cho user trong 1 số trường hợp đặc biêt
)

var PaymentStatusMap = map[PaymentStatus]string{
	PaymentStatusSucceeded: "Thanh toán đã thành công",
	PaymentStatusFailed:    "Thanh toán thất bại link đã vô hiệu lực",
	PaymentStatusPending:   "Thanh toán đang chờ khách hàng thanh toán",
	PaymentStatusExpired:   "Thanh toán đã hết hạn không còn hiệu lực",
	PaymentStatusCanceled:  "Thanh toán đã bị huỷ từ Merchant",
	PaymentStatusRefunded:  "Thanh toán đã được hoàn lại cho user trong 1 số trường hợp đặc biêt",
}

func (s PaymentStatus) ToState() payment_state.PaymentState {
	switch s {
	case PaymentStatusPending:
		return payment_state.Pending
	case PaymentStatusExpired:
		return payment_state.Expired
	case PaymentStatusFailed:
		return payment_state.Failed
	case PaymentStatusSucceeded:
		return payment_state.Success
	case PaymentStatusCanceled:
		return payment_state.Cancelled
	case PaymentStatusRefunded:
		return payment_state.Refunded
	default:
		return payment_state.Unknown
	}
}

func (s PaymentStatus) ToStatus() status4.Status {
	switch s {
	case PaymentStatusPending:
		return status4.S
	case PaymentStatusSucceeded:
		return status4.P
	case PaymentStatusFailed, PaymentStatusExpired, PaymentStatusCanceled, PaymentStatusRefunded:
		return status4.N
	default:
		return status4.S
	}
}

func (s PaymentStatus) ToString() string {
	return string(s)
}

func url(baseUrl, path string) string {
	return baseUrl + path
}

type CreateTransactionResponse struct {
	Url         String `json:"url"`
	Transaction String `json:"transaction"`
}

type CreateTransactionRequest struct {
	Amount             int    `json:"amount"` // min: 2.000vnd max: 10.000.000vnd
	Desc               string `json:"desc"`
	PartnerTransaction string `json:"partnerTransaction"` // unique
	FailedURL          string `json:"failedUrl"`
	RedirectURL        string `json:"redirectUrl"`
}

type GetTransactionRequest struct {
	Transaction string
}

type GetTransactionResponse struct {
	PartnerTransaction String        `json:"partnerTransaction"`
	Amount             Int           `json:"amount"`
	Fee                Int           `json:"fee"`
	Total              Int           `json:"total"`
	State              PaymentStatus `json:"state"`
	Url                String        `json:"url"`
	Desc               String        `json:"desc"`
	Transaction        String        `json:"transaction"`
	ExpiryAt           Time          `json:"expiryAt"`
	CreatedAt          Time          `json:"createdAt"`
	UpdatedAt          Time          `json:"updatedAt"`
}

type CancelTransactionRequest struct {
	Transaction string `json:"transaction"`
}

type CancelTransactionResponse struct {
	Transaction String `json:"transaction"`
}

type CallbackTransaction struct {
	Transaction        String        `json:"transaction"`
	PartnerTransaction String        `json:"partnerTransaction"`
	Amount             Int           `json:"amount"`
	Fee                Int           `json:"fee"`
	Total              Int           `json:"total"`
	State              PaymentStatus `json:"state"`
	Desc               String        `json:"desc"`
	CreatedAt          Time          `json:"createdAt"`
	UpdatedAt          Time          `json:"updatedAt"`
}
