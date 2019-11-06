package aggregate

import (
	"context"

	"etop.vn/api/shopping/suppliering"
	"etop.vn/backend/com/shopping/suppliering/convert"
	"etop.vn/backend/com/shopping/suppliering/model"
	"etop.vn/backend/com/shopping/suppliering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
)

var _ suppliering.Aggregate = &SupplierAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SupplierAggregate struct {
	store sqlstore.SupplierStoreFactory
}

func NewSupplierAggregate(db *cmsql.Database) *SupplierAggregate {
	return &SupplierAggregate{
		store: sqlstore.NewSupplierStore(db),
	}
}

func (a *SupplierAggregate) MessageBus() suppliering.CommandBus {
	b := bus.New()
	return suppliering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SupplierAggregate) CreateSupplier(
	ctx context.Context, args *suppliering.CreateSupplierArgs,
) (*suppliering.ShopSupplier, error) {
	supplier := new(suppliering.ShopSupplier)
	if err := scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	err := a.store(ctx).CreateSupplier(supplier)
	return supplier, err
}

func (a *SupplierAggregate) UpdateSupplier(
	ctx context.Context, args *suppliering.UpdateSupplierArgs,
) (*suppliering.ShopSupplier, error) {
	supplier, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetSupplier()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(args, supplier); err != nil {
		return nil, err
	}
	supplierDB := new(model.ShopSupplier)
	if err := scheme.Convert(supplier, supplierDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateSupplierDB(supplierDB)
	return supplier, err
}

func (a *SupplierAggregate) DeleteSupplier(
	ctx context.Context, ID int64, shopID int64,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return deleted, err
}
