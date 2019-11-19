package sqlstore

import (
	"context"

	"etop.vn/backend/com/main/inventory/convert"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/inventory/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
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
	query   cmsql.QueryFactory
	ft      InventoryVoucherFilters
	preds   []interface{}
	paging  meta.Paging
	filters meta.Filters
}

func (s *InventoryVoucherStore) Paging(page meta.Paging) *InventoryVoucherStore {
	s.paging = page
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

func (s *InventoryVoucherStore) ID(id int64) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *InventoryVoucherStore) IDs(ids ...int64) *InventoryVoucherStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *InventoryVoucherStore) Status(status etop.Status3) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *InventoryVoucherStore) ShopID(id int64) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *InventoryVoucherStore) RefID(id int64) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByRefID(id))
	return s
}

func (s *InventoryVoucherStore) RefType(refType string) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByRefType(refType))
	return s
}

func (s *InventoryVoucherStore) Type(inventoryVoucherType string) *InventoryVoucherStore {
	s.preds = append(s.preds, s.ft.ByType(inventoryVoucherType))
	return s
}

func (s *InventoryVoucherStore) RefIDs(ids ...int64) *InventoryVoucherStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "ref_id", ids))
	return s
}

func (s *InventoryVoucherStore) UpdateInventoryVoucher(inventory *model.InventoryVoucher) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(inventory)
}

func (s *InventoryVoucherStore) UpdateInventoryVoucherAllDB(inventory *model.InventoryVoucher) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(inventory)
}

func (s *InventoryVoucherStore) UpdateInventoryVoucherAll(inventory *inventory.InventoryVoucher) error {
	var updateValue *model.InventoryVoucher
	updateValue = convert.InventoryVoucherToModel(inventory, updateValue)
	return s.UpdateInventoryVoucherAllDB(updateValue)
}

func (s *InventoryVoucherStore) Create(inventoryVoucher *inventory.InventoryVoucher) error {
	var voucherDB *model.InventoryVoucher
	voucherDB = convert.InventoryVoucherToModel(inventoryVoucher, voucherDB)
	return s.CreateDB(voucherDB)
}

func (s *InventoryVoucherStore) CreateDB(voucher *model.InventoryVoucher) error {
	return s.query().ShouldInsert(voucher)
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
	var resultCore *inventory.InventoryVoucher
	resultCore = convert.InventoryVoucherFromModel(resultDB, resultCore)
	return resultCore, nil
}

func (s *InventoryVoucherStore) ListInventoryVoucherDB() ([]*model.InventoryVoucher, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = append(s.paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.paging, SortInventoryVoucher)
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
	resultCore = convert.InventoryVouchersFromModel(resultDB)
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
