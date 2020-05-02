// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package api

import (
	"context"
	"time"

	api "o.o/api/top/services/affiliate"
	cm "o.o/api/top/types/common"
	common "o.o/backend/pkg/common"
	cmwrapper "o.o/backend/pkg/common/apifw/wrapper"
	bus "o.o/backend/pkg/common/bus"
	headers "o.o/backend/pkg/common/headers"
	claims "o.o/backend/pkg/etop/authorize/claims"
	middleware "o.o/backend/pkg/etop/authorize/middleware"
)

func WrapAffiliateService(s func() *AffiliateService, secret string) func() api.AffiliateService {
	return func() api.AffiliateService { return wrapAffiliateService{s: s, secret: secret} }
}

type wrapAffiliateService struct {
	s      func() *AffiliateService
	secret string
}

type AffiliateGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.AffiliateGetProductsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) AffiliateGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *api.AffiliateGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/AffiliateGetProducts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &AffiliateGetProductsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().AffiliateGetProducts(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type CreateOrUpdateAffiliateCommissionSettingEndpoint struct {
	*api.CreateOrUpdateCommissionSettingRequest
	Result  *api.CommissionSetting
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, req *api.CreateOrUpdateCommissionSettingRequest) (resp *api.CommissionSetting, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/CreateOrUpdateAffiliateCommissionSetting"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateAffiliateCommissionSettingEndpoint{CreateOrUpdateCommissionSettingRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().CreateOrUpdateAffiliateCommissionSetting(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type CreateReferralCodeEndpoint struct {
	*api.CreateReferralCodeRequest
	Result  *api.ReferralCode
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) CreateReferralCode(ctx context.Context, req *api.CreateReferralCodeRequest) (resp *api.ReferralCode, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/CreateReferralCode"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateReferralCodeEndpoint{CreateReferralCodeRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().CreateReferralCode(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetCommissionsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.GetCommissionsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetCommissions(ctx context.Context, req *cm.CommonListRequest) (resp *api.GetCommissionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetCommissions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetCommissionsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetCommissions(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetProductPromotionByProductIDEndpoint struct {
	*api.GetProductPromotionByProductIDRequest
	Result  *api.GetProductPromotionByProductIDResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetProductPromotionByProductID(ctx context.Context, req *api.GetProductPromotionByProductIDRequest) (resp *api.GetProductPromotionByProductIDResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetProductPromotionByProductID"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetProductPromotionByProductIDEndpoint{GetProductPromotionByProductIDRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetProductPromotionByProductID(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetReferralCodesEndpoint struct {
	*cm.CommonListRequest
	Result  *api.GetReferralCodesResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetReferralCodes(ctx context.Context, req *cm.CommonListRequest) (resp *api.GetReferralCodesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetReferralCodes"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetReferralCodesEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetReferralCodes(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetReferralsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.GetReferralsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetReferrals(ctx context.Context, req *cm.CommonListRequest) (resp *api.GetReferralsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetReferrals"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetReferralsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetReferrals(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetTransactionsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.GetTransactionsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetTransactions(ctx context.Context, req *cm.CommonListRequest) (resp *api.GetTransactionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetTransactions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetTransactionsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Affiliate = session.Affiliate
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetTransactions(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type NotifyNewShopPurchaseEndpoint struct {
	*api.NotifyNewShopPurchaseRequest
	Result  *api.NotifyNewShopPurchaseResponse
	Context claims.EmptyClaim
}

func (s wrapAffiliateService) NotifyNewShopPurchase(ctx context.Context, req *api.NotifyNewShopPurchaseRequest) (resp *api.NotifyNewShopPurchaseResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/NotifyNewShopPurchase"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, nil, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &NotifyNewShopPurchaseEndpoint{NotifyNewShopPurchaseRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	// Verify secret token
	token := headers.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s().NotifyNewShopPurchase(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

func WrapShopService(s func() *ShopService) func() api.ShopService {
	return func() api.ShopService { return wrapShopService{s: s} }
}

type wrapShopService struct {
	s func() *ShopService
}

type CheckReferralCodeValidEndpoint struct {
	*api.CheckReferralCodeValidRequest
	Result  *api.GetProductPromotionResponse
	Context claims.ShopClaim
}

func (s wrapShopService) CheckReferralCodeValid(ctx context.Context, req *api.CheckReferralCodeValidRequest) (resp *api.GetProductPromotionResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/CheckReferralCodeValid"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CheckReferralCodeValidEndpoint{CheckReferralCodeValidRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().CheckReferralCodeValid(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetProductPromotionEndpoint struct {
	*api.GetProductPromotionRequest
	Result  *api.GetProductPromotionResponse
	Context claims.ShopClaim
}

func (s wrapShopService) GetProductPromotion(ctx context.Context, req *api.GetProductPromotionRequest) (resp *api.GetProductPromotionResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/GetProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetProductPromotionEndpoint{GetProductPromotionRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetProductPromotion(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type ShopGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.ShopGetProductsResponse
	Context claims.ShopClaim
}

func (s wrapShopService) ShopGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *api.ShopGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/ShopGetProducts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &ShopGetProductsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().ShopGetProducts(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

func WrapTradingService(s func() *TradingService) func() api.TradingService {
	return func() api.TradingService { return wrapTradingService{s: s} }
}

type wrapTradingService struct {
	s func() *TradingService
}

type CreateOrUpdateTradingCommissionSettingEndpoint struct {
	*api.CreateOrUpdateTradingCommissionSettingRequest
	Result  *api.SupplyCommissionSetting
	Context claims.ShopClaim
}

func (s wrapTradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, req *api.CreateOrUpdateTradingCommissionSettingRequest) (resp *api.SupplyCommissionSetting, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/CreateOrUpdateTradingCommissionSetting"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateOrUpdateTradingCommissionSettingEndpoint{CreateOrUpdateTradingCommissionSettingRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().CreateOrUpdateTradingCommissionSetting(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type CreateTradingProductPromotionEndpoint struct {
	*api.CreateOrUpdateProductPromotionRequest
	Result  *api.ProductPromotion
	Context claims.ShopClaim
}

func (s wrapTradingService) CreateTradingProductPromotion(ctx context.Context, req *api.CreateOrUpdateProductPromotionRequest) (resp *api.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/CreateTradingProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateTradingProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().CreateTradingProductPromotion(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetTradingProductPromotionByProductIDsEndpoint struct {
	*api.GetTradingProductPromotionByIDsRequest
	Result  *api.GetTradingProductPromotionByIDsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) GetTradingProductPromotionByProductIDs(ctx context.Context, req *api.GetTradingProductPromotionByIDsRequest) (resp *api.GetTradingProductPromotionByIDsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/GetTradingProductPromotionByProductIDs"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetTradingProductPromotionByProductIDsEndpoint{GetTradingProductPromotionByIDsRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetTradingProductPromotionByProductIDs(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type GetTradingProductPromotionsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.GetProductPromotionsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) GetTradingProductPromotions(ctx context.Context, req *cm.CommonListRequest) (resp *api.GetProductPromotionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/GetTradingProductPromotions"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetTradingProductPromotionsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().GetTradingProductPromotions(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type TradingGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *api.SupplyGetProductsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) TradingGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *api.SupplyGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/TradingGetProducts"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &TradingGetProductsEndpoint{CommonListRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().TradingGetProducts(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

type UpdateTradingProductPromotionEndpoint struct {
	*api.CreateOrUpdateProductPromotionRequest
	Result  *api.ProductPromotion
	Context claims.ShopClaim
}

func (s wrapTradingService) UpdateTradingProductPromotion(ctx context.Context, req *api.CreateOrUpdateProductPromotionRequest) (resp *api.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/UpdateTradingProductPromotion"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireShop: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &UpdateTradingProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s().UpdateTradingProductPromotion(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}

func WrapUserService(s func() *UserService) func() api.UserService {
	return func() api.UserService { return wrapUserService{s: s} }
}

type wrapUserService struct {
	s func() *UserService
}

type UpdateReferralEndpoint struct {
	*api.UpdateReferralRequest
	Result  *api.UserReferral
	Context claims.UserClaim
}

func (s wrapUserService) UpdateReferral(ctx context.Context, req *api.UpdateReferralRequest) (resp *api.UserReferral, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.User/UpdateReferral"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		RequireAuth: true,
		RequireUser: true,
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &UpdateReferralEndpoint{UpdateReferralRequest: req}
	if session != nil {
		query.Context.Claim = session.Claim
	}
	query.Context.User = session.User
	query.Context.Admin = session.Admin
	ctx = bus.NewRootContext(ctx)
	err = s.s().UpdateReferral(ctx, query)
	resp = query.Result
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil
}
