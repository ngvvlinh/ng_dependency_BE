package invoice

import (
	"context"

	"o.o/api/subscripting/invoice"
	com "o.o/backend/com/main"
	"o.o/backend/com/subscripting/invoice/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ invoice.QueryService = &InvoiceQuery{}

type InvoiceQuery struct {
	invoiceStore sqlstore.InvoiceStoreFactory
}

func NewInvoiceQuery(db com.MainDB) *InvoiceQuery {
	return &InvoiceQuery{
		invoiceStore: sqlstore.NewInvoiceStore(db),
	}
}

func InvoiceQueryMessageBus(q *InvoiceQuery) invoice.QueryBus {
	b := bus.New()
	return invoice.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *InvoiceQuery) GetInvoiceByID(ctx context.Context, id dot.ID, accountID dot.ID) (*invoice.InvoiceFtLine, error) {
	return q.invoiceStore(ctx).ID(id).OptionalAccountID(accountID).GetInvoiceFtLine()
}

func (q *InvoiceQuery) ListInvoices(ctx context.Context, args *invoice.ListInvoicesArgs) (*invoice.ListInvoicesResponse, error) {
	query := q.invoiceStore(ctx).OptionalAccountID(args.AccountID).WithPaging(args.Paging).Filters(args.Filters)
	res, err := query.ListInvoiceFtLines()
	if err != nil {
		return nil, err
	}
	return &invoice.ListInvoicesResponse{
		Invoices: res,
		Paging:   query.GetPaging(),
	}, nil
}
