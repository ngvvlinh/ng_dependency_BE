// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package pgevent

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	pgevent "etop.vn/backend/pb/services/pgevent"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type EventServiceServer struct {
	EventAPI
}

func NewEventServiceServer(svc EventAPI) Server {
	return &EventServiceServer{
		EventAPI: svc,
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

func (s *EventServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/pgevent.Event/GenerateEvents":
		msg := new(pgevent.GenerateEventsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.EventAPI.GenerateEvents(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	MiscAPI
}

func NewMiscServiceServer(svc MiscAPI) Server {
	return &MiscServiceServer{
		MiscAPI: svc,
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

func (s *MiscServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/pgevent.Misc/VersionInfo":
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