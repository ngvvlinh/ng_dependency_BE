package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/shopping/tradering"
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/suppliering/convert"
	"etop.vn/backend/com/shopping/suppliering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type SupplierStoreFactory func(ctx context.Context) *SupplierStore

func NewSupplierStore(db *cmsql.Database) SupplierStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *SupplierStore {
		return &SupplierStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type SupplierStore struct {
	ft ShopSupplierFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *SupplierStore) Paging(paging meta.Paging) *SupplierStore {
	s.paging = paging
	return s
}

func (s *SupplierStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *SupplierStore) Filters(filters meta.Filters) *SupplierStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *SupplierStore) ID(id int64) *SupplierStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SupplierStore) IDs(ids ...int64) *SupplierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *SupplierStore) ShopID(id int64) *SupplierStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *SupplierStore) OptionalShopID(id int64) *SupplierStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *SupplierStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopSupplier)(nil))
}

func (s *SupplierStore) CreateSupplier(supplier *suppliering.ShopSupplier) error {
	sqlstore.MustNoPreds(s.preds)
	trader := &customeringmodel.ShopTrader{
		ID:     supplier.ID,
		ShopID: supplier.ShopID,
		Type:   tradering.SupplierType,
	}
	supplierDB := new(model.ShopSupplier)
	if err := scheme.Convert(supplier, supplierDB); err != nil {
		return err
	}
	if _, err := s.query().Insert(trader, supplierDB); err != nil {
		return err
	}

	var tempSupplier model.ShopSupplier
	if err := s.query().Where(s.ft.ByID(supplier.ID), s.ft.ByShopID(supplier.ShopID)).ShouldGet(&tempSupplier); err != nil {
		return err
	}

	supplier.CreatedAt = tempSupplier.CreatedAt
	supplier.UpdatedAt = tempSupplier.UpdatedAt

	return nil
}

func (s *SupplierStore) UpdateSupplierDB(supplier *model.ShopSupplier) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ft.ByID(supplier.ID)).UpdateAll().ShouldUpdate(supplier)
	return err
}

func (s *SupplierStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_supplier").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *SupplierStore) DeleteCustomer() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopSupplier)(nil))
	return int(n), err
}

func (s *SupplierStore) GetSupplierDB() (*model.ShopSupplier, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var supplier model.ShopSupplier
	err := query.ShouldGet(&supplier)
	return &supplier, err
}

func (s *SupplierStore) GetSupplier() (supplierResult *suppliering.ShopSupplier, _ error) {
	supplier, err := s.GetSupplierDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_supplieringmodel_ShopSupplier_suppliering_ShopSupplier(supplier, supplierResult), nil
}

func (s *SupplierStore) ListSuppliersDB() ([]*model.ShopSupplier, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortSupplier)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterSupplier)
	if err != nil {
		return nil, err
	}

	var suppliers model.ShopSuppliers
	err = query.Find(&suppliers)
	return suppliers, err
}

func (s *SupplierStore) ListSuppliers() ([]*suppliering.ShopSupplier, error) {
	suppliers, err := s.ListSuppliersDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_supplieringmodel_ShopSuppliers_suppliering_ShopSuppliers(suppliers), nil
}