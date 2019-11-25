// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package integration

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "etop.vn/api/pb/common"
	integration "etop.vn/api/pb/etop/integration"
	"etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type IntegrationServiceServer struct {
	inner IntegrationService
}

func NewIntegrationServiceServer(svc IntegrationService) Server {
	return &IntegrationServiceServer{
		inner: svc,
	}
}

const IntegrationServicePathPrefix = "/integration.Integration/"

func (s *IntegrationServiceServer) PathPrefix() string {
	return IntegrationServicePathPrefix
}

func (s *IntegrationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *IntegrationServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/integration.Integration/GrantAccess":
		msg := &integration.GrantAccessRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GrantAccess(ctx, msg)
		}
		return msg, fn, nil
	case "/integration.Integration/Init":
		msg := &integration.InitRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.Init(ctx, msg)
		}
		return msg, fn, nil
	case "/integration.Integration/LoginUsingToken":
		msg := &integration.LoginUsingTokenRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.LoginUsingToken(ctx, msg)
		}
		return msg, fn, nil
	case "/integration.Integration/Register":
		msg := &integration.RegisterRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.Register(ctx, msg)
		}
		return msg, fn, nil
	case "/integration.Integration/RequestLogin":
		msg := &integration.RequestLoginRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.RequestLogin(ctx, msg)
		}
		return msg, fn, nil
	case "/integration.Integration/SessionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.SessionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	inner MiscService
}

func NewMiscServiceServer(svc MiscService) Server {
	return &MiscServiceServer{
		inner: svc,
	}
}

const MiscServicePathPrefix = "/integration.Misc/"

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
	case "/integration.Misc/VersionInfo":
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