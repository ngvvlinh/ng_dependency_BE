// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package integration

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "o.o/api/top/types/common"
	capi "o.o/capi"
	httprpc "o.o/capi/httprpc"
)

func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() IntegrationService:
		return NewIntegrationServiceServer(builder, hooks...), true
	case func() MiscService:
		return NewMiscServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type IntegrationServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() IntegrationService
}

func NewIntegrationServiceServer(builder func() IntegrationService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &IntegrationServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const IntegrationServicePathPrefix = "/integration.Integration/"

func (s *IntegrationServiceServer) PathPrefix() string {
	return IntegrationServicePathPrefix
}

func (s *IntegrationServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *IntegrationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks)
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.BeforeRequest(ctx, *info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path, hooks, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve(ctx, resp, req, hooks, info, reqMsg, exec)
}

func (s *IntegrationServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/integration.Integration/GrantAccess":
		msg := &GrantAccessRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GrantAccess(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/Init":
		msg := &InitRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Init(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/LoginUsingToken":
		msg := &LoginUsingTokenRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.LoginUsingToken(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/LoginUsingTokenWL":
		msg := &LoginUsingTokenRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.LoginUsingTokenWL(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/Register":
		msg := &RegisterRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Register(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/RequestLogin":
		msg := &RequestLoginRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.RequestLogin(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/integration.Integration/SessionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SessionInfo(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() MiscService
}

func NewMiscServiceServer(builder func() MiscService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &MiscServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const MiscServicePathPrefix = "/integration.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks)
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.BeforeRequest(ctx, *info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path, hooks, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve(ctx, resp, req, hooks, info, reqMsg, exec)
}

func (s *MiscServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/integration.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.VersionInfo(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
