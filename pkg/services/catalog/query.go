package catalog

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/pkg/services/catalog/sqlstore"
)

var _ catalog.QueryService = &QueryService{}

type QueryService struct {
	s *sqlstore.ProductStoreFactory
}

func (s *QueryService) GetProductByID(
	ctx context.Context, args *catalog.GetProductByIDQueryArgs,
) (*catalog.Product, error) {
	panic("implement me")
}

func (s *QueryService) GetProducts(
	ctx context.Context, args *catalog.GetProductsQueryArgs,
) (*catalog.ProductsResonse, error) {
	panic("implement me")
}

func (s *QueryService) GetProductsWithVariants(
	ctx context.Context, args *catalog.GetProductsQueryArgs,
) (*catalog.ProductsWithVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetVariantByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.VariantWithProduct, error) {
	panic("implement me")
}

func (s *QueryService) GetVariants(
	ctx context.Context, args *catalog.GetVariantsQueryArgs,
) (*catalog.VariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetVariantsWithProduct(
	ctx context.Context, args *catalog.GetVariantsQueryArgs,
) (*catalog.VariantsWithProductResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetShopProductByID(
	ctx context.Context, args *catalog.GetShopProductByIDQueryArgs,
) (*catalog.ShopProductExtended, error) {
	panic("implement me")
}

func (s *QueryService) GetShopProducts(
	ctx context.Context, args *catalog.GetShopProductsQueryArgs,
) (*catalog.ShopProductsResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetShopProductsWithVariants(
	ctx context.Context, args *catalog.GetShopProductsQueryArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetShopVariantByID(
	ctx context.Context, args *catalog.GetShopVariantsByIDQueryArgs,
) (*catalog.ShopVariantExtended, error) {
	panic("implement me")
}

func (s *QueryService) GetShopVariants(
	ctx context.Context, args *catalog.GetShopVariantsQueryArgs,
) (*catalog.ShopVariantsResponse, error) {
	panic("implement me")
}
