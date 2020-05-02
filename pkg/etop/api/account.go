package api

import (
	"context"

	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		accountService.UpdateURLSlug,
		accountService.GetPublicPartnerInfo,
		accountService.GetPublicPartners,
	)
}

type AccountService struct{}

var accountService = &AccountService{}

func (s *AccountService) Clone() *AccountService {
	res := *s
	return &res
}

func (s *AccountService) UpdateURLSlug(ctx context.Context, r *UpdateURLSlugEndpoint) error {
	if r.AccountId == 0 {
		return cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}
	if !r.UrlSlug.Valid {
		return cm.Error(cm.InvalidArgument, "Missing url_slug", nil)
	}

	accQuery := &identitymodelx.GetAllAccountRolesQuery{UserID: r.Context.User.ID}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return err
	}

	var account *identitymodel.AccountUserExtended
	for _, acc := range accQuery.Result {
		if acc.Account.ID == r.AccountId {
			account = acc
			break
		}
	}
	if account == nil {
		return cm.Error(cm.InvalidArgument, "Tài khoản yêu cầu không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}

	cmd := &identitymodelx.UpdateAccountURLSlugCommand{
		AccountID: r.AccountId,
		URLSlug:   r.UrlSlug.String,
	}

	r.Result = &pbcm.Empty{}
	return bus.Dispatch(ctx, cmd)
}

func (s *UserService) UpdatePermission(ctx context.Context, q *UpdatePermissionEndpoint) error {
	return cm.ErrTODO
}

func (s *AccountService) GetPublicPartnerInfo(ctx context.Context, q *GetPublicPartnerInfoEndpoint) error {
	partner, err := sqlstore.Partner(ctx).ID(q.Id).Get()
	if err != nil {
		return err
	}
	q.Result = convertpb.PbPublicAccountInfo(partner)
	return nil
}

// GetPublicPartners is a little bit tricky. It handles 3 different cases:
// - get info by given ids
// - list all partners if the account is an admin
// - list all connected partner
func (s *AccountService) GetPublicPartners(ctx context.Context, q *GetPublicPartnersEndpoint) error {
	if len(q.Ids) != 0 {
		partners, err := sqlstore.Partner(ctx).IDs(q.Ids...).IncludeDeleted().List()
		if err != nil {
			return err
		}
		q.Result = &etop.PublicAccountsResponse{
			Accounts: convertpb.PbPublicPartners(partners),
		}
		return nil
	}

	accountIDs := q.Context.AccountIDs
	if isAdmin(accountIDs) {
		partners, err := sqlstore.Partner(ctx).IncludeDeleted().List()
		if err != nil {
			return err
		}
		q.Result = &etop.PublicAccountsResponse{
			Accounts: convertpb.PbPublicPartners(partners),
		}
		return nil
	}

	listAccountIDs := make([]dot.ID, 0, len(accountIDs))
	for id := range accountIDs {
		listAccountIDs = append(listAccountIDs, id)
	}
	query := &identitymodelx.GetPartnersFromRelationQuery{
		AccountIDs: listAccountIDs,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &etop.PublicAccountsResponse{
		Accounts: convertpb.PbPublicPartners(query.Result.Partners),
	}
	return nil
}

func isAdmin(accountIDs map[dot.ID]int) bool {
	for _, typ := range accountIDs {
		if typ == model.TagEtop {
			return true
		}
	}
	return false
}
