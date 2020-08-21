package api

import (
	"context"

	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type AccountService struct {
	session.Session
}

func (s *AccountService) Clone() api.AccountService {
	res := *s
	return &res
}

func (s *AccountService) UpdateURLSlug(ctx context.Context, r *api.UpdateURLSlugRequest) (*pbcm.Empty, error) {
	if r.AccountId == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}
	if !r.UrlSlug.Valid {
		return nil, cm.Error(cm.InvalidArgument, "Missing url_slug", nil)
	}

	accQuery := &identitymodelx.GetAllAccountRolesQuery{UserID: s.SS.User().ID}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return nil, err
	}

	var account *identitymodel.AccountUserExtended
	for _, acc := range accQuery.Result {
		if acc.Account.ID == r.AccountId {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, cm.Error(cm.InvalidArgument, "Tài khoản yêu cầu không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}

	cmd := &identitymodelx.UpdateAccountURLSlugCommand{
		AccountID: r.AccountId,
		URLSlug:   r.UrlSlug.String,
	}

	result := &pbcm.Empty{}
	return result, bus.Dispatch(ctx, cmd)
}

func (s *AccountService) GetPublicPartnerInfo(ctx context.Context, q *pbcm.IDRequest) (*api.PublicAccountInfo, error) {
	partner, err := sqlstore.Partner(ctx).ID(q.Id).Get()
	if err != nil {
		return nil, err
	}
	result := convertpb.PbPublicAccountInfo(partner)
	return result, nil
}

// GetPublicPartners is a little bit tricky. It handles 3 different cases:
// - get info by given ids
// - list all partners if the account is an admin
// - list all connected partner
func (s *AccountService) GetPublicPartners(ctx context.Context, q *pbcm.IDsRequest) (*api.PublicAccountsResponse, error) {
	if len(q.Ids) != 0 {
		partners, err := sqlstore.Partner(ctx).IDs(q.Ids...).IncludeDeleted().List()
		if err != nil {
			return nil, err
		}
		result := &api.PublicAccountsResponse{
			Accounts: convertpb.PbPublicPartners(partners),
		}
		return result, nil
	}

	accountIDs := s.SS.Claim().AccountIDs
	if isAdmin(accountIDs) {
		partners, err := sqlstore.Partner(ctx).IncludeDeleted().List()
		if err != nil {
			return nil, err
		}
		result := &api.PublicAccountsResponse{
			Accounts: convertpb.PbPublicPartners(partners),
		}
		return result, nil
	}

	listAccountIDs := make([]dot.ID, 0, len(accountIDs))
	for id := range accountIDs {
		listAccountIDs = append(listAccountIDs, id)
	}
	query := &identitymodelx.GetPartnersFromRelationQuery{
		AccountIDs: listAccountIDs,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.PublicAccountsResponse{
		Accounts: convertpb.PbPublicPartners(query.Result.Partners),
	}
	return result, nil
}

func isAdmin(accountIDs map[dot.ID]int) bool {
	for _, typ := range accountIDs {
		if typ == account_tag.TagEtop {
			return true
		}
	}
	return false
}
