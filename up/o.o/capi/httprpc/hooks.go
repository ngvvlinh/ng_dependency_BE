package httprpc

import (
	"context"
	"net/http"

	"o.o/capi"
)

type HookInfo struct {
	Route       string
	HTTPRequest *http.Request
	Request     capi.Message
	Response    capi.Message
}

type Hooks struct {
	BeforeRequest  func(ctx context.Context, info HookInfo) (context.Context, error)
	BeforeServing  func(ctx context.Context, info HookInfo, inner interface{}) (context.Context, error)
	BeforeResponse func(ctx context.Context, info HookInfo, respHeaders http.Header) (context.Context, error)
	AfterResponse  func(ctx context.Context, info HookInfo)
	ErrorServing   func(ctx context.Context, info HookInfo, err error) context.Context
}

func ChainHooks(hooks ...*Hooks) *Hooks {
	if len(hooks) == 1 {
		return hooks[0]
	}
	return &Hooks{
		BeforeRequest: func(ctx context.Context, info HookInfo) (_ context.Context, err error) {
			for _, h := range hooks {
				if h != nil && h.BeforeRequest != nil {
					ctx, err = h.BeforeRequest(ctx, info)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		BeforeServing: func(ctx context.Context, info HookInfo, inner interface{}) (_ context.Context, err error) {
			for _, h := range hooks {
				if h != nil && h.BeforeServing != nil {
					ctx, err = h.BeforeServing(ctx, info, inner)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		BeforeResponse: func(ctx context.Context, info HookInfo, respHeaders http.Header) (_ context.Context, err error) {
			for _, h := range hooks {
				if h != nil && h.BeforeResponse != nil {
					ctx, err = h.BeforeResponse(ctx, info, respHeaders)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		AfterResponse: func(ctx context.Context, info HookInfo) {
			for _, h := range hooks {
				if h != nil && h.AfterResponse != nil {
					h.AfterResponse(ctx, info)
				}
			}
		},
		ErrorServing: func(ctx context.Context, info HookInfo, err error) context.Context {
			for _, h := range hooks {
				if h != nil && h.ErrorServing != nil {
					ctx = h.ErrorServing(ctx, info, err)
				}
			}
			return ctx
		},
	}
}

func WrapHooks(hooks *Hooks) (res Hooks) {
	if hooks == nil {
		res = Hooks{}
	} else {
		res = *hooks
	}
	if res.BeforeRequest == nil {
		res.BeforeRequest = func(ctx context.Context, _ HookInfo) (context.Context, error) { return ctx, nil }
	}
	if res.BeforeServing == nil {
		res.BeforeServing = func(ctx context.Context, _ HookInfo, _ interface{}) (context.Context, error) { return ctx, nil }
	}
	if res.BeforeResponse == nil {
		res.BeforeResponse = func(ctx context.Context, _ HookInfo, _ http.Header) (context.Context, error) { return ctx, nil }
	}
	if res.AfterResponse == nil {
		res.AfterResponse = func(ctx context.Context, _ HookInfo) {}
	}
	if res.ErrorServing == nil {
		res.ErrorServing = func(ctx context.Context, _ HookInfo, _ error) context.Context { return ctx }
	}
	return res
}
