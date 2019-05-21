// package shop generated by wrapper_gen. DO NOT EDIT.
package shopW

import (
	"context"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"

	cm "etop.vn/backend/pb/common"
	etop "etop.vn/backend/pb/etop"
	external "etop.vn/backend/pb/external"
	shop "etop.vn/backend/pb/external/shop"
	common "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	cmWrapper "etop.vn/backend/pkg/common/wrapper"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/middleware"
)

var ll = l.New()

var Client Shop

type Shop interface {
	shop.MiscService
	shop.WebhookService
	shop.HistoryService
	shop.ShippingService
}

type ShopClient struct {
	_MiscService     shop.MiscService
	_WebhookService  shop.WebhookService
	_HistoryService  shop.HistoryService
	_ShippingService shop.ShippingService
}

func NewShopClient(addr string, client *http.Client) Shop {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	addr = "http://" + addr
	return &ShopClient{
		_MiscService:     shop.NewMiscServiceProtobufClient(addr, client),
		_WebhookService:  shop.NewWebhookServiceProtobufClient(addr, client),
		_HistoryService:  shop.NewHistoryServiceProtobufClient(addr, client),
		_ShippingService: shop.NewShippingServiceProtobufClient(addr, client),
	}
}

func ConnectShopService(addr string, client *http.Client) error {
	Client = NewShopClient(addr, client)
	bus.AddHandler("client", func(ctx context.Context, q *CurrentAccountEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetLocationListEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *VersionInfoEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *CreateWebhookEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *DeleteWebhookEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetWebhooksEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetChangesEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *CancelOrderEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *CreateAndConfirmOrderEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetFulfillmentEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetOrderEndpoint) error { panic("Unexpected") })
	bus.AddHandler("client", func(ctx context.Context, q *GetShippingServicesEndpoint) error { panic("Unexpected") })
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := Client.VersionInfo(ctx, &cm.Empty{})
	if err == nil {
		ll.S.Infof("Connected to ShopService at %v", addr)
	}
	return err
}

func MustConnectShopService(addr string, client *http.Client) {
	err := ConnectShopService(addr, client)
	if err != nil {
		ll.Fatal("Unable to connect Shop", l.Error(err))
	}
}

type (
	EmptyClaim   = claims.EmptyClaim
	UserClaim    = claims.UserClaim
	AdminClaim   = claims.AdminClaim
	PartnerClaim = claims.PartnerClaim
	ShopClaim    = claims.ShopClaim
)

