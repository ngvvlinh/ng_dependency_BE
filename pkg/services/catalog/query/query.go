package query

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/catalog/sqlstore"
)

var _ catalog.QueryService = &QueryService{}

type QueryService struct {
	product       sqlstore.ProductStoreFactory
	productSource sqlstore.ProductSourceStoreFactory
}

func New(db cmsql.Database) *QueryService {
	return &QueryService{
		product:       sqlstore.NewProductStore(db),
		productSource: sqlstore.NewProductSourceStore(db),
	}
}

func (s *QueryService) GetProductByID(
	ctx context.Context, args *catalog.GetProductByIDQueryArgs,
) (*catalog.ProductWithVariants, error) {
	product, err := s.product(ctx).
		ID(args.ProductID).
		GetProductWithVariants()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) ListProducts(
	ctx context.Context, args *catalog.GetProductsQueryArgs,
) (*catalog.ProductsResonse, error) {
	ps := s.product(ctx)
	if args.IDs != nil {
		ps.IDs(args.IDs...)
	}
	ps.Filters(args.Filters)
	products, err := ps.ListProducts(args.Paging)
	if err != nil {
		return nil, err
	}

	count, err := ps.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ProductsResonse{
		Products: products,
		Count:    int32(count),
	}, nil
}

func (s *QueryService) ListProductsWithVariants(
	ctx context.Context, args *catalog.GetProductsQueryArgs,
) (*catalog.ProductsWithVariantsResponse, error) {
	ps := s.product(ctx)
	if args.IDs != nil {
		ps.IDs(args.IDs...)
	}
	ps.Filters(args.Filters)
	products, err := ps.ListProductsWithVariants(args.Paging)
	if err != nil {
		return nil, err
	}

	count, err := ps.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ProductsWithVariantsResponse{
		Products: products,
		Count:    int32(count),
	}, nil
}

func (s *QueryService) GetVariantByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.VariantWithProduct, error) {
	panic("implement me")
}

func (s *QueryService) ListVariants(
	ctx context.Context, args *catalog.GetVariantsQueryArgs,
) (*catalog.VariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListVariantsWithProduct(
	ctx context.Context, args *catalog.GetVariantsQueryArgs,
) (*catalog.VariantsWithProductResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetShopProductByID(
	ctx context.Context, args *catalog.GetShopProductByIDQueryArgs,
) (*catalog.ShopProductExtended, error) {
	panic("implement me")
}

func (s *QueryService) ListShopProducts(
	ctx context.Context, args *catalog.GetShopProductsQueryArgs,
) (*catalog.ShopProductsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopProductsWithVariants(
	ctx context.Context, args *catalog.GetShopProductsQueryArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) GetShopVariantByID(
	ctx context.Context, args *catalog.GetShopVariantsByIDQueryArgs,
) (*catalog.ShopVariantExtended, error) {
	panic("implement me")
}

func (s *QueryService) ListShopVariants(
	ctx context.Context, args *catalog.GetShopVariantsQueryArgs,
) (*catalog.ShopVariantsResponse, error) {
	panic("implement me")
}
