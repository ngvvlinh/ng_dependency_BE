package subscriptionplan

import (
	"context"

	"o.o/api/subscripting/subscriptionplan"
	com "o.o/backend/com/main"
	"o.o/backend/com/subscripting/subscriptionplan/convert"
	"o.o/backend/com/subscripting/subscriptionplan/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/capi/dot"
)

var _ subscriptionplan.Aggregate = &SubrPlanAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SubrPlanAggregate struct {
	subrPlanStore sqlstore.SubrPlanStoreFactory
}

func NewSubrPlanAggregate(db com.MainDB) *SubrPlanAggregate {
	return &SubrPlanAggregate{
		subrPlanStore: sqlstore.NewSubrPlanStore(db),
	}
}

func SubrPlanAggregateMessageBus(a *SubrPlanAggregate) subscriptionplan.CommandBus {
	b := bus.New()
	return subscriptionplan.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SubrPlanAggregate) CreateSubrPlan(ctx context.Context, args *subscriptionplan.CreateSubrPlanArgs) (*subscriptionplan.SubscriptionPlan, error) {
	if args.ProductID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription product ID")
	}
	if args.Interval == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription plan interval")
	}
	if args.IntervalCount == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription plan interval count")
	}

	var subrPlan subscriptionplan.SubscriptionPlan
	if err := scheme.Convert(args, &subrPlan); err != nil {
		return nil, err
	}
	subrPlan.ID = cm.NewID()
	if err := a.subrPlanStore(ctx).CreateSubscription(&subrPlan); err != nil {
		return nil, err
	}
	return a.subrPlanStore(ctx).ID(subrPlan.ID).GetSubrPlan()
}

func (a *SubrPlanAggregate) UpdateSubrPlan(ctx context.Context, args *subscriptionplan.UpdateSubrPlanArgs) error {
	var subrPlan subscriptionplan.SubscriptionPlan
	if err := scheme.Convert(args, &subrPlan); err != nil {
		return err
	}
	return a.subrPlanStore(ctx).ID(args.ID).UpdateSubrPlan(&subrPlan)
}

func (a *SubrPlanAggregate) DeleteSubrPlan(ctx context.Context, id dot.ID) error {
	_, err := a.subrPlanStore(ctx).ID(id).GetSubrPlan()
	if err != nil {
		return err
	}
	_, err = a.subrPlanStore(ctx).ID(id).SoftDelete()
	return err
}
