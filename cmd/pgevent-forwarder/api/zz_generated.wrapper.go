// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package api

import (
	"context"
	"time"

	cm "etop.vn/api/pb/common"
	pgevent "etop.vn/api/pb/services/pgevent"
	api "etop.vn/api/root/services/pgevent"
	common "etop.vn/backend/pkg/common"
	bus "etop.vn/backend/pkg/common/bus"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapEventService(s *EventService, secret string) api.EventService {
	return wrapEventService{s: s, secret: secret}
}

type wrapEventService struct {
	s      *EventService
	secret string
}

type GenerateEventsEndpoint struct {
	*pgevent.GenerateEventsRequest
	Result  *cm.Empty
	Context claims.EmptyClaim
}

func (s wrapEventService) GenerateEvents(ctx context.Context, req *pgevent.GenerateEventsRequest) (resp *cm.Empty, err error) {
	t0 := time.Now()
	var errs []*cm.Error
	const rpcName = "pgevent.Event/GenerateEvents"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &GenerateEventsEndpoint{GenerateEventsRequest: req}
	// Verify secret token
	token := middleware.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.GenerateEvents(ctx, query)
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
	const rpcName = "pgevent.Misc/VersionInfo"
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
