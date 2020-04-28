package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/backend/com/services/affiliate/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type CustomerPolicyGroupStoreFactory func(ctx context.Context) *CustomerPolicyGroupStore

func NewCustomerPolicyGroupStore(db *cmsql.Database) CustomerPolicyGroupStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CustomerPolicyGroupStore {
		return &CustomerPolicyGroupStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CustomerPolicyGroupStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft CustomerPolicyGroupFilters

	sqlstore.Paging
	filters meta.Filters
}

func (s *CustomerPolicyGroupStore) Pred(pred interface{}) *CustomerPolicyGroupStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *CustomerPolicyGroupStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.CustomerPolicyGroup)(nil))
}

func (s *CustomerPolicyGroupStore) WithPaging(paging meta.Paging) *CustomerPolicyGroupStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *CustomerPolicyGroupStore) ID(id dot.ID) *CustomerPolicyGroupStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CustomerPolicyGroupStore) Name(name string) *CustomerPolicyGroupStore {
	s.preds = append(s.preds, s.ft.ByName(name))
	return s
}

func (s *CustomerPolicyGroupStore) SupplyID(id dot.ID) *CustomerPolicyGroupStore {
	s.preds = append(s.preds, s.ft.BySupplyID(id))
	return s
}

func (s *CustomerPolicyGroupStore) GetCustomerPolicyGroupDB() (*model.CustomerPolicyGroup, error) {
	var customerPolicyGroup model.CustomerPolicyGroup
	err := s.query().Where(s.preds).ShouldGet(&customerPolicyGroup)
	return &customerPolicyGroup, err
}

func (s *CustomerPolicyGroupStore) GetCustomerPolicyGroupsDB() ([]*model.CustomerPolicyGroup, error) {
	var results model.CustomerPolicyGroups
	err := s.query().Where(s.preds).Find(&results)
	return results, err
}

func (s *CustomerPolicyGroupStore) CreateCustomerPolicyGroup(customerPolicyGroup *model.CustomerPolicyGroup) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(customerPolicyGroup)
}

func (s *CustomerPolicyGroupStore) UpdateCustomerPolicyGroup(customerPolicyGroup *model.CustomerPolicyGroup) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().Where(`id = ?`, customerPolicyGroup.ID).ShouldUpdate(customerPolicyGroup)
}
