package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"o.o/api/main/identity"
	webserverinternal "o.o/api/webserver"
	"o.o/backend/cmd/etop-server/config"
	paymentlogaggregate "o.o/backend/com/etc/logging/payment/aggregate"
	paymentaggregate "o.o/backend/com/external/payment/payment/aggregate"
	"o.o/backend/com/external/payment/vtpay"
	vtpaygatewayaggregate "o.o/backend/com/external/payment/vtpay/gateway/aggregate"
	vtpaygatewayserver "o.o/backend/com/external/payment/vtpay/gateway/server"
	"o.o/backend/com/web/ecom/webserver"
	"o.o/backend/pkg/common/apifw/httpx"
	cmservice "o.o/backend/pkg/common/apifw/service"
	cmwrapper "o.o/backend/pkg/common/apifw/wrapper"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/sql/sqltrace"
	api "o.o/backend/pkg/etop/api"
	admin "o.o/backend/pkg/etop/api/admin"
	affiliate "o.o/backend/pkg/etop/api/affiliate"
	crm "o.o/backend/pkg/etop/api/crm"
	integration "o.o/backend/pkg/etop/api/integration"
	sadmin "o.o/backend/pkg/etop/api/sadmin"
	shop "o.o/backend/pkg/etop/api/shop"
	partner "o.o/backend/pkg/etop/apix/partner"
	xshop "o.o/backend/pkg/etop/apix/shop"
	whitelabelapix "o.o/backend/pkg/etop/apix/whitelabel"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/logic/hotfix"
	imcsvghtk "o.o/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "o.o/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "o.o/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	orderimcsv "o.o/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "o.o/backend/pkg/etop/logic/products/imcsv"
	webhookahamove "o.o/backend/pkg/integration/shipnow/ahamove/webhook"
	webhookghn "o.o/backend/pkg/integration/shipping/ghn/webhook"
	webhookghtk "o.o/backend/pkg/integration/shipping/ghtk/webhook"
	webhookvtpost "o.o/backend/pkg/integration/shipping/vtpost/webhook"
	aff "o.o/backend/pkg/services/affiliate/api"
	"o.o/common/jsonx"
	"o.o/common/l"
)

func startServers(webServerQuery webserverinternal.QueryBus) []*http.Server {
	return []*http.Server{
		startEtopServer(),
		startWebServer(webServerQuery),
		startGHNWebhookServer(),
		startGHTKWebhookServer(),
		startVTPostWebhookServer(),
		startAhamoveWebhookServer(),
	}
}

