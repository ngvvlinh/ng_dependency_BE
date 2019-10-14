package sqlstore

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"etop.vn/backend/pkg/common/scheme"

	"etop.vn/backend/pkg/common/sq"

	"etop.vn/backend/com/main/receipting/convert"

	"etop.vn/backend/com/main/receipting/model"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/api/main/receipting"
	"etop.vn/api/meta"

	"etop.vn/backend/pkg/common/cmsql"
)

type ReceiptStoreFactory func(ctx context.Context) *ReceiptStore

func NewReceiptStore(db cmsql.Database) ReceiptStoreFactory {
	return func(ctx context.Context) *ReceiptStore {
		return &ReceiptStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ReceiptStore struct {
	ft ReceiptFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ReceiptStore) Paging(paging meta.Paging) *ReceiptStore {
	s.paging = paging
	return s
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

func (s *ReceiptStore) OrderID(id int64) *ReceiptStore {
	s.preds = append(s.preds, sq.NewExpr("order_ids @> '{?}'", strconv.FormatInt(id, 10)))
	return s
}

func (s *ReceiptStore) OrderIDs(ids ...int64) *ReceiptStore {
	if len(ids) == 0 {
		s.preds = append(s.preds, sq.NewExpr("false"))
		return s
	}

	strConditions := "order_ids && '{"
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

func (s *ReceiptStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
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
	err := s.query().Where(
		s.ft.ByID(receipt.ID),
		s.ft.ByShopID(receipt.ShopID),
	).UpdateAll().ShouldUpdate(receipt)
	return err
}

func (s *ReceiptStore) ListReceiptsDB() ([]*model.Receipt, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
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
