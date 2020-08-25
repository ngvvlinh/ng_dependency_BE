package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type CustomerGroupRelationshipService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *CustomerGroupRelationshipService) Clone() api.CustomerGroupRelationshipService {
	res := *s
	return &res
}

func (s *CustomerGroupRelationshipService) ListRelationships(ctx context.Context, r *externaltypes.ListCustomerGroupRelationshipsRequest) (*externaltypes.CustomerGroupRelationshipsResponse, error) {
	resp, err := s.Shopping.ListRelationshipsGroupCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerGroupRelationshipService) CreateRelationship(ctx context.Context, r *externaltypes.AddCustomerRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.CreateRelationshipGroupCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerGroupRelationshipService) DeleteRelationship(ctx context.Context, r *externaltypes.RemoveCustomerRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteRelationshipGroupCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}
