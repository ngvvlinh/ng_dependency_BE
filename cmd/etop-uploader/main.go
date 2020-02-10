package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/cmd/etop-uploader/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/health"
	"etop.vn/backend/pkg/common/apifw/httpx"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/common/l"
)

type Purpose string

const (
	PurposeDefault             Purpose = "default"
	PurposeAhamoveVerification Purpose = "ahamove_verification"
)

type ImageConfig struct {
	Path      string
	URLPrefix string
}

var imageConfigs = map[Purpose]*ImageConfig{}

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context
	bot *telebot.Channel

	ctxCancel     context.CancelFunc
	healthservice = health.New()
	tokenStore    auth.Validator
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
	}
	imageConfigs = map[Purpose]*ImageConfig{
		PurposeDefault: {
			Path:      cfg.UploadDirImg,
			URLPrefix: cfg.URLPrefix,
		},
		PurposeAhamoveVerification: {
			Path:      cfg.UploadDirAhamoveVerification,
			URLPrefix: cfg.URLPrefixAhamoveVerification,
		},
	}

	cm.SetEnvironment(cfg.Env)
	wl.Init(cm.Env())

	_, err = os.Stat(cfg.UploadDirImg)
	if err != nil {
		ll.Fatal("Unable to open", l.String("upload_dir", cfg.UploadDirImg), l.Error(err))
	}

	_, err = os.Stat(cfg.UploadDirAhamoveVerification)
	if err != nil {
		ll.Fatal("Unable to open", l.String("upload_dir_ahamove_verification", cfg.UploadDirAhamoveVerification), l.Error(err))
	}

	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cm.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
	}

	if cm.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	ctx, ctxCancel = context.WithCancel(context.Background())
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		ll.Info("Received OS signal", l.Stringer("signal", <-osSignal))
		ctxCancel()

		// Wait for maximum 15s
		timer := time.NewTimer(15 * time.Second)
		<-timer.C
		if bot != nil {
			bot.SendMessage("👻 etop-uploader stopped (forced) 👻\n–––")
		}
		ll.Fatal("Force shutdown due to timeout!")
	}()

	bot, err = cfg.ConnectDefault()
	if err != nil {
		ll.Fatal("Connect Telegram", l.Error(err))
	}

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	tokenStore = auth.NewGenerator(redisStore)

	mux := http.NewServeMux()
	rt := httpx.New()
	mux.Handle("/", middleware.CORS(rt))

	l.RegisterHTTPHandler(mux)
	metrics.RegisterHTTPHandler(mux)
	healthservice.RegisterHTTPHandler(mux)
	healthservice.MarkReady()

	rt.Use(httpx.RecoverAndLog(bot, false))
	rt.ServeFiles("/img/*filepath", http.Dir(cfg.UploadDirImg))
	rt.ServeFiles("/ahamove/user_verification/*filepath", http.Dir(cfg.UploadDirAhamoveVerification))

	rt.POST("/upload", UploadHandler, authMiddleware)

	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}
	go func() {
		defer ctxCancel()
		ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	if bot != nil {
		bot.SendMessage("–––\n✨ etop-uploader started ✨\n" + cm.CommitMessage())
		defer bot.SendMessage("👻 etop-uploader stopped 👻\n–––")
	}

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
