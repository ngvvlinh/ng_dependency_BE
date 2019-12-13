package query

import (
	"context"

	"etop.vn/api/main/refund"
	"etop.vn/backend/com/main/refund/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi"
)

var _ refund.QueryService = &RefundQueryService{}

type RefundQueryService struct {
	RefundStore sqlstore.RefundStoreFactory
	EventBus    capi.EventBus
}

func NewQueryRefund(eventBus capi.EventBus, db *cmsql.Database) *RefundQueryService {
	return &RefundQueryService{
		RefundStore: sqlstore.NewRefundStore(db),
		EventBus:    eventBus,
	}
}

func (q *RefundQueryService) MessageBus() refund.QueryBus {
	b := bus.New()
	return refund.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *RefundQueryService) GetRefunds(ctx context.Context, args *refund.GetRefundsArgs) (*refund.GetRefundsResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	query := q.RefundStore(ctx).ShopID(args.ShopID).Filters(args.Filters).Paging(args.Paging)
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	result, err := q.RefundStore(ctx).ShopID(args.ShopID).Filters(args.Filters).Paging(args.Paging).ListRefunds()
	if err != nil {
		return nil, err
	}
	return &refund.GetRefundsResponse{
		PageInfor: query.GetPaging(),
		Refunds:   result,
		Count:     count,
	}, nil
}

func (q *RefundQueryService) GetRefundByID(ctx context.Context, args *refund.GetRefundByIDArgs) (*refund.Refund, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	return q.RefundStore(ctx).ShopID(args.ShopID).ID(args.ID).GetRefund()
}

func (q *RefundQueryService) GetRefundsByIDs(ctx context.Context, args *refund.GetRefundsByIDsArgs) ([]*refund.Refund, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	return q.RefundStore(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListRefunds()
}
