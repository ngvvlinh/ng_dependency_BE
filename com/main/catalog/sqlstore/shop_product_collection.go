package sqlstore

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type ShopProductCollectionStoreFactory func(context.Context) *ShopProductCollectionStore

func NewShopProductCollectionStore(db *cmsql.Database) ShopProductCollectionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopProductCollectionStore {
		return &ShopProductCollectionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type ShopProductCollectionStore struct {
	ftShopProductCollection ShopProductCollectionFilters
	ftShopcollection        ShopCollectionFilters
	ftShopProduct           ShopProductFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging
}

func (s *ShopProductCollectionStore) Paging(paging meta.Paging) *ShopProductCollectionStore {
	s.paging = paging
	return s
}

func (s *ShopProductCollectionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopProductCollectionStore) Filters(filters meta.Filters) *ShopProductCollectionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShopProductCollectionStore) ShopID(id int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByShopID(id))
	return s
}

func (s *ShopProductCollectionStore) IDs(ids []int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("collection_id", ids))
	return s
}

func (s *ShopProductCollectionStore) ProductID(id int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByProductID(id))
	return s
}

func (s *ShopProductCollectionStore) ProductIDs(ids []int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *ShopProductCollectionStore) CollectionID(ids ...int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, sq.In("collection_id", ids))
	return s
}

func (s *ShopProductCollectionStore) OptionalShopID(id int64) *ShopProductCollectionStore {
	s.preds = append(s.preds, s.ftShopProductCollection.ByShopID(id).Optional())
	return s
}

func (s *ShopProductCollectionStore) RemoveProductFromCollection() (int, error) {
	query := s.query().Where(s.preds)
	_deleted, err := query.Table("shop_product_collection").Delete((*model.ShopProductCollection)(nil))
	return int(_deleted), err
}

// AddProductToCollection adds a product to a collection. If the product already exists in the collection, it's a no-op.
func (s *ShopProductCollectionStore) AddProductToCollection(productCollection *catalog.ShopProductCollection) (int, error) {
	sqlstore.MustNoPreds(s.preds)
	var out model.ShopProductCollection
	convert.ShopProductCollectionDB(productCollection, &out)
	created, err := s.query().Suffix("ON CONFLICT ON CONSTRAINT shop_product_collection_constraint DO NOTHING").Insert(&out)
	return int(created), err
}

func (s *ShopProductCollectionStore) ListShopProductCollectionsByProductIDDB() ([]*model.ShopProductCollection, error) {
	query := s.query().Where(s.preds)
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopProductCollection, s.ftShopProductCollection.prefix)
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
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopProductCollection, s.ftShopProductCollection.prefix)
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
