package vtpostimport

import (
	"time"

	"o.o/api/main/moneytx"
)

var dateTimeLayouts = []string{"01-02-06", "02/01/2006 15:04"}

type VTPostMoneyTransactionShippingExternalLine struct {
	ExternalCode     string
	DeliveredAt      time.Time
	TotalCOD         int // tiền hàng
	TotalShippingFee int // tổng tiền cước
	Total            int // Số tiền phải trả
}

func (line *VTPostMoneyTransactionShippingExternalLine) ToModel() *moneytx.MoneyTransactionShippingExternalLine {
	return &moneytx.MoneyTransactionShippingExternalLine{
		ExternalCode:             line.ExternalCode,
		ExternalTotalCOD:         line.TotalCOD,
		ExternalClosedAt:         line.DeliveredAt,
		ExternalTotalShippingFee: line.TotalShippingFee,
	}
}

func ToMoneyTransactionShippingExternalLines(lines []*VTPostMoneyTransactionShippingExternalLine) []*moneytx.MoneyTransactionShippingExternalLine {
	if lines == nil {
		return nil
	}
	res := make([]*moneytx.MoneyTransactionShippingExternalLine, len(lines))
	for i, line := range lines {
		res[i] = line.ToModel()
	}
	return res
}
