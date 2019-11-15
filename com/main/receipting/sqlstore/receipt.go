package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/backend/pkg/common/validate"

	"etop.vn/api/main/receipting"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/receipting/convert"
	"etop.vn/backend/com/main/receipting/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	etopmodel "etop.vn/backend/pkg/etop/model"
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
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ReceiptStore) Paging(paging meta.Paging) *ReceiptStore {
	ss := *s
	ss.paging = paging
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

func (s *ReceiptStore) ID(id int64) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ReceiptStore) IDs(ids ...int64) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *ReceiptStore) ShopID(id int64) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *ReceiptStore) Code(code string) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "code", code))
	return s
}

func (s *ReceiptStore) TraderID(traderID int64) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByTraderID(traderID))
	return s
}

func (s *ReceiptStore) TraderIDs(traderIDs ...int64) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "trader_id", traderIDs))
	return s
}

func (s *ReceiptStore) RefsID(id int64) *ReceiptStore {
	s.preds = append(s.preds, sq.NewExpr("ref_ids @> '{?}'", strconv.FormatInt(id, 10)))
	return s
}

func (s *ReceiptStore) RefIDs(ids ...int64) *ReceiptStore {
	if len(ids) == 0 {
		s.preds = append(s.preds, sq.NewExpr("false"))
		return s
	}

	strConditions := "ref_ids && '{"
	for i, id := range ids {
		strConditions += fmt.Sprintf("%d", id)
		if i < len(ids)-1 {
			strConditions += ","
		}
	}
	strConditions += "}'"

	s.preds = append(s.preds, sq.NewExpr(strConditions))
	return s
}

func (s *ReceiptStore) RefType(refType receipting.ReceiptRefType) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByRefType(string(refType)))
	return s
}

func (s *ReceiptStore) Status(status etopmodel.Status3) *ReceiptStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ReceiptStore) Statuses(statuses ...etop.Status3) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "status", statuses))
	return s
}

func (s *ReceiptStore) LedgerIDs(LedgerIDs ...int64) *ReceiptStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "ledger_id", LedgerIDs))
	return s
}

func (s *ReceiptStore) Count() (_ uint64, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	query, _, err = sqlstore.Filters(query, s.filters, FilterReceipt)
	if err != nil {
		return 0, err
	}

	return query.Count((*model.Receipt)(nil))
}

func (s *ReceiptStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *ReceiptStore) ConfirmReceipt() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"status":       int32(etopmodel.S3Positive),
		"confirmed_at": time.Now(),
	})

	return int(_updated), err
}

func (s *ReceiptStore) CancelReceipt(reason string) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("receipt").UpdateMap(map[string]interface{}{
		"status":           int32(etopmodel.S3Negative),
		"cancelled_reason": reason,
		"cancelled_at":     time.Now(),
	})

	return int(_updated), err
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
	if s.paging.Sort == nil || len(s.paging.Sort) == 0 {
		s.paging.Sort = append(s.paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.paging, SortReceipt)
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

func (s *ReceiptStore) SumAmountReceiptAndPayment() (receipt, payment int64, err error) {
	var sqlReceipt, sqlPayment sql.NullInt64
	query := s.query().Where(s.preds)
	err = query.Table("receipt").Select(
		"SUM(amount) FILTER(WHERE type = 'receipt')",
		"SUM(amount) FILTER(WHERE type = 'payment')").
		Scan(&sqlReceipt, &sqlPayment)

	return sqlReceipt.Int64, sqlPayment.Int64, err
}
