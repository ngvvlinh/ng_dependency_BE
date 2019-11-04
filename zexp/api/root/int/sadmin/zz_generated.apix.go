// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package sadmin

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	sadmin "etop.vn/backend/pb/etop/sadmin"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type MiscServiceServer struct {
	MiscAPI
}

func NewMiscServiceServer(svc MiscAPI) Server {
	return &MiscServiceServer{
		MiscAPI: svc,
	}
}

const MiscServicePathPrefix = "/api/sadmin.Misc/"

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

func (s *MiscServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/sadmin.Misc/VersionInfo":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MiscAPI.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	UserAPI
}

func NewUserServiceServer(svc UserAPI) Server {
	return &UserServiceServer{
		UserAPI: svc,
	}
}

const UserServicePathPrefix = "/api/sadmin.User/"

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

func (s *UserServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/api/sadmin.User/CreateUser":
		msg := new(sadmin.SAdminCreateUserRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.CreateUser(ctx, msg)
		}
		return msg, fn, nil
	case "/api/sadmin.User/LoginAsAccount":
		msg := new(sadmin.LoginAsAccountRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.LoginAsAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/api/sadmin.User/ResetPassword":
		msg := new(sadmin.SAdminResetPasswordRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.ResetPassword(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
