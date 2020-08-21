package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"o.o/backend/cmd/uploader/config"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/common/l"
)

var (
	ll         = l.New()
	cfg        config.Config
	tokenStore auth.Validator
	bucket     storage.Bucket
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("error loading config", l.Error(err))
	}
	for _, purpose := range config.SupportedPurposes() {
		dirCfg, ok := cfg.Dirs[purpose]
		if !ok {
			ll.Fatal("no dir config", l.String("purpose", string(purpose)))
		}
		if err = dirCfg.Validate(); err != nil {
			ll.Fatal("invalid dir config", l.String("purpose", string(purpose)), l.Error(err))
		}
	}

	cmenv.SetEnvironment("uploader", cfg.Env)
	wl.Init(cmenv.Env(), wl.EtopServer)

	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximum 15s
		timer := time.NewTimer(15 * time.Second)
		<-timer.C
		ll.SendMessage("👻 etop-uploader stopped (forced) 👻\n–––")

		ll.Fatal("Force shutdown due to timeout!")
	}()
	cfg.TelegramBot.MustRegister(ctx)
	bucket, err = cfg.StorageDriver.Build(ctx)
	if err != nil {
		ll.Fatal("can not load driver", l.Error(err))
	}

	redisStore := redis.ConnectWithStr(cfg.Redis.ConnectionString())
	tokenStore = auth.NewGenerator(redisStore)

	mux := http.NewServeMux()
	rt := httpx.New()
	mux.Handle("/", middleware.CORS(rt))

	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthService := health.New(redisStore)
	healthService.RegisterHTTPHandler(mux)

	rt.Use(httpx.RecoverAndLog(false))
	rt.POST("/upload", UploadHandler)

	// TODO(vu): support serving files from driver
	fileDriver := cfg.StorageDriver.File
	if fileDriver != nil {
		for _, purpose := range config.SupportedPurposes() {
			dirCfg := cfg.Dirs[purpose]
			if cfg.StorageDriver.File != nil && dirCfg.URLPath != "" {
				dirPath := filepath.Join(fileDriver.RootPath, dirCfg.Path)
				routePath := dirCfg.URLPath + "/*filepath"
				rt.ServeFiles(routePath, http.Dir(dirPath))
			}
		}
	}

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	go func() {
		defer ctxCancel()
		ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())
		err2 := svr.ListenAndServe()
		if err2 != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err2))
		}
		ll.Sync()
	}()

	defer healthService.Shutdown()
	healthService.MarkReady()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Waiting for all requests to finish")

	// Graceful stop
	svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}

func authMiddleware(next httpx.Handler) httpx.Handler {
	return func(c *httpx.Context) error {
		tokenStr, err := auth.FromHTTPHeader(c.Req.Header)
		if err != nil {
			return cm.Error(cm.Unauthenticated, err.Error(), nil)
		}
		if tokenStr == "" {
			return cm.Error(cm.Unauthenticated, "", nil)
		}

		_, err = tokenStore.Validate(tokens.UsageAccessToken, tokenStr, nil)
		if err != nil {
			return cm.Error(cm.Unauthenticated, err.Error(), nil)
		}

		return next(c)
	}
}
