package admin

import (
	"context"

	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/authorize/login"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type MiscService struct {
	session.Session
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
	loginQuery := &login.LoginUserQuery{
		UserID:   s.SS.Claim().UserID,
		Password: q.Password,
	}
	if err := bus.Dispatch(ctx, loginQuery); err != nil {
		return nil, cm.MapError(err).
			Mapf(cm.Unauthenticated, cm.Unauthenticated, "Admin password: %v", err).
			DefaultInternal()
	}

	switch cm.GetTag(q.AccountId) {
	case model.TagShop:
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

	resp, err := api.CreateLoginResponse(ctx, nil, "", userID, nil, accountID, 0, false, adminID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
