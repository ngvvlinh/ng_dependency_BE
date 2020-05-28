package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type ProductService struct {
	Shopping *shopping.Shopping
}

func (s *ProductService) Clone() *ProductService { res := *s; return &res }

func (s *ProductService) GetProduct(ctx context.Context, r *GetProductEndpoint) error {
	resp, err := s.Shopping.GetProduct(ctx, r.Context.Shop.ID, r.GetProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) ListProducts(ctx context.Context, r *ListProductsEndpoint) error {
	resp, err := s.Shopping.ListProducts(ctx, r.Context.Shop.ID, r.ListProductsRequest)
	r.Result = resp
	return err
}

func (s *ProductService) CreateProduct(ctx context.Context, r *CreateProductEndpoint) error {
	resp, err := s.Shopping.CreateProduct(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) UpdateProduct(ctx context.Context, r *UpdateProductEndpoint) error {
	resp, err := s.Shopping.UpdateProduct(ctx, r.Context.Shop.ID, r.UpdateProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) DeleteProduct(ctx context.Context, r *DeleteProductEndpoint) error {
	resp, err := s.Shopping.DeleteProduct(ctx, r.Context.Shop.ID, r.GetProductRequest)
	r.Result = resp
	return err
}
