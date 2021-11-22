package department

import (
	"context"
	"o.o/api/meta"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateDepartment(context.Context, *CreateDepartmentArgs) (*Department, error)
	DeleteDepartment(context.Context, *DeleteDepartmentArgs) error
	UpdateDepartment(context.Context, *UpdateDepartmentArgs) error
}

type QueryService interface {
	ListDepartments(context.Context, *ListDepartmentsArgs) (*ListDepartmentsResponse, error)
	GetDepartmentByID(context.Context, *GetDepartmentByIDArgs) (*Department, error)
}

// +convert:create=Department
type CreateDepartmentArgs struct {
	AccountID   dot.ID
	Name        string
	Description string
}

func (args *CreateDepartmentArgs) Validate() error {
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if args.Name == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing name")
	}
	return nil
}

type ListDepartmentsArgs struct {
	AccountID dot.ID
	Name      filter.FullTextSearch
	Paging    meta.Paging
	Filters   meta.Filters
}

type ListDepartmentsResponse struct {
	Departments []*Department

	Paging meta.PageInfo
}

type GetDepartmentByIDArgs struct {
	ID        dot.ID
	AccountID dot.ID
}

type DeleteDepartmentArgs struct {
	ID        dot.ID
	AccountID dot.ID
}

// +convert:update=Department
type UpdateDepartmentArgs struct {
	ID          dot.ID
	AccountID   dot.ID
	Name        string
	Description string
}
