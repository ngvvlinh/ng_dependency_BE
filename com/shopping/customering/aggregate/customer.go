package aggregate

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ customering.Aggregate = &CustomerAggregate{}

type CustomerAggregate struct {
	store sqlstore.CustomerStoreFactory
}

func NewCustomerAggregate(db cmsql.Database) *CustomerAggregate {
	return &CustomerAggregate{
		store: sqlstore.NewCustomerStore(db),
	}
}

func (a *CustomerAggregate) MessageBus() customering.CommandBus {
	b := bus.New()
	return customering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *CustomerAggregate) CreateCustomer(
	ctx context.Context, args *customering.CreateCustomerArgs,
) (*customering.ShopCustomer, error) {
	customer := convert.CreateShopCustomer(args)
	err := a.store(ctx).CreateCustomer(customer)
	// TODO: created_at, updated_at
	return customer, err
}

func (a *CustomerAggregate) UpdateCustomer(
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

func (a *CustomerAggregate) DeleteCustomer(
	ctx context.Context, id int64, shopID int64,
) (deleted int, _ error) {
	deleted, err := a.store(ctx).ID(id).ShopID(shopID).SoftDelete()
	return deleted, err
}

func (a *CustomerAggregate) BatchSetCustomersStatus(
	ctx context.Context, ids []int64, shopID int64, status int32,
) (*meta.UpdatedResponse, error) {
	update := &model.ShopCustomer{Status: status}
	n, err := a.store(ctx).IDs(ids...).PatchCustomerDB(update)
	return &meta.UpdatedResponse{Updated: int32(n)}, err
}
