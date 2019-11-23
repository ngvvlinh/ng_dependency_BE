// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package api

import (
	"context"
	"time"

	cm "etop.vn/api/pb/common"
	affiliate "etop.vn/api/pb/services/affiliate"
	api "etop.vn/api/root/services/affiliate"
	common "etop.vn/backend/pkg/common"
	bus "etop.vn/backend/pkg/common/bus"
	metrics "etop.vn/backend/pkg/common/metrics"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapAffiliateService(s *AffiliateService, secret string) api.AffiliateService {
	return wrapAffiliateService{s: s, secret: secret}
}

type wrapAffiliateService struct {
	s      *AffiliateService
	secret string
}

type AffiliateGetProductsEndpoint struct {
	*cm.CommonListRequest
	Result  *affiliate.AffiliateGetProductsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) AffiliateGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.AffiliateGetProductsResponse, err error) {
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
	err = s.s.AffiliateGetProducts(ctx, query)
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
	*affiliate.CreateOrUpdateCommissionSettingRequest
	Result  *affiliate.CommissionSetting
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, req *affiliate.CreateOrUpdateCommissionSettingRequest) (resp *affiliate.CommissionSetting, err error) {
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
	err = s.s.CreateOrUpdateAffiliateCommissionSetting(ctx, query)
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
	*affiliate.CreateReferralCodeRequest
	Result  *affiliate.ReferralCode
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) CreateReferralCode(ctx context.Context, req *affiliate.CreateReferralCodeRequest) (resp *affiliate.ReferralCode, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/CreateReferralCode"
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
	query := &CreateReferralCodeEndpoint{CreateReferralCodeRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateReferralCode(ctx, query)
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
	Result  *affiliate.GetCommissionsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetCommissions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetCommissionsResponse, err error) {
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
		Context:          ctx,
		RequireAuth:      true,
		RequireAffiliate: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetCommissionsEndpoint{CommonListRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetCommissions(ctx, query)
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
	*affiliate.GetProductPromotionByProductIDRequest
	Result  *affiliate.GetProductPromotionByProductIDResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetProductPromotionByProductID(ctx context.Context, req *affiliate.GetProductPromotionByProductIDRequest) (resp *affiliate.GetProductPromotionByProductIDResponse, err error) {
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
	err = s.s.GetProductPromotionByProductID(ctx, query)
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
	Result  *affiliate.GetReferralCodesResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetReferralCodes(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetReferralCodesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetReferralCodes"
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
	query := &GetReferralCodesEndpoint{CommonListRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetReferralCodes(ctx, query)
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
	Result  *affiliate.GetReferralsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetReferrals(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetReferralsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Affiliate/GetReferrals"
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
	query := &GetReferralsEndpoint{CommonListRequest: req}
	query.Context.Affiliate = session.Affiliate
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetReferrals(ctx, query)
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
	Result  *affiliate.GetTransactionsResponse
	Context claims.AffiliateClaim
}

func (s wrapAffiliateService) GetTransactions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetTransactionsResponse, err error) {
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
	err = s.s.GetTransactions(ctx, query)
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
	*affiliate.NotifyNewShopPurchaseRequest
	Result  *affiliate.NotifyNewShopPurchaseResponse
	Context claims.EmptyClaim
}

func (s wrapAffiliateService) NotifyNewShopPurchase(ctx context.Context, req *affiliate.NotifyNewShopPurchaseRequest) (resp *affiliate.NotifyNewShopPurchaseResponse, err error) {
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
	err = s.s.NotifyNewShopPurchase(ctx, query)
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

func WrapShopService(s *ShopService) api.ShopService {
	return wrapShopService{s: s}
}

type wrapShopService struct {
	s *ShopService
}

type CheckReferralCodeValidEndpoint struct {
	*affiliate.CheckReferralCodeValidRequest
	Result  *affiliate.GetProductPromotionResponse
	Context claims.ShopClaim
}

func (s wrapShopService) CheckReferralCodeValid(ctx context.Context, req *affiliate.CheckReferralCodeValidRequest) (resp *affiliate.GetProductPromotionResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/CheckReferralCodeValid"
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
	query := &CheckReferralCodeValidEndpoint{CheckReferralCodeValidRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CheckReferralCodeValid(ctx, query)
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
	*affiliate.GetProductPromotionRequest
	Result  *affiliate.GetProductPromotionResponse
	Context claims.ShopClaim
}

func (s wrapShopService) GetProductPromotion(ctx context.Context, req *affiliate.GetProductPromotionRequest) (resp *affiliate.GetProductPromotionResponse, err error) {
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
	err = s.s.GetProductPromotion(ctx, query)
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
	Result  *affiliate.ShopGetProductsResponse
	Context claims.ShopClaim
}

func (s wrapShopService) ShopGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.ShopGetProductsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Shop/ShopGetProducts"
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
	query := &ShopGetProductsEndpoint{CommonListRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.ShopGetProducts(ctx, query)
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

func WrapTradingService(s *TradingService) api.TradingService {
	return wrapTradingService{s: s}
}

type wrapTradingService struct {
	s *TradingService
}

type CreateOrUpdateTradingCommissionSettingEndpoint struct {
	*affiliate.CreateOrUpdateTradingCommissionSettingRequest
	Result  *affiliate.SupplyCommissionSetting
	Context claims.ShopClaim
}

func (s wrapTradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, req *affiliate.CreateOrUpdateTradingCommissionSettingRequest) (resp *affiliate.SupplyCommissionSetting, err error) {
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
	query := &CreateOrUpdateTradingCommissionSettingEndpoint{CreateOrUpdateTradingCommissionSettingRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateOrUpdateTradingCommissionSetting(ctx, query)
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
	*affiliate.CreateOrUpdateProductPromotionRequest
	Result  *affiliate.ProductPromotion
	Context claims.ShopClaim
}

func (s wrapTradingService) CreateTradingProductPromotion(ctx context.Context, req *affiliate.CreateOrUpdateProductPromotionRequest) (resp *affiliate.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/CreateTradingProductPromotion"
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
	query := &CreateTradingProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateTradingProductPromotion(ctx, query)
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
	*affiliate.GetTradingProductPromotionByIDsRequest
	Result  *affiliate.GetTradingProductPromotionByIDsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) GetTradingProductPromotionByProductIDs(ctx context.Context, req *affiliate.GetTradingProductPromotionByIDsRequest) (resp *affiliate.GetTradingProductPromotionByIDsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/GetTradingProductPromotionByProductIDs"
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
	query := &GetTradingProductPromotionByProductIDsEndpoint{GetTradingProductPromotionByIDsRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetTradingProductPromotionByProductIDs(ctx, query)
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
	Result  *affiliate.GetProductPromotionsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) GetTradingProductPromotions(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.GetProductPromotionsResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/GetTradingProductPromotions"
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
	query := &GetTradingProductPromotionsEndpoint{CommonListRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetTradingProductPromotions(ctx, query)
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
	Result  *affiliate.SupplyGetProductsResponse
	Context claims.ShopClaim
}

func (s wrapTradingService) TradingGetProducts(ctx context.Context, req *cm.CommonListRequest) (resp *affiliate.SupplyGetProductsResponse, err error) {
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
	err = s.s.TradingGetProducts(ctx, query)
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
	*affiliate.CreateOrUpdateProductPromotionRequest
	Result  *affiliate.ProductPromotion
	Context claims.ShopClaim
}

func (s wrapTradingService) UpdateTradingProductPromotion(ctx context.Context, req *affiliate.CreateOrUpdateProductPromotionRequest) (resp *affiliate.ProductPromotion, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.Trading/UpdateTradingProductPromotion"
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
	query := &UpdateTradingProductPromotionEndpoint{CreateOrUpdateProductPromotionRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.UpdateTradingProductPromotion(ctx, query)
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

func WrapUserService(s *UserService) api.UserService {
	return wrapUserService{s: s}
}

type wrapUserService struct {
	s *UserService
}

type UpdateReferralEndpoint struct {
	*affiliate.UpdateReferralRequest
	Result  *affiliate.UserReferral
	Context claims.UserClaim
}

func (s wrapUserService) UpdateReferral(ctx context.Context, req *affiliate.UpdateReferralRequest) (resp *affiliate.UserReferral, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "affiliate.User/UpdateReferral"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:     ctx,
		RequireAuth: true,
		RequireUser: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &UpdateReferralEndpoint{UpdateReferralRequest: req}
	query.Context.Claim = session.Claim
	query.Context.User = session.User
	query.Context.Admin = session.Admin
	// Verify that the user has correct service type
	if session.Claim.AuthPartnerID != 0 {
		return nil, common.ErrPermissionDenied
	}
	ctx = bus.NewRootContext(ctx)
	err = s.s.UpdateReferral(ctx, query)
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
