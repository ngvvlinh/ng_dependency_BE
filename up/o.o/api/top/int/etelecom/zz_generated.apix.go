// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package etelecom

import (
	context "context"
	fmt "fmt"
	http "net/http"

	etelecomtypes "o.o/api/top/int/etelecom/types"
	common "o.o/api/top/types/common"
	capi "o.o/capi"
	httprpc "o.o/capi/httprpc"
)

func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() EtelecomService:
		return NewEtelecomServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type EtelecomServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() EtelecomService
}

func NewEtelecomServiceServer(builder func() EtelecomService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &EtelecomServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const EtelecomServicePathPrefix = "/shop.Etelecom/"

const Path_Etelecom_CreateCallLog = "/shop.Etelecom/CreateCallLog"
const Path_Etelecom_CreateExtension = "/shop.Etelecom/CreateExtension"
const Path_Etelecom_CreateUserAndAssignExtension = "/shop.Etelecom/CreateUserAndAssignExtension"
const Path_Etelecom_GetCallLogs = "/shop.Etelecom/GetCallLogs"
const Path_Etelecom_GetExtensions = "/shop.Etelecom/GetExtensions"
const Path_Etelecom_GetHotlines = "/shop.Etelecom/GetHotlines"
const Path_Etelecom_SummaryEtelecom = "/shop.Etelecom/SummaryEtelecom"

func (s *EtelecomServiceServer) PathPrefix() string {
	return EtelecomServicePathPrefix
}

func (s *EtelecomServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *EtelecomServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *EtelecomServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/shop.Etelecom/CreateCallLog":
		msg := &CreateCallLogRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateCallLog(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/CreateExtension":
		msg := &etelecomtypes.CreateExtensionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateExtension(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/CreateUserAndAssignExtension":
		msg := &CreateUserAndAssignExtensionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateUserAndAssignExtension(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/GetCallLogs":
		msg := &etelecomtypes.GetCallLogsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetCallLogs(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/GetExtensions":
		msg := &etelecomtypes.GetExtensionsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetExtensions(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/GetHotlines":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetHotlines(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/SummaryEtelecom":
		msg := &SummaryEtelecomRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SummaryEtelecom(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
