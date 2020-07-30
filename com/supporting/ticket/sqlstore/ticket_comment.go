package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/convert"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TicketCommentStoreFactory func(ctx context.Context) *TicketCommentStore

func NewTicketCommentStore(db *cmsql.Database) TicketCommentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TicketCommentStore {
		return &TicketCommentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TicketCommentStore struct {
	ft TicketCommentFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *TicketCommentStore) WithPaging(paging meta.Paging) *TicketCommentStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *TicketCommentStore) ID(id dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}
func (s *TicketCommentStore) IDs(ids ...dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *TicketCommentStore) CreatedBy(id dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, s.ft.ByCreatedBy(id).Optional())
	return s
}

func (s *TicketCommentStore) TicketID(ticketID dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, s.ft.ByTicketID(ticketID))
	return s
}
func (s *TicketCommentStore) ParentID(parentID dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, s.ft.ByParentID(parentID))
	return s
}

func (s *TicketCommentStore) AccountID(accountID dot.ID) *TicketCommentStore {
	s.preds = append(s.preds, s.ft.ByAccountID(accountID))
	return s
}

func (s *TicketCommentStore) GetTicketCommentDB() (*model.TicketComment, error) {
	query := s.query().Where(s.preds)
	var ticketDB model.TicketComment
	err := query.ShouldGet(&ticketDB)
	return &ticketDB, err
}

func (s *TicketCommentStore) GetTicketComment() (ticketResult *ticket.TicketComment, err error) {
	ticketDB, err := s.GetTicketCommentDB()
	if err != nil {
		return nil, err
	}
	ticketResult = convert.Convert_ticketmodel_TicketComment_ticket_TicketComment(ticketDB, ticketResult)
	return ticketResult, nil
}

func (s *TicketCommentStore) ListTicketCommentsDB() ([]*model.TicketComment, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortTicketComment)
	if err != nil {
		return nil, err
	}
	var labels model.TicketComments
	err = query.Find(&labels)
	return labels, err
}

func (s *TicketCommentStore) ListTicketComments() ([]*ticket.TicketComment, error) {
	ticketsDB, err := s.ListTicketCommentsDB()
	if err != nil {
		return nil, err
	}
	tickets := convert.Convert_ticketmodel_TicketComments_ticket_TicketComments(ticketsDB)
	return tickets, nil
}

func (s *TicketCommentStore) Create(args *ticket.TicketComment) error {
	var label = model.TicketComment{}
	convert.Convert_ticket_TicketComment_ticketmodel_TicketComment(args, &label)
	return s.CreateDB(&label)
}

func (s *TicketCommentStore) CreateDB(TicketComment *model.TicketComment) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(TicketComment)
}

func (s *TicketCommentStore) UpdateTicketCommentDB(args *model.TicketComment) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *TicketCommentStore) SoftDelete(deletedBy dot.ID) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("ticket_comment").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	})
	return _deleted, err
}
