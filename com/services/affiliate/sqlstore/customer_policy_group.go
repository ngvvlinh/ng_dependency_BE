package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
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

	paging  meta.Paging
	filters meta.Filters
}

func (s *CustomerPolicyGroupStore) Pred(pred interface{}) *CustomerPolicyGroupStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *CustomerPolicyGroupStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.CustomerPolicyGroup)(nil))
}

func (s *CustomerPolicyGroupStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *CustomerPolicyGroupStore) Paging(paging meta.Paging) *CustomerPolicyGroupStore {
	s.paging = paging
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
