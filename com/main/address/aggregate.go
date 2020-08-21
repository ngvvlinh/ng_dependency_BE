package address

import (
	"context"

	"o.o/api/main/address"
	"o.o/api/main/location"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/address/convert"
	"o.o/backend/com/main/address/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()
var _ address.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	store      sqlstore.AddressFactory
	eventBus   capi.EventBus
	locationQS location.QueryBus
}

func NewAggregateAddress(
	bus capi.EventBus,
	db com.MainDB,
	locationQS location.QueryBus,
) *Aggregate {
	return &Aggregate{
		eventBus:   bus,
		store:      sqlstore.NewAddressStore(db),
		locationQS: locationQS,
	}
}

func AddressAggregateMessageBus(q *Aggregate) address.CommandBus {
	b := bus.New()
	return address.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a *Aggregate) ValidateLocation(ctx context.Context, in *address.Address) error {
	locationQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: in.ProvinceCode,
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
		Province:     in.Province,
		District:     in.District,
		Ward:         in.Ward,
	}
	if err := a.locationQS.Dispatch(ctx, locationQuery); err != nil {
		return err
	}
	loc := locationQuery.Result
	if loc.Province == nil || loc.District == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp thông tin tỉnh/thành phố và quận/huyện")
	}

	in.Province = loc.Province.Name
	in.ProvinceCode = loc.Province.Code
	in.District = loc.District.Name
	in.DistrictCode = loc.District.Code
	if loc.Ward != nil {
		in.Ward = loc.Ward.Name
		in.WardCode = loc.Ward.Code
	}

	return nil
}

func (a *Aggregate) CreateAddress(ctx context.Context, args *address.CreateAddressArgs) (*address.Address, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	phoneNorm, ok := validate.NormalizePhone(args.Phone)
	if !ok {
		return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
	}

	args.Phone = phoneNorm.String()

	addressCore := &address.Address{}

	if err := scheme.Convert(args, addressCore); err != nil {
		return nil, err
	}

	addressCore.ID = cm.NewID()

	// validate location
	if err := a.ValidateLocation(ctx, addressCore); err != nil {
		return nil, err
	}

	err := a.store(ctx).CreateAddress(addressCore)

	return addressCore, err
}

func (a *Aggregate) UpdateAddress(ctx context.Context, args *address.UpdateAddressArgs) (*address.Address, error) {
	flag, err := args.Validate()
	if err != nil {
		return nil, err
	}
	if args.Phone != "" {
		phoneNorm, ok := validate.NormalizePhone(args.Phone)
		if !ok {
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không hợp lệ", nil)
		}

		args.Phone = phoneNorm.String()
	}

	addressCore := &address.Address{}

	if err := scheme.Convert(args, addressCore); err != nil {
		return nil, err
	}

	if flag == 1 { // validate if there is information location
		if err := a.ValidateLocation(ctx, addressCore); err != nil {
			return nil, err
		}
	}

	return a.store(ctx).UpdateAddress(addressCore)
}

func (a *Aggregate) RemoveAddress(ctx context.Context, q *address.DeleteAddressArgs) error {
	if err := q.Validate(); err != nil {
		return err
	}
	return a.store(ctx).ID(q.ID).AccountID(q.AccountID).Delete()
}
