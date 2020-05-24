package subscription

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/subscripting/subscription/convert"
	"o.o/backend/com/subscripting/subscription/model"
	"o.o/backend/com/subscripting/subscription/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscription.Aggregate = &SubscriptionAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SubscriptionAggregate struct {
	db                    cmsql.Transactioner
	subscriptionStore     sqlstore.SubscriptionStoreFactory
	subscriptionLineStore sqlstore.SubscriptionLineStoreFactory
}

func NewSubscriptionAggregate(db *cmsql.Database) *SubscriptionAggregate {
	return &SubscriptionAggregate{
		db:                    db,
		subscriptionStore:     sqlstore.NewSubscriptionStore(db),
		subscriptionLineStore: sqlstore.NewSubscriptionLineStore(db),
	}
}

func SubscriptionAggregateMessageBus(a *SubscriptionAggregate) subscription.CommandBus {
	b := bus.New()
	return subscription.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SubscriptionAggregate) CreateSubscription(ctx context.Context, args *subscription.CreateSubscriptionArgs) (*subscription.SubscriptionFtLine, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account ID")
	}
	if len(args.Lines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription lines")
	}
	var planIDs []dot.ID
	for _, line := range args.Lines {
		if err := verifySubscriptionLine(line); err != nil {
			return nil, err
		}
		if !cm.IDsContain(planIDs, line.PlanID) {
			planIDs = append(planIDs, line.PlanID)
		}
	}

	var subr subscription.Subscription
	if err := scheme.Convert(args, &subr); err != nil {
		return nil, err
	}
	subr.PlanIDs = planIDs

	subrID := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		subr.ID = subrID
		subr.Status = status3.P
		if err := a.subscriptionStore(ctx).CreateSubscription(&subr); err != nil {
			return err
		}
		for _, line := range args.Lines {
			line.SubscriptionID = subr.ID
			if err := a.subscriptionLineStore(ctx).CreateSubscriptionLine(line); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.subscriptionStore(ctx).ID(subrID).GetSubscriptionFtLine()
}

func verifySubscriptionLine(line *subscription.SubscriptionLine) error {
	if line.PlanID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing plan ID in subscription line")
	}
	if line.Quantity <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing quantity in subscription line")
	}
	return nil
}

func (a *SubscriptionAggregate) UpdateSubscriptionPeriod(ctx context.Context, args *subscription.UpdateSubscriptionPeriodArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID")
	}
	var update subscription.Subscription
	if err := scheme.Convert(args, &update); err != nil {
		return err
	}
	return a.subscriptionStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateSubscription(&update)
}

func (a *SubscriptionAggregate) UpdateSubscripionStatus(ctx context.Context, args *subscription.UpdateSubscriptionStatusArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID")
	}
	if !args.Status.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription status")
	}
	return a.subscriptionStore(ctx).UpdateSubscriptionStatus(args.ID, args.AccountID, args.Status.Enum)
}

func (a *SubscriptionAggregate) UpdateSubscriptionInfo(ctx context.Context, args *subscription.UpdateSubscriptionInfoArgs) error {
	subr, err := a.subscriptionStore(ctx).ID(args.ID).AccountID(args.AccountID).GetSubscriptionFtLine()
	if err != nil {
		return err
	}
	lineIDs := make([]dot.ID, len(subr.Lines))
	for i, line := range subr.Lines {
		lineIDs[i] = line.ID
	}

	var planIDs []dot.ID
	for _, line := range args.Lines {
		if err := verifySubscriptionLine(line); err != nil {
			return err
		}
		if !cm.IDsContain(planIDs, line.PlanID) {
			planIDs = append(planIDs, line.PlanID)
		}
	}

	update := &model.Subscription{
		CancelAtPeriodEnd:    args.CancelAtPeriodEnd,
		BillingCycleAnchorAt: args.BillingCycleAnchorAt,
		Customer:             convert.Convert_subscriptingtypes_CustomerInfo_sharemodel_CustomerInfo(args.Customer, nil),
		PlanIDs:              planIDs,
	}
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.subscriptionStore(ctx).ID(args.ID).AccountID(args.AccountID).UpdateSubscriptionDB(update); err != nil {
			return err
		}
		if len(args.Lines) == 0 {
			return nil
		}

		// update subscription line
		for i, line := range args.Lines {
			if line.ID != 0 {
				if err := a.subscriptionLineStore(ctx).SubscriptionID(subr.ID).ID(line.ID).UpdateSubscriptionLine(line); err != nil {
					return err
				}
				lineIDs = append(lineIDs[:i], lineIDs[i+1:]...)
			} else {
				line.SubscriptionID = subr.ID
				if err := a.subscriptionLineStore(ctx).CreateSubscriptionLine(line); err != nil {
					return err
				}
			}
		}
		if len(lineIDs) > 0 {
			return a.subscriptionLineStore(ctx).IDs(lineIDs).DeleteSubscriptionLine()
		}
		return nil
	})
}

func (a *SubscriptionAggregate) CancelSubscription(ctx context.Context, id dot.ID, accountID dot.ID) error {
	update := &subscription.UpdateSubscriptionStatusArgs{
		ID:        id,
		AccountID: accountID,
		Status:    status3.N.Wrap(),
	}
	return a.UpdateSubscripionStatus(ctx, update)
}

func (a *SubscriptionAggregate) ActivateSubscription(ctx context.Context, id, accountID dot.ID) error {
	update := &subscription.UpdateSubscriptionStatusArgs{
		ID:        id,
		AccountID: accountID,
		Status:    status3.P.Wrap(),
	}
	return a.UpdateSubscripionStatus(ctx, update)
}

func (a *SubscriptionAggregate) DeleteSubscription(ctx context.Context, id dot.ID, accountID dot.ID) error {
	_, err := a.subscriptionStore(ctx).ID(id).AccountID(accountID).GetSubscriptionFtLine()
	if err != nil {
		return err
	}
	_, err = a.subscriptionStore(ctx).ID(id).AccountID(accountID).SoftDelete()
	return err
}
