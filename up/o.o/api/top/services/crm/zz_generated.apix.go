// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package crm

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

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() CrmService:
		return NewCrmServiceServer(builder, hooks...), true
	case func() MiscService:
		return NewMiscServiceServer(builder, hooks...), true
	case func() VhtService:
		return NewVhtServiceServer(builder, hooks...), true
	case func() VtigerService:
		return NewVtigerServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type CrmServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() CrmService
}

func NewCrmServiceServer(builder func() CrmService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &CrmServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const CrmServicePathPrefix = "/crm.Crm/"

const Path_Crm_RefreshFulfillmentFromCarrier = "/crm.Crm/RefreshFulfillmentFromCarrier"
const Path_Crm_SendNotification = "/crm.Crm/SendNotification"

func (s *CrmServiceServer) PathPrefix() string {
	return CrmServicePathPrefix
}

func (s *CrmServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CrmServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CrmServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/crm.Crm/RefreshFulfillmentFromCarrier":
		msg := &RefreshFulfillmentFromCarrierRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.RefreshFulfillmentFromCarrier(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Crm/SendNotification":
		msg := &SendNotificationRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendNotification(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() MiscService
}

func NewMiscServiceServer(builder func() MiscService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &MiscServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const MiscServicePathPrefix = "/crm.Misc/"

const Path_Misc_VersionInfo = "/crm.Misc/VersionInfo"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *MiscServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/crm.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.VersionInfo(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type VhtServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() VhtService
}

func NewVhtServiceServer(builder func() VhtService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &VhtServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const VhtServicePathPrefix = "/crm.Vht/"

const Path_Vht_CreateOrUpdateCallHistoryByCallID = "/crm.Vht/CreateOrUpdateCallHistoryByCallID"
const Path_Vht_CreateOrUpdateCallHistoryBySDKCallID = "/crm.Vht/CreateOrUpdateCallHistoryBySDKCallID"
const Path_Vht_GetCallHistories = "/crm.Vht/GetCallHistories"

func (s *VhtServiceServer) PathPrefix() string {
	return VhtServicePathPrefix
}

func (s *VhtServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *VhtServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *VhtServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/crm.Vht/CreateOrUpdateCallHistoryByCallID":
		msg := &VHTCallLog{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateOrUpdateCallHistoryByCallID(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vht/CreateOrUpdateCallHistoryBySDKCallID":
		msg := &VHTCallLog{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateOrUpdateCallHistoryBySDKCallID(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vht/GetCallHistories":
		msg := &GetCallHistoriesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetCallHistories(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type VtigerServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() VtigerService
}

func NewVtigerServiceServer(builder func() VtigerService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &VtigerServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const VtigerServicePathPrefix = "/crm.Vtiger/"

const Path_Vtiger_CreateOrUpdateContact = "/crm.Vtiger/CreateOrUpdateContact"
const Path_Vtiger_CreateOrUpdateLead = "/crm.Vtiger/CreateOrUpdateLead"
const Path_Vtiger_CreateTicket = "/crm.Vtiger/CreateTicket"
const Path_Vtiger_GetCategories = "/crm.Vtiger/GetCategories"
const Path_Vtiger_GetContacts = "/crm.Vtiger/GetContacts"
const Path_Vtiger_GetTicketStatusCount = "/crm.Vtiger/GetTicketStatusCount"
const Path_Vtiger_GetTickets = "/crm.Vtiger/GetTickets"
const Path_Vtiger_UpdateTicket = "/crm.Vtiger/UpdateTicket"

func (s *VtigerServiceServer) PathPrefix() string {
	return VtigerServicePathPrefix
}

func (s *VtigerServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *VtigerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *VtigerServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/crm.Vtiger/CreateOrUpdateContact":
		msg := &ContactRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateOrUpdateContact(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/CreateOrUpdateLead":
		msg := &LeadRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateOrUpdateLead(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/CreateTicket":
		msg := &CreateOrUpdateTicketRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateTicket(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/GetCategories":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetCategories(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/GetContacts":
		msg := &GetContactsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetContacts(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/GetTicketStatusCount":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetTicketStatusCount(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/GetTickets":
		msg := &GetTicketsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetTickets(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/crm.Vtiger/UpdateTicket":
		msg := &CreateOrUpdateTicketRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateTicket(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
