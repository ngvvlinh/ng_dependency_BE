package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	subcriptingsharemodel "o.o/backend/com/subscripting/sharemodel"
	"o.o/capi/dot"
)

// +sqlgen
type Subscription struct {
	ID        dot.ID
	AccountID dot.ID
	// CancelAtPeriodEnd: hủy/ngưng subscription khi hết hạn
	CancelAtPeriodEnd    dot.NullBool
	CurrentPeriodEndAt   time.Time
	CurrentPeriodStartAt time.Time
	Status               status3.Status
	// BillingCycleAnchorAt: Determines the date of the first full invoice, and, for plans with month or year intervals, the day of the month for subsequent invoices.
	BillingCycleAnchorAt time.Time
	// StartedAt: The date and time when the subscription started
	StartedAt   time.Time
	Customer    *subcriptingsharemodel.CustomerInfo
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
	WLPartnerID dot.ID
	PlanIDs     []dot.ID
}

// +sqlgen
type SubscriptionLine struct {
	ID             dot.ID
	PlanID         dot.ID
	SubscriptionID dot.ID
	Quantity       int
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"update"`
}

type SubscriptionFtLine struct {
	*Subscription
	Lines []*SubscriptionLine
}
