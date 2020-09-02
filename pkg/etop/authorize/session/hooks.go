package session

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/capi/httprpc"
)

type Hook struct {
	perms  map[string]*permission.Decl
	secret string
}

func NewHook(perms map[string]*permission.Decl, opts ...HookOption) (*Hook, error) {
	h := &Hook{
		perms: perms,
	}
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func (h Hook) BuildHooks() httprpc.Hooks {
	return httprpc.Hooks{
		RequestRouted: h.BeforeServing,
	}
}

func (h Hook) BeforeServing(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
	perm, ok := h.perms[info.Route]
	if !ok {
		return ctx, cm.Errorf(cm.Internal, nil, "no permission declaration for route %v", info.Route)
	}
	if perm.Type == permission.Secret {

	}
	_auth, ok := info.Inner.(Sessioner)
	if ok {
		return _auth.StartSession(ctx, *perm, getToken(ctx))
	}
	if perm.Type != permission.Public {
		return ctx, cm.Errorf(cm.Internal, nil, "no session implementation for %v", info.Route)
	}
	return ctx, nil
}

func getToken(ctx context.Context) string {
	return headers.GetBearerTokenFromCtx(ctx)
}
