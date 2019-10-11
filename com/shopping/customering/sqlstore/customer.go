package sqlstore

import (
	"context"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/shopping/tradering"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	"etop.vn/backend/com/shopping/customering/convert"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type CustomerStoreFactory func(context.Context) *CustomerStore

func NewCustomerStore(db cmsql.Database) CustomerStoreFactory {
	return func(ctx context.Context) *CustomerStore {
		return &CustomerStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type CustomerStore struct {
	ft ShopCustomerFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *CustomerStore) Paging(paging meta.Paging) *CustomerStore {
	s.paging = paging
	return s
}

func (s *CustomerStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *CustomerStore) Filters(filters meta.Filters) *CustomerStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *CustomerStore) ID(id int64) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CustomerStore) IDs(ids ...int64) *CustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *CustomerStore) ShopID(id int64) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *CustomerStore) OptionalShopID(id int64) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *CustomerStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopCustomer)(nil))
}

func (s *CustomerStore) CreateCustomer(customer *customering.ShopCustomer) error {
	sqlstore.MustNoPreds(s.preds)
	trader := &model.ShopTrader{
		ID:     customer.ID,
		ShopID: customer.ShopID,
		Type:   tradering.CustomerType,
	}
	customerDB := convert.ShopCustomerDB(customer)
	_, err := s.query().Insert(trader, customerDB)
	return CheckErrorCustomer(err, customerDB.Email, customerDB.Phone)
}

func (s *CustomerStore) UpdateCustomerDB(customer *model.ShopCustomer) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ft.ByID(customer.ID)).UpdateAll().ShouldUpdate(customer)
	return err
}

func (s *CustomerStore) PatchCustomerDB(customer *model.ShopCustomer) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	n, err := query.Update(customer)
	return int(n), err
}

func (s *CustomerStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_customer").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *CustomerStore) DeleteCustomer() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCustomer)(nil))
	return int(n), err
}

func (s *CustomerStore) GetCustomerDB() (*model.ShopCustomer, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var customer model.ShopCustomer
	err := query.ShouldGet(&customer)
	return &customer, err
}

func (s *CustomerStore) GetCustomer() (*customering.ShopCustomer, error) {
	customer, err := s.GetCustomerDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopCustomer(customer), nil
}

func (s *CustomerStore) ListCustomersDB() ([]*model.ShopCustomer, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortCustomer, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomer)
	if err != nil {
		return nil, err
	}

	var customers model.ShopCustomers
	err = query.Find(&customers)
	return customers, err
}

func (s *CustomerStore) ListCustomers() ([]*customering.ShopCustomer, error) {
	customers, err := s.ListCustomersDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopCustomers(customers), nil
}

func CheckErrorCustomer(e error, email string, phone string) error {
	if e != nil {
		errMsg := e.Error()
		switch {
		case strings.Contains(errMsg, "shop_customer_shop_id_email_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "Khách hàng với email %v đã tồn tại", email)
		case strings.Contains(errMsg, "shop_customer_shop_id_phone_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "Khách hàng với số điện thoại %v đã tồn tại", phone)
		}
	}
	return e
}
