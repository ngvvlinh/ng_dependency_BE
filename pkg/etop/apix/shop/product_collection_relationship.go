package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type ProductCollectionRelationshipService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *ProductCollectionRelationshipService) Clone() api.ProductCollectionRelationshipService {
	res := *s
	return &res
}

func (s *ProductCollectionRelationshipService) ListRelationships(ctx context.Context, r *externaltypes.ListProductCollectionRelationshipsRequest) (*externaltypes.ProductCollectionRelationshipsResponse, error) {
	resp, err := s.Shopping.ListRelationshipsProductCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductCollectionRelationshipService) CreateRelationship(ctx context.Context, r *externaltypes.CreateProductCollectionRelationshipRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.CreateRelationshipProductCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductCollectionRelationshipService) DeleteRelationship(ctx context.Context, r *externaltypes.RemoveProductCollectionRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteRelationshipProductCollection(ctx, s.SS.Shop().ID, r)
	return resp, err
}
