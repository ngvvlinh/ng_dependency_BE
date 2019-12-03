// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package pgevent

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "etop.vn/api/pb/common"
	pgevent "etop.vn/api/pb/services/pgevent"
	capi "etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type EventServiceServer struct {
	inner EventService
}

func NewEventServiceServer(svc EventService) Server {
	return &EventServiceServer{
		inner: svc,
	}
}

const EventServicePathPrefix = "/pgevent.Event/"

func (s *EventServiceServer) PathPrefix() string {
	return EventServicePathPrefix
}

func (s *EventServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *EventServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/pgevent.Event/GenerateEvents":
		msg := &pgevent.GenerateEventsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GenerateEvents(ctx, msg)
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

const MiscServicePathPrefix = "/pgevent.Misc/"

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
	case "/pgevent.Misc/VersionInfo":
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
