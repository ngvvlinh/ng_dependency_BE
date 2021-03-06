package admin

import (
	"context"

	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	cm "o.o/backend/pkg/common"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
)

type MiscService struct {
	session.Session

	Login sqlstore.LoginInterface
}

func (s *MiscService) Clone() admin.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop.Admin",
		Version: "0.1",
	}
	return result, nil
}

func (s *MiscService) AdminLoginAsAccount(ctx context.Context, q *admin.LoginAsAccountRequest) (*etop.LoginResponse, error) {
	loginQuery := &sqlstore.LoginUserQuery{
		UserID:   s.SS.Claim().UserID,
		Password: q.Password,
	}
	if err := s.Login.LoginUser(ctx, loginQuery); err != nil {
		return nil, cm.MapError(err).
			Mapf(cm.Unauthenticated, cm.Unauthenticated, "Admin password: %v", err).
			DefaultInternal()
	}

	switch cm.GetTag(q.AccountId) {
	case account_tag.TagShop:
	default:
		return nil, cm.Error(cm.InvalidArgument, "Must be shop account", nil)
	}

	resp, err := s.adminCreateLoginResponse(ctx, s.SS.Claim().UserID, q.UserId, q.AccountId)
	result := resp
	return result, err
}

func (s *MiscService) adminCreateLoginResponse(ctx context.Context, adminID, userID, accountID dot.ID) (*etop.LoginResponse, error) {
	if adminID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing AdminID", nil)
	}

	resp, err := apiroot.CreateLoginResponse(ctx, nil, "", userID, nil, accountID, 0, false, adminID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
