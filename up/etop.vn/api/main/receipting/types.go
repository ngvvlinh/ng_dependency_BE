package receipting

import (
	"time"

	"etop.vn/api/meta"
)

// +gen:event:topic=event/receipting

type ReceiptType string
type ReceiptCreatedType string
type ReceiptRefType string

const (
	ReceiptTypeReceipt ReceiptType = "receipt"
	ReceiptTypePayment ReceiptType = "payment"

	// Created type
	ReceiptCreatedTypeManual ReceiptCreatedType = "manual"
	ReceiptCreatedTypeAuto   ReceiptCreatedType = "auto"

	ReceiptRefTypeOrder         ReceiptRefType = "order"
	ReceiptRefTypeFulfillment   ReceiptRefType = "fulfillment"
	ReceiptRefTypePurchaseOrder ReceiptRefType = "purchase_order"
)

type Receipt struct {
	ID          int64
	ShopID      int64
	TraderID    int64
	Code        string
	CodeNorm    int32
	Title       string
	Type        ReceiptType
	Description string
	Amount      int32
	Status      int32
	LedgerID    int64
	RefIDs      []int64
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time
	CreatedBy   int64
	CreatedType ReceiptCreatedType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReceiptLine struct {
	RefID  int64
	Title  string
	Amount int32
}

type Trader struct {
	ID       int64
	Type     string
	FullName string
	Phone    string
}

func (r *Receipt) GetRefIDs() []int64 {
	ids := make([]int64, 0, len(r.Lines))
	for _, line := range r.Lines {
		if line.RefID != 0 {
			ids = append(ids, line.RefID)
		}
	}
	return ids
}

type MoneyTransactionConfirmedEvent struct {
	meta.EventMeta
	ShopID             int64
	MoneyTransactionID int64
}

type ReceiptConfirmedEvent struct {
	meta.EventMeta
	ShopID    int64
	ReceiptID int64
}

type ReceiptCancelledEvent struct {
	meta.EventMeta
	ShopID    int64
	ReceiptID int64
}

type ReceiptCreatingEvent struct {
	meta.EventMeta
	RefIDs         []int64
	MapRefIDAmount map[int64]int32
	Receipt        *Receipt
}
