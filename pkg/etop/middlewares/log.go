package middlewares

import (
	"context"
	"time"

	cmwrapper "o.o/backend/pkg/common/apifw/wrapper"
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
		BeforeServing: clone.beforeServing,
		ErrorServing:  clone.errorServing,
	}
}

func (log *Logging) beforeServing(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
	log.t0 = time.Now()
	return ctx, nil
}

func (log *Logging) errorServing(ctx context.Context, info httprpc.HookInfo, err error) (context.Context, error) {
	var ss *session.Session
	if _ss, ok := info.Inner.(session.Sessioner); ok {
		ss = _ss.GetSession()
	}
	err = cmwrapper.RecoverAndLog2(ctx, info.Route, ss, info.Request, info.Response, nil, err, cmwrapper.HasErrors(info.Response), log.t0)
	return ctx, err
}
