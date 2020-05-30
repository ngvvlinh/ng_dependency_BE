package build

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/shipnow"
	shippingcore "o.o/api/main/shipping"
	subscriptioncore "o.o/api/subscripting/subscription"
	webserverinternal "o.o/api/webserver"
	"o.o/backend/cmd/etop-server/config"
	vtpaygatewayaggregate "o.o/backend/com/external/payment/vtpay/gateway/aggregate"
	vtpaygatewayserver "o.o/backend/com/external/payment/vtpay/gateway/server"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/com/web/ecom/webserver"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/lifecycle"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/sqltrace"
	"o.o/backend/pkg/etop/api/shop/imports"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/logic/hotfix"
	imcsvghtk "o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	webhookahamove "o.o/backend/pkg/integration/shipnow/ahamove/webhook"
	"o.o/backend/pkg/integration/shipping/ghn"
	webhookghn "o.o/backend/pkg/integration/shipping/ghn/webhook"
	"o.o/backend/pkg/integration/shipping/ghtk"
	webhookghtk "o.o/backend/pkg/integration/shipping/ghtk/webhook"
	"o.o/backend/pkg/integration/shipping/vtpost"
	webhookvtpost "o.o/backend/pkg/integration/shipping/vtpost/webhook"
	"o.o/common/jsonx"
	"o.o/common/l"
)

func NewServers(
	_ *sqlstore.Store, // inject
	_ middleware.Middleware, // inject

	etopServer EtopServer,
	webServer WebServer,
	ghnServer GHNWebhookServer,
	ghtkServer GHTKWebhookServer,
	vtpostServer VTPostWebhookServer,
	ahamoveServer AhamoveWebhookServer,
) []lifecycle.HTTPServer {
	svrs := []lifecycle.HTTPServer{
		{"Main   ", etopServer},
		{"Web    ", webServer},
		{"GHN    ", ghnServer},
		{"GHTK   ", ghtkServer},
		{"VTPOST ", vtpostServer},
		{"AHAMOVE", ahamoveServer},
	}
	return svrs
}

type EtopServer *http.Server

