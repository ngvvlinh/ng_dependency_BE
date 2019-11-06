package model

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenReceipt(&Receipt{})

type Receipt struct {
	ID              int64
	ShopID          int64
	TraderID        int64
	Code            string
	CodeNorm        int32
	Title           string
	Type            string
	Description     string
	Amount          int32
	Status          model.Status3
	RefIDs          []int64
	RefType         string
	Lines           []*ReceiptLine
	LedgerID        int64
	Trader          *Trader
	CancelledReason string
	CreatedType     string
	CreatedBy       int64
	PaidAt          time.Time
	ConfirmedAt     time.Time
	CancelledAt     time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	DeletedAt       time.Time
}

type ReceiptLine struct {
	RefID  int64  `json:"ref_id"`
	Title  string `json:"title"`
	Amount int32  `json:"amount"`
}

type Trader struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}
