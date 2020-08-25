package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type CustomerGroupService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *CustomerGroupService) Clone() api.CustomerGroupService { res := *s; return &res }

func (s *CustomerGroupService) GetGroup(ctx context.Context, r *pbcm.IDRequest) (*externaltypes.CustomerGroup, error) {
	resp, err := s.Shopping.GetGroup(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerGroupService) ListGroups(ctx context.Context, r *externaltypes.ListCustomerGroupsRequest) (*externaltypes.CustomerGroupsResponse, error) {
	resp, err := s.Shopping.ListGroups(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerGroupService) CreateGroup(ctx context.Context, r *externaltypes.CreateCustomerGroupRequest) (*externaltypes.CustomerGroup, error) {
	resp, err := s.Shopping.CreateGroup(ctx, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, r)
	return resp, err
}

func (s *CustomerGroupService) UpdateGroup(ctx context.Context, r *externaltypes.UpdateCustomerGroupRequest) (*externaltypes.CustomerGroup, error) {
	resp, err := s.Shopping.UpdateGroup(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerGroupService) DeleteGroup(ctx context.Context, r *pbcm.IDRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteGroup(ctx, s.SS.Shop().ID, r)
	return resp, err
}
