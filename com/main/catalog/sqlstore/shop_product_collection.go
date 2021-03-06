package sqlstore

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShopProductCollectionStoreFactory func(context.Context) *ShopProductCollectionStore

func NewShopProductCollectionStore(db *cmsql.Database) ShopProductCollectionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopProductCollectionStore {
		return &ShopProductCollectionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopProductCollectionStore struct {
	ftShopProductCollection ShopProductCollectionFilters
	ftShopcollection        ShopCollectionFilters
	ftShopProduct           ShopProductFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *ShopProductCollectionStore) WithPaging(paging meta.Paging) *ShopProductCollectionStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *ShopProductCollectionStore) Filters(filters meta.Filters) *ShopProductCollectionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShopProductCollectionStore) ShopID(id dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByShopID(id))
	return s
}

func (s *ShopProductCollectionStore) IDs(ids []dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("collection_id", ids))
	return s
}

func (s *ShopProductCollectionStore) ProductID(id dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByProductID(id))
	return s
}

func (s *ShopProductCollectionStore) ProductIDs(ids []dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *ShopProductCollectionStore) CollectionID(id dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("collection_id", id))
	return s
}

func (s *ShopProductCollectionStore) OptionalShopID(id dot.ID) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByShopID(id).Optional())
	return s
}

func (s *ShopProductCollectionStore) RemoveProductFromCollection() (int, error) {
	query := s.query().Where(s.preds)
	_deleted, err := query.Table("shop_product_collection").Delete((*model.ShopProductCollection)(nil))
	return _deleted, err
}

// AddProductToCollection adds a product to a collection. If the product already exists in the collection, it's a no-op.
func (s *ShopProductCollectionStore) AddProductToCollection(productCollection *catalog.ShopProductCollection) (int, error) {
	sqlstore.MustNoPreds(s.preds)
	var out model.ShopProductCollection
	if err := scheme.Convert(productCollection, &out); err != nil {
		return 0, err
	}
	created, err := s.query().Suffix("ON CONFLICT ON CONSTRAINT shop_product_collection_constraint DO NOTHING").Insert(&out)
	return created, err
}

func (s *ShopProductCollectionStore) GetShopProductCollectionDB() (*model.ShopProductCollection, error) {
	query := s.query().Where(s.preds)

	var shopProductCollection model.ShopProductCollection
	err := query.ShouldGet(&shopProductCollection)
	return &shopProductCollection, err
}

func (s *ShopProductCollectionStore) GetShopProductCollection() (*catalog.ShopProductCollection, error) {
	shopProductCollectionDB, err := s.GetShopProductCollectionDB()
	if err != nil {
		return nil, err
	}
	result := &catalog.ShopProductCollection{}
	err = scheme.Convert(shopProductCollectionDB, result)
	return result, err
}

func (s *ShopProductCollectionStore) ListShopProductCollectionsByProductIDDB() ([]*model.ShopProductCollection, error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopProductCollection, s.ftShopProductCollection.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var productCollections model.ShopProductCollections
	err = query.Find(&productCollections)
	return productCollections, err
}

func (s *ShopProductCollectionStore) ListShopProductCollectionsByProductID() ([]*catalog.ShopProductCollection, error) {
	productCollections, err := s.ListShopProductCollectionsByProductIDDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopProductCollections_catalog_ShopProductCollections(productCollections), err
}

func (s *ShopProductCollectionStore) ListShopProductCollectionsDB() ([]*model.ShopProductCollection, error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopProductCollection, s.ftShopProductCollection.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var productCollections model.ShopProductCollections
	err = query.Find(&productCollections)
	return productCollections, err
}

func (s *ShopProductCollectionStore) ListShopProductCollections() ([]*catalog.ShopProductCollection, error) {
	productCollections, err := s.ListShopProductCollectionsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopProductCollections_catalog_ShopProductCollections(productCollections), err
}

func (s *ShopProductCollectionStore) DeleteProductCollections() (int, error) {
	query := s.query().Where(s.preds)
	_deleted, err := query.Table("shop_product_collection").Delete((*model.ShopProductCollection)(nil))
	return _deleted, err
}
