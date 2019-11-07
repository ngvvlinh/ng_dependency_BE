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

type ShopVariantSupplierStoreFactory func(ctx context.Context) *SupplierVariantStore

func NewSupplierVariantStore(db *cmsql.Database) ShopVariantSupplierStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *SupplierVariantStore {
		return &SupplierVariantStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type SupplierVariantStore struct {
	ft ShopVariantSupplierFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *SupplierVariantStore) Paging(paging meta.Paging) *SupplierVariantStore {
	s.paging = paging
	return s
}

func (s *SupplierVariantStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *SupplierVariantStore) Filters(filters meta.Filters) *SupplierVariantStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *SupplierVariantStore) VariantID(variantID int64) *SupplierVariantStore {
	s.preds = append(s.preds, s.ft.ByVariantID(variantID))
	return s
}

func (s *SupplierVariantStore) VariantIDs(variantIDs ...int64) *SupplierVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", variantIDs))
	return s
}

func (s *SupplierVariantStore) SupplierID(supplierID int64) *SupplierVariantStore {
	s.preds = append(s.preds, s.ft.BySupplierID(supplierID))
	return s
}

func (s *SupplierVariantStore) SupplierIDs(supplierIDs ...int64) *SupplierVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "Supplier_id", supplierIDs))
	return s
}

func (s *SupplierVariantStore) ShopID(id int64) *SupplierVariantStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *SupplierVariantStore) CreateVariantSupplier(vs *catalog.ShopVariantSupplier) error {
	sqlstore.MustNoPreds(s.preds)

	var vsDB = &model.ShopVariantSupplier{}
	if err := scheme.Convert(vs, vsDB); err != nil {
		return err
	}
	_, err := s.query().Insert(vsDB)
	return err
}

func (s *SupplierVariantStore) DeleteVariantSupplier() error {
	query := s.query().Where(s.preds)
	err := query.ShouldDelete(&model.ShopCategory{})
	return err
}

func (s *SupplierVariantStore) ListVariantSupplier() ([]*model.ShopVariantSupplier, error) {
	query := s.query().Where(s.preds)
	s.paging.Sort = []string{"-created_at"}

	var variantSuppliers model.ShopVariantSuppliers
	err := query.Find(&variantSuppliers)
	return variantSuppliers, err
}
