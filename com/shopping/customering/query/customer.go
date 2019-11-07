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
	store                      sqlstore.CustomerStoreFactory
	customerGroupStore         sqlstore.CustomerGroupStoreFactory
	customerGroupCustomerStore sqlstore.CustomerGroupCustomerStoreFactory
}

func NewCustomerQuery(db *cmsql.Database) *CustomerQuery {
	return &CustomerQuery{
		store:                      sqlstore.NewCustomerStore(db),
		customerGroupStore:         sqlstore.NewCustomerGroupStore(db),
		customerGroupCustomerStore: sqlstore.NewCustomerGroupCustomerStore(db),
	}
}

func (q *CustomerQuery) MessageBus() customering.QueryBus {
	b := bus.New()
	return customering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *CustomerQuery) GetCustomerByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetCustomer()
	if err != nil {
		return nil, err
	}
	q1 := q.customerGroupCustomerStore(ctx).CustomerID(args.ID)
	customerGroups, err := q1.ListShopCustomerGroupsCustomerByCustomerID()
	if err != nil {
		return nil, err
	}
	for _, customerGroup := range customerGroups {
		customer.GroupIDs = append(customer.GroupIDs, customerGroup.GroupID)
	}
	return customer, err
}

func (q *CustomerQuery) GetCustomerByCode(
	ctx context.Context, code string, shopID int64,
) (*customering.ShopCustomer, error) {
	return q.store(ctx).ShopID(shopID).Code(code).GetCustomer()
}

func (q *CustomerQuery) ListCustomers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*customering.CustomersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Paging(args.Paging).Filters(args.Filters)
	customers, err := query.ListCustomers()
	if err != nil {
		return nil, err
	}
	var mapCustomerGroup = make(map[int64][]int64)
	var customerIDs []int64
	for _, customer := range customers {
		customerIDs = append(customerIDs, customer.ID)
	}
	customerGroups, err := q.customerGroupCustomerStore(ctx).CustomerIDs(customerIDs...).ListShopCustomerGroupsCustomer()
	if err != nil {
		return nil, err
	}
	for _, customerGroup := range customerGroups {
		mapCustomerGroup[customerGroup.CustomerID] = append(mapCustomerGroup[customerGroup.CustomerID], customerGroup.GroupID)
	}
	for _, customer := range customers {
		customer.GroupIDs = mapCustomerGroup[customer.ID]
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

func (q *CustomerQuery) GetCustomerGroup(ctx context.Context, args *customering.GetCustomerGroupArgs) (*customering.ShopCustomerGroup, error) {
	customerGroup, err := q.customerGroupStore(ctx).ID(args.ID).GetShopCustomerGroup()
	if err != nil {
		return nil, err
	}
	return customerGroup, err
}

func (q *CustomerQuery) ListCustomerGroups(
	ctx context.Context, args *customering.ListCustomerGroupArgs,
) (*customering.CustomerGroupsResponse, error) {
	query := q.customerGroupStore(ctx).Paging(args.Paging).Filters(args.Filters)
	customerGroup, err := query.ListShopCustomerGroups()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &customering.CustomerGroupsResponse{
		CustomerGroups: customerGroup,
		Count:          int32(count),
		Paging:         query.GetPaging(),
	}, nil
}

func (q *CustomerQuery) GetCustomerByPhone(
	ctx context.Context, phone string, shopID int64,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ShopID(shopID).Phone(phone).GetCustomer()
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (q *CustomerQuery) GetCustomerByEmail(
	ctx context.Context, email string, shopID int64,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ShopID(shopID).Email(email).GetCustomer()
	if err != nil {
		return nil, err
	}
	return customer, nil
}
