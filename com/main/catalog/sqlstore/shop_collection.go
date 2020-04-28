package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
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
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopCollectionStore) WithPaging(paging meta.Paging) *ShopCollectionStore {
	s.Paging.WithPaging(paging)
	return s
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

func (s *ShopCollectionStore) ExternalID(externalID string) *ShopCollectionStore {
	s.preds = append(s.preds, s.ftShopCollection.ByExternalID(externalID))
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
	_deleted, err := query.Table("shop_collection").UpdateMap(map[string]interface{}{
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
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopCollection, s.ftShopCollection.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var collections model.ShopCollections
	err = query.Find(&collections)
	s.Paging.Apply(collections)
	return collections, err
}

func (s *ShopCollectionStore) ListShopCollections() ([]*catalog.ShopCollection, error) {
	collectionsModel, err := s.ListShopCollectionsDB()
	if err != nil {
		return nil, err
	}
	collections := convert.Convert_catalogmodel_ShopCollections_catalog_ShopCollections(collectionsModel)
	for i := 0; i < len(collections); i++ {
		collections[i].Deleted = !collectionsModel[i].DeletedAt.IsZero()
	}
	return collections, nil
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

func (s *ShopCollectionStore) UpdateShopCollection(collection *catalog.ShopCollection) error {
	sqlstore.MustNoPreds(s.preds)
	collectionModel := &model.ShopCollection{}
	if err := scheme.Convert(collection, collectionModel); err != nil {
		return err
	}
	err := s.query().Where(s.ftShopCollection.ByID(collectionModel.ID)).UpdateAll().ShouldUpdate(collectionModel)
	return err
}

func (s *ShopCollectionStore) DeleteShopCollection() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCollection)(nil))
	return n, err
}

func (s *ShopCollectionStore) IncludeDeleted() *ShopCollectionStore {
	s.includeDeleted = true
	return s
}
