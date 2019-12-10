// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package api

import (
	"context"
	"time"

	api "etop.vn/api/top/services/handler"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	cmwrapper "etop.vn/backend/pkg/common/apifw/wrapper"
	bus "etop.vn/backend/pkg/common/bus"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapMiscService(s *MiscService, secret string) api.MiscService {
	return wrapMiscService{s: s, secret: secret}
}

type wrapMiscService struct {
	s      *MiscService
	secret string
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
	const rpcName = "handler.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context: ctx,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &VersionInfoEndpoint{Empty: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := middleware.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
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

func WrapWebhookService(s *WebhookService, secret string) api.WebhookService {
	return wrapWebhookService{s: s, secret: secret}
}

type wrapWebhookService struct {
	s      *WebhookService
	secret string
}

type ResetStateEndpoint struct {
	*api.ResetStateRequest
	Result  *cm.Empty
	Context claims.EmptyClaim
}

func (s wrapWebhookService) ResetState(ctx context.Context, req *api.ResetStateRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "handler.Webhook/ResetState"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context: ctx,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ResetStateEndpoint{ResetStateRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := middleware.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.ResetState(ctx, query)
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
