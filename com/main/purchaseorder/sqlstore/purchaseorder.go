package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/main/purchaseorder/convert"
	"etop.vn/backend/com/main/purchaseorder/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type PurchaseOrderStoreFactory func(ctx context.Context) *PurchaseOrderStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewPurchaseOrderStore(db *cmsql.Database) PurchaseOrderStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *PurchaseOrderStore {
		return &PurchaseOrderStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PurchaseOrderStore struct {
	ft PurchaseOrderFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *PurchaseOrderStore) Paging(paging meta.Paging) *PurchaseOrderStore {
	s.paging = paging
	return s
}

func (s *PurchaseOrderStore) GetPaing() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *PurchaseOrderStore) Filters(filters meta.Filters) *PurchaseOrderStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *PurchaseOrderStore) ID(id dot.ID) *PurchaseOrderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PurchaseOrderStore) IDs(ids ...dot.ID) *PurchaseOrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *PurchaseOrderStore) ShopID(id dot.ID) *PurchaseOrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *PurchaseOrderStore) Status(status status3.Status) *PurchaseOrderStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *PurchaseOrderStore) Statuses(statuses ...status3.Status) *PurchaseOrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "status", statuses))
	return s
}

func (s *PurchaseOrderStore) SupplierIDs(supplierIDs ...dot.ID) *PurchaseOrderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "supplier_id", supplierIDs))
	return s
}

func (s *PurchaseOrderStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	query, _, err := sqlstore.Filters(query, s.filters, FilterPurchaseOrder)
	if err != nil {
		return 0, err
	}

	return query.Count((*model.PurchaseOrder)(nil))
}

func (s *PurchaseOrderStore) GetReceiptByMaximumCodeNorm() (*model.PurchaseOrder, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = query.OrderBy("code_norm desc").Limit(1)

	var purchaseOrder model.PurchaseOrder
	if err := query.ShouldGet(&purchaseOrder); err != nil {
		return nil, err
	}
	return &purchaseOrder, nil
}

func (s *PurchaseOrderStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("purchase_order").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *PurchaseOrderStore) ConfirmPurchaseOrder() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("purchase_order").UpdateMap(map[string]interface{}{
		"status":       int(status5.P),
		"confirmed_at": time.Now(),
	})

	return _updated, err
}

func (s *PurchaseOrderStore) CancelPurchaseOrder(reason string) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("purchase_order").UpdateMap(map[string]interface{}{
		"status":           int(status3.N),
		"cancelled_reason": reason,
		"cancelled_at":     time.Now(),
	})
	return _updated, err
}

func (s *PurchaseOrderStore) GetPurchaseOrderDB() (*model.PurchaseOrder, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var purchaseOrder model.PurchaseOrder
	err := query.ShouldGet(&purchaseOrder)
	return &purchaseOrder, err
}

func (s *PurchaseOrderStore) GetPurchaseOrder() (purchaseOrderResult *purchaseorder.PurchaseOrder, _ error) {
	purchaseOrder, err := s.GetPurchaseOrderDB()
	if err != nil {
		return nil, err
	}
	purchaseOrderResult = convert.Convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(purchaseOrder, purchaseOrderResult)
	return purchaseOrderResult, nil
}

func (s *PurchaseOrderStore) CreatePurchaseOrder(purchaseOrder *purchaseorder.PurchaseOrder) error {
	sqlstore.MustNoPreds(s.preds)
	purchaseOrderDB := new(model.PurchaseOrder)
	if err := scheme.Convert(purchaseOrder, purchaseOrderDB); err != nil {
		return err
	}

	if _, err := s.query().Insert(purchaseOrderDB); err != nil {
		return err
	}

	var tempPurchaseOrder model.PurchaseOrder
	if err := s.query().
		Where(s.ft.ByID(purchaseOrder.ID), s.ft.ByShopID(purchaseOrder.ShopID)).
		ShouldGet(&tempPurchaseOrder); err != nil {
		return err
	}

	purchaseOrder.CreatedAt = tempPurchaseOrder.CreatedAt
	purchaseOrder.UpdatedAt = tempPurchaseOrder.UpdatedAt
	return nil
}

func (s *PurchaseOrderStore) UpdatePurchaseOrderDB(purchaseOrder *model.PurchaseOrder) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(
		s.ft.ByID(purchaseOrder.ID),
		s.ft.ByShopID(purchaseOrder.ShopID),
	).UpdateAll().ShouldUpdate(purchaseOrder)
	return err
}

func (s *PurchaseOrderStore) ListPurchaseOrderDB() ([]*model.PurchaseOrder, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	// default sort by created_at
	if s.paging.Sort == nil || len(s.paging.Sort) == 0 {
		s.paging.Sort = append(s.paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.paging, SortPurchaseOrder)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterPurchaseOrder)
	if err != nil {
		return nil, err
	}

	var purchaseOrders model.PurchaseOrders
	err = query.Find(&purchaseOrders)
	return purchaseOrders, err
}

func (s *PurchaseOrderStore) ListPurchaseOrders() (purchaseOrdersResult []*purchaseorder.PurchaseOrder, _ error) {
	purchaseOrders, err := s.ListPurchaseOrderDB()
	if err != nil {
		return nil, err
	}

	purchaseOrdersResult = convert.Convert_purchaseordermodel_PurchaseOrders_purchaseorder_PurchaseOrders(purchaseOrders)
	return purchaseOrdersResult, nil
}

func (s *PurchaseOrderStore) IncludeDeleted() *PurchaseOrderStore {
	s.includeDeleted = true
	return s
}
