package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/conversion"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/com/main/catalog/convert"

	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
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
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *VariantSupplierStore) Paging(paging meta.Paging) *VariantSupplierStore {
	s.paging = paging
	return s
}

func (s *VariantSupplierStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *VariantSupplierStore) Filters(filters meta.Filters) *VariantSupplierStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *VariantSupplierStore) VariantID(variantID int64) *VariantSupplierStore {
	s.preds = append(s.preds, s.ft.ByVariantID(variantID))
	return s
}

func (s *VariantSupplierStore) VariantIDs(variantIDs ...int64) *VariantSupplierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", variantIDs))
	return s
}

func (s *VariantSupplierStore) SupplierID(supplierID int64) *VariantSupplierStore {
	s.preds = append(s.preds, s.ft.BySupplierID(supplierID))
	return s
}

func (s *VariantSupplierStore) SupplierIDs(supplierIDs ...int64) *VariantSupplierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "Supplier_id", supplierIDs))
	return s
}

func (s *VariantSupplierStore) ShopID(id int64) *VariantSupplierStore {
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
	s.paging.Sort = []string{"-created_at"}

	var variantSuppliers model.ShopVariantSuppliers
	err := query.Find(&variantSuppliers)
	return variantSuppliers, err
}
