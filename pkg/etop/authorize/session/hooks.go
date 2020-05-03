package session

import (
	"context"
	"fmt"

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

func (h Hook) Build() *httprpc.Hooks {
	return &httprpc.Hooks{
		BeforeServing: h.BeforeServing,
	}
}

func (h Hook) BeforeServing(ctx context.Context, info httprpc.HookInfo, inner interface{}) (context.Context, error) {
	perm, ok := h.Permissions[info.Route]
	if !ok {
		return ctx, cm.Errorf(cm.Internal, nil, "no permission declaration for route %v", info.Route)
	}
	_auth, ok := inner.(Sessioner)
	if ok {
		return _auth.StartSession(ctx, *perm, getToken(ctx))
	}
	panic(fmt.Sprintf("%T must be an Authenticator", inner))
}

func getToken(ctx context.Context) string {
	return headers.GetBearerTokenFromCtx(ctx)
}
