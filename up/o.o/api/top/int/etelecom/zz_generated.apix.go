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
	case func() EtelecomUserService:
		return NewEtelecomUserServiceServer(builder, hooks...), true
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

const Path_Etelecom_AssignUserToExtension = "/shop.Etelecom/AssignUserToExtension"
const Path_Etelecom_CreateCallLog = "/shop.Etelecom/CreateCallLog"
const Path_Etelecom_CreateExtension = "/shop.Etelecom/CreateExtension"
const Path_Etelecom_CreateExtensionBySubscription = "/shop.Etelecom/CreateExtensionBySubscription"
const Path_Etelecom_CreateTenant = "/shop.Etelecom/CreateTenant"
const Path_Etelecom_CreateUserAndAssignExtension = "/shop.Etelecom/CreateUserAndAssignExtension"
const Path_Etelecom_ExtendExtension = "/shop.Etelecom/ExtendExtension"
const Path_Etelecom_GetCallLogs = "/shop.Etelecom/GetCallLogs"
const Path_Etelecom_GetExtensions = "/shop.Etelecom/GetExtensions"
const Path_Etelecom_GetHotlines = "/shop.Etelecom/GetHotlines"
const Path_Etelecom_GetTenant = "/shop.Etelecom/GetTenant"
const Path_Etelecom_RemoveUserOfExtension = "/shop.Etelecom/RemoveUserOfExtension"
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
	case "/shop.Etelecom/AssignUserToExtension":
		msg := &AssignUserToExtensionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.AssignUserToExtension(newCtx, msg)
			return
		}
		return msg, fn, nil
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
	case "/shop.Etelecom/CreateExtensionBySubscription":
		msg := &etelecomtypes.CreateExtensionBySubscriptionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateExtensionBySubscription(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/CreateTenant":
		msg := &CreateTenantRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateTenant(newCtx, msg)
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
	case "/shop.Etelecom/ExtendExtension":
		msg := &etelecomtypes.ExtendExtensionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ExtendExtension(newCtx, msg)
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
		msg := &etelecomtypes.GetHotLinesRequest{}
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
	case "/shop.Etelecom/GetTenant":
		msg := &GetTenantRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetTenant(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/shop.Etelecom/RemoveUserOfExtension":
		msg := &RemoveUserOfExtensionRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.RemoveUserOfExtension(newCtx, msg)
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

type EtelecomUserServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() EtelecomUserService
}

func NewEtelecomUserServiceServer(builder func() EtelecomUserService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &EtelecomUserServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const EtelecomUserServicePathPrefix = "/etelecom.User/"

const Path_EtelecomUser_GetUserSetting = "/etelecom.User/GetUserSetting"

func (s *EtelecomUserServiceServer) PathPrefix() string {
	return EtelecomUserServicePathPrefix
}

func (s *EtelecomUserServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *EtelecomUserServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *EtelecomUserServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/etelecom.User/GetUserSetting":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetUserSetting(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
