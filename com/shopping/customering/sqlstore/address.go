package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/shopping/addressing"
	"etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type AddressStoreFactory func(context.Context) *AddressStore

func NewAddressStore(db *cmsql.Database) AddressStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AddressStore {
		return &AddressStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AddressStore struct {
	ft ShopTraderAddressFilters

	query cmsql.QueryFactory
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

func (s *AddressStore) ID(id dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *AddressStore) IDs(ids ...dot.ID) *AddressStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *AddressStore) ShopID(shopID dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *AddressStore) ShopTraderID(shopID, traderID dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	s.preds = append(s.preds, s.ft.ByTraderID(traderID))
	return s
}

func (s *AddressStore) IsDefault(isDefault bool) *AddressStore {
	s.preds = append(s.preds, s.ft.ByIsDefault(isDefault))
	return s
}

func (s *AddressStore) UpdateStatusAddresses(shopID, traderID dot.ID, isDefault bool) error {
	_, err := s.query().Where(
		s.ft.ByShopID(shopID),
		s.ft.ByTraderID(traderID)).
		Table("shop_trader_address").
		UpdateMap(map[string]interface{}{
			"is_default": isDefault,
		})

	return err
}

func (s *AddressStore) SetDefaultAddress(ID, shopID, traderID dot.ID) (int, error) {
	sqlstore.MustNoPreds(s.preds)
	updated, err := s.query().Where(
		s.ft.ByID(ID),
		s.ft.ByShopID(shopID),
		s.ft.ByTraderID(traderID)).
		Table("shop_trader_address").
		UpdateMap(map[string]interface{}{
			"is_default": true,
		})
	return updated, err
}

func (s *AddressStore) CreateAddress(addr *addressing.ShopTraderAddress) error {
	sqlstore.MustNoPreds(s.preds)
	addrDB := &model.ShopTraderAddress{}
	if err := scheme.Convert(addr, addrDB); err != nil {
		return err
	}
	_, err := s.query().Insert(addrDB)
	return err
}

func (s *AddressStore) UpdateAddressDB(addr *model.ShopTraderAddress) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(
		s.ft.ByID(addr.ID),
		s.ft.ByShopID(addr.ShopID),
		s.ft.ByTraderID(addr.TraderID),
	).UpdateAll().ShouldUpdate(addr)
	return err
}

func (s *AddressStore) SoftDelete() (deleted int, _ error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_trader_address").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *AddressStore) GetAddressDB() (*model.ShopTraderAddress, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var address model.ShopTraderAddress
	err := query.ShouldGet(&address)
	return &address, err
}

func (s *AddressStore) GetAddress() (*addressing.ShopTraderAddress, error) {
	address, err := s.GetAddressDB()
	if err != nil {
		return nil, err
	}
	result := &addressing.ShopTraderAddress{}
	err = scheme.Convert(address, result)
	return result, err
}

func (s *AddressStore) ListAddressesDB() ([]*model.ShopTraderAddress, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var addrs model.ShopTraderAddresses
	err := query.Find(&addrs)
	return addrs, err
}

func (s *AddressStore) ListAddresses() (result []*addressing.ShopTraderAddress, err error) {
	addrs, err := s.ListAddressesDB()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(addrs, &result)
	return
}
