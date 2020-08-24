package middlewares

import (
	"context"
	"net/http"
	"time"

	cmwrapper "o.o/backend/pkg/common/apifw/wrapper"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

type Logging struct {
	t0 time.Time
}

func NewLogging() *Logging {
	log := &Logging{}
	return log
}

func (log *Logging) BuildHooks() httprpc.Hooks {
	clone := *log
	return httprpc.Hooks{
		RequestRouted:    clone.requestRouted,
		ResponsePrepared: clone.responsePrepared,
		Error:            clone.error,
	}
}

func (log *Logging) requestRouted(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
	log.t0 = time.Now()
	return ctx, nil
}

func (log *Logging) responsePrepared(ctx context.Context, info httprpc.HookInfo, _ http.Header) (context.Context, error) {
	metrics.APIRequest(info.Route, time.Now().Sub(log.t0), nil)
	return ctx, nil
}

func (log *Logging) error(ctx context.Context, info httprpc.HookInfo, err error) (context.Context, error) {
	var ss *session.Session
	if _ss, ok := info.Inner.(session.Sessioner); ok {
		ss = _ss.GetSession()
	}
	cmwrapper.Censor(info.Request)
	err = cmwrapper.RecoverAndLog2(ctx, info.Route, ss, info.Request, info.Response, nil, err, cmwrapper.HasErrors(info.Response), log.t0)
	return ctx, err
}
