package sqlstore

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/meta"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/inventory/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type InventoryVoucherFactory func(context.Context) *InventoryVoucherStore

func NewInventoryVoucherStore(db *cmsql.Database) InventoryVoucherFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *InventoryVoucherStore {
		return &InventoryVoucherStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type InventoryVoucherStore struct {
	query cmsql.QueryFactory
	ft    InventoryVoucherFilters
	preds []interface{}
	sqlstore.Paging
	filters meta.Filters
}

func (s *InventoryVoucherStore) WithPaging(paging meta.Paging) *InventoryVoucherStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *InventoryVoucherStore) Filters(filters meta.Filters) *InventoryVoucherStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *InventoryVoucherStore) ID(id dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *InventoryVoucherStore) IDs(ids ...dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *InventoryVoucherStore) Status(status status3.Status) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *InventoryVoucherStore) ShopID(id dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *InventoryVoucherStore) VariantID(id dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, sq.NewExpr("variant_ids @> ?", core.Array{V: []dot.ID{id}}))
	return s
}

func (s *InventoryVoucherStore) RefID(id dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByRefID(id))
	return s
}

func (s *InventoryVoucherStore) RefType(refType inventory_voucher_ref.InventoryVoucherRef) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByRefType(refType))
	return s
}

func (s *InventoryVoucherStore) Type(inventoryVoucherType inventory_type.InventoryVoucherType) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByType(inventoryVoucherType))
	return s
}

func (s *InventoryVoucherStore) RefIDs(ids ...dot.ID) *InventoryVoucherStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "ref_id", ids))
	return s
}

func (s *InventoryVoucherStore) UpdateInventoryVoucher(inventoryVoucher *model.InventoryVoucher) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(inventoryVoucher)
}

func (s *InventoryVoucherStore) UpdateInventoryVoucherAllDB(inventoryVoucher *model.InventoryVoucher) error {
	query := s.query().Where(s.preds)
	var variantIDs []dot.ID
	var productIDs []dot.ID
	for _, value := range inventoryVoucher.Lines {
		variantIDs = append(variantIDs, value.VariantID)
		productIDs = append(productIDs, value.ProductID)
	}
	inventoryVoucher.VariantIDs = variantIDs
	inventoryVoucher.ProductIDs = productIDs
	return query.UpdateAll().ShouldUpdate(inventoryVoucher)
}

func (s *InventoryVoucherStore) UpdateInventoryVoucherAll(inventory *inventory.InventoryVoucher) error {
	updateValue := &model.InventoryVoucher{}
	if err := scheme.Convert(inventory, updateValue); err != nil {
		return err
	}

	return s.UpdateInventoryVoucherAllDB(updateValue)
}

func (s *InventoryVoucherStore) Create(inventoryVoucher *inventory.InventoryVoucher) error {
	voucherDB := &model.InventoryVoucher{}
	if err := scheme.Convert(inventoryVoucher, voucherDB); err != nil {
		return err
	}
	return s.CreateDB(voucherDB)
}

func (s *InventoryVoucherStore) CreateDB(inventoryVoucher *model.InventoryVoucher) error {
	var variantIDs []dot.ID
	var productIDs []dot.ID
	for _, value := range inventoryVoucher.Lines {
		variantIDs = append(variantIDs, value.VariantID)
		productIDs = append(productIDs, value.ProductID)
	}
	inventoryVoucher.VariantIDs = variantIDs
	inventoryVoucher.ProductIDs = productIDs
	return s.query().ShouldInsert(inventoryVoucher)
}

func (s *InventoryVoucherStore) GetDB() (*model.InventoryVoucher, error) {
	query := s.query().Where(s.preds)
	var inventoryVoucher model.InventoryVoucher
	err := query.ShouldGet(&inventoryVoucher)
	return &inventoryVoucher, err
}

func (s *InventoryVoucherStore) Get() (*inventory.InventoryVoucher, error) {
	resultDB, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	resultCore := &inventory.InventoryVoucher{}
	if err := scheme.Convert(resultDB, resultCore); err != nil {
		return nil, err
	}
	return resultCore, nil
}

func (s *InventoryVoucherStore) ListInventoryVoucherDB() ([]*model.InventoryVoucher, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortInventoryVoucher)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterInventoryVoucher)
	if err != nil {
		return nil, err
	}
	var result model.InventoryVouchers
	err = query.Find(&result)
	return result, err
}

func (s *InventoryVoucherStore) ListInventoryVoucher() ([]*inventory.InventoryVoucher, error) {
	resultDB, err := s.ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	var resultCore []*inventory.InventoryVoucher
	if err := scheme.Convert(resultDB, &resultCore); err != nil {
		return nil, err
	}
	return resultCore, nil
}

func (s *InventoryVoucherStore) GetInventoryVoucherByMaximumCodeNorm() (*model.InventoryVoucher, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var inventoryVoucher model.InventoryVoucher
	if err := query.ShouldGet(&inventoryVoucher); err != nil {
		return nil, err
	}
	return &inventoryVoucher, nil
}
