package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type CustomerGroupRelationshipService struct{}

func (s *CustomerGroupRelationshipService) Clone() *CustomerGroupRelationshipService {
	res := *s
	return &res
}

func (s *CustomerGroupRelationshipService) ListRelationships(ctx context.Context, r *CustomerGroupListRelationshipsEndpoint) error {
	resp, err := shopping.ListRelationshipsGroupCustomer(ctx, r.Context.Shop.ID, r.ListCustomerGroupRelationshipsRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupRelationshipService) CreateRelationship(ctx context.Context, r *CustomerGroupCreateRelationshipEndpoint) error {
	resp, err := shopping.CreateRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.AddCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupRelationshipService) DeleteRelationship(ctx context.Context, r *CustomerGroupDeleteRelationshipEndpoint) error {
	resp, err := shopping.DeleteRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.RemoveCustomerRequest)
	r.Result = resp
	return err
}
