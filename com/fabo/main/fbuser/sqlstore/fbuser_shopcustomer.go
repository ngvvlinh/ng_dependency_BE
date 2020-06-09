package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalUserShopCustomerStoreFactory func(ctx context.Context) *FbExternalUserShopCustomerStore

func NewFbExternalUserShopCustomerStore(db *cmsql.Database) FbExternalUserShopCustomerStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalUserShopCustomerStore {
		return &FbExternalUserShopCustomerStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalUserShopCustomerStore struct {
	db *cmsql.Database
	ft FbExternalUserShopCustomerFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalUserShopCustomerStore) WithPaging(paging meta.Paging) *FbExternalUserShopCustomerStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalUserShopCustomerStore) Filters(filters meta.Filters) *FbExternalUserShopCustomerStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *FbExternalUserShopCustomerStore) ShopCustomerID(customerID dot.ID) *FbExternalUserShopCustomerStore {
	s.preds = append(s.preds, s.ft.ByCustomerID(customerID))
	return s
}

func (s *FbExternalUserShopCustomerStore) ShopCustomerIDs(customerIDs []dot.ID) *FbExternalUserShopCustomerStore {
	s.preds = append(s.preds, sq.In("customer_id", customerIDs))
	return s
}

func (s *FbExternalUserShopCustomerStore) ShopID(shopID dot.ID) *FbExternalUserShopCustomerStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *FbExternalUserShopCustomerStore) FbExternalUserID(fbExternalUserID string) *FbExternalUserShopCustomerStore {
	s.preds = append(s.preds, s.ft.ByFbExternalUserID(fbExternalUserID))
	return s
}

func (s *FbExternalUserShopCustomerStore) FbExternalUserIDs(fbExternalUserIDs []string) *FbExternalUserShopCustomerStore {
	s.preds = append(s.preds, sq.In("fb_external_user_id", fbExternalUserIDs))
	return s
}

func (s *FbExternalUserShopCustomerStore) ListFbExternalUserDB() ([]*model.FbExternalUserShopCustomer, error) {
	query := s.query().Where(s.preds)
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalUser, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalUser)
	if err != nil {
		return nil, err
	}
	var fbExternalUserShopCustomers model.FbExternalUserShopCustomers
	err = query.Find(&fbExternalUserShopCustomers)
	if err != nil {
		return nil, err
	}
	return fbExternalUserShopCustomers, err
}

func (s *FbExternalUserShopCustomerStore) ListFbExternalUser() ([]*fbusering.FbExternalUserShopCustomer, error) {
	fbExternalUser, err := s.ListFbExternalUserDB()
	if err != nil {
		return nil, err
	}
	var result = []*fbusering.FbExternalUserShopCustomer{}
	err = scheme.Convert(fbExternalUser, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalUserShopCustomerStore) CreateFbExternalUserShopCustomer(FbExternalUserShopCustomer *fbusering.FbExternalUserShopCustomer) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUserDB := new(model.FbExternalUserShopCustomer)
	if err := scheme.Convert(FbExternalUserShopCustomer, fbExternalUserDB); err != nil {
		return err
	}
	_, err := s.query().Insert(fbExternalUserDB)
	return err
}

func (s *FbExternalUserShopCustomerStore) DeleteFbExternalUserShopCustomer() error {
	query := s.query().Where(s.preds)
	_, err := query.Delete(&model.FbExternalUserShopCustomer{})
	return err
}
