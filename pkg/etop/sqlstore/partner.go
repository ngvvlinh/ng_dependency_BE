package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/top/types/etc/account_type"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/backend/pkg/etop/authorize/authkey"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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
	ft    PartnerFilters
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

func (s *PartnerStore) Get() (*model.Partner, error) {
	var item model.Partner
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *PartnerStore) List() ([]*model.Partner, error) {
	var items model.Partners
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).Find(&items)
	return items, err
}

func CreatePartner(ctx context.Context, cmd *model.CreatePartnerCommand) error {

	partner := cmd.Partner
	if partner.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	partner.ID = cm.NewIDWithTag(model.TagPartner)
	if err := partner.BeforeInsert(); err != nil {
		return err
	}

	account := &model.Account{
		ID:       partner.ID,
		OwnerID:  partner.OwnerID,
		Name:     partner.Name,
		Type:     account_type.Partner,
		ImageURL: partner.ImageURL,
		URLSlug:  "",
	}
	accountUser := &model.AccountUser{
		AccountID:            partner.ID,
		UserID:               partner.OwnerID,
		Status:               1,
		ResponseStatus:       0,
		CreatedAt:            time.Time{},
		UpdatedAt:            time.Time{},
		DeletedAt:            time.Time{},
		Permission:           model.Permission{},
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

func GetPartner(ctx context.Context, query *model.GetPartner) error {
	if query.PartnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing PartnerID")
	}

	// TODO: handle disabled partners
	var partner model.Partner
	err := x.Where("id = ?", query.PartnerID).ShouldGet(&partner)
	query.Result.Partner = &partner
	return err
}

func GetPartnerRelationQuery(ctx context.Context, query *model.GetPartnerRelationQuery) error {
	s := x.NewQuery()
	count := 0
	if query.PartnerID != 0 && query.AccountID != 0 {
		count++
		s = s.Where("pr.partner_id = ?", query.PartnerID).
			Where("pr.subject_id = ? AND pr.subject_type = 'account'", query.AccountID)
	}
	if query.AuthKey != "" {
		count++
		s = s.Where("pr.auth_key = ?", query.AuthKey)
	}
	if query.ExternalAccountID != "" {
		count++
		s = s.Where("pr.external_subject_id = ?", query.ExternalAccountID)
	}
	if count == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	var item model.PartnerRelationFtShop
	err := s.ShouldGet(&item)
	query.Result.PartnerRelationFtShop = item
	return err
}

func GetPartnerRelations(ctx context.Context, query *model.GetPartnerRelationsQuery) error {
	if query.PartnerID == 0 || query.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return x.Where("pr.partner_id = ?", query.PartnerID).
		Where("s.owner_id = ?", query.OwnerID).
		Find((*model.PartnerRelationFtShops)(&query.Result.Relations))
}

func GetPartnersFromRelation(ctx context.Context, query *model.GetPartnersFromRelationQuery) error {
	if len(query.AccountIDs) == 0 {
		return nil
	}

	var partnerIDs []dot.ID
	if err := x.SQL(`SELECT array_agg(partner_id) FROM partner_relation`).
		Where("subject_type = ?", model.SubjectTypeAccount).
		In("subject_id", query.AccountIDs).
		Scan(core.ArrayScanner(&partnerIDs)); err != nil {
		return err
	}

	partners, err := Partner(ctx).IncludeDeleted().IDs(partnerIDs...).List()
	query.Result.Partners = partners
	return err
}

// TODO: update old relation if exists
func CreatePartnerRelation(ctx context.Context, cmd *model.CreatePartnerRelationCommand) error {
	if cmd.PartnerID == 0 || cmd.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	key := authkey.GenerateAuthKey(authkey.TypePartnerShopKey, cmd.AccountID)
	rel := &model.PartnerRelation{
		AuthKey:           key,
		PartnerID:         cmd.PartnerID,
		SubjectID:         cmd.AccountID,
		SubjectType:       model.SubjectTypeAccount,
		ExternalSubjectID: cmd.ExternalID,
		Nonce:             cm.NewID(), // TODO: use crypto/rand
		Status:            1,
	}
	cmd.Result.PartnerRelation = rel
	err := x.ShouldInsert(rel)
	return err
}

func UpdatePartnerRelationCommand(ctx context.Context, cmd *model.UpdatePartnerRelationCommand) error {
	if cmd.PartnerID == 0 || cmd.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}

	return x.Where("subject_type = ?", model.SubjectTypeAccount).
		Where("partner_id = ?", cmd.PartnerID).
		Where("subject_id = ?", cmd.AccountID).
		Where("external_subject_id IS NULL").
		ShouldUpdate(&model.PartnerRelation{
			ExternalSubjectID: cmd.ExternalID,
		})
}

func GetPartners(ctx context.Context, query *model.GetPartnersQuery) error {
	var partners []*model.Partner
	s := x.Table("partner").Where("deleted_at IS NULL")
	if query.AvailableFromEtop {
		s = s.Where("available_from_etop = true")
	}
	if err := s.Find((*model.Partners)(&partners)); err != nil {
		return err
	}
	query.Result.Partners = partners
	return nil
}
