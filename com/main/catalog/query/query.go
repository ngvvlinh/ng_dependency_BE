package query

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/api/shopping"
	com "o.o/backend/com/main"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/backend/com/main/catalog/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/capi/dot"
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
	historyStore          historysqlstore.HistoryStoreFactory
}

func New(db com.MainDB) *QueryService {
	return &QueryService{
		shopProduct:           sqlstore.NewShopProductStore(db),
		shopVariant:           sqlstore.NewShopVariantStore(db),
		shopCategory:          sqlstore.NewShopCategoryStore(db),
		shopCollection:        sqlstore.NewShopCollectionStore(db),
		shopProductCollection: sqlstore.NewShopProductCollectionStore(db),
		shopBrand:             sqlstore.NewShopBrandStore(db),
		shopVariantSupplier:   sqlstore.NewVariantSupplierStore(db),
		historyStore:          historysqlstore.NewHistoryStore(db),
	}
}

func QueryServiceMessageBus(s *QueryService) catalog.QueryBus {
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

func (s *QueryService) GetShopProduct(
	ctx context.Context, args *catalog.GetShopProductArgs,
) (*catalog.ShopProduct, error) {
	q := s.shopProduct(ctx).OptionalShopID(args.ShopID)

	count := 0
	if args.ProductID.Int64() != 0 {
		q.ID(args.ProductID)
		count++
	}
	if args.Code != "" {
		q.Code(args.Code)
		count++
	}
	if args.ExternalID != "" {
		q.ExternalID(args.ExternalID)
		count++
	}
	if count == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Arguments are invalid")
	}

	product, err := q.GetShopProduct()
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
	categories, err := q.WithPaging(args.Paging).ListShopCategories()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopCategoriesResponse{
		Categories: categories,
		Paging:     q.GetPaging(),
	}, nil
}

func (s *QueryService) ListShopCategoriesByIDs(
	ctx context.Context, shopID dot.ID, IDs []dot.ID,
) (*catalog.ShopCategoriesByIDsResponse, error) {
	q := s.shopCategory(ctx).OptionalShopID(shopID).IDs(IDs...)
	categories, err := q.ListShopCategories()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopCategoriesByIDsResponse{
		Categories: categories,
		Paging:     q.GetPaging(),
	}, nil
}

