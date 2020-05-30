package sqlstore

import (
	"context"
	"time"

	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	identitysqlstore "o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/etop/authorize/authkey"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		CreatePartner,
		CreatePartnerRelation,
		GetPartner,
		GetPartnerRelationQuery,
		GetPartnerRelations,
		GetPartnersFromRelation,
		UpdatePartnerRelationCommand,
		GetPartners,
	)
}

type PartnerStore struct {
	ctx   context.Context
	ft    identitysqlstore.PartnerFilters
	preds []interface{}

	includeDeleted
}

func Partner(ctx context.Context) *PartnerStore {
	return &PartnerStore{ctx: ctx}
}

func (s *PartnerStore) ID(id dot.ID) *PartnerStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PartnerStore) IDs(ids ...dot.ID) *PartnerStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *PartnerStore) IncludeDeleted() *PartnerStore {
	s.includeDeleted = true
	return s
}

func (s *PartnerStore) Get() (*identitymodel.Partner, error) {
	var item identitymodel.Partner
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *PartnerStore) List() ([]*identitymodel.Partner, error) {
	var items identitymodel.Partners
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).Find(&items)
	return items, err
}

func CreatePartner(ctx context.Context, cmd *identitymodelx.CreatePartnerCommand) error {

	partner := cmd.Partner
	if partner.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	partner.ID = cm.NewIDWithTag(model.TagPartner)
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

	err := x.ShouldInsert(account, partner, accountUser)
	cmd.Result.Partner = cmd.Partner
	return err
}

func GetPartner(ctx context.Context, query *identitymodelx.GetPartner) error {
	if query.PartnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing PartnerID")
	}

	// TODO: handle disabled partners
	var partner identitymodel.Partner
	err := x.Where("id = ?", query.PartnerID).ShouldGet(&partner)
	query.Result.Partner = &partner
	return err
}

func GetPartnerRelationQuery(ctx context.Context, query *identitymodelx.GetPartnerRelationQuery) error {
	s := x.NewQuery()
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

func GetPartnerRelations(ctx context.Context, query *identitymodelx.GetPartnerRelationsQuery) error {
	if query.PartnerID == 0 || query.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return x.Where("pr.partner_id = ?", query.PartnerID).
		Where("s.owner_id = ?", query.OwnerID).
		Find((*identitymodel.PartnerRelationFtShops)(&query.Result.Relations))
}

func GetPartnersFromRelation(ctx context.Context, query *identitymodelx.GetPartnersFromRelationQuery) error {
	if len(query.AccountIDs) == 0 {
		return nil
	}

	var partnerIDs []dot.ID
	if err := x.SQL(`SELECT array_agg(partner_id) FROM partner_relation`).
		Where("subject_type = ?", identitymodel.SubjectTypeAccount).
		In("subject_id", query.AccountIDs).
		Scan(core.ArrayScanner(&partnerIDs)); err != nil {
		return err
	}

	partners, err := Partner(ctx).IncludeDeleted().IDs(partnerIDs...).List()
	query.Result.Partners = partners
	return err
}

// TODO: update old relation if exists
func CreatePartnerRelation(ctx context.Context, cmd *identitymodelx.CreatePartnerRelationCommand) error {
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
			SubjectType:       identitymodel.SubjectTypeAccount,
			ExternalSubjectID: cmd.ExternalID,
			Nonce:             cm.NewID(), // TODO: use crypto/rand
			Status:            1,
		}
		cmd.Result.PartnerRelation = rel
		return x.ShouldInsert(rel)

	case cmd.UserID != 0:
		key := authkey.GenerateAuthKey(authkey.TypePartnerUserKey, cmd.UserID)
		rel := &identitymodel.PartnerRelation{
			AuthKey:           key,
			PartnerID:         cmd.PartnerID,
			SubjectID:         cmd.UserID,
			SubjectType:       identitymodel.SubjectTypeUser,
			ExternalSubjectID: cmd.ExternalID,
			Nonce:             cm.NewID(), // TODO: use crypto/rand
			Status:            1,
		}
		cmd.Result.PartnerRelation = rel
		return x.ShouldInsert(rel)

	default:
		return cm.Errorf(cm.Internal, nil, "Missing account_id or user_id")
	}
}

func UpdatePartnerRelationCommand(ctx context.Context, cmd *identitymodelx.UpdatePartnerRelationCommand) error {
	if cmd.PartnerID == 0 || cmd.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return x.Where("subject_type = ?", identitymodel.SubjectTypeAccount).
		Where("partner_id = ?", cmd.PartnerID).
		Where("subject_id = ?", cmd.AccountID).
		Where("external_subject_id IS NULL").
		ShouldUpdate(&identitymodel.PartnerRelation{
			ExternalSubjectID: cmd.ExternalID,
		})
}

func GetPartners(ctx context.Context, query *identitymodelx.GetPartnersQuery) error {
	var partners []*identitymodel.Partner
	s := x.Table("partner").Where("deleted_at IS NULL")
	if query.AvailableFromEtop {
		s = s.Where("available_from_etop = true")
	}
	if err := s.Find((*identitymodel.Partners)(&partners)); err != nil {
		return err
	}
	query.Result.Partners = partners
	return nil
}
