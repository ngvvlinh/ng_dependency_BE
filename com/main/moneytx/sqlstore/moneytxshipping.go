package sqlstore

import (
	"context"

	"etop.vn/api/main/moneytx"
	"etop.vn/backend/com/main/moneytx/convert"
	"etop.vn/backend/com/main/moneytx/model"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

type MoneyTxShippingStore struct {
	ft    MoneyTransactionShippingFilters
	query func() cmsql.QueryInterface
	preds []interface{}
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

func (s *MoneyTxShippingStore) ID(id dot.ID) *MoneyTxShippingStore {
	s.preds = append(s.preds, s.ft.ByID(id))
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
