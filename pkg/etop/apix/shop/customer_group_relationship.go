package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type CustomerGroupRelationshipService struct {
	Shopping *shopping.Shopping
}

func (s *CustomerGroupRelationshipService) Clone() *CustomerGroupRelationshipService {
	res := *s
	return &res
}

func (s *CustomerGroupRelationshipService) ListRelationships(ctx context.Context, r *CustomerGroupListRelationshipsEndpoint) error {
	resp, err := s.Shopping.ListRelationshipsGroupCustomer(ctx, r.Context.Shop.ID, r.ListCustomerGroupRelationshipsRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupRelationshipService) CreateRelationship(ctx context.Context, r *CustomerGroupCreateRelationshipEndpoint) error {
	resp, err := s.Shopping.CreateRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.AddCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupRelationshipService) DeleteRelationship(ctx context.Context, r *CustomerGroupDeleteRelationshipEndpoint) error {
	resp, err := s.Shopping.DeleteRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.RemoveCustomerRequest)
	r.Result = resp
	return err
}
