package subscriptionplan

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subscription_plan_interval"
	"o.o/capi/dot"
)

type SubscriptionPlan struct {
	ID          dot.ID
	Name        string
	Price       int
	Status      status3.Status
	Description string
	ProductID   dot.ID
	Interval    subscription_plan_interval.SubscriptionPlanInterval
	// IntervalCount: The number of intervals (specified in the interval attribute) between subscription billings. For example, interval=month and interval_count=3 bills every 3 months.
	IntervalCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	WLPartnerID   dot.ID
}