func startEtopServer() *http.Server {
	mux := http.NewServeMux()
	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)
	jsonx.RegisterHTTPHandler(mux)
	sqltrace.RegisterHTTPHandler(mux)

	if *flDocOnly {
		ll.Warn("API IS DISABLED (-doc-only)")
		mux.Handle("/api/", http.NotFoundHandler())

	} else {
		cmwrapper.InitBot(bot)

		// /api/
		apiMux := http.NewServeMux()
		apiMux.Handle("/api/", http.StripPrefix("/api", http.NotFoundHandler()))
		mux.Handle("/api/", http.StripPrefix("/api",
			headers.ForwardHeaders(apiMux)))

		api.NewEtopServer(apiMux)
		sadmin.NewSadminServer(apiMux, ss, hooks)
		admin.NewAdminServer(apiMux)
		shop.NewShopServer(apiMux)
		affiliate.NewAffiliateServer(apiMux)
		integration.NewIntegrationServer(apiMux)
		crm.NewCrmServer(apiMux, cfg.Secret)
		aff.NewAffiliateServer(apiMux, cfg.Secret)
		// /v1/
		v1Mux := http.NewServeMux()
		v1Mux.Handle("/v1/", http.StripPrefix("/v1", http.NotFoundHandler()))
		mux.Handle("/v1/", http.StripPrefix("/v1", headers.ForwardHeaders(v1Mux)))

		partner.NewPartnerServer(v1Mux)
		xshop.NewShopServer(v1Mux)
		whitelabelapix.NewWhiteLabelServer(v1Mux)

		botDefault := cfg.TelegramBot.MustConnectChannel("")
		botImport := cfg.TelegramBot.MustConnectChannel(config.ChannelImport)

		// Register import handlers
		{
			rt := httpx.New()
			mux.Handle("/api/admin.Import/", headers.ForwardHeaders(rt))
			rt.Use(httpx.RecoverAndLog(botImport, false))
			rt.Use(httpx.Auth(permission.EtopAdmin))

			rt.POST("/api/admin.Import/ghn/MoneyTransactions", imcsvghn.HandleImportMoneyTransactions)
			rt.POST("/api/admin.Import/ghtk/MoneyTransactions", imcsvghtk.HandleImportMoneyTransactions)
			rt.POST("/api/admin.Import/vtpost/MoneyTransactions", vtpostimxlsx.HandleImportMoneyTransactions)

			// Hot-fix: trouble on 09/07/2019 (transfer money duplicate)
			_ = hotfix.New(db)
			rt.POST("/api/admin.Import/CreateMoneyTransactionShipping", hotfix.HandleImportMoneyTransactionManual)

			imcsvghtk.Init(moneyTxAggr)
			vtpostimxlsx.Init(moneyTxAggr)
			imcsvghn.Init(moneyTxAggr)
		}

		// Register shop import handlers
		{
			rt := httpx.New()
			mux.Handle("/api/shop.Import/", headers.ForwardHeaders(rt))
			rt.Use(httpx.RecoverAndLog(botImport, false))
			rt.Use(httpx.Auth(permission.Shop))
			rt.POST("/api/shop.Import/Orders", orderimcsv.HandleImportOrders)
			rt.POST("/api/shop.Import/Products", productimcsv.HandleShopImportProducts)
			rt.POST("/api/shop.Import/SampleProducts", productimcsv.HandleShopImportSampleProducts)
		}
		// Register SSE handler
		{
			rt := httpx.New()
			mux.Handle("/api/event-stream",
				headers.ForwardHeaders(rt, headers.Config{
					AllowQueryAuthorization: true,
				}))
			rt.Use(httpx.RecoverAndLog(botDefault, false))
			rt.Use(httpx.Auth(permission.Shop))
			rt.GET("/api/event-stream", eventStreamer.HandleEventStream)
		}
		{
			// Register vtpayClient gateway
			paymentAggr := paymentaggregate.NewAggregate(db).MessageBus()
			paymentLogAggr := paymentlogaggregate.NewAggregate(dbLogs)
			vtpayAggr := vtpay.NewAggregate(db, orderQuery, orderAggr.MessageBus(), paymentAggr, vtpayClient).MessageBus()
			vtpayGatewayAggr := vtpaygatewayaggregate.NewAggregate(orderQuery, orderAggr.MessageBus(), vtpayAggr, vtpayClient)

			vtpayGatewayServer := vtpaygatewayserver.New(vtpayGatewayAggr.MessageBus(), paymentLogAggr)

			buildRoute := vtpaygatewayaggregate.BuildGatewayRoute
			rt := httpx.New()
			mux.Handle("/api/payment/vtpay/", rt)
			rt.Use(httpx.RecoverAndLog(bot, false))
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
	}

	if cfg.ServeDoc || *flDocOnly {
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
	for _, s := range strings.Split("shop,partner", ",") {
		mux.Handle("/doc/ext/"+s, cmservice.RedocHandler())
		mux.Handle("/doc/ext/"+s+"/swagger.json", cmservice.SwaggerHandler("external/"+s+"/swagger.json"))
	}

	handler := middleware.CORS(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: handler,
	}
	ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func startWebServer(webServerQuery webserverinternal.QueryBus) *http.Server {
	ecom := cfg.Ecom
	c := webserver.Config{
		MainSite: ecom.MainSite,
		RootPath: projectpath.GetPath(),
	}
	handler, err := webserver.New(c, webServerQuery)
	if err != nil {
		ll.S.Panicf("error starting web server: %v", err)
	}

	svr := &http.Server{
		Addr:    ecom.HTTP.Address(),
		Handler: handler,
	}
	ll.S.Infof("Web server listening at %v", ecom.HTTP.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("Web server", l.Error(err))
		}
	}()
	return svr
}

func startGHNWebhookServer() *http.Server {
	botWebhook := cfg.TelegramBot.MustConnectChannel(config.ChannelWebhook)

	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(botWebhook, true))

	webhook := webhookghn.New(dbLogs, ghnCarrier)
	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.GHNWebhook.Address(),
		Handler: rt,
	}
	ll.S.Infof("GHN Webhook server listening at %v", cfg.GHNWebhook.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("Webhook server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func startGHTKWebhookServer() *http.Server {
	botWebhook := cfg.TelegramBot.MustConnectChannel(config.ChannelWebhook)

	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(botWebhook, true))

	webhook := webhookghtk.New(dbLogs, ghtkCarrier, shipmentManager)
	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.GHTKWebhook.Address(),
		Handler: rt,
	}
	ll.S.Infof("GHTK Webhook server listening at %v", cfg.GHTKWebhook.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("Webhook GHTK server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func startVTPostWebhookServer() *http.Server {
	botWebhook := cfg.TelegramBot.MustConnectChannel(config.ChannelWebhook)

	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(botWebhook, true))

	webhook := webhookvtpost.New(dbLogs, vtpostCarrier)
	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.VTPostWebhook.Address(),
		Handler: rt,
	}
	ll.S.Infof("VTPost Webhook server listening at %v", cfg.VTPostWebhook.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("Webhook VTPost server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func startAhamoveWebhookServer() *http.Server {
	botWebhook := cfg.TelegramBot.MustConnectChannel(config.ChannelWebhook)

	mux := http.NewServeMux()
	{
		rt := httpx.New()
		rt.Use(httpx.RecoverAndLog(botWebhook, true))
		webhook := webhookahamove.New(db, dbLogs, ahamoveCarrier, shipnowQuery, shipnowAggr, orderAggr.MessageBus(), orderQuery)
		webhook.Register(rt)

		mux.Handle("/webhook/", rt)
	}
	{
		// serve ahamove verification files
		ahamoveRouter := httpx.New()
		path := config.PathAhamoveUserVerification + "/:originname/:filename"
		serveAhamoveVerificationFiles(ahamoveRouter.Router, path, http.Dir(config.PathAhamoveUserVerification))

		mux.Handle(config.PathAhamoveUserVerification+"/", ahamoveRouter)
	}

	svr := &http.Server{
		Addr:    cfg.AhamoveWebhook.Address(),
		Handler: mux,
	}
	ll.S.Infof("Ahamove server listening at %v", cfg.AhamoveWebhook.Address())

	go func() {
		defer ctxCancel()
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("Webhook Ahamove server", l.Error(err))
		}
		ll.Sync()
	}()
	return svr
}

func tryOnDev(err error) {
	if err != nil {
		if cmenv.IsDev() {
			ll.S.Warn("DEVELOPMENT. IGNORED: ", l.Error(err))
		} else {
			ll.S.Fatal(err)
		}
	}
}

func serveAhamoveVerificationFiles(rt *httprouter.Router, path string, root http.FileSystem) {
	// path: <UploadDirAhamoveVerification>/<originname>/<filename>.jpg
	// filepath:
	// user_id_front_<user.id>_<user.create_time>.jpg
	// user_portrait_<user.id>_<user.create_time>.jpg
	regex := regexp.MustCompile(`([0-9]+)_([0-9]+)`)
	rt.GET(path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
}
