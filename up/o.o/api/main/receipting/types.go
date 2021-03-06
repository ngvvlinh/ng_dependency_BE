package receipting

import (
	"time"

	"o.o/api/meta"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	dot "o.o/capi/dot"
)

// +gen:event:topic=event/receipting

type Receipt struct {
	ID           dot.ID
	ShopID       dot.ID
	TraderID     dot.ID
	Code         string
	CodeNorm     int
	Title        string
	CancelReason string
	Type         receipt_type.ReceiptType
	Description  string
	Amount       int
	Status       status3.Status
	LedgerID     dot.ID
	RefIDs       []dot.ID
	RefType      receipt_ref.ReceiptRef
	Lines        []*ReceiptLine
	Trader       *Trader
	PaidAt       time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CreatedBy    dot.ID
	Mode         receipt_mode.ReceiptMode
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Note         string
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

type ReceiptConfirmingEvent struct {
	meta.EventMeta
	ShopID      dot.ID
	ReceiptID   dot.ID
	RefType     receipt_ref.ReceiptRef
	ReceiptType receipt_type.ReceiptType
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
