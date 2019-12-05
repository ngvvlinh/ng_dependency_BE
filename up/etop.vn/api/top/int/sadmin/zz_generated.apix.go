// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package sadmin

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "etop.vn/api/top/types/common"
	capi "etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type MiscServiceServer struct {
	inner MiscService
}

func NewMiscServiceServer(svc MiscService) Server {
	return &MiscServiceServer{
		inner: svc,
	}
}

const MiscServicePathPrefix = "/sadmin.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *MiscServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/sadmin.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	inner UserService
}

func NewUserServiceServer(svc UserService) Server {
	return &UserServiceServer{
		inner: svc,
	}
}

const UserServicePathPrefix = "/sadmin.User/"

func (s *UserServiceServer) PathPrefix() string {
	return UserServicePathPrefix
}

func (s *UserServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *UserServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/sadmin.User/CreateUser":
		msg := &SAdminCreateUserRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateUser(ctx, msg)
		}
		return msg, fn, nil
	case "/sadmin.User/LoginAsAccount":
		msg := &LoginAsAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.LoginAsAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/sadmin.User/ResetPassword":
		msg := &SAdminResetPasswordRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ResetPassword(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
