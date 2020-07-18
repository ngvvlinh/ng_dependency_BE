// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package api

import (
	context "context"
	fmt "fmt"
	http "net/http"

	capi "o.o/capi"
	httprpc "o.o/capi/httprpc"
)

func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	case func() CalcService:
		return NewCalcServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type CalcServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() CalcService
}

func NewCalcServiceServer(builder func() CalcService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &CalcServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const CalcServicePathPrefix = "/calc.Calc/"

func (s *CalcServiceServer) PathPrefix() string {
	return CalcServicePathPrefix
}

func (s *CalcServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CalcServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks)
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

func (s *CalcServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/calc.Calc/Calc":
		msg := &CreateEquationRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Calc(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/calc.Calc/Get":
		msg := &GetRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Get(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/calc.Calc/List":
		msg := &ListEquationRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.List(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/calc.Calc/Update":
		msg := &UpdateEquationRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.Update(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
