package sqlstore

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/moneytx/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type MoneyTxShippingEtopStore struct {
	ft      MoneyTransactionShippingEtopFilters
	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

type MoneyTxShippingEtopStoreFactory func(ctx context.Context) *MoneyTxShippingEtopStore

func NewMoneyTxShippingEtopStore(db *cmsql.Database) MoneyTxShippingEtopStoreFactory {
	return func(ctx context.Context) *MoneyTxShippingEtopStore {
		return &MoneyTxShippingEtopStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *MoneyTxShippingEtopStore) WithPaging(paging meta.Paging) *MoneyTxShippingEtopStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *MoneyTxShippingEtopStore) Filters(filters meta.Filters) *MoneyTxShippingEtopStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *MoneyTxShippingEtopStore) ID(id dot.ID) *MoneyTxShippingEtopStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *MoneyTxShippingEtopStore) IDs(ids ...dot.ID) *MoneyTxShippingEtopStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *MoneyTxShippingEtopStore) Status(status status3.Status) *MoneyTxShippingEtopStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *MoneyTxShippingEtopStore) GetMoneyTxShippingEtopDB() (*model.MoneyTransactionShippingEtop, error) {
	query := s.query().Where(s.preds)
	var moneyTx model.MoneyTransactionShippingEtop
	if err := query.ShouldGet(&moneyTx); err != nil {
		return nil, err
	}
	return &moneyTx, nil
}

func (s *MoneyTxShippingEtopStore) GetMoneyTxShippingEtop() (*moneytx.MoneyTransactionShippingEtop, error) {
	moneyTx, err := s.GetMoneyTxShippingEtopDB()
	if err != nil {
		return nil, err
	}
	var res moneytx.MoneyTransactionShippingEtop
	if err := scheme.Convert(moneyTx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *MoneyTxShippingEtopStore) ListMoneyTxShippingEtopsDB() ([]*model.MoneyTransactionShippingEtop, error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortMoneyTx)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, filterMoneyTxShippingWhitelist)
	if err != nil {
		return nil, err
	}
	var moneyTxs model.MoneyTransactionShippingEtops
	if err := query.Find(&moneyTxs); err != nil {
		return nil, err
	}
	return moneyTxs, nil
}

func (s *MoneyTxShippingEtopStore) ListMoneyTxShippingEtops() (res []*moneytx.MoneyTransactionShippingEtop, _ error) {
	moneyTxs, err := s.ListMoneyTxShippingEtopsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(moneyTxs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingEtopStore) CreateMoneyTxShippingEtopDB(moneyTx *model.MoneyTransactionShippingEtop) error {
	sqlstore.MustNoPreds(s.preds)
	if moneyTx.ID == 0 {
		moneyTx.ID = cm.NewID()
	}
	return s.query().ShouldInsert(moneyTx)
}

func (s *MoneyTxShippingEtopStore) CreateMoneyTxShippingEtop(moneyTx *moneytx.MoneyTransactionShippingEtop) error {
	moneyTxDB := &model.MoneyTransactionShippingEtop{}
	if err := scheme.Convert(moneyTx, moneyTxDB); err != nil {
		return err
	}
	return s.CreateMoneyTxShippingEtopDB(moneyTxDB)
}

func (s *MoneyTxShippingEtopStore) UpdateMoneyTxShippingEtopDB(moneyTx *model.MoneyTransactionShippingEtop) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Where(s.preds).ShouldUpdate(moneyTx)
}

type UpdateMoneyTxShippingEtopStatisticsArgs struct {
	ID                    dot.ID
	TotalCOD              dot.NullInt
	TotalAmount           dot.NullInt
	TotalOrders           dot.NullInt
	TotalFee              dot.NullInt
	TotalMoneyTransaction dot.NullInt
}

func (s *MoneyTxShippingEtopStore) UpdateMoneyTxShippingEtopStatistics(args *UpdateMoneyTxShippingEtopStatisticsArgs) error {
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
	if args.TotalFee.Valid {
		update["total_fee"] = args.TotalFee.Int
	}
	if args.TotalMoneyTransaction.Valid {
		update["total_money_transaction"] = args.TotalMoneyTransaction.Int
	}

	if len(update) > 0 {
		return s.query().Table("money_transaction_shipping_etop").Where(s.ft.ByID(args.ID)).ShouldUpdateMap(update)
	}
	return nil
}

func (s *MoneyTxShippingEtopStore) DeleteMoneyTxShippingEtop(id dot.ID) error {
	query := s.query().Where(s.ft.ByID(id))
	return query.ShouldDelete((&model.MoneyTransactionShippingEtop{}))
}
