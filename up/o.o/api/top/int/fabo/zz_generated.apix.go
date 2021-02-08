// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package fabo

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
	case func() CustomerConversationService:
		return NewCustomerConversationServiceServer(builder, hooks...), true
	case func() CustomerService:
		return NewCustomerServiceServer(builder, hooks...), true
	case func() DemoService:
		return NewDemoServiceServer(builder, hooks...), true
	case func() ExtraShipmentService:
		return NewExtraShipmentServiceServer(builder, hooks...), true
	case func() PageService:
		return NewPageServiceServer(builder, hooks...), true
	case func() ShopService:
		return NewShopServiceServer(builder, hooks...), true
	case func() SummaryService:
		return NewSummaryServiceServer(builder, hooks...), true
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

const Path_CustomerConversation_CreateMessageTemplate = "/fabo.CustomerConversation/CreateMessageTemplate"
const Path_CustomerConversation_CreatePost = "/fabo.CustomerConversation/CreatePost"
const Path_CustomerConversation_DeleteMessageTemplate = "/fabo.CustomerConversation/DeleteMessageTemplate"
const Path_CustomerConversation_GetCustomerConversationByID = "/fabo.CustomerConversation/GetCustomerConversationByID"
const Path_CustomerConversation_HideOrUnHideComment = "/fabo.CustomerConversation/HideOrUnHideComment"
const Path_CustomerConversation_LikeOrUnLikeComment = "/fabo.CustomerConversation/LikeOrUnLikeComment"
const Path_CustomerConversation_ListCommentsByExternalPostID = "/fabo.CustomerConversation/ListCommentsByExternalPostID"
const Path_CustomerConversation_ListCustomerConversations = "/fabo.CustomerConversation/ListCustomerConversations"
const Path_CustomerConversation_ListLiveVideos = "/fabo.CustomerConversation/ListLiveVideos"
const Path_CustomerConversation_ListMessages = "/fabo.CustomerConversation/ListMessages"
const Path_CustomerConversation_MessageTemplateVariables = "/fabo.CustomerConversation/MessageTemplateVariables"
const Path_CustomerConversation_MessageTemplates = "/fabo.CustomerConversation/MessageTemplates"
const Path_CustomerConversation_SearchCustomerConversations = "/fabo.CustomerConversation/SearchCustomerConversations"
const Path_CustomerConversation_SendComment = "/fabo.CustomerConversation/SendComment"
const Path_CustomerConversation_SendMessage = "/fabo.CustomerConversation/SendMessage"
const Path_CustomerConversation_SendPrivateReply = "/fabo.CustomerConversation/SendPrivateReply"
const Path_CustomerConversation_UpdateMessageTemplate = "/fabo.CustomerConversation/UpdateMessageTemplate"
const Path_CustomerConversation_UpdateReadStatus = "/fabo.CustomerConversation/UpdateReadStatus"

func (s *CustomerConversationServiceServer) PathPrefix() string {
	return CustomerConversationServicePathPrefix
}

