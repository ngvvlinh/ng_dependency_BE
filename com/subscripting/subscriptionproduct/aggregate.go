package subscriptionproduct

import (
	"context"

	"o.o/api/subscripting/subscriptionproduct"
	"o.o/backend/com/subscripting/subscriptionproduct/convert"
	"o.o/backend/com/subscripting/subscriptionproduct/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscriptionproduct.Aggregate = &SubrProductAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SubrProductAggregate struct {
	subrProductStore sqlstore.SubrProductStoreFactory
}

func NewSubrProductAggregate(db *cmsql.Database) *SubrProductAggregate {
	return &SubrProductAggregate{
		subrProductStore: sqlstore.NewSubscriptionProductStore(db),
	}
}

func SubrProductAggregateMessageBus(a *SubrProductAggregate) subscriptionproduct.CommandBus {
	b := bus.New()
	return subscriptionproduct.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SubrProductAggregate) CreateSubrProduct(ctx context.Context, args *subscriptionproduct.CreateSubrProductArgs) (*subscriptionproduct.SubscriptionProduct, error) {
	if args.Type == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Subscription product type does not valid")
	}
	var subrProduct subscriptionproduct.SubscriptionProduct
	if err := scheme.Convert(args, &subrProduct); err != nil {
		return nil, err
	}
	subrProduct.ID = cm.NewID()
	if err := a.subrProductStore(ctx).CreateSubscription(&subrProduct); err != nil {
		return nil, err
	}
	return a.subrProductStore(ctx).ID(subrProduct.ID).GetSubrProduct()
}

func (a *SubrProductAggregate) UpdateSubrProduct(ctx context.Context, args *subscriptionproduct.UpdateSubrProductArgs) error {
	var subrProduct subscriptionproduct.SubscriptionProduct
	if err := scheme.Convert(args, &subrProduct); err != nil {
		return err
	}

	return a.subrProductStore(ctx).ID(args.ID).UpdateSubrProduct(&subrProduct)
}

func (a *SubrProductAggregate) DeleteSubrProduct(ctx context.Context, id dot.ID) error {
	_, err := a.subrProductStore(ctx).ID(id).GetSubrProduct()
	if err != nil {
		return err
	}
	_, err = a.subrProductStore(ctx).ID(id).SoftDelete()
	return err
}
