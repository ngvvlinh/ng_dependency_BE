package aggregate

import (
	"context"

	"etop.vn/backend/com/shopping/vendoring/model"

	"etop.vn/backend/pkg/common/scheme"

	"etop.vn/api/shopping/vendoring"
	"etop.vn/backend/com/shopping/vendoring/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ vendoring.Aggregate = &VendorAggregate{}

type VendorAggregate struct {
	store sqlstore.VendorStoreFactory
}

func NewVendorAggregate(db cmsql.Database) *VendorAggregate {
	return &VendorAggregate{
		store: sqlstore.NewVendorStore(db),
	}
}

func (a *VendorAggregate) MessageBus() vendoring.CommandBus {
	b := bus.New()
	return vendoring.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *VendorAggregate) CreateVendor(
	ctx context.Context, args *vendoring.CreateVendorArgs,
) (*vendoring.ShopVendor, error) {
	vendor := new(vendoring.ShopVendor)
	if err := scheme.Convert(args, vendor); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateVendor(vendor)
	return vendor, err
}

func (a *VendorAggregate) UpdateVendor(
	ctx context.Context, args *vendoring.UpdateVendorArgs,
) (*vendoring.ShopVendor, error) {
	vendor, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetVendor()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(args, vendor); err != nil {
		return nil, err
	}
	vendorDB := new(model.ShopVendor)
	if err := scheme.Convert(vendor, vendorDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateVendorDB(vendorDB)
	return vendor, err
}

func (a *VendorAggregate) DeleteVendor(
	ctx context.Context, ID int64, shopID int64,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return deleted, err
}
