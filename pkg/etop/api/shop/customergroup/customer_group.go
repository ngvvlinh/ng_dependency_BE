package customergroup

import (
	"context"

	"o.o/api/shopping/customering"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type CustomerGroupService struct {
	session.Session

	CustomerAggr  customering.CommandBus
	CustomerQuery customering.QueryBus
}

func (s *CustomerGroupService) Clone() api.CustomerGroupService { res := *s; return &res }

func (s *CustomerGroupService) CreateCustomerGroup(ctx context.Context, r *api.CreateCustomerGroupRequest) (*api.CustomerGroup, error) {
	cmd := &customering.CreateCustomerGroupCommand{
		ShopID: s.SS.Shop().ID,
		Name:   r.Name,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCustopmerGroup(cmd.Result)
	return result, nil
}

func (s *CustomerGroupService) GetCustomerGroup(ctx context.Context, q *pbcm.IDRequest) (*api.CustomerGroup, error) {
	query := &customering.GetCustomerGroupQuery{
		ID: q.Id,
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbCustopmerGroup(query.Result)
	return result, nil
}

func (s *CustomerGroupService) GetCustomerGroups(ctx context.Context, q *api.GetCustomerGroupsRequest) (*api.CustomerGroupsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &customering.ListCustomerGroupsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CustomerGroupsResponse{
		Paging:         cmapi.PbPageInfo(paging),
		CustomerGroups: convertpb.PbCustomerGroups(query.Result.CustomerGroups),
	}
	return result, nil
}

func (s *CustomerGroupService) UpdateCustomerGroup(ctx context.Context, r *api.UpdateCustomerGroupRequest) (*api.CustomerGroup, error) {
	cmd := &customering.UpdateCustomerGroupCommand{
		ID:   r.GroupId,
		Name: r.Name,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCustopmerGroup(cmd.Result)
	return result, nil
}