func (s *CustomerConversationServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CustomerConversationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CustomerConversationServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.CustomerConversation/CreateMessageTemplate":
		msg := &CreateMessageTemplateRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateMessageTemplate(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/CreatePost":
		msg := &CreatePostRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreatePost(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/DeleteMessageTemplate":
		msg := &DeleteMessageTemplateRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.DeleteMessageTemplate(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/GetCustomerConversationByID":
		msg := &GetCustomerConversationByIDRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetCustomerConversationByID(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/HideOrUnHideComment":
		msg := &HideOrUnHideCommentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.HideOrUnHideComment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/LikeOrUnLikeComment":
		msg := &LikeOrUnLikeCommentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.LikeOrUnLikeComment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListCommentsByExternalPostID":
		msg := &ListCommentsByExternalPostIDRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListCommentsByExternalPostID(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListCustomerConversations":
		msg := &ListCustomerConversationsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListCustomerConversations(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListLiveVideos":
		msg := &ListLiveVideosRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListLiveVideos(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/ListMessages":
		msg := &ListMessagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListMessages(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/MessageTemplateVariables":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.MessageTemplateVariables(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/MessageTemplates":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.MessageTemplates(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SearchCustomerConversations":
		msg := &SearchCustomerConversationRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SearchCustomerConversations(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SendComment":
		msg := &SendCommentRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendComment(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SendMessage":
		msg := &SendMessageRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendMessage(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/SendPrivateReply":
		msg := &SendPrivateReplyRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SendPrivateReply(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/UpdateMessageTemplate":
		msg := &UpdateMessageTemplateRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateMessageTemplate(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.CustomerConversation/UpdateReadStatus":
		msg := &UpdateReadStatusRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateReadStatus(newCtx, msg)
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

const Path_Customer_CreateFbUserCustomer = "/fabo.Customer/CreateFbUserCustomer"
const Path_Customer_GetFbUser = "/fabo.Customer/GetFbUser"
const Path_Customer_ListCustomersWithFbUsers = "/fabo.Customer/ListCustomersWithFbUsers"
const Path_Customer_ListFbUsers = "/fabo.Customer/ListFbUsers"
const Path_Customer_UpdateTags = "/fabo.Customer/UpdateTags"

func (s *CustomerServiceServer) PathPrefix() string {
	return CustomerServicePathPrefix
}

func (s *CustomerServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *CustomerServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CustomerServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Customer/CreateFbUserCustomer":
		msg := &CreateFbUserCustomerRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateFbUserCustomer(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/GetFbUser":
		msg := &GetFbUserRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetFbUser(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/ListCustomersWithFbUsers":
		msg := &ListCustomersWithFbUsersRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListCustomersWithFbUsers(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/ListFbUsers":
		msg := &ListFbUsersRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListFbUsers(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Customer/UpdateTags":
		msg := &UpdateUserTagsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateTags(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type DemoServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() DemoService
}

func NewDemoServiceServer(builder func() DemoService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &DemoServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const DemoServicePathPrefix = "/fabo.Demo/"

const Path_Demo_ListFeeds = "/fabo.Demo/ListFeeds"
const Path_Demo_ListLiveVideos = "/fabo.Demo/ListLiveVideos"

func (s *DemoServiceServer) PathPrefix() string {
	return DemoServicePathPrefix
}

func (s *DemoServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *DemoServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *DemoServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Demo/ListFeeds":
		msg := &ListFeedsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListFeeds(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Demo/ListLiveVideos":
		msg := &DemoListLiveVideosRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListLiveVideos(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ExtraShipmentServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() ExtraShipmentService
}

func NewExtraShipmentServiceServer(builder func() ExtraShipmentService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &ExtraShipmentServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const ExtraShipmentServicePathPrefix = "/fabo.ExtraShipment/"

const Path_ExtraShipment_CustomerReturnRate = "/fabo.ExtraShipment/CustomerReturnRate"

func (s *ExtraShipmentServiceServer) PathPrefix() string {
	return ExtraShipmentServicePathPrefix
}

func (s *ExtraShipmentServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *ExtraShipmentServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *ExtraShipmentServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.ExtraShipment/CustomerReturnRate":
		msg := &CustomerReturnRateRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CustomerReturnRate(newCtx, msg)
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

const Path_Page_CheckPermissions = "/fabo.Page/CheckPermissions"
const Path_Page_ConnectPages = "/fabo.Page/ConnectPages"
const Path_Page_ListPages = "/fabo.Page/ListPages"
const Path_Page_ListPosts = "/fabo.Page/ListPosts"
const Path_Page_RemovePages = "/fabo.Page/RemovePages"

func (s *PageServiceServer) PathPrefix() string {
	return PageServicePathPrefix
}

func (s *PageServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *PageServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *PageServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Page/CheckPermissions":
		msg := &CheckPagePermissionsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CheckPermissions(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/ConnectPages":
		msg := &ConnectPagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ConnectPages(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/ListPages":
		msg := &ListPagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListPages(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/ListPosts":
		msg := &ListPostsRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.ListPosts(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Page/RemovePages":
		msg := &RemovePagesRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.RemovePages(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShopServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() ShopService
}

func NewShopServiceServer(builder func() ShopService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &ShopServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const ShopServicePathPrefix = "/fabo.Shop/"

const Path_Shop_CreateTag = "/fabo.Shop/CreateTag"
const Path_Shop_DeleteTag = "/fabo.Shop/DeleteTag"
const Path_Shop_GetTags = "/fabo.Shop/GetTags"
const Path_Shop_UpdateTag = "/fabo.Shop/UpdateTag"

func (s *ShopServiceServer) PathPrefix() string {
	return ShopServicePathPrefix
}

func (s *ShopServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *ShopServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *ShopServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Shop/CreateTag":
		msg := &CreateFbShopTagRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.CreateTag(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Shop/DeleteTag":
		msg := &DeleteFbShopTagRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.DeleteTag(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Shop/GetTags":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.GetTags(newCtx, msg)
			return
		}
		return msg, fn, nil
	case "/fabo.Shop/UpdateTag":
		msg := &UpdateFbShopTagRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.UpdateTag(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type SummaryServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() SummaryService
}

func NewSummaryServiceServer(builder func() SummaryService, hooks ...httprpc.HooksBuilder) httprpc.Server {
	return &SummaryServiceServer{
		hooks:   httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const SummaryServicePathPrefix = "/fabo.Summary/"

const Path_Summary_SummaryShop = "/fabo.Summary/SummaryShop"

func (s *SummaryServiceServer) PathPrefix() string {
	return SummaryServicePathPrefix
}

func (s *SummaryServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *SummaryServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *SummaryServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/fabo.Summary/SummaryShop":
		msg := &SummaryShopRequest{}
		fn := func(ctx context.Context) (newCtx context.Context, resp capi.Message, err error) {
			inner := s.builder()
			info.Request, info.Inner = msg, inner
			newCtx, err = hooks.RequestRouted(ctx, *info)
			if err != nil {
				return
			}
			resp, err = inner.SummaryShop(newCtx, msg)
			return
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
