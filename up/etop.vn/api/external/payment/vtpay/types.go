package vtpay

import (
	"context"
	"strconv"
)

type Aggregate interface {
	HandleExternalDataResponse(context.Context, *HandleExternalDataResponseArgs) error
}

type HandleExternalDataResponseArgs struct {
	BillCode string `json:"bill_code"`
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
	// Chuỗi mã hóa tạo ra dựa trên trường dữ liệu truyền sang và được encode UTF-8:
	// access_code + billcode + cust_msisdn +
	// error_code + merchant_code + order_id +
	// payment_status + trans_amount + vt_transaction_id
	CheckSum string `json:"check_sum"`
}

func (r *HandleExternalDataResponseArgs) DataCheckSum() string {
	amountStr := strconv.Itoa(r.TransAmount)
	return r.BillCode + r.CustMsisdn + r.ErrorCode + r.MerchantCode + r.OrderID + r.PaymentStatus + amountStr + r.VtTransactionID
}
