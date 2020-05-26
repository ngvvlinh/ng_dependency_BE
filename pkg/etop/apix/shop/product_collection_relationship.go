package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type ProductCollectionRelationshipService struct{}

func (s *ProductCollectionRelationshipService) Clone() *ProductCollectionRelationshipService {
	res := *s
	return &res
}

func (s *ProductCollectionRelationshipService) ListRelationships(ctx context.Context, r *ProductCollectionListRelationshipsEndpoint) error {
	resp, err := shopping.ListRelationshipsProductCollection(ctx, r.Context.Shop.ID, r.ListProductCollectionRelationshipsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) CreateRelationship(ctx context.Context, r *ProductCollectionCreateRelationshipEndpoint) error {
	resp, err := shopping.CreateRelationshipProductCollection(ctx, r.Context.Shop.ID, r.CreateProductCollectionRelationshipRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) DeleteRelationship(ctx context.Context, r *ProductCollectionDeleteRelationshipEndpoint) error {
	resp, err := shopping.DeleteRelationshipProductCollection(ctx, r.Context.Shop.ID, r.RemoveProductCollectionRequest)
	r.Result = resp
	return err
}
