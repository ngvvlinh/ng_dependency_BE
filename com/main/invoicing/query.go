package invoicing

import (
	"context"

	"o.o/api/main/invoicing"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/invoicing/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ invoicing.QueryService = &InvoiceQuery{}

type InvoiceQuery struct {
	invoiceStore sqlstore.InvoiceStoreFactory
}

func NewInvoiceQuery(db com.MainDB) *InvoiceQuery {
	return &InvoiceQuery{
		invoiceStore: sqlstore.NewInvoiceStore(db),
	}
}

func InvoiceQueryMessageBus(q *InvoiceQuery) invoicing.QueryBus {
	b := bus.New()
	return invoicing.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *InvoiceQuery) GetInvoiceByID(ctx context.Context, id dot.ID, accountID dot.ID) (*invoicing.InvoiceFtLine, error) {
	return q.invoiceStore(ctx).ID(id).OptionalAccountID(accountID).GetInvoiceFtLine()
}

func (q *InvoiceQuery) GetInvoiceByPaymentID(ctx context.Context, paymentID dot.ID) (*invoicing.InvoiceFtLine, error) {
	return q.invoiceStore(ctx).PaymentID(paymentID).GetInvoiceFtLine()
}

func (q *InvoiceQuery) GetInvoiceByReferral(
	ctx context.Context, args *invoicing.GetInvoiceByReferralArgs,
) (*invoicing.InvoiceFtLine, error) {
	return q.invoiceStore(ctx).ReferralType(args.ReferralType).ReferralIDs(args.ReferralIDs).GetInvoiceFtLine()
}

func (q *InvoiceQuery) ListInvoices(ctx context.Context, args *invoicing.ListInvoicesArgs) (*invoicing.ListInvoicesResponse, error) {
	if args.DateTo.Before(args.DateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}

	query := q.invoiceStore(ctx).OptionalAccountID(args.AccountID).WithPaging(args.Paging).Filters(args.Filters)
	if !args.DateFrom.IsZero() {
		query = query.BetweenDateFromAndDateTo(args.DateFrom, args.DateTo)
	}
	if args.RefID != 0 {
		query = query.ReferralID(args.RefID)
	}
	if args.RefType != 0 {
		query = query.ReferralType(args.RefType)
	}
	if args.Type != 0 {
		query = query.ByType(args.Type)
	}

	res, err := query.ListInvoiceFtLines()
	if err != nil {
		return nil, err
	}
	return &invoicing.ListInvoicesResponse{
		Invoices: res,
		Paging:   query.GetPaging(),
	}, nil
}
