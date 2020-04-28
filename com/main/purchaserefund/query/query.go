package query

import (
	"context"

	"o.o/api/main/purchaserefund"
	"o.o/backend/com/main/purchaserefund/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
)

var _ purchaserefund.QueryService = &PurchaseRefundQueryService{}

type PurchaseRefundQueryService struct {
	PurchaseRefundStore sqlstore.PurchaseRefundStoreFactory
	EventBus            capi.EventBus
}

func NewQueryPurchasePurchaseRefund(eventBus capi.EventBus, db *cmsql.Database) *PurchaseRefundQueryService {
	return &PurchaseRefundQueryService{
		PurchaseRefundStore: sqlstore.NewPurchaseRefundStore(db),
		EventBus:            eventBus,
	}
}

func (q *PurchaseRefundQueryService) MessageBus() purchaserefund.QueryBus {
	b := bus.New()
	return purchaserefund.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *PurchaseRefundQueryService) ListPurchaseRefunds(ctx context.Context, args *purchaserefund.GetPurchaseRefundsArgs) (*purchaserefund.GetPurchaseRefundsResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	query := q.PurchaseRefundStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	result, err := q.PurchaseRefundStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging).ListPurchaseRefunds()
	if err != nil {
		return nil, err
	}
	return &purchaserefund.GetPurchaseRefundsResponse{
		PageInfo:        query.Paging.GetPaging(),
		PurchaseRefunds: result,
		Count:           count,
	}, nil
}

func (q *PurchaseRefundQueryService) GetPurchaseRefundByID(ctx context.Context, args *purchaserefund.GetPurchaseRefundByIDArgs) (*purchaserefund.PurchaseRefund, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	return q.PurchaseRefundStore(ctx).ShopID(args.ShopID).ID(args.ID).GetPurchaseRefund()
}

func (q *PurchaseRefundQueryService) GetPurchaseRefundsByIDs(ctx context.Context, args *purchaserefund.GetPurchaseRefundsByIDsArgs) ([]*purchaserefund.PurchaseRefund, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	return q.PurchaseRefundStore(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListPurchaseRefunds()
}

func (q *PurchaseRefundQueryService) GetPurchaseRefundsByPurchaseOrderID(ctx context.Context, args *purchaserefund.GetPurchaseRefundsByPurchaseOrderIDRequest) ([]*purchaserefund.PurchaseRefund, error) {
	var result []*purchaserefund.PurchaseRefund
	if args.ShopID == 0 || args.PurchaseOrderID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.PurchaseRefundStore(ctx).ShopID(args.ShopID).PurchaseOrderID(args.PurchaseOrderID).ListPurchaseRefunds()
	return result, err
}
