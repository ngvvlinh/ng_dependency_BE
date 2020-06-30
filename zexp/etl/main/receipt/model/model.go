package model

import (
	"time"

	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type Receipt struct {
	ID          dot.ID
	ShopID      dot.ID
	TraderID    dot.ID
	Code        string
	Title       string
	Type        receipt_type.ReceiptType `sql_gen:"enum(receipt_type)"`
	Description string
	TraderType  string

	Amount          int
	Status          status3.Status `sql_type:"int2"`
	RefIDs          []dot.ID
	RefType         receipt_ref.ReceiptRef `sql_gen:"enum(receipt_ref_type)"`
	Lines           []*ReceiptLine
	LedgerID        dot.ID
	Trader          *Trader
	CancelledReason string
	CreatedType     receipt_mode.ReceiptMode `sql_gen:"enum(receipt_created_type)"`
	CreatedBy       dot.ID
	PaidAt          time.Time
	ConfirmedAt     time.Time
	CancelledAt     time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Note string
	Rid  dot.ID
}

type ReceiptLine struct {
	RefID  dot.ID `json:"ref_id"`
	Title  string `json:"title"`
	Amount int    `json:"amount"`
}

type Trader struct {
	ID       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
