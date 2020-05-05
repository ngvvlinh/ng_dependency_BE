package subscriptionplan

import (
	"context"

	"o.o/api/top/types/etc/subscription_plan_interval"
	"o.o/api/top/types/etc/subscription_product_type"
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
	ListSubrPlans(context.Context, *ListSubrPlansArgs) ([]*SubscriptionPlan, error)
	GetFreeSubrPlanByProductType(ctx context.Context, ProductType subscription_product_type.ProductSubscriptionType) (*SubscriptionPlan, error)
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

type ListSubrPlansArgs struct {
	ProductIDs []dot.ID
}
