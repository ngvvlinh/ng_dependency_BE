package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type ProductCollectionRelationshipService struct {
	Shopping *shopping.Shopping
}

func (s *ProductCollectionRelationshipService) Clone() *ProductCollectionRelationshipService {
	res := *s
	return &res
}

func (s *ProductCollectionRelationshipService) ListRelationships(ctx context.Context, r *ProductCollectionListRelationshipsEndpoint) error {
	resp, err := s.Shopping.ListRelationshipsProductCollection(ctx, r.Context.Shop.ID, r.ListProductCollectionRelationshipsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) CreateRelationship(ctx context.Context, r *ProductCollectionCreateRelationshipEndpoint) error {
	resp, err := s.Shopping.CreateRelationshipProductCollection(ctx, r.Context.Shop.ID, r.CreateProductCollectionRelationshipRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) DeleteRelationship(ctx context.Context, r *ProductCollectionDeleteRelationshipEndpoint) error {
	resp, err := s.Shopping.DeleteRelationshipProductCollection(ctx, r.Context.Shop.ID, r.RemoveProductCollectionRequest)
	r.Result = resp
	return err
}
