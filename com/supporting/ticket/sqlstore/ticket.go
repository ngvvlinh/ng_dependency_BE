package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/backend/com/supporting/ticket/convert"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type TicketStoreFactory func(ctx context.Context) *TicketStore

func NewTicketStore(db *cmsql.Database) TicketStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketStore {
		return &TicketStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketStore struct {
	ft TicketFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *TicketStore) WithPaging(paging meta.Paging) *TicketStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *TicketStore) ID(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TicketStore) IDs(ids ...dot.ID) *TicketStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *TicketStore) RefTicketID(id dot.ID) *TicketStore {
	s.preds = append(s.preds, sq.In("ref_ticket_id", id))
	return s
}

func (s *TicketStore) AccountID(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *TicketStore) OptionalAccountID(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id).Optional())
	return s
}

func (s *TicketStore) CreatedBy(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByCreatedBy(id).Optional())
	return s
}

func (s *TicketStore) ClosedBy(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByClosedBy(id).Optional())
	return s
}

func (s *TicketStore) ConfirmedBy(id dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByConfirmedBy(id).Optional())
	return s
}

func (s *TicketStore) Source(source ticket_source.TicketSource) *TicketStore {
	s.preds = append(s.preds, s.ft.BySource(source).Optional())
	return s
}

func (s *TicketStore) RefType(refType ticket_ref_type.TicketRefType) *TicketStore {
	s.preds = append(s.preds, s.ft.ByRefType(refType).Optional())
	return s
}

func (s *TicketStore) RefID(refID dot.ID) *TicketStore {
	s.preds = append(s.preds, s.ft.ByRefID(refID).Optional())
	return s
}

func (s *TicketStore) Code(code string) *TicketStore {
	s.preds = append(s.preds, s.ft.ByCode(code).Optional())
	return s
}

func (s *TicketStore) State(state ticket_state.TicketState) *TicketStore {
	s.preds = append(s.preds, s.ft.ByState(state).Optional())
	return s
}

func (s *TicketStore) LabelIDs(labelIDs []dot.ID) *TicketStore {
	s.preds = append(s.preds, sq.NewExpr("label_ids @> ?", core.Array{V: labelIDs}))
	return s
}

func (s *TicketStore) AssignedUserIDs(userIDs []dot.ID) *TicketStore {
	s.preds = append(s.preds, sq.NewExpr("assigned_user_ids @> ?", core.Array{V: userIDs}))
	return s
}

func (s *TicketStore) GetTicketDB() (*model.Ticket, error) {
	query := s.query().Where(s.preds)
	var ticketDB model.Ticket
	err := query.ShouldGet(&ticketDB)
	return &ticketDB, err
}

func (s *TicketStore) GetTicket() (ticketResult *ticket.Ticket, err error) {
	ticketDB, err := s.GetTicketDB()
	if err != nil {
		return nil, err
	}
	ticketResult = convert.Convert_ticketmodel_Ticket_ticket_Ticket(ticketDB, ticketResult)
	return ticketResult, nil
}

func (s *TicketStore) ListTicketsDB() ([]*model.Ticket, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortTicket)
	if err != nil {
		return nil, err
	}

	var ticketExt model.TicketExtendeds
	err = query.Find(&ticketExt)
	if err != nil {
		return nil, err
	}
	var result []*model.Ticket
	for _, v := range ticketExt {
		result = append(result, v.Ticket)
	}
	return result, nil
}

func (s *TicketStore) ListTickets() ([]*ticket.Ticket, error) {
	ticketsDB, err := s.ListTicketsDB()
	if err != nil {
		return nil, err
	}
	tickets := convert.Convert_ticketmodel_Tickets_ticket_Tickets(ticketsDB)
	return tickets, nil
}
func (s *TicketStore) UpdateTicketALL(args *ticket.Ticket) error {
	var ticketDB = &model.Ticket{}
	convert.Convert_ticket_Ticket_ticketmodel_Ticket(args, ticketDB)
	return s.UpdateTicketALLDB(ticketDB)
}

func (s *TicketStore) UpdateTicketALLDB(args *model.Ticket) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *TicketStore) UpdateTicketDB(args *model.Ticket) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *TicketStore) Create(args *ticket.Ticket) error {
	var voucherDB = model.Ticket{}
	convert.Convert_ticket_Ticket_ticketmodel_Ticket(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *TicketStore) CreateDB(Ticket *model.Ticket) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().ShouldInsert(Ticket)
	if err != nil {
		return err
	}
	err = s.query().ShouldInsert(&model.TicketSearch{
		ID:        Ticket.ID,
		TitleNorm: validate.NormalizeSearchCharacter(Ticket.Title),
	})
	return err
}

func (s *TicketStore) GetTicketByMaximumCodeNorm() (*model.Ticket, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0")
	query = query.OrderBy("code_norm desc").Limit(1)

	var ticketDB model.Ticket
	if err := query.ShouldGet(&ticketDB); err != nil {
		return nil, err
	}
	return &ticketDB, nil
}

// Only use this function when get model.ShopExtended
func (s *TicketStore) TitleFullTextSearch(name filter.FullTextSearch) *TicketStore {
	s.preds = append(s.preds, s.ft.Filter(`ts.title_norm @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
}
