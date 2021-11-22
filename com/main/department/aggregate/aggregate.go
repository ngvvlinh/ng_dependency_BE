package aggregate

import (
	"context"

	"o.o/api/main/department"
	"o.o/api/main/identity"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/department/convert"
	"o.o/backend/com/main/department/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
)

var _ department.Aggregate = &DepartmentAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type DepartmentAggregate struct {
	db              *cmsql.Database
	eventBus        capi.EventBus
	departmentStore sqlstore.DepartmentStoreFactory
	identityQuery   identity.QueryBus
}

func NewDepartmentAggregate(db com.MainDB, eventB capi.EventBus, identityQ identity.QueryBus) *DepartmentAggregate {
	return &DepartmentAggregate{
		db:              db,
		eventBus:        eventB,
		departmentStore: sqlstore.NewDepartmentStore(db),
		identityQuery:   identityQ,
	}
}

func DepartmentAggregateMessageBus(a *DepartmentAggregate) department.CommandBus {
	b := bus.New()
	return department.NewAggregateHandler(a).RegisterHandlers(b)
}

func (d *DepartmentAggregate) CreateDepartment(ctx context.Context, args *department.CreateDepartmentArgs) (*department.Department, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	oldDepartment, err := d.departmentStore(ctx).ByAccountID(args.AccountID).ByName(args.Name).GetDepartment()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	if oldDepartment != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Department name have already existed")
	}

	var department department.Department
	if err := scheme.Convert(args, &department); err != nil {
		return nil, err
	}
	department.ID = cm.NewID()
	if err := d.departmentStore(ctx).CreateDepartment(&department); err != nil {
		return nil, err
	}
	return d.departmentStore(ctx).ByID(department.ID).GetDepartment()
}

func (d *DepartmentAggregate) DeleteDepartment(ctx context.Context, args *department.DeleteDepartmentArgs) error {
	query := &identity.ListAccountUsersQuery{
		AccountID:    args.AccountID,
		DepartmentID: args.ID,
	}
	if err := d.identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result.AccountUsers) > 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng chuyển các nhân viên sang phòng ban khác trước khi xóa phòng ban")
	}
	_, err := d.departmentStore(ctx).ByID(args.ID).ByAccountID(args.AccountID).SoftDelete()
	return err
}

func (d *DepartmentAggregate) UpdateDepartment(ctx context.Context, args *department.UpdateDepartmentArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing id")
	}
	if args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing account ID")
	}

	if args.Name != "" {
		oldDepartment, err := d.departmentStore(ctx).ByAccountID(args.AccountID).ByName(args.Name).GetDepartment()
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}
		if oldDepartment != nil {
			return cm.Errorf(cm.FailedPrecondition, nil, "Department name have already existed")
		}
	}

	update := &department.Department{
		Name:        args.Name,
		Description: args.Description,
	}
	return d.departmentStore(ctx).ByID(args.ID).ByAccountID(args.AccountID).UpdateDepartment(update)
}
