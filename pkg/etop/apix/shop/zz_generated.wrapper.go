// +build !generator

// Code generated by generator wrapper. DO NOT EDIT.

package xshop

import (
	"context"
	"time"

	api "etop.vn/api/top/external/shop"
	externaltypes "etop.vn/api/top/external/types"
	etop "etop.vn/api/top/int/etop"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	bus "etop.vn/backend/pkg/common/bus"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
)

func WrapHistoryService(s *HistoryService) api.HistoryService {
	return wrapHistoryService{s: s}
}

type wrapHistoryService struct {
	s *HistoryService
}

type GetChangesEndpoint struct {
	*externaltypes.GetChangesRequest
	Result  *externaltypes.Callback
	Context claims.ShopClaim
}

func (s wrapHistoryService) GetChanges(ctx context.Context, req *externaltypes.GetChangesRequest) (resp *externaltypes.Callback, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.History/GetChanges"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetChangesEndpoint{GetChangesRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetChanges(ctx, query)
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

func WrapMiscService(s *MiscService) api.MiscService {
	return wrapMiscService{s: s}
}

type wrapMiscService struct {
	s *MiscService
}

type CurrentAccountEndpoint struct {
	*cm.Empty
	Result  *etop.PublicAccountInfo
	Context claims.ShopClaim
}

func (s wrapMiscService) CurrentAccount(ctx context.Context, req *cm.Empty) (resp *etop.PublicAccountInfo, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/CurrentAccount"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CurrentAccountEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CurrentAccount(ctx, query)
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

type GetLocationListEndpoint struct {
	*cm.Empty
	Result  *externaltypes.LocationResponse
	Context claims.ShopClaim
}

func (s wrapMiscService) GetLocationList(ctx context.Context, req *cm.Empty) (resp *externaltypes.LocationResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/GetLocationList"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetLocationListEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetLocationList(ctx, query)
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

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context claims.PartnerClaim
}

func (s wrapMiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:        ctx,
		RequireAuth:    true,
		RequireAPIKey:  true,
		RequirePartner: true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &VersionInfoEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.Context.Partner = session.Partner
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.VersionInfo(ctx, query)
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

func WrapShippingService(s *ShippingService) api.ShippingService {
	return wrapShippingService{s: s}
}

type wrapShippingService struct {
	s *ShippingService
}

type CancelOrderEndpoint struct {
	*externaltypes.CancelOrderRequest
	Result  *externaltypes.OrderAndFulfillments
	Context claims.ShopClaim
}

func (s wrapShippingService) CancelOrder(ctx context.Context, req *externaltypes.CancelOrderRequest) (resp *externaltypes.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/CancelOrder"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CancelOrderEndpoint{CancelOrderRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CancelOrder(ctx, query)
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

type CreateAndConfirmOrderEndpoint struct {
	*externaltypes.CreateOrderRequest
	Result  *externaltypes.OrderAndFulfillments
	Context claims.ShopClaim
}

func (s wrapShippingService) CreateAndConfirmOrder(ctx context.Context, req *externaltypes.CreateOrderRequest) (resp *externaltypes.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/CreateAndConfirmOrder"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateAndConfirmOrderEndpoint{CreateOrderRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateAndConfirmOrder(ctx, query)
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

type GetFulfillmentEndpoint struct {
	*externaltypes.FulfillmentIDRequest
	Result  *externaltypes.Fulfillment
	Context claims.ShopClaim
}

func (s wrapShippingService) GetFulfillment(ctx context.Context, req *externaltypes.FulfillmentIDRequest) (resp *externaltypes.Fulfillment, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetFulfillment"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetFulfillmentEndpoint{FulfillmentIDRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetFulfillment(ctx, query)
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

type GetOrderEndpoint struct {
	*externaltypes.OrderIDRequest
	Result  *externaltypes.OrderAndFulfillments
	Context claims.ShopClaim
}

func (s wrapShippingService) GetOrder(ctx context.Context, req *externaltypes.OrderIDRequest) (resp *externaltypes.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetOrder"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetOrderEndpoint{OrderIDRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetOrder(ctx, query)
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

type GetShippingServicesEndpoint struct {
	*externaltypes.GetShippingServicesRequest
	Result  *externaltypes.GetShippingServicesResponse
	Context claims.ShopClaim
}

func (s wrapShippingService) GetShippingServices(ctx context.Context, req *externaltypes.GetShippingServicesRequest) (resp *externaltypes.GetShippingServicesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetShippingServices"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetShippingServicesEndpoint{GetShippingServicesRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetShippingServices(ctx, query)
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

func WrapWebhookService(s *WebhookService) api.WebhookService {
	return wrapWebhookService{s: s}
}

type wrapWebhookService struct {
	s *WebhookService
}

type CreateWebhookEndpoint struct {
	*externaltypes.CreateWebhookRequest
	Result  *externaltypes.Webhook
	Context claims.ShopClaim
}

func (s wrapWebhookService) CreateWebhook(ctx context.Context, req *externaltypes.CreateWebhookRequest) (resp *externaltypes.Webhook, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/CreateWebhook"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &CreateWebhookEndpoint{CreateWebhookRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.CreateWebhook(ctx, query)
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

type DeleteWebhookEndpoint struct {
	*externaltypes.DeleteWebhookRequest
	Result  *externaltypes.WebhooksResponse
	Context claims.ShopClaim
}

func (s wrapWebhookService) DeleteWebhook(ctx context.Context, req *externaltypes.DeleteWebhookRequest) (resp *externaltypes.WebhooksResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/DeleteWebhook"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &DeleteWebhookEndpoint{DeleteWebhookRequest: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.DeleteWebhook(ctx, query)
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

type GetWebhooksEndpoint struct {
	*cm.Empty
	Result  *externaltypes.WebhooksResponse
	Context claims.ShopClaim
}

func (s wrapWebhookService) GetWebhooks(ctx context.Context, req *cm.Empty) (resp *externaltypes.WebhooksResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/GetWebhooks"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		Context:       ctx,
		RequireAuth:   true,
		RequireAPIKey: true,
		RequireShop:   true,
	}
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
	query := &GetWebhooksEndpoint{Empty: req}
	query.Context.Claim = session.Claim
	query.Context.Shop = session.Shop
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	ctx = bus.NewRootContext(ctx)
	err = s.s.GetWebhooks(ctx, query)
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
