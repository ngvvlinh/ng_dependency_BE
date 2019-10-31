package query

import (
	"context"

	"etop.vn/backend/pkg/etop/model"

	"etop.vn/api/shopping"

	"etop.vn/backend/pkg/common/bus"

	"etop.vn/api/main/receipting"
	"etop.vn/backend/com/main/receipting/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ receipting.QueryService = &ReceiptQuery{}

type ReceiptQuery struct {
	store sqlstore.ReceiptStoreFactory
}

func NewReceiptQuery(db *cmsql.Database) *ReceiptQuery {
	return &ReceiptQuery{
		store: sqlstore.NewReceiptStore(db),
	}
}

func (q *ReceiptQuery) MessageBus() receipting.QueryBus {
	b := bus.New()
	return receipting.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *ReceiptQuery) GetReceipt(ctx context.Context, id, ShopID int64) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ID(id).ShopID(ShopID).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) GetReceipts(ctx context.Context, shopID int64) ([]*receipting.Receipt, error) {
	receipts, err := q.store(ctx).ShopID(shopID).ListReceipts()
	return receipts, err
}

func (q *ReceiptQuery) GetReceiptByID(ctx context.Context, args *shopping.IDQueryShopArg) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ID(args.ID).ShopID(args.ShopID).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) ListReceipts(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Filters(args.Filters)
	receipts, err := query.Paging(args.Paging).ListReceipts()
	if err != nil {
		return nil, err
	}

	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	totalAmountConfirmedReceipt, totalAmountConfirmedPayment, err := query.Status(model.S3Positive).SumAmountReceiptAndPayment()
	if err != nil {
		return nil, err
	}

	return &receipting.ReceiptsResponse{
		Receipts:                    receipts,
		Count:                       int32(count),
		TotalAmountConfirmedReceipt: totalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: totalAmountConfirmedPayment,
	}, nil
}

func (q *ReceiptQuery) ListReceiptsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*receipting.ReceiptsResponse, error) {
	panic("implement me")
}

func (q *ReceiptQuery) ListReceiptsByRefIDsAndStatus(
	ctx context.Context, args *receipting.ListReceiptsByRefIDsAndStatusArgs,
) (*receipting.ReceiptsResponse, error) {
	receipts, err := q.store(ctx).ShopID(args.ShopID).RefIDs(args.RefIDs...).Status(model.Status3(args.Status)).ListReceipts()
	if err != nil {
		return nil, err
	}
	return &receipting.ReceiptsResponse{Receipts: receipts}, nil
}

func (q *ReceiptQuery) GetReceiptByCode(ctx context.Context, code string, shopID int64) (*receipting.Receipt, error) {
	receipt, err := q.store(ctx).ShopID(shopID).Code(code).GetReceipt()
	return receipt, err
}

func (q *ReceiptQuery) ListReceiptsByTraderIDs(
	ctx context.Context, shopID int64, traderIDs []int64,
) (*receipting.ReceiptsResponse, error) {
	receipts, err := q.store(ctx).ShopID(shopID).TraderIDs(traderIDs...).ListReceipts()
	return &receipting.ReceiptsResponse{Receipts: receipts}, err
}

func (q *ReceiptQuery) ListReceiptsByLedgerIDs(
	ctx context.Context, args *receipting.ListReceiptsByLedgerIDsArgs,
) (*receipting.ReceiptsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).LedgerIDs(args.LedgerIDs...).Filters(args.Filters)
	receipts, err := query.Paging(args.Paging).ListReceipts()
	if err != nil {
		return nil, err
	}

	count, err := query.Count()
	if err != nil {
		return nil, err
	}

	totalAmountConfirmedReceipt, totalAmountConfirmedPayment, err := query.Status(model.S3Positive).SumAmountReceiptAndPayment()
	if err != nil {
		return nil, err
	}

	return &receipting.ReceiptsResponse{
		Receipts:                    receipts,
		Count:                       int32(count),
		TotalAmountConfirmedReceipt: totalAmountConfirmedReceipt,
		TotalAmountConfirmedPayment: totalAmountConfirmedPayment,
	}, nil
}