func NewEtopServer(
	h EtopHandlers,
	cfg config.Config,
	eventStreamer *eventstream.EventStream,
	vtpayGatewayServer *vtpaygatewayserver.Server,
	ghnIm imcsvghn.Import,
	ghtkIm imcsvghtk.Import,
	vtpostIm vtpostimxlsx.Import,
	shopIm imports.ShopImport,
	healthService *health.Service,
) EtopServer {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService.RegisterHTTPHandler(mux)
	jsonx.RegisterHTTPHandler(mux)
	sqltrace.RegisterHTTPHandler(mux)

	// /api/
	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.StripPrefix("/api", http.NotFoundHandler()))
	mux.Handle("/api/", http.StripPrefix("/api",
		headers.ForwardHeaders(bus.Middleware(apiMux))))

	for _, s := range h.IntHandlers {
		apiMux.Handle(s.PathPrefix(), s)
	}

	// /v1/
	v1Mux := http.NewServeMux()
	v1Mux.Handle("/v1/", http.StripPrefix("/v1", http.NotFoundHandler()))
	mux.Handle("/v1/", http.StripPrefix("/v1", headers.ForwardHeaders(v1Mux)))
	for _, s := range h.ExtHandlers {
		v1Mux.Handle(s.PathPrefix(), s)
	}

	// Register import handlers
	{
		rt := httpx.New()
		mux.Handle("/api/admin.Import/", headers.ForwardHeaders(rt))
		rt.Use(httpx.RecoverAndLog(false))
		rt.Use(httpx.Auth(permission.EtopAdmin))

		rt.POST("/api/admin.Import/ghn/MoneyTransactions", ghnIm.HandleImportMoneyTransactions)
		rt.POST("/api/admin.Import/ghtk/MoneyTransactions", ghtkIm.HandleImportMoneyTransactions)
		rt.POST("/api/admin.Import/vtpost/MoneyTransactions", vtpostIm.HandleImportMoneyTransactions)

		// Hot-fix: trouble on 09/07/2019 (transfer money duplicate)
		rt.POST("/api/admin.Import/CreateMoneyTransactionShipping", hotfix.HandleImportMoneyTransactionManual)
	}

	// Register shop import handlers
	{
		mux.Handle("/api/shop.Import/", headers.ForwardHeaders(shopIm))
		shopIm.Use(httpx.RecoverAndLog(false))
	}
	// Register SSE handler
	{
		rt := httpx.New()
		mux.Handle("/api/event-stream",
			headers.ForwardHeaders(rt, headers.Config{
				AllowQueryAuthorization: true,
			}))
		rt.Use(httpx.RecoverAndLog(false))
		rt.Use(httpx.Auth(permission.Shop))
		rt.GET("/api/event-stream", eventStreamer.HandleEventStream)
	}
	{
		buildRoute := vtpaygatewayaggregate.BuildGatewayRoute
		rt := httpx.New()
		mux.Handle("/api/payment/vtpay/", rt)
		rt.Use(httpx.RecoverAndLog(false))
		rt.POST("/api"+buildRoute(vtpaygatewayaggregate.PathValidateTransaction), vtpayGatewayServer.ValidateTransaction)
		rt.POST("/api"+buildRoute(vtpaygatewayaggregate.PathGetResult), vtpayGatewayServer.GetResult)
	}
	{
		// change path for clearing browser cache and still keep the old/dl
		// path for backward compatible
		mux.Handle("/dl/imports/shop_orders.v1.xlsx",
			cmservice.ServeAssetsByContentGenerator(
				cmservice.MIMEExcel,
				orderimcsv.AssetShopOrderPath,
				5*time.Minute,
				orderimcsv.GenerateImportFile,
			),
		)
		mux.Handle("/dl/imports/shop_orders.v1b.xlsx",
			cmservice.ServeAssetsByContentGenerator(
				cmservice.MIMEExcel,
				orderimcsv.AssetShopOrderPath,
				5*time.Minute,
				orderimcsv.GenerateImportFile,
			),
		)
		mux.Handle("/dl/imports/shop_products.v1.xlsx",
			cmservice.ServeAssetsByContentGenerator(
				cmservice.MIMEExcel,
				productimcsv.AssetShopProductPath,
				5*time.Minute,
				productimcsv.GenerateImportFile,
			),
		)
		mux.Handle("/dl/imports/shop_products.v1.simplified.xlsx",
			cmservice.ServeAssets(
				productimcsv.AssetShopProductSimplifiedPath,
				cmservice.MIMEExcel,
			),
		)
	}
	// }

	if cfg.ServeDoc {
		mux.Handle("/", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		mux.Handle("/doc", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		for _, s := range strings.Split("sadmin,admin,shop,integration,affiliate,services/crm,services/affiliate", ",") {
			swaggerPath := "/doc/" + s + "/swagger.json"
			mux.Handle("/doc/"+s, cmservice.RedocHandler())
			if strings.Contains(s, "/") {
				mux.Handle(swaggerPath, cmservice.SwaggerHandler(s+"/swagger.json"))
			} else {
				mux.Handle(swaggerPath, cmservice.SwaggerHandler("etop/"+s+"/swagger.json"))
			}
		}
		mux.Handle("/doc/etop", cmservice.RedocHandler())
		mux.Handle("/doc/etop/swagger.json", cmservice.SwaggerHandler("etop/swagger.json"))

	} else {
		ll.Warn("DOCUMENTATION IS DISABLED (config.serve_doc = false)")
		mux.Handle("/doc", http.NotFoundHandler())
	}

	// always serve partner documentation
	mux.Handle("/doc/ext", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
	for _, s := range strings.Split("shop,partner,partnercarrier", ",") {
		mux.Handle("/doc/ext/"+s, cmservice.RedocHandler())
		mux.Handle("/doc/ext/"+s+"/swagger.json", cmservice.SwaggerHandler("external/"+s+"/swagger.json"))
	}

	handler := middleware.CORS(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: handler,
	}
	return svr
}

type WebServer *http.Server

func NewWebServer(cfg config.Config, webServerQuery webserverinternal.QueryBus, catalogQuery catalog.QueryBus, subscriptionQuery subscriptioncore.QueryBus, rd redis.Store, locationQueryBus location.QueryBus) WebServer {
	ecom := cfg.Ecom
	c := webserver.Config{
		MainSite: ecom.MainSite,
		CoreSite: cfg.URL.MainSite,
		RootPath: projectpath.GetPath(),
	}
	handler, err := webserver.New(c, webServerQuery, catalogQuery, rd, locationQueryBus, subscriptionQuery, nil) // TODO: order logic
	if err != nil {
		ll.S.Panicf("error starting web server: %v", err)
	}

	svr := &http.Server{
		Addr:    ecom.HTTP.Address(),
		Handler: handler,
	}
	return svr
}

type GHNWebhookServer *http.Server

func NewGHNWebhookServer(
	cfg config.Config,
	shipmentManager *shippingcarrier.ShipmentManager,
	ghnCarrier *ghn.Carrier,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhookghn.Webhook,
) GHNWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.GHNWebhook.Address(),
		Handler: rt,
	}
	return svr
}

