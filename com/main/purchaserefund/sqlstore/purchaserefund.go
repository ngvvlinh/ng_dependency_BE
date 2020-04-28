package sqlstore

import (
	"context"

	"o.o/api/main/purchaserefund"
	"o.o/api/meta"
	"o.o/backend/com/main/purchaserefund/convert"
	"o.o/backend/com/main/purchaserefund/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type PurchaseRefundStoreFactory func(ctx context.Context) *PurchaseRefundStore

func NewPurchaseRefundStore(db *cmsql.Database) PurchaseRefundStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *PurchaseRefundStore {
		return &PurchaseRefundStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PurchaseRefundStore struct {
	ft PurchaseRefundFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	Paging  sqlstore.Paging
}

func (s *PurchaseRefundStore) WithPaging(paging meta.Paging) *PurchaseRefundStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *PurchaseRefundStore) Filters(filters meta.Filters) *PurchaseRefundStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *PurchaseRefundStore) ID(id dot.ID) *PurchaseRefundStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PurchaseRefundStore) IDs(ids ...dot.ID) *PurchaseRefundStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *PurchaseRefundStore) ShopID(id dot.ID) *PurchaseRefundStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *PurchaseRefundStore) Code(code string) *PurchaseRefundStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "code", code))
	return s
}

func (s *PurchaseRefundStore) PurchaseOrderID(purchaseOrderID dot.ID) *PurchaseRefundStore {
	s.preds = append(s.preds, s.ft.ByPurchaseOrderID(purchaseOrderID))
	return s
}

func (s *PurchaseRefundStore) GetPurchaseRefundDB() (*model.PurchaseRefund, error) {
	query := s.query().Where(s.preds)
	var purchaseRefund model.PurchaseRefund
	err := query.ShouldGet(&purchaseRefund)
	return &purchaseRefund, err
}

func (s *PurchaseRefundStore) GetPurchaseRefund() (purchaserefundResult *purchaserefund.PurchaseRefund, err error) {
	purchaseRefund, err := s.GetPurchaseRefundDB()
	if err != nil {
		return nil, err
	}
	purchaserefundResult = convert.Convert_purchaserefundmodel_PurchaseRefund_purchaserefund_PurchaseRefund(purchaseRefund, purchaserefundResult)
	return purchaserefundResult, nil
}

func (s *PurchaseRefundStore) ListPurchaseRefundsDB() ([]*model.PurchaseRefund, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if s.Paging.Sort == nil || len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortPurchaseRefund)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterPurchaseRefund)
	if err != nil {
		return nil, err
	}

	var receipts model.PurchaseRefunds
	err = query.Find(&receipts)
	return receipts, err
}

func (s *PurchaseRefundStore) ListPurchaseRefunds() ([]*purchaserefund.PurchaseRefund, error) {
	purchaserefundsDB, err := s.ListPurchaseRefundsDB()
	if err != nil {
		return nil, err
	}
	purchaserefunds := convert.Convert_purchaserefundmodel_PurchaseRefunds_purchaserefund_PurchaseRefunds(purchaserefundsDB)
	return purchaserefunds, nil
}

func (s *PurchaseRefundStore) UpdatePurchaseRefundDB(args *model.PurchaseRefund) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *PurchaseRefundStore) UpdatePurchaseRefundAll(args *purchaserefund.PurchaseRefund) error {
	var result = &model.PurchaseRefund{}
	result = convert.Convert_purchaserefund_PurchaseRefund_purchaserefundmodel_PurchaseRefund(args, result)
	return s.UpdatePurchaseRefundAllDB(result)
}

func (s *PurchaseRefundStore) UpdatePurchaseRefundAllDB(args *model.PurchaseRefund) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *PurchaseRefundStore) Create(args *purchaserefund.PurchaseRefund) error {
	var voucherDB = model.PurchaseRefund{}
	convert.Convert_purchaserefund_PurchaseRefund_purchaserefundmodel_PurchaseRefund(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *PurchaseRefundStore) CreateDB(PurchaseRefund *model.PurchaseRefund) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(PurchaseRefund)
}

func (s *PurchaseRefundStore) GetPurchaseRefundByMaximumCodeNorm() (*model.PurchaseRefund, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var inventoryVoucher model.PurchaseRefund
	if err := query.ShouldGet(&inventoryVoucher); err != nil {
		return nil, err
	}
	return &inventoryVoucher, nil
}

func (s *PurchaseRefundStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.PurchaseRefund)(nil))
}
