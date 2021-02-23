package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/invoice"
	"o.o/backend/com/subscripting/invoice/convert"
	"o.o/backend/com/subscripting/invoice/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type InvoiceStore struct {
	ft      InvoiceFilters
	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
	ctx              context.Context
	invoiceLineStore *InvoiceLineStore

	includeDeleted sqlstore.IncludeDeleted
}

type InvoiceStoreFactory func(context.Context) *InvoiceStore

func NewInvoiceStore(db *cmsql.Database) InvoiceStoreFactory {
	return func(ctx context.Context) *InvoiceStore {
		return &InvoiceStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			invoiceLineStore: NewInvoiceLineStore(db)(ctx),
			ctx:              ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *InvoiceFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *InvoiceFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *InvoiceStore) WithPaging(paging meta.Paging) *InvoiceStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *InvoiceStore) Filters(filters meta.Filters) *InvoiceStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *InvoiceStore) ID(id dot.ID) *InvoiceStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *InvoiceStore) AccountID(id dot.ID) *InvoiceStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *InvoiceStore) OptionalAccountID(id dot.ID) *InvoiceStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id).Optional())
	return s
}

func (s *InvoiceStore) GetInvoiceFtLineDB() (*model.InvoiceFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var inv model.Invoice
	if err := query.ShouldGet(&inv); err != nil {
		return nil, err
	}

	lines, err := s.invoiceLineStore.InvoiceID(inv.ID).ListInvoiceLinesDB()
	if err != nil {
		return nil, err
	}

	return &model.InvoiceFtLine{
		Invoice: &inv,
		Lines:   lines,
	}, nil
}

func (s *InvoiceStore) GetInvoiceFtLine() (*invoice.InvoiceFtLine, error) {
	subr, err := s.GetInvoiceFtLineDB()
	if err != nil {
		return nil, err
	}
	var res invoice.InvoiceFtLine
	if err = scheme.Convert(subr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *InvoiceStore) ListInvoiceFtLinesDB() ([]*model.InvoiceFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortInvoice, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterInvoice)
	if err != nil {
		return nil, err
	}

	var invoices model.Invoices
	if err = query.Find(&invoices); err != nil {
		return nil, err
	}
	invoiceIDs := make([]dot.ID, len(invoices))
	for i, inv := range invoices {
		invoiceIDs[i] = inv.ID
	}
	invoiceLines, err := s.invoiceLineStore.InvoiceIDs(invoiceIDs...).ListInvoiceLinesDB()
	if err != nil {
		return nil, err
	}
	invoiceLinesMap := make(map[dot.ID][]*model.InvoiceLine)
	for _, line := range invoiceLines {
		invoiceLinesMap[line.InvoiceID] = append(invoiceLinesMap[line.InvoiceID], line)
	}

	var res []*model.InvoiceFtLine
	for _, inv := range invoices {
		res = append(res, &model.InvoiceFtLine{
			Invoice: inv,
			Lines:   invoiceLinesMap[inv.ID],
		})
	}
	return res, nil
}

func (s *InvoiceStore) ListInvoiceFtLines() (res []*invoice.InvoiceFtLine, err error) {
	invs, err := s.ListInvoiceFtLinesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(invs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *InvoiceStore) CreateInvoiceDB(inv *model.Invoice) error {
	sqlstore.MustNoPreds(s.preds)
	if inv.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	inv.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	return s.query().ShouldInsert(inv)
}

func (s *InvoiceStore) CreateInvoice(inv *invoice.Invoice) error {
	var invDB model.Invoice
	if err := scheme.Convert(inv, &invDB); err != nil {
		return err
	}
	return s.CreateInvoiceDB(&invDB)
}

func (s *InvoiceStore) UpdateInvoiceDB(inv *model.Invoice) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must provide preds")
	}
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(inv)
}

func (s *InvoiceStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("invoice").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *InvoiceStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
