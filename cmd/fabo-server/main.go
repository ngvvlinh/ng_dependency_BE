package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"o.o/backend/cmd/fabo-server/config"
	servicefbmessaging "o.o/backend/com/fabo/main/fbmessaging"
	servicefbpage "o.o/backend/com/fabo/main/fbpage"
	servicefbuser "o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/fabo/pkg/fbclient"
	serviceidentity "o.o/backend/com/main/identity"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/health"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	cmwrapper "o.o/backend/pkg/common/apifw/wrapper"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/fabo"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
	"o.o/common/l"
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
	fbPageAggr := servicefbpage.NewExternalFbPageAggregate(db).MessageBus()
	fbPageQuery := servicefbpage.NewFbPageQuery(db).MessageBus()
	fbUserAggr := servicefbuser.NewFbUserAggregate(db, fbPageAggr).MessageBus()
	fbUserQuery := servicefbuser.NewFbUserQuery(db).MessageBus()
	fbMessagingAggr := servicefbmessaging.NewFbExternalMessagingAggregate(db, eventBus).MessageBus()
	fbMessagingQuery := servicefbmessaging.NewFbMessagingQuery(db).MessageBus()

	fbClient := fbclient.New(cfg.FacebookApp, bot)
	if err := fbClient.Ping(); err != nil {
		ll.Fatal("Error while connection Facebook", l.Error(err))
	}

	healthservice.MarkReady()

	mux := http.NewServeMux()
	healthservice.RegisterHTTPHandler(mux)
	svr := &http.Server{
		Addr:    cfg.HTTP.Address(),
		Handler: mux,
	}

	ss := session.New(
		session.OptValidator(tokens.NewTokenStore(redisStore)),
	)
	hooks := httprpc.ChainHooks(
		middlewares.NewLogging(),
		session.NewHook(acl.GetACL()),
	)

	var servers []httprpc.Server
	servers = append(servers, fabo.NewFaboServer(
		hooks, ss,
		fbUserQuery, fbUserAggr,
		fbPageQuery, fbPageAggr,
		fbMessagingQuery, fbMessagingAggr,
		appScopes, fbClient,
	)...)

	mux.Handle("/", http.RedirectHandler("/doc/fabo", http.StatusTemporaryRedirect))
	mux.Handle("/doc", http.RedirectHandler("/doc/fabo", http.StatusTemporaryRedirect))
	mux.Handle("/doc/fabo", cmservice.RedocHandler())
	mux.Handle("/doc/fabo/swagger.json", cmservice.SwaggerHandler("fabo/swagger.json"))

	apiMux := http.NewServeMux()
	apiMux.Handle("/api/", http.StripPrefix("/api", http.NotFoundHandler()))

	mux.Handle("/api/", http.StripPrefix("/api",
		middleware.CORS(headers.ForwardHeaders(bus.Middleware(apiMux)))))

	for _, s := range servers {
		apiMux.Handle(s.PathPrefix(), s)
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

	// Wait for OS signal or any error from services
	<-ctx.Done()
	ll.Info("Gracefully stopped!")

	// Graceful stop
	svr.Shutdown(context.Background())
	ll.Info("Gracefully stopped!")
}
