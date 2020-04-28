package sqlstore

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

type ShopVariantSupplierStoreFactory func(ctx context.Context) *VariantSupplierStore

func NewVariantSupplierStore(db *cmsql.Database) ShopVariantSupplierStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *VariantSupplierStore {
		return &VariantSupplierStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type VariantSupplierStore struct {
	ft ShopVariantSupplierFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *VariantSupplierStore) WithPaging(paging meta.Paging) *VariantSupplierStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *VariantSupplierStore) Filters(filters meta.Filters) *VariantSupplierStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *VariantSupplierStore) VariantID(variantID dot.ID) *VariantSupplierStore {
	s.preds = append(s.preds, s.ft.ByVariantID(variantID))
	return s
}

func (s *VariantSupplierStore) VariantIDs(variantIDs ...dot.ID) *VariantSupplierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", variantIDs))
	return s
}

func (s *VariantSupplierStore) SupplierID(supplierID dot.ID) *VariantSupplierStore {
	s.preds = append(s.preds, s.ft.BySupplierID(supplierID))
	return s
}

func (s *VariantSupplierStore) SupplierIDs(supplierIDs ...dot.ID) *VariantSupplierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "Supplier_id", supplierIDs))
	return s
}

func (s *VariantSupplierStore) ShopID(id dot.ID) *VariantSupplierStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *VariantSupplierStore) CreateVariantSupplier(vs *catalog.ShopVariantSupplier) error {
	sqlstore.MustNoPreds(s.preds)

	var vsDB = &model.ShopVariantSupplier{}
	if err := scheme.Convert(vs, vsDB); err != nil {
		return err
	}
	_, err := s.query().Insert(vsDB)
	return err
}

func (s *VariantSupplierStore) DeleteVariantSupplier() error {
	query := s.query().Where(s.preds)
	err := query.ShouldDelete(&model.ShopVariantSupplier{})
	return err
}

func (s *VariantSupplierStore) ListVariantSupplier() ([]*model.ShopVariantSupplier, error) {
	query := s.query().Where(s.preds)
	s.Paging.Sort = []string{"-created_at"}

	var variantSuppliers model.ShopVariantSuppliers
	err := query.Find(&variantSuppliers)
	return variantSuppliers, err
}
