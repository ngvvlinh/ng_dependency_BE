package model

import (
	"time"

	"o.o/api/top/types/etc/status4"
	subcriptingsharemodel "o.o/backend/com/subscripting/sharemodel"
	"o.o/capi/dot"
)

// +sqlgen
type SubscriptionBill struct {
	ID             dot.ID
	AccountID      dot.ID
	SubscriptionID dot.ID
	TotalAmount    int
	Description    string
	PaymentID      dot.ID
	PaymentStatus  status4.Status
	Status         status4.Status
	Customer       *subcriptingsharemodel.CustomerInfo
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"create"`
	DeletedAt      time.Time
	WLPartnerID    dot.ID
}

// +sqlgen
type SubscriptionBillLine struct {
	ID                 dot.ID
	LineAmount         int
	Price              int
	Quantity           int
	Description        string
	PeriodStartAt      time.Time
	PeriodEndAt        time.Time
	SubscriptionBillID dot.ID
	SubscriptionID     dot.ID
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
}

type SubscriptionBillFtLine struct {
	*SubscriptionBill
	Lines []*SubscriptionBillLine
}
