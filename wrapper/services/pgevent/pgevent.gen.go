// package pgevent generated by wrapper_gen. DO NOT EDIT.
package pgeventW

import (
	"context"
	"net/http"
	"time"

	twirp "github.com/twitchtv/twirp"

	cm "etop.vn/backend/pb/common"
	pgevent "etop.vn/backend/pb/services/pgevent"
	common "etop.vn/backend/pkg/common"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
	bus "etop.vn/common/bus"
	l "etop.vn/common/l"
)

var ll = l.New()

type (
	EmptyClaim   = claims.EmptyClaim
	UserClaim    = claims.UserClaim
	AdminClaim   = claims.AdminClaim
	PartnerClaim = claims.PartnerClaim
	ShopClaim    = claims.ShopClaim
)

type Muxer interface {
	Handle(string, http.Handler)
}

func NewPgeventServer(mux Muxer, hooks *twirp.ServerHooks, secret string) {
	if secret == "" {
		ll.Fatal("Secret is empty")
	}
	bus.Expect(&VersionInfoEndpoint{})
	bus.Expect(&GenerateEventsEndpoint{})
	mux.Handle(pgevent.MiscServicePathPrefix, pgevent.NewMiscServiceServer(MiscService{secret: secret}, hooks))
	mux.Handle(pgevent.EventServicePathPrefix, pgevent.NewEventServiceServer(EventService{secret: secret}, hooks))
}

type PgeventImpl struct {
	MiscService
	EventService
}

type MiscService struct{ secret string }

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context EmptyClaim
}

func (s MiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type EventService struct{ secret string }

type GenerateEventsEndpoint struct {
	*pgevent.GenerateEventsRequest
	Result  *cm.Empty
	Context EmptyClaim
}

func (s EventService) GenerateEvents(ctx context.Context, req *pgevent.GenerateEventsRequest) (resp *cm.Empty, err error) {
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}
