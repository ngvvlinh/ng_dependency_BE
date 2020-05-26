package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type CustomerGroupService struct{}

func (s *CustomerGroupService) Clone() *CustomerGroupService { res := *s; return &res }

func (s *CustomerGroupService) GetGroup(ctx context.Context, r *GetGroupEndpoint) error {
	resp, err := shopping.GetGroup(ctx, r.Context.Shop.ID, r.IDRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) ListGroups(ctx context.Context, r *ListGroupsEndpoint) error {
	resp, err := shopping.ListGroups(ctx, r.Context.Shop.ID, r.ListCustomerGroupsRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) CreateGroup(ctx context.Context, r *CreateGroupEndpoint) error {
	resp, err := shopping.CreateGroup(ctx, r.Context.Shop.ID, 0, r.CreateCustomerGroupRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) UpdateGroup(ctx context.Context, r *UpdateGroupEndpoint) error {
	resp, err := shopping.UpdateGroup(ctx, r.Context.Shop.ID, r.UpdateCustomerGroupRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) DeleteGroup(ctx context.Context, r *DeleteGroupEndpoint) error {
	resp, err := shopping.DeleteGroup(ctx, r.Context.Shop.ID, r.IDRequest)
	r.Result = resp
	return err
}
