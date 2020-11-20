package sqlstore

import (
	"context"
	"time"

	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TicketLabelTicketLabelExternalsStoreFactory func(ctx context.Context) *TicketLabelTicketLabelExternalStore

func NewTicketLabelTicketLabelExternalStore(db *cmsql.Database) TicketLabelTicketLabelExternalsStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketLabelTicketLabelExternalStore {
		return &TicketLabelTicketLabelExternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketLabelTicketLabelExternalStore struct {
	ft TicketLabelTicketLabelExternalFilters

	query cmsql.QueryFactory
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

func (s *TicketLabelTicketLabelExternalStore) TicketLabelID(ticketLabelID dot.ID) *TicketLabelTicketLabelExternalStore {
	s.preds = append(s.preds, s.ft.ByTicketLabelID(ticketLabelID))
	return s
}

func (s *TicketLabelTicketLabelExternalStore) TicketLabelIDs(ticketLabelIDs []dot.ID) *TicketLabelTicketLabelExternalStore {
	s.preds = append(s.preds, sq.In("ticket_label_id", ticketLabelIDs))
	return s
}

func (s *TicketLabelTicketLabelExternalStore) TicketLabelExternalID(ticketLabelExternalID dot.ID) *TicketLabelTicketLabelExternalStore {
	s.preds = append(s.preds, s.ft.ByTicketLabelExternalID(ticketLabelExternalID))
	return s
}

func (s *TicketLabelTicketLabelExternalStore) TicketLabelExternalIDs(ticketLabelExternalIDs []dot.ID) *TicketLabelTicketLabelExternalStore {
	s.preds = append(s.preds, sq.In("ticket_label_external_id", ticketLabelExternalIDs))
	return s
}

func (s *TicketLabelTicketLabelExternalStore) ListDB() ([]*model.TicketLabelTicketLabelExternal, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var results model.TicketLabelTicketLabelExternals
	if err := query.Find(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *TicketLabelTicketLabelExternalStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("ticket_label_ticket_label_external").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
