package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type ShopCategoryStoreFactory func(context.Context) *ShopCategoryStore

func NewShopCategoryStore(db cmsql.Database) ShopCategoryStoreFactory {
	return func(ctx context.Context) *ShopCategoryStore {
		return &ShopCategoryStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopCategoryStore struct {
	ftShopCategory ShopCategoryFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopCategoryStore) Paging(paging meta.Paging) *ShopCategoryStore {
	s.paging = paging
	return s
}

func (s *ShopCategoryStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopCategoryStore) Where(pred sq.FilterQuery) *ShopCategoryStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopCategoryStore) Filters(filters meta.Filters) *ShopCategoryStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShopCategoryStore) ID(id int64) *ShopCategoryStore {
	s.preds = append(s.preds, s.ftShopCategory.ByID(id))
	return s
}
func (s *ShopCategoryStore) OptionalShopID(id int64) *ShopCategoryStore {
	s.preds = append(s.preds, s.ftShopCategory.ByShopID(id).Optional())
	return s
}

func (s *ShopCategoryStore) ShopID(id int64) *ShopCategoryStore {
	s.preds = append(s.preds, s.ftShopCategory.ByShopID(id))
	return s
}

func (s *ShopCategoryStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCategory.NotDeleted())
	_deleted, err := query.Table("shop_category").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *ShopCategoryStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCategory.NotDeleted())
	return query.Count((*model.ShopCategory)(nil))
}

func (s *ShopCategoryStore) GetShopCategoryDB() (*model.ShopCategory, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCategory.NotDeleted())

	var category model.ShopCategory
	err := query.ShouldGet(&category)
	return &category, err
}

func (s *ShopCategoryStore) GetShopCategory() (*catalog.ShopCategory, error) {
	category, err := s.GetShopCategoryDB()
	if err != nil {
		return nil, err
	}
	var out catalog.ShopCategory
	convert.ShopCategory(category, &out)
	return &out, err
}

func (s *ShopCategoryStore) ListShopCategoriesDB() ([]*model.ShopCategory, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopCategory.NotDeleted())
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopCategory, s.ftShopCategory.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var categories model.ShopCategories
	err = query.Find(&categories)
	return categories, err
}

func (s *ShopCategoryStore) ListShopCategories() ([]*catalog.ShopCategory, error) {
	categories, err := s.ListShopCategoriesDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopCategories(categories), nil
}

func (s *ShopCategoryStore) CreateShopCategory(category *catalog.ShopCategory) error {
	sqlstore.MustNoPreds(s.preds)
	var categoryDB model.ShopCategory
	convert.ShopCategoryDB(category, &categoryDB)
	_, err := s.query().Insert(&categoryDB)
	return err
}

func (s *ShopCategoryStore) UpdateShopCategory(category *model.ShopCategory) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ftShopCategory.ByID(category.ID)).UpdateAll().ShouldUpdate(category)
	return err
}

func (s *ShopCategoryStore) DeleteShopCategory() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCategory)(nil))
	return int(n), err
}
