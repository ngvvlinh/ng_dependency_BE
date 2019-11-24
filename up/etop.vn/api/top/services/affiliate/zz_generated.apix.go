// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package affiliate

import (
	context "context"
	fmt "fmt"
	http "net/http"

	common "etop.vn/api/pb/common"
	affiliate "etop.vn/api/pb/services/affiliate"
	"etop.vn/capi"
	httprpc "etop.vn/capi/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type AffiliateServiceServer struct {
	inner AffiliateService
}

func NewAffiliateServiceServer(svc AffiliateService) Server {
	return &AffiliateServiceServer{
		inner: svc,
	}
}

const AffiliateServicePathPrefix = "/affiliate.Affiliate/"

func (s *AffiliateServiceServer) PathPrefix() string {
	return AffiliateServicePathPrefix
}

func (s *AffiliateServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *AffiliateServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Affiliate/AffiliateGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.AffiliateGetProducts(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/CreateOrUpdateAffiliateCommissionSetting":
		msg := &affiliate.CreateOrUpdateCommissionSettingRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateOrUpdateAffiliateCommissionSetting(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/CreateReferralCode":
		msg := &affiliate.CreateReferralCodeRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateReferralCode(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetCommissions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetCommissions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetProductPromotionByProductID":
		msg := &affiliate.GetProductPromotionByProductIDRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetProductPromotionByProductID(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetReferralCodes":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetReferralCodes(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetReferrals":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetReferrals(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetTransactions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetTransactions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/NotifyNewShopPurchase":
		msg := &affiliate.NotifyNewShopPurchaseRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.NotifyNewShopPurchase(ctx, msg)
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

const ShopServicePathPrefix = "/affiliate.Shop/"

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
	case "/affiliate.Shop/CheckReferralCodeValid":
		msg := &affiliate.CheckReferralCodeValidRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CheckReferralCodeValid(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Shop/GetProductPromotion":
		msg := &affiliate.GetProductPromotionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Shop/ShopGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.ShopGetProducts(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type TradingServiceServer struct {
	inner TradingService
}

func NewTradingServiceServer(svc TradingService) Server {
	return &TradingServiceServer{
		inner: svc,
	}
}

const TradingServicePathPrefix = "/affiliate.Trading/"

func (s *TradingServiceServer) PathPrefix() string {
	return TradingServicePathPrefix
}

func (s *TradingServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *TradingServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Trading/CreateOrUpdateTradingCommissionSetting":
		msg := &affiliate.CreateOrUpdateTradingCommissionSettingRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateOrUpdateTradingCommissionSetting(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/CreateTradingProductPromotion":
		msg := &affiliate.CreateOrUpdateProductPromotionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.CreateTradingProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/GetTradingProductPromotionByProductIDs":
		msg := &affiliate.GetTradingProductPromotionByIDsRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetTradingProductPromotionByProductIDs(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/GetTradingProductPromotions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.GetTradingProductPromotions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/TradingGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.TradingGetProducts(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/UpdateTradingProductPromotion":
		msg := &affiliate.CreateOrUpdateProductPromotionRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateTradingProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	inner UserService
}

func NewUserServiceServer(svc UserService) Server {
	return &UserServiceServer{
		inner: svc,
	}
}

const UserServicePathPrefix = "/affiliate.User/"

func (s *UserServiceServer) PathPrefix() string {
	return UserServicePathPrefix
}

func (s *UserServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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

func (s *UserServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.User/UpdateReferral":
		msg := &affiliate.UpdateReferralRequest{}
		fn := func(ctx context.Context) (capi.Message, error) {
			return s.inner.UpdateReferral(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
