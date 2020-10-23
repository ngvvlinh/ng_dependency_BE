package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/receipting"
	"o.o/api/meta"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/receipting/convert"
	"o.o/backend/com/main/receipting/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

type ReceiptStoreFactory func(ctx context.Context) *ReceiptStore

func NewReceiptStore(db *cmsql.Database) ReceiptStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ReceiptStore {
		return &ReceiptStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ReceiptStore struct {
	ft ReceiptFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ReceiptStore) WithPaging(paging meta.Paging) *ReceiptStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *ReceiptStore) Filters(filters meta.Filters) *ReceiptStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ReceiptStore) ID(id dot.ID) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ReceiptStore) IDs(ids ...dot.ID) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *ReceiptStore) ShopID(id dot.ID) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *ReceiptStore) CreatedAtFromAndTo(from, to time.Time) *ReceiptStore {
	s.preds = append(s.preds, sq.NewExpr("created_at >= ? AND created_at < ?", from, to))
	return s
}

func (s *ReceiptStore) Code(code string) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "code", code))
	return s
}

func (s *ReceiptStore) TraderID(traderID dot.ID) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByTraderID(traderID))
	return s
}

func (s *ReceiptStore) TraderIDs(traderIDs ...dot.ID) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "trader_id", traderIDs))
	return s
}

func (s *ReceiptStore) RefsID(id dot.ID) *ReceiptStore {
	s.preds = append(s.preds, sq.NewExpr("ref_ids @> ?", core.Array{V: []dot.ID{id}}))
	return s
}

func (s *ReceiptStore) RefIDs(isContains bool, ids ...dot.ID) *ReceiptStore {
	if len(ids) == 0 {
		s.preds = append(s.preds, sq.NewExpr("false"))
		return s
	}

	// default operator is overlap
	operator := "&&"
	if isContains {
		operator = "@>"
	}
	s.preds = append(s.preds, sq.NewExpr("ref_ids "+operator+" ?", core.Array{V: ids}))
	return s
}

func (s *ReceiptStore) ReceiptType(typ receipt_type.ReceiptType) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByType(typ))
	return s
}

func (s *ReceiptStore) RefType(refType receipt_ref.ReceiptRef) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByRefType(refType))
	return s
}

func (s *ReceiptStore) Status(status status3.Status) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ReceiptStore) Statuses(statuses ...status3.Status) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "status", statuses))
	return s
}

func (s *ReceiptStore) LedgerIDs(LedgerIDs ...dot.ID) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "ledger_id", LedgerIDs))
	return s
}

func (s *ReceiptStore) InPerDate(date time.Time) *ReceiptStore {
	s.preds = append(s.preds, sq.NewExpr("created_at BETWEEN ? AND ?", date, date.Add(24*time.Hour)))
	return s
}

func (s *ReceiptStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ReceiptStore) ConfirmReceipt() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"status":       int(status3.P),
		"confirmed_at": time.Now(),
	})

	return _updated, err
}

func (s *ReceiptStore) CancelReceipt(reason string) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"status":           int(status3.N),
		"cancelled_reason": reason,
		"cancelled_at":     time.Now(),
	})
	return _updated, err
}

func (s *ReceiptStore) GetReceiptDB() (*model.Receipt, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var receipt model.Receipt
	err := query.ShouldGet(&receipt)
	return &receipt, err
}

func (s *ReceiptStore) GetReceipt() (receiptResult *receipting.Receipt, _ error) {
	receipt, err := s.GetReceiptDB()
	if err != nil {
		return nil, err
	}
	receiptResult = convert.Convert_receiptingmodel_Receipt_receipting_Receipt(receipt, receiptResult)
	return receiptResult, nil
}

func (s *ReceiptStore) CreateReceipt(receipt *receipting.Receipt) error {
	sqlstore.MustNoPreds(s.preds)
	receiptDB := new(model.Receipt)
	if err := scheme.Convert(receipt, receiptDB); err != nil {
		return err
	}
	receiptDB.Lines = convert.Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(receipt.Lines)
	if receiptDB.Trader != nil {
		receiptDB.TraderFullNameNorm = validate.NormalizeSearch(receiptDB.Trader.FullName)
		receiptDB.TraderType = receiptDB.Trader.Type
		receiptDB.TraderPhoneNorm = validate.NormalizeSearchPhone(receiptDB.Trader.Phone)
	}
	if _, err := s.query().Insert(receiptDB); err != nil {
		return err
	}

	var tempReceipt model.Receipt
	if err := s.query().Where(s.ft.ByID(receipt.ID), s.ft.ByShopID(receipt.ShopID)).ShouldGet(&tempReceipt); err != nil {
		return err
	}

	receipt.CreatedAt = tempReceipt.CreatedAt
	receipt.UpdatedAt = tempReceipt.UpdatedAt
	return nil
}

func (s *ReceiptStore) UpdateReceiptDB(receipt *model.Receipt) error {
	sqlstore.MustNoPreds(s.preds)
	if receipt.Trader != nil {
		receipt.TraderFullNameNorm = validate.NormalizeSearch(receipt.Trader.FullName)
		receipt.TraderType = receipt.Trader.Type
		receipt.TraderPhoneNorm = validate.NormalizeSearchPhone(receipt.Trader.Phone)
	}
	err := s.query().Where(
		s.ft.ByID(receipt.ID),
		s.ft.ByShopID(receipt.ShopID),
	).UpdateAll().ShouldUpdate(receipt)
	return err
}

func (s *ReceiptStore) GetReceiptByMaximumCodeNorm() (*model.Receipt, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = query.OrderBy("code_norm desc").Limit(1)

	var receipt model.Receipt
	if err := query.ShouldGet(&receipt); err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (s *ReceiptStore) ListReceiptsDB() ([]*model.Receipt, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortReceipt)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterReceipt)
	if err != nil {
		return nil, err
	}

	var receipts model.Receipts
	err = query.Find(&receipts)
	return receipts, err
}

func (s *ReceiptStore) ListReceipts() (receiptsResult []*receipting.Receipt, _ error) {
	receipts, err := s.ListReceiptsDB()
	if err != nil {
		return nil, err
	}

	receiptsResult = convert.Convert_receiptingmodel_Receipts_receipting_Receipts(receipts)
	for i, receipt := range receiptsResult {
		receipt.Lines = convert.Convert_receiptingmodel_ReceiptLines_receipting_ReceiptLines(receipts[i].Lines)
	}

	return receiptsResult, nil
}

func (s *ReceiptStore) SumAmountReceiptAndPayment() (receipt, payment int, err error) {
	var sqlReceipt, sqlPayment core.Int
	query := s.query().Where(s.preds)
	err = query.Table("receipt").Select(
		"SUM(amount) FILTER(WHERE type = 'receipt')",
		"SUM(amount) FILTER(WHERE type = 'payment')").
		Scan(&sqlReceipt, &sqlPayment)

	return int(sqlReceipt), int(sqlPayment), err
}