type GHTKWebhookServer *http.Server

func NewGHTKWebhookServer(
	cfg config.Config,
	shipmentManager *shippingcarrier.ShipmentManager,
	ghtkCarrier *ghtk.Carrier,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhookghtk.Webhook,
) GHTKWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.GHTKWebhook.Address(),
		Handler: rt,
	}
	return svr
}

type VTPostWebhookServer *http.Server

func NewVTPostWebhookServer(
	cfg config.Config,
	shipmentManager *shippingcarrier.ShipmentManager,
	vtpostCarrier *vtpost.Carrier,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhookvtpost.Webhook,
) VTPostWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.VTPostWebhook.Address(),
		Handler: rt,
	}
	return svr
}

type AhamoveWebhookServer *http.Server

func NewAhamoveWebhookServer(
	cfg config.Config,
	shipmentManager *shippingcarrier.ShipmentManager,
	ahamoveCarrier *ahamove.Carrier,
	identityQuery identity.QueryBus,
	shipnowQuery shipnow.QueryBus,
	shipnowAggr shipnow.CommandBus,
	orderAggr ordering.CommandBus,
	orderQuery ordering.QueryBus,
	fileServer AhamoveVerificationFileServer,
	webhook *webhookahamove.Webhook,
) AhamoveWebhookServer {

	mux := http.NewServeMux()
	{
		rt := httpx.New()
		rt.Use(httpx.RecoverAndLog(true))
		webhook.Register(rt)

		mux.Handle("/webhook/", rt)
	}

	// serve ahamove verification files
	mux.Handle(config.PathAhamoveUserVerification+"/", (*httpx.Router)(fileServer))

	svr := &http.Server{
		Addr:    cfg.AhamoveWebhook.Address(),
		Handler: mux,
	}
	return svr
}

type AhamoveVerificationFileServer *httpx.Router

func NewAhamoveVerificationFileServer(ctx context.Context, identityQuery identity.QueryBus) AhamoveVerificationFileServer {
	// path: <UploadDirAhamoveVerification>/<originname>/<filename>.jpg
	// filepath:
	// user_id_front_<user.id>_<user.create_time>.jpg
	// user_portrait_<user.id>_<user.create_time>.jpg
	regex := regexp.MustCompile(`([0-9]+)_([0-9]+)`)

	rt := httpx.New()
	path := config.PathAhamoveUserVerification + "/:originname/:filename"

	rt.Router.GET(path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		originname := ps.ByName("originname")
		fileName := ps.ByName("filename")
		parts := regex.FindStringSubmatch(fileName)
		if len(parts) == 0 {
			http.NotFound(w, req)
			return
		}
		userID := parts[1]
		createTime := parts[2]

		query := &identity.GetExternalAccountAhamoveByExternalIDQuery{
			ExternalID: userID,
		}
		if err := identityQuery.Dispatch(ctx, query); err != nil {
			http.NotFound(w, req)
			return
		}
		accountAhamove := query.Result
		xCreatedAt := accountAhamove.ExternalCreatedAt
		if strconv.FormatInt(xCreatedAt.Unix(), 10) != createTime {
			http.NotFound(w, req)
			return
		}

		url := ""
		if strings.Contains(fileName, "user_id_front") {
			url = accountAhamove.IDCardFrontImg
		} else if strings.Contains(fileName, "user_id_back") {
			url = accountAhamove.IDCardBackImg
		} else if strings.Contains(fileName, "user_portrait") {
			url = accountAhamove.PortraitImg
		}
		if strings.Contains(url, originname) {
			http.Redirect(w, req, url, http.StatusSeeOther)
			return
		}
		http.NotFound(w, req)
	})
	return rt
}
