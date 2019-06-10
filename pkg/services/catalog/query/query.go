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
	variant       sqlstore.VariantStoreFactory
	shopProduct   sqlstore.ShopProductStoreFactory
	productSource sqlstore.ProductSourceStoreFactory
}

func New(db cmsql.Database) *QueryService {
	return &QueryService{
		product:       sqlstore.NewProductStore(db),
		variant:       sqlstore.NewVariantStore(db),
		shopProduct:   sqlstore.NewShopProductStore(db),
		productSource: sqlstore.NewProductSourceStore(db),
	}
}

func (s *QueryService) GetProductByID(
	ctx context.Context, args *catalog.GetProductByIDQueryArgs,
) (*catalog.Product, error) {
	product, err := s.product(ctx).
		ID(args.ProductID).
		GetProduct()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetVariantByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.Variant, error) {
	variant, err := s.variant(ctx).
		ID(args.VariantID).
		GetVariant()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) GetProductWithVariantsByID(
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

func (s *QueryService) GetVariantWithProductByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.VariantWithProduct, error) {
	variant, err := s.variant(ctx).
		ID(args.VariantID).
		GetVariantWithProduct()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) GetShopProductWithVariantsByID(context.Context, *catalog.GetShopProductByIDQueryArgs) (*catalog.ShopProductWithVariants, error) {
	panic("implement me")
}

func (s *QueryService) GetShopProductByID(
	ctx context.Context, args *catalog.GetShopProductByIDQueryArgs,
) (*catalog.ShopProductExtended, error) {
	q := s.shopProduct(ctx).ID(args.ProductID)
	q = q.Where(q.FtShopProduct.ByShopID(args.ShopID).Optional())
	product, err := q.GetShopProduct()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetShopVariantByID(
	ctx context.Context, args *catalog.GetShopVariantByIDQueryArgs,
) (*catalog.ShopVariantExtended, error) {
	panic("implement me")
}

func (s *QueryService) GetShopVariantWithProductByID(context.Context, *catalog.GetShopVariantByIDQueryArgs) (*catalog.ShopVariantWithProduct, error) {
	panic("implement me")
}

func (s *QueryService) ListProducts(
	ctx context.Context, args *catalog.ListProductsQueryArgs,
) (*catalog.ProductsResonse, error) {
	ps := s.product(ctx).Filters(args.Filters)
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
	ctx context.Context, args *catalog.ListProductsQueryArgs,
) (*catalog.ProductsWithVariantsResponse, error) {
	q := s.product(ctx).Filters(args.Filters)
	products, err := q.ListProductsWithVariants(args.Paging)
	if err != nil {
		return nil, err
	}

	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ProductsWithVariantsResponse{
		Products: products,
		Count:    int32(count),
	}, nil
}

func (s *QueryService) ListVariants(
	ctx context.Context, args *catalog.ListVariantsQueryArgs,
) (*catalog.VariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListVariantsWithProduct(
	ctx context.Context, args *catalog.ListVariantsQueryArgs,
) (*catalog.VariantsWithProductResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopProducts(
	ctx context.Context, args *catalog.ListShopProductsQueryArgs,
) (*catalog.ShopProductsResponse, error) {
	q := s.shopProduct(ctx).Filters(args.Filters)
	q = q.Where(q.FtShopProduct.ByShopID(args.ShopID).Optional())
	products, err := q.ListShopProducts(args.Paging)
	if err != nil {
		return nil, err
	}

	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductsResponse{
		Products: products,
		Count:    int32(count),
	}, nil
}

func (s *QueryService) ListShopProductsWithVariants(
	ctx context.Context, args *catalog.ListShopProductsQueryArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopVariants(
	ctx context.Context, args *catalog.ListShopVariantsQueryArgs,
) (*catalog.ShopVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListProductsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.ProductsResonse, error) {
	panic("implement me")
}

func (s *QueryService) ListProductsWithVariantsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.ProductsResonse, error) {
	panic("implement me")
}

func (s *QueryService) ListVariantsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.VariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListVariantsWithProductByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.VariantsWithProductResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopProductsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.ShopProductsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopProductsWithVariantsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopVariantsByIDs(
	context.Context, *catalog.IDsArgs,
) (*catalog.ShopVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListShopVariantsWithProductByIDs(context.Context, *catalog.ListShopVariantsQueryArgs) (*catalog.ShopVariantsWithProductResponse, error) {
	panic("implement me")
}
