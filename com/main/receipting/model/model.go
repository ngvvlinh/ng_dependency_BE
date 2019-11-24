package model

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenReceipt(&Receipt{})

type Receipt struct {
	ID                 dot.ID
	ShopID             dot.ID
	TraderID           dot.ID
	Code               string
	CodeNorm           int32
	Title              string
	Type               string
	Description        string
	TraderFullNameNorm string

	Amount          int32
	Status          model.Status3
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
	Amount int32  `json:"amount"`
}

type Trader struct {
	ID       dot.ID `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
