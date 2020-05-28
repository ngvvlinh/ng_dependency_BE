package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type ProductCollectionService struct {
	Shopping *shopping.Shopping
}

func (s *ProductCollectionService) Clone() *ProductCollectionService { res := *s; return &res }

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *GetCollectionEndpoint) error {
	resp, err := s.Shopping.GetCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *ListCollectionsEndpoint) error {
	resp, err := s.Shopping.ListCollections(ctx, r.Context.Shop.ID, r.ListCollectionsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *CreateCollectionEndpoint) error {
	resp, err := s.Shopping.CreateCollection(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *UpdateCollectionEndpoint) error {
	resp, err := s.Shopping.UpdateCollection(ctx, r.Context.Shop.ID, r.UpdateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *DeleteCollectionEndpoint) error {
	resp, err := s.Shopping.DeleteCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}
