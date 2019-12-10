package gateway

import (
	"context"
	"strconv"
)

// +gen:api

type Aggregate interface {
	ValidateTransaction(context.Context, *ValidateTransactionArgs) (*ValidateTransactionResult, error)

	GetResult(context.Context, *GetResultArgs) (*GetResultResult, error)
}

type ValidateTransactionArgs struct {
	BillCode     string `json:"billcode"`
	MerchantCode string `json:"merchant_code"`
	// Mã giao dịch đối tác đã gửi sang
	OrderID string `json:"order_id"`
	// Chuỗi mã hóa tạo ra dựa trên trường dữ liệu truyền sang:
	// access_code + billcode + merchant_code + order_id + trans_amount
	CheckSum string `json:"check_sum"`
}

func (r *ValidateTransactionArgs) DataCheckSum(transAmount int) string {
	amountStr := strconv.Itoa(transAmount)
	return r.BillCode + r.MerchantCode + r.OrderID + amountStr
}

type ValidateTransactionResult struct {
	BillCode     string `json:"billcode"`
	MerchantCode string `json:"merchant_code"`
	OrderID      string `json:"order_id"`
	TransAmount  string `json:"trans_amount"`
	// Mã lỗi đối tác trả về. Quy định như bên dưới
	// 00: thành công (dữ liệu tồn tại, order_id chưa thanh toán thành công, số tiền khớp với mã giao dịch)
	// 01: không thành công (dữ liệu không chính xác)
	// 02: check_sum gửi sang không đúng
	// 03: có lỗi tại hệ thống đối tác (các loại exception)
	ErrorCode ErrorCode `json:"error_code"`
	// Chuỗi mã hóa tạo ra dựa trên trường dữ liệu truyền sang:
	// access_code + billcode + error_code + merchant_code + order_id + trans_amount
	CheckSum string `json:"check_sum"`
}

func (r *ValidateTransactionResult) DataCheckSum() string {
	return r.BillCode + r.ErrorCode.String() + r.MerchantCode + r.OrderID + r.TransAmount
}

type ErrorCode string

var (
	ErrorCodeSuccess     ErrorCode = "00"
	ErrorCodeInvalidData ErrorCode = "01"
	ErrorCodeCheckSum    ErrorCode = "02"
	ErrorCodeInternal    ErrorCode = "03"
)

func (c ErrorCode) String() string { return string(c) }

type GetResultArgs struct {
	BillCode string `json:"billcode"`
	// Số điện thoại KH thanh toán
	CustMsisdn string `json:"cust_msisdn"`
	// Mã lỗi. Giao dịch thành công sẽ có mã 00.
	ErrorCode    string `json:"error_code"`
	MerchantCode string `json:"merchant_code"`
	OrderID      string `json:"order_id"`
	// Trạng thái thanh toán của KH 1: giao dịch thành công
	PaymentStatus string `json:"payment_status"`
	TransAmount   int    `json:"trans_amount"`
	// Mã giao dịch bên VIETTEL
	VtTransactionID string `json:"vt_transaction_id"`
	// Chuỗi mã hóa tạo ra dựa trên trường dữ liệu truyền sang:
	// access_code + billcode + cust_msisdn + error_code + merchant_code + order_id + payment_status + trans_amount + vt_transaction_id
	CheckSum string `json:"check_sum"`
}

type GetResultResult struct {
	ErrorCode    ErrorCode `json:"error_code"`
	MerchantCode string    `json:"merchant_code"`
	OrderID      string    `json:"order_id"`
	// Url chuyển hướng khách hàng khi thanh toán thành công trên App
	// Trường hợp không có giá trị thì trả về tham số này với giá trị là rỗng.
	ReturnURL string `json:"return_url"`
	// Chuỗi mã hóa tạo ra dựa trên trường dữ liệu truyền sang:
	//  access_code + error_code + merchant_code + order_id
	CheckSum string `json:"check_sum"`
}

func (r *GetResultResult) DataCheckSum() string {
	return r.ErrorCode.String() + r.MerchantCode + r.OrderID
}
