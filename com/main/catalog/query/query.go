package query

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ catalog.QueryService = &QueryService{}

type QueryService struct {
	shopProduct sqlstore.ShopProductStoreFactory
	shopVariant sqlstore.ShopVariantStoreFactory
}

func New(db cmsql.Database) *QueryService {
	return &QueryService{
		shopProduct: sqlstore.NewShopProductStore(db),
		shopVariant: sqlstore.NewShopVariantStore(db),
	}
}

func (s *QueryService) MessageBus() catalog.QueryBus {
	b := bus.New()
	return catalog.NewQueryServiceHandler(s).RegisterHandlers(b)
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
) (*catalog.ShopProduct, error) {
	q := s.shopProduct(ctx).ID(args.ProductID).OptionalShopID(args.ShopID)
	product, err := q.GetShopProduct()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *QueryService) GetShopVariantByID(
	ctx context.Context, args *catalog.GetShopVariantByIDQueryArgs,
) (*catalog.ShopVariant, error) {
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

func (s *QueryService) ListShopProducts(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopProductsResponse, error) {
	q := s.shopProduct(ctx).OptionalShopID(args.ShopID).Filters(args.Filters)
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
		Paging:   q.GetPaging(),
	}, nil
}

func (s *QueryService) ListShopProductsWithVariants(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	q := s.shopProduct(ctx).OptionalShopID(args.ShopID).Filters(args.Filters)
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
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopVariantsResponse, error) {
	return nil, cm.ErrTODO
}

func (s *QueryService) ListShopProductsByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
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
	ctx context.Context, args *shopping.IDsQueryShopArgs,
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
	ctx context.Context, args *shopping.IDsQueryShopArgs,
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
	ctx context.Context, args *shopping.IDsQueryShopArgs,
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
