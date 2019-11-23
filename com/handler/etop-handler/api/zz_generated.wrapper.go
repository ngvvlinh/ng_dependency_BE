// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package api

import (
	"context"
	"time"

	cm "etop.vn/backend/pb/common"
	handler "etop.vn/backend/pb/services/handler"
	common "etop.vn/backend/pkg/common"
	bus "etop.vn/backend/pkg/common/bus"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
	api "etop.vn/backend/zexp/api/root/services/handler"
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
	var errs []*cm.Error
	const rpcName = "handler.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &VersionInfoEndpoint{Empty: req}
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
	*handler.ResetStateRequest
	Result  *cm.Empty
	Context claims.EmptyClaim
}

func (s wrapWebhookService) ResetState(ctx context.Context, req *handler.ResetStateRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var errs []*cm.Error
	const rpcName = "handler.Webhook/ResetState"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &ResetStateEndpoint{ResetStateRequest: req}
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
