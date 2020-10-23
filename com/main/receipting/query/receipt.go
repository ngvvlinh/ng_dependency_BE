package query

import (
	"context"

	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/receipting/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ receipting.QueryService = &ReceiptQuery{}

type ReceiptQuery struct {
	store      sqlstore.ReceiptStoreFactory
	orderQuery ordering.QueryBus
}

func NewReceiptQuery(
	db com.MainDB,
	orderQuery ordering.QueryBus,
) *ReceiptQuery {
	return &ReceiptQuery{
		store:      sqlstore.NewReceiptStore(db),
		orderQuery: orderQuery,
	}
}

func ReceiptQueryMessageBus(q *ReceiptQuery) receipting.QueryBus {
	b := bus.New()
	return receipting.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *ReceiptQuery) GetReceipt(ctx context.Context, id, ShopID dot.ID) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ID(id).ShopID(ShopID).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) GetReceipts(ctx context.Context, shopID dot.ID) ([]*receipting.Receipt, error) {
	receipts, err := q.store(ctx).ShopID(shopID).ListReceipts()
	return receipts, err
}

func (q *ReceiptQuery) GetReceiptByID(ctx context.Context, args *receipting.GetReceiptByIDArg) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ID(args.ID).ShopID(args.ShopID).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) ListReceipts(
	ctx context.Context, args *receipting.ListReceiptsArgs,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Filters(args.Filters)
	receipts, err := query.WithPaging(args.Paging).ListReceipts()
	if err != nil {
		return nil, err
	}
	totalAmountConfirmedReceipt, totalAmountConfirmedPayment, err := query.Status(status3.P).SumAmountReceiptAndPayment()
	if err != nil {
		return nil, err
	}

	return &receipting.ReceiptsResponse{
		Receipts:                    receipts,
		TotalAmountConfirmedReceipt: totalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: totalAmountConfirmedPayment,
	}, nil
}

func (q *ReceiptQuery) ListReceiptsByIDs(context.Context, *receipting.GetReceiptbyIDsArgs) (*receipting.ReceiptsResponse, error) {
	panic("implement me")
}

func (q *ReceiptQuery) ListReceiptsByRefsAndStatus(
	ctx context.Context, args *receipting.ListReceiptsByRefsAndStatusArgs,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).
		RefIDs(args.IsContains, args.RefIDs...).RefType(args.RefType).
		Status(status3.Status(args.Status))
	receipts, err := query.ListReceipts()
	if err != nil {
		return nil, err
	}
	return &receipting.ReceiptsResponse{Receipts: receipts}, nil
}

func (q *ReceiptQuery) ListReceiptsByRefsAndStatusAndType(
	ctx context.Context, args *receipting.ListReceiptsByRefsAndStatusAndTypeArgs,
) ([]*receipting.Receipt, error) {
	query := q.store(ctx).ShopID(args.ShopID).
		RefIDs(args.IsContains, args.RefIDs...).RefType(args.RefType).
		Status(args.Status).ReceiptType(args.ReceiptType)
	receipts, err := query.ListReceipts()
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (q *ReceiptQuery) GetReceiptByCode(ctx context.Context, code string, shopID dot.ID) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ShopID(shopID).Code(code).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) ListReceiptsByTraderIDsAndStatuses(
	ctx context.Context, shopID dot.ID, traderIDs []dot.ID, statuses []status3.Status,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(shopID).TraderIDs(traderIDs...)
	if len(statuses) != 0 {
		query = query.Statuses(statuses...)
	}

	receipts, err := query.ListReceipts()
	return &receipting.ReceiptsResponse{Receipts: receipts}, err
}

func (q *ReceiptQuery) ListReceiptsByLedgerIDs(
	ctx context.Context, args *receipting.ListReceiptsByLedgerIDsArgs,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).LedgerIDs(args.LedgerIDs...).Filters(args.Filters)
	receipts, err := query.WithPaging(args.Paging).ListReceipts()
	if err != nil {
		return nil, err
	}
	totalAmountConfirmedReceipt, totalAmountConfirmedPayment, err := query.Status(status3.P).SumAmountReceiptAndPayment()
	if err != nil {
		return nil, err
	}

	return &receipting.ReceiptsResponse{
		Receipts:                    receipts,
		TotalAmountConfirmedReceipt: totalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: totalAmountConfirmedPayment,
	}, nil
}
