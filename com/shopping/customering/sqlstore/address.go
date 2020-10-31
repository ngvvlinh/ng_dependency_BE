package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/backend/com/shopping/customering/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/capi/filter"
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
	ft              ShopTraderAddressFilters
	addressSearchFt ShopTraderAddressSearchFilters
	query           cmsql.QueryFactory
	preds           []interface{}
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
	includeKhachLe bool
}

func (s *AddressStore) extend() *AddressStore {
	s.ft.prefix = "sta"
	s.addressSearchFt.prefix = "stas"
	return s
}

func (s *AddressStore) WithPaging(paging meta.Paging) *AddressStore {
	s.Paging.WithPaging(paging)
	return s
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

func (s *AddressStore) TraderID(traderID dot.ID) *AddressStore {
	s.preds = append(s.preds, s.ft.ByTraderID(traderID))
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

func (s *AddressStore) Count() (_ int, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopTraderAddressExtendeds)(nil))
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

func (s *AddressStore) IncludeKhachLe() *AddressStore {
	s.includeKhachLe = true
	return s
}

func (s *AddressStore) CreateAddress(addr *addressing.ShopTraderAddress) error {
	sqlstore.MustNoPreds(s.preds)
	addrDB := &model.ShopTraderAddress{}
	if err := scheme.Convert(addr, addrDB); err != nil {
		return err
	}
	_, err := s.query().Insert(addrDB)
	if err != nil {
		return err
	}
	addrSearchDB := &model.ShopTraderAddressSearch{
		ID:        addrDB.ID,
		PhoneNorm: validate.NormalizeSearchPhone(addrDB.Phone),
	}
	_, err = s.query().Insert(addrSearchDB)
	return err
}

func (s *AddressStore) UpdateAddressDB(addr *model.ShopTraderAddress) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(
		s.ft.ByID(addr.ID),
		s.ft.ByShopID(addr.ShopID),
		s.ft.ByTraderID(addr.TraderID),
	).UpdateAll().ShouldUpdate(addr)
	if err != nil {
		return err
	}
	addrSearchDB := &model.ShopTraderAddressSearch{
		PhoneNorm: validate.NormalizeSearchPhone(addr.Phone),
	}
	err = s.query().Where(
		s.ft.ByID(addr.ID),
	).ShouldUpdate(addrSearchDB)
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
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopTraderAddress, s.ft.prefix)
	if err != nil {
		return nil, err
	}

	var addrsEx model.ShopTraderAddressExtendeds
	s.Paging.Apply(addrsEx)
	err = query.Find(&addrsEx)
	if err != nil {
		return nil, err
	}

	var addrs []*model.ShopTraderAddress
	for _, v := range addrsEx {
		// bỏ qua address của Khách Lẻ (https://github.com/etopvn/one/issues/2730)
		if !s.includeKhachLe && v.ShopTraderAddress.TraderID == customering.CustomerAnonymous {
			continue
		}
		addrs = append(addrs, v.ShopTraderAddress)
	}
	return addrs, nil
}

func (s *AddressStore) ListAddresses() (result []*addressing.ShopTraderAddress, err error) {
	addrs, err := s.ListAddressesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(addrs, &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(addrs); i++ {
		result[i].Deleted = !addrs[i].DeletedAt.IsZero()
	}
	return
}

func (s *AddressStore) IncludeDeleted() *AddressStore {
	s.includeDeleted = true
	return s
}

func (s *AddressStore) SearchPhone(phone string) *AddressStore {
	s.preds = append(s.preds, s.addressSearchFt.Filter("phone_norm @@ ?::tsquery", phone))
	return s
}

func (s *AddressStore) FullTextSearchPhone(phone filter.FullTextSearch) *AddressStore {
	s.preds = append(s.preds, s.addressSearchFt.Filter("phone_norm @@ ?::tsquery", validate.NormalizeFullTextSearchQueryAnd(phone)))
	return s
}
