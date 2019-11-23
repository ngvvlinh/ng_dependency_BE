// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package handler

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/api/pb/common"
	handler "etop.vn/api/pb/services/handler"
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

const MiscServicePathPrefix = "/handler.Misc/"

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
	case "/handler.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type WebhookServiceServer struct {
	inner WebhookService
}

func NewWebhookServiceServer(svc WebhookService) Server {
	return &WebhookServiceServer{
		inner: svc,
	}
}

const WebhookServicePathPrefix = "/handler.Webhook/"

func (s *WebhookServiceServer) PathPrefix() string {
	return WebhookServicePathPrefix
}

func (s *WebhookServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *WebhookServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/handler.Webhook/ResetState":
		msg := &handler.ResetStateRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.inner.ResetState(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}