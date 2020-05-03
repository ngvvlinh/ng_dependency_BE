// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package sadmin

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

func NewServer(builder interface{}, hooks ...*httprpc.Hooks) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() MiscService:
		return NewMiscServiceServer(builder, hooks...), true
	case func() UserService:
		return NewUserServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type MiscServiceServer struct {
	hooks   httprpc.Hooks
	builder func() MiscService
}

func NewMiscServiceServer(builder func() MiscService, hooks ...*httprpc.Hooks) httprpc.Server {
	return &MiscServiceServer{
		hooks:   httprpc.WrapHooks(httprpc.ChainHooks(hooks...)),
		builder: builder,
	}
}

const MiscServicePathPrefix = "/sadmin.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx, info := req.Context(), httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := s.hooks.BeforeRequest(ctx, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve(ctx, resp, req, s.hooks, info, reqMsg, exec)
}

func (s *MiscServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/sadmin.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			ctx, err := s.hooks.BeforeServing(ctx, httprpc.HookInfo{Route: path, Request: msg}, inner)
			if err != nil {
				return nil, err
			}
			return inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	hooks   httprpc.Hooks
	builder func() UserService
}

func NewUserServiceServer(builder func() UserService, hooks ...*httprpc.Hooks) httprpc.Server {
	return &UserServiceServer{
		hooks:   httprpc.WrapHooks(httprpc.ChainHooks(hooks...)),
		builder: builder,
	}
}

const UserServicePathPrefix = "/sadmin.User/"

func (s *UserServiceServer) PathPrefix() string {
	return UserServicePathPrefix
}

func (s *UserServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx, info := req.Context(), httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := s.hooks.BeforeRequest(ctx, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve(ctx, resp, req, s.hooks, info, reqMsg, exec)
}

func (s *UserServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/sadmin.User/CreateUser":
		msg := &SAdminCreateUserRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			ctx, err := s.hooks.BeforeServing(ctx, httprpc.HookInfo{Route: path, Request: msg}, inner)
			if err != nil {
				return nil, err
			}
			return inner.CreateUser(ctx, msg)
		}
		return msg, fn, nil
	case "/sadmin.User/LoginAsAccount":
		msg := &LoginAsAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			ctx, err := s.hooks.BeforeServing(ctx, httprpc.HookInfo{Route: path, Request: msg}, inner)
			if err != nil {
				return nil, err
			}
			return inner.LoginAsAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/sadmin.User/ResetPassword":
		msg := &SAdminResetPasswordRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			ctx, err := s.hooks.BeforeServing(ctx, httprpc.HookInfo{Route: path, Request: msg}, inner)
			if err != nil {
				return nil, err
			}
			return inner.ResetPassword(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
