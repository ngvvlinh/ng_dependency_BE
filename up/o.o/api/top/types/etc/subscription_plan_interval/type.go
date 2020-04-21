package subscription_plan_interval

// +enum
// +enum:zero=null
type SubscriptionPlanInterval int

type NullSubscriptionPlanInterval struct {
	Enum  SubscriptionPlanInterval
	Valid bool
}

const (
	// +enum=unknown
	Unknown SubscriptionPlanInterval = 0

	// +enum=day
	Day SubscriptionPlanInterval = 1

	// +enum=week
	Week SubscriptionPlanInterval = 2

	// +enum=month
	Month SubscriptionPlanInterval = 3

	// +enum=year
	Year SubscriptionPlanInterval = 4
)
