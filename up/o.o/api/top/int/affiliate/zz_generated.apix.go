// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package affiliate

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
	case func() AccountService:
		return NewAccountServiceServer(builder, hooks...), true
	case func() MiscService:
		return NewMiscServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type AccountServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() AccountService
}

func NewAccountServiceServer(builder func() AccountService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &AccountServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const AccountServicePathPrefix = "/affiliate.Account/"

func (s *AccountServiceServer) PathPrefix() string {
	return AccountServicePathPrefix
}

func (s *AccountServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks.BuildHooks())
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.BeforeRequest(ctx, *info)
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

func (s *AccountServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Account/DeleteAffiliate":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			ctx, err := hooks.BeforeServing(ctx, *info)
			if err != nil {
				return nil, err
			}
			return inner.DeleteAffiliate(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Account/RegisterAffiliate":
		msg := &RegisterAffiliateRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			ctx, err := hooks.BeforeServing(ctx, *info)
			if err != nil {
				return nil, err
			}
			return inner.RegisterAffiliate(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Account/UpdateAffiliate":
		msg := &UpdateAffiliateRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			ctx, err := hooks.BeforeServing(ctx, *info)
			if err != nil {
				return nil, err
			}
			return inner.UpdateAffiliate(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Account/UpdateAffiliateBankAccount":
		msg := &UpdateAffiliateBankAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			ctx, err := hooks.BeforeServing(ctx, *info)
			if err != nil {
				return nil, err
			}
			return inner.UpdateAffiliateBankAccount(ctx, msg)
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

const MiscServicePathPrefix = "/affiliate.Misc/"

func (s *MiscServiceServer) PathPrefix() string {
	return MiscServicePathPrefix
}

func (s *MiscServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks.BuildHooks())
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.BeforeRequest(ctx, *info)
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
	case "/affiliate.Misc/VersionInfo":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			ctx, err := hooks.BeforeServing(ctx, *info)
			if err != nil {
				return nil, err
			}
			return inner.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
