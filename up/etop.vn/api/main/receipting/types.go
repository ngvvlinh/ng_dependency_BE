package receipting

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	dot "etop.vn/capi/dot"
)

// +gen:event:topic=event/receipting

type ReceiptType string
type ReceiptCreatedType string
type ReceiptRefType string

func (t ReceiptType) String() string        { return string(t) }
func (t ReceiptCreatedType) String() string { return string(t) }
func (t ReceiptRefType) String() string     { return string(t) }

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
	ID          dot.ID
	ShopID      dot.ID
	TraderID    dot.ID
	Code        string
	CodeNorm    int
	Title       string
	Type        ReceiptType
	Description string
	Amount      int
	Status      status3.Status
	LedgerID    dot.ID
	RefIDs      []dot.ID
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time
	CreatedBy   dot.ID
	CreatedType ReceiptCreatedType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReceiptLine struct {
	RefID  dot.ID
	Title  string
	Amount int
}

type Trader struct {
	ID       dot.ID
	Type     string
	FullName string
	Phone    string
}

func (r *Receipt) GetRefIDs() []dot.ID {
	ids := make([]dot.ID, 0, len(r.Lines))
	for _, line := range r.Lines {
		if line.RefID != 0 {
			ids = append(ids, line.RefID)
		}
	}
	return ids
}

type MoneyTransactionConfirmedEvent struct {
	meta.EventMeta
	ShopID             dot.ID
	MoneyTransactionID dot.ID
}

type ReceiptConfirmedEvent struct {
	meta.EventMeta
	ShopID    dot.ID
	ReceiptID dot.ID
}

type ReceiptCancelledEvent struct {
	meta.EventMeta
	ShopID    dot.ID
	ReceiptID dot.ID
}

type ReceiptCreatingEvent struct {
	meta.EventMeta
	RefIDs         []dot.ID
	MapRefIDAmount map[dot.ID]int
	Receipt        *Receipt
}
