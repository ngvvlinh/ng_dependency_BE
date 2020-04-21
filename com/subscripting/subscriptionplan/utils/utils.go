package utils

import (
	"time"

	"o.o/api/top/types/etc/subscription_plan_interval"
	cm "o.o/backend/pkg/common"
)

const (
	oneYear  = 365 * 24 * time.Hour
	oneMonth = 30 * 24 * time.Hour
	oneWeek  = 7 * 24 * time.Hour
	oneDay   = 24 * time.Hour
)

func CalcPlanDuration(interval subscription_plan_interval.SubscriptionPlanInterval, intervalCount int) (time.Duration, error) {
	switch interval {
	case subscription_plan_interval.Day:
		return time.Duration(intervalCount) * oneDay, nil
	case subscription_plan_interval.Week:
		return time.Duration(intervalCount) * oneWeek, nil
	case subscription_plan_interval.Month:
		return time.Duration(intervalCount) * oneMonth, nil
	case subscription_plan_interval.Year:
		return time.Duration(intervalCount) * oneYear, nil
	default:
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Plan does not valid")
	}
}
