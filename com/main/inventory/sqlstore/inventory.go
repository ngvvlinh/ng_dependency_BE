package sqlstore

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/backend/com/main/inventory/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
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
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type InventoryStore struct {
	query cmsql.QueryFactory
	ft    InventoryVariantFilters
	sqlstore.Paging
	preds []interface{}
}

func (s *InventoryStore) VariantID(id dot.ID) *InventoryStore {
	s.preds = append(s.preds, s.ft.ByVariantID(id))
	return s
}

func (s *InventoryStore) WithPaging(paging *cm.Paging) *InventoryStore {
	s.Paging.WithPaging(*paging)
	return s
}

func (s *InventoryStore) VariantIDs(ids ...dot.ID) *InventoryStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", ids))
	return s
}

func (s *InventoryStore) ShopID(id dot.ID) *InventoryStore {
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
	inventoryCore := &model.InventoryVariant{}
	if err := scheme.Convert(args, inventoryCore); err != nil {
		return err
	}
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
	inventoryVariant := &inventory.InventoryVariant{}
	if err := scheme.Convert(result, inventoryVariant); err != nil {
		return nil, err
	}
	return inventoryVariant, nil
}

func (s *InventoryStore) GetDB() (*model.InventoryVariant, error) {
	query := s.query().Where(s.preds)
	var inventoryVariant model.InventoryVariant
	err := query.ShouldGet(&inventoryVariant)
	return &inventoryVariant, err
}

func (s *InventoryStore) ListInventoryDB() ([]*model.InventoryVariant, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.Paging, Sort)
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
	if err := scheme.Convert(result, &inventoryVariants); err != nil {
		return nil, err
	}
	return inventoryVariants, nil
}
