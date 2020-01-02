package query

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ catalog.QueryService = &QueryService{}

type QueryService struct {
	shopProduct           sqlstore.ShopProductStoreFactory
	shopVariant           sqlstore.ShopVariantStoreFactory
	shopCategory          sqlstore.ShopCategoryStoreFactory
	shopCollection        sqlstore.ShopCollectionStoreFactory
	shopProductCollection sqlstore.ShopProductCollectionStoreFactory
	shopBrand             sqlstore.ShopBrandStoreFactory
	shopVariantSupplier   sqlstore.ShopVariantSupplierStoreFactory
}

func New(db *cmsql.Database) *QueryService {
	return &QueryService{
		shopProduct:           sqlstore.NewShopProductStore(db),
		shopVariant:           sqlstore.NewShopVariantStore(db),
		shopCategory:          sqlstore.NewShopCategoryStore(db),
		shopCollection:        sqlstore.NewShopCollectionStore(db),
		shopProductCollection: sqlstore.NewShopProductCollectionStore(db),
		shopBrand:             sqlstore.NewShopBrandStore(db),
		shopVariantSupplier:   sqlstore.NewVariantSupplierStore(db),
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
	q1 := s.shopProductCollection(ctx).OptionalShopID(args.ShopID).ProductID(args.ProductID)
	collections, err := q1.ListShopProductCollectionsByProductID()
	if err != nil {
		return nil, err
	}
	for _, collection := range collections {
		product.CollectionIDs = append(product.CollectionIDs, collection.CollectionID)
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

func (s *QueryService) GetShopCategory(
	ctx context.Context, args *catalog.GetShopCategoryArgs,
) (*catalog.ShopCategory, error) {
	q := s.shopCategory(ctx).ID(args.ID).OptionalShopID(args.ShopID)
	category, err := q.GetShopCategory()
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *QueryService) ListShopCategories(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopCategoriesResponse, error) {
	q := s.shopCategory(ctx).OptionalShopID(args.ShopID).Filters(args.Filters)
	categories, err := q.Paging(args.Paging).ListShopCategories()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopCategoriesResponse{
		Categories: categories,
		Paging:     q.GetPaging(),
	}, nil
}

func (s *QueryService) GetShopVariant(
	ctx context.Context, args *catalog.GetShopVariantQueryArgs,
) (*catalog.ShopVariant, error) {
	q := s.shopVariant(ctx).OptionalShopID(args.ShopID)
	if args.VariantID.Valid {
		q = q.ID(args.VariantID.ID)
	} else {
		q = q.Code(args.Code.String)
	}
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
	return &catalog.ShopProductsResponse{
		Products: products,
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
	var mapProductCollection = make(map[dot.ID][]dot.ID)
	var productIDs []dot.ID
	for _, product := range products {
		productIDs = append(productIDs, product.ProductID)
	}
	productCollections, err := s.shopProductCollection(ctx).OptionalShopID(args.ShopID).ProductIDs(productIDs).ListShopProductCollections()
	if err != nil {
		return nil, err
	}
	for _, productCollection := range productCollections {
		mapProductCollection[productCollection.ProductID] = append(mapProductCollection[productCollection.ProductID], productCollection.CollectionID)
	}
	for _, product := range products {
		product.CollectionIDs = mapProductCollection[product.ProductID]
	}
	return &catalog.ShopProductsWithVariantsResponse{
		Products: products,
		Paging:   q.GetPaging(),
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
	}, nil
}

func (s *QueryService) ListShopVariantsWithProductByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*catalog.ShopVariantsWithProductResponse, error) {
	q := s.shopVariant(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID)
	variants, err := q.IncludeDeleted().ListShopVariantsWithProduct()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopVariantsWithProductResponse{
		Variants: variants,
	}, nil
}
func (s *QueryService) GetShopCollection(
	ctx context.Context, args *catalog.GetShopCollectionArgs,
) (*catalog.ShopCollection, error) {
	q := s.shopCollection(ctx).ID(args.ID).OptionalShopID(args.ShopID)
	collection, err := q.GetShopCollection()
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (s *QueryService) ListShopCollections(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopCollectionsResponse, error) {
	q := s.shopCollection(ctx).OptionalShopID(args.ShopID).Filters(args.Filters)
	collections, err := q.Paging(args.Paging).ListShopCollections()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopCollectionsResponse{
		Collections: collections,
		Paging:      q.GetPaging(),
	}, nil
}

func (s *QueryService) ValidateVariantIDs(ctx context.Context, shopId dot.ID, shopVariantIds []dot.ID) error {
	dbResult, err := s.shopVariant(ctx).IDs(shopVariantIds...).ShopID(shopId).ListShopVariantsDB()
	if err != nil {
		return err
	}
	if len(dbResult) != len(shopVariantIds) {
		return cm.Error(cm.InvalidArgument, "Phiên bản của sản phẩm không còn tồn tại. Vui lòng kiểm tra lại.", nil)
	}
	return nil
}

func (s *QueryService) ListShopCollectionsByProductID(
	ctx context.Context, args *catalog.ListShopCollectionsByProductIDArgs,
) ([]*catalog.ShopCollection, error) {
	q := s.shopProductCollection(ctx).OptionalShopID(args.ShopID).ProductID(args.ProductID)
	productCollections, err := q.ListShopProductCollectionsByProductID()
	if err != nil {
		return nil, err
	}
	collectionIDs := make([]dot.ID, len(productCollections))
	for _, pc := range productCollections {
		collectionIDs = append(collectionIDs, pc.CollectionID)
	}
	qc := s.shopCollection(ctx).OptionalShopID(args.ShopID).IDs(collectionIDs) // qc=querycollection
	collections, err := qc.ListShopCollections()
	return collections, err
}

func (s *QueryService) GetBrandByID(ctx context.Context, id dot.ID, shopID dot.ID) (*catalog.ShopBrand, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := s.shopBrand(ctx).ShopID(shopID).ID(id).GetShopBrand()
	return result, err
}

func (s *QueryService) GetBrandsByIDs(ctx context.Context, ids []dot.ID, shopID dot.ID) ([]*catalog.ShopBrand, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := s.shopBrand(ctx).ShopID(shopID).IDs(ids...).ListShopBrands()
	return result, err
}

func (s *QueryService) ListBrands(ctx context.Context, paging meta.Paging, shopID dot.ID) (*catalog.ListBrandsResult, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	query := s.shopBrand(ctx).ShopID(shopID).Paging(paging)
	result, err := query.ListShopBrands()
	if err != nil {
		return nil, err
	}
	listBrandResult := &catalog.ListBrandsResult{
		ShopBrands: result,
		PageInfo:   query.GetPaging(),
	}
	return listBrandResult, err
}

func (s *QueryService) GetSupplierIDsByVariantID(ctx context.Context, variantID dot.ID, shopID dot.ID) ([]dot.ID, error) {
	if shopID == 0 || variantID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id or supplier_id in request", nil)
	}
	variantSupplier, err := s.shopVariantSupplier(ctx).ShopID(shopID).VariantID(variantID).ListVariantSupplier()
	if err != nil {
		return nil, err
	}
	var listSupplierIDs = make([]dot.ID, len(variantSupplier))
	for _, value := range variantSupplier {
		listSupplierIDs = append(listSupplierIDs, value.SupplierID)
	}
	return listSupplierIDs, nil
}

func (s *QueryService) GetVariantsBySupplierID(ctx context.Context, supplierID dot.ID, shopID dot.ID) (*catalog.ShopVariantsResponse, error) {
	if shopID == 0 || supplierID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id or supplier_id in request", nil)
	}
	variantSupplier, err := s.shopVariantSupplier(ctx).ShopID(shopID).SupplierID(supplierID).ListVariantSupplier()
	if err != nil {
		return nil, err
	}
	var listVariants = make([]dot.ID, len(variantSupplier))
	for _, value := range variantSupplier {
		listVariants = append(listVariants, value.VariantID)
	}
	variants, err := s.shopVariant(ctx).ShopID(shopID).IDs(listVariants...).ListShopVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopVariantsResponse{
		Variants: variants,
	}, err
}
