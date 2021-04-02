package subscription

import (
	"time"

	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

type Subscription struct {
	ID                   dot.ID
	AccountID            dot.ID
	CancelAtPeriodEnd    dot.NullBool
	CurrentPeriodEndAt   time.Time
	CurrentPeriodStartAt time.Time
	Status               status3.Status
	// BillingCycleAnchorAt: Determines the date of the first full invoice, and, for plans with month or year intervals, the day of the month for subsequent invoices.
	BillingCycleAnchorAt time.Time
	StartedAt            time.Time
	Customer             *types.CustomerInfo
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	WLPartnerID          dot.ID
	PlanIDs              []dot.ID
}

type SubscriptionLine struct {
	ID             dot.ID
	PlanID         dot.ID
	SubscriptionID dot.ID
	Quantity       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SubscriptionFtLine struct {
	*Subscription
	Lines []*SubscriptionLine
}
