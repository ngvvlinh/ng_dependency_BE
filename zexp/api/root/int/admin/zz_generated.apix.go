// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package admin

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	admin "etop.vn/backend/pb/etop/admin"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type AccountServiceServer struct {
	AccountAPI
}

func NewAccountServiceServer(svc AccountAPI) Server {
	return &AccountServiceServer{
		AccountAPI: svc,
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

func (s *AccountServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Account/CreatePartner":
		msg := new(admin.CreatePartnerRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AccountAPI.CreatePartner(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Account/GenerateAPIKey":
		msg := new(admin.GenerateAPIKeyRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AccountAPI.GenerateAPIKey(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type CreditServiceServer struct {
	CreditAPI
}

func NewCreditServiceServer(svc CreditAPI) Server {
	return &CreditServiceServer{
		CreditAPI: svc,
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

func (s *CreditServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Credit/ConfirmCredit":
		msg := new(admin.ConfirmCreditRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.ConfirmCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/CreateCredit":
		msg := new(admin.CreateCreditRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.CreateCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/DeleteCredit":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.DeleteCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/GetCredit":
		msg := new(admin.GetCreditRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.GetCredit(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/GetCredits":
		msg := new(admin.GetCreditsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.GetCredits(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Credit/UpdateCredit":
		msg := new(admin.UpdateCreditRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.CreditAPI.UpdateCredit(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type FulfillmentServiceServer struct {
	FulfillmentAPI
}

func NewFulfillmentServiceServer(svc FulfillmentAPI) Server {
	return &FulfillmentServiceServer{
		FulfillmentAPI: svc,
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

func (s *FulfillmentServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Fulfillment/GetFulfillment":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.FulfillmentAPI.GetFulfillment(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/GetFulfillments":
		msg := new(admin.GetFulfillmentsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.FulfillmentAPI.GetFulfillments(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Fulfillment/UpdateFulfillment":
		msg := new(admin.UpdateFulfillmentRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.FulfillmentAPI.UpdateFulfillment(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MiscServiceServer struct {
	MiscAPI
}

func NewMiscServiceServer(svc MiscAPI) Server {
	return &MiscServiceServer{
		MiscAPI: svc,
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

func (s *MiscServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Misc/AdminLoginAsAccount":
		msg := new(admin.LoginAsAccountRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MiscAPI.AdminLoginAsAccount(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Misc/VersionInfo":
		msg := new(common.Empty)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MiscAPI.VersionInfo(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type MoneyTransactionServiceServer struct {
	MoneyTransactionAPI
}

func NewMoneyTransactionServiceServer(svc MoneyTransactionAPI) Server {
	return &MoneyTransactionServiceServer{
		MoneyTransactionAPI: svc,
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

func (s *MoneyTransactionServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.MoneyTransaction/ConfirmMoneyTransaction":
		msg := new(admin.ConfirmMoneyTransactionRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.ConfirmMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingEtop":
		msg := new(admin.ConfirmMoneyTransactionShippingEtopRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.ConfirmMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternal":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.ConfirmMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternals":
		msg := new(common.IDsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.ConfirmMoneyTransactionShippingExternals(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/CreateMoneyTransactionShippingEtop":
		msg := new(common.IDsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.CreateMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/DeleteMoneyTransactionShippingEtop":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.DeleteMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/DeleteMoneyTransactionShippingExternal":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.DeleteMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransaction":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingEtop":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingEtops":
		msg := new(admin.GetMoneyTransactionShippingEtopsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransactionShippingEtops(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingExternal":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactionShippingExternals":
		msg := new(admin.GetMoneyTransactionShippingExternalsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransactionShippingExternals(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/GetMoneyTransactions":
		msg := new(admin.GetMoneyTransactionsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.GetMoneyTransactions(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/RemoveMoneyTransactionShippingExternalLines":
		msg := new(admin.RemoveMoneyTransactionShippingExternalLinesRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.RemoveMoneyTransactionShippingExternalLines(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransaction":
		msg := new(admin.UpdateMoneyTransactionRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.UpdateMoneyTransaction(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransactionShippingEtop":
		msg := new(admin.UpdateMoneyTransactionShippingEtopRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.UpdateMoneyTransactionShippingEtop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.MoneyTransaction/UpdateMoneyTransactionShippingExternal":
		msg := new(admin.UpdateMoneyTransactionShippingExternalRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.MoneyTransactionAPI.UpdateMoneyTransactionShippingExternal(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type NotificationServiceServer struct {
	NotificationAPI
}

func NewNotificationServiceServer(svc NotificationAPI) Server {
	return &NotificationServiceServer{
		NotificationAPI: svc,
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

func (s *NotificationServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Notification/CreateNotifications":
		msg := new(admin.CreateNotificationsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.NotificationAPI.CreateNotifications(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type OrderServiceServer struct {
	OrderAPI
}

func NewOrderServiceServer(svc OrderAPI) Server {
	return &OrderServiceServer{
		OrderAPI: svc,
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

func (s *OrderServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Order/GetOrder":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.OrderAPI.GetOrder(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Order/GetOrders":
		msg := new(admin.GetOrdersRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.OrderAPI.GetOrders(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Order/GetOrdersByIDs":
		msg := new(common.IDsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.OrderAPI.GetOrdersByIDs(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type ShopServiceServer struct {
	ShopAPI
}

func NewShopServiceServer(svc ShopAPI) Server {
	return &ShopServiceServer{
		ShopAPI: svc,
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

func (s *ShopServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/admin.Shop/GetShop":
		msg := new(common.IDRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.ShopAPI.GetShop(ctx, msg)
		}
		return msg, fn, nil
	case "/admin.Shop/GetShops":
		msg := new(admin.GetShopsRequest)
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.ShopAPI.GetShops(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
