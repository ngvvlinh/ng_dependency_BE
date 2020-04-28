package admin

import (
	"context"

	pbcm "o.o/api/top/types/common"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/model"
)

var miscService = &MiscService{}
var userService = &UserService{}

func init() {
	bus.AddHandler("api", miscService.VersionInfo)
	bus.AddHandler("api", userService.CreateUser)
	bus.AddHandler("api", userService.ResetPassword)
	bus.AddHandler("api", userService.LoginAsAccount)
}

type MiscService struct{}
type UserService struct{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.SuperAdmin",
		Version: "0.1",
	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, r *CreateUserEndpoint) error {
	r2 := &api.RegisterEndpoint{
		CreateUserRequest: r.Info,
	}
	if err := bus.Dispatch(ctx, r2); err != nil {
		return err
	}
	r.Result = r2.Result

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
		if err := bus.Dispatch(ctx, roleCmd); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, r *ResetPasswordEndpoint) error {
	if len(r.Password) < 8 {
		return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.Password != r.Confirm {
		return cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &identitymodelx.SetPasswordCommand{
		UserID:   r.UserId,
		Password: r.Password,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.Empty{}
	return nil
}

func (s *UserService) LoginAsAccount(ctx context.Context, r *LoginAsAccountEndpoint) error {
	resp, err := api.CreateLoginResponse(ctx, nil, "", r.UserId, nil, r.AccountId, 0, true, 0)
	r.Result = resp
	return err
}
