package sqlstore

import (
	"context"

	"o.o/api/subscripting/invoice"
	"o.o/backend/com/subscripting/invoice/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type InvoiceLineStore struct {
	ft    InvoiceLineFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type InvoiceLineStoreFactory func(context.Context) *InvoiceLineStore

func NewInvoiceLineStore(db *cmsql.Database) InvoiceLineStoreFactory {
	return func(ctx context.Context) *InvoiceLineStore {
		return &InvoiceLineStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *InvoiceLineStore) InvoiceID(id dot.ID) *InvoiceLineStore {
	s.preds = append(s.preds, s.ft.ByInvoiceID(id))
	return s
}

func (s *InvoiceLineStore) InvoiceIDs(ids ...dot.ID) *InvoiceLineStore {
	s.preds = append(s.preds, sq.In("invoice_id", ids))
	return s
}

func (s *InvoiceLineStore) CreateInvoiceLineDB(line *model.InvoiceLine) error {
	if line.ID == 0 {
		line.ID = cm.NewID()
	}
	if line.InvoiceID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing invoice ID")
	}
	return s.query().ShouldInsert(line)
}

func (s *InvoiceLineStore) CreateInvoiceLine(args *invoice.InvoiceLine) error {
	var line model.InvoiceLine
	if err := scheme.Convert(args, &line); err != nil {
		return err
	}
	return s.CreateInvoiceLineDB(&line)
}

func (s *InvoiceLineStore) ListInvoiceLinesDB() (res []*model.InvoiceLine, err error) {
	err = s.query().Where(s.preds).Find((*model.InvoiceLines)(&res))
	return
}
