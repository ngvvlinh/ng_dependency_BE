package admin

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type AccountUserService struct {
	session.Session

	IdentityQuery identity.QueryBus
}

func (s *AccountUserService) Clone() admin.AccountUserService {
	res := *s
	return &res
}

func (s *AccountUserService) GetAccountUsers(ctx context.Context, r *admin.GetAccountUsersRequest) (*admin.GetAccountUsersResponse, error) {
	// Parse Paging
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}

	query := &identity.ListExtendedAccountUsersQuery{
		Paging: *paging,
	}
	if r.Filter != nil {
		query.FullNameNorm = r.Filter.Name
		query.PhoneNorm = r.Filter.Phone
		query.AccountID = r.Filter.AccountID
		query.Roles = r.Filter.Roles
		query.ExactRoles = r.Filter.ExactRoles
	}
	if err = s.IdentityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &admin.GetAccountUsersResponse{
		AccountUsers: convertpb.Convert_core_ExtendedAccountUsers_To_api_ExtendedAccountUsers(query.Result.AccountUsers),
		Paging:       cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return result, nil
}
