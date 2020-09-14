package session

import (
	"context"
	"reflect"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/capi"
	"o.o/capi/httprpc"
)

type Hook struct {
	perms  map[string]*permission.Decl
	secret string

	// TODO(vu): read permission decl and verify that all required information
	// are provided (secret, captcha, etc.)
	captcha *captcha.Captcha
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

func getCaptchaTokenFromRequest(req capi.Message) string {
	field := reflect.ValueOf(req).Elem().FieldByName("RecaptchaToken")
	if !field.IsValid() {
		ll.S.Panicf("no RecaptchaToken field in request %T", req)
	}
	token, ok := field.Interface().(string)
	if !ok {
		ll.S.Panicf("RecaptchaToken field in request %T is not string (got %T)", req, field.Interface())
	}
	return token
}

func (h Hook) BeforeServing(ctx context.Context, info httprpc.HookInfo) (context.Context, error) {
	perm, ok := h.perms[info.Route]
	if !ok {
		return ctx, cm.Errorf(cm.Internal, nil, "no permission declaration for route %v", info.Route)
	}
	if perm.Captcha != "" {
		err := h.captcha.Verify(getCaptchaTokenFromRequest(info.Request))
		if err != nil {
			return nil, err
		}
	}
	if perm.Type == permission.Secret {
		// TODO(vu): verify permission without requiring sessioner
		token := headers.GetBearerTokenFromCtx(ctx)
		if token != h.secret {
			return nil, cm.ErrUnauthenticated
		}
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
