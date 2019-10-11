package aggregate

import (
	"context"
	"strings"

	"etop.vn/backend/pkg/common/validate"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ customering.Aggregate = &CustomerAggregate{}

type CustomerAggregate struct {
	db                         cmsql.Database
	store                      sqlstore.CustomerStoreFactory
	customerGroupStore         sqlstore.CustomerGroupStoreFactory
	customerGroupCustomerStore sqlstore.CustomerGroupCustomerStoreFactory
}

func NewCustomerAggregate(db cmsql.Database) *CustomerAggregate {
	return &CustomerAggregate{
		db:                         db,
		store:                      sqlstore.NewCustomerStore(db),
		customerGroupStore:         sqlstore.NewCustomerGroupStore(db),
		customerGroupCustomerStore: sqlstore.NewCustomerGroupCustomerStore(db),
	}
}

func (a *CustomerAggregate) MessageBus() customering.CommandBus {
	b := bus.New()
	return customering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *CustomerAggregate) CreateCustomer(
	ctx context.Context, args *customering.CreateCustomerArgs,
) (*customering.ShopCustomer, error) {
	if args.FullName == "" {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập tên đầy đủ", nil)
	}
	if args.Phone == "" {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập số điện thoại", nil)
	}

	phone, isPhone := validate.NormalizePhone(args.Phone)
	if isPhone != true {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
	}
	args.Phone = phone.String()

	if args.Email != "" {
		email, isEmail := validate.NormalizeEmail(args.Email)
		if isEmail != true {
			return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
		}
		args.Email = email.String()
	}
	customer := convert.CreateShopCustomer(args)
	err := a.store(ctx).CreateCustomer(customer)
	if err != nil {
		if strings.Contains(err.Error(), "gender_type") {
			return nil, cm.Error(cm.InvalidArgument, `Giới tính chỉ nằm trong "male", "female", "other"`, err)
		}
	}
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

func (a *CustomerAggregate) CreateCustomerGroup(ctx context.Context, args *customering.CreateCustomerGroupArgs) (*customering.ShopCustomerGroup, error) {
	customerGroup := &customering.ShopCustomerGroup{
		ID:   cm.NewID(),
		Name: args.Name,
	}
	if err := a.customerGroupStore(ctx).CreateShopCustomerGroup(customerGroup); err != nil {
		return nil, err
	}
	return customerGroup, nil
}

func (a *CustomerAggregate) UpdateCustomerGroup(ctx context.Context, args *customering.UpdateCustomerGroupArgs) (*customering.ShopCustomerGroup, error) {
	customerGroupDB, err := a.customerGroupStore(ctx).ID(args.ID).GetShopCustomerGroupDB()
	if err != nil {
		return nil, err
	}
	group := convert.UpdateCustomerGroup(customerGroupDB, args)
	if err := a.customerGroupStore(ctx).UpdateCustomerGroup(group); err != nil {
		return nil, err
	}

	var out customering.ShopCustomerGroup
	convert.ShopCustomerGroup(group, &out)
	return &out, nil
}

func (a *CustomerAggregate) AddCustomersToGroup(ctx context.Context, args *customering.AddCustomerToGroupArgs) (created int, _ error) {
	var err error
	if len(args.CustomerIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "CustomerIDs không được để trống")
	}
	if args.GroupID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "GroupID không được để trống")
	}
	_, err = a.customerGroupStore(ctx).ID(args.GroupID).GetShopCustomerGroup()
	if err != nil {
		return 0, cm.Errorf(cm.InvalidArgument, err, "Group không tồn tại")
	}
	for Customer := range args.CustomerIDs {
		if args.CustomerIDs[Customer] == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, err, "CustomerID không được để trống")
		}
	}
	customers, err := a.store(ctx).IDs(args.CustomerIDs...).ListCustomers()
	if err != nil {
		return 0, err
	}
	if len(customers) != len(args.CustomerIDs) {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Thành viên không tồn tại")
	}
	err = a.db.InTransaction(ctx, func(q cmsql.QueryInterface) error {
		for _, customerID := range args.CustomerIDs {
			customerGroup := &customering.ShopCustomerGroupCustomer{
				CustomerID: customerID,
				GroupID:    args.GroupID,
			}
			lineCreated, err := a.customerGroupCustomerStore(ctx).AddShopCustomerToGroup(customerGroup)
			if err != nil {
				return err
			}
			created += lineCreated
		}
		return nil
	})
	return created, err
}

func (a *CustomerAggregate) RemoveCustomersFromGroup(ctx context.Context, args *customering.RemoveCustomerOutOfGroupArgs) (deleted int, _ error) {

	var err error
	var removedCustomerGroup int
	if len(args.CustomerIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "CustomerIDs không được để trống")
	}
	if args.GroupID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, err, "GroupID không được để trống")
	}
	for customerID := range args.CustomerIDs {
		if args.CustomerIDs[customerID] == 0 {
			return 0, cm.Errorf(cm.InvalidArgument, err, "CustomerID không được để trống")
		}
	}
	removedCustomerGroup, err = a.customerGroupCustomerStore(ctx).CustomerIDs(args.CustomerIDs...).RemoveCustomerFromGroup()
	return removedCustomerGroup, err
}
