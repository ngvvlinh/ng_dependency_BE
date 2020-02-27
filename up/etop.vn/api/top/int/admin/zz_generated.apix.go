// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package admin

import (
	context "context"
	fmt "fmt"
	http "net/http"

	inttypes "etop.vn/api/top/int/types"
	common "etop.vn/api/top/types/common"
	capi "etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type AccountServiceServer struct {
	inner AccountService
}

func NewAccountServiceServer(svc AccountService) Server {
	return &AccountServiceServer{
		inner: svc,
	}
}

const AccountServicePathPrefix = "/admin.Account/"

func (s *AccountServiceServer) PathPrefix() string {
	return AccountServicePathPrefix
}

func (s *AccountServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *AccountServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Account/CreatePartner":
		msg := &CreatePartnerRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreatePartner(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Account/GenerateAPIKey":
		msg := &GenerateAPIKeyRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GenerateAPIKey(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ConnectionServiceServer struct {
	inner ConnectionService
}

func NewConnectionServiceServer(svc ConnectionService) Server {
	return &ConnectionServiceServer{
		inner: svc,
	}
}

const ConnectionServicePathPrefix = "/admin.Connection/"

func (s *ConnectionServiceServer) PathPrefix() string {
	return ConnectionServicePathPrefix
}

func (s *ConnectionServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *ConnectionServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Connection/ConfirmConnection":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmConnection(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Connection/CreateTopshipConnection":
		msg := &inttypes.CreateTopshipConnectionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateTopshipConnection(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Connection/DisableConnection":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DisableConnection(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Connection/GetConnections":
		msg := &common.Empty{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetConnections(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CreditServiceServer struct {
	inner CreditService
}

func NewCreditServiceServer(svc CreditService) Server {
	return &CreditServiceServer{
		inner: svc,
	}
}

const CreditServicePathPrefix = "/admin.Credit/"

func (s *CreditServiceServer) PathPrefix() string {
	return CreditServicePathPrefix
}

func (s *CreditServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *CreditServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Credit/ConfirmCredit":
		msg := &ConfirmCreditRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/CreateCredit":
		msg := &CreateCreditRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/DeleteCredit":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/GetCredit":
		msg := &GetCreditRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/GetCredits":
		msg := &GetCreditsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetCredits(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/UpdateCredit":
		msg := &UpdateCreditRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateCredit(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type FulfillmentServiceServer struct {
	inner FulfillmentService
}

func NewFulfillmentServiceServer(svc FulfillmentService) Server {
	return &FulfillmentServiceServer{
		inner: svc,
	}
}

const FulfillmentServicePathPrefix = "/admin.Fulfillment/"

func (s *FulfillmentServiceServer) PathPrefix() string {
	return FulfillmentServicePathPrefix
}

func (s *FulfillmentServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *FulfillmentServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Fulfillment/GetFulfillment":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/GetFulfillments":
		msg := &GetFulfillmentsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetFulfillments(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/UpdateFulfillment":
		msg := &UpdateFulfillmentRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/UpdateFulfillmentShippingFee":
		msg := &UpdateFulfillmentShippingFeeRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateFulfillmentShippingFee(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/UpdateFulfillmentShippingState":
		msg := &UpdateFulfillmentShippingStateRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateFulfillmentShippingState(ctx, msg)
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

const MiscServicePathPrefix = "/admin.Misc/"

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
	case "/admin.Misc/AdminLoginAsAccount":
		msg := &LoginAsAccountRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.AdminLoginAsAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Misc/VersionInfo":
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

type MoneyTransactionServiceServer struct {
	inner MoneyTransactionService
}

func NewMoneyTransactionServiceServer(svc MoneyTransactionService) Server {
	return &MoneyTransactionServiceServer{
		inner: svc,
	}
}

const MoneyTransactionServicePathPrefix = "/admin.MoneyTransaction/"

func (s *MoneyTransactionServiceServer) PathPrefix() string {
	return MoneyTransactionServicePathPrefix
}

func (s *MoneyTransactionServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *MoneyTransactionServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.MoneyTransaction/ConfirmMoneyTransaction":
		msg := &ConfirmMoneyTransactionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingEtop":
		msg := &ConfirmMoneyTransactionShippingEtopRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternal":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternals":
		msg := &common.IDsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ConfirmMoneyTransactionShippingExternals(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/CreateMoneyTransactionShippingEtop":
		msg := &common.IDsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/DeleteMoneyTransactionShippingEtop":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/DeleteMoneyTransactionShippingExternal":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.DeleteMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransaction":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingEtop":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingEtops":
		msg := &GetMoneyTransactionShippingEtopsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransactionShippingEtops(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingExternal":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingExternals":
		msg := &GetMoneyTransactionShippingExternalsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransactionShippingExternals(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactions":
		msg := &GetMoneyTransactionsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetMoneyTransactions(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/RemoveMoneyTransactionShippingExternalLines":
		msg := &RemoveMoneyTransactionShippingExternalLinesRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.RemoveMoneyTransactionShippingExternalLines(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransaction":
		msg := &UpdateMoneyTransactionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransactionShippingEtop":
		msg := &UpdateMoneyTransactionShippingEtopRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransactionShippingExternal":
		msg := &UpdateMoneyTransactionShippingExternalRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type NotificationServiceServer struct {
	inner NotificationService
}

func NewNotificationServiceServer(svc NotificationService) Server {
	return &NotificationServiceServer{
		inner: svc,
	}
}

const NotificationServicePathPrefix = "/admin.Notification/"

func (s *NotificationServiceServer) PathPrefix() string {
	return NotificationServicePathPrefix
}

func (s *NotificationServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *NotificationServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Notification/CreateNotifications":
		msg := &CreateNotificationsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateNotifications(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type OrderServiceServer struct {
	inner OrderService
}

func NewOrderServiceServer(svc OrderService) Server {
	return &OrderServiceServer{
		inner: svc,
	}
}

const OrderServicePathPrefix = "/admin.Order/"

func (s *OrderServiceServer) PathPrefix() string {
	return OrderServicePathPrefix
}

func (s *OrderServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *OrderServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Order/GetOrder":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Order/GetOrders":
		msg := &GetOrdersRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrders(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Order/GetOrdersByIDs":
		msg := &common.IDsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetOrdersByIDs(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShopServiceServer struct {
	inner ShopService
}

func NewShopServiceServer(svc ShopService) Server {
	return &ShopServiceServer{
		inner: svc,
	}
}

const ShopServicePathPrefix = "/admin.Shop/"

func (s *ShopServiceServer) PathPrefix() string {
	return ShopServicePathPrefix
}

func (s *ShopServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *ShopServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Shop/GetShop":
		msg := &common.IDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetShop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Shop/GetShops":
		msg := &GetShopsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetShops(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
