package query

import (
	"context"
	"o.o/api/main/department"
	"o.o/api/main/identity"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/department/sqlstore"
	"o.o/backend/pkg/common/bus"
)

var _ department.QueryService = &DepartmentQuery{}

type DepartmentQuery struct {
	departmentStore sqlstore.DepartmentStoreFactory
	identityQuery   identity.QueryBus
}

func NewDepartmentQuery(db com.MainDB, identityQ identity.QueryBus) *DepartmentQuery {
	return &DepartmentQuery{
		departmentStore: sqlstore.NewDepartmentStore(db),
		identityQuery:   identityQ,
	}
}
func DepartmentQueryMessageBus(q *DepartmentQuery) department.QueryBus {
	b := bus.New()
	return department.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (d *DepartmentQuery) ListDepartments(ctx context.Context, args *department.ListDepartmentsArgs) (*department.ListDepartmentsResponse, error) {
	query := d.departmentStore(ctx)

	if args.Name != "" {
		query = query.ByNameNorm(args.Name)
	}

	if args.AccountID != 0 {
		query = query.ByAccountID(args.AccountID)
	}

	departments, err := query.WithPaging(args.Paging).ListDepartments()
	if err != nil {
		return nil, err
	}

	cmdQuery := &identity.ListAccountUsersByDepartmentIDsQuery{
		AccountID: args.AccountID,
	}

	if err := d.identityQuery.Dispatch(ctx, cmdQuery); err != nil {
		return nil, err
	}

	results := cmdQuery.Result
	for _, department := range departments {
		for _, result := range results {
			if department.ID == result.DepartmentID {
				department.Count = result.Count
			}
		}
	}

	return &department.ListDepartmentsResponse{
		Departments: departments,
		Paging:      query.GetPaging(),
	}, nil
}

func (d *DepartmentQuery) GetDepartmentByID(ctx context.Context, args *department.GetDepartmentByIDArgs) (*department.Department, error) {
	return d.departmentStore(ctx).ByID(args.ID).ByAccountID(args.AccountID).GetDepartment()
}
