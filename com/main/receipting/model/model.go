package model

import (
	"time"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenReceipt(&Receipt{})

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
	OrderIDs    []int64
	Lines       []*ReceiptLine
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
}

var _ = sqlgenReceiptLine(&ReceiptLine{})

type ReceiptLine struct {
	OrderID int64  `json:"order_id"`
	Title   string `json:"title"`
	Amount  int32  `json:"amount"`
}