func (c *ShopClient) CurrentAccount(ctx context.Context, in *cm.Empty) (*etop.PublicAccountInfo, error) {
	resp, err := c._MiscService.CurrentAccount(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetLocationList(ctx context.Context, in *cm.Empty) (*external.LocationResponse, error) {
	resp, err := c._MiscService.GetLocationList(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) VersionInfo(ctx context.Context, in *cm.Empty) (*cm.VersionInfoResponse, error) {
	resp, err := c._MiscService.VersionInfo(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) CreateWebhook(ctx context.Context, in *external.CreateWebhookRequest) (*external.Webhook, error) {
	resp, err := c._WebhookService.CreateWebhook(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) DeleteWebhook(ctx context.Context, in *external.DeleteWebhookRequest) (*external.WebhooksResponse, error) {
	resp, err := c._WebhookService.DeleteWebhook(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetWebhooks(ctx context.Context, in *cm.Empty) (*external.WebhooksResponse, error) {
	resp, err := c._WebhookService.GetWebhooks(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetChanges(ctx context.Context, in *external.GetChangesRequest) (*external.Callback, error) {
	resp, err := c._HistoryService.GetChanges(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) CancelOrder(ctx context.Context, in *external.CancelOrderRequest) (*external.OrderAndFulfillments, error) {
	resp, err := c._ShippingService.CancelOrder(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) CreateAndConfirmOrder(ctx context.Context, in *external.CreateOrderRequest) (*external.OrderAndFulfillments, error) {
	resp, err := c._ShippingService.CreateAndConfirmOrder(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetFulfillment(ctx context.Context, in *external.FulfillmentIDRequest) (*external.Fulfillment, error) {
	resp, err := c._ShippingService.GetFulfillment(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetOrder(ctx context.Context, in *external.OrderIDRequest) (*external.OrderAndFulfillments, error) {
	resp, err := c._ShippingService.GetOrder(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}
func (c *ShopClient) GetShippingServices(ctx context.Context, in *external.GetShippingServicesRequest) (*external.GetShippingServicesResponse, error) {
	resp, err := c._ShippingService.GetShippingServices(ctx, in)

	node, ok := ctx.(*bus.NodeContext)
	if !ok {
		return resp, err
	}
	newNode := node.WithMessage(map[string]interface{}{
		"Request": in,
		"Result":  resp,
	})
	newNode.Error = err
	return resp, err
}

type Muxer interface {
	Handle(string, http.Handler)
}

func NewShopServer(mux Muxer, hooks *twirp.ServerHooks) {
	bus.Expect(&CurrentAccountEndpoint{})
	bus.Expect(&GetLocationListEndpoint{})
	bus.Expect(&VersionInfoEndpoint{})
	bus.Expect(&CreateWebhookEndpoint{})
	bus.Expect(&DeleteWebhookEndpoint{})
	bus.Expect(&GetWebhooksEndpoint{})
	bus.Expect(&GetChangesEndpoint{})
	bus.Expect(&CancelOrderEndpoint{})
	bus.Expect(&CreateAndConfirmOrderEndpoint{})
	bus.Expect(&GetFulfillmentEndpoint{})
	bus.Expect(&GetOrderEndpoint{})
	bus.Expect(&GetShippingServicesEndpoint{})
	mux.Handle(shop.MiscServicePathPrefix, shop.NewMiscServiceServer(MiscService{}, hooks))
	mux.Handle(shop.WebhookServicePathPrefix, shop.NewWebhookServiceServer(WebhookService{}, hooks))
	mux.Handle(shop.HistoryServicePathPrefix, shop.NewHistoryServiceServer(HistoryService{}, hooks))
	mux.Handle(shop.ShippingServicePathPrefix, shop.NewShippingServiceServer(ShippingService{}, hooks))
}

type ShopImpl struct {
	MiscService
	WebhookService
	HistoryService
	ShippingService
}

func NewShop() Shop {
	return ShopImpl{}
}

type MiscService struct{}

type CurrentAccountEndpoint struct {
	*cm.Empty
	Result  *etop.PublicAccountInfo
	Context ShopClaim
}

func (s MiscService) CurrentAccount(ctx context.Context, req *cm.Empty) (resp *etop.PublicAccountInfo, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/CurrentAccount"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type GetLocationListEndpoint struct {
	*cm.Empty
	Result  *external.LocationResponse
	Context ShopClaim
}

func (s MiscService) GetLocationList(ctx context.Context, req *cm.Empty) (resp *external.LocationResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/GetLocationList"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type VersionInfoEndpoint struct {
	*cm.Empty
	Result  *cm.VersionInfoResponse
	Context PartnerClaim
}

func (s MiscService) VersionInfo(ctx context.Context, req *cm.Empty) (resp *cm.VersionInfoResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Misc/VersionInfo"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type WebhookService struct{}

type CreateWebhookEndpoint struct {
	*external.CreateWebhookRequest
	Result  *external.Webhook
	Context ShopClaim
}

func (s WebhookService) CreateWebhook(ctx context.Context, req *external.CreateWebhookRequest) (resp *external.Webhook, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/CreateWebhook"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type DeleteWebhookEndpoint struct {
	*external.DeleteWebhookRequest
	Result  *external.WebhooksResponse
	Context ShopClaim
}

func (s WebhookService) DeleteWebhook(ctx context.Context, req *external.DeleteWebhookRequest) (resp *external.WebhooksResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/DeleteWebhook"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type GetWebhooksEndpoint struct {
	*cm.Empty
	Result  *external.WebhooksResponse
	Context ShopClaim
}

func (s WebhookService) GetWebhooks(ctx context.Context, req *cm.Empty) (resp *external.WebhooksResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Webhook/GetWebhooks"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type HistoryService struct{}

type GetChangesEndpoint struct {
	*external.GetChangesRequest
	Result  *external.Callback
	Context ShopClaim
}

func (s HistoryService) GetChanges(ctx context.Context, req *external.GetChangesRequest) (resp *external.Callback, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.History/GetChanges"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type ShippingService struct{}

type CancelOrderEndpoint struct {
	*external.CancelOrderRequest
	Result  *external.OrderAndFulfillments
	Context ShopClaim
}

func (s ShippingService) CancelOrder(ctx context.Context, req *external.CancelOrderRequest) (resp *external.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/CancelOrder"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type CreateAndConfirmOrderEndpoint struct {
	*external.CreateOrderRequest
	Result  *external.OrderAndFulfillments
	Context ShopClaim
}

func (s ShippingService) CreateAndConfirmOrder(ctx context.Context, req *external.CreateOrderRequest) (resp *external.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/CreateAndConfirmOrder"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type GetFulfillmentEndpoint struct {
	*external.FulfillmentIDRequest
	Result  *external.Fulfillment
	Context ShopClaim
}

func (s ShippingService) GetFulfillment(ctx context.Context, req *external.FulfillmentIDRequest) (resp *external.Fulfillment, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetFulfillment"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type GetOrderEndpoint struct {
	*external.OrderIDRequest
	Result  *external.OrderAndFulfillments
	Context ShopClaim
}

func (s ShippingService) GetOrder(ctx context.Context, req *external.OrderIDRequest) (resp *external.OrderAndFulfillments, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetOrder"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}

type GetShippingServicesEndpoint struct {
	*external.GetShippingServicesRequest
	Result  *external.GetShippingServicesResponse
	Context ShopClaim
}

func (s ShippingService) GetShippingServices(ctx context.Context, req *external.GetShippingServicesRequest) (resp *external.GetShippingServicesResponse, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "shop.Shipping/GetShippingServices"
	defer func() {
		recovered := recover()
		err = cmWrapper.RecoverAndLog(ctx, rpcName, session, req, resp, recovered, err, errs, t0)
	}()
	defer cmWrapper.Censor(req)
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmWrapper.HasErrors(resp)
	}
	return resp, err
}
