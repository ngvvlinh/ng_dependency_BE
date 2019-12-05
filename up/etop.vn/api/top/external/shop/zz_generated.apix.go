// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package shop

import (
	context "context"
	fmt "fmt"
	http "net/http"

	externaltypes "etop.vn/api/top/external/types"
	common "etop.vn/api/top/types/common"
	capi "etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type HistoryServiceServer struct {
	inner HistoryService
}

func NewHistoryServiceServer(svc HistoryService) Server {
	return &HistoryServiceServer{
		inner: svc,
	}
}

const HistoryServicePathPrefix = "/shop.History/"

func (s *HistoryServiceServer) PathPrefix() string {
	return HistoryServicePathPrefix
}

func (s *HistoryServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *HistoryServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/shop.History/GetChanges":
		msg := &externaltypes.GetChangesRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetChanges(ctx, msg)
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

const MiscServicePathPrefix = "/shop.Misc/"

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
	case "/shop.Misc/CurrentAccount":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CurrentAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Misc/GetLocationList":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetLocationList(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Misc/VersionInfo":
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

type ShippingServiceServer struct {
	inner ShippingService
}

func NewShippingServiceServer(svc ShippingService) Server {
	return &ShippingServiceServer{
		inner: svc,
	}
}

const ShippingServicePathPrefix = "/shop.Shipping/"

func (s *ShippingServiceServer) PathPrefix() string {
	return ShippingServicePathPrefix
}

func (s *ShippingServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *ShippingServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/shop.Shipping/CancelOrder":
		msg := &externaltypes.CancelOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CancelOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Shipping/CreateAndConfirmOrder":
		msg := &externaltypes.CreateOrderRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateAndConfirmOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Shipping/GetFulfillment":
		msg := &externaltypes.FulfillmentIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Shipping/GetOrder":
		msg := &externaltypes.OrderIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Shipping/GetShippingServices":
		msg := &externaltypes.GetShippingServicesRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetShippingServices(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type WebhookServiceServer struct {
	inner WebhookService
}

func NewWebhookServiceServer(svc WebhookService) Server {
	return &WebhookServiceServer{
		inner: svc,
	}
}

const WebhookServicePathPrefix = "/shop.Webhook/"

func (s *WebhookServiceServer) PathPrefix() string {
	return WebhookServicePathPrefix
}

func (s *WebhookServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *WebhookServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/shop.Webhook/CreateWebhook":
		msg := &externaltypes.CreateWebhookRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateWebhook(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Webhook/DeleteWebhook":
		msg := &externaltypes.DeleteWebhookRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteWebhook(ctx, msg)
		}
		return msg, fn, nil
	case "/shop.Webhook/GetWebhooks":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetWebhooks(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
