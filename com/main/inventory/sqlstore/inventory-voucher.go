package sqlstore

import (
	"context"

	"etop.vn/backend/com/main/inventory/convert"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/backend/com/main/inventory/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type InventoryVoucherFactory func(context.Context) *InventoryVoucherStore

func NewInventoryVoucherStore(db *cmsql.Database) InventoryVoucherFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *InventoryVoucherStore {
		return &InventoryVoucherStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type InventoryVoucherStore struct {
	query  func() cmsql.QueryInterface
	ft     InventoryVoucherFilters
	preds  []interface{}
	paging *cm.Paging
}

func (s *InventoryVoucherStore) Paging(page *cm.Paging) *InventoryVoucherStore {
	s.paging = page
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
	query, err := sqlstore.LimitSort(query, s.paging, Sort)
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
