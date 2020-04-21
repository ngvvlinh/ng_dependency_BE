package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subscription_plan_interval"
	"o.o/capi/dot"
)

// +sqlgen
type SubscriptionPlan struct {
	ID            dot.ID
	Name          string
	Price         int
	Status        status3.Status
	Description   string
	ProductID     dot.ID
	Interval      subscription_plan_interval.SubscriptionPlanInterval
	IntervalCount int
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	DeletedAt     time.Time
	WLPartnerID   dot.ID
}
