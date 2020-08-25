package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type ProductCollectionService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *ProductCollectionService) Clone() api.ProductCollectionService { res := *s; return &res }

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *externaltypes.GetCollectionRequest) (*externaltypes.ProductCollection, error) {
	resp, err := s.Shopping.GetCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *externaltypes.ListCollectionsRequest) (*externaltypes.ProductCollectionsResponse, error) {
	resp, err := s.Shopping.ListCollections(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *externaltypes.CreateCollectionRequest) (*externaltypes.ProductCollection, error) {
	resp, err := s.Shopping.CreateCollection(ctx, s.SS.Shop().ID, 0, r)
	return resp, err
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *externaltypes.UpdateCollectionRequest) (*externaltypes.ProductCollection, error) {
	resp, err := s.Shopping.UpdateCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *externaltypes.GetCollectionRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}
