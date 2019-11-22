// +build !generator

// Code generated by generator apix. DO NOT EDIT.

package services_affiliate

import (
	context "context"
	fmt "fmt"
	http "net/http"

	proto "github.com/golang/protobuf/proto"

	common "etop.vn/backend/pb/common"
	affiliate "etop.vn/backend/pb/services/affiliate"
	httprpc "etop.vn/backend/pkg/common/httprpc"
)

type Server interface {
	http.Handler
	PathPrefix() string
}

type AffiliateServiceServer struct {
	AffiliateAPI
}

func NewAffiliateServiceServer(svc AffiliateAPI) Server {
	return &AffiliateServiceServer{
		AffiliateAPI: svc,
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

func (s *AffiliateServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Affiliate/AffiliateGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.AffiliateGetProducts(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/CreateOrUpdateAffiliateCommissionSetting":
		msg := &affiliate.CreateOrUpdateCommissionSettingRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.CreateOrUpdateAffiliateCommissionSetting(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/CreateReferralCode":
		msg := &affiliate.CreateReferralCodeRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.CreateReferralCode(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetCommissions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.GetCommissions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetProductPromotionByProductID":
		msg := &affiliate.GetProductPromotionByProductIDRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.GetProductPromotionByProductID(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetReferralCodes":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.GetReferralCodes(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetReferrals":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.GetReferrals(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/GetTransactions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.GetTransactions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Affiliate/NotifyNewShopPurchase":
		msg := &affiliate.NotifyNewShopPurchaseRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.AffiliateAPI.NotifyNewShopPurchase(ctx, msg)
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

func (s *ShopServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Shop/CheckReferralCodeValid":
		msg := &affiliate.CheckReferralCodeValidRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.ShopAPI.CheckReferralCodeValid(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Shop/GetProductPromotion":
		msg := &affiliate.GetProductPromotionRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.ShopAPI.GetProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Shop/ShopGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.ShopAPI.ShopGetProducts(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type TradingServiceServer struct {
	TradingAPI
}

func NewTradingServiceServer(svc TradingAPI) Server {
	return &TradingServiceServer{
		TradingAPI: svc,
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

func (s *TradingServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.Trading/CreateOrUpdateTradingCommissionSetting":
		msg := &affiliate.CreateOrUpdateTradingCommissionSettingRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.CreateOrUpdateTradingCommissionSetting(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/CreateTradingProductPromotion":
		msg := &affiliate.CreateOrUpdateProductPromotionRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.CreateTradingProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/GetTradingProductPromotionByProductIDs":
		msg := &affiliate.GetTradingProductPromotionByIDsRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.GetTradingProductPromotionByProductIDs(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/GetTradingProductPromotions":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.GetTradingProductPromotions(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/TradingGetProducts":
		msg := &common.CommonListRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.TradingGetProducts(ctx, msg)
		}
		return msg, fn, nil
	case "/affiliate.Trading/UpdateTradingProductPromotion":
		msg := &affiliate.CreateOrUpdateProductPromotionRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.TradingAPI.UpdateTradingProductPromotion(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}

type UserServiceServer struct {
	UserAPI
}

func NewUserServiceServer(svc UserAPI) Server {
	return &UserServiceServer{
		UserAPI: svc,
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

func (s *UserServiceServer) parseRoute(path string) (reqMsg proto.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
	case "/affiliate.User/UpdateReferral":
		msg := &affiliate.UpdateReferralRequest{}
		fn := func(ctx context.Context) (proto.Message, error) {
			return s.UserAPI.UpdateReferral(ctx, msg)
		}
		return msg, fn, nil
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
