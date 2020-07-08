package sqlstore

import (
	"context"

	"o.o/api/main/moneytx"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	identityconvert "o.o/backend/com/main/identity/convert"
	"o.o/backend/com/main/moneytx/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type MoneyTxShippingExternalStore struct {
	ft     MoneyTransactionShippingExternalFilters
	ftLine MoneyTransactionShippingExternalLineFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	ctx context.Context
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
			ctx: ctx,
		}
	}
}

func (s *MoneyTxShippingExternalStore) WithPaging(paging meta.Paging) *MoneyTxShippingExternalStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *MoneyTxShippingExternalStore) ID(id dot.ID) *MoneyTxShippingExternalStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *MoneyTxShippingExternalStore) IDs(ids ...dot.ID) *MoneyTxShippingExternalStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *MoneyTxShippingExternalStore) Filters(filters meta.Filters) *MoneyTxShippingExternalStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalFtLineDB() (*model.MoneyTransactionShippingExternalFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)

	moneyTx := &model.MoneyTransactionShippingExternal{}
	if err := query.ShouldGet(moneyTx); err != nil {
		return nil, err
	}

	lines, err := s.ListMoneyTxShippingExternalLinesDBByMoneyTxID([]dot.ID{moneyTx.ID})
	res := &model.MoneyTransactionShippingExternalFtLine{
		MoneyTransactionShippingExternal: moneyTx,
		Lines:                            lines,
	}
	return res, err
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalDB() (*model.MoneyTransactionShippingExternal, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)

	moneyTx := &model.MoneyTransactionShippingExternal{}
	err := query.ShouldGet(moneyTx)
	return moneyTx, err
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternalFtLine() (*moneytx.MoneyTransactionShippingExternalFtLine, error) {
	moneyTx, err := s.GetMoneyTxShippingExternalFtLineDB()
	if err != nil {
		return nil, err
	}
	var res moneytx.MoneyTransactionShippingExternalFtLine
	if err := scheme.Convert(moneyTx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *MoneyTxShippingExternalStore) GetMoneyTxShippingExternal() (*moneytx.MoneyTransactionShippingExternal, error) {
	moneyTx, err := s.GetMoneyTxShippingExternalDB()
	if err != nil {
		return nil, err
	}
	var res moneytx.MoneyTransactionShippingExternal
	if err := scheme.Convert(moneyTx, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternalsDB() ([]*model.MoneyTransactionShippingExternal, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortMoneyTx)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, filterMoneyTxShippingExternalWhitelist)
	if err != nil {
		return nil, err
	}
	query = s.ByWhiteLabelPartner(s.ctx, query)

	var moneyTxs model.MoneyTransactionShippingExternals
	if err := query.Find(&moneyTxs); err != nil {
		return nil, err
	}
	return moneyTxs, nil
}

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternals() (res []*moneytx.MoneyTransactionShippingExternal, _ error) {
	moneyTxs, err := s.ListMoneyTxShippingExternalsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(moneyTxs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternalsFtLineDB() ([]*model.MoneyTransactionShippingExternalFtLine, error) {
	moneyTxs, err := s.ListMoneyTxShippingExternalsDB()
	if err != nil {
		return nil, err
	}
	moneyTxIDs := make([]dot.ID, len(moneyTxs))
	for i, m := range moneyTxs {
		moneyTxIDs[i] = m.ID
	}
	lines, err := s.ListMoneyTxShippingExternalLinesDBByMoneyTxID(moneyTxIDs)
	if err != nil {
		return nil, err
	}

	linesMoneyTxHash := make(map[dot.ID][]*model.MoneyTransactionShippingExternalLine)
	for _, line := range lines {
		linesMoneyTxHash[line.MoneyTransactionShippingExternalID] = append(linesMoneyTxHash[line.MoneyTransactionShippingExternalID], line)
	}

	result := make([]*model.MoneyTransactionShippingExternalFtLine, len(moneyTxs))
	for i, moneyTx := range moneyTxs {
		result[i] = &model.MoneyTransactionShippingExternalFtLine{
			MoneyTransactionShippingExternal: moneyTx,
			Lines:                            linesMoneyTxHash[moneyTx.ID],
		}
	}

	return result, nil
}

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternalsFtLine() (res []*moneytx.MoneyTransactionShippingExternalFtLine, _ error) {
	moneyTxsDB, err := s.ListMoneyTxShippingExternalsFtLineDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(moneyTxsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingExternalStore) CreateMoneyTxShippingExternal(moneyTx *moneytx.MoneyTransactionShippingExternal) error {
	sqlstore.MustNoPreds(s.preds)
	moneyTxDB := &model.MoneyTransactionShippingExternal{}
	if err := scheme.Convert(moneyTx, moneyTxDB); err != nil {
		return err
	}
	moneyTxDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)

	if err := s.query().ShouldInsert(moneyTxDB); err != nil {
		return err
	}
	return nil
}

func (s *MoneyTxShippingExternalStore) UpdateMoneyTxShippingExternalInfo(args *moneytx.UpdateMoneyTxShippingExternalInfoArgs) error {
	sqlstore.MustNoPreds(s.preds)
	bankAccount := identityconvert.Convert_identitytypes_BankAccount_sharemodel_BankAccount(args.BankAccount, nil)
	update := &model.MoneyTransactionShippingExternal{
		Note:          args.Note,
		InvoiceNumber: args.InvoiceNumber,
		BankAccount:   bankAccount,
	}

	query := s.query().Where(s.ft.ByID(args.MoneyTxShippingExternalID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(update)
}

type UpdateMoneyTxShippingExternalStatisticsArgs struct {
	MoneyTxShippingExternalID dot.ID
	TotalCOD                  dot.NullInt
	TotalOrders               dot.NullInt
}

func (s *MoneyTxShippingExternalStore) UpdateMoneyTxShippingExternalStatistics(args *UpdateMoneyTxShippingExternalStatisticsArgs) error {
	var update = make(map[string]interface{})
	if args.TotalCOD.Valid {
		update["total_cod"] = args.TotalCOD.Int
	}
	if args.TotalOrders.Valid {
		update["total_orders"] = args.TotalOrders.Int
	}
	if args.TotalCOD.Valid && args.TotalCOD.Int == 0 &&
		args.TotalOrders.Valid && args.TotalOrders.Int == 0 {
		// money tx does not have any ffm => update status = -1
		update["status"] = status3.N
	}
	query := s.query().Table("money_transaction_shipping_external").Where(s.ft.ByID(args.MoneyTxShippingExternalID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdateMap(update)
}

func (s *MoneyTxShippingExternalStore) ConfirmMoneyTxShippingExternals(ids []dot.ID) error {
	query := s.query().Table("money_transaction_shipping_external").Where(sq.In("id", ids))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdateMap(map[string]interface{}{
		"status": status3.P,
	})
}

func (s *MoneyTxShippingExternalStore) DeleteMoneyTxShippingExternal(id dot.ID) error {
	query := s.query().Where(s.ft.ByID(id))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldDelete((&model.MoneyTransactionShippingExternal{}))
}

// -- Money Transaction Shipping External Line -- //

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

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternalLinesDBByMoneyTxID(moneyTxShippingExternalIDs []dot.ID) (lines []*model.MoneyTransactionShippingExternalLine, _ error) {
	query := s.query().Where(sq.PrefixedIn(&s.ftLine.prefix, "money_transaction_shipping_external_id", moneyTxShippingExternalIDs))
	if err := query.Find((*model.MoneyTransactionShippingExternalLines)(&lines)); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingExternalStore) Line_by_MoneyTxShippingExternalID(id dot.ID) *MoneyTxShippingExternalStore {
	s.preds = append(s.preds, s.ftLine.ByMoneyTransactionShippingExternalID(id))
	return s
}

func (s *MoneyTxShippingExternalStore) Line_by_LineIDs(ids ...dot.ID) *MoneyTxShippingExternalStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ftLine.prefix, "id", ids))
	return s
}

func (s *MoneyTxShippingExternalStore) ListMoneyTxShippingExternalLinesDB() (lines []*model.MoneyTransactionShippingExternalLine, _ error) {
	query := s.query().Where(s.preds)
	if err := query.Find((*model.MoneyTransactionShippingExternalLines)(&lines)); err != nil {
		return nil, err
	}
	return
}

func (s *MoneyTxShippingExternalStore) DeleteMoneyTxShippingExternalLines() error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	query := s.query().Where(s.preds)
	return query.ShouldDelete(&model.MoneyTransactionShippingExternalLine{})
}

func (s *MoneyTxShippingExternalStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
