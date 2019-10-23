package sqlstore

import (
	"context"

	"etop.vn/backend/com/main/inventory/convert"

	"etop.vn/api/main/inventory"
	"etop.vn/backend/com/main/inventory/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

var Sort = map[string]string{
	"id":         "",
	"created_at": "",
}

type InventoryFactory func(context.Context) *InventoryStore

func NewInventoryStore(db *cmsql.Database) InventoryFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *InventoryStore {
		return &InventoryStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type InventoryStore struct {
	query  func() cmsql.QueryInterface
	ft     InventoryVariantFilters
	paging *cm.Paging
	preds  []interface{}
}

func (s *InventoryStore) VariantID(id int64) *InventoryStore {
	s.preds = append(s.preds, s.ft.ByVariantID(id))
	return s
}

func (s *InventoryStore) Paging(page *cm.Paging) *InventoryStore {
	s.paging = page
	return s
}

func (s *InventoryStore) VariantIDs(ids ...int64) *InventoryStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", ids))
	return s
}

func (s *InventoryStore) ShopID(id int64) *InventoryStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *InventoryStore) UpdateInventoryVariant(args *model.InventoryVariant) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *InventoryStore) UpdateInventoryVariantAllDB(args *model.InventoryVariant) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *InventoryStore) UpdateInventoryVariantAll(args *inventory.InventoryVariant) error {
	var inventoryCore *model.InventoryVariant
	inventoryCore = convert.InventoryVariantToModel(args, inventoryCore)
	return s.UpdateInventoryVariantAllDB(inventoryCore)
}

func (s *InventoryStore) Create(inventory *model.InventoryVariant) error {
	query := s.query().Where(s.preds)
	return query.ShouldInsert(inventory)
}

func (s *InventoryStore) Get() (*inventory.InventoryVariant, error) {
	result, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	var inventoryVariant *inventory.InventoryVariant
	inventoryVariant = convert.InventoryVariantFromModel(result, inventoryVariant)
	return inventoryVariant, err
}

func (s *InventoryStore) GetDB() (*model.InventoryVariant, error) {
	query := s.query().Where(s.preds)
	var inventoryVariant model.InventoryVariant
	err := query.ShouldGet(&inventoryVariant)
	return &inventoryVariant, err
}

func (s *InventoryStore) ListInventoryDB() ([]*model.InventoryVariant, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, s.paging, Sort)
	if err != nil {
		return nil, err
	}
	var addrs model.InventoryVariants
	err = query.Find(&addrs)
	return addrs, err
}

func (s *InventoryStore) ListInventory() ([]*inventory.InventoryVariant, error) {
	result, err := s.ListInventoryDB()
	if err != nil {
		return nil, err
	}
	var inventoryVariants []*inventory.InventoryVariant
	inventoryVariants = convert.InventoryVariantsFromModel(result)
	return inventoryVariants, err
}
