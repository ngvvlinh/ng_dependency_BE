package main

import (
	"net/http"
	"strings"
	"time"

	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/l"
	cmService "etop.vn/backend/pkg/common/service"
	cmWrapper "etop.vn/backend/pkg/common/wrapper"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/permission"
	webhookghn "etop.vn/backend/pkg/integration/ghn/webhook"
	webhookghtk "etop.vn/backend/pkg/integration/ghtk/webhook"
	webhookvtpost "etop.vn/backend/pkg/integration/vtpost/webhook"
	wrapintegration "etop.vn/backend/wrapper/etop/integration"
	wrapxpartner "etop.vn/backend/wrapper/external/partner"
	wrapxshop "etop.vn/backend/wrapper/external/shop"

	imcsvghtk "etop.vn/backend/pkg/etop/logic/money-transaction/ghtk-imcsv"
	imcsvghn "etop.vn/backend/pkg/etop/logic/money-transaction/imcsv"
	vtpostimxlsx "etop.vn/backend/pkg/etop/logic/money-transaction/vtpost-imxlsx"
	wrapetop "etop.vn/backend/wrapper/etop"
	wrapadmin "etop.vn/backend/wrapper/etop/admin"
	wrapsadmin "etop.vn/backend/wrapper/etop/sadmin"
	wrapshop "etop.vn/backend/wrapper/etop/shop"

	_ "etop.vn/backend/pkg/etop/api"
	_ "etop.vn/backend/pkg/etop/api/admin"
	_ "etop.vn/backend/pkg/etop/api/sadmin"
	_ "etop.vn/backend/pkg/etop/api/shop"
	_ "etop.vn/backend/pkg/etop/apix/partner"
	orderimcsv "etop.vn/backend/pkg/etop/logic/orders/imcsv"
	productimcsv "etop.vn/backend/pkg/etop/logic/products/imcsv"
)

func startServers() []*http.Server {
	return []*http.Server{
		startEtopServer(),
		startGHNWebhookServer(),
		startGHTKWebhookServer(),
		startVTPostWebhookServer(),
	}
}

