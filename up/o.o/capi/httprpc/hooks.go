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
	Inner       interface{}
}

type Hooks struct {
	BeforeRequest  func(ctx context.Context, info HookInfo) (context.Context, error)
	BeforeServing  func(ctx context.Context, info HookInfo) (context.Context, error)
	BeforeResponse func(ctx context.Context, info HookInfo, respHeaders http.Header) (context.Context, error)
	AfterResponse  func(ctx context.Context, info HookInfo)
	ErrorServing   func(ctx context.Context, info HookInfo, err error) (context.Context, error)
}

type HooksBuilder interface {
	BuildHooks() Hooks
}

type HooksFunc func() Hooks

func (h HooksFunc) BuildHooks() Hooks { return h() }

type chainHooks []HooksBuilder

func (s chainHooks) BuildHooks() Hooks {
	switch len(s) {
	case 0:
		return Hooks{}
	case 1:
		return s[0].BuildHooks()
	}
	hooks := make([]Hooks, len(s))
	for i, b := range s {
		hooks[i] = b.BuildHooks()
	}
	return Hooks{
		BeforeRequest: func(ctx context.Context, info HookInfo) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.BeforeRequest != nil {
					ctx, err = h.BeforeRequest(ctx, info)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		BeforeServing: func(ctx context.Context, info HookInfo) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.BeforeServing != nil {
					ctx, err = h.BeforeServing(ctx, info)
					if err != nil {
						return ctx, err
					}
				}
			}
			return ctx, nil
		},
		BeforeResponse: func(ctx context.Context, info HookInfo, respHeaders http.Header) (_ context.Context, err error) {
			for _, h := range hooks {
				if h.BeforeResponse != nil {
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
				if h.AfterResponse != nil {
					h.AfterResponse(ctx, info)
				}
			}
		},
		ErrorServing: func(ctx context.Context, info HookInfo, err error) (context.Context, error) {
			for _, h := range hooks {
				if h.ErrorServing != nil {
					ctx, err = h.ErrorServing(ctx, info, err)
				}
			}
			return ctx, err
		},
	}
}

func ChainHooks(hooks ...HooksBuilder) HooksBuilder {
	length := 2 * len(hooks)
	if length < 8 {
		length = 8
	}
	res := make(chainHooks, 0, length)
	for _, h := range hooks {
		if h == nil {
			continue
		}
		if hs, ok := h.(chainHooks); ok {
			res = append(res, hs...)
		} else {
			res = append(res, h)
		}
	}
	switch len(res) {
	case 0:
		return chainHooks{}
	case 1:
		return res[0]
	default:
		return res
	}
}

func WrapHooks(builder HooksBuilder) (res Hooks) {
	var hooks Hooks
	if builder != nil {
		hooks = builder.BuildHooks()
	}
	if hooks.BeforeRequest == nil {
		hooks.BeforeRequest = func(ctx context.Context, _ HookInfo) (context.Context, error) { return ctx, nil }
	}
	if hooks.BeforeServing == nil {
		hooks.BeforeServing = func(ctx context.Context, _ HookInfo) (context.Context, error) { return ctx, nil }
	}
	if hooks.BeforeResponse == nil {
		hooks.BeforeResponse = func(ctx context.Context, _ HookInfo, _ http.Header) (context.Context, error) { return ctx, nil }
	}
	if hooks.AfterResponse == nil {
		hooks.AfterResponse = func(ctx context.Context, _ HookInfo) {}
	}
	if hooks.ErrorServing == nil {
		hooks.ErrorServing = func(ctx context.Context, _ HookInfo, err error) (context.Context, error) { return ctx, err }
	}
	return hooks
}

func WithHooks(servers []Server, hooks ...HooksBuilder) []Server {
	if len(hooks) == 0 {
		return servers
	}
	result := make([]Server, len(servers))
	for i, s := range servers {
		result[i] = s.WithHooks(ChainHooks(hooks...))
	}
	return result
}
