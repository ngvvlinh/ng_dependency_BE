package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/account_type"
	com "o.o/backend/com/main"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	identitysqlstore "o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/authorize/authkey"
	"o.o/capi/dot"
)

type PartnerStoreInterface interface {
	CreatePartner(ctx context.Context, cmd *identitymodelx.CreatePartnerCommand) error

	CreatePartnerRelation(ctx context.Context, cmd *identitymodelx.CreatePartnerRelationCommand) error

	GetPartner(ctx context.Context, query *identitymodelx.GetPartner) error

	GetPartnerRelationQuery(ctx context.Context, query *identitymodelx.GetPartnerRelationQuery) error

	GetPartnerRelations(ctx context.Context, query *identitymodelx.GetPartnerRelationsQuery) error

	GetPartners(ctx context.Context, query *identitymodelx.GetPartnersQuery) error

	GetPartnersFromRelation(ctx context.Context, query *identitymodelx.GetPartnersFromRelationQuery) error

	UpdatePartnerRelationCommand(ctx context.Context, cmd *identitymodelx.UpdatePartnerRelationCommand) error
}

type PartnerStore struct {
	query func() cmsql.QueryInterface
	ft    identitysqlstore.PartnerFilters
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type PartnerStoreFactory func(ctx context.Context) *PartnerStore

func NewPartnerStore(db com.MainDB) PartnerStoreFactory {
	return func(ctx context.Context) *PartnerStore {
		return &PartnerStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func BuildPartnerStore(db com.MainDB) *PartnerStore {
	return NewPartnerStore(db)(context.Background())
}

func BindPartnerStore(s *PartnerStore) (to PartnerStoreInterface) {
	return s
}

func (st *PartnerStore) ID(id dot.ID) *PartnerStore {
	st.preds = append(st.preds, st.ft.ByID(id))
	return st
}

func (st *PartnerStore) IDs(ids ...dot.ID) *PartnerStore {
	st.preds = append(st.preds, sq.In("id", ids))
	return st
}

func (st *PartnerStore) IncludeDeleted() *PartnerStore {
	st.includeDeleted = true
	return st
}

func (st *PartnerStore) Get() (*identitymodel.Partner, error) {
	var item identitymodel.Partner
	err := st.query().Where(st.preds...).Where(st.includeDeleted.FilterDeleted(&st.ft)).ShouldGet(&item)
	return &item, err
}

func (st *PartnerStore) List() ([]*identitymodel.Partner, error) {
	var items identitymodel.Partners
	err := st.query().Where(st.preds...).Where(st.includeDeleted.FilterDeleted(&st.ft)).Find(&items)
	return items, err
}

func (st *PartnerStore) CreatePartner(ctx context.Context, cmd *identitymodelx.CreatePartnerCommand) error {

	partner := cmd.Partner
	if partner.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	partner.ID = cm.NewIDWithTag(account_tag.TagPartner)
	if err := partner.BeforeInsert(); err != nil {
		return err
	}

	account := &identitymodel.Account{
		ID:       partner.ID,
		OwnerID:  partner.OwnerID,
		Name:     partner.Name,
		Type:     account_type.Partner,
		ImageURL: partner.ImageURL,
		URLSlug:  "",
	}
	accountUser := &identitymodel.AccountUser{
		AccountID:            partner.ID,
		UserID:               partner.OwnerID,
		Status:               1,
		ResponseStatus:       0,
		CreatedAt:            time.Time{},
		UpdatedAt:            time.Time{},
		DeletedAt:            time.Time{},
		Permission:           identitymodel.Permission{},
		FullName:             "",
		ShortName:            "",
		Position:             "",
		InvitationSentAt:     time.Time{},
		InvitationSentBy:     0,
		InvitationAcceptedAt: time.Time{},
		InvitationRejectedAt: time.Time{},
		DisabledAt:           time.Time{},
		DisabledBy:           time.Time{},
		DisableReason:        "",
	}

	err := st.query().ShouldInsert(account, partner, accountUser)
	cmd.Result.Partner = cmd.Partner
	return err
}

func (st *PartnerStore) GetPartner(ctx context.Context, query *identitymodelx.GetPartner) error {
	if query.PartnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing PartnerID")
	}

	// TODO: handle disabled partners
	var partner identitymodel.Partner
	err := st.query().Where("id = ?", query.PartnerID).ShouldGet(&partner)
	query.Result.Partner = &partner
	return err
}

func (st *PartnerStore) GetPartnerRelationQuery(ctx context.Context, query *identitymodelx.GetPartnerRelationQuery) error {
	s := st.query()
	count := 0
	if query.PartnerID != 0 && query.ExternalUserID != "" {
		count++ // TODO(vu): correctly handle this
		s = s.Where("pr.partner_id = ?", query.PartnerID).
			Where("pr.external_subject_id = ? AND pr.subject_type = 'user'", query.ExternalUserID)
		var item identitymodel.PartnerRelationFtUser
		err := s.ShouldGet(&item)
		query.Result.PartnerRelation = item.PartnerRelation
		query.Result.User = item.User
		return err
	}
	if query.PartnerID != 0 && query.AccountID != 0 {
		count++
		s = s.Where("pr.partner_id = ?", query.PartnerID).
			Where("pr.subject_id = ? AND pr.subject_type = 'account'", query.AccountID)
	}
	if query.PartnerID != 0 && query.ExternalAccountID != "" {
		count++
		s = s.Where("pr.partner_id = ?", query.PartnerID).
			Where("pr.external_subject_id = ? AND pr.subject_type = 'account'", query.ExternalAccountID)
	}
	if query.AuthKey != "" {
		count++
		s = s.Where("pr.auth_key = ?", query.AuthKey)
	}
	if count == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	var item identitymodel.PartnerRelationFtShop
	err := s.ShouldGet(&item)
	query.Result.PartnerRelationFtShop = item
	return err
}

func (st *PartnerStore) GetPartnerRelations(ctx context.Context, query *identitymodelx.GetPartnerRelationsQuery) error {
	if query.PartnerID == 0 || query.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return st.query().Where("pr.partner_id = ?", query.PartnerID).
		Where("s.owner_id = ?", query.OwnerID).
		Find((*identitymodel.PartnerRelationFtShops)(&query.Result.Relations))
}

func (st *PartnerStore) GetPartnersFromRelation(ctx context.Context, query *identitymodelx.GetPartnersFromRelationQuery) error {
	if len(query.AccountIDs) == 0 {
		return nil
	}

	var partnerIDs []dot.ID
	if err := st.query().SQL(`SELECT array_agg(partner_id) FROM partner_relation`).
		Where("subject_type = ?", identity.SubjectTypeAccount).
		In("subject_id", query.AccountIDs).
		Scan(core.ArrayScanner(&partnerIDs)); err != nil {
		return err
	}

	if len(partnerIDs) == 0 {
		return nil
	}
	partners, err := st.IncludeDeleted().IDs(partnerIDs...).List()
	query.Result.Partners = partners
	return err
}

// TODO: update old relation if exists
func (st *PartnerStore) CreatePartnerRelation(ctx context.Context, cmd *identitymodelx.CreatePartnerRelationCommand) error {
	if cmd.PartnerID == 0 {
		return cm.Errorf(cm.Internal, nil, "Missing partner_id")
	}
	if cmd.AccountID != 0 && cmd.UserID != 0 {
		return cm.Errorf(cm.Internal, nil, "Must not provide both account_id and user_id")
	}
	switch {
	case cmd.AccountID != 0:
		key := authkey.GenerateAuthKey(authkey.TypePartnerShopKey, cmd.AccountID)
		rel := &identitymodel.PartnerRelation{
			AuthKey:           key,
			PartnerID:         cmd.PartnerID,
			SubjectID:         cmd.AccountID,
			SubjectType:       identity.SubjectTypeAccount,
			ExternalSubjectID: cmd.ExternalID,
			Nonce:             cm.NewID(), // TODO: use crypto/rand
			Status:            1,
		}
		cmd.Result.PartnerRelation = rel
		return st.query().ShouldInsert(rel)

	case cmd.UserID != 0:
		key := authkey.GenerateAuthKey(authkey.TypePartnerUserKey, cmd.UserID)
		rel := &identitymodel.PartnerRelation{
			AuthKey:           key,
			PartnerID:         cmd.PartnerID,
			SubjectID:         cmd.UserID,
			SubjectType:       identity.SubjectTypeUser,
			ExternalSubjectID: cmd.ExternalID,
			Nonce:             cm.NewID(), // TODO: use crypto/rand
			Status:            1,
		}
		cmd.Result.PartnerRelation = rel
		return st.query().ShouldInsert(rel)

	default:
		return cm.Errorf(cm.Internal, nil, "Missing account_id or user_id")
	}
}

func (st *PartnerStore) UpdatePartnerRelationCommand(ctx context.Context, cmd *identitymodelx.UpdatePartnerRelationCommand) error {
	if cmd.PartnerID == 0 || cmd.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return st.query().Where("subject_type = ?", identity.SubjectTypeAccount).
		Where("partner_id = ?", cmd.PartnerID).
		Where("subject_id = ?", cmd.AccountID).
		Where("external_subject_id IS NULL").
		ShouldUpdate(&identitymodel.PartnerRelation{
			ExternalSubjectID: cmd.ExternalID,
		})
}

func (st *PartnerStore) GetPartners(ctx context.Context, query *identitymodelx.GetPartnersQuery) error {
	var partners []*identitymodel.Partner
	s := st.query().Table("partner").Where("deleted_at IS NULL")
	if query.AvailableFromEtop {
		s = s.Where("available_from_etop = true")
	}
	if err := s.Find((*identitymodel.Partners)(&partners)); err != nil {
		return err
	}
	query.Result.Partners = partners
	return nil
}
