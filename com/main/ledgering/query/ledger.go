package query

import (
	"context"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/ledgering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi/dot"
)

var _ ledgering.QueryService = &LedgerQuery{}

type LedgerQuery struct {
	store sqlstore.LedgerStoreFactory
}

func NewLedgerQuery(db *cmsql.Database) *LedgerQuery {
	return &LedgerQuery{
		store: sqlstore.NewLedgerStore(db),
	}
}

func (q *LedgerQuery) MessageBus() ledgering.QueryBus {
	b := bus.New()
	return ledgering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *LedgerQuery) GetLedgerByID(ctx context.Context, args *shopping.IDQueryShopArg) (*ledgering.ShopLedger, error) {
	ledger, err := q.store(ctx).ID(args.ID).ShopID(args.ShopID).GetLedger()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "không tìm thấy tài khoản thanh toán").
			Throw()
	}

	return ledger, nil
}

func (q *LedgerQuery) ListLedgers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (_ *ledgering.ShopLedgersResponse, err error) {
	query := q.store(ctx).ShopID(args.ShopID).Paging(args.Paging).Filters(args.Filters)
	ledgers, err := query.ListLedgers()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &ledgering.ShopLedgersResponse{
		Ledgers: ledgers,
		Count:   int(count),
	}, nil
}

func (q *LedgerQuery) ListLedgersByIDs(
	ctx context.Context, shopID dot.ID, IDs []dot.ID,
) (*ledgering.ShopLedgersResponse, error) {
	query := q.store(ctx).ShopID(shopID).IDs(IDs...)
	ledgers, err := query.ListLedgers()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &ledgering.ShopLedgersResponse{
		Ledgers: ledgers,
		Count:   int(count),
	}, nil
}

func (q *LedgerQuery) ListLedgersByType(
	ctx context.Context, ledgerType ledgering.LedgerType, shopID dot.ID,
) (*ledgering.ShopLedgersResponse, error) {
	query := q.store(ctx).ShopID(shopID).Type(string(ledgerType))
	ledgers, err := query.ListLedgers()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &ledgering.ShopLedgersResponse{
		Ledgers: ledgers,
		Count:   int(count),
	}, nil
}

func (q *LedgerQuery) GetLedgerByAccountNumber(
	ctx context.Context, accountNumber string, shopID dot.ID,
) (*ledgering.ShopLedger, error) {
	ledger, err := q.store(ctx).ShopID(shopID).AccountNumber(accountNumber).GetLedger()
	if err != nil {
		return nil, err
	}
	return ledger, nil
}
