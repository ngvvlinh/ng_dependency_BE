package sqlstore

import (
	"context"

	"o.o/api/main/refund"
	"o.o/api/meta"
	"o.o/backend/com/main/refund/convert"
	"o.o/backend/com/main/refund/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type RefundStoreFactory func(ctx context.Context) *RefundStore

func NewRefundStore(db *cmsql.Database) RefundStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *RefundStore {
		return &RefundStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type RefundStore struct {
	ft RefundFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *RefundStore) WithPaging(paging meta.Paging) *RefundStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *RefundStore) Filters(filters meta.Filters) *RefundStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *RefundStore) ID(id dot.ID) *RefundStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *RefundStore) IDs(ids ...dot.ID) *RefundStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *RefundStore) ShopID(id dot.ID) *RefundStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *RefundStore) Code(code string) *RefundStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "code", code))
	return s
}

func (s *RefundStore) OrderID(orderID dot.ID) *RefundStore {
	s.preds = append(s.preds, s.ft.ByOrderID(orderID))
	return s
}

func (s *RefundStore) GetRefundDB() (*model.Refund, error) {
	query := s.query().Where(s.preds)
	var refund model.Refund
	err := query.ShouldGet(&refund)
	return &refund, err
}

func (s *RefundStore) GetRefund() (refundResult *refund.Refund, err error) {
	refund, err := s.GetRefundDB()
	if err != nil {
		return nil, err
	}
	refundResult = convert.Convert_refundmodel_Refund_refund_Refund(refund, refundResult)
	return refundResult, nil
}

func (s *RefundStore) ListRefundsDB() ([]*model.Refund, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortRefund)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterRefund)
	if err != nil {
		return nil, err
	}

	var receipts model.Refunds
	err = query.Find(&receipts)
	return receipts, err
}

func (s *RefundStore) ListRefunds() ([]*refund.Refund, error) {
	refundsDB, err := s.ListRefundsDB()
	if err != nil {
		return nil, err
	}
	refunds := convert.Convert_refundmodel_Refunds_refund_Refunds(refundsDB)
	return refunds, nil
}

func (s *RefundStore) UpdateRefundDB(args *model.Refund) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *RefundStore) UpdateRefundAll(args *refund.Refund) error {
	var result = &model.Refund{}
	result = convert.Convert_refund_Refund_refundmodel_Refund(args, result)
	return s.UpdateRefundAllDB(result)
}

func (s *RefundStore) UpdateRefundAllDB(args *model.Refund) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *RefundStore) Create(args *refund.Refund) error {
	var voucherDB = model.Refund{}
	convert.Convert_refund_Refund_refundmodel_Refund(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *RefundStore) CreateDB(Refund *model.Refund) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(Refund)
}

func (s *RefundStore) GetRefundByMaximumCodeNorm() (*model.Refund, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var inventoryVoucher model.Refund
	if err := query.ShouldGet(&inventoryVoucher); err != nil {
		return nil, err
	}
	return &inventoryVoucher, nil
}
