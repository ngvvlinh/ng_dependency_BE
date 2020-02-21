package query

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
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

func (q *CustomerQuery) GetCustomer(
	ctx context.Context, args *customering.GetCustomerArgs,
) (*customering.ShopCustomer, error) {
	query := q.store(ctx).ShopID(args.ShopID)

	counter := 0
	if args.ID.Int64() != 0 {
		query.ID(args.ID)
		counter++
	}
	if args.Code != "" {
		query.Code(args.Code)
		counter++
	}
	if args.ExternalID != "" {
		query.ExternalID(args.ExternalID)
		counter++
	}
	if counter == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Arguments are invalid", nil)
	}

	customer, err := query.GetCustomer()
	if err != nil {
		return nil, err
	}

	customerGroups, err := q.customerGroupCustomerStore(ctx).CustomerID(customer.ID).ListShopCustomerGroupsCustomerByCustomerID()
	if err != nil {
		return nil, err
	}
	for _, customerGroup := range customerGroups {
		customer.GroupIDs = append(customer.GroupIDs, customerGroup.GroupID)
	}
	return customer, nil

}

func (q *CustomerQuery) GetCustomerByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ID(args.ID).GetCustomer()
	if err != nil {
		return nil, err
	}
	if args.ShopID != 0 && customer.ShopID != 0 && customer.ShopID != args.ShopID {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "khách hàng không thuộc cửa hàng")
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
	ctx context.Context, code string, shopID dot.ID,
) (*customering.ShopCustomer, error) {
	return q.store(ctx).ShopID(shopID).Code(code).GetCustomer()
}

func (q *CustomerQuery) ListCustomers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*customering.CustomersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).WithPaging(args.Paging).Filters(args.Filters)
	customers, err := query.ListCustomers()
	if err != nil {
		return nil, err
	}
	var mapCustomerGroup = make(map[dot.ID][]dot.ID)
	var customerIDs []dot.ID
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
	return &customering.CustomersResponse{
		Customers: customers,
		Paging:    query.GetPaging(),
	}, nil
}

func (q *CustomerQuery) ListCustomersByIDs(
	ctx context.Context, args *customering.ListCustomerByIDsArgs,
) (*customering.CustomersResponse, error) {
	var shopIDs []dot.ID
	count := 0
	if args.ShopID != 0 {
		shopIDs = append(shopIDs, args.ShopID)
		count++
	}
	if args.ShopIDs != nil {
		shopIDs = append(shopIDs, args.ShopIDs...)
		count++
	}
	if count != 1 {
		return nil, cm.Error(cm.InvalidArgument, "Yêu cầu không hợp lệ", nil)
	}
	query := q.store(ctx).ShopIDs(shopIDs...)
	if len(args.IDs) == 0 {
		query = query.WithPaging(args.Paging)
	} else {
		query = query.IDs(args.IDs...)
	}
	customers, err := query.ListCustomers()
	if err != nil {
		return nil, err
	}
	return &customering.CustomersResponse{
		Customers: customers,
		Paging:    query.GetPaging(),
	}, nil
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
	query := q.customerGroupStore(ctx).WithPaging(args.Paging).Filters(args.Filters)
	customerGroup, err := query.ListShopCustomerGroups()
	if err != nil {
		return nil, err
	}
	return &customering.CustomerGroupsResponse{
		CustomerGroups: customerGroup,
		Paging:         query.GetPaging(),
	}, nil
}

func (q *CustomerQuery) GetCustomerByPhone(
	ctx context.Context, phone string, shopID dot.ID,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ShopID(shopID).Phone(phone).GetCustomer()
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (q *CustomerQuery) GetCustomerByEmail(
	ctx context.Context, email string, shopID dot.ID,
) (*customering.ShopCustomer, error) {
	customer, err := q.store(ctx).ShopID(shopID).Email(email).GetCustomer()
	if err != nil {
		return nil, err
	}
	return customer, nil
}
func (q *CustomerQuery) GetCustomerIndependent(context.Context, *meta.Empty) (*customering.ShopCustomer, error) {
	return &customering.ShopCustomer{
		ID:       customering.CustomerAnonymous,
		FullName: "Khách lẻ",
		Type:     customer_type.Independent,
	}, nil
}

func (q *CustomerQuery) ListCustomerGroupsCustomers(ctx context.Context, args *customering.ListCustomerGroupsCustomersArgs) (*customering.CustomerGroupsCustomersResponse, error) {
	query := q.customerGroupCustomerStore(ctx).WithPaging(args.Paging)
	count := 0
	if len(args.GroupIDs) != 0 {
		query = query.IDs(args.GroupIDs...)
		count++
	}
	if len(args.CustomerIDs) != 0 {
		query = query.CustomerIDs(args.CustomerIDs...)
		count++
	}
	if count != 1 {
		return nil, cm.Error(cm.FailedPrecondition, "Request không hợp lệ", nil)
	}
	customerGroupsCustomers, err := query.ListShopCustomerGroupsCustomer()
	if err != nil {
		return nil, err
	}
	var relationships []*customering.CustomerGroupCustomer
	for _, customerGroupCustomer := range customerGroupsCustomers {
		relationships = append(relationships, &customering.CustomerGroupCustomer{
			CustomerID: customerGroupCustomer.CustomerID,
			GroupID:    customerGroupCustomer.GroupID,
		})
	}
	return &customering.CustomerGroupsCustomersResponse{
		CustomerGroupsCustomers: relationships,
		Paging:                  query.GetPaging(),
	}, nil
}
