package admin

import (
	"context"

	"o.o/api/top/int/etop"
	"o.o/api/top/int/sadmin"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
)

type UserService struct {
	session.Sessioner
	ss *session.Session
}

func NewUserService(ss *session.Session) *UserService {
	return &UserService{
		ss: ss,
	}
}

func (s *UserService) Clone() sadmin.UserService {
	res := *s
	res.Sessioner, res.ss = s.ss.Split()
	return &res
}

func (s *UserService) CreateUser(ctx context.Context, r *sadmin.SAdminCreateUserRequest) (*etop.RegisterResponse, error) {
	r2 := &api.RegisterEndpoint{
		CreateUserRequest: r.Info,
	}
	if err := api.UserServiceImpl.Register(ctx, r2); err != nil {
		return r2.Result, err
	}

	if r.IsEtopAdmin {
		if r.Permission != nil {

		}
		roleCmd := &identitymodelx.UpdateRoleCommand{
			AccountID: model.EtopAccountID,
			UserID:    r2.Result.User.Id,
			Permission: identitymodel.Permission{
				Roles:       r.Permission.GetRoles(),
				Permissions: r.Permission.GetPermissions(),
			},
		}
		if err := sqlstore.UpdateRole(ctx, roleCmd); err != nil {
			return nil, err
		}
	}
	return r2.Result, nil
}

func (s *UserService) ResetPassword(ctx context.Context, r *sadmin.SAdminResetPasswordRequest) (*pbcm.Empty, error) {
	if len(r.Password) < 8 {
		return nil, cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.Password != r.Confirm {
		return nil, cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &identitymodelx.SetPasswordCommand{
		UserID:   r.UserId,
		Password: r.Password,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.Empty{}, nil
}

func (s *UserService) LoginAsAccount(ctx context.Context, r *sadmin.LoginAsAccountRequest) (*etop.LoginResponse, error) {
	resp, err := api.CreateLoginResponse(ctx, nil, "", r.UserId, nil, r.AccountId, 0, true, 0)
	return resp, err
}
