package sqlstore

import (
	"context"

	"etop.vn/backend/com/shopping/customering/convert"

	"etop.vn/backend/com/shopping/customering/model"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type CustomerGroupStoreFactory func(context.Context) *CustomerGroupStore

func NewCustomerGroupStore(db *cmsql.Database) CustomerGroupStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CustomerGroupStore {
		return &CustomerGroupStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type CustomerGroupStore struct {
	ft ShopCustomerGroupFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *CustomerGroupStore) Paging(paging meta.Paging) *CustomerGroupStore {
	s.paging = paging
	return s
}

func (s *CustomerGroupStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *CustomerGroupStore) Filters(filters meta.Filters) *CustomerGroupStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *CustomerGroupStore) ID(id int64) *CustomerGroupStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CustomerGroupStore) IDs(ids ...int64) *CustomerGroupStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *CustomerGroupStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopCustomerGroup)(nil))
}

func (s *CustomerGroupStore) CreateShopCustomerGroup(customerGroup *customering.ShopCustomerGroup) error {
	sqlstore.MustNoPreds(s.preds)
	var customerDB model.ShopCustomerGroup
	convert.ShopCustomerGroupDB(customerGroup, &customerDB)
	_, err := s.query().Insert(&customerDB)
	return err
}

func (s *CustomerGroupStore) GetShopCustomerGroupDB() (*model.ShopCustomerGroup, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var customerGroup model.ShopCustomerGroup
	err := query.ShouldGet(&customerGroup)
	return &customerGroup, err
}

func (s *CustomerGroupStore) GetShopCustomerGroup() (*customering.ShopCustomerGroup, error) {
	customerGroupDB, err := s.GetShopCustomerGroupDB()
	if err != nil {
		return nil, err
	}
	var out customering.ShopCustomerGroup
	convert.ShopCustomerGroup(customerGroupDB, &out)
	return &out, err
}

func (s *CustomerGroupStore) ListShopCustomerGroupsDB() ([]*model.ShopCustomerGroup, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopCustomerGroup, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomerGroup)
	if err != nil {
		return nil, err
	}

	var customerGroups model.ShopCustomerGroups
	err = query.Find(&customerGroups)
	return customerGroups, err
}

func (s *CustomerGroupStore) ListShopCustomerGroups() ([]*customering.ShopCustomerGroup, error) {
	customerGroup, err := s.ListShopCustomerGroupsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopCustomerGroups(customerGroup), nil
}

func (s *CustomerGroupStore) UpdateCustomerGroup(customerGroup *model.ShopCustomerGroup) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ft.ByID(customerGroup.ID)).UpdateAll().ShouldUpdate(customerGroup)
	return err
}
