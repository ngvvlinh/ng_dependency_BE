package query

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/catalog/sqlstore"
)

var _ catalog.QueryService = &QueryService{}

type QueryService struct {
	product       sqlstore.ProductStoreFactory
	variant       sqlstore.VariantStoreFactory
	shopProduct   sqlstore.ShopProductStoreFactory
	shopVariant   sqlstore.ShopVariantStoreFactory
	productSource sqlstore.ProductSourceStoreFactory
}

func New(db cmsql.Database) *QueryService {
	return &QueryService{
		product:       sqlstore.NewProductStore(db),
		variant:       sqlstore.NewVariantStore(db),
		shopProduct:   sqlstore.NewShopProductStore(db),
		shopVariant:   sqlstore.NewShopVariantStore(db),
		productSource: sqlstore.NewProductSourceStore(db),
	}
}

func (s *QueryService) MessageBus() catalog.QueryBus {
	b := bus.New()
	return catalog.NewQueryServiceHandler(s).RegisterHandlers(b)
}

func (s *QueryService) GetProductByID(
	ctx context.Context, args *catalog.GetProductByIDQueryArgs,
) (*catalog.Product, error) {
	product, err := s.product(ctx).ID(args.ProductID).GetProduct()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetVariantByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.Variant, error) {
	variant, err := s.variant(ctx).ID(args.VariantID).GetVariant()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) GetProductWithVariantsByID(
	ctx context.Context, args *catalog.GetProductByIDQueryArgs,
) (*catalog.ProductWithVariants, error) {
	product, err := s.product(ctx).ID(args.ProductID).GetProductWithVariants()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetVariantWithProductByID(
	ctx context.Context, args *catalog.GetVariantByIDQueryArgs,
) (*catalog.VariantWithProduct, error) {
	variant, err := s.variant(ctx).ID(args.VariantID).GetVariantWithProduct()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) GetShopProductWithVariantsByID(
	ctx context.Context, args *catalog.GetShopProductByIDQueryArgs,
) (*catalog.ShopProductWithVariants, error) {
	q := s.shopProduct(ctx).ID(args.ProductID).OptionalShopID(args.ShopID)
	product, err := q.GetShopProductWithVariants()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetShopProductByID(
	ctx context.Context, args *catalog.GetShopProductByIDQueryArgs,
) (*catalog.ShopProductExtended, error) {
	q := s.shopProduct(ctx).ID(args.ProductID).OptionalShopID(args.ShopID)
	product, err := q.GetShopProduct()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetShopVariantByID(
	ctx context.Context, args *catalog.GetShopVariantByIDQueryArgs,
) (*catalog.ShopVariantExtended, error) {
	q := s.shopVariant(ctx).ID(args.VariantID).OptionalShopID(args.ShopID)
	variant, err := q.GetShopVariant()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) GetShopVariantWithProductByID(
	ctx context.Context, args *catalog.GetShopVariantByIDQueryArgs,
) (*catalog.ShopVariantWithProduct, error) {
	q := s.shopVariant(ctx).ID(args.VariantID).OptionalShopID(args.ShopID)
	variant, err := q.GetShopVariantWithProduct()
	if err != nil {
		return nil, err
	}
	return variant, nil
}

func (s *QueryService) ListProducts(
	ctx context.Context, args *catalog.ListProductsQueryArgs,
) (*catalog.ProductsResonse, error) {
	q := s.product(ctx).ProductSourceID(args.ProductSourceID).Filters(args.Filters)
	products, err := q.Paging(args.Paging).ListProducts()
	if err != nil {
		return nil, err
	}

	count, err := q.Count()
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
	q := s.product(ctx).ProductSourceID(args.ProductSourceID).Filters(args.Filters)
	products, err := q.Paging(args.Paging).ListProductsWithVariants()
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
	q := s.shopProduct(ctx).ShopID(args.ShopID).Filters(args.Filters)
	products, err := q.Paging(args.Paging).ListShopProducts()
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
	q := s.shopProduct(ctx).ShopID(args.ShopID).Filters(args.Filters)
	products, err := q.Paging(args.Paging).ListShopProductsWithVariants()
	if err != nil {
		return nil, err
	}

	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductsWithVariantsResponse{
		Products: products,
		Count:    int32(count),
	}, nil
}

func (s *QueryService) ListShopVariants(
	ctx context.Context, args *catalog.ListShopVariantsQueryArgs,
) (*catalog.ShopVariantsResponse, error) {
	panic("implement me")
}

func (s *QueryService) ListProductsByIDs(
	ctx context.Context, args *catalog.IDsArgs,
) (*catalog.ProductsResonse, error) {
	q := s.product(ctx).IDs(args.IDs...)
	products, err := q.ListProducts()
	if err != nil {
		return nil, err
	}
	return &catalog.ProductsResonse{
		Products: products,
		Count:    int32(len(products)),
	}, nil
}

func (s *QueryService) ListProductsWithVariantsByIDs(
	ctx context.Context, args *catalog.IDsArgs,
) (*catalog.ProductsWithVariantsResponse, error) {
	q := s.product(ctx).IDs(args.IDs...)
	products, err := q.ListProductsWithVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.ProductsWithVariantsResponse{
		Products: products,
		Count:    int32(len(products)),
	}, nil
}

func (s *QueryService) ListVariantsByIDs(
	ctx context.Context, args *catalog.IDsArgs,
) (*catalog.VariantsResponse, error) {
	q := s.variant(ctx).IDs(args.IDs...)
	variants, err := q.ListVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.VariantsResponse{
		Variants: variants,
		Count:    int32(len(variants)),
	}, nil
}

func (s *QueryService) ListVariantsWithProductByIDs(
	ctx context.Context, args *catalog.IDsArgs,
) (*catalog.VariantsWithProductResponse, error) {
	q := s.variant(ctx).IDs(args.IDs...)
	variants, err := q.ListVariantsWithProduct()
	if err != nil {
		return nil, err
	}
	return &catalog.VariantsWithProductResponse{
		Variants: variants,
		Count:    int32(len(variants)),
	}, nil
}

func (s *QueryService) ListShopProductsByIDs(
	ctx context.Context, args *catalog.IDsShopArgs,
) (*catalog.ShopProductsResponse, error) {
	q := s.shopProduct(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID)
	products, err := q.ListShopProducts()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductsResponse{
		Products: products,
		Count:    int32(len(products)),
	}, nil
}

func (s *QueryService) ListShopProductsWithVariantsByIDs(
	ctx context.Context, args *catalog.IDsShopArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	q := s.shopProduct(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID)
	products, err := q.ListShopProductsWithVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductsWithVariantsResponse{
		Products: products,
		Count:    int32(len(products)),
	}, nil
}

func (s *QueryService) ListShopVariantsByIDs(
	ctx context.Context, args *catalog.IDsShopArgs,
) (*catalog.ShopVariantsResponse, error) {
	q := s.shopVariant(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID)
	variants, err := q.ListShopVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopVariantsResponse{
		Variants: variants,
		Count:    int32(len(variants)),
	}, nil
}

func (s *QueryService) ListShopVariantsWithProductByIDs(
	ctx context.Context, args *catalog.IDsShopArgs,
) (*catalog.ShopVariantsWithProductResponse, error) {
	q := s.shopVariant(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID)
	variants, err := q.ListShopVariantsWithProduct()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopVariantsWithProductResponse{
		Variants: variants,
		Count:    int32(len(variants)),
	}, nil
}
