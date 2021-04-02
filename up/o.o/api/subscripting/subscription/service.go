package subscription

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateSubscription(context.Context, *CreateSubscriptionArgs) (*SubscriptionFtLine, error)
	UpdateSubscriptionPeriod(context.Context, *UpdateSubscriptionPeriodArgs) error
	UpdateSubscripionStatus(context.Context, *UpdateSubscriptionStatusArgs) error
	UpdateSubscriptionInfo(context.Context, *UpdateSubscriptionInfoArgs) error
	CancelSubscription(ctx context.Context, ID dot.ID, AccountID dot.ID) error
	ActivateSubscription(ctx context.Context, ID dot.ID, AccountID dot.ID) error
	DeleteSubscription(ctx context.Context, ID dot.ID, AccountID dot.ID) error
}

type QueryService interface {
	GetSubscriptionByID(ctx context.Context, ID dot.ID, AccountID dot.ID) (*SubscriptionFtLine, error)
	ListSubscriptions(context.Context, *ListSubscriptionsArgs) (*ListSubscriptionsResponse, error)
	GetLastestSubscriptionByProductType(ctx context.Context, AccountID dot.ID, ProductType subscription_product_type.ProductSubscriptionType) (*SubscriptionFtLine, error)
}

// +convert:create=Subscription
type CreateSubscriptionArgs struct {
	AccountID            dot.ID
	CancelAtPeriodEnd    bool
	CurrentPeriodEndAt   time.Time
	CurrentPeriodStartAt time.Time
	Lines                []*SubscriptionLine
	BillingCycleAnchorAt time.Time
	Customer             *types.CustomerInfo
}

// +convert:update=Subscription(ID,AccountID)
type UpdateSubscriptionPeriodArgs struct {
	ID                   dot.ID
	AccountID            dot.ID
	CancelAtPeriodEnd    bool
	CurrentPeriodStartAt time.Time
	CurrentPeriodEndAt   time.Time
	BillingCycleAnchorAt time.Time
	StartedAt            time.Time
}

type UpdateSubscriptionStatusArgs struct {
	ID        dot.ID
	AccountID dot.ID
	Status    status3.NullStatus
}

type UpdateSubscriptionInfoArgs struct {
	ID                   dot.ID
	AccountID            dot.ID
	CancelAtPeriodEnd    dot.NullBool
	BillingCycleAnchorAt time.Time
	Customer             *types.CustomerInfo
	Lines                []*SubscriptionLine
}

type ListSubscriptionsArgs struct {
	AccountID dot.ID
	Paging    meta.Paging
	Filters   meta.Filters
}

type ListSubscriptionsResponse struct {
	Subscriptions []*SubscriptionFtLine
	Paging        meta.PageInfo
}