func (s *QueryService) GetShopVariant(
	ctx context.Context, args *catalog.GetShopVariantQueryArgs,
) (*catalog.ShopVariant, error) {
	q := s.shopVariant(ctx).OptionalShopID(args.ShopID)
	counter := 0
	if args.VariantID != 0 {
		q = q.ID(args.VariantID)
		counter++
	}
	if args.Code != "" {
		q = q.Code(args.Code)
		counter++
	}
	if args.ExternalID != "" {
		q = q.ExternalID(args.ExternalID)
		counter++
	}
	if counter == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Arguments are invalid (Don't allow id equals 0 and code is empty as the same time)")
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
	if args.Name != "" {
		q = q.FullTextSearchName(args.Name)
	}
	products, err := q.WithPaging(args.Paging).ListShopProducts()
	if err != nil {
		return nil, err
	}
	var productsFilter []*catalog.ShopProduct
	for _, v := range products {
		if validate.VerifySearchName(v.Name, args.Name) {
			productsFilter = append(productsFilter, v)
		}
	}
	products = productsFilter
	return &catalog.ShopProductsResponse{
		Products: products,
		Paging:   q.GetPaging(),
	}, nil
}

func (s *QueryService) ListShopProductsWithVariants(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopProductsWithVariantsResponse, error) {
	q := s.shopProduct(ctx).OptionalShopID(args.ShopID).Filters(args.Filters)
	if args.Name != "" {
		q = q.FullTextSearchName(args.Name)
	}
	products, err := q.WithPaging(args.Paging).ListShopProductsWithVariants()
	if err != nil {
		return nil, err
	}
	var productsFilter []*catalog.ShopProductWithVariants
	for _, v := range products {
		if validate.VerifySearchName(v.Name, args.Name) {
			productsFilter = append(productsFilter, v)
		}
	}
	products = []*catalog.ShopProductWithVariants{}
	var productsByVariantFilter []*catalog.ShopProductWithVariants
	isSearchByVariant := false
	if len(args.Name) > 2 {
		// filter full text search variant
		qVariant, err := s.shopVariant(ctx).OptionalShopID(args.ShopID).FullTextSearchName(args.Name).ListShopVariants()
		if err != nil {
			return nil, err
		}
		var productIDs []dot.ID
		for _, v := range qVariant {
			if v.Code == args.Name.String() {
				isSearchByVariant = true
			}
			if !cm.IDsContain(productIDs, v.ProductID) {
				productIDs = append(productIDs, v.ProductID)
			}
		}
		productByVariant, err := s.shopProduct(ctx).ShopID(args.ShopID).IDs(productIDs...).ListShopProductsWithVariants()
		if err != nil {
			return nil, err
		}
		for _, p := range productByVariant {
			if !contains(productsFilter, p.ProductID) {
				productsByVariantFilter = append(productsByVariantFilter, p)
			}
		}
	}
	productsFilter, productsByVariantFilter = ProportionProduct(productsFilter, productsByVariantFilter, 7, 3, args.Paging.Limit)
	// if search have code of variant result of variant is in first
	if isSearchByVariant {
		products = productsByVariantFilter
		products = append(products, productsFilter...)
	} else {
		products = productsFilter
		products = append(products, productsByVariantFilter...)
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

func ProportionProduct(listProduct1, listProduct2 []*catalog.ShopProductWithVariants, proportion1, proportion2, limit int) ([]*catalog.ShopProductWithVariants, []*catalog.ShopProductWithVariants) {
	if len(listProduct1)+len(listProduct2) <= limit {
		return listProduct1, listProduct2
	}
	var ln1 int
	var ln2 int

	if len(listProduct1) <= limit*proportion1/(proportion1+proportion2) {
		ln1 = len(listProduct1)
		ln2 = limit - ln1
	} else if len(listProduct2) <= limit*proportion2/(proportion1+proportion2) {
		ln2 = limit - len(listProduct2)
		ln1 = limit - ln2
	} else {
		ln1 = limit * proportion1 / (proportion1 + proportion2)
		ln2 = limit - ln1
	}
	return listProduct1[:ln1], listProduct2[:ln2]
}

func contains(a []*catalog.ShopProductWithVariants, productID dot.ID) bool {
	for _, v := range a {
		if v.ProductID == productID {
			return true
		}
	}
	return false
}

func (s *QueryService) ListShopVariants(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*catalog.ShopVariantsResponse, error) {
	return nil, cm.ErrTODO
}

func (s *QueryService) ListShopProductsByIDs(
	ctx context.Context, args *catalog.ListShopProductsByIDsArgs,
) (*catalog.ShopProductsResponse, error) {
	query := s.shopProduct(ctx).OptionalShopID(args.ShopID)
	if len(args.IDs) > 0 {
		query = query.IDs(args.IDs...)
	} else {
		query = query.WithPaging(args.Paging)
	}
	if args.IncludeDeleted {
		query = query.IncludeDeleted()
	}
	products, err := query.ListShopProducts()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductsResponse{
		Products: products,
		Paging:   query.GetPaging(),
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
	ctx context.Context, args *catalog.ListShopVariantsByIDsArgs,
) (*catalog.ShopVariantsResponse, error) {
	query := s.shopVariant(ctx).OptionalShopID(args.ShopID)
	if len(args.IDs) > 0 {
		query = query.IDs(args.IDs...)
	} else {
		query = query.WithPaging(args.Paging)
	}
	if args.IncludeDeleted {
		query = query.IncludeDeleted()
	}
	variants, err := query.ListShopVariants()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopVariantsResponse{
		Variants: variants,
		Paging:   query.GetPaging(),
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
	collections, err := q.WithPaging(args.Paging).ListShopCollections()
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
		return cm.Error(cm.InvalidArgument, "Phi??n b???n c???a s???n ph???m kh??ng c??n t???n t???i. Vui l??ng ki???m tra l???i.", nil)
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
	query := s.shopBrand(ctx).ShopID(shopID).WithPaging(paging)
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

func (s *QueryService) ListShopCollectionsByIDs(ctx context.Context, args *catalog.ListShopCollectionsByIDsArg) (*catalog.ShopCollectionsResponse, error) {
	query := s.shopCollection(ctx).ShopID(args.ShopID).WithPaging(args.Paging)
	if len(args.IDs) != 0 {
		query = query.IDs(args.IDs)
	}
	if args.IncludeDeleted {
		query = query.IncludeDeleted()
	}
	collections, err := query.ListShopCollections()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopCollectionsResponse{
		Collections: collections,
		Paging:      query.GetPaging(),
	}, nil
}

func (s *QueryService) GetShopProductCollection(ctx context.Context, args *catalog.GetShopProductCollectionArgs) (*catalog.ShopProductCollection, error) {
	shopProductCollection, err := s.shopProductCollection(ctx).ProductID(args.ProductID).CollectionID(args.CollectionID).GetShopProductCollection()
	if err != nil {
		return nil, err
	}
	return &catalog.ShopProductCollection{
		ProductID:    shopProductCollection.ProductID,
		CollectionID: shopProductCollection.CollectionID,
		ShopID:       shopProductCollection.ShopID,
	}, nil
}

func (q *QueryService) ListShopProductsCollections(ctx context.Context, args *catalog.ListProductsCollections) (*catalog.ShopProductsCollectionResponse, error) {
	query := q.shopProductCollection(ctx).WithPaging(args.Paging)
	count := 0
	if len(args.CollectionIDs) != 0 {
		query = query.IDs(args.CollectionIDs)
		count++
	}
	if len(args.ProductIds) != 0 {
		query = query.ProductIDs(args.ProductIds)
		count++
	}
	if count != 1 {
		return nil, cm.Error(cm.FailedPrecondition, "Request kh??ng h???p l???", nil)
	}
	productCollections, err := query.ListShopProductCollections()
	if err != nil {
		return nil, err
	}
	var relationships []*catalog.ShopProductCollection
	for _, productCollection := range productCollections {
		relationships = append(relationships, &catalog.ShopProductCollection{
			ProductID:    productCollection.ProductID,
			ShopID:       args.ShopID,
			CollectionID: productCollection.CollectionID,
		})
	}

	if args.IncludeDeleted {
		var productCollectionHistories catalogmodel.ProductShopCollectionHistories
		if err := q.historyStore(ctx).ListProductCollectionRelationships(&productCollectionHistories, args.ShopID, args.ProductIds, args.CollectionIDs, historysqlstore.OpDelete); err != nil {
			return nil, err
		}
		for _, arg := range productCollectionHistories {
			temp := catalogmodel.ProductShopCollectionHistory(arg)
			relationships = append(relationships, &catalog.ShopProductCollection{
				ProductID:    temp.ProductID().ID().Apply(0),
				CollectionID: temp.CollectionID().ID().Apply(0),
				ShopID:       args.ShopID,
				Deleted:      true,
			})
		}
	}
	return &catalog.ShopProductsCollectionResponse{
		ProductsCollections: relationships,
		Paging:              query.GetPaging(),
	}, nil
}

func (s *QueryService) ListShopProductWithVariantByCategoriesIDs(ctx context.Context, args *catalog.ListShopProductWithVariantByCategoriesIDsRequest) (*catalog.ShopProductsWithVariantsResponse, error) {
	q := s.shopProduct(ctx).CategoryIDs(args.CategoriesIds...).OptionalShopID(args.ShopID)
	products, err := q.ListShopProductsWithVariants()
	if err != nil {
		return nil, err
	}
	paging := q.GetPaging()
	return &catalog.ShopProductsWithVariantsResponse{
		Products: products,
		Paging:   paging,
	}, nil
}

func (s *QueryService) ListShopProductWithVariantByIDsWithPaging(ctx context.Context, args *catalog.ListShopProductWithVariantByIDsWithPagingRequest) (*catalog.ListShopProductWithVariantByIDsWithPagingResponse, error) {
	q := s.shopProduct(ctx).IDs(args.IDs...).OptionalShopID(args.ShopID).WithPaging(args.Paging)
	products, err := q.ListShopProductsWithVariants()
	if err != nil {
		return nil, err
	}
	paging := q.GetPaging()
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &catalog.ListShopProductWithVariantByIDsWithPagingResponse{
		Products: products,
		Paging:   paging,
		Total:    count,
	}, nil
}

func (s *QueryService) SearchProductByName(ctx context.Context, args *catalog.SearchProductByNameArgs) (*catalog.ListShopProductWithVariantByIDsWithPagingResponse, error) {
	q := s.shopProduct(ctx).OptionalShopID(args.ShopID)
	products, err := q.SearchNameShopProduct(args.Name)
	if err != nil {
		return nil, err
	}
	paging := q.GetPaging()
	count, err := q.CountSearName(args.Name)
	if err != nil {
		return nil, err
	}
	return &catalog.ListShopProductWithVariantByIDsWithPagingResponse{
		Products: products,
		Paging:   paging,
		Total:    count,
	}, nil
}
