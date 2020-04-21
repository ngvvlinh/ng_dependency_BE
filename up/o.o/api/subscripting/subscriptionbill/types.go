package subscriptionbill

import (
	"time"

	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/status4"
	"o.o/capi/dot"
)

// +gen:event:topic=event/subscription_bill

type SubscriptionBill struct {
	ID             dot.ID
	AccountID      dot.ID
	SubscriptionID dot.ID
	TotalAmount    int
	Description    string
	PaymentID      dot.ID
	PaymentStatus  status4.Status
	Status         status4.Status
	Customer       *types.CustomerInfo
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	WLPartnerID    dot.ID
}

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
	SubscriptionLineID dot.ID
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type SubscriptionBillFtLine struct {
	*SubscriptionBill
	Lines []*SubscriptionBillLine
}

type SubscriptionBillPaidEvent struct {
	ID        dot.ID
	AccountID dot.ID
}
