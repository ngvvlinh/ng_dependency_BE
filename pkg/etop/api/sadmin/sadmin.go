package admin

import (
	"context"

	pbcm "etop.vn/backend/pb/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/model"
	wrapetop "etop.vn/backend/wrapper/etop"
	wrapadmin "etop.vn/backend/wrapper/etop/sadmin"
)

func init() {
	bus.AddHandler("api", VersionInfo)
	bus.AddHandler("api", CreateUser)
	bus.AddHandler("api", ResetPassword)
	bus.AddHandler("api", LoginAsAccount)

	bus.Expect(&model.UpdateRoleCommand{})
	bus.Expect(&model.SetPasswordCommand{})
}

func VersionInfo(ctx context.Context, q *wrapadmin.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.SuperAdmin",
		Version: "0.1",
	}
	return nil
}

func CreateUser(ctx context.Context, r *wrapadmin.CreateUserEndpoint) error {
	r2 := &wrapetop.RegisterEndpoint{
		CreateUserRequest: r.Info,
	}
	if err := bus.Dispatch(ctx, r2); err != nil {
		return err
	}
	r.Result = r2.Result

	if r.IsEtopAdmin {
		if r.Permission != nil {

		}
		roleCmd := &model.UpdateRoleCommand{
			AccountID: model.EtopAccountID,
			UserID:    r2.Result.User.Id,
			Permission: model.Permission{
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

func ResetPassword(ctx context.Context, r *wrapadmin.ResetPasswordEndpoint) error {
	if len(r.Password) < 8 {
		return cm.Error(cm.InvalidArgument, "Mật khẩu phải có ít nhất 8 ký tự", nil)
	}
	if r.Password != r.Confirm {
		return cm.Error(cm.InvalidArgument, "Mật khẩu không khớp", nil)
	}

	cmd := &model.SetPasswordCommand{
		UserID:   r.UserId,
		Password: r.Password,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.Empty{}
	return nil
}

func LoginAsAccount(ctx context.Context, r *wrapadmin.LoginAsAccountEndpoint) error {
	resp, err := api.CreateLoginResponse(ctx, nil, "", r.UserId, nil, r.AccountId, 0, true, 0)
	r.Result = resp
	return err
}
