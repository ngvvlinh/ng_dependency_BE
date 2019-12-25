// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package admin

import (
	"context"
	"time"

	etop "etop.vn/api/top/int/etop"
	api "etop.vn/api/top/int/sadmin"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	cmwrapper "etop.vn/backend/pkg/common/apifw/wrapper"
	bus "etop.vn/backend/pkg/common/bus"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapMiscService(s *MiscService) api.MiscService {
	return wrapMiscService{s: s}
}

type wrapMiscService struct {
	s *MiscService
}

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context claims.EmptyClaim
}

func (s wrapMiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &VersionInfoEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = s.s.VersionInfo(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

func WrapUserService(s *UserService) api.UserService {
	return wrapUserService{s: s}
}

type wrapUserService struct {
	s *UserService
}

type CreateUserEndpoint struct {
	*api.SAdminCreateUserRequest
	Result  *etop.RegisterResponse
	Context claims.EmptyClaim
}

func (s wrapUserService) CreateUser(ctx context.Context, req *api.SAdminCreateUserRequest) (resp *etop.RegisterResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/CreateUser"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateUserEndpoint{SAdminCreateUserRequest: req}
	query.Context.Claim = session.Claim
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateUser(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type LoginAsAccountEndpoint struct {
	*api.LoginAsAccountRequest
	Result  *etop.LoginResponse
	Context claims.EmptyClaim
}

func (s wrapUserService) LoginAsAccount(ctx context.Context, req *api.LoginAsAccountRequest) (resp *etop.LoginResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/LoginAsAccount"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &LoginAsAccountEndpoint{LoginAsAccountRequest: req}
	query.Context.Claim = session.Claim
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = s.s.LoginAsAccount(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type ResetPasswordEndpoint struct {
	*api.SAdminResetPasswordRequest
	Result  *cm.Empty
	Context claims.EmptyClaim
}

func (s wrapUserService) ResetPassword(ctx context.Context, req *api.SAdminResetPasswordRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "sadmin.User/ResetPassword"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:           ctx,
		RequireAuth:       true,
		RequireSuperAdmin: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ResetPasswordEndpoint{SAdminResetPasswordRequest: req}
	query.Context.Claim = session.Claim
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	ctx = bus.NewRootContext(ctx)
	err = s.s.ResetPassword(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}
