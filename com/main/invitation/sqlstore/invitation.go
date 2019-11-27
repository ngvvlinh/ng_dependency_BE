package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/invitation"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/invitation/convert"
	"etop.vn/backend/com/main/invitation/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type InvitationStoreFactory func(ctx context.Context) *InvitationStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewInvitationStore(db *cmsql.Database) InvitationStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *InvitationStore {
		return &InvitationStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type InvitationStore struct {
	ft InvitationFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *InvitationStore) Paging(paging meta.Paging) *InvitationStore {
	s.paging = paging
	return s
}

func (s *InvitationStore) GetPaing() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *InvitationStore) Filters(filters meta.Filters) *InvitationStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *InvitationStore) ID(id dot.ID) *InvitationStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *InvitationStore) AccountID(id dot.ID) *InvitationStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *InvitationStore) Token(token string) *InvitationStore {
	s.preds = append(s.preds, s.ft.ByToken(token))
	return s
}

func (s *InvitationStore) Email(email string) *InvitationStore {
	s.preds = append(s.preds, s.ft.ByEmail(email))
	return s
}

func (s *InvitationStore) Status(status etop.Status3) *InvitationStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *InvitationStore) NotExpires(t time.Time) *InvitationStore {
	s.preds = append(s.preds, sq.NewExpr("expires_at > ?", t))
	return s
}

func (s *InvitationStore) ExpiresAt(t time.Time) *InvitationStore {
	s.preds = append(s.preds, sq.NewExpr("expires_at <= ?", t))
	return s
}

func (s *InvitationStore) RejectedAt(t *time.Time) *InvitationStore {
	s.preds = append(s.preds, sq.NewExpr("rejected_at <= ? or rejected_at IS NULL", t))
	return s
}

func (s *InvitationStore) AcceptedAt(t *time.Time) *InvitationStore {
	s.preds = append(s.preds, sq.NewExpr("accepted_at <= ? or accepted_at IS NULL", t))
	return s
}

func (s *InvitationStore) Accepted() *InvitationStore {
	s.preds = append(s.preds, sq.NewExpr("accepted_at IS NOT NULL"))
	return s
}

func (s *InvitationStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("invitation").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *InvitationStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	query, _, err := sqlstore.Filters(query, s.filters, FilterInvitation)
	if err != nil {
		return 0, err
	}

	return query.Count((*model.Invitation)(nil))
}

func (s *InvitationStore) Accept() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("invitation").UpdateMap(map[string]interface{}{
		"status":      int(etop.S3Positive),
		"accepted_at": time.Now(),
	})

	return int(_updated), err
}

func (s *InvitationStore) Reject() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_updated, err := query.Table("invitation").UpdateMap(map[string]interface{}{
		"status":      int(etop.S3Negative),
		"rejected_at": time.Now(),
	})

	return int(_updated), err
}

func (s *InvitationStore) GetInvitationDB() (*model.Invitation, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var invitation model.Invitation
	err := query.ShouldGet(&invitation)
	return &invitation, err
}

func (s *InvitationStore) GetInvitation() (invitationResult *invitation.Invitation, _ error) {
	invitation, err := s.GetInvitationDB()
	if err != nil {
		return nil, err
	}
	invitationResult = convert.Convert_invitationmodel_Invitation_invitation_Invitation(invitation, invitationResult)
	return invitationResult, nil
}

func (s *InvitationStore) CreateInvitation(invitation *invitation.Invitation) error {
	sqlstore.MustNoPreds(s.preds)
	invitationDB := new(model.Invitation)
	if err := scheme.Convert(invitation, invitationDB); err != nil {
		return err
	}
	if _, err := s.query().Insert(invitationDB); err != nil {
		return err
	}

	var tempInvitation model.Invitation
	if err := s.query().Where(s.ft.ByID(invitation.ID), s.ft.ByAccountID(invitation.AccountID)).
		ShouldGet(&tempInvitation); err != nil {
		return err
	}

	invitation.CreatedAt = tempInvitation.CreatedAt
	invitation.UpdatedAt = tempInvitation.UpdatedAt
	return nil
}

func (s *InvitationStore) ListInvitationsDB() ([]*model.Invitation, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortInvitation)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterInvitation)
	if err != nil {
		return nil, err
	}

	var invitations model.Invitations
	err = query.Find(&invitations)
	return invitations, err
}

func (s *InvitationStore) ListInvitations() (invitationsResult []*invitation.Invitation, _ error) {
	invitations, err := s.ListInvitationsDB()
	if err != nil {
		return nil, err
	}

	for _, invitationEl := range invitations {
		var invitationResult *invitation.Invitation
		invitationResult = convert.Convert_invitationmodel_Invitation_invitation_Invitation(invitationEl, invitationResult)
		invitationsResult = append(invitationsResult, invitationResult)
	}
	return invitationsResult, nil
}

func (s *InvitationStore) IncludeDeleted() *InvitationStore {
	s.includeDeleted = true
	return s
}
