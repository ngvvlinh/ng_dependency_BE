package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type ProductCollectionService struct{}

func (s *ProductCollectionService) Clone() *ProductCollectionService { res := *s; return &res }

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *GetCollectionEndpoint) error {
	resp, err := shopping.GetCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *ListCollectionsEndpoint) error {
	resp, err := shopping.ListCollections(ctx, r.Context.Shop.ID, r.ListCollectionsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *CreateCollectionEndpoint) error {
	resp, err := shopping.CreateCollection(ctx, r.Context.Shop.ID, 0, r.CreateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *UpdateCollectionEndpoint) error {
	resp, err := shopping.UpdateCollection(ctx, r.Context.Shop.ID, r.UpdateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *DeleteCollectionEndpoint) error {
	resp, err := shopping.DeleteCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}
