package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/scheme"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type CustomerGroupCustomerStoreFactory func(context.Context) *CustomerGroupCustomerStore

func NewCustomerGroupCustomerStore(db *cmsql.Database) CustomerGroupCustomerStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CustomerGroupCustomerStore {
		return &CustomerGroupCustomerStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CustomerGroupCustomerStore struct {
	ft ShopCustomerGroupCustomerFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging
}

func (s *CustomerGroupCustomerStore) Paging(paging meta.Paging) *CustomerGroupCustomerStore {
	s.paging = paging
	return s
}

func (s *CustomerGroupCustomerStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *CustomerGroupCustomerStore) Filters(filters meta.Filters) *CustomerGroupCustomerStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *CustomerGroupCustomerStore) ID(id int64) *CustomerGroupCustomerStore {
	s.preds = append(s.preds, s.ft.ByGroupID(id))
	return s
}

func (s *CustomerGroupCustomerStore) CustomerID(id int64) *CustomerGroupCustomerStore {
	s.preds = append(s.preds, s.ft.ByCustomerID(id))
	return s
}

func (s *CustomerGroupCustomerStore) IDs(ids ...int64) *CustomerGroupCustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "group_id", ids))
	return s
}
func (s *CustomerGroupCustomerStore) CustomerIDs(ids ...int64) *CustomerGroupCustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "customer_id", ids))
	return s
}

// AddCustomerToGroup add a customer to a group. If the customer already exists in the group, it's a no-op.
func (s *CustomerGroupCustomerStore) AddShopCustomerToGroup(groupCustomer *customering.ShopCustomerGroupCustomer) (int, error) {
	sqlstore.MustNoPreds(s.preds)
	out := &model.ShopCustomerGroupCustomer{}
	if err := scheme.Convert(groupCustomer, out); err != nil {
		return 0, err
	}
	created, err := s.query().Suffix("ON CONFLICT ON CONSTRAINT shop_customer_group_customer_constraint DO NOTHING").Insert(out)
	return int(created), err
}

func (s *CustomerGroupCustomerStore) RemoveCustomerFromGroup() (int, error) {
	query := s.query().Where(s.preds)
	_deleted, err := query.Table("shop_customer_group_customer").Delete((*model.ShopCustomerGroupCustomer)(nil))
	return int(_deleted), err
}

func (s *CustomerGroupCustomerStore) GetShopCustomerToGroupDB() (*model.ShopCustomerGroupCustomer, error) {
	query := s.query().Where(s.preds)

	var groupCustomer model.ShopCustomerGroupCustomer
	err := query.ShouldGet(&groupCustomer)
	return &groupCustomer, err
}

func (s *CustomerGroupCustomerStore) ListShopCustomerGroupsCustomerByCustomerIDDB() ([]*model.ShopCustomerGroupCustomer, error) {
	query := s.query().Where(s.preds)
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopCustomerGroupCustomer, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomerGroupCustomer)
	if err != nil {
		return nil, err
	}

	var customerGroups model.ShopCustomerGroupCustomers
	err = query.Find(&customerGroups)
	return customerGroups, err
}

func (s *CustomerGroupCustomerStore) ListShopCustomerGroupsCustomerByCustomerID() ([]*customering.ShopCustomerGroupCustomer, error) {
	customerGroups, err := s.ListShopCustomerGroupsCustomerByCustomerIDDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_customeringmodel_ShopCustomerGroupCustomers_customering_ShopCustomerGroupCustomers(customerGroups), err
}

func (s *CustomerGroupCustomerStore) ListShopCustomerGroupsCustomerDB() ([]*model.ShopCustomerGroupCustomer, error) {
	query := s.query().Where(s.preds)
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopCustomerGroupCustomer, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomerGroupCustomer)
	if err != nil {
		return nil, err
	}

	var customerGroupsCustomer model.ShopCustomerGroupCustomers
	err = query.Find(&customerGroupsCustomer)
	return customerGroupsCustomer, err
}

func (s *CustomerGroupCustomerStore) ListShopCustomerGroupsCustomer() ([]*customering.ShopCustomerGroupCustomer, error) {
	customerGroupsCustomer, err := s.ListShopCustomerGroupsCustomerDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_customeringmodel_ShopCustomerGroupCustomers_customering_ShopCustomerGroupCustomers(customerGroupsCustomer), err
}
