package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/com/main/catalog/convert"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type ShopCollectionStoreFactory func(context.Context) *ShopCollectionStore

func NewShopCollectionStore(db *cmsql.Database) ShopCollectionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopCollectionStore {
		return &ShopCollectionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopCollectionStore struct {
	ftShopCollection ShopCollectionFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopCollectionStore) Paging(paging meta.Paging) *ShopCollectionStore {
	s.paging = paging
	return s
}

func (s *ShopCollectionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopCollectionStore) Where(pred sq.FilterQuery) *ShopCollectionStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopCollectionStore) Filters(filters meta.Filters) *ShopCollectionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShopCollectionStore) ID(id dot.ID) *ShopCollectionStore {
	s.preds = append(s.preds, s.ftShopCollection.ByID(id))
	return s
}

func (s *ShopCollectionStore) IDs(ids []dot.ID) *ShopCollectionStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShopCollectionStore) OptionalShopID(id dot.ID) *ShopCollectionStore {
	s.preds = append(s.preds, s.ftShopCollection.ByShopID(id).Optional())
	return s
}

func (s *ShopCollectionStore) ShopID(id dot.ID) *ShopCollectionStore {
	s.preds = append(s.preds, s.ftShopCollection.ByShopID(id))
	return s
}

func (s *ShopCollectionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCollection.NotDeleted())
	_deleted, err := query.Table("shop_Collection").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ShopCollectionStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCollection.NotDeleted())
	return query.Count((*model.ShopCollection)(nil))
}

func (s *ShopCollectionStore) GetShopCollectionDB() (*model.ShopCollection, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCollection.NotDeleted())

	var collection model.ShopCollection
	err := query.ShouldGet(&collection)
	return &collection, err
}

func (s *ShopCollectionStore) GetShopCollection() (*catalog.ShopCollection, error) {
	collection, err := s.GetShopCollectionDB()
	if err != nil {
		return nil, err
	}
	var out catalog.ShopCollection
	err = scheme.Convert(collection, &out)
	return &out, err
}

func (s *ShopCollectionStore) ListShopCollectionsDB() ([]*model.ShopCollection, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCollection.NotDeleted())
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopCollection, s.ftShopCollection.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var collections model.ShopCollections
	err = query.Find(&collections)
	return collections, err
}

func (s *ShopCollectionStore) ListShopCollections() ([]*catalog.ShopCollection, error) {
	collectionsModel, err := s.ListShopCollectionsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopCollections_catalog_ShopCollections(collectionsModel), nil
}

func (s *ShopCollectionStore) CreateShopCollection(Collection *catalog.ShopCollection) error {
	sqlstore.MustNoPreds(s.preds)
	var collectionDB model.ShopCollection
	if err := scheme.Convert(Collection, &collectionDB); err != nil {
		return err
	}
	_, err := s.query().Insert(&collectionDB)
	return err
}

func (s *ShopCollectionStore) UpdateShopCollection(Collection *model.ShopCollection) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ftShopCollection.ByID(Collection.ID)).UpdateAll().ShouldUpdate(Collection)
	return err
}

func (s *ShopCollectionStore) DeleteShopCollection() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCollection)(nil))
	return n, err
}
