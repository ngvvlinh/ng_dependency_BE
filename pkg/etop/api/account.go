package api

import (
	"context"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrapetop "etop.vn/backend/wrapper/etop"
)

func init() {
	bus.AddHandlers("api",
		UpdateURLSlug,
		UpdatePermission,
		GetPublicPartnerInfo,
		GetPublicPartners,
	)
}

func UpdateURLSlug(ctx context.Context, r *wrapetop.UpdateURLSlugEndpoint) error {
	if r.AccountId == 0 {
		return cm.Error(cm.InvalidArgument, "Missing account_id", nil)
	}
	if r.UrlSlug == nil {
		return cm.Error(cm.InvalidArgument, "Missing url_slug", nil)
	}

	accQuery := &model.GetAllAccountRolesQuery{UserID: r.Context.User.ID}
	if err := bus.Dispatch(ctx, accQuery); err != nil {
		return err
	}

	var account *model.AccountUserExtended
	for _, acc := range accQuery.Result {
		if acc.Account.ID == r.AccountId {
			account = acc
			break
		}
	}
	if account == nil {
		return cm.Error(cm.InvalidArgument, "Tài khoản yêu cầu không hợp lệ. Vui lòng kiểm tra lại.", nil)
	}

	cmd := &model.UpdateAccountURLSlugCommand{
		AccountID: r.AccountId,
		URLSlug:   *r.UrlSlug,
	}

	r.Result = &pbcm.Empty{}
	return bus.Dispatch(ctx, cmd)
}

func UpdatePermission(ctx context.Context, q *wrapetop.UpdatePermissionEndpoint) error {
	return cm.ErrTODO
}

func GetPublicPartnerInfo(ctx context.Context, q *wrapetop.GetPublicPartnerInfoEndpoint) error {
	partner, err := sqlstore.Partner(ctx).ID(q.Id).Get()
	if err != nil {
		return err
	}
	q.Result = pbetop.PbPublicAccountInfo(partner)
	return nil
}

// GetPublicPartners is a little bit tricky. It handles 3 different cases:
// - get info by given ids
// - list all partners if the account is an admin
// - list all connected partner
func GetPublicPartners(ctx context.Context, q *wrapetop.GetPublicPartnersEndpoint) error {
	if len(q.Ids) != 0 {
		partners, err := sqlstore.Partner(ctx).IDs(q.Ids...).IncludeDeleted().List()
		if err != nil {
			return err
		}
		q.Result = &pbetop.PublicAccountsResponse{
			Accounts: pbetop.PbPublicPartners(partners),
		}
		return nil
	}

	accountIDs := q.Context.AccountIDs
	if isAdmin(accountIDs) {
		partners, err := sqlstore.Partner(ctx).IncludeDeleted().List()
		if err != nil {
			return err
		}
		q.Result = &pbetop.PublicAccountsResponse{
			Accounts: pbetop.PbPublicPartners(partners),
		}
		return nil
	}

	listAccountIDs := make([]int64, 0, len(accountIDs))
	for id := range accountIDs {
		listAccountIDs = append(listAccountIDs, id)
	}
	query := &model.GetPartnersFromRelationQuery{
		AccountIDs: listAccountIDs,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.PublicAccountsResponse{
		Accounts: pbetop.PbPublicPartners(query.Result.Partners),
	}
	return nil
}

func isAdmin(accountIDs map[int64]int) bool {
	for _, typ := range accountIDs {
		if typ == model.TagEtop {
			return true
		}
	}
	return false
}
