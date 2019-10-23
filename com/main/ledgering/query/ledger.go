package query

import (
	"context"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/shopping"
	"etop.vn/backend/com/main/ledgering/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
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
			Wrap(cm.NotFound, "không tìm thấy sổ quỹ").
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
		Ledger: ledgers,
		Count:  int32(count),
	}, nil
}
