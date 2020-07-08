package ghtkimport

import (
	"time"

	"o.o/api/main/moneytx"
)

type GHTKMoneyTransactionShippingExternalLine struct {
	ExternalCode     string
	ShopCode         string
	Customer         string
	TotalCOD         int
	InsuranceFee     int
	ShippingFee      int
	ReturnFee        int
	Discount         int
	ChangeAddressFee int
	Total            int // after sub fees
	CreatedAt        time.Time
	DeliveredAt      time.Time
}

func (line *GHTKMoneyTransactionShippingExternalLine) ToModel() *moneytx.MoneyTransactionShippingExternalLine {
	return &moneytx.MoneyTransactionShippingExternalLine{
		ExternalCode:         line.ExternalCode,
		ExternalCustomer:     line.Customer,
		ExternalTotalCOD:     line.TotalCOD,
		ExternalCreatedAt:    line.CreatedAt,
		ExternalClosedAt:     line.DeliveredAt,
		EtopFulfillmentIDRaw: line.ShopCode,
	}
}

func ToMoneyTransactionShippingExternalLines(lines []*GHTKMoneyTransactionShippingExternalLine) []*moneytx.MoneyTransactionShippingExternalLine {
	if lines == nil {
		return nil
	}
	res := make([]*moneytx.MoneyTransactionShippingExternalLine, len(lines))
	for i, line := range lines {
		res[i] = line.ToModel()
	}
	return res
}

const (
	ExternalCode     = "Mã đơn hàng"
	ShopCode         = "Mã đơn hàng shop"
	CustomerInfo     = "Thông tin khách hàng"
	TotalCOD         = "Tổng tiền thu hộ"
	InsuranceFee     = "Phí bảo hiểm"
	ShippingFee      = "Phí dịch vụ"
	ReturnFee        = "Phí chuyển hoàn"
	Discount         = "Khuyến mãi"
	ChangeAddressFee = "Phí thay đổi địa chỉ giao"
	Total            = "Thanh toán"
	CreatedAt        = "Ngày tạo"
	DeliveredAt      = "Ngày hoàn thành"
)

var (
	HeaderStrings = []string{ExternalCode, ShopCode, CustomerInfo, TotalCOD,
		InsuranceFee, ShippingFee, ReturnFee, Discount, ChangeAddressFee, Total, CreatedAt, DeliveredAt}
)
