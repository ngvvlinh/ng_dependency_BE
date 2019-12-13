package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenReceipt(&Receipt{})

type Receipt struct {
	ID                 dot.ID
	ShopID             dot.ID
	TraderID           dot.ID
	Code               string
	CodeNorm           int
	Title              string
	Type               string
	Description        string
	TraderFullNameNorm string

	Amount          int
	Status          status3.Status
	RefIDs          []dot.ID
	RefType         string
	Lines           []*ReceiptLine
	LedgerID        dot.ID
	Trader          *Trader
	CancelledReason string
	CreatedType     string
	CreatedBy       dot.ID
	PaidAt          time.Time
	ConfirmedAt     time.Time
	CancelledAt     time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	DeletedAt       time.Time
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
