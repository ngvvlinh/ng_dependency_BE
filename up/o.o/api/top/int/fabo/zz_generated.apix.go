// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package fabo

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
	case func() CustomerConversationService:
		return NewCustomerConversationServiceServer(builder, hooks...), true
	case func() CustomerService:
		return NewCustomerServiceServer(builder, hooks...), true
	case func() PageService:
		return NewPageServiceServer(builder, hooks...), true
	default:
		return nil, false
	}
}

type CustomerConversationServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() CustomerConversationService
}

func NewCustomerConversationServiceServer(builder func() CustomerConversationService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &CustomerConversationServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const CustomerConversationServicePathPrefix = "/fabo.CustomerConversation/"

func (s *CustomerConversationServiceServer) PathPrefix() string {
	return CustomerConversationServicePathPrefix
}

func (s *CustomerConversationServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CustomerConversationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CustomerConversationServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.CustomerConversation/ListCommentsByExternalPostID":
		msg := &ListCommentsByExternalPostIDRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListCommentsByExternalPostID(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListCustomerConversations":
		msg := &ListCustomerConversationsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListCustomerConversations(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListMessages":
		msg := &ListMessagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListMessages(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SendComment":
		msg := &SendCommentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendComment(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SendMessage":
		msg := &SendMessageRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendMessage(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/UpdateReadStatus":
		msg := &UpdateReadStatusRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateReadStatus(ctx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CustomerServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() CustomerService
}

func NewCustomerServiceServer(builder func() CustomerService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &CustomerServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const CustomerServicePathPrefix = "/fabo.Customer/"

func (s *CustomerServiceServer) PathPrefix() string {
	return CustomerServicePathPrefix
}

func (s *CustomerServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CustomerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CustomerServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Customer/CreateFbUserCustomer":
		msg := &CreateFbUserCustomerRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateFbUserCustomer(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/GetFbUser":
		msg := &GetFbUserRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetFbUser(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/ListFbUsers":
		msg := &ListFbUsersRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListFbUsers(ctx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type PageServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() PageService
}

func NewPageServiceServer(builder func() PageService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &PageServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const PageServicePathPrefix = "/fabo.Page/"

func (s *PageServiceServer) PathPrefix() string {
	return PageServicePathPrefix
}

func (s *PageServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *PageServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *PageServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Page/ConnectPages":
		msg := &ConnectPagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ConnectPages(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/ListPages":
		msg := &ListPagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListPages(ctx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/RemovePages":
		msg := &RemovePagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.BeforeServing(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.RemovePages(ctx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
