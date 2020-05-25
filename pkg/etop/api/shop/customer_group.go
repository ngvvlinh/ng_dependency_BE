package shop

import (
	"context"

	"o.o/api/shopping/customering"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type CustomerGroupService struct {
	CustomerAggr  customering.CommandBus
	CustomerQuery customering.QueryBus
}

func (s *CustomerGroupService) Clone() *CustomerGroupService { res := *s; return &res }

func (s *CustomerGroupService) CreateCustomerGroup(ctx context.Context, r *CreateCustomerGroupEndpoint) error {
	cmd := &customering.CreateCustomerGroupCommand{
		ShopID: r.Context.Shop.ID,
		Name:   r.Name,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustopmerGroup(cmd.Result)
	return nil
}

func (s *CustomerGroupService) GetCustomerGroup(ctx context.Context, q *GetCustomerGroupEndpoint) error {
	query := &customering.GetCustomerGroupQuery{
		ID: q.Id,
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbCustopmerGroup(query.Result)
	return nil
}

func (s *CustomerGroupService) GetCustomerGroups(ctx context.Context, q *GetCustomerGroupsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &customering.ListCustomerGroupsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CustomerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CustomerGroupsResponse{
		Paging:         cmapi.PbPageInfo(paging),
		CustomerGroups: convertpb.PbCustomerGroups(query.Result.CustomerGroups),
	}
	return nil
}

func (s *CustomerGroupService) UpdateCustomerGroup(ctx context.Context, r *UpdateCustomerGroupEndpoint) error {
	cmd := &customering.UpdateCustomerGroupCommand{
		ID:   r.GroupId,
		Name: r.Name,
	}
	if err := s.CustomerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustopmerGroup(cmd.Result)
	return nil
}
