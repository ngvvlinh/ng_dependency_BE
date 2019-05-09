package admin

import (
	"context"

	cmP "etop.vn/backend/pb/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/model"
	etopW "etop.vn/backend/wrapper/etop"
	sadminW "etop.vn/backend/wrapper/etop/sadmin"
)

func init() {
	bus.AddHandler("api", VersionInfo)
	bus.AddHandler("api", CreateUser)
	bus.AddHandler("api", ResetPassword)
	bus.AddHandler("api", LoginAsAccount)

	bus.Expect(&model.UpdateRoleCommand{})
	bus.Expect(&model.SetPasswordCommand{})
}

func VersionInfo(ctx context.Context, q *sadminW.VersionInfoEndpoint) error {
	q.Result = &cmP.VersionInfoResponse{
		Service: "etop.SuperAdmin",
		Version: "0.1",
	}
	return nil
}

func CreateUser(ctx context.Context, r *sadminW.CreateUserEndpoint) error {
	r2 := &etopW.RegisterEndpoint{
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

func ResetPassword(ctx context.Context, r *sadminW.ResetPasswordEndpoint) error {
	if len(r.Password) < 8 {
		return cm.Error(cm.InvalidArgument, "Password is too short", nil)
	}
	if r.Password != r.Confirm {
		return cm.Error(cm.InvalidArgument, "Password does not match", nil)
	}

	cmd := &model.SetPasswordCommand{
		UserID:   r.UserId,
		Password: r.Password,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &cmP.Empty{}
	return nil
}

func LoginAsAccount(ctx context.Context, r *sadminW.LoginAsAccountEndpoint) error {
	resp, err := api.CreateLoginResponse(ctx, nil, "", r.UserId, nil, r.AccountId, 0, true, 0)
	r.Result = resp
	return err
}
