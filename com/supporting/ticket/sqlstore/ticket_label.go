package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/convert"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TicketLabelStoreFactory func(ctx context.Context) *TicketLabelStore

func NewTicketLabelStore(db *cmsql.Database) TicketLabelStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketLabelStore {
		return &TicketLabelStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketLabelStore struct {
	ft TicketLabelFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *TicketLabelStore) WithPaging(paging meta.Paging) *TicketLabelStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *TicketLabelStore) ID(id dot.ID) *TicketLabelStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TicketLabelStore) IDs(ids ...dot.ID) *TicketLabelStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *TicketLabelStore) Code(code string) *TicketLabelStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *TicketLabelStore) GetTicketLabelDB() (*model.TicketLabel, error) {
	query := s.query().Where(s.preds)
	var ticketDB model.TicketLabel
	err := query.ShouldGet(&ticketDB)
	return &ticketDB, err
}

func (s *TicketLabelStore) GetTicketLabel() (ticketResult *ticket.TicketLabel, err error) {
	ticketDB, err := s.GetTicketLabelDB()
	if err != nil {
		return nil, err
	}
	ticketResult = convert.Convert_ticketmodel_TicketLabel_ticket_TicketLabel(ticketDB, ticketResult)
	return ticketResult, nil
}

func (s *TicketLabelStore) ListTicketLabelsDB() ([]*model.TicketLabel, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}

	var labels model.TicketLabels
	err := query.Find(&labels)
	return labels, err
}

func (s *TicketLabelStore) ListTicketLabels() ([]*ticket.TicketLabel, error) {
	ticketsDB, err := s.ListTicketLabelsDB()
	if err != nil {
		return nil, err
	}
	tickets := convert.Convert_ticketmodel_TicketLabels_ticket_TicketLabels(ticketsDB)
	return tickets, nil
}

func (s *TicketLabelStore) Create(args *ticket.TicketLabel) error {
	var label = model.TicketLabel{}
	convert.Convert_ticket_TicketLabel_ticketmodel_TicketLabel(args, &label)
	return s.CreateDB(&label)
}

func (s *TicketLabelStore) CreateDB(TicketLabel *model.TicketLabel) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(TicketLabel)
}

func (s *TicketLabelStore) UpdateTicketLabel(args *ticket.TicketLabel) error {
	var result = &model.TicketLabel{}
	result = convert.Convert_ticket_TicketLabel_ticketmodel_TicketLabel(args, result)
	return s.UpdateTicketLabelDB(result)
}

func (s *TicketLabelStore) UpdateTicketLabelDB(args *model.TicketLabel) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *TicketLabelStore) Delete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Delete((*model.TicketLabel)(nil))
}
