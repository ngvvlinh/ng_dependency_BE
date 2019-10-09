package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ customering.QueryService = &CustomerQuery{}

type CustomerQuery struct {
	store sqlstore.CustomerStoreFactory
}

func NewCustomerQuery(db cmsql.Database) *CustomerQuery {
	return &CustomerQuery{
		store: sqlstore.NewCustomerStore(db),
	}
}

func (q *CustomerQuery) MessageBus() customering.QueryBus {
	b := bus.New()
	return customering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *CustomerQuery) GetCustomerByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*customering.ShopCustomer, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetCustomer()
}

func (q *CustomerQuery) ListCustomers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*customering.CustomersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Paging(args.Paging).Filters(args.Filters)
	customers, err := query.ListCustomers()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &customering.CustomersResponse{
		Customers: customers,
		Count:     int32(count),
		Paging:    query.GetPaging(),
	}, nil
}

func (q *CustomerQuery) ListCustomersByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*customering.CustomersResponse, error) {
	customers, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListCustomers()
	if err != nil {
		return nil, err
	}
	return &customering.CustomersResponse{Customers: customers}, nil
}
