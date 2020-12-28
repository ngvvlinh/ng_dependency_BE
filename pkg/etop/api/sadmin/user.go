package sadmin

import (
	"context"

	"o.o/api/top/int/etop"
	"o.o/api/top/int/sadmin"
	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etc/idutil"
	apiroot "o.o/backend/pkg/etop/api/root"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
)

type UserService struct {
	session.Session

	AccountUserStore *sqlstore.AccountUserStore
	UserStore        sqlstore.UserStoreInterface
}

func (s *UserService) Clone() sadmin.UserService {
	res := *s
	return &res
}

func (s *UserService) CreateUser(ctx context.Context, r *sadmin.SAdminCreateUserRequest) (*etop.RegisterResponse, error) {
	resp, err := apiroot.UserServiceImpl.Register(ctx, r.Info)
	if err != nil {
		return nil, err
	}
	if r.IsEtopAdmin {
		if r.Permission != nil {

		}
		roleCmd := &identitymodelx.UpdateRoleCommand{
			AccountID: idutil.EtopAccountID,
			UserID:    resp.User.Id,
			Permission: identitymodel.Permission{
				Roles:       r.Permission.GetRoles(),
				Permissions: r.Permission.GetPermissions(),
			},
		}
		if err := s.AccountUserStore.UpdateRole(ctx, roleCmd); err != nil {
			return nil, err
		}
	}
	return resp, nil
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
	if err := s.UserStore.SetPassword(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.Empty{}, nil
}

func (s *UserService) LoginAsAccount(ctx context.Context, r *sadmin.LoginAsAccountRequest) (*etop.LoginResponse, error) {
	resp, err := apiroot.CreateLoginResponse(ctx, nil, "", r.UserId, nil, r.AccountId, 0, true, 0)
	return resp, err
}