func startEtopServer() *http.Server {
	mux := http.NewServeMux()
	healthservice.RegisterHTTP(mux)

	if *flDocOnly {
		ll.Warn("API IS DISABLED (-doc-only)")
		mux.Handle("/api/", http.NotFoundHandler())

	} else {
		cmWrapper.InitBot(bot)

		// /api/
		apiMux := http.NewServeMux()
		apiMux.Handle("/api/", http.NotFoundHandler())
		mux.Handle("/api/", middleware.ForwardHeaders(apiMux))

		wrapetop.NewEtopServer(apiMux, nil)
		wrapsadmin.NewSadminServer(apiMux, nil)
		wrapadmin.NewAdminServer(apiMux, nil)
		wrapshop.NewShopServer(apiMux, nil)
		wrapintegration.NewIntegrationServer(apiMux, nil)

		// /v1/
		v1Mux := http.NewServeMux()
		v1Mux.Handle("/v1/", http.NotFoundHandler())
		mux.Handle("/v1/", middleware.ForwardHeaders(v1Mux))

		wrapxpartner.NewPartnerServer(v1Mux, nil)
		wrapxshop.NewShopServer(v1Mux, nil)

		botDefault := cfg.TelegramBot.MustConnectChannel("")
		botImport := cfg.TelegramBot.MustConnectChannel(config.ChannelImport)

		// Register import handlers
		{
			rt := httpx.New()
			mux.Handle("/api/import/", middleware.ForwardHeaders(rt))
			rt.Use(httpx.RecoverAndLog(botImport, false))
			rt.Use(httpx.Auth(permission.EtopAdmin))

			// @deprecated
			rt.POST("/api/import/MoneyTransactions", imcsvghn.HandleImportMoneyTransactions)
		}
		{
			rt := httpx.New()
			mux.Handle("/api/admin.Import/", middleware.ForwardHeaders(rt))
			rt.Use(httpx.RecoverAndLog(botImport, false))
			rt.Use(httpx.Auth(permission.EtopAdmin))

			// @deprecated
			rt.POST("/api/admin.Import/MoneyTransactions", imcsvghn.HandleImportMoneyTransactions)

			rt.POST("/api/admin.Import/ghn/MoneyTransactions", imcsvghn.HandleImportMoneyTransactions)
			rt.POST("/api/admin.Import/ghtk/MoneyTransactions", imcsvghtk.HandleImportMoneyTransactions)
			rt.POST("/api/admin.Import/vtpost/MoneyTransactions", vtpostimxlsx.HandleImportMoneyTransactions)
		}

		// Register shop import handlers
		{
			rt := httpx.New()
			mux.Handle("/api/shop.Import/", middleware.ForwardHeaders(rt))
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
				middleware.ForwardHeaders(rt, middleware.Config{
					AllowQueryAuthorization: true,
				}))
			rt.Use(httpx.RecoverAndLog(botDefault, false))
			rt.Use(httpx.Auth(permission.Shop))
			rt.GET("/api/event-stream", eventStreamer.HandleEventStream)
		}
		{
			// change path for clearing browser cache and still keep the old
			// path for backward compatible
			mux.Handle("/dl/imports/shop_orders.v1.xlsx",
				cmService.ServeAssetsByContentGenerator(
					cmService.MIMEExcel,
					orderimcsv.AssetShopOrderPath,
					5*time.Minute,
					orderimcsv.GenerateImportFile,
				),
			)
			mux.Handle("/dl/imports/shop_orders.v1b.xlsx",
				cmService.ServeAssetsByContentGenerator(
					cmService.MIMEExcel,
					orderimcsv.AssetShopOrderPath,
					5*time.Minute,
					orderimcsv.GenerateImportFile,
				),
			)
			mux.Handle("/dl/imports/shop_products.v1.xlsx",
				cmService.ServeAssetsByContentGenerator(
					cmService.MIMEExcel,
					productimcsv.AssetShopProductPath,
					5*time.Minute,
					productimcsv.GenerateImportFile,
				),
			)
			mux.Handle("/dl/imports/shop_products.v1.simplified.xlsx",
				cmService.ServeAssets(
					productimcsv.AssetShopProductSimplifiedPath,
					cmService.MIMEExcel,
				),
			)
		}
	}

	if cfg.ServeDoc || *flDocOnly {
		mux.Handle("/", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		mux.Handle("/doc", http.RedirectHandler("/doc/etop", http.StatusTemporaryRedirect))
		for _, s := range strings.Split("sadmin,admin,shop", ",") {
			mux.Handle("/doc/"+s, cmService.RedocHandler())
			mux.Handle("/doc/"+s+"/swagger.json", cmService.SwaggerHandler("etop/"+s+"/"+s+".swagger.json"))
		}
		mux.Handle("/doc/etop", cmService.RedocHandler())
		mux.Handle("/doc/etop/swagger.json", cmService.SwaggerHandler("etop/etop.swagger.json"))

	} else {
		ll.Warn("DOCUMENTATION IS DISABLED (config.serve_doc = false)")
		mux.Handle("/doc", http.NotFoundHandler())
	}

	// always serve partner documentation
	mux.Handle("/doc/ext", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
	for _, s := range strings.Split("shop,partner", ",") {
		mux.Handle("/doc/ext/"+s, cmService.RedocHandler())
		mux.Handle("/doc/ext/"+s+"/swagger.json", cmService.SwaggerHandler("external/"+s+"/"+s+".swagger.json"))
	}

	var handler http.Handler = mux
	if cm.IsDev() {
		ll.Warn("Enabled CORS for local development")
		handler = middleware.DevelopmentCORS(mux)
	}
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

	webhook := webhookghtk.New(dbLogs, ghtkCarrier)
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

func tryOnDev(err error) {
	if err != nil {
		if cm.IsDev() {
			ll.S.Warn("DEVELOPMENT. IGNORED: ", l.Error(err))
		} else {
			ll.S.Fatal(err)
		}
	}
}
