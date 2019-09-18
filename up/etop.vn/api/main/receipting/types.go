package receipting

import (
	"time"
)

type Receipt struct {
	ID          int64
	ShopID      int64
	TraderID    int64
	UserID      int64
	OrderIDs    []int64
	Code        string
	Title       string
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
