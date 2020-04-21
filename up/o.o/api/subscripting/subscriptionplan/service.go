package subscriptionplan

import (
	"context"

	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/subscription_plan_interval"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSubrPlan(context.Context, *CreateSubrPlanArgs) (*SubscriptionPlan, error)
	UpdateSubrPlan(context.Context, *UpdateSubrPlanArgs) error
	DeleteSubrPlan(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetSubrPlanByID(ctx context.Context, ID dot.ID) (*SubscriptionPlan, error)
	ListSubrPlans(context.Context, *cm.Empty) ([]*SubscriptionPlan, error)
}

// +convert:create=SubscriptionPlan
type CreateSubrPlanArgs struct {
	Name          string
	Price         int
	Description   string
	ProductID     dot.ID
	Interval      subscription_plan_interval.SubscriptionPlanInterval
	IntervalCount int
}

// +convert:update=SubscriptionPlan(ID)
type UpdateSubrPlanArgs struct {
	ID            dot.ID
	Name          string
	Price         int
	Description   string
	ProductID     dot.ID
	Interval      subscription_plan_interval.SubscriptionPlanInterval
	IntervalCount int
}
