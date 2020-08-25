package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type ProductService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *ProductService) Clone() api.ProductService { res := *s; return &res }

func (s *ProductService) GetProduct(ctx context.Context, r *externaltypes.GetProductRequest) (*externaltypes.ShopProduct, error) {
	resp, err := s.Shopping.GetProduct(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductService) ListProducts(ctx context.Context, r *externaltypes.ListProductsRequest) (*externaltypes.ShopProductsResponse, error) {
	resp, err := s.Shopping.ListProducts(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductService) CreateProduct(ctx context.Context, r *externaltypes.CreateProductRequest) (*externaltypes.ShopProduct, error) {
	resp, err := s.Shopping.CreateProduct(ctx, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, r)
	return resp, err
}

func (s *ProductService) UpdateProduct(ctx context.Context, r *externaltypes.UpdateProductRequest) (*externaltypes.ShopProduct, error) {
	resp, err := s.Shopping.UpdateProduct(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *ProductService) DeleteProduct(ctx context.Context, r *externaltypes.GetProductRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteProduct(ctx, s.SS.Shop().ID, r)
	return resp, err
}
