package aggregate

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/customering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/validate"
)

var _ customering.Aggregate = &CustomerAggregate{}

type CustomerAggregate struct {
	db                         *cmsql.Database
	store                      sqlstore.CustomerStoreFactory
	customerGroupStore         sqlstore.CustomerGroupStoreFactory
	customerGroupCustomerStore sqlstore.CustomerGroupCustomerStoreFactory
	addressStore               sqlstore.AddressStoreFactory
}

func NewCustomerAggregate(db *cmsql.Database) *CustomerAggregate {
	return &CustomerAggregate{
		db:                         db,
		store:                      sqlstore.NewCustomerStore(db),
		customerGroupStore:         sqlstore.NewCustomerGroupStore(db),
		customerGroupCustomerStore: sqlstore.NewCustomerGroupCustomerStore(db),
		addressStore:               sqlstore.NewAddressStore(db),
	}
}

func (a *CustomerAggregate) MessageBus() customering.CommandBus {
	b := bus.New()
	return customering.NewAggregateHandler(a).RegisterHandlers(b)
}

const (
	codeRegex  = "^KH[0-9]{6}$"
	codePrefix = "KH"
)

var reCode = regexp.MustCompile(codeRegex)

func (a *CustomerAggregate) CreateCustomer(
	ctx context.Context, args *customering.CreateCustomerArgs,
) (_ *customering.ShopCustomer, err error) {
	if args.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập tên đầy đủ")
	}
	if args.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập số điện thoại")
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
	if customer.Code != "" { // check code exists
		_, err := a.store(ctx).ShopID(args.ShopID).Code(customer.Code).GetCustomerDB()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Mã khách hàng đã tồn tại")
		case cm.NotFound:
			if codeNorm, err := convert.ParseCodeNorm(customer.Code); err {
				customer.CodeNorm = int32(codeNorm)
			}
		default:
			return nil, err
		}
	} else {
		var maxCodeNorm int32
		customerTemp, err := a.store(ctx).ShopID(args.ShopID).IncludeDeleted().GetCustomerByMaximumCodeNorm()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			maxCodeNorm = customerTemp.CodeNorm
		case cm.NotFound:
			// no-op
		default:
			return nil, err
		}

		if maxCodeNorm >= convert.MaxCodeNorm {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
		}
		codeNorm := maxCodeNorm + 1
		customer.Code = convert.GenerateCode(int(codeNorm))
		customer.CodeNorm = codeNorm
	}

	if true {
		fmt.Printf("%v", customer)
	}

	err = a.store(ctx).CreateCustomer(customer)
	if err != nil {
		if strings.Contains(err.Error(), "gender_type") {
			return nil, cm.Error(cm.InvalidArgument, `Giới tính chỉ nằm trong "male", "female", "other"`, err)
		}
	}
	// TODO: created_at, updated_at
	var customerResult customering.ShopCustomer
	return convert.Convert_customeringmodel_ShopCustomer_customering_ShopCustomer(customer, &customerResult), err
}

func (a *CustomerAggregate) UpdateCustomer(
	ctx context.Context, args *customering.UpdateCustomerArgs,
) (*customering.ShopCustomer, error) {
	customer, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetCustomer()
	if err != nil {
		return nil, err
	}
	// Verify phone
	if args.Phone.Valid {
		if args.Phone.String == "" {
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại không thể rỗng", err)
		} else {
			phone, isPhone := validate.NormalizePhone(args.Phone.String)
			if isPhone != true {
				return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng số điện thoại", nil)
			}
			args.Phone.String = phone.String()
		}
	}

	// Verify email
	if args.Email.Valid {
		if args.Email.String == "" {
			return nil, cm.Error(cm.InvalidArgument, "Email không thể rỗng", err)
		} else {
			email, isEmail := validate.NormalizeEmail(args.Email.String)
			if isEmail != true {
				return nil, cm.Error(cm.InvalidArgument, "Vui lòng nhập đúng định dạng email", nil)
			}
			args.Email.String = email.String()
		}
	}

	// Check phone
	if args.Phone.Valid && customer.Phone != args.Phone.String {
		_, err = a.store(ctx).ShopID(args.ShopID).Phone(args.Phone.String).GetCustomerDB()
		if err == nil {
			return nil, cm.Error(cm.InvalidArgument, "Số điện thoại đã tồn tại", err)
		}
	}

	// Check email
	if args.Email.Valid && customer.Email != args.Email.String {
		_, err = a.store(ctx).ShopID(args.ShopID).Phone(args.Email.String).GetCustomerDB()
		if err == nil {
			return nil, cm.Error(cm.InvalidArgument, "Email đã tồn tại", err)
		}
	}

	customer = convert.Apply_customering_UpdateCustomerArgs_customering_ShopCustomer(args, customer)
	customerModel := &model.ShopCustomer{}
	if err = scheme.Convert(customer, customerModel); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdateCustomerDB(customerModel)

	if args.Code.Valid && args.Code.String != "" {
		customerTemp, err := a.store(ctx).ShopID(args.ShopID).Code(args.Code.String).GetCustomerDB()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			if customerTemp.ID != customer.ID {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Mã khách hàng đã tồn tại")
			}
		case cm.NotFound:
			// no-op
		default:
			return nil, err
		}
		if codeNorm, err := convert.ParseCodeNorm(args.Code.String); err {
			customerModel.CodeNorm = int32(codeNorm)
		}
	}

	err = a.store(ctx).UpdateCustomerDB(customerModel)
	return customer, err
}

func (a *CustomerAggregate) DeleteCustomer(
	ctx context.Context, id int64, shopID int64,
) (deleted int, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var errTr error
		deleted, errTr = a.store(ctx).ID(id).ShopID(shopID).SoftDelete()
		if errTr != nil {
			return errTr
		}
		_, errTr = a.addressStore(ctx).ShopTraderID(shopID, id).SoftDelete()
		return errTr
	})
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

	out := &customering.ShopCustomerGroup{}
	err = scheme.Convert(group, out)
	return out, err
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
