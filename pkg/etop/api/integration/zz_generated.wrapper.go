// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package integration

import (
	"context"
	"time"

	cm "etop.vn/backend/pb/common"
	api "etop.vn/backend/pb/etop/integration"
	common "etop.vn/backend/pkg/common"
	bus "etop.vn/backend/pkg/common/bus"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
	model "etop.vn/backend/pkg/etop/model"
)

func WrapIntegrationService(s *IntegrationService) api.IntegrationService {
	return wrapIntegrationService{s: s}
}

type wrapIntegrationService struct {
	s *IntegrationService
}

type GrantAccessEndpoint struct {
	*api.GrantAccessRequest
	Result     *api.GrantAccessResponse
	Context    claims.UserClaim
	CtxPartner *model.Partner
}

func (s wrapIntegrationService) GrantAccess(ctx context.Context, req *api.GrantAccessRequest) (resp *api.GrantAccessResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "integration.Integration/GrantAccess"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireUser: true,
		AuthPartner: 2,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GrantAccessEndpoint{GrantAccessRequest: req}
	query.Context.Claim = session.Claim
	query.Context.User = session.User
	query.Context.Admin = session.Admin
	query.CtxPartner = session.CtxPartner
	ctx = bus.NewRootContext(ctx)
	err = s.s.GrantAccess(ctx, query)
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

type InitEndpoint struct {
	*api.InitRequest
	Result  *api.LoginResponse
	Context claims.EmptyClaim
}

func (s wrapIntegrationService) Init(ctx context.Context, req *api.InitRequest) (resp *api.LoginResponse, err error) {
	t0 := time.Now()
	var errs []*cm.Error
	const rpcName = "integration.Integration/Init"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &InitEndpoint{InitRequest: req}
	ctx = bus.NewRootContext(ctx)
	err = s.s.Init(ctx, query)
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

type LoginUsingTokenEndpoint struct {
	*api.LoginUsingTokenRequest
	Result     *api.LoginResponse
	Context    claims.EmptyClaim
	CtxPartner *model.Partner
}

func (s wrapIntegrationService) LoginUsingToken(ctx context.Context, req *api.LoginUsingTokenRequest) (resp *api.LoginResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "integration.Integration/LoginUsingToken"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		AuthPartner: 2,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &LoginUsingTokenEndpoint{LoginUsingTokenRequest: req}
	query.Context.Claim = session.Claim
	query.CtxPartner = session.CtxPartner
	ctx = bus.NewRootContext(ctx)
	err = s.s.LoginUsingToken(ctx, query)
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

type RegisterEndpoint struct {
	*api.RegisterRequest
	Result     *api.RegisterResponse
	Context    claims.EmptyClaim
	CtxPartner *model.Partner
}

func (s wrapIntegrationService) Register(ctx context.Context, req *api.RegisterRequest) (resp *api.RegisterResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "integration.Integration/Register"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		AuthPartner: 2,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &RegisterEndpoint{RegisterRequest: req}
	query.Context.Claim = session.Claim
	query.CtxPartner = session.CtxPartner
	ctx = bus.NewRootContext(ctx)
	err = s.s.Register(ctx, query)
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

type RequestLoginEndpoint struct {
	*api.RequestLoginRequest
	Result     *api.RequestLoginResponse
	Context    claims.EmptyClaim
	CtxPartner *model.Partner
}

func (s wrapIntegrationService) RequestLogin(ctx context.Context, req *api.RequestLoginRequest) (resp *api.RequestLoginResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "integration.Integration/RequestLogin"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		AuthPartner: 2,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &RequestLoginEndpoint{RequestLoginRequest: req}
	query.Context.Claim = session.Claim
	query.CtxPartner = session.CtxPartner
	// Verify captcha token
	if err := middleware.VerifyCaptcha(ctx, req.RecaptchaToken); err != nil {
		return nil, err
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.RequestLogin(ctx, query)
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

type SessionInfoEndpoint struct {
	*cm.Empty
	Result     *api.LoginResponse
	Context    claims.EmptyClaim
	CtxPartner *model.Partner
}

func (s wrapIntegrationService) SessionInfo(ctx context.Context, req *cm.Empty) (resp *api.LoginResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "integration.Integration/SessionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		AuthPartner: 2,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &SessionInfoEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.CtxPartner = session.CtxPartner
	ctx = bus.NewRootContext(ctx)
	err = s.s.SessionInfo(ctx, query)
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
	var errs []*cm.Error
	const rpcName = "integration.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &VersionInfoEndpoint{Empty: req}
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
