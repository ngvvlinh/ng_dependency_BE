package sqlstore

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/moneytx/convert"
	"etop.vn/backend/com/main/moneytx/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

const MaxLimitedNumber = 10000

type MoneyTxShippingStore struct {
	ft      MoneyTransactionShippingFilters
	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

type MoneyTxShippingStoreFactory func(ctx context.Context) *MoneyTxShippingStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewMoneyTxShippingStore(db *cmsql.Database) MoneyTxShippingStoreFactory {
	return func(ctx context.Context) *MoneyTxShippingStore {
		return &MoneyTxShippingStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *MoneyTxShippingStore) extend() *MoneyTxShippingStore {
	s.ft.prefix = "m"
	return s
}

func (s *MoneyTxShippingStore) WithPaging(paging meta.Paging) *MoneyTxShippingStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *MoneyTxShippingStore) ID(id dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *MoneyTxShippingStore) IDs(ids ...dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *MoneyTxShippingStore) OptionalMoneyTxShippingEtopID(id dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionShippingEtopID(id).Optional())
	return s
}

func (s *MoneyTxShippingStore) Filters(filters meta.Filters) *MoneyTxShippingStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *MoneyTxShippingStore) ShopID(shopID dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *MoneyTxShippingStore) OptionalShopID(shopID dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID).Optional())
	return s
}

func (s *MoneyTxShippingStore) MoneyTxShippingExternalID(id dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionShippingExternalID(id))
	return s
}

func (s *MoneyTxShippingStore) MoneyTxShippingEtopID(id dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionShippingEtopID(id))
	return s
}

func (s *MoneyTxShippingStore) GetMoneyTxShippingDB() (*model.MoneyTransactionShipping, error) {
	query := s.query().Where(s.preds)
	var moneyTx model.MoneyTransactionShipping
	err := query.ShouldGet(&moneyTx)
	return &moneyTx, err
}

func (s *MoneyTxShippingStore) GetMoneyTxShipping() (*moneytx.MoneyTransactionShipping, error) {
	moneyTx, err := s.GetMoneyTxShippingDB()
	if err != nil {
		return nil, err
	}
	var res moneytx.MoneyTransactionShipping
	if err := scheme.Convert(moneyTx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *MoneyTxShippingStore) ListMoneyTxShippingsDB() ([]*model.MoneyTransactionShipping, error) {
	var res []*model.MoneyTransactionShipping
	if isSupportSpecialFilters(s.filters) {
		moneyTxs, err := s.ListMoneyTxShippingFtShopsDB()
		if err != nil {
			return nil, err
		}
		for _, moneyTx := range moneyTxs {
			res = append(res, moneyTx.MoneyTransactionShipping)
		}
		return res, nil
	}
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	if s.Paging.Limit == 0 {
		s.Paging.Limit = MaxLimitedNumber
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortMoneyTx, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, filterMoneyTxShippingWhitelist)
	if err != nil {
		return nil, err
	}
	if err := query.Find((*model.MoneyTransactionShippings)(&res)); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *MoneyTxShippingStore) ListMoneyTxShippingFtShopsDB() ([]*model.MoneyTransactionShippingFtShop, error) {
	query := s.extend().query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}

	query, err := sqlstore.LimitSort(query, &s.Paging, SortMoneyTx, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, filterMoneyTxShippingExtendedWhitelist)
	if err != nil {
		return nil, err
	}

	var moneyTxs model.MoneyTransactionShippingFtShops
	if err := query.Find(&moneyTxs); err != nil {
		return nil, err
	}
	return moneyTxs, nil
}

var specialShopFilters = []string{
	"shop.money_transaction_rrule",
	"shop.bank_account",
}

func isSupportSpecialFilters(filters meta.Filters) bool {
	for _, filter := range filters.ListFilters() {
		if cm.StringsContain(specialShopFilters, filter.Name) {
			return true
		}
	}
	return false
}

func (s *MoneyTxShippingStore) ListMoneyTxShippings() (res []*moneytx.MoneyTransactionShipping, _ error) {
	moneyTxs, err := s.ListMoneyTxShippingsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(moneyTxs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingStore) CreateMoneyTxShippingDB(moneyTx *model.MoneyTransactionShipping) error {
	sqlstore.MustNoPreds(s.preds)
	if moneyTx.ID == 0 {
		moneyTx.ID = cm.NewID()
	}
	return s.query().ShouldInsert(moneyTx)
}

func (s *MoneyTxShippingStore) CreateMoneyTxShipping(moneyTx *moneytx.MoneyTransactionShipping) error {
	sqlstore.MustNoPreds(s.preds)
	moneyTxDB := &model.MoneyTransactionShipping{}
	if err := scheme.Convert(moneyTx, moneyTxDB); err != nil {
		return err
	}
	return s.CreateMoneyTxShippingDB(moneyTxDB)
}

func (s *MoneyTxShippingStore) UpdateMoneyTxShippingDB(moneyTx *model.MoneyTransactionShipping) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Where(s.preds).ShouldUpdate(moneyTx)
}

func (s *MoneyTxShippingStore) DeleteMoneyTxShipping(id dot.ID) error {
	query := s.query().Where(s.ft.ByID(id))
	return query.ShouldDelete((&model.MoneyTransactionShipping{}))
}

type UpdateMoneyTxShippingStatisticsArgs struct {
	ID          dot.ID
	TotalCOD    dot.NullInt
	TotalAmount dot.NullInt
	TotalOrders dot.NullInt
}

func (s *MoneyTxShippingStore) UpdateMoneyTxShippingStatistics(args *UpdateMoneyTxShippingStatisticsArgs) error {
	var update = make(map[string]interface{})
	if args.TotalAmount.Valid {
		update["total_amount"] = args.TotalAmount.Int
	}
	if args.TotalCOD.Valid {
		update["total_cod"] = args.TotalCOD.Int
	}
	if args.TotalOrders.Valid {
		update["total_orders"] = args.TotalOrders.Int
	}
	if args.TotalOrders.Valid && args.TotalOrders.Int == 0 {
		update["status"] = status3.Z
	}
	if len(update) > 0 {
		return s.query().Table("money_transaction_shipping").Where(s.ft.ByID(args.ID)).ShouldUpdateMap(update)
	}
	return nil
}

func (s *MoneyTxShippingStore) RemoveMoneyTxShippingMoneyTxShippingEtopID() error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Table("money_transaction_shipping").Where(s.preds).ShouldUpdateMap(map[string]interface{}{
		"money_transaction_shipping_etop_id": nil,
	})
}
