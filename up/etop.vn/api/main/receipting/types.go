package receipting

import (
	"time"
)

const (
	ReceiptType = "receipt"
	PaymentType = "payment"
)

type Receipt struct {
	ID          int64
	ShopID      int64
	TraderID    int64
	UserID      int64
	Code        string
	Title       string
	Type        string
	Description string
	Amount      int32
	Status      int32
	Lines       []*ReceiptLine
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReceiptLine struct {
	OrderID        int64
	Title          string
	Amount         int32
	ReceivedAmount int32
}

func (r *Receipt) GetOrderIDs() []int64 {
	ids := make([]int64, 0, len(r.Lines))
	for _, line := range r.Lines {
		if line.OrderID != 0 {
			ids = append(ids, line.OrderID)
		}
	}
	return ids
}
