// package affiliate generated by wrapper_gen. DO NOT EDIT.
package affiliateW

import (
	"context"
	"net/http"
	"time"

	twirp "github.com/twitchtv/twirp"

	cm "etop.vn/backend/pb/common"
	affiliate "etop.vn/backend/pb/services/affiliate"
	common "etop.vn/backend/pkg/common"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
	bus "etop.vn/common/bus"
	l "etop.vn/common/l"
)

var ll = l.New()

type (
	EmptyClaim     = claims.EmptyClaim
	UserClaim      = claims.UserClaim
	AdminClaim     = claims.AdminClaim
	PartnerClaim   = claims.PartnerClaim
	ShopClaim      = claims.ShopClaim
	AffiliateClaim = claims.AffiliateClaim
)

type Muxer interface {
	Handle(string, http.Handler)
}

func NewAffiliateServer(mux Muxer, hooks *twirp.ServerHooks, secret string) {
	if secret == "" {
		ll.Fatal("Secret is empty")
	}
	bus.Expect(&CreateOrUpdateTradingCommissionSettingEndpoint{})
	bus.Expect(&CreateProductPromotionEndpoint{})
	bus.Expect(&GetProductPromotionsEndpoint{})
	bus.Expect(&TradingGetProductsEndpoint{})
	bus.Expect(&UpdateProductPromotionEndpoint{})
	bus.Expect(&GetProductPromotionEndpoint{})
	bus.Expect(&AffiliateGetProductsEndpoint{})
	bus.Expect(&CreateOrUpdateAffiliateCommissionSettingEndpoint{})
	bus.Expect(&GetCommissionsEndpoint{})
	bus.Expect(&GetProductPromotionByProductIDEndpoint{})
	bus.Expect(&GetTransactionsEndpoint{})
	bus.Expect(&NotifyNewShopPurchaseEndpoint{})
	mux.Handle(affiliate.TradingServicePathPrefix, affiliate.NewTradingServiceServer(TradingService{secret: secret}, hooks))
	mux.Handle(affiliate.ShopServicePathPrefix, affiliate.NewShopServiceServer(ShopService{secret: secret}, hooks))
	mux.Handle(affiliate.AffiliateServicePathPrefix, affiliate.NewAffiliateServiceServer(AffiliateService{secret: secret}, hooks))
}

type AffiliateImpl struct {
	TradingService
	ShopService
	AffiliateService
}

type TradingService struct{ secret string }

type CreateOrUpdateTradingCommissionSettingEndpoint struct {
	*affiliate.CreateOrUpdateCommissionSettingRequest
	Result  *affiliate.CommissionSetting
	Context ShopClaim
}

func (s TradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, req *affiliate.CreateOrUpdateCommissionSettingRequest) (resp *affiliate.CommissionSetting, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/CreateOrUpdateTradingCommissionSetting"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateTradingCommissionSettingEndpoint{CreateOrUpdateCommissionSettingRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type CreateProductPromotionEndpoint struct {
	*affiliate.CreateOrUpdateProductPromotionRequest
	Result  *affiliate.ProductPromotion
	Context ShopClaim
}

func (s TradingService) CreateProductPromotion(ctx context.Context, req *affiliate.CreateOrUpdateProductPromotionRequest) (resp *affiliate.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/CreateProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type GetProductPromotionsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.GetProductPromotionsResponse
	Context ShopClaim
}

func (s TradingService) GetProductPromotions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetProductPromotionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/GetProductPromotions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetProductPromotionsEndpoint{CommonListRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type TradingGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.ShopGetProductsResponse
	Context ShopClaim
}

func (s TradingService) TradingGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.ShopGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/TradingGetProducts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &TradingGetProductsEndpoint{CommonListRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type UpdateProductPromotionEndpoint struct {
	*affiliate.CreateOrUpdateProductPromotionRequest
	Result  *affiliate.ProductPromotion
	Context ShopClaim
}

func (s TradingService) UpdateProductPromotion(ctx context.Context, req *affiliate.CreateOrUpdateProductPromotionRequest) (resp *affiliate.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/UpdateProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &UpdateProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type ShopService struct{ secret string }

type GetProductPromotionEndpoint struct {
	*affiliate.GetProductPromotionRequest
	Result  *affiliate.GetProductPromotionResponse
	Context ShopClaim
}

func (s ShopService) GetProductPromotion(ctx context.Context, req *affiliate.GetProductPromotionRequest) (resp *affiliate.GetProductPromotionResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/GetProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetProductPromotionEndpoint{GetProductPromotionRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type AffiliateService struct{ secret string }

type AffiliateGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.AffiliateGetProductsResponse
	Context AffiliateClaim
}

func (s AffiliateService) AffiliateGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.AffiliateGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/AffiliateGetProducts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:          ctx,
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &AffiliateGetProductsEndpoint{CommonListRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type CreateOrUpdateAffiliateCommissionSettingEndpoint struct {
	*affiliate.CreateOrUpdateCommissionSettingRequest
	Result  *affiliate.CommissionSetting
	Context AffiliateClaim
}

func (s AffiliateService) CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, req *affiliate.CreateOrUpdateCommissionSettingRequest) (resp *affiliate.CommissionSetting, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/CreateOrUpdateAffiliateCommissionSetting"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:          ctx,
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateAffiliateCommissionSettingEndpoint{CreateOrUpdateCommissionSettingRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type GetCommissionsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.GetCommissionsResponse
	Context ShopClaim
}

func (s AffiliateService) GetCommissions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetCommissionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetCommissions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireShop: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetCommissionsEndpoint{CommonListRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type GetProductPromotionByProductIDEndpoint struct {
	*affiliate.GetProductPromotionByProductIDRequest
	Result  *affiliate.GetProductPromotionByProductIDResponse
	Context AffiliateClaim
}

func (s AffiliateService) GetProductPromotionByProductID(ctx context.Context, req *affiliate.GetProductPromotionByProductIDRequest) (resp *affiliate.GetProductPromotionByProductIDResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetProductPromotionByProductID"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:          ctx,
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetProductPromotionByProductIDEndpoint{GetProductPromotionByProductIDRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type GetTransactionsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.GetTransactionsResponse
	Context AffiliateClaim
}

func (s AffiliateService) GetTransactions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetTransactionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetTransactions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:          ctx,
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetTransactionsEndpoint{CommonListRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}

type NotifyNewShopPurchaseEndpoint struct {
	*affiliate.NotifyNewShopPurchaseRequest
	Result  *affiliate.NotifyNewShopPurchaseResponse
	Context EmptyClaim
}

func (s AffiliateService) NotifyNewShopPurchase(ctx context.Context, req *affiliate.NotifyNewShopPurchaseRequest) (resp *affiliate.NotifyNewShopPurchaseResponse, err error) {
	t0 := time.Now()
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/NotifyNewShopPurchase"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	query := &NotifyNewShopPurchaseEndpoint{NotifyNewShopPurchaseRequest: req}
	// Verify secret token
	token := middleware.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}
