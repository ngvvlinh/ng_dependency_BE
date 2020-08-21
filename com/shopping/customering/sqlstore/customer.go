package sqlstore

import (
	"context"
	"strings"
	"time"

	"o.o/api/meta"
	"o.o/api/shopping/customering"
	"o.o/api/shopping/tradering"
	"o.o/api/top/types/etc/customer_type"
	"o.o/backend/com/shopping/customering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type CustomerStoreFactory func(context.Context) *CustomerStore

func NewCustomerStore(db *cmsql.Database) CustomerStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CustomerStore {
		return &CustomerStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CustomerStore struct {
	ft ShopCustomerFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *CustomerStore) WithPaging(paging meta.Paging) *CustomerStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *CustomerStore) Filters(filters meta.Filters) *CustomerStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *CustomerStore) FullTextSearchFullName(name filter.FullTextSearch) *CustomerStore {
	s.preds = append(s.preds, s.ft.Filter(`full_name_norm @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
}

func (s *CustomerStore) ID(id dot.ID) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CustomerStore) IDs(ids ...dot.ID) *CustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *CustomerStore) ShopID(id dot.ID) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *CustomerStore) ShopIDs(ids ...dot.ID) *CustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "shop_id", ids))
	return s
}

func (s *CustomerStore) Type(typ customer_type.CustomerType) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByType(typ))
	return s
}

func (s *CustomerStore) Code(code string) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *CustomerStore) ExternalID(externalID string) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *CustomerStore) CodeNorm(codeNorm int) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByCodeNorm(codeNorm))
	return s
}

func (s *CustomerStore) Phone(phone string) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByPhone(phone))
	return s
}

func (s *CustomerStore) Phones(phones ...string) *CustomerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "phone", phones))
	return s
}

func (s *CustomerStore) Email(email string) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByEmail(email))
	return s
}

func (s *CustomerStore) OptionalShopID(id dot.ID) *CustomerStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *CustomerStore) Count() (_ int, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomer)
	if err != nil {
		return 0, err
	}
	return query.Count((*model.ShopCustomer)(nil))
}

func (s *CustomerStore) CreateCustomer(customer *model.ShopCustomer) error {
	sqlstore.MustNoPreds(s.preds)
	trader := &model.ShopTrader{
		ID:     customer.ID,
		ShopID: customer.ShopID,
		Type:   tradering.CustomerType,
	}
	_, err := s.query().Insert(trader, customer)
	return CheckErrorCustomer(err, customer.Email, customer.Phone, customer.ExternalID, customer.ExternalCode)
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
	return n, err
}

func (s *CustomerStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_customer").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *CustomerStore) DeleteCustomer() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCustomer)(nil))
	return n, err
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
	result := &customering.ShopCustomer{}
	err = scheme.Convert(customer, result)
	if err != nil {
		return nil, err
	}
	result.Deleted = !customer.DeletedAt.IsZero()
	return result, err
}

func (s *CustomerStore) GetCustomerByMaximumCodeNorm() (*model.ShopCustomer, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var customer model.ShopCustomer
	if err := query.ShouldGet(&customer); err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerStore) ListCustomersDB() ([]*model.ShopCustomer, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortCustomer, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCustomer)
	if err != nil {
		return nil, err
	}

	var customers model.ShopCustomers
	err = query.Find(&customers)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(customers)
	return customers, nil
}

func (s *CustomerStore) ListCustomers() (result []*customering.ShopCustomer, err error) {
	customers, err := s.ListCustomersDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(customers, &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(customers); i++ {
		result[i].Deleted = !customers[i].DeletedAt.IsZero()
	}
	return
}

func CheckErrorCustomer(e error, email, phone, externalID, externalCode string) error {
	if e != nil {
		errMsg := e.Error()
		switch {
		case strings.Contains(errMsg, "shop_customer_shop_id_email_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "Khách hàng với email %v đã tồn tại", email)
		case strings.Contains(errMsg, "shop_customer_shop_id_phone_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "Khách hàng với số điện thoại %v đã tồn tại", phone)
		case strings.Contains(errMsg, "shop_customer_shop_id_external_id_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "external_id %v đã tồn tại", externalID)
		case strings.Contains(errMsg, "shop_customer_shop_id_external_code_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "external_code %v đa tồn tại", externalCode)
		}
	}
	return e
}

func (s *CustomerStore) IncludeDeleted() *CustomerStore {
	s.includeDeleted = true
	return s
}
