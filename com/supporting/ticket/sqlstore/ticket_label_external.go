package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/convert"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TicketLabelExternalStoreFactory func(ctx context.Context) *TicketLabelExternalStore

func NewTicketLabelExternalStore(db *cmsql.Database) TicketLabelExternalStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketLabelExternalStore {
		return &TicketLabelExternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketLabelExternalStore struct {
	ft TicketLabelExternalFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *TicketLabelExternalStore) WithPaging(paging meta.Paging) *TicketLabelExternalStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *TicketLabelExternalStore) ID(id dot.ID) *TicketLabelExternalStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TicketLabelExternalStore) ExternalID(externalID string) *TicketLabelExternalStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *TicketLabelExternalStore) ConnectionID(connectionID dot.ID) *TicketLabelExternalStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connectionID))
	return s
}

func (s *TicketLabelExternalStore) GetTicketLabelExternalDB() (*model.TicketLabelExternal, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var ticketLabelExternalDB model.TicketLabelExternal
	err := query.ShouldGet(&ticketLabelExternalDB)
	return &ticketLabelExternalDB, err
}

func (s *TicketLabelExternalStore) GetTicketLabelExternal() (ticketResult *ticket.TicketLabelExternal, err error) {
	ticketLabelExternalDB, err := s.GetTicketLabelExternalDB()
	if err != nil {
		return nil, err
	}
	ticketResult = convert.Convert_ticketmodel_TicketLabelExternal_ticket_TicketLabelExternal(ticketLabelExternalDB, ticketResult)
	return ticketResult, nil
}

func (s *TicketLabelExternalStore) ListTicketLabelExternalDB() ([]*model.TicketLabelExternal, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortTicketComment)
	if err != nil {
		return nil, err
	}
	var tickets model.TicketLabelExternals
	err = query.Find(&tickets)
	return tickets, err
}

func (s *TicketLabelExternalStore) ListTicketLabelExternals() ([]*ticket.TicketLabelExternal, error) {
	ticketsDB, err := s.ListTicketLabelExternalDB()
	if err != nil {
		return nil, err
	}
	tickets := convert.Convert_ticketmodel_TicketLabelExternals_ticket_TicketLabelExternals(ticketsDB)
	return tickets, nil
}

func (s *TicketLabelExternalStore) CreateDB(ticketLabelExternal *model.TicketLabelExternal) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(ticketLabelExternal)
}

func (s *TicketLabelExternalStore) Create(args *ticket.TicketLabelExternal) error {
	ticketLabelExternalDB := convert.Convert_ticket_TicketLabelExternal_ticketmodel_TicketLabelExternal(args, nil)
	return s.CreateDB(ticketLabelExternalDB)
}

func (s *TicketLabelExternalStore) UpdateTicketLabelExternalDB(args *model.TicketLabelExternal) error {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(args)
}

func (s *TicketLabelExternalStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("ticket_label_external").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
