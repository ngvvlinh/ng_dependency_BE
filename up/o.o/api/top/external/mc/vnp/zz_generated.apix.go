// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package vnp

import (
	context "context"
	fmt "fmt"
	http "net/http"

	externaltypes "o.o/api/top/external/types"
	common "o.o/api/top/types/common"
	capi "o.o/capi"
	httprpc "o.o/capi/httprpc"
)

func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() ShipnowService:
		return NewShipnowServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type ShipnowServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() ShipnowService
}

func NewShipnowServiceServer(builder func() ShipnowService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &ShipnowServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const ShipnowServicePathPrefix = "/vnposts/"

func (s *ShipnowServiceServer) PathPrefix() string {
	return ShipnowServicePathPrefix
}

func (s *ShipnowServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *ShipnowServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks)
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.RequestReceived(ctx, *info)
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

func (s *ShipnowServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/vnposts/cancelordervnpost":
		msg := &externaltypes.CancelShipnowFulfillmentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CancelShipnowFulfillment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/vnposts/createordervnpost":
		msg := &externaltypes.CreateShipnowFulfillmentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateShipnowFulfillment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/vnposts/getordervnpost":
		msg := &externaltypes.FulfillmentIDRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetShipnowFulfillment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/vnposts/getservicesvnpost":
		msg := &externaltypes.GetShipnowServicesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetShipnowServices(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/vnposts/ping":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Ping(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
