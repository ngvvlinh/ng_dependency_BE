package sqlstore

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/moneytx/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type MoneyTxShippingExternalStore struct {
	ft     MoneyTransactionShippingExternalFilters
	ftLine MoneyTransactionShippingExternalLineFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging
}

func (s *MoneyTxShippingExternalStore) extend() *MoneyTxShippingExternalStore {
	s.ftLine.prefix = "m"
	return s
}

type MoneyTxShippingExternalStoreFactory func(ctx context.Context) *MoneyTxShippingExternalStore

func NewMoneyTxShippingExternalStore(db *cmsql.Database) MoneyTxShippingExternalStoreFactory {
	return func(ctx context.Context) *MoneyTxShippingExternalStore {
		return &MoneyTxShippingExternalStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *MoneyTxShippingExternalStore) ID(id dot.ID) *MoneyTxShippingExternalStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalExtendedDB() (*model.MoneyTransactionShippingExternalExtended, error) {
	query := s.query().Where(s.preds)
	moneyTx := &model.MoneyTransactionShippingExternal{}
	err := query.ShouldGet(moneyTx)

	lines, err := s.GetMoneyTxShippingExternalLineExtendedDBs([]dot.ID{moneyTx.ID})
	res := &model.MoneyTransactionShippingExternalExtended{
		MoneyTransactionShippingExternal: moneyTx,
		Lines:                            lines,
	}
	return res, err
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalExtended() (*moneytx.MoneyTransactionShippingExternalExtended, error) {
	moneyTx, err := s.GetMoneyTxShippingExternalExtendedDB()
	if err != nil {
		return nil, err
	}
	var res moneytx.MoneyTransactionShippingExternalExtended
	if err := scheme.Convert(moneyTx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *MoneyTxShippingExternalStore) CreateMoneyTxShippingExternal(moneyTx *moneytx.MoneyTransactionShippingExternal) error {
	sqlstore.MustNoPreds(s.preds)
	moneyTxDB := &model.MoneyTransactionShippingExternal{}
	if err := scheme.Convert(moneyTx, moneyTxDB); err != nil {
		return err
	}
	if err := s.query().ShouldInsert(moneyTxDB); err != nil {
		return err
	}
	return nil
}

func (s *MoneyTxShippingExternalStore) CreateMoneyTxShippingExternalLine(line *moneytx.MoneyTransactionShippingExternalLine) error {
	sqlstore.MustNoPreds(s.preds)
	lineDB := &model.MoneyTransactionShippingExternalLine{}
	if err := scheme.Convert(line, lineDB); err != nil {
		return err
	}
	if err := s.query().ShouldInsert(lineDB); err != nil {
		return err
	}
	return nil
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalLineExtendedDBs(moneyTxShippingExternalIDs []dot.ID) (lines []*model.MoneyTransactionShippingExternalLineExtended, _ error) {
	query := s.extend().query().Where(sq.PrefixedIn(&s.ftLine.prefix, "money_transaction_shipping_external_id", moneyTxShippingExternalIDs))
	if err := query.Find((*model.MoneyTransactionShippingExternalLineExtendeds)(&lines)); err != nil {
		return nil, err
	}
	return
}
