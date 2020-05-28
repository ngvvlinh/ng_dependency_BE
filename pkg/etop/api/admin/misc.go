package admin

import (
	"context"

	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/authorize/login"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type MiscService struct{}

func (s *MiscService) Clone() *MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Admin",
		Version: "0.1",
	}
	return nil
}

func (s *MiscService) AdminLoginAsAccount(ctx context.Context, q *AdminLoginAsAccountEndpoint) error {
	loginQuery := &login.LoginUserQuery{
		UserID:   q.Context.UserID,
		Password: q.Password,
	}
	if err := bus.Dispatch(ctx, loginQuery); err != nil {
		return cm.MapError(err).
			Mapf(cm.Unauthenticated, cm.Unauthenticated, "Admin password: %v", err).
			DefaultInternal()
	}

	switch cm.GetTag(q.AccountId) {
	case model.TagShop:
	default:
		return cm.Error(cm.InvalidArgument, "Must be shop account", nil)
	}

	resp, err := s.adminCreateLoginResponse(ctx, q.Context.UserID, q.UserId, q.AccountId)
	q.Result = resp
	return err
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
