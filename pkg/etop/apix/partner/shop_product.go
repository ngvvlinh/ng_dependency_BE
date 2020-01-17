package partner

import (
	"context"

	"etop.vn/backend/pkg/etop/apix/shopping"
)

func (s *ProductService) GetProduct(ctx context.Context, r *GetProductEndpoint) error {
	resp, err := shopping.GetProduct(ctx, r.Context.Shop.ID, r.GetProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) ListProducts(ctx context.Context, r *ListProductsEndpoint) error {
	resp, err := shopping.ListProducts(ctx, r.Context.Shop.ID, r.ListProductsRequest)
	r.Result = resp
	return err
}

func (s *ProductService) CreateProduct(ctx context.Context, r *CreateProductEndpoint) error {
	resp, err := shopping.CreateProduct(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) UpdateProduct(ctx context.Context, r *UpdateProductEndpoint) error {
	resp, err := shopping.UpdateProduct(ctx, r.Context.Shop.ID, r.UpdateProductRequest)
	r.Result = resp
	return err
}

func (s *ProductService) DeleteProduct(ctx context.Context, r *DeleteProductEndpoint) error {
	resp, err := shopping.DeleteProduct(ctx, r.Context.Shop.ID, r.GetProductRequest)
	r.Result = resp
	return err
}
