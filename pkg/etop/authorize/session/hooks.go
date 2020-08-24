package session

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/capi/httprpc"
)

type Hook struct {
	Permissions map[string]*permission.Decl
}

func NewHook(perms map[string]*permission.Decl) Hook {
	return Hook{
		Permissions: perms,
	}
}

func (h Hook) BuildHooks() httprpc.Hooks {
	return httprpc.Hooks{
		RequestRouted: h.BeforeServing,
	}
}

func (h Hook) BeforeServing(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
	perm, ok := h.Permissions[info.Route]
	if !ok {
		return ctx, cm.Errorf(cm.Internal, nil, "no permission declaration for route %v", info.Route)
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
