package aggregate

import (
	"context"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/backend/pkg/shopping/customering/model"

	"etop.vn/common/bus"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/shopping/customering/convert"
	"etop.vn/backend/pkg/shopping/customering/sqlstore"
)

var _ customering.Aggregate = &Aggregate{}

type Aggregate struct {
	store sqlstore.CustomerStoreFactory
}

func New(db cmsql.Database) *Aggregate {
	return &Aggregate{
		store: sqlstore.NewCustomerStore(db),
	}
}

func (q *Aggregate) MessageBus() customering.CommandBus {
	b := bus.New()
	return customering.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a *Aggregate) CreateCustomer(
	ctx context.Context, args *customering.CreateCustomerArgs,
) (*customering.ShopCustomer, error) {
	customer := convert.CreateShopCustomer(args)
	err := a.store(ctx).CreateCustomer(customer)
	// TODO: created_at, updated_at
	return customer, err
}

func (a *Aggregate) UpdateCustomer(
	ctx context.Context, args *customering.UpdateCustomerArgs,
) (*customering.ShopCustomer, error) {
	customer, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetCustomer()
	if err != nil {
		return nil, err
	}
	updated := convert.UpdateShopCustomer(customer, args)
	err = a.store(ctx).UpdateCustomerDB(convert.ShopCustomerDB(updated))
	return updated, err
}

func (a *Aggregate) DeleteCustomer(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*meta.Empty, error) {
	return nil, cm.ErrTODO
}

func (a *Aggregate) BatchSetCustomersStatus(
	ctx context.Context, args *customering.BatchSetCustomersStatusArgs,
) (*meta.UpdatedResponse, error) {
	update := &model.ShopCustomer{Status: args.Status}
	n, err := a.store(ctx).IDs(args.IDs...).PatchCustomerDB(update)
	return &meta.UpdatedResponse{Updated: int32(n)}, err
}
