package location

import (
	"context"

	"o.o/api/main/location"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/location/convert"
	"o.o/backend/com/main/location/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/capi/dot"
)

var _ location.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	customRegionStore sqlstore.CustomRegionFactory
}

func NewAggregate(db com.MainDB) *Aggregate {
	return &Aggregate{
		customRegionStore: sqlstore.NewCustomRegionStore(db),
	}
}

func AggregateMessageBus(a *Aggregate) location.CommandBus {
	b := bus.New()
	return location.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateCustomRegion(ctx context.Context, args *location.CreateCustomRegionArgs) (*location.CustomRegion, error) {
	if len(args.ProvinceCodes) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing province_codes")
	}
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	var region location.CustomRegion
	if err := scheme.Convert(args, &region); err != nil {
		return nil, err
	}
	return a.customRegionStore(ctx).CreateCustomRegion(&region)
}

func (a *Aggregate) UpdateCustomRegion(ctx context.Context, args *location.CustomRegion) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return a.customRegionStore(ctx).UpdateCustomRegion(args)
}

func (a *Aggregate) DeleteCustomRegion(ctx context.Context, id dot.ID) error {
	_, err := a.customRegionStore(ctx).ID(id).SoftDelete()
	return err
}
