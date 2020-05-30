package subscriptionbill

import (
	"context"

	"o.o/api/subscripting/subscriptionbill"
	com "o.o/backend/com/main"
	"o.o/backend/com/subscripting/subscriptionbill/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ subscriptionbill.QueryService = &SubrBillQuery{}

type SubrBillQuery struct {
	subrBillStore sqlstore.SubrBillStoreFactory
}

func NewSubrBillQuery(db com.MainDB) *SubrBillQuery {
	return &SubrBillQuery{
		subrBillStore: sqlstore.NewSubrBillStore(db),
	}
}

func SubrBillQueryMessageBus(q *SubrBillQuery) subscriptionbill.QueryBus {
	b := bus.New()
	return subscriptionbill.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SubrBillQuery) GetSubscriptionBillByID(ctx context.Context, id dot.ID, accountID dot.ID) (*subscriptionbill.SubscriptionBillFtLine, error) {
	return q.subrBillStore(ctx).ID(id).OptionalAccountID(accountID).GetSubrBillFtLine()
}

func (q *SubrBillQuery) ListSubscriptionBills(ctx context.Context, args *subscriptionbill.ListSubscriptionBillsArgs) (*subscriptionbill.ListSubscriptionBillsResponse, error) {
	query := q.subrBillStore(ctx).OptionalAccountID(args.AccountID).OptionalSubscriptionID(args.SubscriptionID).WithPaging(args.Paging).Filters(args.Filters)
	res, err := query.ListSubrBillFtLines()
	if err != nil {
		return nil, err
	}
	return &subscriptionbill.ListSubscriptionBillsResponse{
		SubscriptionBills: res,
		Paging:            query.GetPaging(),
	}, nil
}
