package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"etop.vn/backend/cmd/fabo-server/config"
	servicefbpage "etop.vn/backend/com/fabo/main/fbpage"
	servicefbuser "etop.vn/backend/com/fabo/main/fbuser"
	"etop.vn/backend/com/fabo/pkg/fbclient"
	serviceidentity "etop.vn/backend/com/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/health"
	cmservice "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	cmwrapper "etop.vn/backend/pkg/common/apifw/wrapper"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/headers"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/authorize/middleware"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/fabo"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	ctx context.Context
	bot *telebot.Channel

	ctxCancel     context.CancelFunc
	healthservice = health.New()

	appScopes = map[string]string{
		"manage_pages":    "Quản lý các trang của bạn",
		"pages_show_list": "Hiển thị các trang do tài khoản quản lý",
		"publish_pages":   "Đăng nội dung lên trang do bạn quản lý",
		"pages_messaging": "Quản lý và truy cập các cuộc trò chuyện của trang",
		"public_profile":  "Hiển thị thông tin cơ bản của tài khoản",
	}
)

func main() {
	cc.InitFlags()
	cc.ParseFlags()

	var err error
	cfg, err = config.Load()
	if err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	cmenv.SetEnvironment(cfg.Env)
	ll.Info("Service started with config", l.String("commit", cm.CommitMessage()))
	if cmenv.IsDev() {
		ll.Info("config", l.Object("cfg", cfg))
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
		ll.Fatal("Force shutdown due to timeout!")
	}()

	if cmenv.IsDev() {
		ll.Warn("DEVELOPMENT MODE ENABLED")
	}

	wl.Init(cmenv.Env())

	bot, err = cfg.TelegramBot.ConnectDefault()
	if err != nil {
		ll.Fatal("Unable to connect to Telegram", l.Error(err))
	}

	if bot != nil {
		bot.SendMessage("–––\n✨ fabo-app started ✨\n" + cm.CommitMessage())
		defer bot.SendMessage("👹 fabo-app stopped 👹\n–––")
	}

	redisStore := redis.Connect(cfg.Redis.ConnectionString())
	tokens.Init(redisStore)
	db, err := cmsql.Connect(cfg.Postgres)
	if err != nil {
		ll.Fatal("Unable to connect to Postgres", l.Error(err))
	}
	cmwrapper.InitBot(bot)
	eventBus := bus.New()
	sqlstore.Init(db)
	sqlstore.AddEventBus(eventBus)

	_ = serviceidentity.NewQueryService(db).MessageBus()
	fbpageaggregate := servicefbpage.NewFbPageAggregate(db).MessageBus()
	fbpagequery := servicefbpage.NewFbPageQuery(db).MessageBus()
	fbuseraggregate := servicefbuser.NewFbUserAggregate(db, fbpageaggregate).MessageBus()
	fbuserquery := servicefbuser.NewFbUserQuery(db).MessageBus()

	fbClient := fbclient.New(cfg.FacebookApp, bot)
	if err := fbClient.Ping(); err != nil {
		ll.Fatal("Error while connection Facebook", l.Error(err))
	}
	middleware.NewFabo(fbpagequery, fbuserquery)

	fabo.Init(
		fbuserquery,
		fbuseraggregate,
		fbpagequery,
		fbpageaggregate,
		fbClient,
		appScopes,
	)

	healthservice.MarkReady()

	mux := http.NewServeMux()
	healthservice.RegisterHTTPHandler(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	mux.Handle("/", http.RedirectHandler("/doc/fabo", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/fabo", http.StatusTemporaryRedirect))
	mux.Handle("/doc/fabo", cmservice.RedocHandler())
	mux.Handle("/doc/fabo/swagger.json", cmservice.SwaggerHandler("fabo/swagger.json"))

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.StripPrefix("/api", http.NotFoundHandler()))

	mux.Handle("/api/", http.StripPrefix("/api",
		headers.ForwardHeaders(apiMux)))

	fabo.NewFaboServer(apiMux)

	go func() {
		defer ctxCancel()
		ll.S.Infof("HTTP server listening at %v", cfg.HTTP.Address())
		err := svr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error("HTTP server", l.Error(err))
		}
		ll.Sync()
	}()

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Gracefully stopped!")

	// Graceful stop
	svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
