package main

import (
	"net/http"

	cmP "etop.vn/backend/pb/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/l"
	cmService "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/kiotviet/cmd/kiotviet-sync-service/config"
	"etop.vn/backend/pkg/integration/kiotviet/webhook"
	kiotvietW "etop.vn/backend/wrapper/integration/kiotviet"

	_ "etop.vn/backend/pkg/integration/kiotviet/api"
)

func startServers() []*http.Server {
	return []*http.Server{
		startServiceServer(),
		startWebhookServer(),
	}
}

func startServiceServer() *http.Server {
	mux := http.NewServeMux()
	healthservice.RegisterHTTP(mux)

	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	sqlstore.Init(db)
	// redisStore := redis.Connect(cfg.Redis.ConnectionString())

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.NotFoundHandler())
	kiotvietW.NewKiotvietServer(apiMux, nil, cfg.Secret)

	mux.Handle("/api/", middleware.ForwardHeaders(apiMux))

	if cfg.ServeDoc {
		mux.Handle("/", http.RedirectHandler("/doc/kiotviet", http.StatusTemporaryRedirect))
		mux.Handle("/doc/kiotviet", cmService.RedocHandler())
		mux.Handle("/doc/kiotviet/swagger", cmService.SwaggerHandler("integration/kiotviet/kiotviet.swagger.json"))

	} else {
		ll.Warn("DOCUMENTATION IS DISABLED (config.serve_doc = false)")
		mux.Handle("/doc", http.NotFoundHandler())
	}

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
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

func startWebhookServer() *http.Server {
	botWebhook := cfg.TelegramBot.MustConnectChannel(config.ChannelWebhook)

	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(botWebhook, true))

	webhook.Register(rt)
	rt.GET("/", func(c *httpx.Context) error {
		resp := cmP.VersionInfoResponse{
			Service: "kiotviet/webhook",
		}
		c.SetResult(resp)
		return nil
	})

	svr := &http.Server{
		Addr:    cfg.Webhook.Address(),
		Handler: rt,
	}
	ll.S.Infof("Webhook server listening at %v", cfg.Webhook.Address())

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
