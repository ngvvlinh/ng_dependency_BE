package department

import (
	"context"

	"o.o/api/main/department"
	api "o.o/api/top/int/shop"
	shoptypes "o.o/api/top/int/shop/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type DepartmentService struct {
	session.Session
	DepartmentGroupAggr  department.CommandBus
	DepartmentGroupQuery department.QueryBus
}

func (s *DepartmentService) Clone() api.DepartmentService { res := *s; return &res }

func (s *DepartmentService) CreateDepartment(ctx context.Context, r *api.CreateDepartmentRequest) (*shoptypes.Department, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	cmd := &department.CreateDepartmentCommand{
		AccountID:   s.SS.Shop().ID,
		Name:        r.Name,
		Description: r.Description,
	}
	if err := s.DepartmentGroupAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.Convert_core_Department_To_api_Department(cmd.Result), nil
}

func (s *DepartmentService) GetDepartments(ctx context.Context, r *api.GetDepartmentsRequest) (*api.GetDepartmentsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	query := &department.ListDepartmentsQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
	}

	if r.Filter != nil {
		query.Name = r.Filter.Name
	}

	if err = s.DepartmentGroupQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &api.GetDepartmentsResponse{
		Departments: convertpb.Convert_core_Departments_To_api_Departments(query.Result.Departments),
		Paging:      cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *DepartmentService) DeleteDepartment(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &department.DeleteDepartmentCommand{
		ID:        r.Id,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.DepartmentGroupAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: 1}, nil
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, r *api.UpdateDepartmentRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &department.UpdateDepartmentCommand{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		AccountID:   s.SS.Shop().ID,
	}
	if err := s.DepartmentGroupAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &pbcm.UpdatedResponse{
		Updated: 1,
	}, nil
}

func (s *DepartmentService) GetDepartment(ctx context.Context, r *pbcm.IDRequest) (*shoptypes.Department, error) {
	query := &department.GetDepartmentByIDQuery{
		ID:        r.Id,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.DepartmentGroupQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpb.Convert_core_Department_To_api_Department(query.Result)
	return res, nil
}
